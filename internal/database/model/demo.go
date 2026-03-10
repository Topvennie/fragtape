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
	DemoStatusQueuedDownload DemoStatus = "queued_download"
	DemoStatusDownloading    DemoStatus = "downloading"
	DemoStatusQueuedParse    DemoStatus = "queued_parse"
	DemoStatusParsing        DemoStatus = "parsing"
	DemoStatusQueuedRender   DemoStatus = "queued_render"
	DemoStatusRendering      DemoStatus = "rendering"
	DemoStatusQueuedFinalize DemoStatus = "queued_finalize"
	DemoStatusFinalizing     DemoStatus = "finalizing"
	DemoStatusFinished       DemoStatus = "finished"
	DemoStatusFailed         DemoStatus = "failed"
)

type Demo struct {
	ID              int
	Source          DemoSource
	SourceID        string
	SourceURL       string
	Status          DemoStatus
	FileID          string
	DataID          string
	Attempts        int
	Error           string
	PlayedAt        time.Time
	StatusUpdatedAt time.Time
	CreatedAt       time.Time

	// Non db fields
	Stats []Stat
}

func DemoModel(d sqlc.Demo) *Demo {
	return &Demo{
		ID:              int(d.ID),
		Source:          DemoSource(d.Source),
		SourceID:        d.SourceID,
		SourceURL:       fromString(d.SourceUrl),
		Status:          DemoStatus(d.Status),
		FileID:          fromString(d.FileID),
		DataID:          fromString(d.DataID),
		Attempts:        int(d.Attempts),
		Error:           fromString(d.Error),
		PlayedAt:        fromTime(d.PlayedAt),
		StatusUpdatedAt: d.StatusUpdatedAt.Time,
		CreatedAt:       d.CreatedAt.Time,
	}
}

type DemoFilterResult struct {
	Demos []Demo
	Total int
}

type DemoFilter struct {
	UserID        int
	Source        *DemoSource
	Result        *Result
	PlayedAtStart time.Time
	PlayedAtEnd   time.Time
	HasHighlight  *bool
	Limit         int
	Offset        int
}
