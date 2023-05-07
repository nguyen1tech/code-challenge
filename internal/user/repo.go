package user

import (
	"context"
	"fmt"
	"strconv"

	"code-challenge/internal/entity"

	"gorm.io/gorm"
)

type (
	Repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, user *entity.User) (string, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return "", result.Error
	}
	return fmt.Sprint(user.ID), nil
}

func (r *Repo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user = entity.User{}
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	intID, _ := strconv.ParseUint(id, 10, 64)
	var user = entity.User{ID: intID}
	result := r.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
