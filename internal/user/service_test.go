package user

import (
	"context"
	"fmt"
	"testing"

	"code-challenge/internal/entity"
	"code-challenge/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestService_Register(t *testing.T) {
	mockRegisterRequest := &RegisterRequest{
		Username: "test",
		Password: "123",
	}
	mockRepo := new(MockRepo)
	logger := log.New()
	svc := NewService(mockRepo, logger)

	t.Run("Register user success", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.On("GetByUsername", ctx, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", ctx, mock.Anything).Return("100", nil).Once()

		actualID, err := svc.Register(ctx, mockRegisterRequest)
		assert.NoError(t, err)
		assert.Equalf(t, "100", actualID, "expected id: %v but got: %v", "100", actualID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Register user when user already exists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("GetByUsername", ctx, mock.Anything).Return(&entity.User{ID: 100}, nil).Once()

		errMsg := "User already exists"
		_, err := svc.Register(ctx, mockRegisterRequest)
		assert.EqualErrorf(t, err, errMsg, "expected error: %v but got:  %v", errMsg, err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Register user when there is a database problem", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("GetByUsername", ctx, mock.Anything).
			Return(nil, fmt.Errorf("mock db issue")).Once()

		errMsg := "We encountered an error while processing your request."
		_, err := svc.Register(ctx, mockRegisterRequest)
		assert.EqualErrorf(t, err, errMsg, "expected error: %v but got:  %v", errMsg, err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetByID(t *testing.T) {
	mockRepo := new(MockRepo)
	logger := log.New()
	svc := NewService(mockRepo, logger)

	t.Run("GetByID success", func(t *testing.T) {
		ctx := context.Background()

		mockUser := &entity.User{ID: 100}
		mockRepo.On("GetByID", ctx, mock.Anything).Return(mockUser, nil).Once()

		user, err := svc.GetByID(ctx, "100")
		assert.NoError(t, err)
		assert.Equalf(t, mockUser, user, "expected user: %v but got user:  %v", mockUser, user)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID when user does not exist", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()

		errMsg := "User not found"
		_, err := svc.GetByID(ctx, "100")
		assert.EqualErrorf(t, err, errMsg, "expected error: %v but got:  %v", errMsg, err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID when there is a database problem", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, mock.Anything).Return(nil, fmt.Errorf("mock err")).Once()

		errMsg := "We encountered an error while processing your request."
		_, err := svc.GetByID(ctx, "100")
		assert.EqualErrorf(t, err, errMsg, "expected error: %v but got:  %v", errMsg, err.Error())

		mockRepo.AssertExpectations(t)
	})
}

// Repo is an autogenerated mock type for the repo type
type MockRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *MockRepo) Create(ctx context.Context, _a1 *entity.User) (string, error) {
	ret := _m.Called(ctx, _a1)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) (string, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) string); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *entity.User) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *MockRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: ctx, username
func (_m *MockRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
