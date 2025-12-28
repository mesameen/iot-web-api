package telematics

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/gettelematicsdata", h.gettelematicsdata)
}
