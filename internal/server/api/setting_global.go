package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/internal/server/service"
)

type SettingGlobal struct {
	router fiber.Router

	setting service.SettingGlobal
}

func NewSettingGlobal(router fiber.Router, service service.Service) *SettingGlobal {
	api := &SettingGlobal{
		router:  router.Group("/setting/global"),
		setting: *service.NewSettingGlobal(),
	}

	api.createRoutes()

	return api
}

func (s *SettingGlobal) createRoutes() {
	s.router.Get("/", s.get)
	s.router.Post("/", s.update)
}

func (s *SettingGlobal) get(c *fiber.Ctx) error {
	setting, err := s.setting.Get(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(setting)
}

func (s *SettingGlobal) update(c *fiber.Ctx) error {
	var setting dto.SettingGlobal
	if err := c.BodyParser(&setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := dto.Validate.Struct(setting); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := s.setting.Update(c.Context(), setting, userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
