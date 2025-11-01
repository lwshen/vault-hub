package echoapp

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/auth"
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
	apiGroup.POST("/auth/login", loginHandler)
	apiGroup.POST("/auth/signup", signupHandler)
	apiGroup.GET("/auth/logout", logoutHandler)
	apiGroup.GET("/auth/login/oidc", loginOIDCHandler)
	apiGroup.GET("/auth/callback/oidc", loginOIDCCallbackHandler)
	apiGroup.POST("/auth/password/reset/request", requestPasswordResetHandler)
	apiGroup.POST("/auth/password/reset/confirm", confirmPasswordResetHandler)
	apiGroup.POST("/auth/magic-link/request", requestMagicLinkHandler)
	apiGroup.GET("/auth/magic-link/token", consumeMagicLinkHandler)
	apiGroup.GET("/vaults", getVaultsHandler)
	apiGroup.GET("/vaults/:uniqueId", getVaultHandler)
	apiGroup.POST("/vaults", createVaultHandler)
	apiGroup.PUT("/vaults/:uniqueId", updateVaultHandler)
	apiGroup.DELETE("/vaults/:uniqueId", deleteVaultHandler)
	apiGroup.GET("/api-keys", getAPIKeysHandler)
	apiGroup.POST("/api-keys", createAPIKeyHandler)
	apiGroup.PATCH("/api-keys/:id", updateAPIKeyHandler)
	apiGroup.DELETE("/api-keys/:id", deleteAPIKeyHandler)
	apiGroup.GET("/audit-logs", getAuditLogsHandler)
	apiGroup.GET("/audit-logs/metrics", getAuditMetricsHandler)
	cliGroup := apiGroup.Group("/cli")
	cliGroup.GET("/vaults", getCLIVaultsHandler)
	cliGroup.GET("/vault/:uniqueId", getCLIVaultHandler)
	cliGroup.GET("/vault/name/:name", getCLIVaultByNameHandler)
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

func loginHandler(c echo.Context) error {
	var input api.LoginRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())
	resp, apiErr := api.LoginWithPassword(input.Email, input.Password, clientInfo)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func signupHandler(c echo.Context) error {
	var input api.SignupRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())
	resp, apiErr := api.SignupWithPassword(input, clientInfo)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func logoutHandler(c echo.Context) error {
	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())
	if user, ok := c.Get("user").(*model.User); ok {
		api.RecordLogoutAudit(user, clientInfo)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Successfully logged out",
	})
}

