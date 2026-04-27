package connections

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/data", h.getConnections)
	r.POST("/recent", h.getRecentConnections)
}
