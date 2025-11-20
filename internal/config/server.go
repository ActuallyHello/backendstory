package config

import "time"

type ServerConfig struct {
	Addr            string        `mapstructure:"addr"`
	TimeoutConfig   TimeoutConfig `mapstructure:"timeout"`
	StaticFilesPath string        `mapstructure:"static"`
}

type TimeoutConfig struct {
	Idle  time.Duration `mapstructure:"idle"`
	Read  time.Duration `mapstructure:"read"`
	Write time.Duration `mapstructure:"write"`
}
