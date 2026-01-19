package dto

import "github.com/topvennie/fragtape/internal/database/model"

type StatsDemo struct {
	Map      string `json:"map"`
	RoundsCT int    `json:"rounds_ct"`
	RoundsT  int    `json:"rounds_t"`
}

func StatsDemoDTO(s *model.StatsDemo) StatsDemo {
	return StatsDemo{
		Map:      s.Map,
		RoundsCT: s.RoundsCT,
		RoundsT:  s.RoundsT,
	}
}
