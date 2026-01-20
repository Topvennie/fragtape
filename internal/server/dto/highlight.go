package dto

import (
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
)

type Highlight struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Round     int    `json:"round"`
	DurationS int    `json:"duration_s"`
}

func HighlightDTO(h *model.Highlight) Highlight {
	return Highlight{
		ID:        h.ID,
		Title:     h.Title,
		Round:     h.Round,
		DurationS: int(h.Duration / time.Second),
	}
}
