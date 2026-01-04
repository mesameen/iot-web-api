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
	// f, err := os.Open("telematics.json")
	// if err != nil {
	// 	return nil, err
	// }
	// contentBytes, err := io.ReadAll(f)
	// if err != nil {
	// 	return nil, err
	// }
	// var posts []*model.TelematicsData

	// if err := json.Unmarshal(contentBytes, &posts); err != nil {
	// 	return nil, err
	// }
	telematcsRecords, err := c.db.GetTelematicsData(ctx)
	if err != nil {
		return nil, err
	}
	return telematcsRecords, nil
}
