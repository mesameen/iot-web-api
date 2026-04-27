package telematics

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/data", h.getTelematicsData)
	r.POST("/recent", h.getRecentTelematicsData)
}
