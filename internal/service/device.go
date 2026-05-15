package service

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/redis"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DeviceService struct {
	log *zap.Logger
}

func NewDeviceService(log *zap.Logger) *DeviceService {
	return &DeviceService{log: log}
}

func (s *DeviceService) GetDevice(deviceID string) (*model.Device, error) {
	var device model.Device
	err := database.DB.Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *DeviceService) GetOnlineDevices(page, count int) (*utils.PageInfo[model.Device], error) {
	return s.getAllDevicesTyped("", boolPtr(true), page, count)
}

func (s *DeviceService) GetAllDevices(query string, online *bool, page, count int) (*utils.PageInfo[model.Device], error) {
	return s.getAllDevicesTyped(query, online, page, count)
}

func (s *DeviceService) getAllDevicesTyped(query string, online *bool, page, count int) (*utils.PageInfo[model.Device], error) {
	var devices []model.Device
	var total int64

	db := database.DB.Model(&model.Device{})
	if online != nil {
		db = db.Where("on_line = ?", *online)
	}
	if query != "" {
		db = db.Where("device_id LIKE ? OR name LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("id DESC").Find(&devices).Error; err != nil {
		return nil, err
	}

	return utils.NewPageInfo[model.Device](total, devices, page, count), nil
}

func (s *DeviceService) UpdateDevice(device *model.Device) error {
	return database.DB.Save(device).Error
}

func (s *DeviceService) DeleteDevice(deviceID string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("device_id = ?", deviceID).Delete(&model.Device{}).Error; err != nil {
			return err
		}
		return tx.Where("device_id = ?", deviceID).Delete(&model.DeviceChannel{}).Error
	})
}

func (s *DeviceService) SetDeviceOnline(device *model.Device) {
	device.Online = true
	database.DB.Model(&model.Device{}).Where("device_id = ?", device.DeviceID).Update("on_line", true)
	redis.Client.Set(redis.Ctx, redis.DeviceOnline+device.DeviceID, "1", 0)
}

func (s *DeviceService) SetDeviceOffline(device *model.Device) {
	device.Online = false
	database.DB.Model(&model.Device{}).Where("device_id = ?", device.DeviceID).Update("on_line", false)
	redis.Client.Del(redis.Ctx, redis.DeviceOnline+device.DeviceID)
}

func (s *DeviceService) UpdateKeepalive(deviceID string) {
	database.DB.Model(&model.Device{}).Where("device_id = ?", deviceID).Update("on_line", true)
	redis.Client.Set(redis.Ctx, redis.DeviceOnline+deviceID, "1", 0)
}

func boolPtr(b bool) *bool {
	return &b
}
