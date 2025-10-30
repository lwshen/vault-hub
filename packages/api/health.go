package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (Server) Health(ctx echo.Context) error {
	status := "ok"
	time := time.Now()
	resp := HealthCheckResponse{
		Status:    &status,
		Timestamp: &time,
	}

	return ctx.JSON(http.StatusOK, resp)
}
