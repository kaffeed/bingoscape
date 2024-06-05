package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kaffeed/bingoscape/db"
)

type State int

func (e State) String() string {
	switch e {
	case SUBMITTED:
		return "Submitted"
	case NEEDS_REVIEW:
		return "Needs review"
	case ACCEPTED:
		return "Accepted"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}
func (u *State) Scan(value interface{}) error { *u = State(value.(int64)); return nil }
func (u State) Value() (driver.Value, error)  { return int64(u), nil }

const (
	SUBMITTED State = iota
	NEEDS_REVIEW
	ACCEPTED
)

type BingoService struct {
	Store db.Store
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
	Ready       bool
	CodePhrase  string
}

type Comments []Comment
type Comment struct {
	Id        int
	Comment   string
	CreatedAt time.Time
	By        string
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
	Submissions Submissions
}

type TileStats struct {
	Submitted      int
	NeedReview     int
	Accepted       int
	State          State
	HasSubmissions bool
}

func (t Tile) Stats(loginId int) TileStats {
	stat := TileStats{
		Submitted:  0,
		NeedReview: 0,
		Accepted:   0,
		State:      -1,
	}

	for _, val := range t.Submissions {
		for _, s := range val {
			if s.LoginId == loginId {
				stat.State = s.State
			}
			switch s.State {
			case SUBMITTED:
				stat.Submitted++
			case NEEDS_REVIEW:
				stat.NeedReview++
			case ACCEPTED:
				stat.Accepted++
			}
		}
	}
	stat.HasSubmissions = stat.Submitted > 0 || stat.Accepted > 0 || stat.NeedReview > 0
	return stat
}

type Tiles []Tile

type Submission struct {
	Id          int
	SubmittedAt time.Time
	TileId      int
	LoginId     int
	State       State
	ImagePaths  []string
	Comments    Comments
}

func (s *Submission) LoadComments(db *sql.DB) error {
	query := `SELECT c.submission_id, l.name, c.comment, c.created_at FROM public.submission_comments c
	JOIN public.logins l on c.login_id = l.id
	WHERE c.submission_id = $1`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(s.Id)
	if err != nil {
		return err
	}
	comments := Comments{}
	for rows.Next() {
		var c Comment
		rows.Scan(&c.Id, &c.By, &c.Comment, &c.CreatedAt)
		comments = append(comments, c)
	}

	s.Comments = comments

	return nil
}

func (bs *BingoService) CreateSubmissionComment(submissionId, uid int, comment string) error {
	query := `INSERT INTO public.submission_comments (submission_id, login_id, comment, created_at) values ($1, $2, $3, $4)`
	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(submissionId, uid, comment, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *Submission) LoadImages(db *sql.DB) error {
	query := `SELECT path FROM public.submission_images WHERE submission_id = $1`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	rows, err := stmt.Query(s.Id)
	s.ImagePaths = []string{}
	for rows.Next() {
		var path string
		err = rows.Scan(&path)
		if err != nil {
			log.Printf("Problem during path scan %v", err)
			continue
		}
		s.ImagePaths = append(s.ImagePaths, path)
	}
	return nil
}

type Submissions map[string][]Submission

func (tiles Tiles) BulkInsert(db *sql.DB) error {
	var (
		placeholders []string
		vals         []interface{}
	)

	for index, tile := range tiles {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)",
			index*4+1,
			index*4+2,
			index*4+3,
			index*4+4,
		))

		vals = append(vals, tile.Title, tile.ImagePath, tile.Description, tile.BingoId)
	}

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
		Store: store,
	}
}

func (bs *BingoService) DeleteBingo(id int) error {
	query := `DELETE FROM bingos WHERE id = $1`
	stmt, err := bs.Store.Db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(id); err != nil {
		return err
	}
	return nil
}

func (bs *BingoService) LoadUserSubmissions(tileId int, loginId int) (Submissions, error) {
	return bs.loadSubmissions(tileId, &loginId)
}

func (bs *BingoService) LoadAllSubmissionsForTile(tileId int) (Submissions, error) {
	return bs.loadSubmissions(tileId, nil)
}

