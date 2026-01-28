package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/internal/server/service"
)

type User struct {
	router fiber.Router

	user service.User
}

func NewUser(router fiber.Router, service service.Service) *User {
	api := &User{
		router: router.Group("/user"),
		user:   *service.NewUser(),
	}

	api.routes()

	return api
}

func (u *User) routes() {
	u.router.Get("/me", u.getMe)
	u.router.Get("/filtered", u.getFiltered)
	u.router.Get("/admin", u.getAmin)
	u.router.Post("/admin/:id", u.createAdmin)
	u.router.Delete("/admin/:id", u.deleteAdmin)
}

func (u *User) getMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, err := u.user.Get(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (u *User) getFiltered(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	name := c.Query("name", "")

	var real *bool
	if v := c.Query("real"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			real = &b
		}
	}

	var admin *bool
	if v := c.Query("admin"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			admin = &b
		}
	}

	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	if limit < 1 || page < 1 {
		return fiber.ErrBadRequest
	}

	users, err := u.user.GetFiltered(c.Context(), userID, dto.UserFilter{
		Name:   name,
		Admin:  admin,
		Real:   real,
		Limit:  limit,
		Offset: (page - 1) * limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (u *User) getAmin(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	admins, err := u.user.GetAdmin(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(admins)
}

func (u *User) createAdmin(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	adminID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := u.user.CreateAdmin(c.Context(), userID, adminID); err != nil {
		return err
	}

	return nil
}

func (u *User) deleteAdmin(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	adminID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := u.user.DeleteAdmin(c.Context(), userID, adminID); err != nil {
		return err
	}

	return nil
}
