package services

import (
	"github.com/kaffeed/bingoscape/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(store db.Store) *UserServices {
	return &UserServices{
		UserStore: store,
	}
}

type User struct {
	Id           int    `json:"id,omitempty"`
	Password     string `json:"password,omitempty"`
	Username     string `json:"username,omitempty"`
	IsManagement bool   `json:"is_management,omitempty"`
}

type UserServices struct {
	User      User
	UserStore db.Store
}

func (us *UserServices) CreateUser(u User) error {
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

func (us *UserServices) CheckUsername(username string) (User, error) {

	query := `SELECT id, password, username FROM users
		WHERE username = ?`

	stmt, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	// us.User.Email = email
	// err = stmt.QueryRow(
	// 	us.User.Email,
	// ).Scan(
	// 	&us.User.ID,
	// 	&us.User.Email,
	// 	&us.User.Password,
	// 	&us.User.Username,
	// )

	// if err != nil {
	// 	return User{}, err
	// }

	return User{}, nil // TODO:
}
