package test

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	historyData "github.com/kiminodare/HOVARLAY-BE/ent/generated/history"
	user2 "github.com/kiminodare/HOVARLAY-BE/ent/generated/user"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/history"
	dtoHistory "github.com/kiminodare/HOVARLAY-BE/internal/modules/history/dto"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/testutils"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHistoryRepository_Create_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	// Setup user
	userRepo := user.NewUserRepository(client)
	testUser, err := userRepo.Create(ctx, &dtoUser.Request{
		Email:    fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8]),
		Name:     faker.Name(),
		Password: faker.Password(),
	})
	require.NoError(t, err)

	// Test Create Success
	historyRepo := history.NewHistoryRepository(client)

	data := dtoHistory.CreateHistoryRequest{}
	err = faker.FakeData(&data, options.WithTagName("custom"))
	if err != nil {
		t.Fatal(err)
	}

	created, err := historyRepo.Create(ctx, testUser.ID, data.Text, data.Voice, data.Rate, data.Pitch, data.Volume)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, data.Text, created.Text)
	require.Equal(t, data.Voice, created.Voice)
	require.Equal(t, data.Rate, created.Rate)
	require.Equal(t, data.Pitch, created.Pitch)
	require.Equal(t, data.Volume, created.Volume)
	fetchedHistories, err := client.History.Query().
		Where(historyData.HasUserWith(user2.ID(testUser.ID))).
		All(ctx)
	require.NoError(t, err)
	require.Len(t, fetchedHistories, 1)
	require.Equal(t, testUser.ID, fetchedHistories[0].QueryUser().Where(user2.ID(testUser.ID)).OnlyX(ctx).ID)
}

func TestHistoryRepository_Create_InvalidUserID(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	historyRepo := history.NewHistoryRepository(client)

	randomUUID, err := uuid.Parse(faker.UUIDHyphenated())
	require.NoError(t, err)
	require.NotNil(t, randomUUID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = historyRepo.Create(ctx, randomUUID, "text", "voice", 1.0, 1.0, 1.0)
	require.Error(t, err)
}

func TestHistoryRepository_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	// Setup
	userRepo := user.NewUserRepository(client)
	testUser, err := userRepo.Create(ctx, &dtoUser.Request{
		Email:    faker.Email(),
		Name:     faker.Name(),
		Password: faker.Password(),
	})
	require.NoError(t, err)

	historyRepo := history.NewHistoryRepository(client)
	created, err := historyRepo.Create(ctx, testUser.ID, "test", "voice", 1.0, 1.0, 1.0)
	require.NoError(t, err)

	// Test GetByID Success
	found, err := historyRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, found.ID)
	require.Equal(t, "test", found.Text)
}

func TestHistoryRepository_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	historyRepo := history.NewHistoryRepository(client)
	randomUUID, err := uuid.Parse(faker.UUIDHyphenated())
	require.NoError(t, err)
	require.NotNil(t, randomUUID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = historyRepo.GetByID(ctx, randomUUID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "history not found") // atau pesan error ent
}

func TestHistoryRepository_GetByUser_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	// Setup user
	userRepo := user.NewUserRepository(client)
	testUser, err := userRepo.Create(ctx, &dtoUser.Request{
		Email:    faker.Email(),
		Name:     faker.Name(),
		Password: faker.Password(),
	})
	require.NoError(t, err)

	// Buat beberapa history
	historyRepo := history.NewHistoryRepository(client)
	for i := 0; i < 3; i++ {
		_, err := historyRepo.Create(ctx, testUser.ID, faker.Sentence(), faker.Word(), 1.0, 1.0, 1.0)
		require.NoError(t, err)
	}

	// Test GetByUser
	results, err := historyRepo.GetByUser(ctx, testUser.ID, 0, 10)
	require.NoError(t, err)
	require.Len(t, results, 3)
}

func TestHistoryRepository_GetByUser_Empty(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	historyRepo := history.NewHistoryRepository(client)

	// Test dengan user yang tidak punya history
	results, err := historyRepo.GetByUser(ctx, uuid.New(), 0, 10)
	require.NoError(t, err)
	require.Len(t, results, 0)
}

func TestHistoryRepository_Update_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	// Setup
	userRepo := user.NewUserRepository(client)
	testUser, err := userRepo.Create(ctx, &dtoUser.Request{
		Email:    faker.Email(),
		Name:     faker.Name(),
		Password: faker.Password(),
	})
	require.NoError(t, err)

	historyRepo := history.NewHistoryRepository(client)
	created, err := historyRepo.Create(ctx, testUser.ID, "old text", "old voice", 1.0, 1.0, 1.0)
	require.NoError(t, err)

	data := dtoHistory.CreateHistoryRequest{}
	err = faker.FakeData(&data, options.WithTagName("custom"))
	if err != nil {
		t.Fatal(err)
	}
	// Test Update Success
	err = historyRepo.Update(ctx, created.ID, data.Text, data.Voice, data.Rate, data.Pitch, data.Volume)
	require.NoError(t, err)

	// Verify update
	updated, err := historyRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, data.Text, updated.Text)
	require.Equal(t, data.Voice, updated.Voice)
	require.Equal(t, data.Rate, updated.Rate)
	require.Equal(t, data.Pitch, updated.Pitch)
}

func TestHistoryRepository_Update_NotFound(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	historyRepo := history.NewHistoryRepository(client)
	data := dtoHistory.CreateHistoryRequest{}
	err := faker.FakeData(&data)
	if err != nil {
		t.Fatal(err)
	}
	// Test update dengan ID yang tidak ada
	err = historyRepo.Update(ctx, uuid.New(), data.Text, data.Voice, data.Rate, data.Pitch, data.Volume)
	require.Error(t, err)
	require.Error(t, err)
}

func TestHistoryRepository_Delete_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	// Setup
	userRepo := user.NewUserRepository(client)
	testUser, err := userRepo.Create(ctx, &dtoUser.Request{
		Email:    faker.Email(),
		Name:     faker.Name(),
		Password: faker.Password(),
	})
	require.NoError(t, err)

	historyRepo := history.NewHistoryRepository(client)
	data := dtoHistory.CreateHistoryRequest{}
	err = faker.FakeData(&data, options.WithTagName("custom"))
	if err != nil {
		t.Fatal(err)
	}
	created, err := historyRepo.Create(ctx, testUser.ID, data.Text, data.Voice, data.Rate, data.Pitch, data.Volume)
	require.NoError(t, err)
	require.NoError(t, err)

	// Test Delete Success
	err = historyRepo.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = historyRepo.GetByID(ctx, created.ID)
	require.Error(t, err)
}

func TestHistoryRepository_Delete_NotFound(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	historyRepo := history.NewHistoryRepository(client)

	// Test delete dengan ID yang tidak ada
	err := historyRepo.Delete(ctx, uuid.New())
	require.Error(t, err)
}
