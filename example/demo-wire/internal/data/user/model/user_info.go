package model

type UserInfo struct {
	Id    uint   `gorm:"column:id;primaryKey"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}
