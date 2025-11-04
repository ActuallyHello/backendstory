package config

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	Deployment     string          `mapstructure:"deployment"`
	LogLevel       string          `mapstructure:"log-level"`
	DatabaseConfig *DatabaseConfig `mapstructure:"database"`
	ServerConfig   *ServerConfig   `mapstructure:"server"`
	KeycloakConfig *KeycloakConfig `mapstructure:"keycloak"`
}

func MustLoadConfig(path string) *ApplicationConfig {
	var (
		commonConfigPath = "."
		configName       = "application"
		configType       = "yaml"
		configRoot       = "app"
	)

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			slog.Error("Env file for project was not found!", "err", err)
			log.Fatal(err)
		}
	} else {
		slog.Info("No .env file found, using environment variables")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(commonConfigPath)
	if path != "" {
		viper.AddConfigPath(path)
	}

	// viper.SetEnvPrefix(configRoot)
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config file: %v", err)
	}

	replaceEnv()

	var applicationConfig ApplicationConfig
	if err := viper.UnmarshalKey(configRoot, &applicationConfig); err != nil {
		log.Fatalf("error while unmarshal config file: %v", err)
	}

	return &applicationConfig
}

func replaceEnv() {
	keys := viper.AllKeys()
	for _, key := range keys {
		value := viper.GetString(key)
		if strings.Contains(value, "${") {
			expandedVal := expandEnv(value)
			viper.Set(key, expandedVal)
		}
	}
}

func expandEnv(s string) string {
	return os.Expand(s, func(value string) string {
		// Обрабатываем значения по умолчанию ${VAR:default}
		if strings.Contains(s, ":") {
			parts := strings.SplitN(value, ":", 2)
			envValue := os.Getenv(parts[0])
			if envValue == "" {
				return parts[1]
			}
			return envValue
		}
		return os.Getenv(value)
	})
}
