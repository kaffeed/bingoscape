package services

import (
	"database/sql"
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

func (bs *BingoService) GetBingos(isManagement bool, userId int) ([]Bingo, error) {

	var query string
	if isManagement {
		query = `SELECT b.id, b.title, b.validFrom, b.validTo, b.size FROM bingos b`
	} else {
		query = `SELECT b.id, b.title, b.validFrom, b.validTo, b.size FROM bingos b
		JOIN bingos_logins bl ON b.id = bl.bingos_id
		JOIN logins l ON bl.logins_id = l.id
		WHERE b.id = $1`
	}

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return []Bingo{}, err
	}

	defer stmt.Close()

	bingos := make([]Bingo, 0)
	var rows *sql.Rows
	if isManagement {
		rows, err = stmt.Query()
	} else {
		rows, err = stmt.Query(userId)
	}

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
