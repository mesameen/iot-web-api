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

func (c *Controller) GetTelematicsData(ctx context.Context, req *model.GetTelematicsDataRequest) ([]*model.TelematicsData, error) {
	telematcsRecords, err := c.db.GetTelematicsData(ctx, req)
	if err != nil {
		return nil, err
	}
	return telematcsRecords, nil
}

func (c *Controller) GetRecentTelematicsData(ctx context.Context, req *model.GetTelematicsDataRequest) ([]*model.TelematicsData, error) {
	telematcsRecords, err := c.db.GetRecentTelematicsData(ctx, req)
	if err != nil {
		return nil, err
	}
	return telematcsRecords, nil
}

func (c *Controller) GetConnectionsData(ctx context.Context, req *model.GetConnectionsDataRequest) ([]*model.ConnectionsData, error) {
	records, err := c.db.GetConnectionEvents(ctx, req)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c *Controller) GetRecentConnectionsData(ctx context.Context, req *model.GetConnectionsDataRequest) ([]*model.ConnectionsData, error) {
	records, err := c.db.GetRecentConnectionEvents(ctx, req)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c *Controller) GetRegisteredDevices(ctx context.Context, req *model.GetRegisteredDevicesRequest) ([]*model.RegisteredDevice, error) {
	records, err := c.db.GetRegisteredDevices(ctx, req)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c *Controller) GetEntitiesData(ctx context.Context) ([]*model.ConnectionsData, error) {
	return nil, nil
}

func (c *Controller) GetCommands(ctx context.Context) ([]*model.Command, error) {
	records, err := c.db.GetCommands(ctx)
	if err != nil {
		return nil, err
	}
	return records, nil
}
