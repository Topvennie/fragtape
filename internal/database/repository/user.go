package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
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
