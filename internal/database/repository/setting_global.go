package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
)

type SettingGlobal struct {
	repo Repository
}

func (r *Repository) NewSettingGlobal() *SettingGlobal {
	return &SettingGlobal{
		repo: *r,
	}
}

func (s *SettingGlobal) Get(ctx context.Context) (*model.SettingGlobal, error) {
	setting, err := s.repo.queries(ctx).SettingGlobalGet(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no setting global found")
		}
		return nil, fmt.Errorf("get setting global %w", err)
	}

	return model.SettingGlobalModel(setting), nil
}

func (s *SettingGlobal) Update(ctx context.Context, setting model.SettingGlobal) error {
	if err := s.repo.queries(ctx).SettingGlobalUpdate(ctx, sqlc.SettingGlobalUpdateParams{
		DemoUpload:     toBool(&setting.DemoUpload),
		CustomCriteria: toBool(&setting.CustomCriteria),
		ChatCommand:    toBool(&setting.ChatCommand),
		ChatTrigger:    toString(setting.ChatTrigger),
	}); err != nil {
		return fmt.Errorf("update setting global %+v | %w", setting, err)
	}

	return nil
}
