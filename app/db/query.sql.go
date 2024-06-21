// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBingo = `-- name: CreateBingo :one
INSERT INTO bingos (title, validFrom, validTo, rows, cols, description, active, codephrase) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, title, validfrom, validto, rows, cols, description, codephrase, active
`

type CreateBingoParams struct {
	Title       string
	Validfrom   pgtype.Timestamp
	Validto     pgtype.Timestamp
	Rows        int32
	Cols        int32
	Description string
	Active      bool
	Codephrase  string
}

func (q *Queries) CreateBingo(ctx context.Context, arg CreateBingoParams) (Bingo, error) {
	row := q.db.QueryRow(ctx, createBingo,
		arg.Title,
		arg.Validfrom,
		arg.Validto,
		arg.Rows,
		arg.Cols,
		arg.Description,
		arg.Active,
		arg.Codephrase,
	)
	var i Bingo
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Validfrom,
		&i.Validto,
		&i.Rows,
		&i.Cols,
		&i.Description,
		&i.Codephrase,
		&i.Active,
	)
	return i, err
}

const createBingoParticipant = `-- name: CreateBingoParticipant :exec
INSERT INTO bingos_logins (login_id, bingo_id) VALUES ($1, $2)
`

type CreateBingoParticipantParams struct {
	LoginID int32
	BingoID int32
}

func (q *Queries) CreateBingoParticipant(ctx context.Context, arg CreateBingoParticipantParams) error {
	_, err := q.db.Exec(ctx, createBingoParticipant, arg.LoginID, arg.BingoID)
	return err
}

const createLogin = `-- name: CreateLogin :exec
INSERT INTO logins(password, name, is_management) VALUES($1, $2, $3) returning id, name, is_management, password
`

type CreateLoginParams struct {
	Password     string
	Name         string
	IsManagement bool
}

func (q *Queries) CreateLogin(ctx context.Context, arg CreateLoginParams) error {
	_, err := q.db.Exec(ctx, createLogin, arg.Password, arg.Name, arg.IsManagement)
	return err
}

const createSubmission = `-- name: CreateSubmission :one
INSERT INTO public.submissions (login_id, tile_id, state) values ($1, $2, $3) returning id, login_id, tile_id, date, state
`

type CreateSubmissionParams struct {
	LoginID int32
	TileID  int32
	State   Submissionstate
}

func (q *Queries) CreateSubmission(ctx context.Context, arg CreateSubmissionParams) (Submission, error) {
	row := q.db.QueryRow(ctx, createSubmission, arg.LoginID, arg.TileID, arg.State)
	var i Submission
	err := row.Scan(
		&i.ID,
		&i.LoginID,
		&i.TileID,
		&i.Date,
		&i.State,
	)
	return i, err
}

const createSubmissionComment = `-- name: CreateSubmissionComment :exec
INSERT INTO public.submission_comments (submission_id, login_id, comment) values ($1, $2, $3)
`

type CreateSubmissionCommentParams struct {
	SubmissionID int32
	LoginID      int32
	Comment      string
}

func (q *Queries) CreateSubmissionComment(ctx context.Context, arg CreateSubmissionCommentParams) error {
	_, err := q.db.Exec(ctx, createSubmissionComment, arg.SubmissionID, arg.LoginID, arg.Comment)
	return err
}

const createSubmissionImage = `-- name: CreateSubmissionImage :exec
INSERT INTO submission_images(path, submission_id) VALUES ($1, $2)
`

type CreateSubmissionImageParams struct {
	Path         string
	SubmissionID int32
}

func (q *Queries) CreateSubmissionImage(ctx context.Context, arg CreateSubmissionImageParams) error {
	_, err := q.db.Exec(ctx, createSubmissionImage, arg.Path, arg.SubmissionID)
	return err
}

const createTemplateTile = `-- name: CreateTemplateTile :one
INSERT INTO template_tiles(title, imagepath, description, weight, secondary_image_path) VALUES ($1, $2, $3, $4, $5) returning id, title, imagepath, description, weight, secondary_image_path
`

type CreateTemplateTileParams struct {
	Title              string
	Imagepath          string
	Description        string
	Weight             int32
	SecondaryImagePath string
}

func (q *Queries) CreateTemplateTile(ctx context.Context, arg CreateTemplateTileParams) (TemplateTile, error) {
	row := q.db.QueryRow(ctx, createTemplateTile,
		arg.Title,
		arg.Imagepath,
		arg.Description,
		arg.Weight,
		arg.SecondaryImagePath,
	)
	var i TemplateTile
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Imagepath,
		&i.Description,
		&i.Weight,
		&i.SecondaryImagePath,
	)
	return i, err
}

const createTile = `-- name: CreateTile :one
INSERT INTO tiles(title, imagepath, description, bingo_id, weight, secondary_image_path) VALUES ($1, $2, $3, $4, $5, $6) returning id, title, imagepath, description, bingo_id, weight, secondary_image_path
`

type CreateTileParams struct {
	Title              string
	Imagepath          string
	Description        string
	BingoID            int32
	Weight             int32
	SecondaryImagePath string
}

func (q *Queries) CreateTile(ctx context.Context, arg CreateTileParams) (Tile, error) {
	row := q.db.QueryRow(ctx, createTile,
		arg.Title,
		arg.Imagepath,
		arg.Description,
		arg.BingoID,
		arg.Weight,
		arg.SecondaryImagePath,
	)
	var i Tile
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Imagepath,
		&i.Description,
		&i.BingoID,
		&i.Weight,
		&i.SecondaryImagePath,
	)
	return i, err
}

const deleteBingoById = `-- name: DeleteBingoById :exec
DELETE FROM bingos WHERE id = $1
`

func (q *Queries) DeleteBingoById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteBingoById, id)
	return err
}

const deleteBingoParticipant = `-- name: DeleteBingoParticipant :exec
DELETE FROM bingos_logins WHERE login_id = $1 AND bingo_id = $2
`

type DeleteBingoParticipantParams struct {
	LoginID int32
	BingoID int32
}

func (q *Queries) DeleteBingoParticipant(ctx context.Context, arg DeleteBingoParticipantParams) error {
	_, err := q.db.Exec(ctx, deleteBingoParticipant, arg.LoginID, arg.BingoID)
	return err
}

const deleteLogin = `-- name: DeleteLogin :exec
DELETE FROM logins WHERE id = $1
`

func (q *Queries) DeleteLogin(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteLogin, id)
	return err
}

const deleteSubmissionById = `-- name: DeleteSubmissionById :exec
delete from submissions where id = $1
`

func (q *Queries) DeleteSubmissionById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteSubmissionById, id)
	return err
}

const deleteTemplateById = `-- name: DeleteTemplateById :exec
delete from template_tiles where id = $1
`

func (q *Queries) DeleteTemplateById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteTemplateById, id)
	return err
}

const getAllLogins = `-- name: GetAllLogins :many
select id, name, is_management, password from public.logins order by id asc
`

func (q *Queries) GetAllLogins(ctx context.Context) ([]Login, error) {
	rows, err := q.db.Query(ctx, getAllLogins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Login
	for rows.Next() {
		var i Login
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsManagement,
			&i.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBingoById = `-- name: GetBingoById :one
