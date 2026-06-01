package routes

import (
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProjects(r *gin.Engine, h *handlers.ProjectHandler) {
	projects := r.Group("/projects", middleware.Required())
	{
		projects.GET("", h.Index)
		projects.GET("/new", h.New)
		projects.POST("", h.Create)
		projects.GET("/:id", h.Show)
		projects.POST("/:id/delete", h.Delete)
		projects.GET("/:id/members", h.ShowMembers)
		projects.POST("/:id/members", h.AddMember)
		projects.POST("/:id/members/:memberID/remove", h.RemoveMember)
	}
}
