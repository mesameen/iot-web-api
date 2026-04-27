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
	Address                     string `toml:"address"`
	TelematicsDataTable         string `toml:"telematicsdata_table"`
	RecentTelematicsDataTable   string `toml:"recent_telematicsdata_table"`
	ConnectionEventsTable       string `toml:"connection_events_table"`
	RecentConnectionEventsTable string `toml:"recent_connection_events_table"`
	RegisteredDevicesTable      string `toml:"registered_devices_table"`
	CommandsTable               string `toml:"commands_table"`
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
	log.Printf("ServerConfig: %+v\n", Config.Server)
	log.Printf("Common Config: %+v\n", Config.Common)
	log.Printf("Postgres Config: %+v\n", Config.Postgres)
	return nil
}
