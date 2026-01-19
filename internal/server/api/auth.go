// Package api contains all api routes
package api

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/steam"
	"github.com/shareed2k/goth_fiber"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/internal/server/service"
	"github.com/topvennie/fragtape/pkg/config"
	"go.uber.org/zap"
)

type Auth struct {
	router fiber.Router

	user service.User

	redirectURL string
}

func NewAuth(router fiber.Router, service service.Service) *Auth {
	goth.UseProviders(
		steam.New(
			config.GetString("server.auth.steam.api_key"),
			config.GetString("server.auth.steam.callback_url"),
		),
	)

	api := &Auth{
		router:      router.Group("/auth"),
		user:        *service.NewUser(),
		redirectURL: config.GetDefaultString("server.auth.redirect_url", "/"),
	}

	api.routes()

	return api
}

func (r *Auth) routes() {
	r.router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	r.router.Get("/callback/:provider", r.loginCallback)
	r.router.Post("/logout", r.logout)
}

func (r *Auth) loginCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		zap.S().Errorf("Failed to complete user auth %v", err)
		return fiber.ErrInternalServerError
	}

	userID, err := strconv.Atoi(user.UserID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	dtoUser, err := r.user.GetByUID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {

			// New user
			dtoUser = dto.User{
				UID:         userID,
				Name:        user.Name,
				DisplayName: user.NickName,
				AvatarURL:   user.AvatarURL,
			}

			dtoUser, err = r.user.Create(c.Context(), dtoUser)
		}
		if err != nil {
			return err
		}
	}
	if dtoUser.Name != user.Name {
		dtoUser.Name = user.Name
		dtoUser.DisplayName = user.NickName
		dtoUser.AvatarURL = user.AvatarURL

		if dtoUser, err = r.user.Update(c.Context(), dtoUser); err != nil {
			return err
		}
	}

	if err := storeInSession(c, "userID", dtoUser.ID); err != nil {
		zap.S().Errorf("Failed to store user id in session %v", err)
		return fiber.ErrInternalServerError
	}
	if err = storeInSession(c, "steamUID", dtoUser.UID); err != nil {
		zap.S().Errorf("Failed to store steam id in session %v", err)
		return fiber.ErrInternalServerError
	}

	return c.Redirect(r.redirectURL)
}

func (r *Auth) logout(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		zap.S().Errorf("Failed to logout %v", err)
	}

	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("Failed to get session %v", err)
		return fiber.ErrInternalServerError
	}
	if err := session.Destroy(); err != nil {
		zap.S().Errorf("Failed to destroy %v", err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
