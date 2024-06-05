package services

import (
	"context"

	"github.com/kaffeed/bingoscape/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(store db.Queries) *UserService {
	return &UserService{
		UserStore: store,
	}
}

type UserService struct {
	UserStore db.Queries
}

func (us *UserService) GetAllUsers() ([]db.Login, error) {
	return us.UserStore.GetAllLogins(context.Background())
}

func (us *UserService) CreateUser(u db.CreateLoginParams) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return us.UserStore.CreateLogin(context.Background(), u)
}

func (us *UserService) UpdatePassword(uid int, p string) error {
	_ = `SELECT id, password, name, is_management FROM logins`

	return nil
}

func (us *UserService) DeleteUser(uid int) error {
	_ = `SELECT id, password, name, is_management FROM logins`

	return nil
}

func (us *UserService) CheckUsername(username string) (User, error) {

	query := `SELECT id, password, name, is_management FROM logins
		WHERE name = $1`

	stmt, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	u := User{}
	err = stmt.QueryRow(
		username,
	).Scan(
		&u.Id,
		&u.Password,
		&u.Username,
		&u.IsManagement,
	)

	if err != nil {
		return User{}, err
	}

	return u, nil // TODO:
}
