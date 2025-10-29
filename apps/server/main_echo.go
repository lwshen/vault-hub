//go:build echo

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/server/echoapp"
	"github.com/lwshen/vault-hub/internal/version"
	"github.com/lwshen/vault-hub/model"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting VaultHub Echo Server", "version", version.Version, "commit", version.Commit)

	if err := model.Open(logger); err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	e, err := echoapp.NewServer(echoapp.Options{
		Logger: logger,
	})
	if err != nil {
		logger.Error("Failed to initialize Echo server", "error", err)
		os.Exit(1)
	}

	echoapp.RegisterRoutes(e)

	if err := echoapp.MountStatic(e, logger); err != nil {
		logger.Error("Failed to mount static assets", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := e.Start(":" + config.AppPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Echo server exited", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Failed to gracefully shutdown Echo server", "error", err)
		os.Exit(1)
	}

	logger.Info("Echo server stopped gracefully")
}
