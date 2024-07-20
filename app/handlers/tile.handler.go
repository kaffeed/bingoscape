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
	"strings"

	"github.com/kaffeed/bingoscape/app/db"
	"github.com/kaffeed/bingoscape/app/internal/util"
	"github.com/kaffeed/bingoscape/app/services"
	"github.com/kaffeed/bingoscape/app/views"
	authviews "github.com/kaffeed/bingoscape/app/views/auth"
	components "github.com/kaffeed/bingoscape/app/views/components"
	"github.com/labstack/echo/v4"
)

type TileHandler struct {
	UserService *services.UserService
	TileService *services.TileService
}

func NewTileHandler(ts *services.TileService, us *services.UserService) *TileHandler {
	return &TileHandler{
		TileService: ts,
		UserService: us,
	}
}

func handleFileUpload(c echo.Context, fn, hfn, tip string) (string, error) {
	file, err := c.FormFile(fn)

	var f string

	if err != nil {
		var imagePath string
		err := echo.
			FormFieldBinder(c).
			String(hfn, &imagePath).
			BindError()
		if !util.IsEmptyOrWhitespace(tip) || err != nil {
			f = tip
		} else {
			f = imagePath
		}
	} else {
		f, err = util.SaveFile(file)
		if err != nil {
			return "", err
		}
	}

	return f, nil
}

func (th *TileHandler) handleDeleteTemplate(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	if !isAuthenticated {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Need to be authenticated")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Errorf("invalid type for key '"+mgmnt_key+"'"))
	}

	if !isManagement {
		return echo.NewHTTPError(echo.ErrForbidden.Code, fmt.Errorf("management only"))
	}

	var tmplId int32
	err := echo.PathParamsBinder(c).Int32("templateId", &tmplId).BindError()
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err)
	}

	err = th.TileService.Store.DeleteTemplateById(context.TODO(), tmplId)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.Redirect(http.StatusSeeOther, "/tiles/templates")
}

func (th *TileHandler) handleLoadFromTemplate(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	if !isAuthenticated {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Need to be authenticated")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, fmt.Errorf("invalid type for key '"+mgmnt_key+"'"))
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
		return fmt.Errorf("need valid tile id: %w", err)
	}

	var templateId int32
	err = echo.QueryParamsBinder(c).Int32("templateId", &templateId).BindError()

	if err != nil {
		return fmt.Errorf("need valid template id: %w", err)
	}

	tile, err := th.TileService.LoadTile(int(tileId))
	if err != nil {
		return err
	}
	fmt.Printf("Tile: %#v", tile)

	uid, _ := c.Get(user_id_key).(int32)
	var templates []db.TemplateTile
	var submissions views.Submissions

	if isManagement {
		submissions, _ = th.TileService.LoadAllSubmissionsForTile(tile.ID)
	} else {
		submissions, _ = th.TileService.LoadUserSubmissions(tile.ID, uid)
	}
	templates, err = th.TileService.Store.GetTemplateTiles(context.Background())
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(templates, func(myTempl db.TemplateTile) bool {
		return myTempl.ID == templateId
	})

	if idx != -1 {
		templateTile := templates[idx]
		tile.Title = templateTile.Title
		tile.Description = templateTile.Description
		tile.Imagepath = templateTile.Imagepath
		tile.Weight = templateTile.Weight
		tile.SecondaryImagePath = templateTile.SecondaryImagePath
	}

	tm := views.TileModel{
		Tile:        tile,
		Templates:   templates,
		Submissions: submissions,
	}

	editForm := authviews.Tile(isManagement, tm, uid)
	return render(c, editForm) // TODO: Change to tileEditForm
}

func (th *TileHandler) handleGetTemplateTiles(c echo.Context) error {
	if isManagement, _ := c.Get(mgmnt_key).(bool); !isManagement {
		setFlashmessages(c, "error", "Only management accounts can view and edit template tiles!")
		return c.Redirect(http.StatusUnauthorized, "/")
	}

	if isAuthenticated, _ := c.Get("ISAUTHENTICATED").(bool); !isAuthenticated {
		setFlashmessages(c, "error", "Log in to view and edit template tiles")
		return c.Redirect(http.StatusUnauthorized, "/")
	}

	tt, err := th.TileService.Store.GetTemplateTiles(context.TODO())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "problem while loading template tiles")
	}

	tv := authviews.TemplateTiles(tt)

	c.Set("ISERROR", false)

	return render(c, authviews.TemplateTilesIndex(
		"| Template Tiles",
		"", // TODO: set someday
		true,
		true,
		c.Get("ISERROR").(bool),
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		tv,
	))
}

