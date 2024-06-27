package services

import (
	"context"
	"fmt"

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

func (bs *BingoService) GetPossibleParticipants(bingoId int32) (views.PossibleBingoParticipants, error) {
	return bs.Store.GetPossibleBingoParticipants(context.Background(), bingoId)
}

func (bs *BingoService) GetParticipants(bingoId int32) (views.BingoParticipants, error) {
	bp, err := bs.Store.GetBingoParticipants(context.Background(), int32(bingoId))
	if err != nil {
		return nil, err
	}
	bps := views.BingoParticipants{}
	for _, p := range bp {
		stats, err := bs.Store.GetStatsByLoginAndBingo(context.TODO(), db.GetStatsByLoginAndBingoParams{
			LoginID: p.ID,
			BingoID: bingoId,
		})

		if err != nil {
			return nil, err
		}
		bps[p] = stats
	}

	return bps, nil
}

func (bs *BingoService) GetBingos(isManagement bool, userId int32) ([]db.Bingo, error) {
	if isManagement {
		return bs.Store.GetBingos(context.Background())
	}
	return bs.Store.GetBingosForLogin(context.Background(), userId)
}

func (bs *BingoService) GetBingo(bingoId int32) (db.Bingo, error) {
	return bs.Store.GetBingoById(context.Background(), bingoId)
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

func (bs *BingoService) RemoveParticipation(pId, bId int32) error {
	return bs.Store.DeleteBingoParticipant(context.Background(), db.DeleteBingoParticipantParams{
		LoginID: pId,
		BingoID: bId,
	})
}
func (bs *BingoService) AddParticipantToBingo(pId, bId int32) error {
	return bs.Store.CreateBingoParticipant(context.Background(), db.CreateBingoParticipantParams{
		LoginID: pId,
		BingoID: bId,
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
	if err != nil {
		return views.BingoDetailModel{}, err
	}
	tiles := make([]db.Tile, b.Rows*b.Cols)
	for i := 0; i < int(b.Rows*b.Cols); i++ {
		tiles[i], err = qtx.CreateTile(context.TODO(), db.CreateTileParams{
			Title:              fmt.Sprintf("Tile %d", i+1),
			Imagepath:          "https://i.ibb.co/7N9Pjcs/image.png",
			SecondaryImagePath: "https://i.ibb.co/7N9Pjcs/image.png",
			Description:        fmt.Sprintf("This is tile %d", i),
			BingoID:            bingo.ID,
			Weight:             int32(1),
		})
		if err != nil {
			return views.BingoDetailModel{}, err
		}
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
