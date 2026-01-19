package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type Result string

const (
	ResultWin  Result = "win"
	ResultLoss Result = "loss"
	ResultTie  Result = "tie"
)

type Team string

const (
	TeamCT = "ct"
	TeamT  = "t"
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
	return &Stat{
		ID:        int(s.ID),
		DemoID:    int(s.DemoID),
		UserID:    int(s.UserID),
		Result:    Result(s.Result),
		StartTeam: Team(s.StartTeam),
		Kills:     int(s.Kills),
		Assists:   int(s.Assists),
		Deaths:    int(s.Deaths),
	}
}
