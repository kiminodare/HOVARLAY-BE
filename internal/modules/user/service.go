package user

import (
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"golang.org/x/net/context"
)

type Service struct {
	repo *Repository
}

func NewUserService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, userEntity *dtoUser.Request) (*generated.User, error) {
	return s.repo.Create(ctx, userEntity)
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*generated.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*generated.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
