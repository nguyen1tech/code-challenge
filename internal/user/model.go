package user

import "code-challenge/internal/entity"

type (
	RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	RegisterResponse struct {
		ID string `json:"id"`
	}

	GetByIDResponse struct {
		*entity.User
	}
)
