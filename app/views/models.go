package views

import "github.com/kaffeed/bingoscape/app/db"

type Submissions map[string][]SubmissionData
type PossibleBingoParticipants []db.GetPossibleBingoParticipantsRow
type BingoParticipants map[db.GetBingoParticipantsRow]db.GetStatsByLoginAndBingoRow

type TeamSubmissionModel struct {
	Submissions       []SubmissionData
	BingoID           int32
	Name              string
	SubmissionsClosed bool
}

type SubmissionData struct {
	db.Submission
	Comments []db.GetCommentsForSubmissionRow
	Images   []string
	Tile     db.Tile
}

type BingoDetailModel struct {
	db.Bingo
	Tiles                []TileModel
	PossibleParticipants PossibleBingoParticipants
	Participants         BingoParticipants
	Leaderboard          []db.GetBingoLeaderboardRow
}

type TileModel struct {
	db.Tile
	Submissions      Submissions
	Templates        []db.TemplateTile
	SubmissionClosed bool
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
