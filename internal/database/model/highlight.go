package model

import (
	"time"

	"github.com/topvennie/fragtape/pkg/sqlc"
)

type Highlight struct {
	ID        int
	UserID    int
	DemoID    int
	FileID    string
	FileWebID string
	Title     string
	CreatedAt time.Time

	// Non db fields
	Segments []HighlightSegment
}

func HighlightModel(h sqlc.Highlight) *Highlight {
	return &Highlight{
		ID:        int(h.ID),
		DemoID:    int(h.DemoID),
		UserID:    int(h.UserID),
		FileID:    fromString(h.FileID),
		FileWebID: fromString(h.FileWebID),
		Title:     h.Title,
		CreatedAt: h.CreatedAt.Time,
	}
}

type HighlightSegment struct {
	ID          int
	HighlightID int
	StartTick   int
	EndTick     int
}

func HighlightSegmentModel(h sqlc.HighlightSegment) *HighlightSegment {
	return &HighlightSegment{
		ID:          int(h.ID),
		HighlightID: int(h.HighlightID),
		StartTick:   int(h.StartTick),
		EndTick:     int(h.EndTick),
	}
}