SELECT id, title, validfrom, validto, rows, cols, description, codephrase, active FROM bingos WHERE id = $1
`

func (q *Queries) GetBingoById(ctx context.Context, id int32) (Bingo, error) {
	row := q.db.QueryRow(ctx, getBingoById, id)
	var i Bingo
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Validfrom,
		&i.Validto,
		&i.Rows,
		&i.Cols,
		&i.Description,
		&i.Codephrase,
		&i.Active,
	)
	return i, err
}

const getBingoLeaderboard = `-- name: GetBingoLeaderboard :many
select l.name, sum(t.weight) as points from submissions s
JOIN logins as l on l.id = s.login_id
JOIN tiles as t ON s.tile_id = t.id
JOIN bingos_logins as bl on bl.login_id = l.id
WHERE bl.bingo_id = $1 and s.state = 'Accepted'::SUBMISSIONSTATE
GROUP BY l.name
ORDER BY points desc
`

type GetBingoLeaderboardRow struct {
	Name   string
	Points int64
}

func (q *Queries) GetBingoLeaderboard(ctx context.Context, bingoID int32) ([]GetBingoLeaderboardRow, error) {
	rows, err := q.db.Query(ctx, getBingoLeaderboard, bingoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBingoLeaderboardRow
	for rows.Next() {
		var i GetBingoLeaderboardRow
		if err := rows.Scan(&i.Name, &i.Points); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBingoParticipants = `-- name: GetBingoParticipants :many
SELECT l.Id, l.name FROM public.logins l
	JOIN bingos_logins bl ON l.id = bl.login_id
	WHERE bl.bingo_id = $1
`

type GetBingoParticipantsRow struct {
	ID   int32
	Name string
}

func (q *Queries) GetBingoParticipants(ctx context.Context, bingoID int32) ([]GetBingoParticipantsRow, error) {
	rows, err := q.db.Query(ctx, getBingoParticipants, bingoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBingoParticipantsRow
	for rows.Next() {
		var i GetBingoParticipantsRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBingos = `-- name: GetBingos :many
