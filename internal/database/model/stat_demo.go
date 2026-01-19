package model

import "github.com/topvennie/fragtape/pkg/sqlc"

type StatsDemo struct {
	ID       int
	DemoID   int
	Map      string
	RoundsCT int
	RoundsT  int
}

func StatsDemoModel(s sqlc.StatsDemo) *StatsDemo {
	return &StatsDemo{
		ID:       int(s.ID),
		DemoID:   int(s.DemoID),
		Map:      s.Map,
		RoundsCT: int(s.RoundsCt),
		RoundsT:  int(s.RoundsT),
	}
}