func loginOIDCHandler(c echo.Context) error {
	redirectURL, err := auth.AuthCodeURL(baseURL(c))
	if err != nil {
		slog.Error("Failed to get OIDC URL", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

func loginOIDCCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if err := auth.VerifyState(state); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	token, err := auth.Verify(c.Request().Context(), code)
	if err != nil {
		slog.Error("Failed to verify OIDC token", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	userInfo, err := auth.UserInfo(c.Request().Context(), token)
	if err != nil {
		slog.Error("Failed to fetch OIDC user info", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	var claims map[string]any
	if err := userInfo.Claims(&claims); err != nil {
		slog.Error("Failed to extract OIDC claims", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	emailClaim, ok := claims["email"].(string)
	if !ok || emailClaim == "" {
		slog.Error("OIDC userInfo missing email claim", slog.Any("claims", claims))
		return c.NoContent(http.StatusBadRequest)
	}

	user := model.User{Email: emailClaim}
	if err := user.GetByEmail(); err != nil {
		name := ""
		if nameClaim, ok := claims["name"].(string); ok {
			name = nameClaim
		}

		createParams := model.CreateUserParams{
			Email:    emailClaim,
			Password: nil,
			Name:     name,
		}

		newUser, createErr := createParams.Create()
		if createErr != nil {
			slog.Error("Failed to create user from OIDC", "error", createErr, "email", emailClaim)
			return c.NoContent(http.StatusInternalServerError)
		}
		user = *newUser
		slog.Info("User created from OIDC", "email", emailClaim, "name", name)
	}

	jwtToken, err := user.GenerateToken()
	if err != nil {
		slog.Error("Failed to generate token for OIDC user", "error", err, "userID", user.ID)
		return c.NoContent(http.StatusInternalServerError)
	}

	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())
	if err := model.LogUserAction(model.ActionLoginUser, user.ID, model.SourceWeb, clientInfo.IP, clientInfo.UserAgent); err != nil {
		slog.Error("Failed to create audit log for OIDC login", "error", err, "userID", user.ID)
	}

	redirectFragment := "/login#token=" + url.QueryEscape(jwtToken) + "&source=oidc"
	return c.Redirect(http.StatusFound, redirectFragment)
}

func requestPasswordResetHandler(c echo.Context) error {
	var input api.PasswordResetRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	outcome, apiErr := api.RequestPasswordResetEmail(input.Email, baseURL(c))
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	if outcome.RetryAfterHeader != nil {
		c.Response().Header().Set(echo.HeaderRetryAfter, *outcome.RetryAfterHeader)
	}

	return c.JSON(outcome.Status, emailOutcomeToResponse(outcome))
}

func confirmPasswordResetHandler(c echo.Context) error {
	var input api.PasswordResetConfirmRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	if apiErr := api.ConfirmPasswordResetToken(input.Token, input.NewPassword); apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.NoContent(http.StatusOK)
}

func requestMagicLinkHandler(c echo.Context) error {
	var input api.MagicLinkRequest
	if err := c.Bind(&input); err != nil {
		return sendError(c, http.StatusBadRequest, err.Error())
	}

	outcome, apiErr := api.RequestMagicLinkEmail(input.Email, baseURL(c))
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	if outcome.RetryAfterHeader != nil {
		c.Response().Header().Set(echo.HeaderRetryAfter, *outcome.RetryAfterHeader)
	}

	return c.JSON(outcome.Status, emailOutcomeToResponse(outcome))
}

func consumeMagicLinkHandler(c echo.Context) error {
	token := c.QueryParam("token")
	acceptsJSON := strings.Contains(c.Request().Header.Get(echo.HeaderAccept), echo.MIMEApplicationJSON)

	jwtToken, apiErr := api.ConsumeMagicLinkToken(token)
	if apiErr != nil {
		if acceptsJSON {
			return c.JSON(apiErr.Status, map[string]any{
				"error": apiErr.Message,
				"code":  "email_token_failed",
			})
		}
		return c.NoContent(apiErr.Status)
	}

	redirectFragment := "/login#token=" + url.QueryEscape(jwtToken) + "&source=magic"
	if acceptsJSON {
		return c.JSON(http.StatusOK, map[string]any{
			"token":       jwtToken,
			"redirectUrl": fmt.Sprintf("%s/dashboard", baseURL(c)),
			"code":        "email_token_sent",
			"success":     true,
		})
	}

	return c.Redirect(http.StatusFound, redirectFragment)
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

func getAuditLogsHandler(c echo.Context) error {
	params, apiErr := parseAuditLogQueryParams(c)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	user, ok := c.Get("user").(*model.User)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "user not found in context")
	}

	resp, apiErr := api.GetAuditLogsForUser(user, params)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func getAuditMetricsHandler(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "user not found in context")
	}

	resp, apiErr := api.GetAuditMetricsForUser(user)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func getCLIVaultsHandler(c echo.Context) error {
	apiKey, ok := c.Get("api_key").(*model.APIKey)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "API key not found in context")
	}

	vaults, apiErr := api.GetVaultsForAPIKey(apiKey)
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, vaults)
}

func getCLIVaultHandler(c echo.Context) error {
	apiKey, ok := c.Get("api_key").(*model.APIKey)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "API key not found in context")
	}

	headerVal := c.Request().Header.Get("X-Enable-Client-Encryption")
	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())
	resp, apiErr := api.GetVaultByAPIKeyWithLookup(apiKey, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByUniqueID(c.Param("uniqueId"), apiKey.UserID)
		return &vault, err
	}, c.Param("uniqueId"), strings.EqualFold(headerVal, "true"), clientInfo, c.Request().Header.Get(echo.HeaderAuthorization), headerVal != "")
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func getCLIVaultByNameHandler(c echo.Context) error {
	apiKey, ok := c.Get("api_key").(*model.APIKey)
	if !ok {
		return sendError(c, http.StatusUnauthorized, "API key not found in context")
	}

	headerVal := c.Request().Header.Get("X-Enable-Client-Encryption")
	clientInfo := api.ExtractClientInfo(c.Request().Header.Get, c.RealIP())
	resp, apiErr := api.GetVaultByAPIKeyWithLookup(apiKey, func(apiKey *model.APIKey) (*model.Vault, error) {
		var vault model.Vault
		err := vault.GetByName(c.Param("name"), apiKey.UserID)
		return &vault, err
	}, c.Param("name"), strings.EqualFold(headerVal, "true"), clientInfo, c.Request().Header.Get(echo.HeaderAuthorization), headerVal != "")
	if apiErr != nil {
		return sendError(c, apiErr.Status, apiErr.Message)
	}

	return c.JSON(http.StatusOK, resp)
}

func emailOutcomeToResponse(outcome api.EmailTokenOutcome) map[string]any {
	return map[string]any{
		"success": outcome.Success,
		"code":    outcome.Code,
	}
}

func baseURL(c echo.Context) string {
	req := c.Request()
	scheme := c.Scheme()
	if forwardedProto := req.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = strings.TrimSpace(strings.Split(forwardedProto, ",")[0])
	}

	host := req.Host
	if forwardedHost := req.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = strings.TrimSpace(strings.Split(forwardedHost, ",")[0])
	}
	if host == "" && req.URL != nil {
		host = req.URL.Host
	}
	if host == "" {
		host = "localhost"
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

func parseAuditLogQueryParams(c echo.Context) (api.GetAuditLogsParams, *api.APIError) {
	params := api.GetAuditLogsParams{}

	if start := c.QueryParam("startDate"); start != "" {
		timeVal, err := parseISOTime(start)
		if err != nil {
			return params, api.NewAPIError(http.StatusBadRequest, "invalid startDate")
		}
		params.StartDate = timeVal
	}

	if end := c.QueryParam("endDate"); end != "" {
		timeVal, err := parseISOTime(end)
		if err != nil {
			return params, api.NewAPIError(http.StatusBadRequest, "invalid endDate")
		}
		params.EndDate = timeVal
	}

	if vaultID := c.QueryParam("vaultUniqueId"); vaultID != "" {
		params.VaultUniqueId = &vaultID
	}

	pageSizeStr := c.QueryParam("pageSize")
	if pageSizeStr == "" {
		return params, api.NewAPIError(http.StatusBadRequest, "pageSize is required")
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return params, api.NewAPIError(http.StatusBadRequest, "invalid pageSize")
	}
	params.PageSize = pageSize

	pageIndexStr := c.QueryParam("pageIndex")
	if pageIndexStr == "" {
		return params, api.NewAPIError(http.StatusBadRequest, "pageIndex is required")
	}
	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		return params, api.NewAPIError(http.StatusBadRequest, "invalid pageIndex")
	}
	params.PageIndex = pageIndex

	return params, nil
}

func parseISOTime(value string) (*time.Time, error) {
	layouts := []string{time.RFC3339Nano, time.RFC3339}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return &parsed, nil
		}
	}
	return nil, fmt.Errorf("invalid time format")
}
