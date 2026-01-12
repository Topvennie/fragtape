package model

import (
	"time"

	"github.com/topvennie/fragtape/pkg/sqlc"
)

type DemoSource string

const (
	DemoSourceManual DemoSource = "manual"
	DemoSourceSteam  DemoSource = "steam"
	DemoSourceFaceit DemoSource = "faceit"
)

type DemoStatus string

const (
	DemoStatusQueuedParse  DemoStatus = "queued_parse"
	DemoStatusParsing      DemoStatus = "parsing"
	DemoStatusQueuedRender DemoStatus = "queued_render"
	DemoStatusRendering    DemoStatus = "rendering"
	DemoStatusRendered     DemoStatus = "rendered"
	DemoStatusCompleted    DemoStatus = "completed"
	DemoStatusFailed       DemoStatus = "failed"
)

type Demo struct {
	ID              int
	UserID          int
	Source          DemoSource
	SourceID        string
	Status          DemoStatus
	FileID          string
	Attempts        int
	Error           string
	CreatedAt       time.Time
	StatusUpdatedAt time.Time
	DeletedAt       time.Time
}

func DemoModel(d sqlc.Demo) *Demo {
	return &Demo{
		ID:              int(d.ID),
		UserID:          int(d.UserID),
		Source:          DemoSource(d.Source),
		SourceID:        fromString(d.SourceID),
		Status:          DemoStatus(d.Status),
		FileID:          fromString(d.FileID),
		Attempts:        int(d.Attempts),
		Error:           fromString(d.Error),
		CreatedAt:       d.CreatedAt.Time,
		StatusUpdatedAt: d.StatusUpdatedAt.Time,
		DeletedAt:       fromTime(d.DeletedAt),
	}
}
