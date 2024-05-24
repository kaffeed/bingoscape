package services

import (
	"time"

	"github.com/kaffeed/bingoscape/db"
)

type BingoService struct {
	store db.Store
}

type Bingo struct {
	Id    int
	Title string
	From  time.Time
	To    time.Time
	Size  int
}

func NewBingoService(store db.Store) *BingoService {
	return &BingoService{
		store: store,
	}
}

func (bs *BingoService) GetBingos() ([]Bingo, error) {
	query := `SELECT id, title, "from", "to", size FROM bingos`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return []Bingo{}, err
	}

	defer stmt.Close()

	bingos := make([]Bingo, 0)
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		bingo := Bingo{}
		if err := rows.Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Size); err != nil {
			return nil, err
		}
		bingos = append(bingos, bingo)
	}

	if err != nil {
		return []Bingo{}, err
	}

	return bingos, nil // TODO:
}