func (th *TileHandler) handleDeleteSubmission(c echo.Context) error {
	isManagement, _ := c.Get(mgmnt_key).(bool)

	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	if !isAuthenticated {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	var tileId, submissionId int32
	err := echo.PathParamsBinder(c).
		Int32("tileId", &tileId).
		Int32("submissionId", &submissionId).
		BindError()

	s, err := th.TileService.Store.GetSubmissionById(context.TODO(), submissionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	uid, _ := c.Get(user_id_key).(int32)
	if !isManagement && s.LoginID != uid {
		setFlashmessages(c, "error", "can't delete another teams submission!")
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tileId))
	}

	t, err := th.TileService.LoadTile(int(tileId))
	if err != nil {
		return echo.ErrInternalServerError
	}

	sc, err := th.TileService.Store.GetSubmissionClosedStatusForBingo(context.TODO(), t.BingoID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if sc {
		setFlashmessages(c, "error", "Can't delete submissions after submissions are closed")
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tileId))
	}

	err = th.TileService.Store.DeleteSubmissionById(context.TODO(), s.ID)
	if err != nil {
		setFlashmessages(c, "error", "Could not delete submission")
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tileId))
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tileId))
}

func (th *TileHandler) handlePutSubmissionStatus(c echo.Context) error {
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

	s, err := th.TileService.Store.UpdateSubmissionState(context.Background(), db.UpdateSubmissionStateParams{
		ID:    int32(submissionId),
		State: state,
	})

	if err != nil {
		c.Set("ISERROR", true)
		setFlashmessages(c, "error", "Could not update submission state")
		return err
	}

	if !util.IsEmptyOrWhitespace(comment) { // FIXME: bleh
		uid, ok := c.Get(user_id_key).(int32)

		if ok {
			err = th.TileService.Store.CreateSubmissionComment(context.TODO(), db.CreateSubmissionCommentParams{
				SubmissionID: int32(submissionId),
				LoginID:      int32(uid),
				Comment:      comment,
			})
			if err != nil {
				return err
			}
		}
	}

	sc, err := th.TileService.Store.GetCommentsForSubmission(context.Background(), s.ID)
	if err != nil {
		c.Set("ISERROR", true)
		setFlashmessages(c, "error", "Could not update submission state")
		return err
	}

	ip, err := th.TileService.Store.GetImagesForSubmission(context.Background(), s.ID)
	if err != nil {
		c.Set("ISERROR", true)
		setFlashmessages(c, "error", "Could not update submission state")
		return err
	}

	c.Set("ISERROR", false)

	data := views.SubmissionData{
		Submission: s,
		Comments:   sc,
		Images:     ip,
	}

	return render(c, components.SubmissionHeader(isManagement, data))
}

