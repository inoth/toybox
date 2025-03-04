package model

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	Name string `json:"name"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}
