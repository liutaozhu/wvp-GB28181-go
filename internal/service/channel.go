package service

import (
	"fmt"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
)

type ChannelService struct {
	log *zap.Logger
}

func NewChannelService(log *zap.Logger) *ChannelService {
	return &ChannelService{log: log}
}

func (s *ChannelService) GetChannel(deviceID, channelID string) (*model.DeviceChannel, error) {
	var ch model.DeviceChannel
	err := database.DB.Where("device_id = ? AND gb_device_id = ?", deviceID, channelID).First(&ch).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *ChannelService) GetChannelByID(id uint) (*model.DeviceChannel, error) {
	var ch model.DeviceChannel
	err := database.DB.Where("gb_id = ?", id).First(&ch).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *ChannelService) GetChannelByGBDeviceID(gbDeviceID string) (*model.DeviceChannel, error) {
	var ch model.DeviceChannel
	err := database.DB.Where("gb_device_id = ?", gbDeviceID).First(&ch).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *ChannelService) GetChannels(deviceID string, page, count int) (*utils.PageInfo[model.DeviceChannel], error) {
	var channels []model.DeviceChannel
	var total int64

	db := database.DB.Model(&model.DeviceChannel{})
	if deviceID != "" {
		db = db.Where("device_id = ?", deviceID)
	}

	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("gb_id ASC").Find(&channels).Error; err != nil {
		return nil, err
	}

	return utils.NewPageInfo[model.DeviceChannel](total, channels, page, count), nil
}

