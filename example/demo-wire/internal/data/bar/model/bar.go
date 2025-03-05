package model

type Bar struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}

func (Bar) TableName() string {
	return "bar"
}
