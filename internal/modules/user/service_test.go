// internal/modules/user/test/service_test.go
package user_test

import (
	"context"
	"errors"
	"github.com/go-faker/faker/v4"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	mockrepo "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/mock"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestUserService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	req := &dtoUser.Request{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &generated.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: "hashed-password",
	}

	// Setup expectation
	mockRepo.EXPECT().
		Create(gomock.Any(), req).
		Return(expectedUser, nil)

	result, err := service.Register(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedUser.ID, result.ID)
	require.Equal(t, expectedUser.Email, result.Email)
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	req := &dtoUser.Request{
		Name:     faker.Name(),
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	// Setup expectation untuk constraint error
	mockRepo.EXPECT().
		Create(gomock.Any(), req).
		Return(nil, &generated.ConstraintError{})

	result, err := service.Register(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, utils.ErrEmailAlreadyExists, err)
}

func TestUserService_Register_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	req := &dtoUser.Request{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), req).
		Return(nil, errors.New("database error"))

	result, err := service.Register(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "database error")
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	userID := uuid.New()
	expectedUser := &generated.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Setup expectation
	mockRepo.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(expectedUser, nil)

	result, err := service.GetUserByID(context.Background(), userID)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedUser.ID, result.ID)
	require.Equal(t, expectedUser.Email, result.Email)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	userID := uuid.New()

	// Setup expectation
	mockRepo.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found"))

	result, err := service.GetUserByID(context.Background(), userID)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "user not found")
}

func TestUserService_GetUserByEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	email := "test@example.com"
	expectedUser := &generated.User{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: email,
	}

	// Setup expectation
	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), email).
		Return(expectedUser, nil)

	result, err := service.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedUser.Email, result.Email)
	require.Equal(t, expectedUser.ID, result.ID)
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockRepositoryInterface(ctrl)
	service := user.NewUserService(mockRepo)

	email := "nonexistent@example.com"

	// Setup expectation
	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), email).
		Return(nil, errors.New("user not found"))

	result, err := service.GetUserByEmail(context.Background(), email)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "user not found")
}
