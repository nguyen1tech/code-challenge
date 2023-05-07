package user

import (
	"context"
	"errors"
	"time"

	"code-challenge/internal/entity"
	respErrors "code-challenge/internal/errors"
	"code-challenge/pkg/log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type repo interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) (string, error)
}

type Service struct {
	repo repo

	logger log.Logger
}

func NewService(repo repo, logger log.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (svc *Service) Register(ctx context.Context, req *RegisterRequest) (string, error) {
	svc.logger.With(ctx).Info("Registering user: ", req.Username)

	// TODO: Validate password constraint such as length, number, characters, special chars, etc.
	user, err := svc.repo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		svc.logger.With(ctx).Error("Failed to get user by username: %s, error: %+v", req.Username, err)
		return "", respErrors.InternalServerError("")
	}
	// The user already exists
	if user != nil {
		svc.logger.With(ctx).Errorf("User %s already exists", req.Username)
		return "", respErrors.BadRequest("User already exists")
	}
	// Encrypt user's password
	hash, err := hashPassword(req.Password)
	if err != nil {
		svc.logger.With(ctx).Errorf("Failed to hash password for user: %s, error: %+v", req.Username, err)
		return "", respErrors.InternalServerError("")
	}

	now := time.Now()
	return svc.repo.Create(ctx, &entity.User{
		Username: req.Username,
		Password: hash,

		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (svc *Service) GetByID(ctx context.Context, id string) (*entity.User, error) {
	svc.logger.With(ctx).Info("Getting user by id: ", id)
	user, err := svc.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, respErrors.NotFound("User not found")
		}
		return nil, respErrors.InternalServerError("")
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
