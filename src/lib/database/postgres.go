package database

import (
	"context"

	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Postgres struct {
	telem *telemetryservice.Service
}

func connectPostgres(ctx context.Context, telem *telemetryservice.Service) (*Postgres, error) {
	return &Postgres{
		telem,
	}, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	return nil
}