func (s *ChannelService) GetAllChannels(query string, online, hasSub, hasStream *bool, page, count int) (*utils.PageInfo[model.DeviceChannel], error) {
	var channels []model.DeviceChannel
	var total int64

	db := database.DB.Model(&model.DeviceChannel{})
	if query != "" {
		db = db.Where("gb_device_id LIKE ? OR gb_name LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if online != nil {
		db = db.Where("status = ?", func() string {
			if *online {
				return "ON"
			}
			return "OFF"
		}())
	}
	if hasSub != nil && *hasSub {
		db = db.Where("gb_parental = 1")
	}
	if hasStream != nil && *hasStream {
		db = db.Where("stream_id IS NOT NULL AND stream_id != ''")
	}

	db.Count(&total)

	// Also count stream proxies
	var proxyCount int64
	proxyDB := database.DB.Model(&model.StreamProxy{})
	if query != "" {
		proxyDB = proxyDB.Where("name LIKE ? OR app LIKE ? OR stream LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	if online != nil {
		if *online {
			proxyDB = proxyDB.Where("pulling = ?", true)
		} else {
			proxyDB = proxyDB.Where("pulling = ?", false)
		}
	}
	proxyDB.Count(&proxyCount)
	total += proxyCount

	if err := db.Offset((page - 1) * count).Limit(count).Order("gb_id ASC").Find(&channels).Error; err != nil {
		return nil, err
	}

	// If current page doesn't fill up from device channels, supplement with proxy channels
	remaining := count - len(channels)
	deviceChannelCount := int64(0)
	database.DB.Model(&model.DeviceChannel{}).Count(&deviceChannelCount)

	if remaining > 0 {
		proxyOffset := (page-1)*count - int(deviceChannelCount)
		if proxyOffset < 0 {
			proxyOffset = 0
		}
		var proxies []model.StreamProxy
		pdb := database.DB.Model(&model.StreamProxy{})
		if query != "" {
			pdb = pdb.Where("name LIKE ? OR app LIKE ? OR stream LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
		}
		if online != nil {
			if *online {
				pdb = pdb.Where("pulling = ?", true)
			} else {
				pdb = pdb.Where("pulling = ?", false)
			}
		}
		pdb.Offset(proxyOffset).Limit(remaining).Order("id ASC").Find(&proxies)

		for _, p := range proxies {
			status := "OFF"
			if p.Pulling {
				status = "ON"
			}
			ch := model.DeviceChannel{
				DeviceID: "proxy",
				Name:     p.Name,
				Status:   status,
				StreamID: fmt.Sprintf("%s/%s", p.App, p.Stream),
			}
			ch.GBID = 100000 + uint(p.ID)
			ch.GBDeviceID = p.DeviceID
			ch.GBName = p.Name
			ch.GBStatus = status
			ch.DataType = 3 // proxy type
			channels = append(channels, ch)
		}
	}

	return utils.NewPageInfo[model.DeviceChannel](total, channels, page, count), nil
}

func (s *ChannelService) UpdateChannel(ch *model.DeviceChannel) error {
	return database.DB.Save(ch).Error
}

func (s *ChannelService) DeleteChannel(id uint) error {
	return database.DB.Delete(&model.DeviceChannel{}, id).Error
}

func (s *ChannelService) GetCivilCodeList() ([]string, error) {
	var civilCodes []string
	err := database.DB.Model(&model.DeviceChannel{}).
		Distinct("gb_civil_code").
		Where("gb_civil_code IS NOT NULL AND gb_civil_code != ''").
		Pluck("gb_civil_code", &civilCodes).Error
	return civilCodes, err
}

func (s *ChannelService) GetParentChannels() ([]model.DeviceChannel, error) {
	var channels []model.DeviceChannel
	err := database.DB.Where("gb_parental = 1").Find(&channels).Error
	return channels, err
}

func (s *ChannelService) StartPlay(channelID uint, stream string) error {
	return database.DB.Model(&model.DeviceChannel{}).Where("gb_id = ?", channelID).Update("stream_id", stream).Error
}

func (s *ChannelService) StopPlay(channelID uint) error {
	return database.DB.Model(&model.DeviceChannel{}).Where("gb_id = ?", channelID).Update("stream_id", "").Error
}

func (s *ChannelService) UpdateChannelGPS(device *model.Device, channel *model.DeviceChannel, pos *model.MobilePosition) error {
	// Convert GCJ02 to WGS84 if needed
	if device.GeoCoordSys == "GCJ02" {
		pos.Longitude, pos.Latitude = gcj02ToWGS84(pos.Longitude, pos.Latitude)
	}

	updates := map[string]interface{}{
		"gb_longitude": pos.Longitude,
		"gb_latitude":  pos.Latitude,
		"gps_altitude": pos.Altitude,
		"gps_speed":    pos.Speed,
		"gps_direction": pos.Direction,
		"gps_time":     pos.Time,
	}
	return database.DB.Model(&model.DeviceChannel{}).Where("gb_id = ?", channel.GBID).Updates(updates).Error
}

func (s *ChannelService) GetRegionTree() ([]map[string]interface{}, error) {
	var civilCodes []string
	if err := database.DB.Model(&model.DeviceChannel{}).
		Distinct("gb_civil_code").
		Where("gb_civil_code IS NOT NULL AND gb_civil_code != ''").
		Pluck("gb_civil_code", &civilCodes).Error; err != nil {
		return nil, err
	}

	tree := buildCivilCodeTree(civilCodes)
	return tree, nil
}

func buildCivilCodeTree(codes []string) []map[string]interface{} {
	root := make(map[string]map[string]interface{})
	for _, code := range codes {
		if code == "" {
			continue
		}
		parts := []rune(code)
		current := root
		for i := 2; i <= len(parts); i += 2 {
			key := string(parts[:i])
			if _, ok := current[key]; !ok {
				current[key] = map[string]interface{}{
					"id":       key,
					"name":     key,
					"children": make(map[string]map[string]interface{}),
				}
			}
			current = current[key]["children"].(map[string]map[string]interface{})
		}
	}
	return flattenTree(root)
}

func flattenTree(root map[string]map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(root))
	for _, node := range root {
		children := flattenTree(node["children"].(map[string]map[string]interface{}))
		n := map[string]interface{}{
			"id":   node["id"],
			"name": node["id"],
		}
		if len(children) > 0 {
			n["children"] = children
		}
		result = append(result, n)
	}
	return result
}

// gcj02ToWGS84 converts GCJ02 coordinates to WGS84
func gcj02ToWGS84(lng, lat float64) (float64, float64) {
	// Simplified conversion - in production use a proper library
	dlat := transformLat(lng-105.0, lat-35.0)
	dlng := transformLng(lng-105.0, lat-35.0)
	radlat := lat / 180.0 * 3.14159265358979323846
	magic := 1.0 - 0.00669342162296594323 * sin(radlat) * sin(radlat)
	sqrtmagic := mathSqrt(magic)
	dlat = (dlat * 180.0) / ((6378245.0 * (1.0 - 0.00669342162296594323)) / (magic * sqrtmagic) * 3.14159265358979323846)
	dlng = (dlng * 180.0) / (6378245.0 / sqrtmagic * cos(radlat) * 3.14159265358979323846)
	return lng - dlng, lat - dlat
}

func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*mathSqrt(mathAbs(x))
	ret += (20.0*sin(6.0*x*3.14159265358979323846) + 20.0*sin(2.0*x*3.14159265358979323846)) * 2.0 / 3.0
	ret += (20.0*sin(y*3.14159265358979323846) + 40.0*sin(y/3.0*3.14159265358979323846)) * 2.0 / 3.0
	ret += (160.0*sin(y/12.0*3.14159265358979323846) + 320*sin(y*3.14159265358979323846/30.0)) * 2.0 / 3.0
	return ret
}

func transformLng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*mathSqrt(mathAbs(x))
	ret += (20.0*sin(6.0*x*3.14159265358979323846) + 20.0*sin(2.0*x*3.14159265358979323846)) * 2.0 / 3.0
	ret += (20.0*sin(x*3.14159265358979323846) + 40.0*sin(x/3.0*3.14159265358979323846)) * 2.0 / 3.0
	ret += (150.0*sin(x/12.0*3.14159265358979323846) + 300.0*sin(x/30.0*3.14159265358979323846)) * 2.0 / 3.0
	return ret
}

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func mathSqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func sin(x float64) float64 {
	// Taylor series approximation
	x = x - float64(int(x/(2*3.14159265358979323846)))*2*3.14159265358979323846
	result := x
	term := x
	for i := 1; i <= 10; i++ {
		term *= -x * x / (float64(2*i) * float64(2*i+1))
		result += term
	}
	return result
}

func cos(x float64) float64 {
	return sin(x + 3.14159265358979323846/2)
}
