package services

import (
	"context"

	"github.com/kaffeed/bingoscape/app/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(store *db.Queries) *UserService {
	return &UserService{
		UserStore: store,
	}
}

type UserService struct {
	UserStore *db.Queries
}

func (us *UserService) GetAllUsers() ([]db.Login, error) {
	return us.UserStore.GetAllLogins(context.Background())
}

func (us *UserService) CreateUser(params db.CreateLoginParams) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 8)
	if err != nil {
		return err
	}

	params.Password = string(hashedPassword)
	return us.UserStore.CreateLogin(context.Background(), params)
}

func (us *UserService) UpdatePassword(uid int32, p string) (db.Login, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), 8)
	if err != nil {
		return db.Login{}, err
	}
	return us.UserStore.UpdateLoginPassword(context.Background(), db.UpdateLoginPasswordParams{
		ID:       uid,
		Password: string(hashedPassword),
	})
}

func (us *UserService) DeleteUser(uid int32) error {
	return us.UserStore.DeleteLogin(context.Background(), uid)
}

func (us *UserService) CheckUsername(username string) (db.Login, error) {
	return us.UserStore.GetLoginByName(context.Background(), username)
}
