package api

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/fragtape/internal/server/service"
)

type Highlight struct {
	router fiber.Router

	highlight service.Highlight
}

func NewHighlight(router fiber.Router, service service.Service) *Highlight {
	api := &Highlight{
		router:    router.Group("/highlight"),
		highlight: *service.NewHighlight(),
	}

	api.createRoutes()

	return api
}

func (h *Highlight) createRoutes() {
	h.router.Get("/video/:id", h.getVideo)
}

func (h *Highlight) getVideo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "no id found")
	}

	data, err := h.highlight.GetVideo(c.Context(), id)
	if err != nil {
		return err
	}

	totalSize := int64(len(data))

	hash := md5.Sum(data)
	etag := fmt.Sprintf(`W/"%x-%x"`, hash, totalSize)

	c.Set("Cache-Control", "public, max-age=3600, stale-while-revalidate=120")
	c.Set("ETag", etag)

	rangeHeader := c.Get("Range")

	if rangeHeader == "" {
		if match := c.Get("If-None-Match"); match == etag {
			return c.SendStatus(fiber.StatusNotModified)
		}

		c.Set("Content-Type", MimeWEBM)
		return c.Send(data)
	}

	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeHeader, "-")

	start, _ := strconv.ParseInt(parts[0], 10, 64)
	end := totalSize - 1

	if len(parts) > 1 && parts[1] != "" {
		end, _ = strconv.ParseInt(parts[1], 10, 64)
	}

	if start > end || start >= totalSize {
		c.Set("Content-Range", fmt.Sprintf("bytes */%d", totalSize))
		return c.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
	}

	if end >= totalSize {
		end = totalSize - 1
	}

	chunkSize := (end - start) + 1

	c.Status(fiber.StatusPartialContent)
	c.Set("Content-Type", MimeWEBM)
	c.Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	c.Set("Accept-Ranges", "bytes")
	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))

	return c.Send(data[start : end+1])
}
