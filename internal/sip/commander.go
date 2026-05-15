package sip

import (
	"fmt"
	"time"

	"wvp-pro-go/internal/config"
	sipxml "wvp-pro-go/internal/sip/xml"

	"go.uber.org/zap"
)

// Commander sends SIP commands to GB28181 devices
type Commander struct {
	cfg       config.SIPConfig
	userCfg   config.UserSettingConfig
	sender    *Sender
	subscribe *Subscribe
	sdp       *SDPBuilder
	ssrc      *SSRCManager
	session   *SessionManager
	log       *zap.Logger
	cseq      int
}

func NewCommander(cfg config.SIPConfig, userCfg config.UserSettingConfig, sender *Sender,
	subscribe *Subscribe, ssrc *SSRCManager, session *SessionManager, log *zap.Logger) *Commander {
	return &Commander{
		cfg:       cfg,
		userCfg:   userCfg,
		sender:    sender,
		subscribe: subscribe,
		sdp:       NewSDPBuilder(cfg.IP),
		ssrc:      ssrc,
		session:   session,
		log:       log,
	}
}

func (c *Commander) nextCSEQ() int {
	c.cseq++
	return c.cseq
}

// PTZCmd sends PTZ control command to a device
func (c *Commander) PTZCmd(deviceID, channelID string, cmdCode byte, speed byte) error {
	sn := c.nextCSEQ()
	ptzBytes := []byte{cmdCode, speed, speed, 0}
	body, err := sipxml.MarshalPTZCmd(sn, channelID, ptzBytes)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// FrontEndCmd sends generic front-end control command
func (c *Commander) FrontEndCmd(deviceID, channelID string, cmdCode int, param1, param2, combineCode2 int) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalPTZCmd(sn, channelID, []byte{byte(cmdCode), byte(param1), byte(param2), byte(combineCode2)})
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// PlayStreamCmd sends INVITE for live streaming
func (c *Commander) PlayStreamCmd(deviceID, channelID, ssrc string, sdpIP string, port int) (*EventResult, error) {
	sn := c.nextCSEQ()
	sdpBody := NewSDPBuilder(sdpIP).WithSSRC(ssrc).WithMediaServer(sdpIP, port, "").BuildPlaySDP(channelID)

	msg := c.sender.BuildInvite(deviceID, channelID, sdpBody)
	key := fmt.Sprintf("%s:%d:%s", deviceID, sn, channelID)
	event := c.subscribe.AddSubscribe(key, time.Duration(c.userCfg.PlayTimeout)*time.Millisecond)

	go func() {
		result, err := c.sender.Send(deviceID, msg, 3000)
		if err != nil {
			c.subscribe.Notify(key, &EventResult{StatusCode: 500, Msg: err.Error()})
			return
		}
		if result != nil && result.StatusCode >= 200 {
			c.subscribe.Notify(key, result)
		}
	}()

	return c.subscribe.WaitForResult(event)
}

// PlaybackStreamCmd sends INVITE for video playback
func (c *Commander) PlaybackStreamCmd(deviceID, channelID, startTime, endTime, ssrc string, sdpIP string, port int) (*EventResult, error) {
	sn := c.nextCSEQ()
	sdpBody := NewSDPBuilder(sdpIP).WithSSRC(ssrc).WithMediaServer(sdpIP, port, "").BuildPlaybackSDP(channelID, startTime, endTime)

	msg := c.sender.BuildInvite(deviceID, channelID, sdpBody)
	key := fmt.Sprintf("%s:%d:%s:playback", deviceID, sn, channelID)
	event := c.subscribe.AddSubscribe(key, time.Duration(c.userCfg.PlayTimeout)*time.Millisecond)

	go func() {
		result, err := c.sender.Send(deviceID, msg, 3000)
		if err != nil {
			c.subscribe.Notify(key, &EventResult{StatusCode: 500, Msg: err.Error()})
			return
		}
		if result != nil {
			c.subscribe.Notify(key, result)
		}
	}()

	return c.subscribe.WaitForResult(event)
}

// DownloadStreamCmd sends INVITE for video download
func (c *Commander) DownloadStreamCmd(deviceID, channelID, startTime, endTime, ssrc string, sdpIP string, port int) (*EventResult, error) {
	sn := c.nextCSEQ()
	sdpBody := NewSDPBuilder(sdpIP).WithSSRC(ssrc).WithMediaServer(sdpIP, port, "").BuildDownloadSDP(channelID, startTime, endTime)

	msg := c.sender.BuildInvite(deviceID, channelID, sdpBody)
	key := fmt.Sprintf("%s:%d:%s:download", deviceID, sn, channelID)
	event := c.subscribe.AddSubscribe(key, time.Duration(c.userCfg.PlayTimeout)*time.Millisecond)

	go func() {
		result, err := c.sender.Send(deviceID, msg, 3000)
		if err != nil {
			c.subscribe.Notify(key, &EventResult{StatusCode: 500, Msg: err.Error()})
			return
		}
		if result != nil {
			c.subscribe.Notify(key, result)
		}
	}()

	return c.subscribe.WaitForResult(event)
}

// StreamByeCmd sends BYE to stop a stream
func (c *Commander) StreamByeCmd(deviceID, channelID, callID string) error {
	msg := c.sender.BuildBye(deviceID, channelID, callID)
	_, err := c.sender.Send(deviceID, msg, 3000)
	return err
}

// TalkStreamCmd sends INVITE for audio talk
func (c *Commander) TalkStreamCmd(deviceID, channelID, ssrc string, sdpIP string, port int) (*EventResult, error) {
	sn := c.nextCSEQ()
	sdpBody := NewSDPBuilder(sdpIP).WithSSRC(ssrc).WithMediaServer(sdpIP, port, "").BuildBroadcastSDP(channelID)

	msg := c.sender.BuildInvite(deviceID, channelID, sdpBody)
	key := fmt.Sprintf("%s:%d:%s:talk", deviceID, sn, channelID)
	event := c.subscribe.AddSubscribe(key, 10*time.Second)

	go func() {
		result, err := c.sender.Send(deviceID, msg, 3000)
		if err != nil {
			c.subscribe.Notify(key, &EventResult{StatusCode: 500, Msg: err.Error()})
			return
		}
		if result != nil {
			c.subscribe.Notify(key, result)
		}
	}()

	return c.subscribe.WaitForResult(event)
}

// AudioBroadcastCmd sends audio broadcast command
func (c *Commander) AudioBroadcastCmd(deviceID, channelID string) error {
	_ = c.nextCSEQ()
	// Audio broadcast uses MESSAGE with specific XML
	msg := c.sender.BuildMessage(deviceID, "", "MESSAGE")
	_, err := c.sender.Send(deviceID, msg, 3000)
	return err
}

// CatalogQuery sends catalog query to device
func (c *Commander) CatalogQuery(deviceID string) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalCatalogQuery(sn, deviceID)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// RecordInfoQuery sends record info query
func (c *Commander) RecordInfoQuery(deviceID, channelID, startTime, endTime string) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalRecordInfoQuery(sn, channelID, startTime, endTime)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// DeviceInfoQuery sends device info query
func (c *Commander) DeviceInfoQuery(deviceID string) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalDeviceInfoQuery(sn, deviceID)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// PresetQuery sends preset position query
func (c *Commander) PresetQuery(deviceID, channelID string) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalPresetQuery(sn, channelID)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// KeepaliveQuery sends keepalive query
func (c *Commander) KeepaliveQuery(deviceID string) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalKeepalive(sn, deviceID)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// MobilePositionSubscribe subscribes to mobile position updates
func (c *Commander) MobilePositionSubscribe(deviceID, channelID string, interval int) error {
	sn := c.nextCSEQ()
	body, err := sipxml.MarshalMobilePositionSubscribe(sn, channelID, interval)
	if err != nil {
		return err
	}

	msg := c.sender.BuildMessage(deviceID, string(body), "MESSAGE")
	_, err = c.sender.Send(deviceID, msg, 3000)
	return err
}

// PlaybackControl sends playback control (pause, resume, seek, speed)
func (c *Commander) PlaybackControl(deviceID, channelID, controlType string, param string) error {
	sn := c.nextCSEQ()
	// Build INFO message for playback control
	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Control><CmdType>DeviceControl</CmdType><SN>%d</SN><DeviceId>%s</DeviceId><%s>%s</%s></Control>`,
		sn, channelID, controlType, param, controlType)

	msg := c.sender.BuildMessage(deviceID, body, "INFO")
	_, err := c.sender.Send(deviceID, msg, 3000)
	return err
}
