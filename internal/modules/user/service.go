package user

import (
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user/interface"
	"github.com/kiminodare/HOVARLAY-BE/internal/utils"
	"golang.org/x/net/context"
)

type Service struct {
	repo _interface.RepositoryInterface
}

var _ _interface.ServiceInterface = (*Service)(nil)

func NewUserService(repo _interface.RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req *dtoUser.Request) (*generated.User, error) {
	user, err := s.repo.Create(ctx, req)
	if err != nil {
		if generated.IsConstraintError(err) {
			return nil, utils.ErrEmailAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*generated.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*generated.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
