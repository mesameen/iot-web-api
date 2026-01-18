package database

import (
	"context"
	"errors"

	"github.com/mesameen/iot-web-api/src/config"
	"github.com/mesameen/iot-web-api/src/model"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Repo interface {
	GetTelematicsData(ctx context.Context) ([]*model.TelematicsData, error)
	GetConnectionSnapshotsData(ctx context.Context) ([]*model.ConnectionsData, error)
	GetRegistereddevices(ctx context.Context) ([]*model.RegisteredDevice, error)
	GetCommands(ctx context.Context) ([]*model.Command, error)
	GetEntities(ctx context.Context) ([]*model.ConnectionsData, error)
	Close(ctx context.Context) error
}

func New(ctx context.Context, telem *telemetryservice.Service) (Repo, error) {
	if config.Config.Common.DBName == "postgres" {
		return connectPostgres(ctx, telem)
	}
	return nil, errors.New("Unknown db configured")
}
