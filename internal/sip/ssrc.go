package sip

import (
	"fmt"
	"strconv"
	"strings"

	"wvp-pro-go/internal/redis"

	"go.uber.org/zap"
)

// SSRCManager manages SSRC pool and allocation
type SSRCManager struct {
	log    *zap.Logger
	prefix string // SSRC prefix from SIP domain
}

func NewSSRCManager(log *zap.Logger, domain string) *SSRCManager {
	// Use last 10 characters of domain as SSRC prefix (GB28181 rule)
	prefix := domain
	if len(domain) > 10 {
		prefix = domain[len(domain)-10:]
	}
	return &SSRCManager{
		log:    log,
		prefix: prefix,
	}
}

// GenerateSSRC generates a new SSRC
func (m *SSRCManager) GenerateSSRC() string {
	// SSRC format: domain_prefix + 6-digit serial
	return m.prefix + "000001"
}

// AllocateSSRC allocates an SSRC from the pool
func (m *SSRCManager) AllocateSSRC(streamType string) (string, error) {
	// In production, this would atomically pop from a Redis set pool
	// For now, generate based on incrementing counter
	count, _ := redis.Client.Incr(redis.Ctx, redis.CSEQKey).Result()
	ssrc := m.prefix + fmt.Sprintf("%06d", count%1000000)
	return ssrc, nil
}

// ReleaseSSRC releases an SSRC back to the pool
func (m *SSRCManager) ReleaseSSRC(ssrc string) error {
	return nil
}

// ParseSSRC parses SSRC from SDP (y= field)
func (m *SSRCManager) ParseSSRC(sdp string) string {
	lines := strings.Split(sdp, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "y=") {
			return strings.TrimPrefix(line, "y=")
		}
	}
	return ""
}

// ParseStreamInfo parses stream info from SDP f= fields
func (m *SSRCManager) ParseStreamInfo(sdp string) (map[string]string, error) {
	info := make(map[string]string)
	lines := strings.Split(sdp, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "f=") {
			val := strings.TrimPrefix(line, "f=")
			parts := strings.SplitN(val, " ", 2)
			if len(parts) == 2 {
				info[parts[0]] = parts[1]
			}
		}
	}
	return info, nil
}

// ExtractSSRCFromInvite parses SSRC from INVITE SDP response
func ExtractSSRCFromSDP(sdp string) (ssrc string, stream string) {
	lines := strings.Split(sdp, "\r\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "y="):
			ssrc = strings.TrimPrefix(line, "y=")
		case strings.HasPrefix(line, "f="):
			val := strings.TrimPrefix(line, "f=")
			parts := strings.Split(val, " ")
			if len(parts) >= 2 {
				stream = parts[1]
			}
		}
	}
	return
}

// ExtractFromSDP extracts key info from SDP
func ExtractFromSDP(sdp string) (ip string, port int, ssrc string) {
	lines := strings.Split(sdp, "\r\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "c=IN IP4"):
			ip = strings.TrimPrefix(line, "c=IN IP4 ")
		case strings.HasPrefix(line, "m=video"):
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				port, _ = strconv.Atoi(parts[2])
			}
		case strings.HasPrefix(line, "y="):
			ssrc = strings.TrimPrefix(line, "y=")
		}
	}
	return
}
