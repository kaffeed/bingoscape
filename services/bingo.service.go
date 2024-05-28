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
	Id          int
	Title       string
	From        time.Time
	To          time.Time
	Rows        int
	Cols        int
	Description string
	Tiles       []Tile
	IsActive    bool
}

type Participant struct {
	Id   int
	Name string
}

type Participants []Participant

type Tile struct {
	Id          int
	Title       string
	ImagePath   string
	Description string
	BingoId     int
}

type TemplateTile Tile

type Tiles []Tile

func (tiles Tiles) BulkInsert(db *sql.DB) error {
	var (
		placeholders []string
		vals         []interface{}
	)

	fmt.Printf("VALUES: %#v", tiles)

	for index, tile := range tiles {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)",
			index*4+1,
			index*4+2,
			index*4+3,
			index*4+4,
		))

		vals = append(vals, tile.Title, tile.ImagePath, tile.Description, tile.BingoId)
	}

	fmt.Printf("VALUES: %#v", vals)
	insertStatement := fmt.Sprintf("INSERT INTO tiles(title, imagepath, description, bingo_id) VALUES %s", strings.Join(placeholders, ","))
	stmt, err := db.Prepare(insertStatement)
	if err != nil {
		return fmt.Errorf("Problem during statement preparation: %w", err)
	}
	fmt.Printf("db vals: %+v", vals)
	_, err = stmt.Exec(vals...)

	if err != nil {
		return fmt.Errorf("Error during batch insert: %w", err)
	}

	return nil
}

func NewBingoService(store db.Store) *BingoService {
	return &BingoService{
		store: store,
	}
}

func (bs *BingoService) CreateTemplateTile(t Tile) error {
	query := "INSERT INTO template_tiles(title, imagepath, description) VALUES ($1, $1, $1)"
	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Could not prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Title, t.ImagePath, t.Description)

	if err != nil {
		return fmt.Errorf("Error template tile insert: %w", err)
	}

	return nil
}

func (bs *BingoService) GetPossibleParticipants(bingoId int) (Participants, error) {
	query := `SELECT l.id, l.name FROM public.logins l
	WHERE l.id NOT IN (SELECT login_id from public.bingos_logins WHERE bingo_id = $1)
	AND not l.is_management`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("Error during query prep: %w", err)
	}
	defer stmt.Close()

	participants := Participants{}
	rows, err := stmt.Query(bingoId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := Participant{}
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func (bs *BingoService) GetParticipants(bingoId int) (Participants, error) {
	query := `SELECT l.Id, l.name FROM public.logins l
	JOIN bingos_logins bl ON l.id = bl.login_id
	WHERE bl.bingo_id = $1`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("Error during query prep: %w", err)
	}
	defer stmt.Close()

	participants := Participants{}
	rows, err := stmt.Query(bingoId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := Participant{}
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func (bs *BingoService) GetBingos(isManagement bool, userId int) ([]Bingo, error) {
	var query string
	if isManagement {
		query = `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols", b.description, b.isactive FROM bingos b`
	} else {
		query = `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols", b.description, b.isactive FROM bingos b
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
		if err := rows.Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Rows, &bingo.Cols, &bingo.Description, &bingo.IsActive); err != nil {
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
	query := `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols", b.description, b.isactive FROM bingos b WHERE b.id = $1`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return Bingo{}, err
	}
	defer stmt.Close()

	var bingo Bingo
	err = stmt.QueryRow(bingoId).Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Rows, &bingo.Cols, &bingo.Description, &bingo.IsActive)
	if err != nil {
		return Bingo{}, err
	}

	err = bingo.loadTiles(bs.store.Db)
	if err != nil {
		return Bingo{}, fmt.Errorf("Error during loading tiles: %w", err)
	}

	return bingo, nil
}

func (bs *BingoService) UpdateTile(t Tile) error {
	query := `UPDATE tiles SET title = $1, imagepath = $2, description = $3 WHERE id = $4`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Error during statement preparation: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Title, t.ImagePath, t.Description, t.Id)
	if err != nil {
		return fmt.Errorf("Could not update tile: %w", err)
	}
	return nil
}

func (bs *BingoService) LoadTile(id int) (Tile, error) {
	query := `SELECT id, title, imagepath, description, bingo_id FROM tiles WHERE id = $1`
	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return Tile{}, fmt.Errorf("Error during statement preparation: %w", err)
	}
	defer stmt.Close()

	var t Tile
	err = stmt.QueryRow(id).Scan(&t.Id, &t.Title, &t.ImagePath, &t.Description, &t.BingoId)
	if err != nil {
		return Tile{}, fmt.Errorf("Error during query row: %w", err)
	}
	return t, nil
}

func (b *Bingo) loadTiles(db *sql.DB) error {
	query := `SELECT id, imagepath, description, bingo_id FROM tiles WHERE bingo_id = $1 ORDER BY id ASC`
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Error during statement preparation: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(b.Id)

	var tiles Tiles
	for rows.Next() {
		var t Tile
		rows.Scan(&t.Id, &t.ImagePath, &t.Description, &t.BingoId)
		tiles = append(tiles, t)
	}
	b.Tiles = tiles
	return nil
}

func (bs *BingoService) RemoveParticipation(pId, bId int) error {
	query := `DELETE FROM bingos_logins WHERE login_id = $1 AND bingo_id = $2`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(pId, bId); err != nil {
		log.Printf("Error during execution of bingos participation")
		return err
	}
	log.Printf("Added team to bingo")
	return nil

}
func (bs *BingoService) AddParticipantToBingo(pId, bId int) error {
	query := `INSERT INTO bingos_logins (login_id, bingo_id) VALUES ($1, $2)`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(pId, bId); err != nil {
		log.Printf("Error during execution of bingos participation")
		return err
	}
	log.Printf("Added team to bingo")
	return nil

}

func (bs *BingoService) CreateBingo(b Bingo) (Bingo, error) {
	query := `INSERT INTO bingos (title, validFrom, validTo, rows, cols, description, isactive) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	stmt, err := bs.store.Db.Prepare(query)
	if err != nil {
		return b, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(
		b.Title,
		b.From,
		b.To,
		b.Rows,
		b.Cols,
		b.Description,
		b.IsActive,
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
