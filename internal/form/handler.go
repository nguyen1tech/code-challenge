package form

import (
	"context"
	"mime/multipart"
	"net/http"

	"code-challenge/pkg/log"

	"github.com/gin-gonic/gin"
)

const fileLimit = 1024 * 1024 * 8

type service interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
}

type Handler struct {
	svc service

	logger log.Logger
}

func NewHandler(svc service, logger log.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func RegisterHandlers(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *Handler) {
	routerGroup.GET("/login", h.Login)
	routerGroup.GET("/upload", h.Upload)
	routerGroup.POST("/upload", authMiddleware, h.Upload)
}

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "info.html", gin.H{})
}

func (h *Handler) Login(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
}

func (h *Handler) Upload(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
		return
	}

	// Handle image upload form data
	_ = c.Request.ParseMultipartForm(10 << 20) // 10 MB
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusOK, "info.html", gin.H{"Message": "Failed to upload image, err: " + err.Error()})
		return
	}
	if file.Size > fileLimit {
		c.HTML(http.StatusOK, "info.html", gin.H{"Message": "File too large(8 MB max)"})
		return
	}

	path, err := h.svc.Upload(c.Request.Context(), file)
	if err != nil {
		c.HTML(http.StatusOK, "info.html", gin.H{"Message": "Failed to save image, err: " + err.Error()})
		return
	}

	c.HTML(http.StatusOK, "info.html", gin.H{"Message": "Uploaded file to: " + path})
}
