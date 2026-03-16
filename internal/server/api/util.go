package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

const (
	MimeWEBP = "image/webp"
	MimeWEBM = "video/webm"
)

func storeInSession[T any](ctx *fiber.Ctx, key string, value T) error {
	session, err := goth_fiber.SessionStore.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get session %w", err)
	}

	session.Set(key, value)

	return session.Save()
}

// Uncomment when used
// func getFromSession[T any](ctx *fiber.Ctx, key string) (T, error) {
// 	session, err := goth_fiber.SessionStore.Get(ctx)
// 	if err != nil {
// 		return *new(T), fmt.Errorf("failed to get session %w", err)
// 	}
//
// 	value, ok := session.Get(key).(T)
// 	if !ok {
// 		return *new(T), fmt.Errorf("unable to convert session value to T %v", value)
// 	}
//
// 	return value, nil
// }
