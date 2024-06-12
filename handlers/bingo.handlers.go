package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kaffeed/bingoscape/db"
	"github.com/kaffeed/bingoscape/internal/util"
	"github.com/kaffeed/bingoscape/services"
	"github.com/kaffeed/bingoscape/views"
	authviews "github.com/kaffeed/bingoscape/views/auth"
	components "github.com/kaffeed/bingoscape/views/components"
	"github.com/labstack/echo/v4"
)

const (
	ISO8601 string = "2006-01-02T15:04"
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

	var bingoId int
	echo.PathParamsBinder(c).Int("bingoId", &bingoId)

	bingo, err := bh.BingoService.GetBingo(bingoId)

	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return fmt.Errorf("Could not get participants! %w", err)
	}

	pp, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return fmt.Errorf("Could not get possible participants! %w", err)
	}
	participantTable := components.BingoTeams(isManagement, views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                []views.TileModel{},
		PossibleParticipants: pp,
		Participants:         p,
	})

	return render(c, participantTable)
}

func (bh *BingoHandler) handleLoadFromTemplate(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	if !isAuthenticated {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Need to be authenticated")
	}
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		return fmt.Errorf("Invalid type for key '" + mgmnt_key + "'")
	}

	if !isManagement {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Need to be management for this")
	}

	var tileId int32

	err := echo.
		PathParamsBinder(c).
		Int32("tileId", &tileId).
		BindError()

	if err != nil {
		return fmt.Errorf("Need valid tile id: %w", err)
	}

	var templateId int32
	err = echo.QueryParamsBinder(c).Int32("templateId", &templateId).BindError()

	if err != nil {
		return fmt.Errorf("Need valid template id: %w", err)
	}

	log.Printf("####################################################################################################")
	log.Printf("can load queryparams")
	log.Printf("####################################################################################################")

	tile, err := bh.BingoService.LoadTile(int(tileId))
	if err != nil {
		return err
	}
	fmt.Printf("Tile: %#v", tile)

	uid, ok := c.Get(user_id_key).(int32)
	var templates []db.TemplateTile
	var submissions views.Submissions

	if isManagement {
		submissions, err = bh.BingoService.LoadAllSubmissionsForTile(tile.ID)
		templates, err = bh.BingoService.Store.GetTemplateTiles(context.Background())
	} else {
		submissions, err = bh.BingoService.LoadUserSubmissions(tile.ID, uid)
	}
	templates, err = bh.BingoService.Store.GetTemplateTiles(context.Background())
	if err != nil {
		return err
	}

	log.Printf("####################################################################################################")
	log.Printf("can load templates")
	log.Printf("####################################################################################################")

	idx := slices.IndexFunc(templates, func(myTempl db.TemplateTile) bool {
		return myTempl.ID == templateId
	})

	if idx != -1 {
		templateTile := templates[idx]
		tile.Title = templateTile.Title
		tile.Description = templateTile.Description
		tile.Imagepath = templateTile.Imagepath

		log.Printf("####################################################################################################")
		log.Printf("Could update from template")
		log.Printf("####################################################################################################")
	}

	tm := views.TileModel{
		Tile:        tile,
		Templates:   templates,
		Submissions: submissions,
	}

	editForm := authviews.Tile(isManagement, tm, uid)
	return render(c, editForm) // TODO: Change to tileEditForm
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

	tzone, _ := c.Get(tzone_key).(string)
	_, err := time.LoadLocation(tzone)
	if err != nil {
		setFlashmessages(c, "error", "Could not load timezone information!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/"))
	}

	createView := authviews.CreateBingo(isManagement)
	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		bingo := db.CreateBingoParams{}

		var validFrom, validTo time.Time
		err := echo.FormFieldBinder(c).
			String("title", &bingo.Title).
			Int32("rows", &bingo.Rows).
			Int32("cols", &bingo.Cols).
			Time("validfrom", &validFrom, ISO8601).
			Time("validto", &validTo, ISO8601).
			String("description", &bingo.Description).
			String("codephrase", &bingo.Codephrase).
			BindError()

		bingo.Validfrom = pgtype.Timestamp{
			Time:             validFrom,
			InfinityModifier: 0,
			Valid:            true,
		}

		bingo.Validto = pgtype.Timestamp{
			Time:             validTo,
			InfinityModifier: 0,
			Valid:            true,
		}
		if err != nil {
			return err
		}

		viewModel, err := bh.BingoService.CreateBingo(bingo)
		if err != nil {

			setFlashmessages(c, "error", fmt.Sprintf(
				"something went wrong: %s",
				err,
			))

			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "Created a new bingo!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/bingos/%d", viewModel.Bingo.ID))
	}

	return render(c, authviews.CreateUserIndex( // FIXME: Hmm... Really CreateUserIndex?
		"| Create Login",
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
	var bingoId int
	echo.PathParamsBinder(c).Int("bingoId", &bingoId)

	err := bh.BingoService.Store.DeleteBingoById(context.Background(), int32(bingoId))
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (bh *BingoHandler) handlePutSubmissionStatus(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok || !isManagement {
		return c.Redirect(http.StatusUnauthorized, c.Path())
	}

	var submissionId int
	var parsedState string

	var comments []string
	echo.FormFieldBinder(c).
		Strings("comment", &comments) // FIXME: Why the fuck is that multiple strings
	comment := strings.Join(comments, "\n")

	err := echo.PathParamsBinder(c).
		Int("submissionId", &submissionId).
		String("state", &parsedState).
		BindError()

	if err != nil {
		return err
	}

	state := db.Submissionstate(parsedState)

	s, err := bh.BingoService.Store.UpdateSubmissionState(context.Background(), db.UpdateSubmissionStateParams{
		ID:    int32(submissionId),
		State: state,
	})

	if err != nil {
		c.Set("ISERROR", true)
		setFlashmessages(c, "error", "Could not update submission state")
		return err
	}

	if comment != "" {
		uid, ok := c.Get(user_id_key).(int32)

		if ok {
			err = bh.BingoService.Store.CreateSubmissionComment(context.Background(), db.CreateSubmissionCommentParams{
				SubmissionID: int32(submissionId),
				LoginID:      int32(uid),
				Comment:      comment,
			})
			if err != nil {
				return err
			}
		}
	}

	sc, err := bh.BingoService.Store.GetCommentsForSubmission(context.Background(), s.ID)
	if err != nil {
		c.Set("ISERROR", true)
		setFlashmessages(c, "error", "Could not update submission state")
		return err
	}

	ip, err := bh.BingoService.Store.GetImagesForSubmission(context.Background(), s.ID)
	if err != nil {
		c.Set("ISERROR", true)
		setFlashmessages(c, "error", "Could not update submission state")
		return err
	}

	c.Set("ISERROR", false)
	loc := getLocation(c)

	data := views.SubmissionData{
		Submission: s,
		Comments:   sc,
		Images:     ip,
	}

	return render(c, components.SubmissionHeader(isManagement, data, loc))
}

func (bh *BingoHandler) handleTileSubmission(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	if !isAuthenticated {
		c.Redirect(http.StatusUnauthorized, "/login")
	}

	var tileId int
	echo.PathParamsBinder(c).Int("tileId", &tileId)

	tile, err := bh.BingoService.LoadTile(tileId)
	if err != nil {
		return err
	}
	fmt.Printf("Tile: %#v", tile)

	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		form, err := c.MultipartForm()

		if err != nil {
			return err
		}
		files := form.File["files"]
		filePaths := []string{}
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			// Destination

			p := filepath.Join(os.Getenv("IMAGE_PATH"), "submissions")
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
			filePaths = append(filePaths, util.StaticPath(dst.Name()))
		}

		u := c.Get(user_id_key).(int32)
		err = bh.BingoService.CreateSubmission(tileId, u, filePaths)

		if err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		setFlashmessages(c, "success", "Submission successful!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tile.ID))
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tile.ID))
}

