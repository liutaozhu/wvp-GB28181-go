package sip

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"wvp-pro-go/internal/config"
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/event"
	"wvp-pro-go/internal/model"
	sipxml "wvp-pro-go/internal/sip/xml"

	"go.uber.org/zap"
)

// ReceiveHandler handles incoming SIP messages from devices
type ReceiveHandler struct {
	cfg        config.SIPConfig
	userCfg    config.UserSettingConfig
	log        *zap.Logger
	eventBus   *event.Bus
	sessionMgr *SessionManager
	ssrcMgr    *SSRCManager
	subscribe  *Subscribe
	server     *Server
	mu         sync.Mutex
	devices    map[string]*deviceState
}

type deviceState struct {
	lastKeepalive time.Time
	lastRegister  time.Time
	ip            string
	port          int
	transport     string
}

func NewReceiveHandler(cfg config.SIPConfig, userCfg config.UserSettingConfig, log *zap.Logger,
	eventBus *event.Bus, sessionMgr *SessionManager, ssrcMgr *SSRCManager,
	subscribe *Subscribe, server *Server) *ReceiveHandler {
	return &ReceiveHandler{
		cfg:        cfg,
		userCfg:    userCfg,
		log:        log,
		eventBus:   eventBus,
		sessionMgr: sessionMgr,
		ssrcMgr:    ssrcMgr,
		subscribe:  subscribe,
		server:     server,
		devices:    make(map[string]*deviceState),
	}
}

// Start starts listening for incoming SIP messages
func (h *ReceiveHandler) Start() error {
	ips := h.getListenIPs()
	for _, ip := range ips {
		addr := fmt.Sprintf("%s:%d", ip, h.cfg.Port)

		// Start UDP listener
		if err := h.startUDP(addr); err != nil {
			h.log.Error("failed to start SIP UDP listener", zap.String("addr", addr), zap.Error(err))
			return err
		}

		// Start TCP listener
		if err := h.startTCP(addr); err != nil {
			h.log.Error("failed to start SIP TCP listener", zap.String("addr", addr), zap.Error(err))
			return err
		}
	}

	h.log.Info("SIP message receiver started",
		zap.Int("port", h.cfg.Port),
		zap.String("domain", h.cfg.Domain),
	)
	return nil
}

func (h *ReceiveHandler) startUDP(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("resolve UDP addr %s: %w", addr, err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen UDP %s: %w", addr, err)
	}

	go h.handleUDP(conn)
	return nil
}

func (h *ReceiveHandler) startTCP(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return fmt.Errorf("resolve TCP addr %s: %w", addr, err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("listen TCP %s: %w", addr, err)
	}

	go h.handleTCP(listener)
	return nil
}

func (h *ReceiveHandler) handleUDP(conn *net.UDPConn) {
	buf := make([]byte, 4096)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			h.log.Error("UDP read error", zap.Error(err))
			continue
		}

		raw := string(buf[:n])
		h.log.Debug("received SIP message",
			zap.String("from", remoteAddr.String()),
			zap.Int("size", n),
		)

		go h.processMessage(raw, remoteAddr.IP.String(), remoteAddr.Port, "UDP", conn, remoteAddr)
	}
}

func (h *ReceiveHandler) handleTCP(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			h.log.Error("TCP accept error", zap.Error(err))
			continue
		}

		go h.handleTCPConn(conn)
	}
}

func (h *ReceiveHandler) handleTCPConn(conn *net.TCPConn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				h.log.Error("TCP read error", zap.Error(err))
			}
			return
		}

		// Read full message (SIP messages end with \r\n\r\n for headers + body)
		if strings.HasSuffix(raw, "\r\n") {
			// Check if this is a complete message or needs more
			msg := raw
			if strings.Contains(msg, "Content-Length:") {
				// Parse content length and read body
				lines := strings.Split(msg, "\r\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Content-Length:") {
						cl := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
						length, _ := strconv.Atoi(cl)
						if length > 0 {
							body := make([]byte, length)
							reader.Read(body)
							msg += string(body)
						}
						break
					}
				}
			}

			remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
			go h.processMessage(msg, remoteAddr.IP.String(), remoteAddr.Port, "TCP", conn, remoteAddr)
		}
	}
}

