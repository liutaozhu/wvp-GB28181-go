package sip

import (
	"fmt"
	"strings"
	"time"
)

// SDPBuilder builds SDP for GB28181
type SDPBuilder struct {
	ssrc        string
	mediaServer string
	sdpIP       string
	port        int
	streamMode  string // active/passive
}

func NewSDPBuilder(sdpIP string) *SDPBuilder {
	return &SDPBuilder{
		sdpIP: sdpIP,
	}
}

// WithSSRC sets the SSRC
func (b *SDPBuilder) WithSSRC(ssrc string) *SDPBuilder {
	b.ssrc = ssrc
	return b
}

// WithMediaServer sets the media server info
func (b *SDPBuilder) WithMediaServer(ip string, rtpPort int, streamID string) *SDPBuilder {
	b.sdpIP = ip
	b.port = rtpPort
	return b
}

// BuildPlaySDP builds SDP for live play
func (b *SDPBuilder) BuildPlaySDP(channelID string) string {
	var sb strings.Builder
	sb.WriteString("v=0\r\n")
	sb.WriteString(fmt.Sprintf("o=%s 0 0 IN IP4 %s\r\n", channelID, b.sdpIP))
	sb.WriteString("s=Play\r\n")
	sb.WriteString(fmt.Sprintf("c=IN IP4 %s\r\n", b.sdpIP))
	sb.WriteString("t=0 0\r\n")
	sb.WriteString(fmt.Sprintf("m=video %d RTP/AVP 96\r\n", b.port))
	sb.WriteString("a=recvonly\r\n")
	sb.WriteString("a=rtpmap:96 PS/90000\r\n")
	sb.WriteString(fmt.Sprintf("y=%s\r\n", b.ssrc))
	sb.WriteString("f=v/2/0/0\r\n")
	return sb.String()
}

// BuildPlaybackSDP builds SDP for video playback
func (b *SDPBuilder) BuildPlaybackSDP(channelID, startTime, endTime string) string {
	var sb strings.Builder
	sb.WriteString("v=0\r\n")
	sb.WriteString(fmt.Sprintf("o=%s 0 0 IN IP4 %s\r\n", channelID, b.sdpIP))
	sb.WriteString("s=Playback\r\n")
	sb.WriteString(fmt.Sprintf("c=IN IP4 %s\r\n", b.sdpIP))
	sb.WriteString(fmt.Sprintf("t=%s %s\r\n", timeToNTP(startTime), timeToNTP(endTime)))
	sb.WriteString(fmt.Sprintf("m=video %d RTP/AVP 96\r\n", b.port))
	sb.WriteString("a=recvonly\r\n")
	sb.WriteString("a=rtpmap:96 PS/90000\r\n")
	sb.WriteString(fmt.Sprintf("y=%s\r\n", b.ssrc))
	sb.WriteString("f=v/2/0/0\r\n")
	return sb.String()
}

// BuildDownloadSDP builds SDP for video download
func (b *SDPBuilder) BuildDownloadSDP(channelID, startTime, endTime string) string {
	var sb strings.Builder
	sb.WriteString("v=0\r\n")
	sb.WriteString(fmt.Sprintf("o=%s 0 0 IN IP4 %s\r\n", channelID, b.sdpIP))
	sb.WriteString("s=Download\r\n")
	sb.WriteString(fmt.Sprintf("c=IN IP4 %s\r\n", b.sdpIP))
	sb.WriteString(fmt.Sprintf("t=%s %s\r\n", timeToNTP(startTime), timeToNTP(endTime)))
	sb.WriteString(fmt.Sprintf("m=video %d RTP/AVP 96\r\n", b.port))
	sb.WriteString("a=recvonly\r\n")
	sb.WriteString("a=rtpmap:96 PS/90000\r\n")
	sb.WriteString(fmt.Sprintf("y=%s\r\n", b.ssrc))
	sb.WriteString("f=v/2/0/0\r\n")
	return sb.String()
}

// BuildBroadcastSDP builds SDP for audio broadcast
func (b *SDPBuilder) BuildBroadcastSDP(channelID string) string {
	var sb strings.Builder
	sb.WriteString("v=0\r\n")
	sb.WriteString(fmt.Sprintf("o=%s 0 0 IN IP4 %s\r\n", channelID, b.sdpIP))
	sb.WriteString("s=Audio Broadcast\r\n")
	sb.WriteString(fmt.Sprintf("c=IN IP4 %s\r\n", b.sdpIP))
	sb.WriteString("t=0 0\r\n")
	sb.WriteString(fmt.Sprintf("m=audio %d RTP/AVP 8\r\n", b.port))
	sb.WriteString("a=sendrecv\r\n")
	sb.WriteString("a=rtpmap:8 PCMA/8000\r\n")
	sb.WriteString(fmt.Sprintf("y=%s\r\n", b.ssrc))
	return sb.String()
}

func timeToNTP(timeStr string) string {
	// Convert "2024-01-01 00:00:00" to NTP timestamp
	// NTP epoch starts from 1900-01-01
	// For simplicity, return the string as-is (GB28181 allows this format)
	if timeStr == "" {
		return "0"
	}
	// Try to parse and convert
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return "0"
	}
	ntpTime := t.Unix() + 2208988800 // NTP offset
	return fmt.Sprintf("%d", ntpTime)
}
