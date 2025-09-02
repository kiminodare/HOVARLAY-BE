// internal/modules/auth/test/service_test.go
package auth_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	dtoAuth "github.com/kiminodare/HOVARLAY-BE/internal/modules/auth/dto"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto" // Mock user service
	mockuser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/mock"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	MockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")

	authService := auth.NewAuthService(MockUserService, jwtUtil)

	// Setup mock behavior
	email := "test@example.com"
	password := "password123"
	userId := uuid.New()

	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	req := &dtoAuth.Request{
		Email:    email,
		Password: password,
	}

	userDetail := &generated.User{
		ID:       userId,
		Email:    email,
		Password: hashedPassword,
		Name:     "Test User",
	}

	// Setup mock behavior
	MockUserService.EXPECT().GetUserByEmail(gomock.Any(), email).Return(userDetail, nil)

	// Test
	response, err := authService.Login(context.Background(), req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotEmpty(t, response.Token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	MockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")
	authService := auth.NewAuthService(MockUserService, jwtUtil)

	req := &dtoAuth.Request{
		Email:    faker.Email(),
		Password: "password123",
	}

	// Setup mock behavior
	MockUserService.EXPECT().GetUserByEmail(gomock.Any(), req.Email).Return(nil, utils.ErrUserNotFound)

	// Test
	response, err := authService.Login(context.Background(), req)

	// Assertions
	require.Error(t, err)
	require.Nil(t, response)
	require.Equal(t, utils.ErrInvalidCredentials, err)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	MockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")
	authService := auth.NewAuthService(MockUserService, jwtUtil)

	email := "test@example.com"
	password := "correct-password"
	wrongPassword := "wrong-password"

	hashedPassword, _ := utils.HashPassword(password)
	userDetail := &generated.User{
		ID:       uuid.New(),
		Email:    email,
		Password: hashedPassword,
	}

	req := &dtoAuth.Request{
		Email:    email,
		Password: wrongPassword,
	}

	// mock userService return user with hashed password
	MockUserService.EXPECT().GetUserByEmail(gomock.Any(), email).Return(userDetail, nil)

	response, err := authService.Login(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, response)
	require.Equal(t, utils.ErrInvalidCredentials, err)
}

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	MockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")
	authService := auth.NewAuthService(MockUserService, jwtUtil)

	// Setup mock behavior
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

	// Setup mock behavior
	MockUserService.EXPECT().
		Register(gomock.Any(), gomock.Any()). // ✅ Accept any second parameter
		DoAndReturn(func(ctx context.Context, userReq *dtoUser.Request) (*generated.User, error) {
			require.NotEqual(t, req.Password, userReq.Password)
			require.True(t, len(userReq.Password) > 20)
			return expectedUser, nil
		})

	// Test
	result, err := authService.Register(context.Background(), req)

	// Assertions
	fmt.Println(result)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedUser.ID, result.ID)
	require.Equal(t, expectedUser.Email, result.Email)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")
	authService := auth.NewAuthService(mockUserService, jwtUtil)

	req := &dtoUser.Request{
		Name:     "Test User",
		Email:    "duplicate@example.com",
		Password: "password123",
	}

	// Setup mock expectation
	mockUserService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, utils.ErrEmailAlreadyExists)

	// Test
	result, err := authService.Register(context.Background(), req)

	// Assertions
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "already exists") // Auth service mungkin wrap error
}

func TestAuthService_Register_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")
	authService := auth.NewAuthService(mockUserService, jwtUtil)

	req := &dtoUser.Request{
		Name:     "", // ✅ Invalid data
		Email:    "invalid-email",
		Password: "123", // ✅ Too short
	}

	// Setup mock expectation untuk validation error
	mockUserService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, utils.ErrInvalidData) // ✅ Simulate validation error

	// Test
	result, err := authService.Register(context.Background(), req)

	// Assertions
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, utils.ErrInvalidData, err)
}

func TestAuthService_Register_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockuser.NewMockServiceInterface(ctrl)
	jwtUtil := utils.NewAESJWTUtil("test-secret", "test-aes-key-16-chars")
	authService := auth.NewAuthService(mockUserService, jwtUtil)

	req := &dtoUser.Request{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Setup mock expectation
	mockUserService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("database connection failed"))

	// Test
	result, err := authService.Register(context.Background(), req)

	// Assertions
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "database connection failed")
}
