package handlers

import (
	"github.com/kaffeed/bingoscape/services"
)

type BingoHandler struct {
	BingoService *services.BingoService
}

func NewBingoHandler(bingoservice *services.BingoService) *BingoHandler {
	return &BingoHandler{
		BingoService: bingoservice,
	}
}
