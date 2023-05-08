package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"code-challenge/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Verify token success", func(t *testing.T) {
		mockJWTService := new(MockJWTService)
		logger := log.New()
		middleware := Middleware(mockJWTService, logger)

		r := SetUpRouter()
		r.GET("/middleware", middleware, func(c *gin.Context) {
			c.Status(200)
		})

		mockJWTService.On("ValidateToken", "token").Return(&jwt.Token{
			Valid: true,
			Claims: jwt.MapClaims{
				"id": "id",
			},
		}, nil)

		request, _ := http.NewRequest(http.MethodGet, "/middleware", nil)
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "token",
			MaxAge: 300,
		}
		request.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, 200, w.Code)

		mockJWTService.AssertExpectations(t)
	})

	t.Run("Verify token fails, token not provided", func(t *testing.T) {
		mockJWTService := new(MockJWTService)
		logger := log.New()
		middleware := Middleware(mockJWTService, logger)

		r := SetUpRouter()
		r.GET("/middleware", middleware, func(c *gin.Context) {
			c.Status(200)
		})

		request, _ := http.NewRequest(http.MethodGet, "/middleware", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, 401, w.Code)

		mockJWTService.AssertExpectations(t)
	})

	t.Run("Verify token fails, validate token fails", func(t *testing.T) {
		mockJWTService := new(MockJWTService)
		logger := log.New()
		middleware := Middleware(mockJWTService, logger)

		r := SetUpRouter()
		r.GET("/middleware", middleware, func(c *gin.Context) {
			c.Status(200)
		})

		mockJWTService.On("ValidateToken", "token").Return(nil, fmt.Errorf("mock error"))

		request, _ := http.NewRequest(http.MethodGet, "/middleware", nil)
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "token",
			MaxAge: 300,
		}
		request.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, 401, w.Code)

		mockJWTService.AssertExpectations(t)
	})

	t.Run("Verify token fails, nil claims", func(t *testing.T) {
		mockJWTService := new(MockJWTService)
		logger := log.New()
		middleware := Middleware(mockJWTService, logger)

		r := SetUpRouter()
		r.GET("/middleware", middleware, func(c *gin.Context) {
			c.Status(200)
		})

		mockJWTService.On("ValidateToken", "token").Return(&jwt.Token{
			Valid:  true,
			Claims: Claims{},
		}, nil)

		request, _ := http.NewRequest(http.MethodGet, "/middleware", nil)
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "token",
			MaxAge: 300,
		}
		request.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, 500, w.Code)

		mockJWTService.AssertExpectations(t)
	})

	t.Run("Verify token fails, user id not in claim", func(t *testing.T) {
		mockJWTService := new(MockJWTService)
		logger := log.New()
		middleware := Middleware(mockJWTService, logger)

		r := SetUpRouter()
		r.GET("/middleware", middleware, func(c *gin.Context) {
			c.Status(200)
		})

		mockJWTService.On("ValidateToken", "token").Return(&jwt.Token{
			Valid: true,
			Claims: jwt.MapClaims{
				"test": "test",
			},
		}, nil)

		request, _ := http.NewRequest(http.MethodGet, "/middleware", nil)
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "token",
			MaxAge: 300,
		}
		request.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)

		assert.Equal(t, 400, w.Code)

		mockJWTService.AssertExpectations(t)
	})
}
