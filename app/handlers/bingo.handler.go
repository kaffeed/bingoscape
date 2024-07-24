package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kaffeed/bingoscape/app/db"
	"github.com/kaffeed/bingoscape/app/services"
	"github.com/kaffeed/bingoscape/app/views"
	authviews "github.com/kaffeed/bingoscape/app/views/auth"
	components "github.com/kaffeed/bingoscape/app/views/components"
	"github.com/labstack/echo/v4"
)

const (
	ISO8601 string = "2006-01-02T15:04"
)

type BingoHandler struct {
	BingoService *services.BingoService
	UserService  *services.UserService
	TileService  *services.TileService
}

func NewBingoHandler(bs *services.BingoService, us *services.UserService, ts *services.TileService) *BingoHandler {
	return &BingoHandler{
		BingoService: bs,
		UserService:  us,
		TileService:  ts,
	}
}

func (bh *BingoHandler) handleGetBingoParticipationTable(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok || !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, _ := c.Get(mgmnt_key).(bool)

	var bingoId int32
	echo.PathParamsBinder(c).Int32("bingoId", &bingoId)

	bingo, err := bh.BingoService.GetBingo(bingoId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	pp, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	participantTable := components.BingoTeams(isManagement, views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                []views.TileModel{},
		PossibleParticipants: pp,
		Participants:         p,
	})

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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("problem during bingo deletion %w", err))
	}

	return c.Redirect(http.StatusSeeOther, "/")
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

func (bh *BingoHandler) handleBingoParticipation(c echo.Context) error {
	var bingoId int32
	err := echo.PathParamsBinder(c).Int32("bingoId", &bingoId).BindError()

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	log.Printf("ISAUTHENTICATED: %+v", isAuthenticated)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/")
	}
	isManagement, _ := c.Get(mgmnt_key).(bool)
	if !isManagement {
		return err
	}
	bingo, err := bh.BingoService.GetBingo(bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if c.Request().Method == "POST" {
		var pId int32
		err := echo.FormFieldBinder(c).Int32("team", &pId).BindError()
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

	l, err := bh.BingoService.Store.GetBingoLeaderboard(context.TODO(), int32(bingoId))

	bingoView := components.BingoTeams(isManagement, views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                []views.TileModel{},
		PossibleParticipants: possible,
		Participants:         p,
		Leaderboard:          l,
	})
	c.Response().Header().Add("HX-Trigger", "updateLeaderboard")
	return render(c, bingoView)
}

func (bh *BingoHandler) removeBingoParticipation(c echo.Context) error {
	var bingoId, pId int32
	err := echo.PathParamsBinder(c).Int32("bingoId", &bingoId).Int32("pId", &pId).BindError()

	if err != nil {
		return echo.ErrBadRequest
	}

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !isAuthenticated {
		return errors.New("need to be authenticated")
	}
	log.Printf("Removing Bingo participant %d", pId)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok || !isManagement {
		log.Fatalf("Needs to be management!!")
		return errors.New("needs to be management")
	}

	err = bh.BingoService.RemoveParticipation(pId, bingoId)
	if err != nil {
		return fmt.Errorf("could not remove participation %d from bingo %d", pId, bingoId)
	}

	bingo, err := bh.BingoService.GetBingo(bingoId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("could not load bingo! %w", err))
	}

	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("could not get participants! %w", err))
	}

	pp, err := bh.BingoService.GetPossibleParticipants(bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("could not get participants! %w", err))
	}

	return render(c, components.BingoTeams(isManagement, views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                []views.TileModel{},
		PossibleParticipants: pp,
		Participants:         p,
	}))
}

func (bh *BingoHandler) handleGetBingoBoard(c echo.Context) error {
	var bingoId, forUser int32
	err := echo.
		PathParamsBinder(c).
		Int32("bingoId", &bingoId).
		BindError()

	if err != nil {
		return echo.ErrBadRequest
	}

	err = echo.QueryParamsBinder(c).Int32("forUser", &forUser).BindError()

	if err != nil {
		return echo.ErrBadRequest
	}

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	log.Printf("ISAUTHENTICATED: %+v", isAuthenticated)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}

	uid := c.Get(user_id_key).(int32)
	isDifferentUser := uid != forUser

	if !isManagement && isDifferentUser {
		return echo.ErrForbidden
	}

	bingo, err := bh.BingoService.GetBingo(bingoId)
	if err != nil {
		return err
	}
	p, err := bh.BingoService.GetParticipants(bingoId)
	if err != nil {
		return err
	}

	// possible, err := bh.BingoService.GetPossibleParticipants(bingoId)
	// if err != nil {
	// 	return err
	// }

	tiles, err := bh.TileService.LoadTilesForBingo(bingoId)
	if err != nil {
		return echo.ErrNotFound
	}

	bm := views.BingoDetailModel{
		Bingo:        bingo,
		Tiles:        tiles,
		Participants: p,
	}
	board := authviews.BingoBoard(isManagement, isDifferentUser, bm, forUser)
	return render(c, board)
}

