package controller

import (
	"context"

	"github.com/mesameen/iot-web-api/src/lib/database"
	"github.com/mesameen/iot-web-api/src/model"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Controller struct {
	telem *telemetryservice.Service
	db    database.Repo
}

func New(telem *telemetryservice.Service, db database.Repo) *Controller {
	return &Controller{
		telem,
		db,
	}
}

func (c *Controller) GetTelematicsData(ctx context.Context) ([]*model.TelematicsData, error) {
	telematcsRecords, err := c.db.GetTelematicsData(ctx)
	if err != nil {
		return nil, err
	}
	return telematcsRecords, nil
}

func (c *Controller) GetConnectionsData(ctx context.Context) ([]*model.ConnectionsData, error) {
	records, err := c.db.GetConnectionSnapshotsData(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c *Controller) GetRegisteredDevices(ctx context.Context) ([]*model.RegisteredDevice, error) {
	records, err := c.db.GetRegistereddevices(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c *Controller) GetEntitiesData(ctx context.Context) ([]*model.ConnectionsData, error) {
	records, err := c.db.GetConnectionSnapshotsData(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c *Controller) GetCommands(ctx context.Context) ([]*model.Command, error) {
	records, err := c.db.GetCommands(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}
