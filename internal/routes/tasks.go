package routes

import (
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTasks(r *gin.Engine, h *handlers.TaskHandler) {
	tasks := r.Group("/projects/:id/tasks", middleware.Required())
	{
		tasks.GET("/new", h.New)
		tasks.POST("", h.Create)
		tasks.GET("/:taskID/edit", h.Edit)
		tasks.POST("/:taskID", h.Update)
		tasks.POST("/:taskID/status", h.UpdateStatus)
		tasks.POST("/:taskID/delete", h.Delete)
		tasks.POST("/:taskID/comments", h.CreateComment)
		tasks.POST("/:taskID/title", h.UpdateTitle)
	}
}
