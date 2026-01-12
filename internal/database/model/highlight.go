package model

import (
	"time"

	"github.com/topvennie/fragtape/pkg/sqlc"
)

type Highlight struct {
	ID        int
	DemoID    int
	FileID    string
	FileWebID string
	Title     string
	CreatedAt time.Time
}

func HighlightModel(h sqlc.Highlight) *Highlight {
	return &Highlight{
		ID:        int(h.ID),
		DemoID:    int(h.DemoID),
		FileID:    fromString(h.FileID),
		FileWebID: fromString(h.FileWebID),
		Title:     h.Title,
		CreatedAt: h.CreatedAt.Time,
	}
}
