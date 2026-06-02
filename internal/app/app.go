package app

import (
	"html/template"
	"log/slog"
	"net/http"
	"taskflow/internal/config"
	"taskflow/internal/database"
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"
	"taskflow/internal/repositories"
	"taskflow/internal/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewApp(cfg *config.Config) (*App, error) {
	gin.SetMode(cfg.GinMode)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	slog.Info("banco de dados conectado")

	if err := database.Migrate(db); err != nil {
		return nil, err
	}
	slog.Info("migrações aplicadas")

	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	authHandler := handlers.NewAuthHandler(userRepo)
	projectHandler := handlers.NewProjectHandler(projectRepo, userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo, projectRepo, userRepo)

	r := gin.Default()

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	r.Use(sessions.Sessions("taskflow_session", store))

	authMW := middleware.NewAuthMiddleware(userRepo)
	r.Use(authMW.LoadUser())

	r.HTMLRender = buildRenderer()
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		if _, ok := c.Get("current_user"); ok {
			c.Redirect(http.StatusFound, "/projects")
		} else {
			c.Redirect(http.StatusFound, "/auth/login")
		}
	})

	routes.RegisterAuth(r, authHandler)
	routes.RegisterProjects(r, projectHandler)
	routes.RegisterTasks(r, taskHandler)

	return &App{Router: r, DB: db}, nil
}

var funcMap = template.FuncMap{
	"deref": func(p *uint) uint {
		if p == nil {
			return 0
		}
		return *p
	},
}

type htmlRenderer map[string]*template.Template

func (r htmlRenderer) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: r[name].Lookup("base"),
		Data:     data,
	}
}

func buildRenderer() render.HTMLRender {
	r := make(htmlRenderer)
	base := "templates/layout/base.html"

	add := func(name, path string) {
		r[name] = template.Must(template.New(name).ParseFiles(base, path))
	}

	add("login", "templates/auth/login.html")
	add("register", "templates/auth/register.html")
	add("projects/index", "templates/projects/index.html")
	add("projects/new", "templates/projects/new.html")
	add("projects/show", "templates/projects/show.html")
	add("projects/members", "templates/projects/members.html")
	add("tasks/new", "templates/tasks/new.html")
	r["tasks/edit"] = template.Must(template.New("tasks/edit").Funcs(funcMap).ParseFiles(base, "templates/tasks/edit.html"))
	add("error", "templates/error.html")

	return r
}
