package service

import (
	"example/internal/model"

	"github.com/inoth/toybox/component/database"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *database.MysqlComponent) *UserService {
	return &UserService{
		db: db.GetDB(),
	}
}

// func NewUserService() *UserService {
// 	return &UserService{}
// }

func (us *UserService) Query(uid string) *model.UserInfo {
	var user model.UserInfo
	if err := us.db.Model(model.UserInfo{}).Where("id = ?", uid).First(&user).Error; err != nil {
		return nil
	}
	// return fmt.Sprintf("hello %v", uid)
	return &user
}
