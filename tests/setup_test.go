package tests

import (
	"net/http"
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"
	"taskflow/internal/models"
	"taskflow/internal/routes"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testEnv struct {
	router       *gin.Engine
	userRepo     *mockUserRepo
	projectRepo  *mockProjectRepo
	taskRepo     *mockTaskRepo
	commentRepo  *mockCommentRepo
	activityRepo *mockActivityRepo
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	userRepo := newMockUserRepo()
	projectRepo := newMockProjectRepo()
	taskRepo := newMockTaskRepo()
	commentRepo := newMockCommentRepo()
	activityRepo := newMockActivityRepo()

	r := gin.New()
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("taskflow_session", store))

	authMW := middleware.NewAuthMiddleware(userRepo)
	r.Use(authMW.LoadUser())

	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	authH := handlers.NewAuthHandler(userRepo)
	projectH := handlers.NewProjectHandler(projectRepo, userRepo, activityRepo)
	taskH := handlers.NewTaskHandler(taskRepo, projectRepo, userRepo, commentRepo, activityRepo)

	routes.RegisterAuth(r, authH)
	routes.RegisterProjects(r, projectH)
	routes.RegisterTasks(r, taskH)

	return &testEnv{
		router:       r,
		userRepo:     userRepo,
		projectRepo:  projectRepo,
		taskRepo:     taskRepo,
		commentRepo:  commentRepo,
		activityRepo: activityRepo,
	}
}

func (e *testEnv) seedUser(name, email, password string) *models.User {
	u := &models.User{Name: name, Email: email, Password: password}
	_ = e.userRepo.Create(u)
	return u
}

func (e *testEnv) seedProject(title string, ownerID uint) *models.Project {
	p := &models.Project{Title: title, OwnerID: ownerID}
	_ = e.projectRepo.Create(p)
	_ = e.projectRepo.AddMember(p.ID, ownerID)
	return p
}
