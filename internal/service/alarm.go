package service

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
)

type AlarmService struct {
	log *zap.Logger
}

func NewAlarmService(log *zap.Logger) *AlarmService {
	return &AlarmService{log: log}
}

func (s *AlarmService) List(deviceID string, page, count int) (*utils.PageInfo[model.DeviceAlarm], error) {
	var alarms []model.DeviceAlarm
	var total int64

	db := database.DB.Model(&model.DeviceAlarm{})
	if deviceID != "" {
		db = db.Where("device_id = ?", deviceID)
	}

	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("alarm_time DESC").Find(&alarms).Error; err != nil {
		return nil, err
	}

	return utils.NewPageInfo[model.DeviceAlarm](total, alarms, page, count), nil
}

func (s *AlarmService) GetOne(id uint) (*model.DeviceAlarm, error) {
	var a model.DeviceAlarm
	err := database.DB.Where("id = ?", id).First(&a).Error
	return &a, err
}

func (s *AlarmService) Delete(id uint) error {
	return database.DB.Delete(&model.DeviceAlarm{}, id).Error
}

func (s *AlarmService) CleanupOldAlarms(days int) error {
	return nil
}
