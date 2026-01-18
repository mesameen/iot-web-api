package devices

import (
	"net/http"

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

func (h *Handler) getDevices(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_registered_devices")
	defer span.End()
	records, err := h.ctrl.GetRegisteredDevices(ctx)
	if err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to get telematics data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.telem.Infof(ctx, "No of devices returned: %d", len(records))
	c.JSON(http.StatusOK, records)
}
