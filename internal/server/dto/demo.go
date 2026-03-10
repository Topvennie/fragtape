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
	PlayedAt        time.Time        `json:"played_at"`
	StatusUpdatedAt time.Time        `json:"status_updated_at"`
}

func DemoDTO(d *model.Demo) Demo {
	playedAt := d.CreatedAt
	if !d.PlayedAt.IsZero() {
		playedAt = d.PlayedAt
	}

	return Demo{
		ID:              d.ID,
		Source:          d.Source,
		Players:         []DemoPlayer{},
		Status:          d.Status,
		PlayedAt:        playedAt,
		StatusUpdatedAt: d.StatusUpdatedAt,
	}
}

type DemoPlayer struct {
	User `json:"user"`
	Stat `json:"stat"`

	Highlights []Highlight `json:"highlights,omitzero"`
}

type DemoFilterResult struct {
	Demos []Demo `json:"demos"`
	Total int    `json:"total"`
}

type DemoFilter struct {
	UserID int
	Source *model.DemoSource
	Limit  int
	Offset int
}

func (d *DemoFilter) ToModel() *model.DemoFilter {
	return &model.DemoFilter{
		UserID: d.UserID,
		Source: d.Source,
		Limit:  d.Limit,
		Offset: d.Offset,
	}
}
