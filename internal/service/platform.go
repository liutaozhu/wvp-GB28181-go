package service

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
)

type PlatformService struct {
	log *zap.Logger
}

func NewPlatformService(log *zap.Logger) *PlatformService {
	return &PlatformService{log: log}
}

func (s *PlatformService) GetPlatform(id uint) (*model.Platform, error) {
	var p model.Platform
	err := database.DB.Where("id = ?", id).First(&p).Error
	return &p, err
}

func (s *PlatformService) ListPlatforms(page, count int) (*utils.PageInfo[any], error) {
	var platforms []model.Platform
	var total int64

	db := database.DB.Model(&model.Platform{})
	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("id DESC").Find(&platforms).Error; err != nil {
		return nil, err
	}

	list := make([]any, len(platforms))
	for i := range platforms {
		list[i] = platforms[i]
	}
	return utils.NewPageInfo[any](total, list, page, count), nil
}

func (s *PlatformService) AddPlatform(p *model.Platform) error {
	return database.DB.Create(p).Error
}

func (s *PlatformService) UpdatePlatform(p *model.Platform) error {
	return database.DB.Save(p).Error
}

func (s *PlatformService) DeletePlatform(id uint) error {
	return database.DB.Delete(&model.Platform{}, id).Error
}

func (s *PlatformService) GetPlatformChannels(platformID uint, page, count int) (*utils.PageInfo[any], error) {
	var channels []model.PlatformChannel
	var total int64

	db := database.DB.Model(&model.PlatformChannel{}).Where("platform_id = ?", platformID)
	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Find(&channels).Error; err != nil {
		return nil, err
	}

	list := make([]any, len(channels))
	for i := range channels {
		list[i] = channels[i]
	}
	return utils.NewPageInfo[any](total, list, page, count), nil
}

func (s *PlatformService) AddPlatformChannel(pc *model.PlatformChannel) error {
	return database.DB.Create(pc).Error
}

func (s *PlatformService) RemovePlatformChannel(platformID, channelID uint) error {
	return database.DB.Where("platform_id = ? AND gb_id = ?", platformID, channelID).Delete(&model.PlatformChannel{}).Error
}

func (s *PlatformService) PushChannelToPlatform(platformID, channelID uint) error {
	// Push channel data to upper platform via SIP
	return nil
}
