package controller

import (
	"github.com/mesameen/iot-web-api/src/lib/database"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Controller struct {
	telem *telemetryservice.Service
	db    *database.Repo
}

func New(telem *telemetryservice.Service, db *database.Repo) *Controller {
	return &Controller{
		telem,
		db,
	}
}
