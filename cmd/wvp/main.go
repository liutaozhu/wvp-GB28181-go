package main

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"wvp-pro-go/internal/config"
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/depcheck"
	"wvp-pro-go/internal/event"
	"wvp-pro-go/internal/handler"
	"wvp-pro-go/internal/middleware"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/redis"
	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/sip"
	"wvp-pro-go/internal/task"
	"wvp-pro-go/internal/zlm"
	"wvp-pro-go/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	configPath := ""
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := config.InitLogger(cfg.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("WVP-PRO-GO starting", zap.String("version", "2.7.4"))

	// --- Dependency Check & Auto-Start ---
	depChecker := depcheck.New(logger)
	mysqlAddr := fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port)
	redisAddr := cfg.Redis.Addr()
	zlmAddr := fmt.Sprintf("%s:%d", cfg.Media.IP, cfg.Media.HTTPPort)

	results := depChecker.CheckAll(mysqlAddr, redisAddr, zlmAddr)

	// Check for missing dependencies
	var missing []depcheck.CheckResult
	var notRunning []depcheck.CheckResult
	for _, r := range results {
		if !r.Installed {
			missing = append(missing, r)
		} else if !r.Running {
			notRunning = append(notRunning, r)
		}
	}

	if len(missing) > 0 {
		depcheck.PrintInstallGuide(results)
		fmt.Println(depcheck.QuickDockerInstall())
		fmt.Println()
		logger.Fatal("Missing dependencies, please install them before running WVP-PRO-GO")
		return
	}

	if len(notRunning) > 0 {
		fmt.Println()
		fmt.Println("=====================================================")
		fmt.Println("  Dependencies not running, attempting auto-start...")
		fmt.Println("=====================================================")
		fmt.Println()

		for _, r := range notRunning {
			logger.Info("starting " + r.Name + "...")
			if err := depChecker.TryStart(r); err != nil {
				logger.Warn(r.Name+" auto-start failed", zap.Error(err))
				fmt.Printf("  %s: auto-start failed, please start manually:\n", r.Name)
				switch r.Name {
				case "MySQL":
					fmt.Printf("    %s\n", depcheck.ServiceCommand("mysql"))
				case "Redis":
					fmt.Printf("    %s\n", depcheck.ServiceCommand("redis"))
				case "ZLMediaKit":
					fmt.Printf("    docker start zlmediakit\n")
				}
			} else {
				logger.Info(r.Name + " started successfully")
				// Wait a moment for service to be ready
				time.Sleep(2 * time.Second)
			}
		}
		fmt.Println()
	}

	// --- Initialize Database ---
	if err := database.Init(cfg.Database, logger); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
		return
	}

	if err := database.AutoMigrate(
		&model.Device{},
		&model.DeviceChannel{},
		&model.Platform{},
		&model.PlatformChannel{},
		&model.Group{},
		&model.Region{},
		&model.StreamProxy{},
		&model.MediaServer{},
		&model.MobilePosition{},
		&model.DeviceAlarm{},
		&model.JTTerminal{},
		&model.JTChannel{},
		&model.User{},
		&model.Role{},
		&model.Record{},
	); err != nil {
		logger.Fatal("Failed to auto-migrate", zap.Error(err))
		return
	}
	logger.Info("database schema migrated")

	// --- Create default admin user if not exists ---
	if err := initDefaultAdmin(logger); err != nil {
		logger.Fatal("Failed to init default admin user", zap.Error(err))
		return
	}

	// --- Initialize Redis ---
	if err := redis.Init(cfg.Redis, logger); err != nil {
		logger.Fatal("Failed to init redis", zap.Error(err))
		return
	}
	defer redis.Client.Close()

	// --- Initialize Event Bus ---
	eventBus := event.NewBus()

	// --- Initialize SIP Server ---
	sipServer := sip.NewServer(cfg.SIP, logger)
	if err := sipServer.Start(); err != nil {
		logger.Warn("SIP server failed to start (non-critical for API-only mode)", zap.Error(err))
	}

	subscribe := sip.NewSubscribe(logger, cfg.UserSetting)
	ssrcManager := sip.NewSSRCManager(logger, cfg.SIP.Domain)
	sessionManager := sip.NewSessionManager(logger)
	sender := sip.NewSender(cfg.SIP, subscribe, logger)
	commander := sip.NewCommander(cfg.SIP, cfg.UserSetting, sender, subscribe, ssrcManager, sessionManager, logger)

	receiveHandler := sip.NewReceiveHandler(cfg.SIP, cfg.UserSetting, logger, eventBus,
		sessionManager, ssrcManager, subscribe, sipServer)
	if err := receiveHandler.Start(); err != nil {
		logger.Warn("SIP receive handler failed to start (non-critical for API-only mode)", zap.Error(err))
	}

	// --- Initialize ZLMediaKit Client ---
	zlmClient := zlm.NewClient(cfg.Media.IP, cfg.Media.HTTPPort, cfg.Media.Secret)
	zlmServer := zlm.NewServer(logger, zlmClient)

	// --- Ensure default MediaServer record exists in DB ---
	if err := initDefaultMediaServer(cfg, logger); err != nil {
		logger.Warn("Failed to init default media server", zap.Error(err))
	}

	// --- Initialize Services ---
	svcs := service.InitServices(commander, zlmClient, zlmServer,
		ssrcManager, sessionManager, subscribe, eventBus, logger)

	// --- Initialize Task Scheduler ---
	scheduler := task.NewScheduler(logger)
	scheduler.Start()
	defer scheduler.Stop()

	// --- Setup HTTP Server ---
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(logger))

	hookHandler := zlm.NewHookHandler(logger, cfg.Media, eventBus, sessionManager, ssrcManager, commander)
	hookGroup := r.Group("/index/hook")
	hookHandler.Register(hookGroup)

	// Public API routes (no auth required)
	public := r.Group("/api")
	{
		authHandler := handler.NewAuthHandler(svcs)
		public.GET("/user/login", authHandler.Login)
		public.GET("/user/logout", authHandler.Logout)
	}

	// Protected API routes (JWT auth required)
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		authHandler := handler.NewAuthHandler(svcs)
		api.POST("/user/userInfo", authHandler.GetUserInfo)
		api.GET("/user/users", authHandler.QueryUsers)
		api.POST("/user/add", authHandler.AddUser)
		api.DELETE("/user/delete", authHandler.DeleteUser)
		api.POST("/user/changePassword", authHandler.ChangePassword)

		handler.NewChannelHandler(svcs).Register(api.Group("/common/channel"))
		handler.NewPlayHandler(svcs).Register(api.Group("/play"))
		handler.NewPtzHandler(svcs).Register(api.Group("/front-end"))
		handler.NewStreamProxyHandler(svcs).Register(api.Group("/proxy"))
		handler.NewPlatformHandler(svcs).Register(api.Group("/platform"))
		handler.NewRegionHandler(svcs).Register(api.Group("/region"))
		handler.NewGroupHandler(svcs).Register(api.Group("/group"))
		handler.NewServerHandler(svcs).Register(api.Group("/server"))
	}

	// --- Serve embedded frontend static files ---
	staticFS, _ := fs.Sub(web.StaticFS, "dist")
	indexHTML, _ := fs.ReadFile(web.StaticFS, "dist/index.html")
	fileServer := http.FileServer(http.FS(staticFS))
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API and hook paths (should not reach here, but just in case)
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/index/hook") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		// Try to serve the static file
		filePath := strings.TrimPrefix(path, "/")
		if filePath != "" {
			if f, err := staticFS.Open(filePath); err == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// Root path or SPA fallback: serve index.html directly
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	_ = eventBus
	_ = sipServer
	_ = zlmServer

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	localIP := getLocalIP()
	accessURL := fmt.Sprintf("http://%s:%d", localIP, cfg.Server.Port)
	logger.Info("WVP-PRO-GO is ready!", zap.String("addr", addr))
	fmt.Println()
	fmt.Println("=====================================================")
	fmt.Printf("  WVP-PRO-GO started successfully!\n")
	fmt.Printf("  Browser access URL: %s\n", accessURL)
	fmt.Printf("  Local access URL:   http://localhost:%d\n", cfg.Server.Port)
	fmt.Println("=====================================================")

	// --- Unified lifecycle management ---
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Check if port is available
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Port is already in use, cannot start HTTP server", zap.String("addr", addr), zap.Error(err))
		fmt.Fprintf(os.Stderr, "\nError: port %d is already in use. Please stop the other process first.\n", cfg.Server.Port)
		os.Exit(1)
	}
	ln.Close()

	// Start HTTP server in goroutine
	srv := &http.Server{Addr: addr, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down...")
	cancel()

	// Graceful shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("HTTP server shutdown error", zap.Error(err))
	}

	logger.Info("all services stopped")
}

