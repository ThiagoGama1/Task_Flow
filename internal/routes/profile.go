package routes

import (
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProfile(r *gin.Engine, h *handlers.ProfileHandler) {
	profile := r.Group("/profile", middleware.Required())
	{
		profile.GET("", h.Show)
		profile.POST("", h.Update)
	}
}
