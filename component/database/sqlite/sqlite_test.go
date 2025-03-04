package database

import (
	"os"
	"testing"

	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Name  string
	Email string
}

func (u *UserInfo) TableName() string {
	return "ts_user_info"
}

func TestNewSqliteDBConnect(t *testing.T) {
	sc := &SqliteComponent{
		DbMap: make(map[string]*gorm.DB),
		Configs: []Config{{
			DBPath:          "test.db",
			DbName:          "test",
			MaxIdleConns:    100,
			MaxOpenConns:    100,
			ConnMaxIdletime: 30,
			ConnMaxLifetime: 60,
		}},
	}
	sc.initConnect()

	db := sc.GetDatabase("test")
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		if err := os.Remove("test.db"); err != nil {
			t.Error(err.Error())
		}
	}()

	err := db.AutoMigrate(&UserInfo{})
	if err != nil {
		t.Error(err.Error())
	}
	err = db.AutoMigrate(&UserInfo{})
	if err != nil {
		t.Error(err.Error())
	}

	user := UserInfo{
		Name:  "test",
		Email: "test@test.com",
	}

	result := db.Create(&user)
	if result.Error != nil {
		t.Errorf("db.Create(&user) err = %v", result.Error)
	} else {
		t.Log("User created successfully, ID:", user.ID)
	}

	var resUser UserInfo
	if err := db.First(&resUser, user.ID).Error; err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("select user %+v", resUser)

	db.Model(&resUser).Update("Email", "newemail@example.com")
	t.Logf("User updated, new email: %s\n", resUser.Email)

	db.Delete(&resUser)
	t.Log("User deleted successfully")
}
