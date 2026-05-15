package model

// CommonGBChannel maps to wvp_device_channel table (GB standard fields)
type CommonGBChannel struct {
	GBID                     uint    `gorm:"column:gb_id;type:bigint;index" json:"gbId"`
	GBDeviceID               string  `gorm:"column:gb_device_id;type:varchar(50)" json:"gbDeviceId"`
	GBName                   string  `gorm:"column:gb_name;type:varchar(255)" json:"gbName"`
	GBManufacturer           string  `gorm:"column:gb_manufacturer;type:varchar(50)" json:"gbManufacturer"`
	GBModel                  string  `gorm:"column:gb_model;type:varchar(50)" json:"gbModel"`
	GBOwner                  string  `gorm:"column:gb_owner;type:varchar(50)" json:"gbOwner"`
	GBCivilCode              string  `gorm:"column:gb_civil_code;type:varchar(50)" json:"gbCivilCode"`
	GBBlock                  string  `gorm:"column:gb_block;type:varchar(50)" json:"gbBlock"`
	GBAddress                string  `gorm:"column:gb_address;type:varchar(50)" json:"gbAddress"`
	GBParental               int     `gorm:"column:gb_parental" json:"gbParental"`
	GBParentID               string  `gorm:"column:gb_parent_id;type:varchar(50)" json:"gbParentId"`
	GBSafetyWay              int     `gorm:"column:gb_safety_way" json:"gbSafetyWay"`
	GBRegisterWay            int     `gorm:"column:gb_register_way" json:"gbRegisterWay"`
	GBCertNum                string  `gorm:"column:gb_cert_num;type:varchar(50)" json:"gbCertNum"`
	GBCertifiable            int     `gorm:"column:gb_certifiable" json:"gbCertifiable"`
	GBErrCode                int     `gorm:"column:gb_err_code" json:"gbErrCode"`
	GBEndTime                string  `gorm:"column:gb_end_time;type:varchar(50)" json:"gbEndTime"`
	GBSecrecy                int     `gorm:"column:gb_secrecy" json:"gbSecrecy"`
	GBIPAddress              string  `gorm:"column:gb_ip_address;type:varchar(50)" json:"gbIpAddress"`
	GBPort                   int     `gorm:"column:gb_port" json:"gbPort"`
	GBPassword               string  `gorm:"column:gb_password;type:varchar(255)" json:"gbPassword"`
	GBStatus                 string  `gorm:"column:gb_status;type:varchar(50)" json:"gbStatus"`
	GBLongitude              float64 `gorm:"column:gb_longitude" json:"gbLongitude"`
	GBLatitude               float64 `gorm:"column:gb_latitude" json:"gbLatitude"`
	GPSAltitude              float64 `gorm:"column:gps_altitude" json:"gpsAltitude"`
	GPSSpeed                 float64 `gorm:"column:gps_speed" json:"gpsSpeed"`
	GPSDirection             float64 `gorm:"column:gps_direction" json:"gpsDirection"`
	GPSTime                  string  `gorm:"column:gps_time;type:varchar(50)" json:"gpsTime"`
	GBBusinessGroupID        string  `gorm:"column:gb_business_group_id;type:varchar(50)" json:"gbBusinessGroupId"`
	GBPtzType                int     `gorm:"column:gb_ptz_type" json:"gbPtzType"`
	GBPositionType           int     `gorm:"column:gb_position_type" json:"gbPositionType"`
	GBRoomType               int     `gorm:"column:gb_room_type" json:"gbRoomType"`
	GBUseType                int     `gorm:"column:gb_use_type" json:"gbUseType"`
	GBSupplyLightType        int     `gorm:"column:gb_supply_light_type" json:"gbSupplyLightType"`
	GBDirectionType          int     `gorm:"column:gb_direction_type" json:"gbDirectionType"`
	GBResolution             string  `gorm:"column:gb_resolution;type:varchar(50)" json:"gbResolution"`
	GBDownloadSpeed          string  `gorm:"column:gb_download_speed;type:varchar(50)" json:"gbDownloadSpeed"`
	GBSvcSpaceSupportMod     int     `gorm:"column:gb_svc_space_support_mod" json:"gbSvcSpaceSupportMod"`
	GBSvcTimeSupportMode     int     `gorm:"column:gb_svc_time_support_mode" json:"gbSvcTimeSupportMode"`
	RecordPLan               string  `gorm:"column:record_plan;type:varchar(50)" json:"recordPLan"`
	DataType                 int     `gorm:"column:data_type;default:1" json:"dataType"`
	DataDeviceID             int     `gorm:"column:data_device_id" json:"dataDeviceId"`
	StreamIdentification     string  `gorm:"column:stream_identification;type:varchar(50)" json:"streamIdentification"`
	EnableBroadcast          int     `gorm:"column:enable_broadcast;default:0" json:"enableBroadcast"`
	MapLevel                 int     `gorm:"column:map_level;default:0" json:"mapLevel"`
}

