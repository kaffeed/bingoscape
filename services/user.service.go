package services

import (
	"github.com/kaffeed/topez-bingomania/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(u User, store db.Store) *UserServices {
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

	stmt := `INSERT INTO users(password, username) VALUES($1, $2, $3)`

	_, err = us.UserStore.Db.Exec(
		stmt,
		string(hashedPassword),
		u.Username,
	)

	return err
}

// func (us *UserServices) CheckEmail(email string) (User, error) {
//
// 	query := `SELECT id, email, password, username FROM users
// 		WHERE email = ?`
//
// 	stmt, err := us.UserStore.Db.Prepare(query)
// 	if err != nil {
// 		return User{}, err
// 	}
//
// 	defer stmt.Close()
//
// 	us.User.Email = email
// 	err = stmt.QueryRow(
// 		us.User.Email,
// 	).Scan(
// 		&us.User.ID,
// 		&us.User.Email,
// 		&us.User.Password,
// 		&us.User.Username,
// 	)
// 	if err != nil {
// 		return User{}, err
// 	}
//
// 	return us.User, nil
// }
