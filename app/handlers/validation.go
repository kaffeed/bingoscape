package handlers

import (
	"crypto/subtle"

	"github.com/kaffeed/bingoscape/app/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func basicAuthValidatorFunc(us *services.UserService) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		user, err := us.CheckUsername(username)
		if err != nil {
			return false, err
		}

		if subtle.ConstantTimeCompare([]byte(password), []byte(user.Password)) == 1 {
			c.Set(user_id_key, user.ID)
			return true, nil
		}

		return false, nil
	}
}
