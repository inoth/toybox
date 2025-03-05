package biz

type UserInfo struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRepo interface {
	CreateUser(user *UserInfo) (uint, error)
	GetUserById(id uint) (*UserInfo, error)
	GetAllUsers() ([]*UserInfo, error)
}

type UserUsecase struct {
	userRepo UserRepo
}

func NewUserUsecase(userRepo UserRepo) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) CreateUser(user *UserInfo) (uint, error) {
	return uc.userRepo.CreateUser(user)
}

func (uc *UserUsecase) GetUserById(id uint) (*UserInfo, error) {
	return uc.userRepo.GetUserById(id)
}

func (uc *UserUsecase) GetAllUsers() ([]*UserInfo, error) {
	return uc.userRepo.GetAllUsers()
}
