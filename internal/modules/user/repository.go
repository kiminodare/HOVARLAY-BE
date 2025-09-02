package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated/user"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user/interface"
)

type Repository struct {
	client *generated.Client
}

var _ _interface.RepositoryInterface = (*Repository)(nil)

func NewUserRepository(client *generated.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Create(ctx context.Context, userEntity *dtoUser.Request) (*generated.User, error) {
	if err := userEntity.Validate(); err != nil {
		return nil, err
	}

	return r.client.User.Create().
		SetName(userEntity.Name).
		SetEmail(userEntity.Email).
		SetPassword(userEntity.Password).
		Save(ctx)
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*generated.User, error) {
	userDetail, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return userDetail, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*generated.User, error) {
	userDetail, err := r.client.User.Query().Where(user.Email(email)).First(ctx)
	if err != nil {
		return nil, err
	}
	return userDetail, nil
}
