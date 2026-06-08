package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"taskflow/internal/models"
)

func main() {
	_ = godotenv.Load()

	dsn := getEnv("DATABASE_URL", "postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("falha ao conectar ao banco: %v", err)
	}

	log.Println("→ AutoMigrate...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
		&models.Comment{},
		&models.ActivityLog{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// ── Usuários ─────────────────────────────────────────────────────────────
	u1 := seedUser(db, "Demo User", "demo@taskflow.app", "demo123")
	u2 := seedUser(db, "Ana Lima", "ana@taskflow.app", "demo123")
	u3 := seedUser(db, "Carlos Souza", "carlos@taskflow.app", "demo123")

	// ── Projetos ─────────────────────────────────────────────────────────────
	p1 := seedProject(db, "Site Institucional",
		"Redesign completo do site da empresa com nova identidade visual.",
		u1, []models.User{*u1, *u2})

	p2 := seedProject(db, "App Mobile TaskFlow",
		"Desenvolvimento do aplicativo para Android e iOS.",
		u2, []models.User{*u1, *u2, *u3})

	// ── Datas ─────────────────────────────────────────────────────────────────
	now := time.Now()
	pt := func(t time.Time) *time.Time {
		d := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, time.UTC)
		return &d
	}
	overdue1 := pt(now.AddDate(0, 0, -5))  // 5 dias atrás
	overdue2 := pt(now.AddDate(0, 0, -1))  // ontem
	today    := pt(now)
	week1    := pt(now.AddDate(0, 0, 4))
	week2    := pt(now.AddDate(0, 0, 10))
	month1   := pt(now.AddDate(0, 1, 0))

	// ── Tarefas do Projeto 1 ─────────────────────────────────────────────────
	t1a := seedTask(db, p1.ID, taskInput{
		Title:      "Criar wireframes das páginas",
		Desc:       "Elaborar wireframes de todas as páginas principais no Figma.",
		Status:     "done",
		Priority:   "high",
		DueDate:    pt(now.AddDate(0, -1, 0)),
		AssigneeID: &u2.ID,
	})
	t1b := seedTask(db, p1.ID, taskInput{
		Title:      "Implementar página de contato",
		Desc:       "Formulário com validação e envio de e-mail via SMTP.",
		Status:     "todo",
		Priority:   "high",
		DueDate:    overdue1,
		AssigneeID: &u1.ID,
	})
	t1c := seedTask(db, p1.ID, taskInput{
		Title:      "Otimizar imagens do banner",
		Desc:       "Converter para WebP e ajustar tamanhos responsivos.",
		Status:     "todo",
		Priority:   "medium",
		DueDate:    today,
		AssigneeID: &u1.ID,
	})
	seedTask(db, p1.ID, taskInput{
		Title:    "Integrar Google Analytics",
		Desc:     "Adicionar tag GTM e configurar eventos de conversão.",
		Status:   "in_progress",
		Priority: "medium",
		DueDate:  week1,
	})
	seedTask(db, p1.ID, taskInput{
		Title:    "Revisar textos de SEO",
		Desc:     "Meta descriptions e títulos de todas as páginas.",
		Status:   "todo",
		Priority: "low",
	})
	seedTask(db, p1.ID, taskInput{
		Title:    "Configurar CDN para assets",
		Desc:     "Cloudfront + S3 para servir imagens e CSS.",
		Status:   "todo",
		Priority: "low",
		DueDate:  month1,
	})

	// ── Tarefas do Projeto 2 ─────────────────────────────────────────────────
	seedTask(db, p2.ID, taskInput{
		Title:      "Definir arquitetura do app",
		Desc:       "Escolha de framework React Native vs Flutter e padrão de estado.",
		Status:     "done",
		Priority:   "high",
		DueDate:    pt(now.AddDate(0, -1, -15)),
		AssigneeID: &u3.ID,
	})
	t2b := seedTask(db, p2.ID, taskInput{
		Title:      "Tela de login e cadastro",
		Desc:       "Fluxo completo com validação, feedback visual e biometria.",
		Status:     "in_progress",
		Priority:   "high",
		DueDate:    overdue2,
		AssigneeID: &u1.ID,
	})
	t2c := seedTask(db, p2.ID, taskInput{
		Title:      "Notificações push",
		Desc:       "Implementar FCM para Android e APNS para iOS.",
		Status:     "todo",
		Priority:   "medium",
		DueDate:    today,
		AssigneeID: &u2.ID,
	})
	seedTask(db, p2.ID, taskInput{
		Title:      "Testes de usabilidade",
		Desc:       "Conduzir sessões com 5 usuários reais e documentar achados.",
		Status:     "todo",
		Priority:   "medium",
		DueDate:    week2,
		AssigneeID: &u3.ID,
	})
	seedTask(db, p2.ID, taskInput{
		Title:    "Documentação da API",
		Desc:     "Gerar Swagger e exemplos de integração para parceiros.",
		Status:   "todo",
		Priority: "low",
	})
	seedTask(db, p2.ID, taskInput{
		Title:    "Pipeline CI/CD",
		Desc:     "GitHub Actions para build, teste e deploy nas stores.",
		Status:   "todo",
		Priority: "high",
		DueDate:  week2,
	})

	// ── Comentários ───────────────────────────────────────────────────────────
	seedComment(db, t1a.ID, u1.ID, "Wireframes aprovados pelo cliente! Ficaram ótimos.")
	seedComment(db, t1a.ID, u2.ID, "Vou iniciar o desenvolvimento com base nessa versão.")
	seedComment(db, t1b.ID, u1.ID, "Precisamos de validação server-side além do client-side.")
	seedComment(db, t1b.ID, u2.ID, "Posso ajudar nessa parte, já tenho a lib configurada.")
	seedComment(db, t1c.ID, u1.ID, "WebP reduz ~60% do tamanho. Vale muito a pena.")
	seedComment(db, t2b.ID, u3.ID, "Prototipei no Figma, link enviado no Slack.")
	seedComment(db, t2b.ID, u1.ID, "Aprovado! Pode começar a implementar.")
	seedComment(db, t2b.ID, u2.ID, "Lembrar de testar no iOS 15 também, teve quebras antes.")
	seedComment(db, t2c.ID, u2.ID, "FCM configurado no Firebase Console. Credenciais no Vault.")

	// ── Histórico de atividades ────────────────────────────────────────────────
	seedActivity(db, p1.ID, u2.ID, `criou a tarefa "Criar wireframes das páginas"`)
	seedActivity(db, p1.ID, u2.ID, `moveu "Criar wireframes das páginas" para Concluído`)
	seedActivity(db, p1.ID, u1.ID, `criou a tarefa "Implementar página de contato"`)
	seedActivity(db, p1.ID, u1.ID, `comentou na tarefa "Implementar página de contato"`)
	seedActivity(db, p1.ID, u1.ID, `criou a tarefa "Otimizar imagens do banner"`)
	seedActivity(db, p1.ID, u1.ID, `adicionou membro ana@taskflow.app ao projeto`)

	seedActivity(db, p2.ID, u3.ID, `criou a tarefa "Definir arquitetura do app"`)
	seedActivity(db, p2.ID, u3.ID, `moveu "Definir arquitetura do app" para Concluído`)
	seedActivity(db, p2.ID, u1.ID, `criou a tarefa "Tela de login e cadastro"`)
	seedActivity(db, p2.ID, u2.ID, `moveu "Tela de login e cadastro" para Em Andamento`)
	seedActivity(db, p2.ID, u3.ID, `comentou na tarefa "Tela de login e cadastro"`)
	seedActivity(db, p2.ID, u2.ID, `adicionou membro carlos@taskflow.app ao projeto`)

	// ── Resumo ────────────────────────────────────────────────────────────────
	fmt.Println("\n============================================================")
	fmt.Println("  SEED CONCLUÍDO — TaskFlow")
	fmt.Println("============================================================")
	fmt.Println("\n  CREDENCIAIS DE ACESSO:")
	fmt.Println("  ┌─────────────────────────────────────┬──────────┐")
	fmt.Println("  │ E-mail                              │ Senha    │")
	fmt.Println("  ├─────────────────────────────────────┼──────────┤")
	fmt.Println("  │ demo@taskflow.app   (Demo User)     │ demo123  │")
	fmt.Println("  │ ana@taskflow.app    (Ana Lima)      │ demo123  │")
	fmt.Println("  │ carlos@taskflow.app (Carlos Souza)  │ demo123  │")
	fmt.Println("  └─────────────────────────────────────┴──────────┘")
	fmt.Printf("\n  Projetos criados: \"%s\" (ID=%d), \"%s\" (ID=%d)\n",
		p1.Title, p1.ID, p2.Title, p2.ID)
	fmt.Println("\n  Tarefas:")
	fmt.Println("  • 2 ATRASADAS  (vermelhas no dashboard + banner de alerta)")
	fmt.Println("  • 2 VENCEM HOJE (amarelas no dashboard)")
	fmt.Println("  • Várias com prazo futuro, sem prazo e já concluídas")
	fmt.Println("  • Comentários em 5 tarefas (badges nos cards do Kanban)")
	fmt.Println("  • 12 entradas no histórico de atividades")
	fmt.Println("\n  Acesse: http://localhost:3000/auth/login")
	fmt.Println("============================================================\n")
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func seedUser(db *gorm.DB, name, email, password string) *models.User {
	var u models.User
	if err := db.Where("email = ?", email).First(&u).Error; err == nil {
		log.Printf("  [skip] usuário já existe: %s", email)
		return &u
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u = models.User{Name: name, Email: email, Password: string(hashed)}
	if err := db.Create(&u).Error; err != nil {
		log.Fatalf("criar usuário %s: %v", email, err)
	}
	log.Printf("  [+] usuário: %s", email)
	return &u
}

func seedProject(db *gorm.DB, title, desc string, owner *models.User, members []models.User) *models.Project {
	var p models.Project
	if err := db.Where("title = ? AND owner_id = ?", title, owner.ID).First(&p).Error; err == nil {
		log.Printf("  [skip] projeto já existe: %s", title)
		return &p
	}
	p = models.Project{Title: title, Description: desc, OwnerID: owner.ID, Members: members}
	if err := db.Create(&p).Error; err != nil {
		log.Fatalf("criar projeto %s: %v", title, err)
	}
	log.Printf("  [+] projeto: %s (ID=%d)", title, p.ID)
	return &p
}

type taskInput struct {
	Title      string
	Desc       string
	Status     string
	Priority   string
	DueDate    *time.Time
	AssigneeID *uint
}

func seedTask(db *gorm.DB, projectID uint, in taskInput) *models.Task {
	var t models.Task
	if err := db.Where("title = ? AND project_id = ?", in.Title, projectID).First(&t).Error; err == nil {
		return &t
	}
	t = models.Task{
		Title:       in.Title,
		Description: in.Desc,
		Status:      in.Status,
		Priority:    in.Priority,
		DueDate:     in.DueDate,
		ProjectID:   projectID,
		AssigneeID:  in.AssigneeID,
	}
	if err := db.Create(&t).Error; err != nil {
		log.Fatalf("criar tarefa %s: %v", in.Title, err)
	}
	return &t
}

func seedComment(db *gorm.DB, taskID, userID uint, content string) {
	var n int64
	db.Model(&models.Comment{}).
		Where("task_id = ? AND user_id = ? AND content = ?", taskID, userID, content).
		Count(&n)
	if n > 0 {
		return
	}
	db.Create(&models.Comment{TaskID: taskID, UserID: userID, Content: content})
}

func seedActivity(db *gorm.DB, projectID, userID uint, action string) {
	var n int64
	db.Model(&models.ActivityLog{}).
		Where("project_id = ? AND user_id = ? AND action = ?", projectID, userID, action).
		Count(&n)
	if n > 0 {
		return
	}
	db.Create(&models.ActivityLog{ProjectID: projectID, UserID: userID, Action: action})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
