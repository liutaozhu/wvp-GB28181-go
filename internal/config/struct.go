package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	GlobalConfig *Config
	Logger       *zap.Logger
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	SIP         SIPConfig         `mapstructure:"sip"`
	Media       MediaConfig       `mapstructure:"media"`
	UserSetting UserSettingConfig `mapstructure:"user-settings"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Log         LogConfig         `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`   // mysql, postgres
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Charset  string `mapstructure:"charset"`
}

func (d *DatabaseConfig) DSN() string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.DBName, d.Charset)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.Port, d.Username, d.Password, d.DBName)
	default:
		return ""
	}
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type SIPConfig struct {
	IP       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Domain   string `mapstructure:"domain"`
	ID       string `mapstructure:"id"`
	Password string `mapstructure:"password"`
	Alarm    bool   `mapstructure:"alarm"`
}

type MediaConfig struct {
	ID        string `mapstructure:"id"`
	IP        string `mapstructure:"ip"`
	WanIP     string `mapstructure:"wan-ip"`
	HookIP    string `mapstructure:"hook-ip"`
	HTTPPort  int    `mapstructure:"http-port"`
	Secret    string `mapstructure:"secret"`
	RTP       RTPConfig `mapstructure:"rtp"`
	Type      string `mapstructure:"type"` // zlm, abl
}

type RTPConfig struct {
	Enable        bool   `mapstructure:"enable"`
	PortRange     string `mapstructure:"port-range"`
	SendPortRange string `mapstructure:"send-port-range"`
}

type UserSettingConfig struct {
	PlayTimeout            int  `mapstructure:"play-timeout"`
	AutoApplyPlay          bool `mapstructure:"auto-apply-play"`
	RecordPushLive         bool `mapstructure:"record-push-live"`
	RecordSIP              bool `mapstructure:"record-sip"`
	StreamOnDemand         bool `mapstructure:"stream-on-demand"`
	InterfaceAuthentication bool `mapstructure:"interface-authentication"`
	UseSourceIpAsStreamIP  bool `mapstructure:"use-source-ip-as-stream-ip"`
}

type JWTConfig struct {
	Secret   string `mapstructure:"secret"`
	Expire   int    `mapstructure:"expire"`   // hours
	Issuer   string `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file-path"`
	MaxSize    int    `mapstructure:"max-size"`    // MB
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`      // days
	Compress   bool   `mapstructure:"compress"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// Environment variable support
	v.SetEnvPrefix("WVP")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 18080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "release"
	}
	if cfg.Database.Charset == "" {
		cfg.Database.Charset = "utf8mb4"
	}
	if cfg.SIP.Port == 0 {
		cfg.SIP.Port = 8116
	}
	if cfg.SIP.Domain == "" {
		cfg.SIP.Domain = "4101050000"
	}
	if cfg.SIP.ID == "" {
		cfg.SIP.ID = "41010500002000000001"
	}
	if cfg.Media.HookIP == "" {
		cfg.Media.HookIP = "127.0.0.1"
	}
	if cfg.UserSetting.PlayTimeout == 0 {
		cfg.UserSetting.PlayTimeout = 180000
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

func InitLogger(cfg LogConfig) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	if cfg.FilePath != "" {
		zapCfg.OutputPaths = []string{cfg.FilePath, "stdout"}
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	Logger = logger
	return logger, nil
}

// GetConfigDir returns the directory of the config file
func GetConfigDir() string {
	if GlobalConfig == nil {
		return "."
	}
	if path := viper.ConfigFileUsed(); path != "" {
		dir, _ := os.Stat(path)
		if dir != nil && !dir.IsDir() {
			return ""
		}
	}
	return "."
}
