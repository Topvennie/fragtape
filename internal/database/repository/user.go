package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/sqlc"
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

func (u *User) GetByUID(ctx context.Context, uid string) (*model.User, error) {
	user, err := u.repo.queries(ctx).UserGetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user with uid %s | %w", uid, err)
	}

	return model.UserModel(user), nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	id, err := u.repo.queries(ctx).UserCreate(ctx, sqlc.UserCreateParams{
		Uid:         user.UID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		AvatarUrl:   user.AvatarURL,
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
		Name:        user.Name,
		DisplayName: user.DisplayName,
		AvatarUrl:   user.AvatarURL,
	}); err != nil {
		return fmt.Errorf("update user %+v | %w", user, err)
	}

	return nil
}
