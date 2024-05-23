package services

import "github.com/kaffeed/bingoscape/db"

type BingoService struct {
	store db.Store
}

type Bingo struct {
}

func NewBingoService(store db.Store) *BingoService {
	return &BingoService{
		store: store,
	}
}
