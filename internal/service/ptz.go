package service

import (
	"wvp-pro-go/internal/sip"

	"go.uber.org/zap"
)

type PTZService struct {
	cmd *sip.Commander
	log *zap.Logger
}

func NewPTZService(cmd *sip.Commander, log *zap.Logger) *PTZService {
	return &PTZService{cmd: cmd, log: log}
}

func (s *PTZService) PTZControl(deviceID, channelID string, cmdCode, speed int) error {
	s.log.Info("PTZ control", zap.String("deviceID", deviceID), zap.String("channelID", channelID),
		zap.Int("cmdCode", cmdCode), zap.Int("speed", speed))
	return s.cmd.PTZCmd(deviceID, channelID, byte(cmdCode), byte(speed))
}

func (s *PTZService) FrontEndCommand(deviceID, channelID string, cmdCode, parameter1, parameter2, combineCode2 int) error {
	s.log.Info("Front-end command", zap.String("deviceID", deviceID), zap.String("channelID", channelID),
		zap.Int("cmdCode", cmdCode))
	return s.cmd.FrontEndCmd(deviceID, channelID, cmdCode, parameter1, parameter2, combineCode2)
}

func (s *PTZService) QueryPresets(deviceID, channelID string) error {
	return s.cmd.PresetQuery(deviceID, channelID)
}

func (s *PTZService) PresetCommand(deviceID, channelID string, cmdCode, presetID, speed int) error {
	return s.cmd.FrontEndCmd(deviceID, channelID, cmdCode, presetID, speed, 0)
}

func (s *PTZService) CruiseCommand(deviceID, channelID string, cmdCode, routeID, param int) error {
	return s.cmd.FrontEndCmd(deviceID, channelID, cmdCode, param, routeID, 0)
}

func (s *PTZService) ScanCommand(deviceID, channelID string, cmdCode, param int) error {
	return s.cmd.FrontEndCmd(deviceID, channelID, cmdCode, param, 0, 0)
}

func (s *PTZService) Wiper(deviceID, channelID string, on bool) error {
	param := 1
	if !on {
		param = 0
	}
	return s.cmd.FrontEndCmd(deviceID, channelID, 0, param, 0, 0)
}

func (s *PTZService) Auxiliary(deviceID, channelID string, switchID int, on bool) error {
	param := 1
	if !on {
		param = 0
	}
	return s.cmd.FrontEndCmd(deviceID, channelID, 0, switchID, param, 0)
}
