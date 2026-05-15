package model

// Group maps to wvp_common_group table
type Group struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID       string `gorm:"column:device_id;type:varchar(50);uniqueIndex" json:"deviceId"`
	Name           string `gorm:"type:varchar(255)" json:"name"`
	ParentID       uint   `gorm:"column:parent_id" json:"parentId"`
	ParentDeviceID string `gorm:"column:parent_device_id;type:varchar(50)" json:"parentDeviceId"`
	BusinessGroup  string `gorm:"column:business_group;type:varchar(50)" json:"businessGroup"`
	CivilCode      string `gorm:"type:varchar(50)" json:"civilCode"`
	Alias          string `gorm:"type:varchar(255)" json:"alias"`
	CreateTime     string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime     string `gorm:"type:varchar(50)" json:"updateTime"`
}

func (Group) TableName() string {
	return "wvp_common_group"
}

// Region maps to wvp_common_region table
type Region struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID       string `gorm:"column:device_id;type:varchar(50);uniqueIndex" json:"deviceId"`
	Name           string `gorm:"type:varchar(255)" json:"name"`
	ParentID       uint   `gorm:"column:parent_id" json:"parentId"`
	ParentDeviceID string `gorm:"column:parent_device_id;type:varchar(50)" json:"parentDeviceId"`
	CreateTime     string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime     string `gorm:"type:varchar(50)" json:"updateTime"`
}

func (Region) TableName() string {
	return "wvp_common_region"
}
