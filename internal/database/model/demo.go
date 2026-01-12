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
	Source          DemoSource
	SourceID        string
	Status          DemoStatus
	FileID          string
	Attempts        int
	Error           string
	StatusUpdatedAt time.Time
	CreatedAt       time.Time
}

func DemoModel(d sqlc.Demo) *Demo {
	return &Demo{
		ID:              int(d.ID),
		Source:          DemoSource(d.Source),
		SourceID:        fromString(d.SourceID),
		Status:          DemoStatus(d.Status),
		FileID:          fromString(d.FileID),
		Attempts:        int(d.Attempts),
		Error:           fromString(d.Error),
		StatusUpdatedAt: d.StatusUpdatedAt.Time,
		CreatedAt:       d.CreatedAt.Time,
	}
}

type DemoUser struct {
	ID        int
	DemoID    int
	UserID    int
	DeletedAt time.Time
}

func DemoUserModel(d sqlc.DemoUser) *DemoUser {
	return &DemoUser{
		ID:        int(d.ID),
		DemoID:    int(d.DemoID),
		UserID:    int(d.UserID),
		DeletedAt: fromTime(d.DeletedAt),
	}
}
