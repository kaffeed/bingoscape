package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kaffeed/bingoscape/services"
	authviews "github.com/kaffeed/bingoscape/views/auth"
	components "github.com/kaffeed/bingoscape/views/components"
	"github.com/labstack/echo/v4"
)

type BingoHandler struct {
	BingoService *services.BingoService
	UserService  *services.UserService
}

func NewBingoHandler(bs *services.BingoService, us *services.UserService) *BingoHandler {
	return &BingoHandler{
		BingoService: bs,
		UserService:  us,
	}
}

func (bh *BingoHandler) handleGetActiveBingo(c echo.Context) error {
	if c.Request().Method != "GET" {
		return errors.New("Invalid http method!")
	}

	return nil
}

func (bh *BingoHandler) RegisterHandler(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		return fmt.Errorf("Invalid type for key '" + mgmnt_key + "'")
	}
	canCreateTeams := isAuthenticated && isManagement
	registerView := authviews.Register(canCreateTeams)
	// isError = false
	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		user := services.User{
			Password: c.FormValue("password"),
			Username: c.FormValue("username"),
		}

		err := bh.UserService.CreateUser(user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the email is already in use")
				setFlashmessages(c, "error", fmt.Sprintf(
					"something went wrong: %s",
					err,
				))

				return c.Redirect(http.StatusSeeOther, "/register")
			}

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "You have successfully registered!!")

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return render(c, authviews.RegisterIndex(
		"| Register",
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		registerView,
	))
}

func (bh *BingoHandler) handleGetAllBingos(c echo.Context) error {
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}
	u := c.Get(user_id_key).(int)
	bingos, _ := bh.BingoService.GetBingos(isManagement, u)
	bingoTable := components.BingoTable(isManagement, bingos)

	return render(c, bingoTable)
}
