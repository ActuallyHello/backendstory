package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ActuallyHello/backendstory/pkg/config"
	"github.com/ActuallyHello/backendstory/pkg/container"
	"github.com/ActuallyHello/backendstory/pkg/server"
	// "github.com/ActuallyHello/backendstory/internal/config"
	// "github.com/ActuallyHello/backendstory/internal/core/container"
	// "github.com/ActuallyHello/backendstory/internal/server/router"
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelError = "error"

	deploymentLocal = "local"
	deploymentDev   = "dev"
	deploymentProd  = "prod"
)

// @title BackendStory API
// @version 1.0
// @description API для управления перечислениями, пользователями и аутентификацией
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @x-extension-openapi {"example": "value"}
func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	config := config.MustLoadConfig(".")

	appLogger := MustSetupLogger(config.Deployment, config.LogLevel)
	slog.SetDefault(appLogger)

	container, err := container.NewAppContainer(ctx, config)
	if err != nil {
		slog.Error("failed to init container", "error", err)
		os.Exit(1)
	}
	defer container.Close()

	slog.Info("Application ready to start!")

	// Получаем текущую рабочую директорию
	workDir, err := os.Getwd()
	if err != nil {
		slog.Error("failed to get working directory", "error", err)
		os.Exit(1)
	}
	// Создаем абсолютный путь к статическим файлам
	staticFilesPath := filepath.Join(workDir, config.ServerConfig.StaticFilesPath)
	// Создаем директорию если она не существует
	staticFullPath := filepath.Join(staticFilesPath)
	if err := os.MkdirAll(staticFullPath, 0755); err != nil {
		slog.Error("failed to create static directory", "error", err)
		os.Exit(1)
	}

	router, err := server.SetupRouter(container, staticFullPath)
	if err != nil {
		slog.Error("failed to setup routes", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    config.ServerConfig.Addr,
		Handler: router,
	}

	go func() {
		slog.Info("Starting server on port " + config.ServerConfig.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	slog.Info("Application stopped gracefully")
}

func MustSetupLogger(deployment, level string) *slog.Logger {
	var handlerOptions *slog.HandlerOptions
	switch strings.ToLower(level) {
	case levelDebug:
		handlerOptions = &slog.HandlerOptions{Level: slog.LevelDebug}
	case levelInfo:
		handlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	case levelError:
		handlerOptions = &slog.HandlerOptions{Level: slog.LevelError}
	default:
		handlerOptions = &slog.HandlerOptions{Level: slog.LevelInfo}
	}

	var logger *slog.Logger
	switch strings.ToLower(deployment) {
	case deploymentLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	case deploymentDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	case deploymentProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	}

	return logger
}