func (bh *BingoHandler) handleGetTileSubmissions(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	if !isAuthenticated {
		return errors.New("Need to be authenticated")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}

	uid := c.Get(user_id_key).(int32)

	var tileId int32
	err := echo.PathParamsBinder(c).Int32("tileId", &tileId).BindError()
	if err != nil {
		return fmt.Errorf("Need valid tile id: %w", err)
	}

	var submissionMap views.Submissions
	if isManagement {
		submissionMap, err = bh.BingoService.LoadAllSubmissionsForTile(tileId)
	} else {
		submissionMap, err = bh.BingoService.LoadUserSubmissions(tileId, uid)
	}

	if err != nil {
		return err
	}
	loc := getLocation(c)

	return render(c, components.Submissions(isManagement, submissionMap, loc))
}

func (bh *BingoHandler) handleTile(c echo.Context) error {
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

		var f string

		if err != nil {
			var imagePath string
			err := echo.FormFieldBinder(c).String("imagepath", &imagePath).BindError()
			if err != nil {
				f = tile.Imagepath
			} else {
				f = imagePath
			}
		} else {
			f, err = util.SaveFile(file)

			if err != nil {
				f = tile.Imagepath
			}
		}

		t := db.UpdateTileParams{ // TODO: Read from config
			ID:          int32(tileId),
			Title:       c.FormValue("title"),
			Description: c.FormValue("description"),
			Imagepath:   f,
		}

		_, err = bh.BingoService.Store.UpdateTile(context.Background(), t)

		if err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		if c.FormValue("saveAsTemplate") == "on" {
			_, _ = bh.BingoService.Store.CreateTemplateTile(context.Background(), db.CreateTemplateTileParams{
				Title:       t.Title,
				Imagepath:   t.Imagepath,
				Description: t.Description,
			})
		}

		setFlashmessages(c, "success", "Updated tile!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/bingos/%d", tile.BingoID))
	}

	uid := c.Get(user_id_key).(int32)
	var submissions views.Submissions
	var templates []db.TemplateTile

	if isManagement {
		submissions, err = bh.BingoService.LoadAllSubmissionsForTile(tile.ID)
		templates, err = bh.BingoService.Store.GetTemplateTiles(context.Background())
	} else {
		submissions, err = bh.BingoService.LoadUserSubmissions(tile.ID, uid)
	}

	tm := views.TileModel{
		Tile:        tile,
		Submissions: submissions,
		Templates:   templates,
	}

	editTileView := authviews.Tile(isManagement, tm, uid)

	return render(c, authviews.TileIndex(
		"| Tile",
		"", // TODO: set someday
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		editTileView,
	))
}

