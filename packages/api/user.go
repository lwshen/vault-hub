package api

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (Server) GetCurrentUser(c echo.Context) error {
	userVal := c.Get("user")
	user, ok := userVal.(*model.User)
	if !ok {
		return handler.SendError(c, http.StatusUnauthorized, "user not found in context")
	}

	resp := GetUserResponse{
		Email:  openapi_types.Email(user.Email),
		Avatar: user.Avatar,
		Name:   user.Name,
	}

	return c.JSON(http.StatusOK, resp)
}
