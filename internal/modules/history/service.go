package history

import (
	"context"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	userID uuid.UUID,
	text string,
	voice string,
	rate, pitch, volume float64,
) (*generated.History, error) {
	return s.repo.Create(ctx, userID, text, voice, rate, pitch, volume)
}

func (s *Service) GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*generated.History, error) {
	return s.repo.GetByUser(ctx, userID, offset, limit)
}

func (s *Service) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.CountByUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*generated.History, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, text string, voice string, rate, pitch, volume float64) error {
	return s.repo.Update(ctx, id, text, voice, rate, pitch, volume)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
