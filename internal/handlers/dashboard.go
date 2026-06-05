package handlers

import (
	"net/http"
	"taskflow/internal/models"
	"taskflow/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	taskRepo repositories.TaskRepository
}

func NewDashboardHandler(taskRepo repositories.TaskRepository) *DashboardHandler {
	return &DashboardHandler{taskRepo: taskRepo}
}

func (h *DashboardHandler) Index(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)

	tasks, _ := h.taskRepo.FindAssignedTo(user.ID)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tomorrow := today.Add(24 * time.Hour)

	var overdue, dueToday, upcoming, noDue []models.Task
	for _, t := range tasks {
		if t.Status == "done" {
			continue
		}
		if t.DueDate == nil {
			noDue = append(noDue, t)
			continue
		}
		if t.DueDate.Before(today) {
			overdue = append(overdue, t)
		} else if t.DueDate.Before(tomorrow) {
			dueToday = append(dueToday, t)
		} else {
			upcoming = append(upcoming, t)
		}
	}

	c.HTML(http.StatusOK, "dashboard", gin.H{
		"Title":    "Dashboard",
		"User":     user,
		"Overdue":  overdue,
		"DueToday": dueToday,
		"Upcoming": upcoming,
		"NoDue":    noDue,
	})
}
