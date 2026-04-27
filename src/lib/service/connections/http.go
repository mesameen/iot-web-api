package connections

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

func (h *Handler) getConnections(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_connections_data")
	defer span.End()
	records, err := h.ctrl.GetConnectionsData(ctx)
	if err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to get telematics data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.telem.Infof(ctx, "No of connections data returned: %d", len(records))
	c.JSON(http.StatusOK, records)
}

func (h *Handler) getRecentConnections(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_recent_connections_data")
	defer span.End()
	records, err := h.ctrl.GetRecentConnectionsData(ctx)
	if err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to get recent telematics data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.telem.Infof(ctx, "No of recent connections data returned: %d", len(records))
	c.JSON(http.StatusOK, records)
}
