package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gulmix/Social-Network/internal/config"
	"github.com/gulmix/Social-Network/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try Authorization header first, then ?token= query param (for WebSocket upgrades).
		tokenString := utils.ExtractTokenFromHeader(c.GetHeader("Authorization"))
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			ctx := context.WithValue(c.Request.Context(), "config", cfg)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, "config", cfg)
		c.Request = c.Request.WithContext(ctx)

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

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

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}
