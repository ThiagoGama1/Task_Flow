package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"taskflow/internal/models"
)

func TestTaskCreationLogic(t *testing.T) {
	env := newTestEnv(t)
	owner := env.seedUser("Owner", "taskowner@test.com", "hash")
	project := env.seedProject("Projeto Tarefas", owner.ID)

	task := &models.Task{
		Title:     "Implementar login",
		Status:    "todo",
		ProjectID: project.ID,
	}
	if err := env.taskRepo.Create(task); err != nil {
		t.Fatalf("erro ao criar tarefa: %v", err)
	}
	if task.ID == 0 {
		t.Error("tarefa deveria ter ID após criação")
	}
}

func TestTaskStatusUpdate(t *testing.T) {
	env := newTestEnv(t)
	owner := env.seedUser("Owner", "status@test.com", "hash")
	project := env.seedProject("Proj", owner.ID)

	task := &models.Task{Title: "Tarefa", Status: "todo", ProjectID: project.ID}
	_ = env.taskRepo.Create(task)

	task.Status = "in_progress"
	if err := env.taskRepo.Update(task); err != nil {
		t.Fatalf("erro ao atualizar status: %v", err)
	}

	updated, err := env.taskRepo.FindByID(task.ID)
	if err != nil {
		t.Fatalf("tarefa não encontrada: %v", err)
	}
	if updated.Status != "in_progress" {
		t.Errorf("esperado status 'in_progress', obteve %q", updated.Status)
	}
}

func TestTaskDelete(t *testing.T) {
	env := newTestEnv(t)
	owner := env.seedUser("Owner", "deltask@test.com", "hash")
	project := env.seedProject("Proj", owner.ID)

	task := &models.Task{Title: "Tarefa", Status: "todo", ProjectID: project.ID}
	_ = env.taskRepo.Create(task)

	if err := env.taskRepo.Delete(task.ID); err != nil {
		t.Fatalf("erro ao deletar: %v", err)
	}
	_, err := env.taskRepo.FindByID(task.ID)
	if err == nil {
		t.Error("tarefa deveria ter sido excluída")
	}
}

func TestTaskStatusValidationInRoute(t *testing.T) {
	env := newTestEnv(t)
	cookie := loginAndGetCookie(t, env, "taskroute@test.com", "senha123")

	projects, _ := env.projectRepo.FindByMemberID(1)
	if len(projects) == 0 {
		// Cria projeto
		form := url.Values{"title": {"Meu Projeto"}, "description": {""}}
		req := httptest.NewRequest(http.MethodPost, "/projects",
			strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		w := httptest.NewRecorder()
		env.router.ServeHTTP(w, req)

		projects, _ = env.projectRepo.FindByMemberID(1)
	}

	if len(projects) == 0 {
		t.Fatal("nenhum projeto disponível para o teste")
	}

	projectID := projects[0].ID

	form := url.Values{
		"title":  {"Tarefa"},
		"status": {"invalido"},
	}
	req := httptest.NewRequest(http.MethodPost,
		"/projects/"+strconv.Itoa(int(projectID))+"/tasks",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusFound {
		t.Error("tarefa com status inválido não deveria ser criada")
	}
}

func TestNewTaskRouteRequiresAuth(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/tasks/new", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("GET /tasks/new sem auth: esperado 302, obteve %d", w.Code)
	}
}

func TestTaskAllStatusTransitions(t *testing.T) {
	env := newTestEnv(t)
	task := &models.Task{Title: "T", Status: "todo", ProjectID: 1}
	_ = env.taskRepo.Create(task)

	transitions := []string{"in_progress", "done", "todo"}
	for _, s := range transitions {
		task.Status = s
		_ = env.taskRepo.Update(task)
		got, _ := env.taskRepo.FindByID(task.ID)
		if got.Status != s {
			t.Errorf("transição para %q falhou, obteve %q", s, got.Status)
		}
	}
}

