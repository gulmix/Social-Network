package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gulmix/Social-Network/internal/config"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {

		}

		ctx := context.WithValue(c.Request.Context(), "config", cfg)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}
