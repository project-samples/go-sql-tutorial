package app

import (
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
)

type Root struct {
	Server     ServerConfig   `mapstructure:"server"`
	DB         DatabaseConfig `mapstructure:"db"`
	Log        log.Config     `mapstructure:"log"`
	MiddleWare mid.LogConfig  `mapstructure:"middleware"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver         string `mapstructure:"driver"`
	DataSourceName string `mapstructure:"data_source_name"`
}
