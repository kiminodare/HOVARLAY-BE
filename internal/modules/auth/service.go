package auth

import (
	"context"
	"errors"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	dtoAuth "github.com/kiminodare/HOVARLAY-BE/internal/modules/auth/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
	"os"
)

type Service struct {
	userService *user.Service
	jwtUtil     *utils.AESJWTUtil
}

func NewAuthService(userService *user.Service, jwtUtil *utils.AESJWTUtil) *Service {
	return &Service{
		userService: userService,
		jwtUtil:     jwtUtil,
	}
}

func (s *Service) Login(ctx context.Context, req *dtoAuth.Request) (*dtoAuth.Response, error) {
	userDetail, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	err = utils.ComparePassword(req.Password, userDetail.Password)
	if err = utils.ComparePassword(req.Password, userDetail.Password); err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	jwtUtils := utils.NewAESJWTUtil(os.Getenv("JWT_SECRET"), os.Getenv("AES_KEY"))
	token, err := jwtUtils.GenerateToken(userDetail.ID, userDetail.Email)
	if err != nil {
		return nil, err
	}

	return &dtoAuth.Response{
		Token: token,
	}, nil
}

func (s *Service) Register(ctx context.Context, req *dtoUser.Request) (*generated.User, error) {
	// Hash password dengan Argon2
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	userReq := &dtoUser.Request{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword, // Password sudah di-hash
	}

	createdUser, err := s.userService.CreateUser(ctx, userReq)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	return createdUser, nil
}
