package database

import (
	"context"
	"errors"

	"github.com/mesameen/iot-web-api/src/config"
	"github.com/mesameen/iot-web-api/src/model"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Repo interface {
	GetTelematicsData(ctx context.Context, req *model.GetTelematicsDataRequest) ([]*model.TelematicsData, error)
	GetRecentTelematicsData(ctx context.Context, req *model.GetTelematicsDataRequest) ([]*model.TelematicsData, error)
	GetConnectionEvents(ctx context.Context, req *model.GetConnectionsDataRequest) ([]*model.ConnectionsData, error)
	GetRecentConnectionEvents(ctx context.Context, req *model.GetConnectionsDataRequest) ([]*model.ConnectionsData, error)
	GetRegisteredDevices(ctx context.Context, req *model.GetRegisteredDevicesRequest) ([]*model.RegisteredDevice, error)
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
