package service

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/pkg/faceit"
	"github.com/topvennie/fragtape/pkg/steam"
	"go.uber.org/zap"
)

type SettingUser struct {
	setting repository.SettingUser
	user    repository.User
}

func (s *Service) NewSettingUser() *SettingUser {
	return &SettingUser{
		setting: *s.repo.NewSettingUser(),
		user:    *s.repo.NewUser(),
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

func (s *SettingUser) SteamConnect(ctx context.Context, settingSteam dto.SettingUserSteam, userID int) error {
	user, err := s.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if setting == nil {
		return fiber.ErrInternalServerError
	}
	if setting.SteamMatchToken != "" || setting.SteamAuthenticationToken != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Steam is already configured")
	}

	// Verify the credentials
	nextDemo, err := steam.S.NextDemo(ctx, steam.NextDemoParams{
		SteamID:                  user.UID,
		SteamAuthenticationToken: settingSteam.AuthenticationToken,
		SteamMatchToken:          settingSteam.MatchToken,
	})
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if nextDemo.Error != nil {
		if nextDemo.Code == http.StatusForbidden || nextDemo.Code == http.StatusPreconditionFailed {
			// Invalid credentials
			return fiber.NewError(fiber.StatusBadRequest, "invalid steam credentials")
		}
	}

	// Credentials are verified

	setting.SteamMatchToken = settingSteam.MatchToken
	setting.SteamAuthenticationToken = settingSteam.AuthenticationToken

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

func (s *SettingUser) FaceitConnect(ctx context.Context, userID int) error {
	user, err := s.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if setting == nil {
		return fiber.ErrInternalServerError
	}
	if setting.FaceitID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Faceit is already configured")
	}

	// Get faceit id
	faceitID, err := faceit.F.GetUserID(ctx, user.UID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if faceitID == "" {
		return fiber.ErrNotFound
	}

	setting.FaceitID = faceitID
	if err := s.setting.Update(ctx, *setting); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *SettingUser) FaceitDisconnect(ctx context.Context, userID int) error {
	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if setting == nil {
		return fiber.ErrInternalServerError
	}

	setting.FaceitID = ""

	if err := s.setting.Update(ctx, *setting); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *SettingUser) FirstTimeWizard(ctx context.Context, userID int, firstTimeWizard dto.SettingUserFirsTimeWizard) error {
	setting, err := s.setting.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if setting == nil {
		return fiber.ErrInternalServerError
	}

	setting.FirstTimeWizard = *firstTimeWizard.FirsTimeWizard

	if err := s.setting.Update(ctx, *setting); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
