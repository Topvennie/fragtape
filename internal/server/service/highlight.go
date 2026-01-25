package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/storage"
	"go.uber.org/zap"
)

type Highlight struct {
	highlight repository.Highlight
}

func (s *Service) NewHighlight() *Highlight {
	return &Highlight{
		highlight: *s.repo.NewHighlight(),
	}
}

func (h *Highlight) GetVideo(ctx context.Context, id int) ([]byte, error) {
	highlight, err := h.highlight.Get(ctx, id)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if highlight == nil {
		return nil, fiber.ErrNotFound
	}
	if highlight.FileID == "" {
		return nil, fiber.ErrBadRequest
	}

	video, err := storage.S.Get(highlight.FileID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return video, nil
}
