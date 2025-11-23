package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActuallyHello/backendstory/internal/config"
	"github.com/ActuallyHello/backendstory/internal/core/container"
	"github.com/ActuallyHello/backendstory/internal/server/router"
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
	slog.Info("Loading backendstory application...")

	config := config.MustLoadConfig(".")

	slog.Info("Config was loaded...")

	// appLogger := MustSetupLogger(config.Deployment, config.LogLevel)
	// slog.SetDefault(appLogger)

	container := container.NewAppContainer(config)

	slog.Info("Dependency container uploaded! Application ready to start!")

	// Получаем текущую рабочую директорию
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get working directory: %v", err))
	}
	// Создаем абсолютный путь к статическим файлам
	staticFilesPath := filepath.Join(workDir, config.ServerConfig.StaticFilesPath)
	// Создаем директорию если она не существует
	fullPath := filepath.Join(staticFilesPath, "imgs", "products")
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create static directory: %v", err))
	}
	slog.Info("Static directories created", "path", fullPath)

	r := router.SetupRouter(container, config.ServerConfig.StaticFilesPath)

	slog.Info("Starting server on port " + config.ServerConfig.Addr)

	if err := http.ListenAndServe(config.ServerConfig.Addr, r); err != nil {
		log.Fatal(err)
	}
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
