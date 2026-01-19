package dto

import "github.com/topvennie/fragtape/internal/database/model"

type Stat struct {
	Result    model.Result `json:"result"`
	StartTeam model.Team   `json:"start_team"`
	Kills     int          `json:"kills"`
	Assists   int          `json:"assists"`
	Deaths    int          `json:"deaths"`
}

func StatDTO(s *model.Stat) Stat {
	return Stat{
		Result:    s.Result,
		StartTeam: s.StartTeam,
		Kills:     s.Kills,
		Assists:   s.Assists,
		Deaths:    s.Deaths,
	}
}
