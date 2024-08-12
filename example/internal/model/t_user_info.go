package model

import "time"

type UserInfo struct {
	Id          string    `gorm:"id"`
	Nickname    string    `gorm:"nickname"`
	Username    string    `gorm:"username"`
	Passwd      string    `gorm:"passwd"`
	Avatar      string    `gorm:"avatar"`
	AccountType string    `gorm:"account_type"`
	CreatedTime time.Time `gorm:"created_time"`
}

func (u *UserInfo) TableName() string {
	return "t_user_info"
}