// DeviceChannel maps to wvp_device_channel (extends CommonGBChannel)
type DeviceChannel struct {
	CommonGBChannel

	ID             uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID       string  `gorm:"column:device_id;type:varchar(50);index" json:"deviceId"`
	Name           string  `gorm:"type:varchar(255)" json:"name"`
	Manufacturer   string  `gorm:"type:varchar(50)" json:"manufacturer"`
	Model          string  `gorm:"type:varchar(50)" json:"model"`
	Owner          string  `gorm:"type:varchar(50)" json:"owner"`
	CivilCode      string  `gorm:"type:varchar(50)" json:"civilCode"`
	Block          string  `gorm:"type:varchar(50)" json:"block"`
	Address        string  `gorm:"type:varchar(50)" json:"address"`
	Parental       int     `json:"parental"`
	ParentID       string  `gorm:"column:parent_id;type:varchar(50)" json:"parentId"`
	SafetyWay      int     `json:"safetyWay"`
	RegisterWay    int     `json:"registerWay"`
	CertNum        string  `gorm:"column:cert_num;type:varchar(50)" json:"certNum"`
	Certifiable    int     `json:"certifiable"`
	ErrCode        int     `gorm:"column:err_code" json:"errCode"`
	EndTime        string  `gorm:"column:end_time;type:varchar(50)" json:"endTime"`
	Secrecy        int     `json:"secrecy"`
	IPAddress      string  `gorm:"column:ip_address;type:varchar(50)" json:"ipAddress"`
	Port           int     `json:"port"`
	Password       string  `gorm:"type:varchar(255)" json:"password"`
	Status         string  `gorm:"type:varchar(50)" json:"status"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	PTZType        int     `gorm:"column:ptz_type" json:"ptzType"`
	PositionType   int     `gorm:"column:position_type" json:"positionType"`
	RoomType       int     `gorm:"column:room_type" json:"roomType"`
	UseType        int     `gorm:"column:use_type" json:"useType"`
	SupplyLightType int    `gorm:"column:supply_light_type" json:"supplyLightType"`
	DirectionType  int     `gorm:"column:direction_type" json:"directionType"`
	Resolution     string  `gorm:"type:varchar(50)" json:"resolution"`

	// Extra fields
	ParentName     string `gorm:"-" json:"parentName"`
	PtzTypeText    string `gorm:"-" json:"ptzTypeText"`
	SubCount       int    `gorm:"-" json:"subCount"`
	HasAudio       bool   `gorm:"-" json:"hasAudio"`
	GPSTime        string `gorm:"-" json:"gpsTime"`
	ChannelType    int    `gorm:"-" json:"channelType"`
	StreamID       string `gorm:"column:stream_id;type:varchar(50)" json:"streamId"`
}

func (DeviceChannel) TableName() string {
	return "wvp_device_channel"
}
