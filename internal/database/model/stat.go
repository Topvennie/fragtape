package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type Stat struct {
	ID      int
	DemoID  int
	UserID  int
	Kills   int
	Assists int
	Deaths  int
}

func StatModel(s sqlc.Stat) *Stat {
	return &Stat{
		ID:      int(s.ID),
		DemoID:  int(s.DemoID),
		UserID:  int(s.UserID),
		Kills:   int(s.Kills),
		Assists: int(s.Assists),
		Deaths:  int(s.Deaths),
	}
}
