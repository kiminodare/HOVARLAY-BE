package _interface

//go:generate mockgen -source=$GOFILE -destination=../mock/service_mock.go

import (
	"context"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
)

type ServiceInterface interface {
	Create(
		ctx context.Context,
		userID uuid.UUID,
		text string,
		voice string,
		rate, pitch, volume float64,
	) (*generated.History, error)

	GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*generated.History, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*generated.History, error)
	Update(ctx context.Context, id uuid.UUID, text string, voice string, rate, pitch, volume float64) error
	Delete(ctx context.Context, id uuid.UUID) error
}
