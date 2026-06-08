package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"taskflow/internal/models"
	"taskflow/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskRepo     repositories.TaskRepository
	projectRepo  repositories.ProjectRepository
	userRepo     repositories.UserRepository
	commentRepo  repositories.CommentRepository
	activityRepo repositories.ActivityRepository
}

func NewTaskHandler(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	userRepo repositories.UserRepository,
	commentRepo repositories.CommentRepository,
	activityRepo repositories.ActivityRepository,
) *TaskHandler {
	return &TaskHandler{
		taskRepo:     taskRepo,
		projectRepo:  projectRepo,
		userRepo:     userRepo,
		commentRepo:  commentRepo,
		activityRepo: activityRepo,
	}
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

	if err := h.taskRepo.Create(task); err == nil {
		h.logActivity(projectID, user.ID, fmt.Sprintf(`criou a tarefa "%s"`, task.Title))
	}
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

	comments, _ := h.commentRepo.FindByTaskID(taskID)

	c.HTML(http.StatusOK, "tasks/edit", gin.H{
		"Title":    "Editar Tarefa",
		"User":     user,
		"Project":  project,
		"Task":     task,
		"Comments": comments,
	})
}

func (h *TaskHandler) CreateComment(c *gin.Context) {
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

	content := c.PostForm("content")
	if len(content) < 1 {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id")+"/tasks/"+c.Param("taskID")+"/edit")
		return
	}

	comment := &models.Comment{
		Content: content,
		TaskID:  taskID,
		UserID:  user.ID,
	}
	if err := h.commentRepo.Create(comment); err == nil {
		h.logActivity(projectID, user.ID, fmt.Sprintf(`comentou na tarefa "%s"`, task.Title))
	}
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id")+"/tasks/"+c.Param("taskID")+"/edit")
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
		comments, _ := h.commentRepo.FindByTaskID(taskID)
		c.HTML(http.StatusUnprocessableEntity, "tasks/edit", gin.H{
			"Title":    "Editar Tarefa",
			"User":     user,
			"Project":  project,
			"Task":     task,
			"Comments": comments,
			"Error":    "Dados inválidos.",
		})
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Priority = input.Priority
	task.AssigneeID = parseOptionalID(input.AssigneeID)
	task.Assignee = nil
	task.DueDate = parseDate(input.DueDate)

	if err := h.taskRepo.Update(task); err != nil {
		project, _ := h.projectRepo.WithMembers(projectID)
		comments, _ := h.commentRepo.FindByTaskID(taskID)
		c.HTML(http.StatusInternalServerError, "tasks/edit", gin.H{
			"Title":    "Editar Tarefa",
			"User":     user,
			"Project":  project,
			"Task":     task,
			"Comments": comments,
			"Error":    "Erro ao salvar. Tente novamente.",
		})
		return
	}
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

	oldStatus := task.Status
	task.Status = input.Status
	if err := h.taskRepo.Update(task); err == nil && oldStatus != input.Status {
		label := map[string]string{"todo": "A Fazer", "in_progress": "Em Andamento", "done": "Concluído"}
		h.logActivity(projectID, user.ID,
			fmt.Sprintf(`moveu "%s" para %s`, task.Title, label[input.Status]))
	}
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

	title := task.Title
	if err := h.taskRepo.Delete(task.ID); err == nil {
		h.logActivity(projectID, user.ID, fmt.Sprintf(`excluiu a tarefa "%s"`, title))
	}
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id"))
}

func (h *TaskHandler) UpdateTitle(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projectID := parseID(c.Param("id"))
	taskID := parseID(c.Param("taskID"))

	if !h.projectRepo.IsMember(projectID, user.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil || task.ProjectID != projectID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var body struct {
		Title string `json:"title" binding:"required,min=2,max=200"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "título inválido"})
		return
	}

	task.Title = body.Title
	task.Assignee = nil
	if err := h.taskRepo.Update(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"title": task.Title})
}

func (h *TaskHandler) logActivity(projectID, userID uint, action string) {
	_ = h.activityRepo.Create(&models.ActivityLog{
		ProjectID: projectID,
		UserID:    userID,
		Action:    action,
	})
}

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
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
