package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/version"
)

func (Server) GetVersion(ctx *fiber.Ctx) error {
	resp := VersionResponse{
		Version: version.Version,
		Commit:  version.Commit,
	}

	return ctx.
		Status(http.StatusOK).
		JSON(resp)
}