func (h *ReceiveHandler) processMessage(raw string, ip string, port int, transport string,
	conn interface{}, addr net.Addr) {
	msg, err := ParseSIPMessage(raw)
	if err != nil {
		h.log.Warn("failed to parse SIP message", zap.Error(err), zap.String("raw", raw[:min(len(raw), 100)]))
		return
	}

	if msg.IsRequest {
		h.handleRequest(msg, ip, port, transport, conn, addr)
	} else {
		h.handleResponse(msg, ip, port, transport)
	}
}

func (h *ReceiveHandler) handleRequest(msg *Message, ip string, port int, transport string,
	conn interface{}, addr net.Addr) {
	switch msg.Method {
	case "REGISTER":
		h.handleRegister(msg, ip, port, transport, conn, addr)
	case "MESSAGE":
		h.handleMessage(msg, ip, port, transport, conn, addr)
	case "INVITE":
		h.handleInvite(msg, ip, port, transport, conn, addr)
	case "BYE":
		h.handleBye(msg, ip, port, transport)
	case "ACK":
		h.handleAck(msg, ip, port, transport)
	case "NOTIFY":
		h.handleNotify(msg, ip, port, transport, conn, addr)
	default:
		h.log.Info("unsupported SIP method", zap.String("method", msg.Method))
		h.sendResponse(conn, addr, transport, msg, 405, "Method Not Allowed", "")
	}
}

func (h *ReceiveHandler) handleResponse(msg *Message, ip string, port int, transport string) {
	switch msg.StatusCode {
	case 200:
		// 200 OK response to our INVITE/MESSAGE/etc
		cseq := msg.GetHeader("CSeq")
		parts := strings.Fields(cseq)
		if len(parts) >= 2 {
			method := parts[1]
			switch method {
			case "INVITE":
				h.handleInviteResponse(msg, ip, port, transport)
			case "MESSAGE":
				h.handleMessageResponse(msg, ip, port, transport)
			}
		}
	case 401, 407:
		// Authentication challenge
		h.log.Debug("received authentication challenge", zap.Int("status", msg.StatusCode))
	default:
		h.log.Warn("received SIP error response",
			zap.Int("status", msg.StatusCode),
			zap.String("reason", msg.Reason),
		)
	}
}

func (h *ReceiveHandler) handleRegister(msg *Message, ip string, port int, transport string,
	conn interface{}, addr net.Addr) {
	// Extract device ID from From header
	from := msg.GetHeader("From")
	deviceID := extractDeviceIDFromSIPHeader(from)
	if deviceID == "" {
		h.log.Warn("cannot extract device ID from REGISTER")
		h.sendResponse(conn, addr, transport, msg, 400, "Bad Request", "")
		return
	}

	// Check if this is a registration or unregistration
	expires := msg.GetHeader("Expires")
	isRegister := expires != "0"

	h.mu.Lock()
	h.devices[deviceID] = &deviceState{
		lastRegister: time.Now(),
		lastKeepalive: time.Now(),
		ip:        ip,
		port:      port,
		transport: transport,
	}
	h.mu.Unlock()

	if isRegister {
		h.log.Info("device registered",
			zap.String("deviceID", deviceID),
			zap.String("ip", ip),
			zap.Int("port", port),
			zap.String("transport", transport),
		)

		// Send 200 OK
		h.sendResponse(conn, addr, transport, msg, 200, "OK", "")

		// Update device in database
		h.updateDeviceRegister(deviceID, ip, port, transport)

		// Publish event
		h.eventBus.Publish(event.EventDeviceRegister, event.DeviceEvent{
			DeviceID: deviceID,
			IP:       ip,
			Port:     port,
		})

		// Query device info after registration
		go h.queryDeviceInfo(deviceID)
	} else {
		h.log.Info("device unregistered", zap.String("deviceID", deviceID))
		h.sendResponse(conn, addr, transport, msg, 200, "OK", "")
	}
}

