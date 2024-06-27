package services

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaffeed/bingoscape/app/db"
	"github.com/kaffeed/bingoscape/app/views"
)

type TileService struct {
	Db    *pgxpool.Pool
	Store *db.Queries
}

func NewTileService(store *db.Queries, pool *pgxpool.Pool) *TileService {
	return &TileService{
		Db:    pool,
		Store: store,
	}
}

func (s *TileService) LoadUserSubmissions(tileId int32, loginId int32) (views.Submissions, error) {
	return s.loadSubmissions(tileId, &loginId)
}

func (s *TileService) LoadAllSubmissionsForTile(tileId int32) (views.Submissions, error) {
	return s.loadSubmissions(tileId, nil)
}

func (ts *TileService) loadSubmissions(tileId int32, loginId *int32) (views.Submissions, error) {

	var subs []db.GetSubmissionsForTileAndLoginRow
	var err error
	if loginId != nil {
		subs, err = ts.Store.GetSubmissionsForTileAndLogin(context.Background(), db.GetSubmissionsForTileAndLoginParams{
			TileID:  tileId,
			LoginID: *loginId,
		})
		if err != nil {
			return nil, err
		}
	} else {
		tmp, err := ts.Store.GetSubmissionsForTile(context.Background(), int32(tileId))
		if err != nil {
			return nil, err
		}
		for _, s := range tmp {
			subs = append(subs, db.GetSubmissionsForTileAndLoginRow(s))
		}
	}

	submissions := make(map[string][]views.SubmissionData)
	for _, s := range subs {
		submissionData := views.SubmissionData{
			Submission: s.Submission,
			Comments:   []db.GetCommentsForSubmissionRow{},
			Images:     []string{},
		}
		c, _ := ts.Store.GetCommentsForSubmission(context.Background(), s.Submission.ID)
		i, _ := ts.Store.GetImagesForSubmission(context.Background(), s.Submission.ID)
		submissionData.Comments = c
		submissionData.Images = i

		subs, ok := submissions[s.Login.Name]

		if !ok {
			submissions[s.Login.Name] = []views.SubmissionData{submissionData}
		} else {
			submissions[s.Login.Name] = append(subs, submissionData)
		}
	}

	return submissions, nil
}

// TODO: change submission loading logic
func (bs *TileService) CreateSubmission(tileId int, loginId int32, filePaths []string) error {
	fail := func(err error) error {
		return fmt.Errorf("CreateSubmission: %w", err)
	}
	tx, err := bs.Db.Begin(context.Background())
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback(context.Background())
	qtx := bs.Store.WithTx(tx)

	submissionId, err := qtx.GetSubmissionIdForTileAndLogin(context.Background(), db.GetSubmissionIdForTileAndLoginParams{
		TileID:  int32(tileId),
		LoginID: int32(loginId),
	})
	if err != nil {
		s, err := qtx.CreateSubmission(context.Background(), db.CreateSubmissionParams{
			LoginID: int32(loginId),
			TileID:  int32(tileId),
			State:   db.SubmissionstateSubmitted,
		})
		if err != nil {
			return err
		}
		submissionId = s.ID
	}

	for _, path := range filePaths {
		_ = qtx.CreateSubmissionImage(context.Background(), db.CreateSubmissionImageParams{
			Path:         path,
			SubmissionID: submissionId,
		})
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fail(err)
	}
	return err
}

func (bs *TileService) LoadTile(id int) (db.Tile, error) {
	return bs.Store.GetTileById(context.Background(), int32(id))
}

func (bs *TileService) LoadTilesForBingo(bingoId int32) ([]views.TileModel, error) {
	tiles, err := bs.Store.GetTilesForBingo(context.Background(), int32(bingoId))
	if err != nil {
		return nil, err
	}

	tChan := make(chan views.TileModel)
	go func() {
		wg := sync.WaitGroup{}

		for _, t := range tiles {
			wg.Add(1)
			tmp := t
			go func() {
				tm := views.TileModel{}
				defer wg.Done()
				s, _ := bs.LoadAllSubmissionsForTile(tmp.ID)
				tm.Tile = tmp
				tm.Submissions = s
				tChan <- tm
			}()
		}
		wg.Wait()
		close(tChan)
	}()

	insertAt := func(data []views.TileModel, i int, v views.TileModel) []views.TileModel {
		if i == len(data) {
			return append(data, v)
		}
		data = append(data[:i+1], data[i:]...)

		data[i] = v

		return data
	}

	insertSorted := func(data []views.TileModel, v views.TileModel) []views.TileModel {
		i := sort.Search(len(data), func(i int) bool { return data[i].Tile.ID >= v.Tile.ID })
		return insertAt(data, i, v)
	}

	var res []views.TileModel
	for t := range tChan {
		res = insertSorted(res, t)
	}

	return res, nil
}
