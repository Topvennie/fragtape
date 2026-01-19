package dto

import (
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
)

type Demo struct {
	ID              int              `json:"id"`
	Source          model.DemoSource `json:"source"`
	Status          model.DemoStatus `json:"status"`
	Players         []DemoPlayer     `json:"players"`
	Stats           StatsDemo        `json:"stats"`
	CreatedAt       time.Time        `json:"created_at"`
	StatusUpdatedAt time.Time        `json:"status_updated_at"`
}

func DemoDTO(d *model.Demo) Demo {
	return Demo{
		ID:              d.ID,
		Source:          d.Source,
		Players:         []DemoPlayer{},
		Status:          d.Status,
		CreatedAt:       d.CreatedAt,
		StatusUpdatedAt: d.StatusUpdatedAt,
	}
}

type DemoPlayer struct {
	User `json:"user"`
	Stat `json:"stat"`

	Highlights []Highlight `json:"highlights,omitzero"`
}
