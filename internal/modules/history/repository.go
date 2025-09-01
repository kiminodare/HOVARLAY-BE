package history

import (
	"context"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated/history"
	user2 "github.com/kiminodare/HOVARLAY-BE/ent/generated/user"

	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
)

type Repository struct {
	client *generated.Client
}

func NewHistoryRepository(client *generated.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Create(
	ctx context.Context,
	userID uuid.UUID,
	text string,
	voice string,
	rate, pitch, volume float64,
) (*generated.History, error) {
	return r.client.History.Create().
		SetText(text).
		SetVoice(voice).
		SetRate(rate).
		SetPitch(pitch).
		SetVolume(volume).
		SetUserID(userID).
		Save(ctx)
}

func (r *Repository) GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*generated.History, error) {
	user, err := r.client.History.Query().
		Where(history.HasUserWith(user2.ID(userID))).
		Order(generated.Desc(history.FieldUpdatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	return r.client.History.Query().
		Where(history.HasUserWith(user2.ID(userID))).
		Count(ctx)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*generated.History, error) {
	return r.client.History.Get(ctx, id)
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, text string, voice string, rate, pitch, volume float64) error {
	return r.client.History.UpdateOneID(id).
		SetText(text).
		SetVoice(voice).
		SetRate(rate).
		SetPitch(pitch).
		SetVolume(volume).
		Exec(ctx)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.History.DeleteOneID(id).Exec(ctx)
}