func (bs *BingoService) loadSubmissions(tileId int, loginId *int) (Submissions, error) {
	fail := func(err error) error {
		return fmt.Errorf("loadSubmissions: %w", err)
	}

	query := `SELECT s.id, s.login_id, l.name, s.tile_id, s.date, s.state
	FROM public.Submissions s 
	JOIN public.logins l ON l.id = s.login_id
	WHERE tile_id = $1`
	if loginId != nil {
		query = query + " AND login_id = $2"
	}

	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return nil, fail(err)
	}
	defer stmt.Close()

	var rows *sql.Rows

	if loginId != nil {
		rows, err = stmt.Query(tileId, loginId)
		if err != nil {
			return nil, fail(err)
		}
	} else {
		rows, err = stmt.Query(tileId)
	}

	submissions := make(map[string][]Submission)
	for rows.Next() {
		var s Submission
		var team string
		err := rows.Scan(&s.Id, &s.LoginId, &team, &s.TileId, &s.SubmittedAt, &s.State)
		if err != nil {
			log.Printf("Something happened,")
		}
		subs, ok := submissions[team]

		s.LoadImages(bs.Store.Db)
		s.LoadComments(bs.Store.Db)

		if !ok {
			submissions[team] = []Submission{s}
		} else {
			submissions[team] = append(subs, s)
		}
	}

	return submissions, nil
}

func (bs *BingoService) CreateSubmission(tileId int, loginId int, filePaths []string) error {
	fail := func(err error) error {
		return fmt.Errorf("CreateSubmission: %w", err)
	}
	tx, err := bs.Store.Db.BeginTx(context.Background(), nil)
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback()

	var submissionId int
	query := `SELECT id FROM public.submissions WHERE tile_id = $1 AND login_id = $2`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fail(err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(tileId, loginId).Scan(&submissionId); err != nil {
		query = `INSERT INTO public.submissions (login_id, tile_id, date) values ($1,$2,$3) returning id`
		stmt, err = tx.Prepare(query)
		if err != nil {
			return fail(err)
		}
		defer stmt.Close()

		if err := stmt.QueryRow(
			loginId, tileId, time.Now().UTC(),
		).Scan(&submissionId); err != nil {
			return fail(err)
		}
	}

	var (
		placeholders []string
		vals         []interface{}
	)

	for index, path := range filePaths {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d)",
			index*2+1,
			index*2+2,
		))

		vals = append(vals, path, submissionId)
	}

	insertStatement := fmt.Sprintf("INSERT INTO submission_images(path, submission_id) VALUES %s", strings.Join(placeholders, ","))
	stmt, err = tx.Prepare(insertStatement)
	if err != nil {
		return fail(err)
	}
	_, err = stmt.Exec(vals...)
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	_, err = bs.UpdateSubmissionState(submissionId, SUBMITTED)
	return err
}

func (bs *BingoService) CreateTemplateTile(t Tile) error {
	query := "INSERT INTO template_tiles(title, imagepath, description) VALUES ($1, $1, $1)"
	stmt, err := bs.Store.Db.Prepare(query)
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

	stmt, err := bs.Store.Db.Prepare(query)
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

	stmt, err := bs.Store.Db.Prepare(query)
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
	query := `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols", b.description, b.ready, b.codephrase FROM bingos b `
	if !isManagement {
		query = query + `
		JOIN bingos_logins bl ON b.id = bl.bingo_id
		JOIN logins l ON bl.login_id = l.id
		WHERE l.id = $1`
	}

	stmt, err := bs.Store.Db.Prepare(query)
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
		if err := rows.Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Rows, &bingo.Cols, &bingo.Description, &bingo.Ready, &bingo.CodePhrase); err != nil {
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
	query := `SELECT b.id, b.title, b.validFrom, b.validTo, b."rows", b."cols", b.description, b.ready, b.codephrase FROM bingos b WHERE b.id = $1`

	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return Bingo{}, err
	}
	defer stmt.Close()

	var bingo Bingo
	err = stmt.QueryRow(bingoId).Scan(&bingo.Id, &bingo.Title, &bingo.From, &bingo.To, &bingo.Rows, &bingo.Cols, &bingo.Description, &bingo.Ready, &bingo.CodePhrase)
	if err != nil {
		return Bingo{}, err
	}

	err = bingo.loadTiles(bs)
	if err != nil {
		return Bingo{}, fmt.Errorf("Error during loading tiles: %w", err)
	}

	return bingo, nil
}

