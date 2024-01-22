package service

import (
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/component/mysql"

	"gorm.io/gorm"
)

type UserInfo struct {
	Id   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (UserInfo) TableName() string {
	return "ts_user_info"
}

type FooService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewFooService(conf toybox.ConfigMate) *FooService {
	var err error
	fs := &FooService{
		log: logger.GetLogger(logger.LoggerConfig{ServerName: "FooService"}),
	}
	fs.db, err = mysql.GetDB(mysql.SetName("mysql2"), mysql.SetConfig(conf))
	if err != nil {
		panic(err)
	}
	return fs
}

func (fs *FooService) SayHi() string {
	var user UserInfo
	err := fs.db.Model(UserInfo{}).Where("id = ?", 157).First(&user).Error
	if err != nil {
		return "Hi, I'm FooService"
	}
	fs.log.Info("Hi, I'm FooService")
	return user.Name
}
