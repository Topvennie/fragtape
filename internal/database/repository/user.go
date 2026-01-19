package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
	"github.com/topvennie/fragtape/pkg/utils"
)

type User struct {
	repo Repository
}

func (r *Repository) NewUser() *User {
	return &User{
		repo: *r,
	}
}

func (u *User) Get(ctx context.Context, id int) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGet(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with id %d | %w", id, err)
	}

	return model.UserModel(user), nil
}

func (u *User) GetByUID(ctx context.Context, uid int) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGetByUid(ctx, int32(uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with uid %d | %w", uid, err)
	}

	return model.UserModel(user), nil
}

func (u *User) GetByIDs(ctx context.Context, ids []int) ([]*model.User, error) {
	users, err := u.repo.queries(ctx).UserGetByIds(ctx, utils.SliceMap(ids, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get users with id %+v | %w", ids, err)
	}

	return utils.SliceMap(users, model.UserModel), nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	id, err := u.repo.queries(ctx).UserCreate(ctx, sqlc.UserCreateParams{
		Uid:         int32(user.UID),
		Name:        toString(user.Name),
		DisplayName: user.DisplayName,
		AvatarUrl:   toString(user.AvatarURL),
		Crosshair:   toString(user.Crosshair),
	})
	if err != nil {
		return fmt.Errorf("create user %+v | %w", *user, err)
	}

	user.ID = int(id)

	return nil
}

func (u *User) Update(ctx context.Context, user model.User) error {
	if err := u.repo.queries(ctx).UserUpdate(ctx, sqlc.UserUpdateParams{
		ID:          int32(user.ID),
		Name:        toString(user.Name),
		DisplayName: user.DisplayName,
		AvatarUrl:   toString(user.AvatarURL),
		Crosshair:   toString(user.Crosshair),
	}); err != nil {
		return fmt.Errorf("update user %+v | %w", user, err)
	}

	return nil
}
