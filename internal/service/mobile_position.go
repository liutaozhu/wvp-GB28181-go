package service

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
)

type MobilePositionService struct {
	log *zap.Logger
}

func NewMobilePositionService(log *zap.Logger) *MobilePositionService {
	return &MobilePositionService{log: log}
}

func (s *MobilePositionService) List(deviceID, channelID string, page, count int) (*utils.PageInfo[any], error) {
	var positions []model.MobilePosition
	var total int64

	db := database.DB.Model(&model.MobilePosition{})
	if deviceID != "" {
		db = db.Where("device_id = ?", deviceID)
	}
	if channelID != "" {
		db = db.Where("channel_id = ?", channelID)
	}

	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("time DESC").Find(&positions).Error; err != nil {
		return nil, err
	}

	list := make([]any, len(positions))
	for i := range positions {
		list[i] = positions[i]
	}
	return utils.NewPageInfo[any](total, list, page, count), nil
}

func (s *MobilePositionService) Save(pos *model.MobilePosition) error {
	return database.DB.Create(pos).Error
}

func (s *MobilePositionService) GetLatest(deviceID string) (*model.MobilePosition, error) {
	var pos model.MobilePosition
	err := database.DB.Where("device_id = ?", deviceID).Order("time DESC").First(&pos).Error
	return &pos, err
}

func (s *MobilePositionService) CleanupOldPositions(days int) error {
	return nil
}
