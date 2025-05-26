package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (Server) Health(ctx *fiber.Ctx) error {
	status := "ok"
	time := time.Now()
	resp := HealthCheckResponse{
		Status:    &status,
		Timestamp: &time,
	}

	return ctx.
		Status(http.StatusOK).
		JSON(resp)
}
