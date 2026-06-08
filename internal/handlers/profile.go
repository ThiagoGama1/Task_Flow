package handlers

import (
	"log/slog"
	"net/http"
	"taskflow/internal/models"
	"taskflow/internal/repositories"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	userRepo repositories.UserRepository
}

func NewProfileHandler(userRepo repositories.UserRepository) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

type ProfileInput struct {
	Name            string `form:"name"             binding:"required,min=2,max=100"`
	CurrentPassword string `form:"current_password" binding:"required"`
	NewPassword     string `form:"new_password"`
	PasswordConfirm string `form:"password_confirm"`
}

func (h *ProfileHandler) Show(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)
	c.HTML(http.StatusOK, "profile", gin.H{
		"Title": "Meu Perfil",
		"User":  user,
	})
}

func (h *ProfileHandler) Update(c *gin.Context) {
	user := c.MustGet("current_user").(*models.User)

	var input ProfileInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "profile", gin.H{
			"Title": "Meu Perfil",
			"User":  user,
			"Error": "Nome obrigatório (mínimo 2 caracteres) e senha atual são obrigatórios.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.HTML(http.StatusUnprocessableEntity, "profile", gin.H{
			"Title": "Meu Perfil",
			"User":  user,
			"Error": "Senha atual incorreta.",
		})
		return
	}

	user.Name = input.Name

	if input.NewPassword != "" {
		if len(input.NewPassword) < 6 {
			c.HTML(http.StatusUnprocessableEntity, "profile", gin.H{
				"Title": "Meu Perfil",
				"User":  user,
				"Error": "A nova senha deve ter pelo menos 6 caracteres.",
			})
			return
		}
		if input.NewPassword != input.PasswordConfirm {
			c.HTML(http.StatusUnprocessableEntity, "profile", gin.H{
				"Title": "Meu Perfil",
				"User":  user,
				"Error": "As novas senhas não coincidem.",
			})
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("bcrypt error", "error", err)
			c.HTML(http.StatusInternalServerError, "profile", gin.H{
				"Title": "Meu Perfil",
				"User":  user,
				"Error": "Erro interno. Tente novamente.",
			})
			return
		}
		user.Password = string(hashed)
	}

	if err := h.userRepo.Update(user); err != nil {
		c.HTML(http.StatusInternalServerError, "profile", gin.H{
			"Title": "Meu Perfil",
			"User":  user,
			"Error": "Erro ao atualizar perfil.",
		})
		return
	}

	c.HTML(http.StatusOK, "profile", gin.H{
		"Title":   "Meu Perfil",
		"User":    user,
		"Success": "Perfil atualizado com sucesso!",
	})
}
