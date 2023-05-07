package auth

import (
	"context"
	"errors"
	"net/http"

	respErrors "code-challenge/internal/errors"
	"code-challenge/pkg/log"

	"github.com/gin-gonic/gin"
)

type service interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type Handler struct {
	svc service

	logger log.Logger
}

func NewHandler(svc service, logger log.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func RegisterHandlers(routerGroup *gin.RouterGroup, h *Handler) {
	routerGroup.POST("/login", h.Login)
}

func (h *Handler) Login(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		h.logger.With(c.Request.Context()).Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	token, err := h.svc.Authenticate(c.Request.Context(), credentials.Username, credentials.Password)
	if err != nil {
		var respErr *respErrors.ErrorResponse
		switch {
		case errors.As(err, &respErr):
			c.JSON(respErr.Status, respErr)
		default:
			c.JSON(http.StatusInternalServerError, respErrors.InternalServerError(""))
		}
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{token})
}
