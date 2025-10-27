package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginResponse represents the response for successful login
type EchoLoginResponse struct {
	Token string `json:"token"`
}

// LoginOidcEcho handles OIDC login using Echo context
// LoginOidcEcho handles OIDC login using Echo context
func LoginOidcEcho(c echo.Context) error {
	// TODO: Implement OIDC login with Echo context
	// For now, return a placeholder
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "OIDC login not yet implemented in Echo migration",
	})
}

// LoginOidcCallbackEcho handles OIDC callback using Echo context
func LoginOidcCallbackEcho(c echo.Context) error {
	// TODO: Implement OIDC callback with Echo context
	// For now, return a placeholder
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "OIDC callback not yet implemented in Echo migration",
	})
}

// getEchoBaseURL extracts the base URL from Echo context
func getEchoBaseURL(c echo.Context) string {
	req := c.Request()
	scheme := "https"
	if req.TLS == nil {
		scheme = "http"
	}

	return scheme + "://" + req.Host
}

// getEchoClientInfo extracts IP address and User-Agent from Echo request
func getEchoClientInfo(c echo.Context) (string, string) {
	req := c.Request()

	// Get IP address (check for forwarded headers first)
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.RealIP()
	}

	// Get User-Agent
	userAgent := req.UserAgent()
	return ip, userAgent
}