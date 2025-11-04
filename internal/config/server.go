package config

import "time"

type ServerConfig struct {
	Addr          string        `mapstructure:"addr"`
	TimeoutConfig TimeoutConfig `mapstructure:"timeout"`
}

type TimeoutConfig struct {
	Idle  time.Duration `mapstructure:"idle"`
	Read  time.Duration `mapstructure:"read"`
	Write time.Duration `mapstructure:"write"`
}
