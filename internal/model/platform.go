package model

// Platform maps to wvp_platform table
type Platform struct {
	ID                        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Enable                    int    `gorm:"default:0" json:"enable"`
	Name                      string `gorm:"type:varchar(255)" json:"name"`
	ServerGBID                string `gorm:"column:server_gb_id;type:varchar(50);uniqueIndex" json:"serverGBId"`
	ServerGBDomain            string `gorm:"column:server_gb_domain;type:varchar(50)" json:"serverGBDomain"`
	ServerIP                  string `gorm:"column:server_ip;type:varchar(50)" json:"serverIp"`
	ServerPort                int    `gorm:"column:server_port" json:"serverPort"`
	DeviceGBID                string `gorm:"column:device_gb_id;type:varchar(50)" json:"deviceGBId"`
	DeviceIP                  string `gorm:"column:device_ip;type:varchar(50)" json:"deviceIp"`
	DevicePort                int    `gorm:"column:device_port" json:"devicePort"`
	Username                  string `gorm:"type:varchar(255)" json:"username"`
	Password                  string `gorm:"type:varchar(255)" json:"password"`
	Expires                   int    `json:"expires"`
	KeepTimeout               int    `gorm:"column:keep_timeout" json:"keepTimeout"`
	Transport                 string `gorm:"type:varchar(50)" json:"transport"`
	CharacterSet              string `gorm:"column:character_set;type:varchar(50)" json:"characterSet"`
	PTZ                       int    `gorm:"column:ptz" json:"ptz"`
	RTCP                      int    `gorm:"column:rtcp" json:"rtcp"`
	Status                    string `gorm:"type:varchar(50)" json:"status"`
	ChannelCount              int    `gorm:"column:channel_count" json:"channelCount"`
	CatalogSubscribe          int    `gorm:"column:catalog_subscribe" json:"catalogSubscribe"`
	AlarmSubscribe            int    `gorm:"column:alarm_subscribe" json:"alarmSubscribe"`
	MobilePositionSubscribe   int    `gorm:"column:mobile_position_subscribe" json:"mobilePositionSubscribe"`
	CatalogGroup              int    `gorm:"column:catalog_group" json:"catalogGroup"`
	AsMessageChannel          bool   `gorm:"default:false" json:"asMessageChannel"`
	SendStreamIP              string `gorm:"column:send_stream_ip;type:varchar(50)" json:"sendStreamIp"`
	AutoPushChannel           bool   `gorm:"default:false" json:"autoPushChannel"`
	CatalogWithPlatform       int    `gorm:"column:catalog_with_platform" json:"catalogWithPlatform"`
	CatalogWithGroup          int    `gorm:"column:catalog_with_group" json:"catalogWithGroup"`
	CatalogWithRegion         int    `gorm:"column:catalog_with_region" json:"catalogWithRegion"`
	CivilCode                 string `gorm:"type:varchar(50)" json:"civilCode"`
	Manufacturer              string `gorm:"type:varchar(255)" json:"manufacturer"`
	Model                     string `gorm:"type:varchar(255)" json:"model"`
	Address                   string `gorm:"type:varchar(255)" json:"address"`
	RegisterWay               int    `gorm:"column:register_way" json:"registerWay"`
	Secrecy                   int    `json:"secrecy"`
	ServerID                  string `gorm:"column:server_id;type:varchar(50)" json:"serverId"`
}

func (Platform) TableName() string {
	return "wvp_platform"
}

// PlatformChannel extends CommonGBChannel with custom fields per platform
type PlatformChannel struct {
	CommonGBChannel

	ID                        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PlatformID                uint   `gorm:"column:platform_id;index" json:"platformId"`
	CustomDeviceID            string `gorm:"column:custom_device_id;type:varchar(50)" json:"customDeviceId"`
	CustomName                string `gorm:"column:custom_name;type:varchar(255)" json:"customName"`
	CustomManufacturer        string `gorm:"column:custom_manufacturer;type:varchar(50)" json:"customManufacturer"`
	CustomModel               string `gorm:"column:custom_model;type:varchar(50)" json:"customModel"`
	CustomOwner               string `gorm:"column:custom_owner;type:varchar(50)" json:"customOwner"`
	CustomCivilCode           string `gorm:"column:custom_civil_code;type:varchar(50)" json:"customCivilCode"`
	CustomBlock               string `gorm:"column:custom_block;type:varchar(50)" json:"customBlock"`
	CustomAddress             string `gorm:"column:custom_address;type:varchar(50)" json:"customAddress"`
	CustomParental            int    `gorm:"column:custom_parental" json:"customParental"`
	CustomParentID            string `gorm:"column:custom_parent_id;type:varchar(50)" json:"customParentId"`
	CustomSafetyWay           int    `gorm:"column:custom_safety_way" json:"customSafetyWay"`
	CustomRegisterWay         int    `gorm:"column:custom_register_way" json:"customRegisterWay"`
	CustomCertNum             string `gorm:"column:custom_cert_num;type:varchar(50)" json:"customCertNum"`
	CustomCertifiable         int    `gorm:"column:custom_certifiable" json:"customCertifiable"`
	CustomErrCode             int    `gorm:"column:custom_err_code" json:"customErrCode"`
	CustomEndTime             string `gorm:"column:custom_end_time;type:varchar(50)" json:"customEndTime"`
	CustomSecrecy             int    `gorm:"column:custom_secrecy" json:"customSecrecy"`
	CustomIPAddress           string `gorm:"column:custom_ip_address;type:varchar(50)" json:"customIpAddress"`
	CustomPort                int    `gorm:"column:custom_port" json:"customPort"`
	CustomPassword            string `gorm:"column:custom_password;type:varchar(255)" json:"customPassword"`
	CustomStatus              string `gorm:"column:custom_status;type:varchar(50)" json:"customStatus"`
	CustomLongitude           float64 `gorm:"column:custom_longitude" json:"customLongitude"`
	CustomLatitude            float64 `gorm:"column:custom_latitude" json:"customLatitude"`
	CustomBusinessGroupID     string `gorm:"column:custom_business_group_id;type:varchar(50)" json:"customBusinessGroupId"`
	CustomPtzType             int    `gorm:"column:custom_ptz_type" json:"customPtzType"`
	CustomPositionType        int    `gorm:"column:custom_position_type" json:"customPositionType"`
	CustomRoomType            int    `gorm:"column:custom_room_type" json:"customRoomType"`
	CustomUseType             int    `gorm:"column:custom_use_type" json:"customUseType"`
	CustomSupplyLightType     int    `gorm:"column:custom_supply_light_type" json:"customSupplyLightType"`
	CustomDirectionType       int    `gorm:"column:custom_direction_type" json:"customDirectionType"`
	CustomResolution          string `gorm:"column:custom_resolution;type:varchar(50)" json:"customResolution"`
}

func (PlatformChannel) TableName() string {
	return "wvp_platform_channel"
}
