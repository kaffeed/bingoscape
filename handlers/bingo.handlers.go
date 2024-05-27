package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (bh *BingoHandler) handleCreateBingo(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		return fmt.Errorf("Invalid type for key '" + mgmnt_key + "'")
	}

	createView := authviews.CreateBingo(isManagement)
	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		rows, _ := strconv.Atoi(c.FormValue("rows"))
		cols, _ := strconv.Atoi(c.FormValue("cols"))
		bingo := services.Bingo{ // TODO: Read from config
			Title: c.FormValue("title"),
			From:  time.Now(),
			To:    time.Now().AddDate(0, 30, 0),
			Rows:  rows,
			Cols:  cols,
		}

		bingo, err := bh.BingoService.CreateBingo(bingo)
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

		for i := 0; i < bingo.Rows*bingo.Cols; i++ {
		}

		setFlashmessages(c, "success", "Created a new bingo!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/bingos/%d", bingo.Id))
	}

	return render(c, authviews.RegisterIndex(
		"| Register",
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		createView,
	))
}

func (bh *BingoHandler) handleDeleteBingo(c echo.Context) error {
	if isManagement, _ := c.Get(mgmnt_key).(bool); !isManagement {
		setFlashmessages(c, "error", "Only management accounts can delete bingos!")
		return c.Redirect(http.StatusUnauthorized, "/")
	}
	if isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool); !isAuthenticated {
		setFlashmessages(c, "error", "You're not authorized to delete bingos")
		return c.Redirect(http.StatusUnauthorized, "/")
	}

	query := `DELETE FROM bingos WHERE id = $1`
	stmt, err := bh.UserService.UserStore.Db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	bId, err := strconv.Atoi(c.Param("bingoId"))

	if _, err = stmt.Exec(bId); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
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
	registerView := authviews.Register(isManagement)
	// isError = false
	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		user := services.User{
			Password:     c.FormValue("password"),
			Username:     c.FormValue("username"),
			IsManagement: c.FormValue("IsManagement") == "on",
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

func (bh *BingoHandler) handleGetBingoDetail(c echo.Context) error {
	bingoId, err := strconv.Atoi(c.Param("bingoId"))
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	log.Printf("ISAUTHENTICATED: %+v", isAuthenticated)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, _ := c.Get(mgmnt_key).(bool)
	bingo, err := bh.BingoService.GetBingo(bingoId)
	if err != nil {
		return err
	}
	bingoView := authviews.BingoDetail(isAuthenticated, bingo)
	c.Set("ISERROR", false)
	// fmt.Printf("\033[31mFROMPROTECTED = %t\n\033[0m", fromProtected)

	return render(c, authviews.BingoDetailIndex(
		"| Bingo",
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		bingoView,
	))
}