func (h *ReceiveHandler) handleMessage(msg *Message, ip string, port int, transport string,
	conn interface{}, addr net.Addr) {
	// Parse XML body
	body := msg.Body
	if body == "" {
		h.sendResponse(conn, addr, transport, msg, 200, "OK", "")
		return
	}

	// Determine message type from XML
	cmdType, err := sipxml.GetCmdType(body)
	if err != nil {
		h.log.Warn("failed to parse XML in MESSAGE", zap.Error(err))
		h.sendResponse(conn, addr, transport, msg, 200, "OK", "")
		return
	}

	// Extract device ID from From header
	from := msg.GetHeader("From")
	deviceID := extractDeviceIDFromSIPHeader(from)

	switch cmdType {
	case "Keepalive":
		h.handleKeepalive(body, deviceID, msg, conn, addr, transport)
	case "Catalog":
		h.handleCatalog(body, deviceID, msg, conn, addr, transport)
	case "DeviceInfo":
		h.handleDeviceInfo(body, deviceID, msg, conn, addr, transport)
	case "Alarm":
		h.handleAlarm(body, deviceID, msg, conn, addr, transport)
	case "MobilePosition":
		h.handleMobilePosition(body, deviceID, msg, conn, addr, transport)
	case "MediaStatus":
		h.handleMediaStatus(body, deviceID, msg, conn, addr, transport)
	default:
		h.log.Info("unknown MESSAGE type", zap.String("cmdType", cmdType))
	}

	h.sendResponse(conn, addr, transport, msg, 200, "OK", "")
}

func (h *ReceiveHandler) handleInvite(msg *Message, ip string, port int, transport string,
	conn interface{}, addr net.Addr) {
	// Incoming INVITE (device initiating stream to us)
	// This is uncommon in GB28181 - usually server sends INVITE to device
	h.log.Info("received incoming INVITE from device", zap.String("ip", ip))

	// For now, respond with 488 Not Acceptable Here
	h.sendResponse(conn, addr, transport, msg, 488, "Not Acceptable Here", "")
}

func (h *ReceiveHandler) handleBye(msg *Message, ip string, port int, transport string) {
	// Device is hanging up a stream
	callID := msg.GetHeader("Call-ID")
	h.log.Info("received BYE from device", zap.String("callID", callID))

	// Clean up session
	session, err := h.sessionMgr.GetByCallID(callID)
	if err == nil && session != nil {
		h.sessionMgr.DeleteByCallID(callID)
		h.ssrcMgr.ReleaseSSRC(session.SSRC)

		h.log.Info("session cleaned up on BYE",
			zap.String("deviceID", session.DeviceID),
			zap.String("channelID", session.ChannelID),
		)
	}
}

func (h *ReceiveHandler) handleAck(msg *Message, ip string, port int, transport string) {
	// ACK is part of the SIP three-way handshake for INVITE
	// In a full implementation, this would confirm the session is established
	h.log.Debug("received ACK", zap.String("callID", msg.GetHeader("Call-ID")))
}

func (h *ReceiveHandler) handleNotify(msg *Message, ip string, port int, transport string,
	conn interface{}, addr net.Addr) {
	// NOTIFY is used for subscription updates (catalog changes, etc.)
	body := msg.Body
	if body != "" {
		cmdType, err := sipxml.GetCmdType(body)
		if err == nil {
			switch cmdType {
			case "Catalog":
				h.handleCatalog(body, extractDeviceIDFromSIPHeader(msg.GetHeader("From")), msg, conn, addr, transport)
			case "MobilePosition":
				h.handleMobilePosition(body, extractDeviceIDFromSIPHeader(msg.GetHeader("From")), msg, conn, addr, transport)
			}
		}
	}
	h.sendResponse(conn, addr, transport, msg, 200, "OK", "")
}

// Message type handlers

