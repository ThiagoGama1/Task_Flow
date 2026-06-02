package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"taskflow/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func seedHashedUser(t *testing.T, env *testEnv, email, plainPassword string) *models.User {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	u := &models.User{Name: "Teste", Email: email, Password: string(hashed)}
	if err := env.userRepo.Create(u); err != nil {
		t.Fatal(err)
	}
	return u
}

func TestLoginPageExists(t *testing.T) {
	env := newTestEnv(t)
	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Errorf("rota GET /auth/login não encontrada (404)")
	}
}

func TestRegisterPageExists(t *testing.T) {
	env := newTestEnv(t)
	req := httptest.NewRequest(http.MethodGet, "/auth/register", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Errorf("rota GET /auth/register não encontrada (404)")
	}
}

func TestRegisterCreatesUserAndRedirects(t *testing.T) {
	env := newTestEnv(t)

	form := url.Values{
		"name":             {"Arthur Vieira"},
		"email":            {"arthur@test.com"},
		"password":         {"senha123"},
		"password_confirm": {"senha123"},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/register",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("esperado 302 após cadastro, obteve %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/projects" {
		t.Errorf("redirect esperado para /projects, obteve %q", loc)
	}

	// Confirma que o usuário foi persistido no mock
	user, err := env.userRepo.FindByEmail("arthur@test.com")
	if err != nil {
		t.Fatal("usuário não foi salvo após cadastro")
	}
	if user.Name != "Arthur Vieira" {
		t.Errorf("nome do usuário incorreto: %q", user.Name)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	env := newTestEnv(t)
	env.seedUser("Existente", "dup@test.com", "hash")

	form := url.Values{
		"name":             {"Novo"},
		"email":            {"dup@test.com"},
		"password":         {"senha123"},
		"password_confirm": {"senha123"},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/register",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusFound {
		t.Error("cadastro não deveria ser aceito com e-mail duplicado")
	}
}

func TestRegisterRejectsPasswordMismatch(t *testing.T) {
	env := newTestEnv(t)

	form := url.Values{
		"name":             {"Alguém"},
		"email":            {"alguem@test.com"},
		"password":         {"abc123"},
		"password_confirm": {"diferente"},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/register",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusFound {
		t.Error("cadastro não deveria prosseguir com senhas diferentes")
	}
}

func TestLoginSuccessRedirects(t *testing.T) {
	env := newTestEnv(t)
	seedHashedUser(t, env, "login@test.com", "senha123")

	form := url.Values{
		"email":    {"login@test.com"},
		"password": {"senha123"},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/login",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("esperado 302 após login, obteve %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/projects" {
		t.Errorf("redirect esperado para /projects, obteve %q", loc)
	}
}

func TestLoginFailsWithWrongPassword(t *testing.T) {
	env := newTestEnv(t)
	seedHashedUser(t, env, "user@test.com", "correta")

	form := url.Values{
		"email":    {"user@test.com"},
		"password": {"errada"},
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/login",
		strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	if w.Code == http.StatusFound {
		t.Error("login não deveria ser aceito com senha errada")
	}
}

func TestLogoutRequiresAuthentication(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()
	env.router.ServeHTTP(w, req)

	// Sem sessão ativa deve redirecionar para login
	if w.Code != http.StatusFound {
		t.Errorf("esperado redirect (302), obteve %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/auth/login" {
		t.Errorf("esperado redirect para /auth/login, obteve %q", loc)
	}
}
