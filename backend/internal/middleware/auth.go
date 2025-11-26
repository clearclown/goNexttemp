package middleware

import (
	"net/http"
	"strings"

	"github.com/ablaze/gonexttemp-backend/internal/auth"
	"github.com/ablaze/gonexttemp-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ContextUserID       = "userID"
	ContextUserEmail    = "userEmail"
)

func AuthMiddleware(jwt *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				response.CodeUnauthorized,
				"Authorization header is required",
			))
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				response.CodeUnauthorized,
				"Invalid authorization header format",
			))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			if err == auth.ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
					response.CodeTokenExpired,
					"Access token has expired",
				))
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(
				response.CodeTokenInvalid,
				"Invalid access token",
			))
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserEmail, claims.Email)
		c.Next()
	}
}