func (bh *BingoHandler) CreateLoginHandler(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		return fmt.Errorf("Invalid type for key '" + mgmnt_key + "'")
	}
	registerView := authviews.CreateUser(isManagement)
	// isError = false
	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		user := db.CreateLoginParams{
			Password:     c.FormValue("password"),
			Name:         c.FormValue("username"),
			IsManagement: c.FormValue("management") == "on",
		}

		err := bh.UserService.CreateUser(user)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the username is already in use")
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

	return render(c, authviews.CreateUserIndex(
		"| Create User",
		"",
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		registerView,
	))
}

func getLocation(c echo.Context) *time.Location {
	locName, ok := c.Get(tzone_key).(string)
	if !ok {
		return time.Now().Location()
	}

	loc, err := time.LoadLocation(locName)
	if err != nil {
		return time.Now().Location()
	}
	return loc
}

func (bh *BingoHandler) handleGetAllBingos(c echo.Context) error {
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}
	u := c.Get(user_id_key).(int32)
	loc := getLocation(c)
	bingos, _ := bh.BingoService.GetBingos(isManagement, u)
	bingoTable := components.BingoTable(isManagement, bingos, loc)

	return render(c, bingoTable)
}

type ParticipantId struct {
	id int `form:"possibleparticipants,omitempty"`
}

func (bh *BingoHandler) handleBingoParticipation(c echo.Context) error {
	var bingoId int
	err := echo.PathParamsBinder(c).Int("bingoId", &bingoId).BindError()

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	log.Printf("ISAUTHENTICATED: %+v", isAuthenticated)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/")
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
		var pId int
		err := echo.FormFieldBinder(c).Int("team", &pId).BindError()
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

	bingoView := components.BingoTeams(isManagement, views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                []views.TileModel{},
		PossibleParticipants: possible,
		Participants:         p,
	})
	return render(c, bingoView)
}

func (bh *BingoHandler) removeBingoParticipation(c echo.Context) error {
	var bingoId int
	err := echo.PathParamsBinder(c).Int("bingoId", &bingoId).BindError()

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

	return render(c, components.BingoTeams(isManagement, views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                []views.TileModel{},
		PossibleParticipants: pp,
		Participants:         p,
	}))
}

func (bh *BingoHandler) handleGetBingoDetail(c echo.Context) error {
	var bingoId int
	err := echo.PathParamsBinder(c).Int("bingoId", &bingoId).BindError()
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

	tiles, err := bh.BingoService.LoadTilesForBingo(bingoId)

	uid := c.Get(user_id_key).(int32)
	bm := views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                tiles,
		PossibleParticipants: possible,
		Participants:         p,
	}

	bingoView := authviews.BingoDetail(isManagement, bm, uid)
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

func (bh *BingoHandler) handleBingoState(c echo.Context) error {

	isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool)
	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, _ := c.Get(mgmnt_key).(bool)

	if !isManagement {
		return c.Redirect(http.StatusUnauthorized, c.Request().URL.RequestURI()) // FIXME: is this the right way?
	}

	var bingoId int32
	err := echo.PathParamsBinder(c).Int32("bingoId", &bingoId).BindError()

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	bingoReady, err := bh.BingoService.Store.ToggleBingoState(context.Background(), bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return render(c, authviews.BingoStateButton(bingoId, bingoReady))
}
