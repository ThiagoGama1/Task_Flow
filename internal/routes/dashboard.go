package routes

import (
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDashboard(r *gin.Engine, h *handlers.DashboardHandler) {
	r.GET("/dashboard", middleware.Required(), h.Index)
}
