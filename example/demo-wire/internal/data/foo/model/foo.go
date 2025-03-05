package model

type Foo struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}

func (Foo) TableName() string {
	return "foo"
}
