package user

import (
	"demo-wire/internal/biz"
	"demo-wire/internal/data/user/model"

	database "github.com/inoth/toybox/component/database/sqlite"
	"github.com/inoth/toybox/util/convert"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *database.SqliteComponent) biz.UserRepo {
	user := db.GetDatabase("user", &model.UserInfo{})
	return &userRepo{db: user}
}

func (r *userRepo) CreateUser(user *biz.UserInfo) (uint, error) {
	add, err := convert.ConvertByUnsafe[biz.UserInfo, model.UserInfo](user)
	if err != nil {
		return 0, err
	}
	if err := r.db.Create(add).Error; err != nil {
		return 0, err
	}
	return add.Id, nil
}

func (r *userRepo) GetUserById(id uint) (*biz.UserInfo, error) {
	var user model.UserInfo
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return convert.ConvertByUnsafe[model.UserInfo, biz.UserInfo](&user)
}

func (r *userRepo) GetAllUsers() ([]*biz.UserInfo, error) {
	var users []model.UserInfo
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	res := make([]*biz.UserInfo, 0, len(users))
	for _, user := range users {
		r, err := convert.ConvertByUnsafe[model.UserInfo, biz.UserInfo](&user)
		if err != nil {
			continue
		}
		res = append(res, r)
	}
	return res, nil
}
