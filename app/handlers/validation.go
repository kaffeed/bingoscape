package handlers

import (
	"github.com/kaffeed/bingoscape/app/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

func basicAuthValidatorFunc(us *services.UserService) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		user, err := us.CheckUsername(username)
		if err != nil {
			return false, err
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(password),
		)

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return false, nil
		}

		c.Set(user_key, user)
		return true, nil
	}
}
