package dto

import "github.com/topvennie/fragtape/internal/database/model"

type Stat struct {
	Kills   int `json:"kills"`
	Assists int `json:"assists"`
	Deaths  int `json:"deaths"`
}

func StatDTO(s *model.Stat) Stat {
	return Stat{
		Kills:   s.Kills,
		Assists: s.Assists,
		Deaths:  s.Deaths,
	}
}
