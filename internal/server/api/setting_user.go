package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/internal/server/service"
)

type SettingUser struct {
	router fiber.Router

	setting service.SettingUser
}

func NewSettingUser(router fiber.Router, service service.Service) *SettingUser {
	api := &SettingUser{
		router:  router.Group("/setting/user"),
		setting: *service.NewSettingUser(),
	}

	api.createRoutes()

	return api
}

func (s *SettingUser) createRoutes() {
	s.router.Get("/", s.get)
	s.router.Post("/steam", s.steamConnect)
	s.router.Delete("/steam", s.steamDisconnect)
	s.router.Post("/faceit", s.faceitConnect)
	s.router.Delete("/faceit", s.faceitDisconnect)
}

func (s *SettingUser) get(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	setting, err := s.setting.GetByUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(setting)
}

func (s *SettingUser) steamConnect(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var steam dto.SettingUserSteam
	if err := c.BodyParser(&steam); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(steam); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.setting.SteamConnect(c.Context(), steam, userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *SettingUser) steamDisconnect(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := s.setting.SteamDisconnect(c.Context(), userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *SettingUser) faceitConnect(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := s.setting.FaceitConnect(c.Context(), userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *SettingUser) faceitDisconnect(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := s.setting.FaceitDisconnect(c.Context(), userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