SELECT id, title, validfrom, validto, rows, cols, description, codephrase, active FROM bingos
`

func (q *Queries) GetBingos(ctx context.Context) ([]Bingo, error) {
	rows, err := q.db.Query(ctx, getBingos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bingo
	for rows.Next() {
		var i Bingo
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Validfrom,
			&i.Validto,
			&i.Rows,
			&i.Cols,
			&i.Description,
			&i.Codephrase,
			&i.Active,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBingosForLogin = `-- name: GetBingosForLogin :many
SELECT b.id, b.title, b.validfrom, b.validto, b.rows, b.cols, b.description, b.codephrase, b.active FROM bingos b
JOIN bingos_logins bl ON b.id = bl.bingo_id
JOIN logins l ON bl.login_id = l.id
WHERE l.id = $1 and b.active
`

func (q *Queries) GetBingosForLogin(ctx context.Context, id int32) ([]Bingo, error) {
	rows, err := q.db.Query(ctx, getBingosForLogin, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bingo
	for rows.Next() {
		var i Bingo
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Validfrom,
			&i.Validto,
			&i.Rows,
			&i.Cols,
			&i.Description,
			&i.Codephrase,
			&i.Active,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCommentsForSubmission = `-- name: GetCommentsForSubmission :many
SELECT c.submission_id, l.name as managementuser, c.comment, c.created_at FROM public.submission_comments c
JOIN public.logins l on c.login_id = l.id
WHERE c.submission_id = $1
`

type GetCommentsForSubmissionRow struct {
	SubmissionID   int32
	Managementuser string
	Comment        string
	CreatedAt      pgtype.Timestamptz
}

func (q *Queries) GetCommentsForSubmission(ctx context.Context, submissionID int32) ([]GetCommentsForSubmissionRow, error) {
	rows, err := q.db.Query(ctx, getCommentsForSubmission, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsForSubmissionRow
	for rows.Next() {
		var i GetCommentsForSubmissionRow
		if err := rows.Scan(
			&i.SubmissionID,
			&i.Managementuser,
			&i.Comment,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getImagesForSubmission = `-- name: GetImagesForSubmission :many
SELECT path FROM public.submission_images WHERE submission_id = $1
`

func (q *Queries) GetImagesForSubmission(ctx context.Context, submissionID int32) ([]string, error) {
	rows, err := q.db.Query(ctx, getImagesForSubmission, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		items = append(items, path)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLoginById = `-- name: GetLoginById :one
SELECT id, name, is_management, password FROM logins WHERE id = $1
`

func (q *Queries) GetLoginById(ctx context.Context, id int32) (Login, error) {
	row := q.db.QueryRow(ctx, getLoginById, id)
	var i Login
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsManagement,
		&i.Password,
	)
	return i, err
}

const getLoginByName = `-- name: GetLoginByName :one
SELECT id, name, is_management, password FROM logins WHERE name = $1
`

func (q *Queries) GetLoginByName(ctx context.Context, name string) (Login, error) {
	row := q.db.QueryRow(ctx, getLoginByName, name)
	var i Login
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsManagement,
		&i.Password,
	)
	return i, err
}

const getPossibleBingoParticipants = `-- name: GetPossibleBingoParticipants :many
SELECT l.id, l.name FROM public.logins l
	WHERE l.id NOT IN (SELECT login_id from public.bingos_logins WHERE bingo_id = $1)
	AND not l.is_management
`

type GetPossibleBingoParticipantsRow struct {
	ID   int32
	Name string
}

func (q *Queries) GetPossibleBingoParticipants(ctx context.Context, bingoID int32) ([]GetPossibleBingoParticipantsRow, error) {
	rows, err := q.db.Query(ctx, getPossibleBingoParticipants, bingoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPossibleBingoParticipantsRow
	for rows.Next() {
		var i GetPossibleBingoParticipantsRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubmissionById = `-- name: GetSubmissionById :one
SELECT id, login_id, tile_id, date, state FROM submissions WHERE id = $1
`

func (q *Queries) GetSubmissionById(ctx context.Context, id int32) (Submission, error) {
	row := q.db.QueryRow(ctx, getSubmissionById, id)
	var i Submission
	err := row.Scan(
		&i.ID,
		&i.LoginID,
		&i.TileID,
		&i.Date,
		&i.State,
	)
	return i, err
}

const getSubmissionIdForTileAndLogin = `-- name: GetSubmissionIdForTileAndLogin :one
SELECT id FROM public.submissions WHERE tile_id = $1 AND login_id = $2
`

type GetSubmissionIdForTileAndLoginParams struct {
	TileID  int32
	LoginID int32
}

func (q *Queries) GetSubmissionIdForTileAndLogin(ctx context.Context, arg GetSubmissionIdForTileAndLoginParams) (int32, error) {
	row := q.db.QueryRow(ctx, getSubmissionIdForTileAndLogin, arg.TileID, arg.LoginID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getSubmissionsByBingoAndLogin = `-- name: GetSubmissionsByBingoAndLogin :many
