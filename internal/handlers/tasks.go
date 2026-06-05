package handlers

import (
	"net/http"
	"strconv"
	"taskflow/internal/models"
	"taskflow/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

func NewTaskHandler(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	userRepo repositories.UserRepository,
) *TaskHandler {
	return &TaskHandler{taskRepo: taskRepo, projectRepo: projectRepo, userRepo: userRepo}
}

type TaskInput struct {
	Title       string `form:"title"       binding:"required,min=2,max=200"`
	Description string `form:"description"`
	AssigneeID  string `form:"assignee_id"`
	Status      string `form:"status"      binding:"required,oneof=todo in_progress done"`
	Priority    string `form:"priority"    binding:"required,oneof=low medium high"`
	DueDate     string `form:"due_date"`
}

func (h *TaskHandler) New(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	id := parseID(c.Param("id"))

	project, err := h.projectRepo.WithMembers(id)
	if err != nil || !h.projectRepo.IsMember(id, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	c.HTML(http.StatusOK, "tasks/new", gin.H{
		"Title":   "Nova Tarefa",
		"User":    user,
		"Project": project,
	})
}

func (h *TaskHandler) Create(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projectID := parseID(c.Param("id"))

	if !h.projectRepo.IsMember(projectID, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	var input TaskInput
	if err := c.ShouldBind(&input); err != nil {
		project, _ := h.projectRepo.WithMembers(projectID)
		c.HTML(http.StatusUnprocessableEntity, "tasks/new", gin.H{
			"Title":   "Nova Tarefa",
			"User":    user,
			"Project": project,
			"Error":   "Título obrigatório (mínimo 2 caracteres) e status inválido.",
			"Input":   input,
		})
		return
	}

	task := &models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		ProjectID:   projectID,
		AssigneeID:  parseOptionalID(input.AssigneeID),
		DueDate:     parseDate(input.DueDate),
	}

	_ = h.taskRepo.Create(task)
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
}

func (h *TaskHandler) Edit(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projectID := parseID(c.Param("id"))
	taskID := parseID(c.Param("taskID"))

	project, err := h.projectRepo.WithMembers(projectID)
	if err != nil || !h.projectRepo.IsMember(projectID, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil || task.ProjectID != projectID {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
		return
	}

	c.HTML(http.StatusOK, "tasks/edit", gin.H{
		"Title":   "Editar Tarefa",
		"User":    user,
		"Project": project,
		"Task":    task,
	})
}

func (h *TaskHandler) Update(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projectID := parseID(c.Param("id"))
	taskID := parseID(c.Param("taskID"))

	if !h.projectRepo.IsMember(projectID, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil || task.ProjectID != projectID {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
		return
	}

	var input TaskInput
	if err := c.ShouldBind(&input); err != nil {
		project, _ := h.projectRepo.WithMembers(projectID)
		c.HTML(http.StatusUnprocessableEntity, "tasks/edit", gin.H{
			"Title":   "Editar Tarefa",
			"User":    user,
			"Project": project,
			"Task":    task,
			"Error":   "Dados inválidos.",
		})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Priority = input.Priority
	task.AssigneeID = parseOptionalID(input.AssigneeID)
	task.Assignee = nil // evita que GORM sobrescreva AssigneeID com a associação antiga
	task.DueDate = parseDate(input.DueDate)

	_ = h.taskRepo.Update(task)
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projectID := parseID(c.Param("id"))
	taskID := parseID(c.Param("taskID"))

	if !h.projectRepo.IsMember(projectID, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil || task.ProjectID != projectID {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
		return
	}

	var input struct {
		Status string `form:"status" binding:"required,oneof=todo in_progress done"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
		return
	}

	task.Status = input.Status
	_ = h.taskRepo.Update(task)
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
}

func (h *TaskHandler) Delete(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projectID := parseID(c.Param("id"))
	taskID := parseID(c.Param("taskID"))

	if !h.projectRepo.IsMember(projectID, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil || task.ProjectID != projectID {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
		return
	}

	_ = h.taskRepo.Delete(task.ID)
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
}

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	// Armazena ao meio-dia UTC para evitar deslocamentos de fuso
	noon := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, time.UTC)
	return &noon
}

func parseOptionalID(s string) *uint {
	if s == "" {
		return nil
	}
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	uid := uint(id)
	return &uid
}
