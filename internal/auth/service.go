package auth

import (
	"context"
	"errors"

	"code-challenge/internal/entity"
	respErrors "code-challenge/internal/errors"
	"code-challenge/pkg/log"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type repo interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type jwtService interface {
	GenerateToken(identity Identity) (string, error)
	ValidateToken(authToken string) (*jwt.Token, error)
}

type Service struct {
	repo repo

	jwtSvc jwtService

	logger log.Logger
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
}

func NewService(repo repo, jwtSvc jwtService, logger log.Logger) *Service {
	return &Service{
		repo:   repo,
		jwtSvc: jwtSvc,
		logger: logger,
	}
}

// Authenticate authenticates a user by username and password
func (svc *Service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := svc.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			svc.logger.With(ctx).Error("Failed to get user by username: %s, error: %+v", username, err)
			return "", respErrors.NotFound("username not found")
		}
		svc.logger.With(ctx).Errorf("Failed to get user by username: %s, error: %+v", username, err)
		return "", respErrors.InternalServerError("")
	}

	if valid := checkPasswordHash(password, user.Password); !valid {
		svc.logger.With(ctx).Errorf("Failed to authenticate user: %s, error: incorrect username or password", username)
		return "", respErrors.Unauthorized("incorrect username or password")
	}

	signedToken, err := svc.jwtSvc.GenerateToken(user)
	if err != nil {
		svc.logger.With(ctx).Errorf("Failed to generate JWT, error: %s", err)
		return "", respErrors.InternalServerError("Failed to generate JWT")
	}
	return signedToken, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
