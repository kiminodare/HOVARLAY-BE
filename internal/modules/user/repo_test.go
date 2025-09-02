package user_test

import (
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/testutils"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/user"
	dtoUser "github.com/kiminodare/HOVARLAY-BE/internal/modules/user/dto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserRepository_Create_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := faker.Name()
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// Create user
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, req.Email, created.Email)

	// Get by ID
	fetched, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
}

func TestUserRepository_Create_InvalidEmail(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	repo := user.NewUserRepository(client)
	email := "invalid-email"
	name := faker.Name()
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// Create user
	created, err := repo.Create(ctx, req)
	require.Error(t, err)
	require.Nil(t, created)
	require.Contains(t, err.Error(), "email")
}

func TestUserRepository_Create_InvalidName(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := "i"
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}
	// Create user
	created, err := repo.Create(ctx, req)
	require.Error(t, err)
	require.Nil(t, created)
	require.Contains(t, err.Error(), "min")
}

func TestUserRepository_Create_EmptyName(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	repo := user.NewUserRepository(client)
	email := faker.Email()
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     "",
		Email:    email,
		Password: password,
	}

	// Create user
	created, err := repo.Create(ctx, req)
	require.Error(t, err)
	require.Nil(t, created)
	require.Contains(t, err.Error(), "required")
}

func TestUserRepository_Create_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := faker.Name()
	password := "i"
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}
	// Create user
	created, err := repo.Create(ctx, req)
	require.Error(t, err)
	require.Nil(t, created)
	require.Contains(t, err.Error(), "min")
}

func TestUserRepository_Create_EmptyPassword(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)

	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := faker.Name()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: "", // ✅ Kosongkan password
	}

	// Create user
	created, err := repo.Create(ctx, req)
	require.Error(t, err)
	require.Nil(t, created)
	require.Contains(t, err.Error(), "required") // ✅ Check 'required' tag
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := faker.Name()
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}
	// Create user
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, req.Email, created.Email)
	// Create user with same email
	created, err = repo.Create(ctx, req)
	require.Error(t, err)
	require.Nil(t, created)
	require.Contains(t, err.Error(), "email")
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := faker.Name()
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}
	// Create user
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, req.Email, created.Email)
	// Get by ID
	fetched, err := repo.GetUserByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	// Get by ID
	fetched, err := repo.GetUserByID(ctx, uuid.New())
	require.Error(t, err)
	require.Nil(t, fetched)
	require.Contains(t, err.Error(), "not found")
}

func TestUserRepository_GetUserByEmail_Success(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	email := faker.Email()
	name := faker.Name()
	password := faker.Password()
	req := &dtoUser.Request{
		Name:     name,
		Email:    email,
		Password: password,
	}
	// Create user
	created, err := repo.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, req.Email, created.Email)
	// Get by Email
	fetched, err := repo.GetUserByEmail(ctx, created.Email)
	require.NoError(t, err)
	require.Equal(t, created.Email, fetched.Email)
	require.Equal(t, created.Name, fetched.Name)
	require.Equal(t, created.ID, fetched.ID)
}

func TestUserRepository_GetUserByEmail_NotFound(t *testing.T) {
	ctx := context.Background()
	client := testutils.NewTestDB(t)
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			t.Fatalf("failed closing connection: %v", err)
		}
	}(client)
	repo := user.NewUserRepository(client)
	// Get by Email
	fetched, err := repo.GetUserByEmail(ctx, faker.Email())
	require.Error(t, err)
	require.Nil(t, fetched)
	require.Contains(t, err.Error(), "not found")
}
