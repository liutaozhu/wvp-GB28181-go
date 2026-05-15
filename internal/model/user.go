package model

// User maps to wvp_user table
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	Email    string `gorm:"type:varchar(255)" json:"email"`
	Enable   bool   `gorm:"default:true" json:"enable"`
	CreateTime string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime string `gorm:"type:varchar(50)" json:"updateTime"`
}

func (User) TableName() string {
	return "wvp_user"
}

// Role maps to wvp_user_role table
type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	CreateTime  string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime  string `gorm:"type:varchar(50)" json:"updateTime"`
}

func (Role) TableName() string {
	return "wvp_user_role"
}
