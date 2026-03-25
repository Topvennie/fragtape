// Package status bundles some status related logic
package status

import (
	"context"
	"slices"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

var (
	Demo        *demo
	MaxAttempts = 3
	statusses   = []model.DemoStatus{
		model.DemoStatusQueuedDownload,
		model.DemoStatusDownloading,
		model.DemoStatusQueuedParse,
		model.DemoStatusParsing,
		model.DemoStatusQueuedRender,
		model.DemoStatusRendering,
		model.DemoStatusQueuedFinalize,
		model.DemoStatusFinalizing,
		model.DemoStatusFinished,
		model.DemoStatusFailed,
	}
)

type demo struct {
	repo repository.Demo
}

func Init(repo repository.Repository) {
	Demo = &demo{
		repo: *repo.NewDemo(),
	}
}

func (d *demo) Get(ctx context.Context, status model.DemoStatus) (*model.Demo, error) {
	demos, err := d.repo.GetByStatusUpdateAtomic(ctx, status, d.nextStatus(status), 1)
	if err != nil {
		return nil, err
	}
	if len(demos) == 0 {
		return nil, nil
	}

	zap.S().Infof("Demo %d is now %s", demos[0].ID, demos[0].Status)

	return demos[0], nil
}

func (d *demo) Succes(ctx context.Context, demo *model.Demo) error {
	zap.S().Infof("Demo %d is going to the next status %s", demo.ID, d.nextStatus(demo.Status))
	demo.Status = d.nextStatus(demo.Status)
	demo.Attempts = 0

	return d.repo.UpdateStatus(ctx, *demo)
}

func (d *demo) Fail(ctx context.Context, demo *model.Demo, err error) error {
	zap.S().Warnf("Demo %+v failed %v", *demo, err)
	demo.Error = err.Error()
	demo.Status = d.prevStatus(demo.Status)
	if demo.Attempts > MaxAttempts {
		demo.Status = model.DemoStatusFailed

		if demo.FileID != "" {
			// Best effort
			_ = storage.S.Delete(demo.FileID)
			demo.FileID = ""
			_ = d.repo.UpdateFile(ctx, *demo)
		}
	}

	return d.repo.UpdateStatus(ctx, *demo)
}

func (d *demo) Reset(ctx context.Context, newStatus model.DemoStatus) error {
	zap.S().Infof("Resetting demo statusses from %s to %s", d.nextStatus(newStatus), newStatus)
	return d.repo.ResetStatusAll(ctx, d.nextStatus(newStatus), newStatus)
}

func (d *demo) nextStatus(status model.DemoStatus) model.DemoStatus {
	switch status {
	case model.DemoStatusFinished:
		return model.DemoStatusFinished
	case model.DemoStatusFailed:
		return model.DemoStatusFailed
	default:
		if idx := slices.Index(statusses, status); idx != -1 {
			return statusses[min(idx+1, len(statusses)-1)]
		}

		return model.DemoStatusFailed
	}
}

func (d *demo) prevStatus(status model.DemoStatus) model.DemoStatus {
	switch status {
	case model.DemoStatusFinished:
		return model.DemoStatusFinished
	case model.DemoStatusFailed:
		return model.DemoStatusFailed
	default:
		if idx := slices.Index(statusses, status); idx != -1 {
			return statusses[max(idx-1, 0)]
		}

		return model.DemoStatusFailed
	}
}
