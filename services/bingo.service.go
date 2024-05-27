package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	Rows  int
	Cols  int
	Tiles []Tile
}

type Tile struct {
	Id          int
	ImagePath   string
	Description string
	BingoId     int
}

type Tiles []Tile

func (tiles Tiles) BulkInsert(db *sql.DB) error {
	log.Printf("Tiles: %+v", tiles)
	var (
		placeholders []string
		vals         []interface{}
	)

	fmt.Printf("VALUES: %#v", tiles)

	for index, tile := range tiles {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d)",
			index*3+1,
			index*3+2,
			index*3+3,
		))

		vals = append(vals, tile.ImagePath, tile.Description, tile.BingoId)
	}

	fmt.Printf("VALUES: %#v", vals)
	insertStatement := fmt.Sprintf("INSERT INTO tiles(image, description, bingo_id) VALUES %s", strings.Join(placeholders, ","))
	stmt, err := db.Prepare(insertStatement)
	if err != nil {
		return fmt.Errorf("Problem during statement preparation: %w", err)
	}
	fmt.Printf("db vals: %+v", vals)
	_, err = stmt.Exec(vals...)

	if err != nil {
		return fmt.Errorf("Error during batch insert: %w", err)
	}

	log.Printf("Should've successfully batch-inserted tiles")

	return nil
}

func NewBingoService(store db.Store) *BingoService {
	return &BingoService{
		store: store,
	}
}

func (bs *BingoService) GetBingos(isManagement bool, userId int) ([]Bingo, error) {
	var query string
	if isManagement {
		query = `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols" FROM bingos b`
	} else {
		query = `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols" FROM bingos b
		JOIN bingos_logins bl ON b.id = bl.bingo_id
		JOIN logins l ON bl.login_id = l.id
		WHERE l.id = $1`
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
		if err := rows.Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Rows, &bingo.Cols); err != nil {
			return nil, err
		}
		bingos = append(bingos, bingo)
	}

	if err != nil {
		return []Bingo{}, err
	}

	return bingos, nil
}

func (bs *BingoService) GetBingo(bingoId int) (Bingo, error) {
	query := `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols" FROM bingos b WHERE b.id = $1`

	stmt, err := bs.store.Db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return Bingo{}, err
	}

	var bingo Bingo
	err = stmt.QueryRow(bingoId).Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Rows, &bingo.Cols)
	if err != nil {
		return Bingo{}, err
	}

	err = bingo.loadTiles(bs.store.Db)
	if err != nil {
		return Bingo{}, fmt.Errorf("Error during loading tiles: %w", err)
	}

	return bingo, nil
}

func (b *Bingo) loadTiles(db *sql.DB) error {
	query := `SELECT id, image, description, bingo_id FROM tiles WHERE bingo_id = $1`
	stmt, err := db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("Error during statement preparation: %w", err)
	}

	rows, err := stmt.Query(b.Id)

	var tiles Tiles
	for rows.Next() {
		var t Tile
		rows.Scan(&t.Id, &t.ImagePath, &t.Description, &t.BingoId)
		tiles = append(tiles, t)
	}
	log.Printf("tiles: %+v", tiles)
	b.Tiles = tiles
	return nil
}

func (bs *BingoService) CreateBingo(b Bingo) (Bingo, error) {
	query := `INSERT INTO bingos (title, validFrom, validTo, rows, cols) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	stmt, err := bs.store.Db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return b, err
	}

	if err := stmt.QueryRow(
		b.Title,
		b.From,
		b.To,
		b.Rows,
		b.Cols,
	).Scan(&b.Id); err != nil {
		return Bingo{}, err
	}

	tiles := make(Tiles, b.Rows*b.Cols)
	for i := 0; i < b.Rows*b.Cols; i++ {
		tiles[i] = Tile{
			ImagePath:   "https://i.ibb.co/7N9Pjcs/image.png",
			Description: fmt.Sprintf("This is sample tile %d", i),
			BingoId:     b.Id,
		}
	}

	err = tiles.BulkInsert(bs.store.Db)
	if err != nil {
		return Bingo{}, fmt.Errorf("Error during bulk inserting tiles: %w", err)
	}

	b.Tiles = tiles
	return b, nil
}
