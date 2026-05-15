package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
)

type UserService struct {
	log *zap.Logger
}

func NewUserService(log *zap.Logger) *UserService {
	return &UserService{log: log}
}

func (s *UserService) Login(username, password string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ? AND enable = ?", username, true).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Frontend already sends MD5-hashed password, compare directly
	if password != user.Password {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return &user, nil
}

func (s *UserService) GetUsers(page, count int) (*utils.PageInfo[any], error) {
	var users []model.User
	var total int64

	db := database.DB.Model(&model.User{})
	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("id DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	list := make([]any, len(users))
	for i := range users {
		list[i] = users[i]
	}
	return utils.NewPageInfo[any](total, list, page, count), nil
}

func (s *UserService) GetUser(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (s *UserService) CreateUser(user *model.User) error {
	user.Password = hashPassword(user.Password)
	return database.DB.Create(user).Error
}

func (s *UserService) UpdateUser(user *model.User) error {
	if user.Password != "" {
		user.Password = hashPassword(user.Password)
	}
	return database.DB.Save(user).Error
}

func (s *UserService) DeleteUser(id uint) error {
	return database.DB.Delete(&model.User{}, id).Error
}

func hashPassword(password string) string {
	h := md5.Sum([]byte(password))
	return hex.EncodeToString(h[:])
}

// UserCreateRequest is the request body for creating a user
type UserCreateRequest struct {
	Username string
	Password string
}

// CreateUserFromHandler creates a user from handler request
func (s *UserService) CreateUserFromHandler(req *UserCreateRequest) error {
	user := &model.User{
		Username:   req.Username,
		Password:   hashPassword(req.Password),
		Enable:     true,
		CreateTime: fmt.Sprintf("%d", 0),
	}
	return database.DB.Create(user).Error
}

// ChangePassword changes user password after verifying old password
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user model.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	if hashPassword(oldPassword) != user.Password {
		return fmt.Errorf("旧密码错误")
	}

	return database.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", hashPassword(newPassword)).Error
}
