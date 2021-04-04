package app

type Root struct {
	Server ServerConfig   `mapstructure:"server"`
	DB     DatabaseConfig `mapstructure:"db"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver         string `mapstructure:"driver"`
	DataSourceName string `mapstructure:"data_source_name"`
}