select bingos_logins.bingo_id, submissions.id, submissions.login_id, submissions.tile_id, submissions.date, submissions.state, tiles.id, tiles.title, tiles.imagepath, tiles.description, tiles.bingo_id, tiles.weight, tiles.secondary_image_path from submissions  
join tiles on submissions.tile_id = tiles.id
join bingos_logins on tiles.bingo_id = bingos_logins.bingo_id
where submissions.login_id = $1 and bingos_logins.bingo_id = $2
ORDER BY tiles.id asc
`

type GetSubmissionsByBingoAndLoginParams struct {
	LoginID int32
	BingoID int32
}

type GetSubmissionsByBingoAndLoginRow struct {
	BingoID    int32
	Submission Submission
	Tile       Tile
}

func (q *Queries) GetSubmissionsByBingoAndLogin(ctx context.Context, arg GetSubmissionsByBingoAndLoginParams) ([]GetSubmissionsByBingoAndLoginRow, error) {
	rows, err := q.db.Query(ctx, getSubmissionsByBingoAndLogin, arg.LoginID, arg.BingoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubmissionsByBingoAndLoginRow
	for rows.Next() {
		var i GetSubmissionsByBingoAndLoginRow
		if err := rows.Scan(
			&i.BingoID,
			&i.Submission.ID,
			&i.Submission.LoginID,
			&i.Submission.TileID,
			&i.Submission.Date,
			&i.Submission.State,
			&i.Tile.ID,
			&i.Tile.Title,
			&i.Tile.Imagepath,
			&i.Tile.Description,
			&i.Tile.BingoID,
			&i.Tile.Weight,
			&i.Tile.SecondaryImagePath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubmissionsForTile = `-- name: GetSubmissionsForTile :many
SELECT submissions.id, submissions.login_id, submissions.tile_id, submissions.date, submissions.state, logins.id, logins.name, logins.is_management, logins.password
FROM submissions
JOIN logins ON logins.id = submissions.login_id
WHERE submissions.tile_id = $1
`

type GetSubmissionsForTileRow struct {
	Submission Submission
	Login      Login
}

func (q *Queries) GetSubmissionsForTile(ctx context.Context, tileID int32) ([]GetSubmissionsForTileRow, error) {
	rows, err := q.db.Query(ctx, getSubmissionsForTile, tileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubmissionsForTileRow
	for rows.Next() {
		var i GetSubmissionsForTileRow
		if err := rows.Scan(
			&i.Submission.ID,
			&i.Submission.LoginID,
			&i.Submission.TileID,
			&i.Submission.Date,
			&i.Submission.State,
			&i.Login.ID,
			&i.Login.Name,
			&i.Login.IsManagement,
			&i.Login.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubmissionsForTileAndLogin = `-- name: GetSubmissionsForTileAndLogin :many
SELECT submissions.id, submissions.login_id, submissions.tile_id, submissions.date, submissions.state, logins.id, logins.name, logins.is_management, logins.password
FROM submissions
JOIN logins ON logins.id = submissions.login_id
WHERE submissions.tile_id = $1 AND submissions.login_id = $2
`

type GetSubmissionsForTileAndLoginParams struct {
	TileID  int32
	LoginID int32
}

type GetSubmissionsForTileAndLoginRow struct {
	Submission Submission
	Login      Login
}

func (q *Queries) GetSubmissionsForTileAndLogin(ctx context.Context, arg GetSubmissionsForTileAndLoginParams) ([]GetSubmissionsForTileAndLoginRow, error) {
	rows, err := q.db.Query(ctx, getSubmissionsForTileAndLogin, arg.TileID, arg.LoginID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubmissionsForTileAndLoginRow
	for rows.Next() {
		var i GetSubmissionsForTileAndLoginRow
		if err := rows.Scan(
			&i.Submission.ID,
			&i.Submission.LoginID,
			&i.Submission.TileID,
			&i.Submission.Date,
			&i.Submission.State,
			&i.Login.ID,
			&i.Login.Name,
			&i.Login.IsManagement,
			&i.Login.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTemplateImagePath = `-- name: GetTemplateImagePath :one
