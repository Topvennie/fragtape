package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type Result string

const (
	ResultWin     Result = "win"
	ResultLoss    Result = "loss"
	ResultTie     Result = "tie"
	ResultUnknown Result = ""
)

type Team string

const (
	TeamCT      Team = "ct"
	TeamT       Team = "t"
	TeamUnknown Team = ""
)

type Stat struct {
	ID        int
	DemoID    int
	UserID    int
	Result    Result
	StartTeam Team
	Kills     int
	Assists   int
	Deaths    int
}

func StatModel(s sqlc.Stat) *Stat {
	result := ResultUnknown
	if s.Result.Valid {
		result = Result(s.Result.Result)
	}
	startTeam := TeamUnknown
	if s.StartTeam.Valid {
		startTeam = Team(s.StartTeam.Team)
	}

	return &Stat{
		ID:        int(s.ID),
		DemoID:    int(s.DemoID),
		UserID:    int(s.UserID),
		Result:    result,
		StartTeam: startTeam,
		Kills:     fromInt(s.Kills),
		Assists:   fromInt(s.Assists),
		Deaths:    fromInt(s.Deaths),
	}
}
