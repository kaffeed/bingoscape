package views

import "github.com/kaffeed/bingoscape/db"

type Submissions map[string][]SubmissionData
type PossibleBingoParticipants []db.GetPossibleBingoParticipantsRow
type BingoParticipants []db.GetBingoParticipantsRow

type SubmissionData struct {
	Submission db.Submission
	Comments   []db.GetCommentsForSubmissionRow
	Images     []string
}
type BingoDetailModel struct {
	db.Bingo
	Tiles                []TileModel
	PossibleParticipants PossibleBingoParticipants
	Participants         BingoParticipants
}

type TileModel struct {
	db.Tile
	Submissions Submissions
}

type TileStats struct {
	Submitted      int
	NeedReview     int
	Accepted       int
	State          db.Submissionstate
	HasSubmissions bool
}

func (t TileModel) Stats(loginId int32) TileStats {
	stat := TileStats{
		Submitted:  0,
		NeedReview: 0,
		Accepted:   0,
		State:      "",
	}

	for _, val := range t.Submissions {
		for _, s := range val {
			if s.Submission.LoginID == loginId {
				stat.State = s.Submission.State
			}
			switch s.Submission.State {
			case db.SubmissionstateSubmitted:
				stat.Submitted++
			case db.SubmissionstateActionRequired:
				stat.NeedReview++
			case db.SubmissionstateAccepted:
				stat.Accepted++
			}
		}
	}
	stat.HasSubmissions = stat.Submitted > 0 || stat.Accepted > 0 || stat.NeedReview > 0
	return stat
}
