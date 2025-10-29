package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck builds the response payload used by health endpoints.
func HealthCheck() HealthCheckResponse {
	status := "ok"
	now := time.Now()
	return HealthCheckResponse{
		Status:    &status,
		Timestamp: &now,
	}
}

func (Server) Health(ctx *fiber.Ctx) error {
	resp := HealthCheck()
	return ctx.
		Status(http.StatusOK).
		JSON(resp)
}
