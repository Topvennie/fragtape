package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"go.uber.org/zap"
)

type User struct {
	service Service

	user repository.User
}

func (s *Service) NewUser() *User {
	return &User{
		service: *s,
		user:    *s.repo.NewUser(),
	}
}

func (u *User) Get(ctx context.Context, id int) (dto.User, error) {
	user, err := u.user.Get(ctx, id)
	if err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(user), nil
}
