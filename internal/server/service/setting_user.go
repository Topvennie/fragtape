package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"go.uber.org/zap"
)

type SettingUser struct {
	setting repository.SettingUser
}

func (s *Service) NewSettingUser() *SettingUser {
	return &SettingUser{
		setting: *s.repo.NewSettingUser(),
	}
}

func (s *SettingUser) GetByUser(ctx context.Context, userID int) (dto.SettingUser, error) {
	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return dto.SettingUser{}, fiber.ErrInternalServerError
	}
	if setting == nil {
		return dto.SettingUser{}, fiber.ErrNotFound
	}

	return dto.SettingUserDTO(setting), nil
}

func (s *SettingUser) CreateIfNotExist(ctx context.Context, userID int) error {
	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if setting != nil {
		return nil
	}

	setting = &model.SettingUser{
		UserID: userID,
	}

	if err := s.setting.Create(ctx, setting); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *SettingUser) SteamConnect(ctx context.Context, steam dto.SettingUserSteam, userID int) error {
	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if setting == nil {
		return fiber.ErrInternalServerError
	}

	setting.SteamMatchToken = steam.MatchToken
	setting.SteamAuthenticationToken = steam.AuthenticationToken

	if err := s.setting.Update(ctx, *setting); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *SettingUser) SteamDisconnect(ctx context.Context, userID int) error {
	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if setting == nil {
		return fiber.ErrInternalServerError
	}

	setting.SteamMatchToken = ""
	setting.SteamAuthenticationToken = ""

	if err := s.setting.Update(ctx, *setting); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
