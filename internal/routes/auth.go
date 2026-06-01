package routes

import (
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(r *gin.Engine, h *handlers.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.GET("/login", h.ShowLogin)
		auth.POST("/login", h.Login)
		auth.GET("/register", h.ShowRegister)
		auth.POST("/register", h.Register)
		auth.POST("/logout", middleware.Required(), h.Logout)
	}
}
