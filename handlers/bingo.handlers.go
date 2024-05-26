package handlers

import (
	"errors"

	"github.com/kaffeed/bingoscape/services"
	components "github.com/kaffeed/bingoscape/views/components"
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
