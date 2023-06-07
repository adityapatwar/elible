package middleware

import (
	"strings"

	"elible/internal/app/services"
	"elible/internal/app/utils"
	"elible/internal/config"
	errors "elible/internal/pkg"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware(cfg *config.Config, tokenService *services.AdminService, useLocalValidation bool, useDatabaseValidation bool, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			errors.WriteErrorResponse(c.Writer, 401, "Authorization header required")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			errors.WriteErrorResponse(c.Writer, 401, "Invalid Authorization header format")
			c.Abort()
			return
		}
		token := tokenParts[1]

		if useLocalValidation {
			_, err := utils.ValidateJWTWithLocalSecret(token, cfg)
			if err != nil {
				errors.WriteErrorResponse(c.Writer, 401, "Invalid Authorization token (local validation)")
				c.Abort()
				return
			}
		}

		if useDatabaseValidation {
			dbToken, err := tokenService.GetAdminByToken(token)
			if err != nil || dbToken == nil {
				errors.WriteErrorResponse(c.Writer, 401, "Invalid Authorization token (database validation)")
				c.Abort()
				return
			}
		}

		next(c)
		c.Next()
	}
}
