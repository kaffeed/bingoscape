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

type BingoService struct {
	Db    *pgxpool.Pool
	Store *db.Queries
}

func NewBingoService(store *db.Queries, pool *pgxpool.Pool) *BingoService {
	return &BingoService{
		Store: store,
		Db:    pool,
	}
}

func (bs *BingoService) LoadUserSubmissions(tileId int32, loginId int32) (views.Submissions, error) {
	return bs.loadSubmissions(tileId, &loginId)
}

func (bs *BingoService) LoadAllSubmissionsForTile(tileId int32) (views.Submissions, error) {
	return bs.loadSubmissions(tileId, nil)
}

// TODO: change submission loading logic
func (bs *BingoService) loadSubmissions(tileId int32, loginId *int32) (views.Submissions, error) {

	var subs []db.GetSubmissionsForTileAndLoginRow
	var err error
	if loginId != nil {
		subs, err = bs.Store.GetSubmissionsForTileAndLogin(context.Background(), db.GetSubmissionsForTileAndLoginParams{
			TileID:  tileId,
			LoginID: *loginId,
		})
		if err != nil {
			return nil, err
		}
	} else {
		tmp, err := bs.Store.GetSubmissionsForTile(context.Background(), int32(tileId))
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
		c, _ := bs.Store.GetCommentsForSubmission(context.Background(), s.Submission.ID)
		i, _ := bs.Store.GetImagesForSubmission(context.Background(), s.Submission.ID)
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

func (bs *BingoService) CreateSubmission(tileId int, loginId int32, filePaths []string) error {
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

func (bs *BingoService) GetPossibleParticipants(bingoId int) (views.PossibleBingoParticipants, error) {
	return bs.Store.GetPossibleBingoParticipants(context.Background(), int32(bingoId))
}

func (bs *BingoService) GetParticipants(bingoId int) (views.BingoParticipants, error) {
	return bs.Store.GetBingoParticipants(context.Background(), int32(bingoId))
}

func (bs *BingoService) GetBingos(isManagement bool, userId int32) ([]db.Bingo, error) {
	if isManagement {
		return bs.Store.GetBingos(context.Background())
	}
	return bs.Store.GetBingosForLogin(context.Background(), userId)
}

func (bs *BingoService) GetBingo(bingoId int) (db.Bingo, error) {
	return bs.Store.GetBingoById(context.Background(), int32(bingoId))
}

func (bs *BingoService) loadSubmissionById(submissionId int) (db.Submission, error) {
	return bs.Store.GetSubmissionById(context.Background(), int32(submissionId))
}

func (bs *BingoService) UpdateSubmissionState(submissionId int, state db.Submissionstate) (db.Submission, error) {
	return bs.Store.UpdateSubmissionState(context.Background(), db.UpdateSubmissionStateParams{
		ID:    int32(submissionId),
		State: state,
	})
}

func (bs *BingoService) LoadTile(id int) (db.Tile, error) {
	return bs.Store.GetTileById(context.Background(), int32(id))
}

func (bs *BingoService) LoadTilesForBingo(bingoId int) ([]views.TileModel, error) {
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

func (bs *BingoService) RemoveParticipation(pId, bId int) error {
	return bs.Store.DeleteBingoParticipant(context.Background(), db.DeleteBingoParticipantParams{
		LoginID: int32(pId),
		BingoID: int32(bId),
	})
}
func (bs *BingoService) AddParticipantToBingo(pId, bId int) error {
	return bs.Store.CreateBingoParticipant(context.Background(), db.CreateBingoParticipantParams{
		LoginID: int32(pId),
		BingoID: int32(bId),
	})
}

func (bs *BingoService) CreateBingo(b db.CreateBingoParams) (views.BingoDetailModel, error) {
	tx, err := bs.Db.Begin(context.Background())
	if err != nil {
		return views.BingoDetailModel{}, err
	}
	defer tx.Rollback(context.Background())
	qtx := bs.Store.WithTx(tx)

	bingo, err := qtx.CreateBingo(context.Background(), b)

	tiles := make([]db.Tile, b.Rows*b.Cols)
	for i := 0; i < int(b.Rows*b.Cols); i++ {
		tiles[i], err = qtx.CreateTile(context.Background(), db.CreateTileParams{
			Title:       fmt.Sprintf("Tile %d", i+1),
			Imagepath:   "https://i.ibb.co/7N9Pjcs/image.png",
			Description: fmt.Sprintf("This is tile %d", i),
			BingoID:     bingo.ID,
		})
	}

	if err := tx.Commit(context.Background()); err != nil {
		return views.BingoDetailModel{}, err
	}
	tm := make([]views.TileModel, len(tiles))
	for _, t := range tiles {
		tm = append(tm, views.TileModel{
			Tile:        t,
			Submissions: map[string][]views.SubmissionData{},
		})
	}

	res := views.BingoDetailModel{
		Bingo: bingo,
		Tiles: tm,
	}
	return res, nil
}
