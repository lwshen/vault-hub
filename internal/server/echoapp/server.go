package echoapp

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lwshen/vault-hub/internal/embed"
)

// Options configures the baseline Echo server instance.
type Options struct {
	Logger *slog.Logger
}

// NewServer constructs a baseline Echo server with shared middleware and static
// asset serving configured. Further route registration occurs in later phases.
func NewServer(opts Options) (*echo.Echo, error) {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(requestLogger(logger))
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(SecurityMiddleware())

	if err := mountStatic(e, logger); err != nil {
		return nil, err
	}

	return e, nil
}

func mountStatic(e *echo.Echo, logger *slog.Logger) error {
	distFS, err := embed.GetDistFS()
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(distFS))
	e.GET("/*", echo.WrapHandler(fileServer))
	logger.Info("Echo static assets mounted", "path", "/*")
	return nil
}

func requestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			elapsed := time.Since(start)

			status := c.Response().Status
			attrs := []any{
				"method", c.Request().Method,
				"path", c.Path(),
				"status", status,
				"latency", elapsed.String(),
				"remote_ip", c.RealIP(),
				"user_agent", c.Request().UserAgent(),
			}

			var logErr error
			if err != nil {
				attrs = append(attrs, "error", err.Error())
				logErr = err
			}

			switch {
			case status >= 500 || errors.Is(logErr, echo.ErrInternalServerError):
				logger.Error("http_request", attrs...)
			case status >= 400:
				logger.Warn("http_request", attrs...)
			default:
				logger.Info("http_request", attrs...)
			}

			return err
		}
	}
}
