package echoapp

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
	api "github.com/lwshen/vault-hub/packages/api"
)

// RegisterRoutes wires Echo route groups with handlers that replicate the
// existing Fiber behavior.
func RegisterRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api")

	apiGroup.GET("/health", healthHandler)
	apiGroup.GET("/config", configHandler)
	apiGroup.GET("/status", statusHandler)
	apiGroup.GET("/user", currentUserHandler)
	apiGroup.GET("/vaults", getVaultsHandler)
	apiGroup.GET("/vaults/:uniqueId", getVaultHandler)
	apiGroup.POST("/vaults", createVaultHandler)
	apiGroup.PUT("/vaults/:uniqueId", updateVaultHandler)
	apiGroup.DELETE("/vaults/:uniqueId", deleteVaultHandler)
	apiGroup.GET("/api-keys", getAPIKeysHandler)
	apiGroup.POST("/api-keys", createAPIKeyHandler)
	apiGroup.PATCH("/api-keys/:id", updateAPIKeyHandler)
	apiGroup.DELETE("/api-keys/:id", deleteAPIKeyHandler)
}

func healthHandler(c echo.Context) error {
	resp := api.HealthCheck()
	return c.JSON(http.StatusOK, resp)
}

func configHandler(c echo.Context) error {
	resp := api.PublicConfig()
	return c.JSON(http.StatusOK, resp)
}

func statusHandler(c echo.Context) error {
	resp := api.BuildStatusResponse()
	return c.JSON(http.StatusOK, resp)
}

func currentUserHandler(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "user not found in context")
	}

	resp, apiErr := api.BuildCurrentUserResponse(user)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func getVaultsHandler(c echo.Context) error {
	params, apiErr := parseVaultQueryParams(c)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	resp, err := api.GetVaultsForUser(user, params)
	if err != nil {
		return sendError(c, err.Status, err.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func getVaultHandler(c echo.Context) error {
	uniqueID := c.Param("uniqueId")
	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	resp, err := api.GetVaultForUser(user, uniqueID, clientInfo)
	if err != nil {
		return sendError(c, err.Status, err.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func createVaultHandler(c echo.Context) error {
	var input api.CreateVaultRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	resp, apiErr := api.CreateVaultForUser(user, input, clientInfo)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusCreated, resp)
}

func updateVaultHandler(c echo.Context) error {
	uniqueID := c.Param("uniqueId")

	var input api.UpdateVaultRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	resp, apiErr := api.UpdateVaultForUser(user, uniqueID, input, clientInfo)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func deleteVaultHandler(c echo.Context) error {
	uniqueID := c.Param("uniqueId")
	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	if apiErr := api.DeleteVaultForUser(user, uniqueID, clientInfo); apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.NoContent(http.StatusNoContent)
}

func parseVaultQueryParams(c echo.Context) (api.GetVaultsParams, *api.APIError) {
	params := api.GetVaultsParams{}

	if pageSize := c.QueryParam("pageSize"); pageSize != "" {
		size, err := strconv.Atoi(pageSize)
		if err != nil {
			return params, api.NewAPIError(http.StatusBadRequest, "invalid pageSize")
		}
		params.PageSize = &size
	}

	if pageIndex := c.QueryParam("pageIndex"); pageIndex != "" {
		index, err := strconv.Atoi(pageIndex)
		if err != nil {
			return params, api.NewAPIError(http.StatusBadRequest, "invalid pageIndex")
		}
		params.PageIndex = &index
	}

	return params, nil
}

func getAPIKeysHandler(c echo.Context) error {
	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	params := api.GetAPIKeysParams{}
	if pageSize := c.QueryParam("pageSize"); pageSize != "" {
		size, err := strconv.Atoi(pageSize)
		if err != nil {
			return sendError(c, http.StatusBadRequest, "invalid pageSize")
		}
		params.PageSize = size
	}
	if pageIndex := c.QueryParam("pageIndex"); pageIndex != "" {
		index, err := strconv.Atoi(pageIndex)
		if err != nil {
			return sendError(c, http.StatusBadRequest, "invalid pageIndex")
		}
		params.PageIndex = index
	}

	resp, apiErr := api.GetAPIKeysForUser(user, params)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func createAPIKeyHandler(c echo.Context) error {
	var req api.CreateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	resp, apiErr := api.CreateAPIKeyForUser(user, req, clientInfo)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusCreated, resp)
}

func updateAPIKeyHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return sendError(c, http.StatusBadRequest, "invalid API key id")
	}

	var req api.UpdateAPIKeyRequest
	if err := c.Bind(&req); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	resp, apiErr := api.UpdateAPIKeyForUser(user, id, req, clientInfo)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func deleteAPIKeyHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return sendError(c, http.StatusBadRequest, "invalid API key id")
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())

	var user *model.User
	if ctxUser, ok := c.Get("user").(*model.User); ok {
		user = ctxUser
	}

	if apiErr := api.DeleteAPIKeyForUser(user, id, clientInfo); apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.NoContent(http.StatusNoContent)
}
