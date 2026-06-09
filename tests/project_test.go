package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func loginAndGetCookie(t *testing.T, env *testEnv, email, password string) *http.Cookie {
	t.Helper()

	form := url.Values{
		"name":             {"Usuário Teste"},
		"email":            {email},
		"password":         {password},
		"password_confirm": {password},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/register",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	for _, c := range w.Result().Cookies() {
		if c.Name == "taskflow_session" {
			return c
		}
	}
	t.Fatal("cookie de sessão não encontrado após registro")
	return nil
}

func TestProjectsIndexRequiresAuth(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("GET /projects sem auth: esperado 302, obteve %d", w.Code)
	}
}

func TestNewProjectPageRequiresAuth(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/projects/new", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("GET /projects/new sem auth: esperado 302, obteve %d", w.Code)
	}
}

func TestCreateProjectRequiresAuth(t *testing.T) {
	env := newTestEnv(t)

	form := url.Values{"title": {"Projeto"}, "description": {""}}
	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("POST /projects sem auth: esperado 302, obteve %d", w.Code)
	}
}

func TestCreateProjectWithValidDataRedirects(t *testing.T) {
	env := newTestEnv(t)
	cookie := loginAndGetCookie(t, env, "criador@test.com", "senha123")

	form := url.Values{
		"title":       {"Trabalho de Redes"},
		"description": {"Projeto para a disciplina"},
	}
	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("esperado redirect 302, obteve %d", w.Code)
	}

	projects, _ := env.projectRepo.FindByMemberID(1)
	if len(projects) == 0 {
		t.Error("projeto deveria ter sido criado")
	}
}

func TestCreateProjectValidationRejectsShortTitle(t *testing.T) {
	env := newTestEnv(t)
	cookie := loginAndGetCookie(t, env, "val@test.com", "senha123")

	form := url.Values{"title": {"AB"}, "description": {""}}
	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusFound {
		t.Error("título com < 3 caracteres não deveria ser aceito")
	}
}

func TestProjectMembershipLogic(t *testing.T) {
	env := newTestEnv(t)
	owner := env.seedUser("Owner", "owner@test.com", "hash")
	member := env.seedUser("Rafael", "rafael@test.com", "hash")
	project := env.seedProject("Projeto X", owner.ID)

	if !env.projectRepo.IsMember(project.ID, owner.ID) {
		t.Error("owner deveria ser membro automaticamente")
	}

	if env.projectRepo.IsMember(project.ID, member.ID) {
		t.Error("Rafael não deveria ser membro ainda")
	}

	if err := env.projectRepo.AddMember(project.ID, member.ID); err != nil {
		t.Fatalf("erro ao adicionar membro: %v", err)
	}
	if !env.projectRepo.IsMember(project.ID, member.ID) {
		t.Error("Rafael deveria ser membro após adição")
	}

	_ = env.projectRepo.RemoveMember(project.ID, member.ID)
	if env.projectRepo.IsMember(project.ID, member.ID) {
		t.Error("Rafael não deveria mais ser membro")
	}
}

func TestDeleteProjectRemovesFromRepo(t *testing.T) {
	env := newTestEnv(t)
	owner := env.seedUser("Owner", "del@test.com", "hash")
	project := env.seedProject("Para Deletar", owner.ID)

	if err := env.projectRepo.Delete(project.ID); err != nil {
		t.Fatalf("erro ao deletar: %v", err)
	}

	_, err := env.projectRepo.FindByID(project.ID)
	if err == nil {
		t.Error("projeto deveria ter sido excluído")
	}
}
