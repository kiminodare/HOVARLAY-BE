package _interface

//go:generate mockgen -source=$GOFILE -destination=../mock/user_service_mock.go

import (
	"context"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
)

type ServiceInterface interface {
	Register(ctx context.Context, req *dtoUser.Request) (*generated.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*generated.User, error)
	GetUserByEmail(ctx context.Context, email string) (*generated.User, error)
}
