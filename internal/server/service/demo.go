package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

type Demo struct {
	demo repository.Demo
}

func (s *Service) NewDemo() *Demo {
	return &Demo{
		demo: *s.repo.NewDemo(),
	}
}

func (d *Demo) GetAll(ctx context.Context, userID int) ([]dto.Demo, error) {
	demos, err := d.demo.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(demos, dto.DemoDTO), nil
}

func (d *Demo) Upload(ctx context.Context, userID int, file []byte) error {
	demo := &model.Demo{
		UserID: userID,
		Source: model.DemoSourceManual,
		FileID: uuid.NewString(),
	}

	if err := storage.S.Set(demo.FileID, file, 0); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if err := d.demo.Create(ctx, demo); err != nil {
		zap.S().Error(err)

		if err := storage.S.Delete(demo.FileID); err != nil {
			zap.S().Error(err)
		}

		return fiber.ErrInternalServerError
	}

	return nil
}

func (d *Demo) Delete(ctx context.Context, userID, demoID int) error {
	demo, err := d.demo.Get(ctx, demoID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if demo.UserID != userID {
		return fiber.ErrForbidden
	}
	if !demo.DeletedAt.IsZero() {
		return fiber.ErrBadRequest
	}

	if err := d.demo.Delete(ctx, demoID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
