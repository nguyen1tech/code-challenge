package auth

import (
	"context"
	"net/http"
	"strings"

	respErrors "code-challenge/internal/errors"
	"code-challenge/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Middleware(jwtService jwtService, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Errorf("Authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, respErrors.Unauthorized(""))
			return
		}

		strs := strings.Split(authHeader, "Bearer ")
		if len(strs) != 2 {
			logger.Errorf("Malformed Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, respErrors.Unauthorized(""))
			return
		}

		authToken := strs[1]
		token, err := jwtService.ValidateToken(authToken)
		if err != nil {
			logger.Errorf("Failed to validate token, err: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, respErrors.Unauthorized("Invalid token"))
			return
		}

		if !token.Valid {
			logger.Errorf("Invalid token provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, respErrors.Unauthorized("Invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Errorf("Failed to parse claims, err: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, respErrors.InternalServerError(""))
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			logger.Error("Failed to retrieve user id from claim")
			c.AbortWithStatusJSON(http.StatusBadRequest, respErrors.BadRequest(""))
			return
		}

		ctx := context.WithValue(c.Request.Context(), "userID", userID)
		c.Request = c.Request.WithContext(ctx)
	}
}
