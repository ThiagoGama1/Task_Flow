package handlers

import (
	"log/slog"
	"net/http"
	"taskflow/internal/models"
	"taskflow/internal/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo repositories.UserRepository
}

func NewAuthHandler(userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

type RegisterInput struct {
	Name            string `form:"name"             binding:"required,min=2,max=100"`
	Email           string `form:"email"            binding:"required,email"`
	Password        string `form:"password"         binding:"required,min=6"`
	PasswordConfirm string `form:"password_confirm" binding:"required"`
}

type LoginInput struct {
	Email    string `form:"email"    binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{"Title": "Login"})
}

func (h *AuthHandler) ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{"Title": "Cadastro"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "register", gin.H{
			"Title": "Cadastro",
			"Error": "Preencha todos os campos corretamente (senha mínimo 6 caracteres).",
			"Input": input,
		})
		return
	}

	if input.Password != input.PasswordConfirm {
		c.HTML(http.StatusUnprocessableEntity, "register", gin.H{
			"Title": "Cadastro",
			"Error": "As senhas não coincidem.",
			"Input": input,
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("bcrypt error", "error", err)
		c.HTML(http.StatusInternalServerError, "register", gin.H{
			"Title": "Cadastro",
			"Error": "Erro interno. Tente novamente.",
		})
		return
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	if err := h.userRepo.Create(user); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "register", gin.H{
			"Title": "Cadastro",
			"Error": "Este e-mail já está cadastrado.",
			"Input": input,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", int(user.ID))
	if err := session.Save(); err != nil {
		slog.Error("session save error", "error", err)
	}

	c.Redirect(http.StatusFound, "/projects")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "login", gin.H{
			"Title": "Login",
			"Error": "Preencha e-mail e senha.",
			"Input": input,
		})
		return
	}

	user, err := h.userRepo.FindByEmail(input.Email)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login", gin.H{
			"Title": "Login",
			"Error": "E-mail ou senha incorretos.",
			"Input": input,
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login", gin.H{
			"Title": "Login",
			"Error": "E-mail ou senha incorretos.",
			"Input": input,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", int(user.ID))
	if err := session.Save(); err != nil {
		slog.Error("session save error", "error", err)
	}

	c.Redirect(http.StatusFound, "/projects")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	c.Redirect(http.StatusFound, "/auth/login")
}
