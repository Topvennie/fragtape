package api

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/server/service"
	"go.uber.org/zap"
)

type Demo struct {
	router fiber.Router
	demo   service.Demo
}

func NewDemo(router fiber.Router, service service.Service) *Demo {
	api := &Demo{
		router: router.Group("/demo"),
		demo:   *service.NewDemo(),
	}

	api.createRoutes()

	return api
}

func (d *Demo) createRoutes() {
	d.router.Get("/", d.getAll)
	d.router.Post("/upload", d.upload)
}

func (d *Demo) getAll(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	demos, err := d.demo.GetAll(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(demos)
}

func (d *Demo) upload(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	file, err := fileHeader.Open()
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	defer func() {
		// nolint:errcheck // Unlucky if it fails
		_ = file.Close()
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if err := d.demo.Upload(c.Context(), userID, data); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
