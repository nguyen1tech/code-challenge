package user

import (
	"context"
	"errors"
	"net/http"

	"code-challenge/internal/entity"
	respErrors "code-challenge/internal/errors"
	"code-challenge/pkg/log"

	"github.com/gin-gonic/gin"
)

type (
	service interface {
		Register(ctx context.Context, req *RegisterRequest) (string, error)
		GetByID(ctx context.Context, id string) (*entity.User, error)
	}

	Handler struct {
		svc service

		logger log.Logger
	}
)

func NewHandler(svc service, logger log.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func RegisterHandlers(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *Handler) {
	routerGroup.POST("/register", h.Register)
	routerGroup.GET("/me", authMiddleware, h.Me)
}

// Register registers a new user with the given credentials.
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.With(c.Request.Context()).Errorf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, respErrors.BadRequest(""))
		c.Abort()
		return
	}

	id, err := h.svc.Register(c.Request.Context(), &req)
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

	resp := RegisterResponse{
		ID: id,
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Request.Context().Value("userID").(string)
	user, err := h.svc.GetByID(c.Request.Context(), userID)
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

	c.JSON(http.StatusOK, GetByIDResponse{user})
}
