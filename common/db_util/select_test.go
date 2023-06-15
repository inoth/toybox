package dbutil

import (
	"fmt"
	"testing"
)

type User struct {
	Id   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func TestSelect(t *testing.T) {
	user, err := Select(nil, User{Id: "1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	println(user.Name)
	t.Log()
}

func TestSelectList(t *testing.T) {
	users, err := SelectList(nil, User{Id: "1"}, Paginate(1, 10))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", users)
	t.Log()
}