func (h *ReceiveHandler) handleKeepalive(body string, deviceID string, msg *Message,
	conn interface{}, addr net.Addr, transport string) {
	h.log.Debug("keepalive received", zap.String("deviceID", deviceID))

	h.mu.Lock()
	if state, ok := h.devices[deviceID]; ok {
		state.lastKeepalive = time.Now()
	}
	h.mu.Unlock()

	// Update keepalive time in database
	database.DB.Model(&model.Device{}).
		Where("device_id = ?", deviceID).
		Update("keepalive_time", time.Now().Format("2006-01-02 15:04:05"))
}

func (h *ReceiveHandler) handleCatalog(body string, deviceID string, msg *Message,
	conn interface{}, addr net.Addr, transport string) {
	catalog, err := sipxml.ParseCatalog(body)
	if err != nil {
		h.log.Warn("failed to parse catalog XML", zap.Error(err))
		return
	}

	h.log.Info("catalog received",
		zap.String("deviceID", deviceID),
		zap.Int("channels", len(catalog.DeviceList)),
	)

	// Save channels to database
	for _, ch := range catalog.DeviceList {
		channel := model.DeviceChannel{
			DeviceID:     deviceID,
			Name:         ch.Name,
			Manufacturer: ch.Manufacturer,
			Model:        ch.Model,
			Owner:        ch.Owner,
			CivilCode:    ch.CivilCode,
			Address:      ch.Address,
			Parental:     ch.Parental,
			ParentID:     ch.ParentID,
			SafetyWay:    ch.SafetyWay,
			RegisterWay:  ch.RegisterWay,
			CertNum:      ch.CertNum,
			Status:       ch.Status,
		}
		channel.GBDeviceID = ch.DeviceID
		channel.GBName = ch.Name
		channel.GBManufacturer = ch.Manufacturer
		channel.GBModel = ch.Model
		channel.GBOwner = ch.Owner
		channel.GBCivilCode = ch.CivilCode
		channel.GBAddress = ch.Address
		channel.GBParental = ch.Parental
		channel.GBParentID = ch.ParentID
		channel.GBSafetyWay = ch.SafetyWay
		channel.GBRegisterWay = ch.RegisterWay
		channel.GBCertNum = ch.CertNum
		channel.GBStatus = ch.Status
		channel.GBLongitude = ch.Longitude
		channel.GBLatitude = ch.Latitude

		// Upsert channel
		var existing model.DeviceChannel
		database.DB.Where("device_id = ? AND gb_device_id = ?", deviceID, ch.DeviceID).First(&existing)
		if existing.ID == 0 {
			database.DB.Create(&channel)
		}
	}

	h.log.Info("catalog saved",
		zap.String("deviceID", deviceID),
		zap.Int("channels", len(catalog.DeviceList)),
	)

	h.eventBus.Publish(event.EventCatalogUpdated, event.DeviceEvent{
		DeviceID: deviceID,
	})
}

func (h *ReceiveHandler) handleDeviceInfo(body string, deviceID string, msg *Message,
	conn interface{}, addr net.Addr, transport string) {
	info, err := sipxml.ParseDeviceInfo(body)
	if err != nil {
		h.log.Warn("failed to parse device info XML", zap.Error(err))
		return
	}

	h.log.Info("device info received", zap.String("deviceID", deviceID))

	// Update device in database
	updates := map[string]interface{}{
		"name":          info.Name,
		"manufacturer":  info.Manufacturer,
		"model":         info.Model,
		"firmware":      info.Firmware,
		"channels":      info.Channels,
	}
	database.DB.Model(&model.Device{}).Where("device_id = ?", deviceID).Updates(updates)
}

