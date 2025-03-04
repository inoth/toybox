package service

import (
	"demo-wire/internal/model"

	database "github.com/inoth/toybox/component/database/sqlite"
	"gorm.io/gorm"
)

type UserInfoService struct {
	db *gorm.DB
}

func NewUserInfoService(db *database.SqliteComponent) *UserInfoService {
	user := db.GetDatabase("user")
	user.AutoMigrate(&model.UserInfo{})
	return &UserInfoService{db: user}
}

func (s *UserInfoService) GetUserInfoByName(name string) (*model.UserInfo, error) {
	var user model.UserInfo
	if err := s.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserInfoService) CreateUserInfo(name string) (uint, error) {
	user := model.UserInfo{
		Name: name,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}
