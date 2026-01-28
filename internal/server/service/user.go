package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/pkg/utils"
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

func (u *User) GetByUID(ctx context.Context, uid int) (dto.User, error) {
	user, err := u.user.GetByUID(ctx, uid)
	if err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}
	if user == nil {
		return dto.User{}, fiber.ErrNotFound
	}

	return dto.UserDTO(user), nil
}

func (u *User) GetAdmin(ctx context.Context, userID int) ([]dto.User, error) {
	user, err := u.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if user == nil {
		return nil, fiber.ErrBadRequest
	}
	if !user.Admin {
		return nil, fiber.ErrForbidden
	}

	admins, err := u.user.GetAdmin(ctx)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(admins, dto.UserDTO), nil
}

func (u *User) GetFiltered(ctx context.Context, userID int, filter dto.UserFilter) (dto.UserFilterResult, error) {
	user, err := u.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return dto.UserFilterResult{}, fiber.ErrInternalServerError
	}
	if user == nil {
		return dto.UserFilterResult{}, fiber.ErrBadRequest
	}
	if !user.Admin {
		return dto.UserFilterResult{}, fiber.ErrForbidden
	}

	filterModel := filter.ToModel()
	zap.S().Debugf("%+v", *filterModel)

	result, err := u.user.GetFiltered(ctx, *filterModel)
	if err != nil {
		zap.S().Error(err)
		return dto.UserFilterResult{}, fiber.ErrInternalServerError
	}
	if result == nil {
		return dto.UserFilterResult{Users: []dto.User{}}, nil
	}

	return dto.UserFilterResultDTO(result), nil
}

func (u *User) Create(ctx context.Context, userSave dto.User) (dto.User, error) {
	user := userSave.ToModel()

	if err := u.user.Create(ctx, user); err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}

	return dto.UserDTO(user), nil
}

func (u *User) CreateAdmin(ctx context.Context, userID, adminID int) error {
	user, err := u.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if !user.Admin {
		return fiber.ErrForbidden
	}

	admin, err := u.user.Get(ctx, adminID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if admin == nil {
		return fiber.ErrBadRequest
	}
	if admin.Admin {
		return fiber.ErrBadRequest
	}

	admin.Admin = true
	if err := u.user.Update(ctx, *admin); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *User) Update(ctx context.Context, userSave dto.User) (dto.User, error) {
	user := userSave.ToModel()

	if err := u.user.Update(ctx, *user); err != nil {
		zap.S().Error(err)
		return dto.User{}, fiber.ErrInternalServerError
	}

	return dto.UserDTO(user), nil
}

func (u *User) DeleteAdmin(ctx context.Context, userID, adminID int) error {
	user, err := u.user.Get(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if !user.Admin {
		return fiber.ErrForbidden
	}

	admin, err := u.user.Get(ctx, adminID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if admin == nil {
		return fiber.ErrBadRequest
	}
	if !admin.Admin {
		return fiber.ErrBadRequest
	}

	if user.ID == admin.ID {
		return fiber.ErrBadRequest
	}

	admin.Admin = false
	if err := u.user.Update(ctx, *admin); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
