package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Common    *commonConfig    `toml:"common"`
	Server    *ServerConfig    `toml:"server"`
	Postgres  *PostgresConfig  `toml:"postgres"`
	Log       *LogConfig       `toml:"log"`
	Telemetry *telemetryConfig `toml:"telemetry"`
}

type commonConfig struct {
	AppName       string `toml:"app_name"`
	DBName        string `toml:"db_name"`
	APIRouteGroup string `toml:"api_route_group"`
}

type ServerConfig struct {
	Address string `toml:"address"`
}

type PostgresConfig struct {
	Address             string `toml:"address"`
	TelematicsDataTable string `toml:"telematics_data_table"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

type telemetryConfig struct {
	Address string `toml:"endpoint"`
}

var Config Configuration

// InitConfig loads config
func InitConfig() error {
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		log.Panicf("Failed to load config %v\n", err)
	}
	log.Printf("Config: %+v\n", Config)
	return nil
}
