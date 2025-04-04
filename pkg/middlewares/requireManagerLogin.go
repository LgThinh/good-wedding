package middlewares

import (
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"good-template-go/pkg/errors"
	"good-template-go/pkg/security"
	"good-template-go/pkg/utils"
	"good-template-go/pkg/utils/logger"
	"strings"
)

// AuthJWTMiddleware is a function that validates the jwt owner token
func AuthManagerJWTMiddleware() gin.HandlerFunc {
	log := logger.Tag("AuthManagerJWTMiddleware")
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.LogError(log, errors2.New("unauthorized"), "missing authorization header")
			appErr := errors.FeAppError(errors.VnMissingAuthorizationHeader, errors.MissingAuthorizationHeader)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.LogError(log, errors2.New("unauthorized"), "invalid authorization format")
			appErr := errors.FeAppError(errors.VnInvalidAuthorizationFormat, errors.InvalidAuthorizationFormat)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// Parse token
		claims, err := security.ParseManagerJWT(tokenString)
		if err != nil {
			logger.LogError(log, err, "error parsing token")
			appErr := errors.FeAppError(errors.VnTokenInvalid, errors.TokenInvalid)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		if claims.Role != utils.Manager {
			appErr := errors.FeAppError(errors.VnPermissionDenied, errors.PermissionDenied)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		c.Set("id", claims.ID)
		c.Set("role", claims.Role)
		c.Set("token_type", claims.TokenType)

		c.Next()
	}
}