func (bs *BingoService) loadSubmissionById(submissionId int) (Submission, error) {
	query := `SELECT s.id, s.login_id, s.tile_id, s.date, s.state
	FROM public.submissions s 
	WHERE s.id = $1`

	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return Submission{}, err
	}
	defer stmt.Close()

	var s Submission
	if err := stmt.QueryRow(submissionId).Scan(&s.Id, &s.LoginId, &s.TileId, &s.SubmittedAt, &s.State); err != nil {
		return Submission{}, err
	}

	return s, nil
}

func (bs *BingoService) UpdateSubmissionState(submissionId int, state State) (Submission, error) {
	query := `UPDATE public.submissions SET state = $1 WHERE id = $2`

	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return Submission{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(state, submissionId)
	if err != nil {
		return Submission{}, err
	}

	s, err := bs.loadSubmissionById(submissionId)
	if err != nil {
		return Submission{}, err
	}

	return s, nil
}

func (bs *BingoService) UpdateTile(t Tile) error {
	query := `UPDATE tiles SET title = $1, imagepath = $2, description = $3 WHERE id = $4`

	stmt, err := bs.Store.Db.Prepare(query)
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
	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return Tile{}, fmt.Errorf("Error during statement preparation: %w", err)
	}
	defer stmt.Close()

	var t Tile
	err = stmt.QueryRow(id).Scan(&t.Id, &t.Title, &t.ImagePath, &t.Description, &t.BingoId)
	if err != nil {
		return Tile{}, fmt.Errorf("Error during query row: %w", err)
	}

	s, _ := bs.LoadAllSubmissionsForTile(t.Id)
	t.Submissions = s

	return t, nil
}

func (b *Bingo) loadTiles(bs *BingoService) error {
	query := `SELECT t.id, t.imagepath, t.description, t.bingo_id
	FROM tiles t 
	WHERE bingo_id = $1 ORDER BY id ASC`
	stmt, err := bs.Store.Db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Error during statement preparation: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(b.Id)

	var tiles Tiles

	tChan := make(chan Tile)
	go func() {
		wg := sync.WaitGroup{}
		for rows.Next() {
			var t Tile
			rows.Scan(&t.Id, &t.ImagePath, &t.Description, &t.BingoId)
			wg.Add(1)
			go func() {
				defer wg.Done()
				s, _ := bs.LoadAllSubmissionsForTile(t.Id)
				t.Submissions = s
				tChan <- t
			}()
		}
		wg.Wait()
		close(tChan)
	}()

	insertAt := func(data Tiles, i int, v Tile) Tiles {
		if i == len(data) {
			// Insert at end is the easy case.
			return append(data, v)
		}

		// make space for the inserted element by shifting
		// values at the insertion index up one index. the call
		// to append does not allocate memory when cap(data) is
		// greater ​than len(data).
		data = append(data[:i+1], data[i:]...)

		// Insert the new element.
		data[i] = v

		// Return the updated slice.
		return data
	}

	insertSorted := func(data Tiles, v Tile) Tiles {
		i := sort.Search(len(data), func(i int) bool { return data[i].Id >= v.Id })
		return insertAt(data, i, v)
	}
	for t := range tChan {
		tiles = insertSorted(tiles, t)
	}
	b.Tiles = tiles
	log.Printf("###########################################################")
	log.Printf("No error")
	log.Printf("Bingo... %#v", b)
	log.Printf("###########################################################")
	return nil
}

func (bs *BingoService) RemoveParticipation(pId, bId int) error {
	query := `DELETE FROM bingos_logins WHERE login_id = $1 AND bingo_id = $2`

	stmt, err := bs.Store.Db.Prepare(query)
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

	stmt, err := bs.Store.Db.Prepare(query)
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
	query := `INSERT INTO bingos (title, validFrom, validTo, rows, cols, description, ready, codephrase) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	stmt, err := bs.Store.Db.Prepare(query)
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
		b.Ready,
		b.CodePhrase,
	).Scan(&b.Id); err != nil {
		return Bingo{}, err
	}

	tiles := make(Tiles, b.Rows*b.Cols)
	for i := 0; i < b.Rows*b.Cols; i++ {
		tiles[i] = Tile{
			Title:       fmt.Sprintf("Tile %d", i),
			ImagePath:   "https://i.ibb.co/7N9Pjcs/image.png",
			Description: fmt.Sprintf("This is tile %d", i),
			BingoId:     b.Id,
		}
	}

	err = tiles.BulkInsert(bs.Store.Db)
	if err != nil {
		return Bingo{}, fmt.Errorf("Error during bulk inserting tiles: %w", err)
	}

	b.Tiles = tiles
	return b, nil
}
