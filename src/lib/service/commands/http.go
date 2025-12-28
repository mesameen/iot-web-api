package commands

import (
	"github.com/gin-gonic/gin"
	"github.com/mesameen/iot-web-api/src/lib/controller"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

type Handler struct {
	telem telemetryservice.Repo
	ctrl  *controller.Controller
}

func NewHandler(telem telemetryservice.Repo, ctrl *controller.Controller) *Handler {
	return &Handler{
		telem,
		ctrl,
	}
}

func (h *Handler) getCommands(c *gin.Context) {

}
