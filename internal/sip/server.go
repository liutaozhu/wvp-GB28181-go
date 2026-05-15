package sip

import (
	"fmt"
	"net"
	"sync"

	"wvp-pro-go/internal/config"

	"go.uber.org/zap"
)

// Server manages SIP UDP/TCP listeners
type Server struct {
	cfg          config.SIPConfig
	log          *zap.Logger
	tcpProviders map[string]interface{}
	udpProviders map[string]interface{}
	mu           sync.RWMutex
}

func NewServer(cfg config.SIPConfig, logger *zap.Logger) *Server {
	return &Server{
		cfg:          cfg,
		log:          logger,
		tcpProviders: make(map[string]interface{}),
		udpProviders: make(map[string]interface{}),
	}
}

func (s *Server) Start() error {
	// Determine IPs to listen on
	ips := s.getListenIPs()

	for _, ip := range ips {
		addr := fmt.Sprintf("%s:%d", ip, s.cfg.Port)

		// Start UDP listener
		if err := s.startUDP(addr, ip); err != nil {
			s.log.Error("failed to start SIP UDP listener", zap.String("addr", addr), zap.Error(err))
			return err
		}

		// Start TCP listener
		if err := s.startTCP(addr, ip); err != nil {
			s.log.Error("failed to start SIP TCP listener", zap.String("addr", addr), zap.Error(err))
			return err
		}
	}

	s.log.Info("SIP server started",
		zap.Int("port", s.cfg.Port),
		zap.String("domain", s.cfg.Domain),
		zap.String("id", s.cfg.ID),
	)
	return nil
}

func (s *Server) startUDP(addr, ip string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// In a full implementation, this would use net.ListenUDP and handle SIP messages
	// For the skeleton, we just register the IP
	s.udpProviders[ip] = addr

	s.log.Info("SIP UDP listener configured", zap.String("addr", addr))
	return nil
}

func (s *Server) startTCP(addr, ip string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// In a full implementation, this would use net.ListenTCP and handle SIP messages
	s.tcpProviders[ip] = addr

	s.log.Info("SIP TCP listener configured", zap.String("addr", addr))
	return nil
}

func (s *Server) getListenIPs() []string {
	if s.cfg.IP != "" {
		return []string{s.cfg.IP}
	}

	// Auto-detect IPs from network interfaces
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		s.log.Warn("failed to list network interfaces", zap.Error(err))
		return []string{"0.0.0.0"}
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}

	if len(ips) == 0 {
		return []string{"0.0.0.0"}
	}
	return ips
}

func (s *Server) GetProvider(ip, transport string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if transport == "TCP" {
		return s.tcpProviders[ip]
	}
	return s.udpProviders[ip]
}
