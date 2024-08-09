package handlers

import (
	"net/http"

	"github.com/kaffeed/bingoscape/app/services"
	"github.com/labstack/echo/v4"
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

	return c.JSON(http.StatusOK, nil)
}