func (h *ReceiveHandler) handleAlarm(body string, deviceID string, msg *Message,
	conn interface{}, addr net.Addr, transport string) {
	alarm, err := sipxml.ParseAlarm(body)
	if err != nil {
		h.log.Warn("failed to parse alarm XML", zap.Error(err))
		return
	}

	h.log.Info("alarm received",
		zap.String("deviceID", deviceID),
		zap.String("channelID", alarm.DeviceID),
		zap.String("type", alarm.AlarmType),
	)

	// Save alarm to database
	dbAlarm := model.DeviceAlarm{
		DeviceID:       deviceID,
		ChannelID:      alarm.DeviceID,
		AlarmPriority:  alarm.AlarmPriority,
		AlarmMethod:    alarm.AlarmMethod,
		AlarmTime:      alarm.AlarmTime,
		AlarmDescription: alarm.Description,
		AlarmType:      alarm.AlarmType,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
	}
	database.DB.Create(&dbAlarm)

	h.eventBus.Publish(event.EventAlarmReceived, event.AlarmEvent{
		DeviceID:  deviceID,
		ChannelID: alarm.DeviceID,
		Type:      alarm.AlarmType,
	})
}

func (h *ReceiveHandler) handleMobilePosition(body string, deviceID string, msg *Message,
	conn interface{}, addr net.Addr, transport string) {
	pos, err := sipxml.ParseMobilePosition(body)
	if err != nil {
		h.log.Warn("failed to parse mobile position XML", zap.Error(err))
		return
	}

	h.log.Debug("mobile position received",
		zap.String("deviceID", deviceID),
		zap.Float64("lng", pos.Longitude),
		zap.Float64("lat", pos.Latitude),
	)

	// Save position to database
	dbPos := model.MobilePosition{
		DeviceID:  deviceID,
		ChannelID: pos.DeviceID,
		Longitude: pos.Longitude,
		Latitude:  pos.Latitude,
		Speed:     pos.Speed,
		Direction: pos.Direction,
		Time:      pos.Time,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	database.DB.Create(&dbPos)
}

func (h *ReceiveHandler) handleMediaStatus(body string, deviceID string, msg *Message,
	conn interface{}, addr net.Addr, transport string) {
	status, err := sipxml.ParseMediaStatus(body)
	if err != nil {
		h.log.Warn("failed to parse media status XML", zap.Error(err))
		return
	}

	h.log.Debug("media status received",
		zap.String("deviceID", deviceID),
		zap.String("notifyType", status.NotifyType),
	)

	// Handle media status notifications (close, etc.)
	if status.NotifyType == "close" {
		// Stream was closed by device
		h.log.Info("stream closed by device", zap.String("deviceID", deviceID))
	}
}

func (h *ReceiveHandler) handleInviteResponse(msg *Message, ip string, port int, transport string) {
	// 200 OK response to our INVITE contains SDP
	sdp := msg.Body
	if sdp == "" {
		return
	}

	callID := msg.GetHeader("Call-ID")
	ssrc, stream := ExtractSSRCFromSDP(sdp)

	h.log.Info("INVITE 200 OK received",
		zap.String("callID", callID),
		zap.String("ssrc", ssrc),
		zap.String("stream", stream),
	)

	// Extract IP and port from SDP
	sdpIP, _, _ := ExtractFromSDP(sdp)

	// Notify subscribe about the response
	cseq := msg.GetHeader("CSeq")
	key := callID
	if cseq != "" {
		key = cseq
	}

	event := &EventResult{
		StatusCode: 200,
		Msg:        "OK",
		Source:     sdpIP,
		Raw:        sdp,
	}

	h.subscribe.Notify(key, event)
}

func (h *ReceiveHandler) handleMessageResponse(msg *Message, ip string, port int, transport string) {
	// 200 OK response to our MESSAGE
	callID := msg.GetHeader("Call-ID")
	h.log.Debug("MESSAGE 200 OK received", zap.String("callID", callID))

	cseq := msg.GetHeader("CSeq")
	if cseq != "" {
		event := &EventResult{
			StatusCode: 200,
			Msg:        "OK",
		}
		h.subscribe.Notify(cseq, event)
	}
}

// Helper methods

func (h *ReceiveHandler) sendResponse(conn interface{}, addr net.Addr, transport string,
	req *Message, statusCode int, reason string, body string) {
	response := buildResponseMessage(req, statusCode, reason, body, h.cfg)

	switch transport {
	case "UDP":
		if udpConn, ok := conn.(*net.UDPConn); ok {
			udpAddr := &net.UDPAddr{
				IP:   addr.(*net.UDPAddr).IP,
				Port: addr.(*net.UDPAddr).Port,
			}
			udpConn.WriteTo([]byte(response), udpAddr)
		}
	case "TCP":
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.Write([]byte(response))
		}
	}
}