func (th *TileHandler) handleTileSubmission(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}
	if !isAuthenticated {
		c.Redirect(http.StatusUnauthorized, "/login")
	}

	var tileId int
	echo.PathParamsBinder(c).Int("tileId", &tileId)

	tile, err := th.TileService.LoadTile(tileId)
	if err != nil {
		return err
	}
	fmt.Printf("Tile: %#v", tile)

	c.Set("ISERROR", false)

	if c.Request().Method == "POST" {
		sc, err := th.TileService.Store.GetSubmissionClosedStatusForBingo(context.TODO(), tile.BingoID)
		if err != nil {
			return echo.ErrInternalServerError
		}

		if sc {
			setFlashmessages(c, "error", "Submissions are closed")
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tileId))
		}

		form, err := c.MultipartForm()

		if err != nil {
			return err
		}
		files := form.File["files"]
		filePaths := []string{}
		for _, file := range files {
			err = util.ValidateImageFile(file)
			if err != nil {
				setFlashmessages(c, "error", "Can only upload images!")
				c.Set("ISERROR", true)
				return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tiles/%d", tile.ID))
			}
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
		err = th.TileService.CreateSubmission(tileId, u, filePaths)

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

func (bh *TileHandler) handleGetTileSubmissions(c echo.Context) error {
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

	tile, err := bh.TileService.LoadTile(int(tileId))
	if err != nil {
		return fmt.Errorf("Could not load tile: %w", err)
	}

	var submissionMap views.Submissions
	if isManagement {
		submissionMap, err = bh.TileService.LoadAllSubmissionsForTile(tileId)
	} else {
		submissionMap, err = bh.TileService.LoadUserSubmissions(tileId, uid)
	}

	if err != nil {
		return err
	}

	sc, err := bh.TileService.Store.GetSubmissionClosedStatusForBingo(context.TODO(), tile.BingoID)
	if err != nil {
		return err
	}

	return render(c, components.Submissions(isManagement, sc, submissionMap))
}

type imagePathSelectorFunc func(string) string

func imageUrlSelectorFunction(c echo.Context, fieldName string) imagePathSelectorFunc {
	return func(def string) string {
		var imageDefaultValue string
		err := echo.
			FormFieldBinder(c).
			String(fieldName, &imageDefaultValue).
			BindError()

		if err != nil || util.IsEmptyOrWhitespace(imageDefaultValue) {
			return def
		}
		return imageDefaultValue
	}
}

func (th *TileHandler) handleTile(c echo.Context) error {
	isAuthenticated, ok := c.Get("ISAUTHENTICATED").(bool)
	if !ok {
		return errors.New("invalid type for key 'ISAUTHENTICATED'")
	}

	isManagement, ok := c.Get(mgmnt_key).(bool)
	if !ok {
		isManagement = false
	}

	var tileId int
	err := echo.PathParamsBinder(c).Int("tileId", &tileId).BindError()
	if err != nil {
		return echo.ErrNotFound
	}

	tile, err := th.TileService.LoadTile(tileId)
	if err != nil {
		return echo.ErrNotFound
	}

	c.Set("ISERROR", false)

	if c.Request().Method == "PUT" {
		primaryImageSelector := imageUrlSelectorFunction(c, "primaryImageUrl")
		primaryImage, err := handleFileUpload(c, "file", "imagepath", primaryImageSelector(tile.Imagepath))
		if err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"Wrong image filetype: %s",
					err,
				))
		}
		secondaryImagePathSelector := imageUrlSelectorFunction(c, "secondaryImageUrl")
		secondaryImage, err := handleFileUpload(c, "secondaryfile", "secondaryimage", secondaryImagePathSelector(tile.SecondaryImagePath))

		if err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"Wrong image filetype: %s",
					err,
				))
		}

		t := db.UpdateTileParams{ // TODO: Read from config
			ID:                 int32(tileId),
			Imagepath:          primaryImage,
			SecondaryImagePath: secondaryImage,
		}

		err = echo.
			FormFieldBinder(c).
			String("title", &t.Title).
			String("description", &t.Description).
			Int32("weight", &t.Weight).
			BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error while binding form fields")
		}

		_, err = th.TileService.Store.UpdateTile(context.Background(), t)

		if err != nil {
			return echo.NewHTTPError(
				echo.ErrInternalServerError.Code,
				fmt.Sprintf(
					"something went wrong: %s",
					err,
				))
		}

		if c.FormValue("saveAsTemplate") == "on" {
			_, _ = th.TileService.Store.CreateTemplateTile(context.Background(), db.CreateTemplateTileParams{
				Title:              t.Title,
				Imagepath:          primaryImage,
				Description:        t.Description,
				Weight:             t.Weight,
				SecondaryImagePath: secondaryImage,
			})
		}

		setFlashmessages(c, "success", "Updated tile!")

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/bingos/%d", tile.BingoID))
	}

	uid := c.Get(user_id_key).(int32)
	var submissions views.Submissions
	var templates []db.TemplateTile

	if isManagement {
		submissions, err = th.TileService.LoadAllSubmissionsForTile(tile.ID)
		templates, err = th.TileService.Store.GetTemplateTiles(context.TODO())
	} else {
		submissions, err = th.TileService.LoadUserSubmissions(tile.ID, uid)
	}

	sc, err := th.TileService.Store.GetSubmissionClosedStatusForBingo(context.TODO(), tile.BingoID)

	tm := views.TileModel{
		Tile:             tile,
		Submissions:      submissions,
		Templates:        templates,
		SubmissionClosed: sc,
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
