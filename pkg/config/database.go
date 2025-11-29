package config

import "time"

type DatabaseConfig struct {
	Host             string           `mapstructure:"host"`
	Port             string           `mapstructure:"port"`
	Username         string           `mapstructure:"username"`
	Password         string           `mapstructure:"password"`
	Database         string           `mapstructure:"database"`
	ConnectionConfig ConnectionConfig `mapstructure:"connection"`
}

type ConnectionConfig struct {
	MaxOpen     int           `mapstructure:"max-open"`
	MaxIdle     int           `mapstructure:"max-idle"`
	MaxLifetime time.Duration `mapstructure:"max-life-time"`
	MaxIdleTime time.Duration `mapstructure:"max-idle-time"`
}
