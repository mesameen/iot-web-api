package entities

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

func (h *Handler) getentities(c *gin.Context) {
	ctx, span := h.telem.TraceStart(c.Request.Context(), "get_entities_data")
	defer span.End()
	records, err := h.ctrl.GetEntitiesData(ctx)
	if err != nil {
		h.telem.Errorf(c.Request.Context(), "Failed to get entities data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.telem.Infof(ctx, "No of entites data returned: %d", len(records))
	c.JSON(http.StatusOK, records)
}
