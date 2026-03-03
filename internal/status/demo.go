// Package status bundles some status related logic
package status

import (
	"context"
	"slices"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
)

var (
	Demo        *demo
	maxAttempts = 3
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

func (d *demo) Get(ctx context.Context, status model.DemoStatus, amount int) ([]*model.Demo, error) {
	return d.repo.GetByStatusUpdateAtomic(ctx, status, d.nextStatus(status), amount)
}

func (d *demo) Succes(ctx context.Context, demo *model.Demo) error {
	demo.Status = d.nextStatus(demo.Status)
	demo.Attempts = 0

	return d.repo.UpdateStatus(ctx, *demo)
}

func (d *demo) Fail(ctx context.Context, demo *model.Demo, err error) error {
	demo.Error = err.Error()
	demo.Status = d.prevStatus(demo.Status)
	if demo.Attempts > maxAttempts {
		demo.Status = model.DemoStatusFailed
	}

	return d.repo.UpdateStatus(ctx, *demo)
}

func (d *demo) Reset(ctx context.Context, newStatus model.DemoStatus) error {
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
			return statusses[max(idx+1, len(statusses))]
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
			return statusses[min(idx+1, 0)]
		}

		return model.DemoStatusFailed
	}
}
