package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
)

type SettingUser struct {
	repo Repository
}

func (r *Repository) NewSettingUser() *SettingUser {
	return &SettingUser{
		repo: *r,
	}
}

func (s *SettingUser) GetByUser(ctx context.Context, userID int) (*model.SettingUser, error) {
	setting, err := s.repo.queries(ctx).SettingUserGetByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user setting by user %d | %w", userID, err)
	}

	return model.SettingUserModel(setting), nil
}

func (s *SettingUser) Create(ctx context.Context, setting *model.SettingUser) error {
	id, err := s.repo.queries(ctx).SettingUserCreate(ctx, sqlc.SettingUserCreateParams{
		UserID:                   int32(setting.UserID),
		SteamMatchToken:          toString(setting.SteamMatchToken),
		SteamAuthenticationToken: toString(setting.SteamAuthenticationToken),
	})
	if err != nil {
		return fmt.Errorf("create user setting %+v | %w", *setting, err)
	}

	setting.ID = int(id)

	return nil
}

func (s *SettingUser) Update(ctx context.Context, setting model.SettingUser) error {
	if err := s.repo.queries(ctx).SettingUserUpdate(ctx, sqlc.SettingUserUpdateParams{
		ID:                       int32(setting.ID),
		SteamMatchToken:          toString(setting.SteamMatchToken),
		SteamAuthenticationToken: toString(setting.SteamAuthenticationToken),
	}); err != nil {
		return fmt.Errorf("update user setting %+v | %w", setting, err)
	}

	return nil
}
