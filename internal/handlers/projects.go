package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"taskflow/internal/models"
	"taskflow/internal/repositories"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectRepo  repositories.ProjectRepository
	userRepo     repositories.UserRepository
	activityRepo repositories.ActivityRepository
}

func NewProjectHandler(
	projectRepo repositories.ProjectRepository,
	userRepo repositories.UserRepository,
	activityRepo repositories.ActivityRepository,
) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo, userRepo: userRepo, activityRepo: activityRepo}
}

type CreateProjectInput struct {
	Title       string `form:"title"       binding:"required,min=3,max=200"`
	Description string `form:"description" binding:"max=1000"`
}

type AddMemberInput struct {
	Email string `form:"email" binding:"required,email"`
}

func (h *ProjectHandler) Index(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	projects, err := h.projectRepo.FindByMemberID(user.ID)
	if err != nil {
		projects = []models.Project{}
	}
	c.HTML(http.StatusOK, "projects/index", gin.H{
		"Title":    "Meus Projetos",
		"User":     user,
		"Projects": projects,
	})
}

func (h *ProjectHandler) New(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	c.HTML(http.StatusOK, "projects/new", gin.H{
		"Title": "Novo Projeto",
		"User":  user,
	})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)

	var input CreateProjectInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "projects/new", gin.H{
			"Title": "Novo Projeto",
			"User":  user,
			"Error": "O título é obrigatório (mínimo 3 caracteres).",
			"Input": input,
		})
		return
	}

	project := &models.Project{
		Title:       input.Title,
		Description: input.Description,
		OwnerID:     user.ID,
	}

	if err := h.projectRepo.Create(project); err != nil {
		c.HTML(http.StatusInternalServerError, "projects/new", gin.H{
			"Title": "Novo Projeto",
			"User":  user,
			"Error": "Erro ao criar projeto. Tente novamente.",
		})
		return
	}

	_ = h.projectRepo.AddMember(project.ID, user.ID)

	c.Redirect(http.StatusFound, "/projects/"+strconv.Itoa(int(project.ID)))
}

func (h *ProjectHandler) Show(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	id := parseID(c.Param("id"))

	project, err := h.projectRepo.WithMembers(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error", gin.H{"User": user, "Error": "Projeto não encontrado."})
		return
	}

	if !h.projectRepo.IsMember(project.ID, user.ID) {
		c.HTML(http.StatusForbidden, "error", gin.H{"User": user, "Error": "Você não é membro deste projeto."})
		return
	}

	activities, _ := h.activityRepo.FindByProjectID(id, 20)

	c.HTML(http.StatusOK, "projects/show", gin.H{
		"Title":      project.Title,
		"User":       user,
		"Project":    project,
		"Todo":       filterByStatus(project.Tasks, "todo"),
		"InProgress": filterByStatus(project.Tasks, "in_progress"),
		"Done":       filterByStatus(project.Tasks, "done"),
		"Activities": activities,
	})
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	id := parseID(c.Param("id"))

	project, err := h.projectRepo.FindByID(id)
	if err != nil || project.OwnerID != user.ID {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	_ = h.projectRepo.Delete(id)
	c.Redirect(http.StatusFound, "/projects")
}

func (h *ProjectHandler) ShowMembers(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	id := parseID(c.Param("id"))

	project, err := h.projectRepo.WithMembers(id)
	if err != nil || !h.projectRepo.IsMember(project.ID, user.ID) {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	c.HTML(http.StatusOK, "projects/members", gin.H{
		"Title":   "Membros — " + project.Title,
		"User":    user,
		"Project": project,
	})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	id := parseID(c.Param("id"))

	project, err := h.projectRepo.FindByID(id)
	if err != nil || project.OwnerID != user.ID {
		c.Redirect(http.StatusFound, "/projects")
		return
	}

	var input AddMemberInput
	if err := c.ShouldBind(&input); err != nil {
		h.renderMembersWithError(c, user, id, "E-mail inválido.")
		return
	}

	newMember, err := h.userRepo.FindByEmail(input.Email)
	if err != nil {
		h.renderMembersWithError(c, user, id, "Usuário não encontrado.")
		return
	}

	_ = h.projectRepo.AddMember(project.ID, newMember.ID)
	h.logActivity(id, user.ID, fmt.Sprintf("adicionou %s ao projeto", newMember.Name))
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id")+"/members")
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	id := parseID(c.Param("id"))
	memberID := parseID(c.Param("memberID"))

	project, err := h.projectRepo.FindByID(id)
	if err != nil || project.OwnerID != user.ID || memberID == project.OwnerID {
		c.Redirect(http.StatusFound, "/projects/"+c.Param("id")+"/members")
		return
	}

	_ = h.projectRepo.RemoveMember(id, memberID)
	c.Redirect(http.StatusFound, "/projects/"+c.Param("id")+"/members")
}

func (h *ProjectHandler) renderMembersWithError(c *gin.Context, user *models.User, id uint, errMsg string) {
	project, _ := h.projectRepo.WithMembers(id)
	c.HTML(http.StatusUnprocessableEntity, "projects/members", gin.H{
		"Title":   "Membros — " + project.Title,
		"User":    user,
		"Project": project,
		"Error":   errMsg,
	})
}

func (h *ProjectHandler) logActivity(projectID, userID uint, action string) {
	_ = h.activityRepo.Create(&models.ActivityLog{
		ProjectID: projectID,
		UserID:    userID,
		Action:    action,
	})
}

func filterByStatus(tasks []models.Task, status string) []models.Task {
	var out []models.Task
	for _, t := range tasks {
		if t.Status == status {
			out = append(out, t)
		}
	}
	return out
}

func parseID(s string) uint {
	id, _ := strconv.ParseUint(s, 10, 64)
	return uint(id)
}
