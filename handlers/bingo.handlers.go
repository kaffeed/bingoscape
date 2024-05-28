package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

func (bh *BingoHandler) handleGetBingoParticipationTable(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok || !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, _ := c.Get(mgmnt_key).(bool)

	bingoId, err := strconv.Atoi(c.Param("bingoId"))

	if err != nil {
		return fmt.Errorf("No valid parameter bingoId: %w", err)
	}
	bingo, err := bh.BingoService.GetBingo(bingoId)

	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return fmt.Errorf("Could not get participants! %w", err)
	}

	pp, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return fmt.Errorf("Could not get possible participants! %w", err)
	}
	participantTable := components.BingoTeams(isManagement, bingo, p, pp)

	return render(c, participantTable)
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
		bingo := services.Bingo{
			Title:       c.FormValue("title"),
			From:        time.Now(),
			To:          time.Now().AddDate(0, 30, 0),
			Rows:        rows,
			Cols:        cols,
			Description: c.FormValue("Description"),
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

func (bh *BingoHandler) handleEditTile(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}

	tileId, err := strconv.Atoi(c.Param("tileId"))
	if err != nil {
		return fmt.Errorf("Need valid tile id: %w", err)
	}

	tile, err := bh.BingoService.LoadTile(tileId)
	fmt.Printf("Tile: %#v", tile)

	c.Set("ISERROR", false)

	if c.Request().Method == "PUT" {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination

		p := filepath.Join(os.Getenv("IMAGE_PATH"), "tiles")
		dst, err := os.CreateTemp(p, fmt.Sprintf("bingoscape-*%s", path.Ext(file.Filename)))

		log.Printf("Copying file to %s", dst.Name())
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		t := services.Tile{ // TODO: Read from config
			Id:          tileId,
			Title:       c.FormValue("title"),
			Description: c.FormValue("description"),
			ImagePath:   staticPath(dst.Name(), p),
			BingoId:     tile.BingoId,
		}

		err = bh.BingoService.UpdateTile(t)

		if c.FormValue("saveAsTemplate") == "on" {
			bh.BingoService.CreateTemplateTile(t)
		}

		if err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "Created a new bingo!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/bingos/%d", t.BingoId))
	}

	editTileView := authviews.Tile(isManagement, tile)

	return render(c, authviews.TileIndex(
		"| Tile",
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		editTileView,
	))
}

func staticPath(imgPath, basePath string) string {
	return strings.Replace(imgPath, basePath, "/img/tiles", -1)
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

		return c.Redirect(http.StatusSeeOther, "/")
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

type ParticipantId struct {
	id int `form:"possibleparticipants,omitempty"`
}

func (bh *BingoHandler) handleBingoParticipation(c echo.Context) error {
	bingoId, err := strconv.Atoi(c.Param("bingoId"))

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	log.Printf("ISAUTHENTICATED: %+v", isAuthenticated)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !isManagement {
		return err
	}
	bingo, err := bh.BingoService.GetBingo(bingoId)
	if err != nil {
		return err
	}

	if c.Request().Method == "POST" {
		pId, err := strconv.Atoi(c.FormValue("team"))
		if err != nil {
			log.Printf("Error during getting team id")
			return err
		}

		bh.BingoService.AddParticipantToBingo(pId, bingoId)
	}

	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return err
	}

	possible, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return err
	}

	bingoView := components.BingoTeams(isManagement, bingo, p, possible)
	return render(c, bingoView)
}

func (bh *BingoHandler) removeBingoParticipation(c echo.Context) error {
	bingoId, err := strconv.Atoi(c.Param("bingoId"))
	if err != nil {
		return errors.New("Can't parse bingoId from params")
	}
	pId, err := strconv.Atoi(c.Param("pId"))
	if err != nil {
		return errors.New("Can't parse participationId from params")
	}

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !isAuthenticated {
		return errors.New("Need to be authenticated")
	}
	log.Printf("Removing Bingo participant %d", pId)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok || !isManagement {
		log.Fatalf("Needs to be management!!")
		return errors.New("Needs to be management")
	}

	err = bh.BingoService.RemoveParticipation(pId, bingoId)
	if err != nil {
		return fmt.Errorf("Could not remove participation %d from bingo %d", pId, bingoId)
	}

	bingo, err := bh.BingoService.GetBingo(bingoId)

	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return fmt.Errorf("Could not get participants! %w", err)
	}

	pp, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return fmt.Errorf("Could not get possible participants! %w", err)
	}

	return render(c, components.BingoTeams(isManagement, bingo, p, pp))
}

func (bh *BingoHandler) handleGetBingoDetail(c echo.Context) error {
	bingoId, err := strconv.Atoi(c.Param("bingoId"))
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	log.Printf("ISAUTHENTICATED: %+v", isAuthenticated)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}
	bingo, err := bh.BingoService.GetBingo(bingoId)
	if err != nil {
		return err
	}
	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return err
	}

	possible, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return err
	}

	bingoView := authviews.BingoDetail(isManagement, bingo, p, possible)
	c.Set("ISERROR", false)

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

func (bh *BingoHandler) handleUpdateActiveState(c echo.Context) error {

	return nil
}
