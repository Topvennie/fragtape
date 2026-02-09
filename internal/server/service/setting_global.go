package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"go.uber.org/zap"
)

type SettingGlobal struct {
	setting repository.SettingGlobal
	user    repository.User
}

func (s *Service) NewSettingGlobal() *SettingGlobal {
	return &SettingGlobal{
		setting: *s.repo.NewSettingGlobal(),
		user:    *s.repo.NewUser(),
	}
}

func (s *SettingGlobal) Get(ctx context.Context) (dto.SettingGlobal, error) {
	setting, err := s.setting.Get(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.SettingGlobal{}, fiber.ErrInternalServerError
	}

	return dto.SettingGlobalDTO(setting), nil
}

func (s *SettingGlobal) Update(ctx context.Context, setting dto.SettingGlobal, userID int) error {
	user, err := s.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if !user.Admin {
		return fiber.ErrForbidden
	}

	if err := s.setting.Update(ctx, *setting.ToModel()); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
