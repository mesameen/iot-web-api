package telematics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/iot-web-api/src/lib/controller"
	"github.com/mesameen/iot-web-api/src/model"
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

func (h *Handler) getTelematicsData(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_telematics_data")
	defer span.End()
	var req model.GetTelematicsDataRequest
	if err := c.BindJSON(&req); err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to decode request data. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	records, err := h.ctrl.GetTelematicsData(ctx, &req)
	if err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to get telematics data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.telem.Infof(ctx, "No of telematics data returned: %d", len(records))
	c.JSON(http.StatusOK, records)
}

func (h *Handler) getRecentTelematicsData(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_recent_telematics_data")
	defer span.End()
	var req model.GetTelematicsDataRequest
	if err := c.BindJSON(&req); err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to decode request data. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	records, err := h.ctrl.GetRecentTelematicsData(ctx, &req)
	if err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to get recent telematics data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.telem.Infof(ctx, "No of telematics data returned: %d", len(records))
	c.JSON(http.StatusOK, records)
}