func (h *ReceiveHandler) updateDeviceRegister(deviceID, ip string, port int, transport string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	device := model.Device{
		DeviceID:   deviceID,
		Name:       deviceID,
		IP:         ip,
		Port:       port,
		Transport:  transport,
		Online:     true,
		CreateTime: now,
		UpdateTime: now,
	}

	var existing model.Device
	result := database.DB.Where("device_id = ?", deviceID).First(&existing)
	if result.Error != nil {
		// Device not found, create new
		database.DB.Create(&device)
	} else {
		// Update existing device
		database.DB.Model(&model.Device{}).Where("device_id = ?", deviceID).Updates(map[string]interface{}{
			"on_line":     true,
			"ip":          ip,
			"port":        port,
			"transport":   transport,
			"update_time": now,
		})
	}

	h.eventBus.Publish(event.EventDeviceOnline, event.DeviceEvent{
		DeviceID: deviceID,
		IP:       ip,
		Port:     port,
	})
}

func (h *ReceiveHandler) queryDeviceInfo(deviceID string) {
	// After registration, query device info and catalog
	time.Sleep(1 * time.Second)

	// These would use the Commander, but we avoid circular dependency
	// In production, inject Commander or use a callback
	h.log.Info("would query device info and catalog after registration",
		zap.String("deviceID", deviceID))
}

func (h *ReceiveHandler) getListenIPs() []string {
	if h.cfg.IP != "" {
		return []string{h.cfg.IP}
	}
	return []string{"0.0.0.0"}
}

func buildResponseMessage(req *Message, statusCode int, reason string, body string, cfg config.SIPConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SIP/2.0 %d %s\r\n", statusCode, reason))

	// Copy Via headers
	if vias, ok := req.Headers["Via"]; ok {
		for _, via := range vias {
			sb.WriteString(fmt.Sprintf("Via: %s\r\n", via))
		}
	}

	// Copy Call-ID
	if callID, ok := req.Headers["Call-ID"]; ok && len(callID) > 0 {
		sb.WriteString(fmt.Sprintf("Call-ID: %s\r\n", callID[0]))
	}

	// Copy From
	if from, ok := req.Headers["From"]; ok && len(from) > 0 {
		sb.WriteString(fmt.Sprintf("From: %s\r\n", from[0]))
	}

	// Copy To (add tag if not present)
	if to, ok := req.Headers["To"]; ok && len(to) > 0 {
		toVal := to[0]
		if !strings.Contains(toVal, "tag=") {
			toVal = fmt.Sprintf("%s;tag=%s", toVal, generateTag())
		}
		sb.WriteString(fmt.Sprintf("To: %s\r\n", toVal))
	}

	// Copy CSeq
	if cseq, ok := req.Headers["CSeq"]; ok && len(cseq) > 0 {
		sb.WriteString(fmt.Sprintf("CSeq: %s\r\n", cseq[0]))
	}

	if body != "" {
		sb.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	} else {
		sb.WriteString("Content-Length: 0\r\n")
	}

	sb.WriteString("\r\n")
	if body != "" {
		sb.WriteString(body)
	}

	return sb.String()
}

func generateTag() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func extractDeviceIDFromSIPHeader(header string) string {
	// From: <sip:34020000001110000001@192.168.1.100:5060>;tag=xxx
	// Extract the SIP user part (between sip: and @)
	start := strings.Index(header, "sip:")
	if start < 0 {
		return ""
	}
	start += 4
	end := strings.Index(header[start:], "@")
	if end < 0 {
		// Try closing angle bracket
		end = strings.Index(header[start:], ">")
		if end < 0 {
			return ""
		}
	}
	return header[start : start+end]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
