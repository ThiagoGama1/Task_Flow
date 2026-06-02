package middleware

import (
	"net/http"
	"taskflow/internal/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userRepo repositories.UserRepository
}

func NewAuthMiddleware(userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{userRepo: userRepo}
}

func (m *AuthMiddleware) LoadUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		val := session.Get("user_id")
		if val != nil {
			if id, ok := val.(int); ok {
				if user, err := m.userRepo.FindByID(uint(id)); err == nil {
					c.Set("current_user", user)
				}
			}
		}
		c.Next()
	}
}

func Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("current_user"); !exists {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
