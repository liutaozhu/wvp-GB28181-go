package sip

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"wvp-pro-go/internal/config"

	"go.uber.org/zap"
)

// Sender sends SIP messages to devices
type Sender struct {
	cfg       config.SIPConfig
	subscribe *Subscribe
	log       *zap.Logger
	client    *http.Client
	cseq      int
}

func NewSender(cfg config.SIPConfig, subscribe *Subscribe, log *zap.Logger) *Sender {
	return &Sender{
		cfg:       cfg,
		subscribe: subscribe,
		log:       log,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *Sender) nextCSEQ() int {
	s.cseq++
	return s.cseq
}

// Send sends a SIP message and returns the result
func (s *Sender) Send(deviceID string, message string, timeoutMs int) (*EventResult, error) {
	// In a full implementation, this would use gosip's transport layer
	// to send the raw SIP message over UDP/TCP
	// For the skeleton, we log and simulate the send

	s.log.Debug("sending SIP message",
		zap.String("deviceID", deviceID),
		zap.Int("msgLen", len(message)),
	)

	// The actual SIP transmission would go through the gosip transport layer
	// which handles retransmissions, timeout, and response routing

	return nil, nil
}

// BuildMessage builds a SIP MESSAGE request
func (s *Sender) BuildMessage(deviceID string, body string, method string) string {
	builder := NewMessageBuilder(s.cfg)
	extraHeaders := map[string]string{
		"Subject": buildSubject(deviceID, s.cfg.ID, s.nextCSEQ()),
	}
	return builder.BuildRequest(method, deviceID, body, extraHeaders)
}

// BuildInvite builds a SIP INVITE request with SDP
func (s *Sender) BuildInvite(deviceID string, channelID string, sdpBody string) string {
	builder := NewMessageBuilder(s.cfg)
	extraHeaders := map[string]string{
		"Subject":    buildInviteSubject(channelID, s.cfg.ID, s.nextCSEQ()),
		"Contact":    fmt.Sprintf("<sip:%s@%s:%d>", s.cfg.ID, s.cfg.IP, s.cfg.Port),
		"Allow":      "INVITE,ACK,BYE,CANCEL,OPTIONS,PRACK,MESSAGE,NOTIFY,INFO,SUBSCRIBE",
		"Supported":  "timer,100rel",
	}
	return builder.BuildRequest("INVITE", channelID, sdpBody, extraHeaders)
}

// BuildBye builds a SIP BYE request
func (s *Sender) BuildBye(deviceID string, channelID string, callID string) string {
	builder := NewMessageBuilder(s.cfg)
	extraHeaders := map[string]string{
		"Call-ID": callID,
	}
	return builder.BuildRequest("BYE", channelID, "", extraHeaders)
}

// BuildAck builds a SIP ACK request
func (s *Sender) BuildAck(deviceID string, channelID string, callID string, cseq int) string {
	builder := NewMessageBuilder(s.cfg)
	extraHeaders := map[string]string{
		"Call-ID": callID,
		"CSeq":    fmt.Sprintf("%d ACK", cseq),
	}
	return builder.BuildRequest("ACK", channelID, "", extraHeaders)
}

func buildSubject(deviceID, serverID string, cseq int) string {
	return fmt.Sprintf("%s:%s,%s:%d", deviceID, serverID, serverID, cseq)
}

func buildInviteSubject(channelID, serverID string, cseq int) string {
	return fmt.Sprintf("%s:0,%s:0", channelID, serverID)
}

// GetRemoteAddr extracts remote address from SIP message
func GetRemoteAddr(msg string) string {
	lines := strings.Split(msg, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Via:") {
			// Parse Via header for address
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				addr := parts[2]
				// Remove semicolon parameters
				if idx := strings.Index(addr, ";"); idx > 0 {
					addr = addr[:idx]
				}
				return addr
			}
		}
	}
	return ""
}
