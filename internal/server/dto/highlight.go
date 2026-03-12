package dto

import (
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
)

type Highlight struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Round     int    `json:"round"`
	Kills     int    `json:"kills"`
	DurationS int    `json:"duration_s"`
	Generated bool   `json:"generated"`
}

func HighlightDTO(h *model.Highlight) Highlight {
	return Highlight{
		ID:        h.ID,
		Title:     h.Title,
		Round:     h.Round,
		Kills:     h.Kills,
		DurationS: int(h.Duration / time.Second),
		Generated: h.FileID != "",
	}
}