func (bh *BingoHandler) handleGetBingoDetail(c echo.Context) error {
	var bingoId int32
	err := echo.
		PathParamsBinder(c).
		Int32("bingoId", &bingoId).
		BindError()

	if err != nil {
		return echo.ErrBadRequest
	}

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

	tiles, err := bh.TileService.LoadTilesForBingo(bingoId)
	if err != nil {
		return echo.ErrNotFound
	}

	l, _ := bh.BingoService.Store.GetBingoLeaderboard(context.TODO(), bingo.ID)
	uid := c.Get(user_id_key).(int32)
	bm := views.BingoDetailModel{
		Bingo:                bingo,
		Tiles:                tiles,
		PossibleParticipants: possible,
		Participants:         p,
		Leaderboard:          l,
	}

	bingoView := authviews.BingoDetail(isManagement, false, bm, uid)
	c.Set("ISERROR", false)

	c.Response().Header().Add("HX-Trigger", "updateLeaderboard")

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

	return render(c, authviews.BingoVisibilityButton(bingoId, bingoReady))
}

func (bh *BingoHandler) handleBingoToggleSubmissionClosed(c echo.Context) error {

	isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool)
	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, _ := c.Get(mgmnt_key).(bool)

	if !isManagement {
		return c.Redirect(http.StatusForbidden, c.Request().URL.RequestURI()) // FIXME: is this the right way?
	}

	var bingoId int32
	err := echo.PathParamsBinder(c).Int32("bingoId", &bingoId).BindError()

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	submissionsClosed, err := bh.BingoService.Store.ToggleSubmissionsClosed(context.Background(), bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return render(c, authviews.BingoSubmissionsClosedButton(bingoId, submissionsClosed))
}

func (bh *BingoHandler) handleBingoToggleLeaderboardPublic(c echo.Context) error {

	isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool)
	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, _ := c.Get(mgmnt_key).(bool)

	if !isManagement {
		return c.Redirect(http.StatusForbidden, c.Request().URL.RequestURI()) // FIXME: is this the right way?
	}

	var bingoId int32
	err := echo.PathParamsBinder(c).Int32("bingoId", &bingoId).BindError()

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	leaderboardPublic, err := bh.BingoService.Store.ToggleLeaderboardPublic(context.Background(), bingoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return render(c, authviews.BingoLeaderboardPublicButton(bingoId, leaderboardPublic))
}

func (bh *BingoHandler) handleTeamSubmissions(c echo.Context) error {
	isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool)
	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	isManagement, _ := c.Get(mgmnt_key).(bool)

	if !isManagement {
		return echo.ErrForbidden // FIXME: is this the right way?
	}

	var bingoId, loginId int32
	err := echo.
		PathParamsBinder(c).
		Int32("bingoId", &bingoId).
		Int32("loginId", &loginId).
		BindError()

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	submissions, err := bh.
		BingoService.
		Store.
		GetSubmissionsByBingoAndLogin(context.TODO(), db.GetSubmissionsByBingoAndLoginParams{
			LoginID: loginId,
			BingoID: bingoId,
		})

	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "could not load submissons")
	}

	sd := make([]views.SubmissionData, len(submissions))
	for i, s := range submissions {
		sc, err := bh.BingoService.Store.GetCommentsForSubmission(context.Background(), s.Submission.ID)
		if err != nil {
			c.Set("ISERROR", true)
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, "could not load submisson comments")
		}

		ip, err := bh.BingoService.Store.GetImagesForSubmission(context.Background(), s.Submission.ID)
		if err != nil {
			c.Set("ISERROR", true)
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, "could not load submisson images")
		}

		sd[i] = views.SubmissionData{
			Submission: s.Submission,
			Comments:   sc,
			Images:     ip,
			Tile:       s.Tile,
		}
	}

	l, err := bh.BingoService.Store.GetLoginById(context.TODO(), loginId)
	if err != nil {
		c.Set("ISERROR", true)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "no login found")
	}
	sc, err := bh.BingoService.Store.GetSubmissionClosedStatusForBingo(context.TODO(), bingoId)

	m := views.TeamSubmissionModel{
		Submissions:       sd,
		BingoID:           bingoId,
		Name:              l.Name,
		SubmissionsClosed: sc,
	}

	submissionView := authviews.TeamSubmissions(isManagement, m)

	c.Set("ISERROR", false)

	return render(c, authviews.TeamSubmissionsIndex(
		"| Team Submissions",
		"", // TODO: set someday
		isAuthenticated,
		isManagement,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		submissionView,
	))
}
