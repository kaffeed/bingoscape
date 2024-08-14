package handlers

import (
	"net/http"

	"github.com/kaffeed/bingoscape/app/services"
	"github.com/kaffeed/bingoscape/app/views"
	"github.com/labstack/echo/v4"
)

const (
	user_key string = "user_key"
)

type ApiHandler struct {
	BingoService *services.BingoService
	TileService  *services.TileService
	UserService  *services.UserService
}

func NewApiHandler(ts *services.TileService, us *services.UserService, bs *services.BingoService) *ApiHandler {
	return &ApiHandler{
		TileService:  ts,
		UserService:  us,
		BingoService: bs,
	}
}

func (ah *ApiHandler) handleGetBingo(c echo.Context) error {
	var bingoId int32
	err := echo.PathParamsBinder(c).Int32("bingoId", &bingoId).BindError()
	if err != nil {
		return err
	}
	b, err := ah.BingoService.GetBingo(bingoId)
	if err != nil {
		return err
	}
	tiles, err := ah.TileService.LoadTilesForBingo(b.ID)
	if err != nil {
		return err
	}

	bv := views.BingoDetailModel{
		Bingo: b,
		Tiles: tiles,
	}

	return c.JSON(http.StatusOK, bv)
}
