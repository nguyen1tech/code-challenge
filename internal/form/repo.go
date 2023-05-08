package form

import (
	"context"
	"fmt"

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

func (r *Repo) Create(_ context.Context, metadata *entity.Metadata) (string, error) {
	result := r.db.Create(&metadata)
	if result.Error != nil {
		return "", result.Error
	}
	return fmt.Sprint(metadata.ID), nil
}
