package handlers

import (
	"errors"

	"github.com/kaffeed/bingoscape/services"
	"github.com/labstack/echo/v4"
)

type BingoHandler struct {
	BingoService *services.BingoService
}

func NewBingoHandler(bingoservice *services.BingoService) *BingoHandler {
	return &BingoHandler{
		BingoService: bingoservice,
	}
}

func (bh *BingoHandler) handleGetActiveBingo(c echo.Context) error {
	if c.Request().Method != "GET" {
		return errors.New("Invalid http method!")
	}

	return nil
}
