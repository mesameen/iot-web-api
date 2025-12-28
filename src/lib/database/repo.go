package database

import (
	"context"
	"errors"

	"github.com/mesameen/iot-web-api/src/config"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Repo interface {
	Close(ctx context.Context) error
}

func New(ctx context.Context, telem *telemetryservice.Service) (Repo, error) {
	if config.Config.Common.DBName == "postgres" {
		return connectPostgres(ctx, telem)
	}
	return nil, errors.New("Unknown db configured")
}
