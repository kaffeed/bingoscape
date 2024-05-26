package services

import (
	"github.com/kaffeed/bingoscape/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(store db.Store) *UserService {
	return &UserService{
		UserStore: store,
	}
}

type User struct {
	Id           int    `json:"id,omitempty"`
	Password     string `json:"password,omitempty"`
	Username     string `json:"username,omitempty"`
	IsManagement bool   `json:"is_management,omitempty"`
}

type UserService struct {
	UserStore db.Store
}

func (us *UserService) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(password, username, ismanagement) VALUES($1, $2, $3)`

	_, err = us.UserStore.Db.Exec(
		stmt,
		string(hashedPassword),
		u.Username,
		u.IsManagement,
	)

	return err
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