SELECT imagepath from template_tiles where id = $1
`

func (q *Queries) GetTemplateImagePath(ctx context.Context, id int32) (string, error) {
	row := q.db.QueryRow(ctx, getTemplateImagePath, id)
	var imagepath string
	err := row.Scan(&imagepath)
	return imagepath, err
}

const getTemplateTiles = `-- name: GetTemplateTiles :many
SELECT id, title, imagepath, description, weight, secondary_image_path FROM template_tiles
`

func (q *Queries) GetTemplateTiles(ctx context.Context) ([]TemplateTile, error) {
	rows, err := q.db.Query(ctx, getTemplateTiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TemplateTile
	for rows.Next() {
		var i TemplateTile
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Imagepath,
			&i.Description,
			&i.Weight,
			&i.SecondaryImagePath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTileById = `-- name: GetTileById :one
SELECT id, title, imagepath, description, bingo_id, weight, secondary_image_path FROM tiles WHERE id = $1
`

func (q *Queries) GetTileById(ctx context.Context, id int32) (Tile, error) {
	row := q.db.QueryRow(ctx, getTileById, id)
	var i Tile
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Imagepath,
		&i.Description,
		&i.BingoID,
		&i.Weight,
		&i.SecondaryImagePath,
	)
	return i, err
}

const getTilesForBingo = `-- name: GetTilesForBingo :many
SELECT id, title, imagepath, description, bingo_id, weight, secondary_image_path
FROM tiles 
WHERE bingo_id = $1 ORDER BY id ASC
`

func (q *Queries) GetTilesForBingo(ctx context.Context, bingoID int32) ([]Tile, error) {
	rows, err := q.db.Query(ctx, getTilesForBingo, bingoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tile
	for rows.Next() {
		var i Tile
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Imagepath,
			&i.Description,
			&i.BingoID,
			&i.Weight,
			&i.SecondaryImagePath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const makeUserManagement = `-- name: MakeUserManagement :exec
UPDATE logins SET is_management = true WHERE id = $1
`

func (q *Queries) MakeUserManagement(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, makeUserManagement, id)
	return err
}

const toggleBingoState = `-- name: ToggleBingoState :one
UPDATE bingos SET active = NOT active WHERE id = $1 returning active
`

func (q *Queries) ToggleBingoState(ctx context.Context, id int32) (bool, error) {
	row := q.db.QueryRow(ctx, toggleBingoState, id)
	var active bool
	err := row.Scan(&active)
	return active, err
}

const updateLoginPassword = `-- name: UpdateLoginPassword :one
UPDATE logins SET password = $2 WHERE id = $1 returning id, name, is_management, password
`

type UpdateLoginPasswordParams struct {
	ID       int32
	Password string
}

func (q *Queries) UpdateLoginPassword(ctx context.Context, arg UpdateLoginPasswordParams) (Login, error) {
	row := q.db.QueryRow(ctx, updateLoginPassword, arg.ID, arg.Password)
	var i Login
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsManagement,
		&i.Password,
	)
	return i, err
}

const updateSubmissionState = `-- name: UpdateSubmissionState :one
UPDATE public.submissions SET state = $2 WHERE id = $1 returning id, login_id, tile_id, date, state
`

type UpdateSubmissionStateParams struct {
	ID    int32
	State Submissionstate
}

func (q *Queries) UpdateSubmissionState(ctx context.Context, arg UpdateSubmissionStateParams) (Submission, error) {
	row := q.db.QueryRow(ctx, updateSubmissionState, arg.ID, arg.State)
	var i Submission
	err := row.Scan(
		&i.ID,
		&i.LoginID,
		&i.TileID,
		&i.Date,
		&i.State,
	)
	return i, err
}

const updateTile = `-- name: UpdateTile :one
UPDATE tiles SET title = $2, imagepath = $3, description = $4, weight = $5, secondary_image_path = $6 WHERE id = $1 returning id, title, imagepath, description, bingo_id, weight, secondary_image_path
`

type UpdateTileParams struct {
	ID                 int32
	Title              string
	Imagepath          string
	Description        string
	Weight             int32
	SecondaryImagePath string
}

func (q *Queries) UpdateTile(ctx context.Context, arg UpdateTileParams) (Tile, error) {
	row := q.db.QueryRow(ctx, updateTile,
		arg.ID,
		arg.Title,
		arg.Imagepath,
		arg.Description,
		arg.Weight,
		arg.SecondaryImagePath,
	)
	var i Tile
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Imagepath,
		&i.Description,
		&i.BingoID,
		&i.Weight,
		&i.SecondaryImagePath,
	)
	return i, err
}
