package sip

import (
	"crypto/rand"
	"fmt"
	"net"
	"strconv"
	"strings"

	"wvp-pro-go/internal/config"

	"github.com/google/uuid"
)

// Message represents a SIP message
type Message struct {
	Method     string
	URI        string
	Headers    map[string][]string
	Body       string
	IsRequest  bool
	StatusCode int
	Reason     string
}

// MessageBuilder constructs SIP messages
type MessageBuilder struct {
	cfg    config.SIPConfig
	cseq   int
}

func NewMessageBuilder(cfg config.SIPConfig) *MessageBuilder {
	return &MessageBuilder{cfg: cfg}
}

// BuildRequest builds a SIP request message
func (b *MessageBuilder) BuildRequest(method, recipient string, body string, extraHeaders map[string]string) string {
	callID := b.generateCallID()
	cseq := b.nextCSEQ()
	via := fmt.Sprintf("SIP/2.0/UDP %s:%d;rport;branch=%s",
		b.cfg.IP, b.cfg.Port, b.generateBranch())
	from := fmt.Sprintf("<sip:%s@%s>;tag=%s", b.cfg.ID, b.cfg.Domain, b.generateTag())
	to := fmt.Sprintf("<sip:%s>", recipient)
	maxForwards := "70"

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s sip:%s SIP/2.0\r\n", method, recipient))
	sb.WriteString(fmt.Sprintf("Via: %s\r\n", via))
	sb.WriteString(fmt.Sprintf("Call-ID: %s\r\n", callID))
	sb.WriteString(fmt.Sprintf("CSeq: %d %s\r\n", cseq, method))
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	sb.WriteString(fmt.Sprintf("Max-Forwards: %s\r\n", maxForwards))
	sb.WriteString(fmt.Sprintf("User-Agent: WVP-PRO-GO/2.7.4\r\n"))

	// Content type based on method
	switch method {
	case "MESSAGE":
		sb.WriteString("Content-Type: Application/MANSCDP+xml\r\n")
	case "INVITE":
		sb.WriteString("Content-Type: Application/SDP\r\n")
	}

	if body != "" {
		sb.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	} else {
		sb.WriteString("Content-Length: 0\r\n")
	}

	// Extra headers
	for k, v := range extraHeaders {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	sb.WriteString("\r\n")
	if body != "" {
		sb.WriteString(body)
	}

	return sb.String()
}

// BuildResponse builds a SIP response message
func (b *MessageBuilder) BuildResponse(requestMsg string, statusCode int, reason string, body string) string {
	// Parse the request to extract Via, Call-ID, etc.
	lines := strings.Split(requestMsg, "\r\n")
	callID := ""
	via := ""
	from := ""
	to := ""

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Call-ID:"):
			callID = strings.TrimSpace(strings.TrimPrefix(line, "Call-ID:"))
		case strings.HasPrefix(line, "Via:"):
			via = strings.TrimSpace(strings.TrimPrefix(line, "Via:"))
		case strings.HasPrefix(line, "From:"):
			from = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		case strings.HasPrefix(line, "To:"):
			to = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SIP/2.0 %d %s\r\n", statusCode, reason))
	if via != "" {
		sb.WriteString(fmt.Sprintf("Via: %s\r\n", via))
	}
	if callID != "" {
		sb.WriteString(fmt.Sprintf("Call-ID: %s\r\n", callID))
	}
	if from != "" {
		sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	}
	if to != "" {
		// Add tag to To header in response
		if !strings.Contains(to, "tag=") {
			to = fmt.Sprintf("%s;tag=%s", to, b.generateTag())
		}
		sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	}

	cseq := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "CSeq:") {
			cseq = strings.TrimSpace(strings.TrimPrefix(line, "CSeq:"))
			break
		}
	}
	if cseq != "" {
		sb.WriteString(fmt.Sprintf("CSeq: %s\r\n", cseq))
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

func (b *MessageBuilder) nextCSEQ() int {
	b.cseq++
	return b.cseq
}

func (b *MessageBuilder) generateCallID() string {
	return uuid.New().String()
}

func (b *MessageBuilder) generateBranch() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return fmt.Sprintf("z9hG4bK%x", buf)
}

func (b *MessageBuilder) generateTag() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return fmt.Sprintf("%x", buf)
}

// GetLocalIP returns the best local IP for SIP
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// ParseSIPMessage parses a raw SIP message
func ParseSIPMessage(raw string) (*Message, error) {
	lines := strings.Split(raw, "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty SIP message")
	}

	msg := &Message{
		Headers: make(map[string][]string),
	}

	// Parse start line
	firstLine := lines[0]
	if strings.HasPrefix(firstLine, "SIP/2.0") {
		// Response
		msg.IsRequest = false
		parts := strings.SplitN(firstLine, " ", 3)
		if len(parts) >= 2 {
			msg.StatusCode, _ = strconv.Atoi(parts[1])
		}
		if len(parts) >= 3 {
			msg.Reason = parts[2]
		}
	} else {
		// Request
		msg.IsRequest = true
		parts := strings.SplitN(firstLine, " ", 3)
		if len(parts) >= 2 {
			msg.Method = parts[0]
			msg.URI = parts[1]
		}
	}

	// Parse headers
	bodyStart := false
	var bodyLines []string
	for _, line := range lines[1:] {
		if line == "" && !bodyStart {
			bodyStart = true
			continue
		}
		if bodyStart {
			bodyLines = append(bodyLines, line)
			continue
		}

		if colonIdx := strings.Index(line, ":"); colonIdx > 0 {
			key := strings.TrimSpace(line[:colonIdx])
			value := strings.TrimSpace(line[colonIdx+1:])
			msg.Headers[key] = append(msg.Headers[key], value)
		}
	}

	msg.Body = strings.Join(bodyLines, "\r\n")
	return msg, nil
}

// GetHeader returns the first value for a header
func (m *Message) GetHeader(name string) string {
	if values, ok := m.Headers[name]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}
