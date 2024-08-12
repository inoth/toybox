package service

import (
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

// func (us *UserService) Query(uid string) []any{

// }