// initDefaultAdmin creates a default admin user (admin/admin123) if no users exist
func initDefaultAdmin(logger *zap.Logger) error {
	var count int64
	database.DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		logger.Info("users already exist, skip default admin creation")
		return nil
	}

	// MD5 hash of "admin123"
	h := model.User{
		Username:     "admin",
		Password:     "0192023a7bbd73250516f069df18b500", // md5("admin123")
		Name:         "Administrator",
		Enable:       true,
		CreateTime:   fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
	if err := database.DB.Create(&h).Error; err != nil {
		return err
	}

	logger.Info("default admin user created", zap.String("username", "admin"), zap.String("password", "admin123"))
	return nil
}

// initDefaultMediaServer ensures the default media server record exists in DB
func initDefaultMediaServer(cfg *config.Config, logger *zap.Logger) error {
	var ms model.MediaServer
	err := database.DB.Where("id = ?", cfg.Media.ID).First(&ms).Error

	streamIP := cfg.Media.IP
	if cfg.Media.WanIP != "" {
		streamIP = cfg.Media.WanIP
	}
	// For local docker deployment, use the host IP for stream access
	if streamIP == "127.0.0.1" || streamIP == "localhost" {
		// Try to get local network IP
		streamIP = getLocalIP()
	}

	if err != nil {
		// Create new record
		ms = model.MediaServer{
			ID:           cfg.Media.ID,
			IP:           cfg.Media.IP,
			HookIP:       cfg.Media.HookIP,
			SDPIP:        streamIP,
			StreamIP:     streamIP,
			HTTPPort:     cfg.Media.HTTPPort,
			RTSPPort:     8554,
			RTMPPort:     1935,
			FLVPort:      cfg.Media.HTTPPort,
			RTPProxyPort: 10000,
			Secret:       cfg.Media.Secret,
			RTPEnable:    cfg.Media.RTP.Enable,
			RTPPortRange: cfg.Media.RTP.PortRange,
			SendRTPPortRange: cfg.Media.RTP.SendPortRange,
			Status:       true,
			DefaultServer: true,
			Type:         cfg.Media.Type,
			CreateTime:   fmt.Sprintf("%d", time.Now().UnixMilli()),
			UpdateTime:   fmt.Sprintf("%d", time.Now().UnixMilli()),
		}
		if err := database.DB.Create(&ms).Error; err != nil {
			return err
		}
		logger.Info("default media server created", zap.String("id", ms.ID), zap.String("streamIP", streamIP))
	} else {
		// Update existing record with current config
		ms.IP = cfg.Media.IP
		ms.HookIP = cfg.Media.HookIP
		ms.StreamIP = streamIP
		ms.SDPIP = streamIP
		ms.HTTPPort = cfg.Media.HTTPPort
		ms.RTSPPort = 8554
		ms.RTMPPort = 1935
		ms.FLVPort = cfg.Media.HTTPPort
		ms.RTPProxyPort = 10000
		ms.Secret = cfg.Media.Secret
		ms.RTPEnable = cfg.Media.RTP.Enable
		ms.Status = true
		ms.UpdateTime = fmt.Sprintf("%d", time.Now().UnixMilli())
		if err := database.DB.Save(&ms).Error; err != nil {
			return err
		}
		logger.Info("default media server updated", zap.String("id", ms.ID), zap.String("streamIP", streamIP))
	}
	return nil
}

// getLocalIP returns the first non-loopback IPv4 address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
