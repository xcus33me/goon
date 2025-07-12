package repo_test

import (
	"auth/internal/entity"
	"auth/internal/repo"
	"auth/internal/repo/persistent"
	"auth/pkg/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestAuthRepo(t *testing.T) *persistent.AuthRepo {
	t.Helper()

	embeddedDB := postgres.NewEmbedded(t)

	pg := &postgres.Postgres{Db: embeddedDB.DB}
	return persistent.New(pg)
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	r := NewTestAuthRepo(t)

	tests := []struct {
		name string
		user entity.User
		want error
	}{
		{
			name: "minimal data",
			user: entity.User{
				Login:        "testuser1",
				PasswordHash: "hash123",
			},
			want: nil,
		},
		{
			name: "another user",
			user: entity.User{
				Login:        "testuser2",
				PasswordHash: "hash456",
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := r.CreateUser(ctx, &tc.user)
			//if tc.want == nil && err != nil {
			//	t.Fatalf("expected no error, got: %v", err)
			//}
			//if tc.want != nil && !errors.Is(err, tc.want) {
			//	t.Fatalf("want error: %v, got: %v", tc.want, err)
			//}
			//
			//// Проверяем, что ID был присвоен
			//if tc.want == nil && tc.user.ID == 0 {
			//	t.Error("expected user ID to be set after creation")
			//}
			require.NoError(t, err, "should be succeed")
			assert.NotZero(t, tc.user.ID, "ID should be set")
			assert.Equal(t, tc.user.Login, tc.user.Login, "Login should match")
			assert.Equal(t, tc.user.PasswordHash, tc.user.PasswordHash, "PasswordHash should match")
			assert.NotZero(t, tc.user.CreatedAt, "CreatedAt should be set")
			assert.NotZero(t, tc.user.UpdatedAt, "UpdatedAt should be set")
		})
	}
}

func TestCreateUser_DuplicateLogin(t *testing.T) {
	ctx := context.Background()
	r := NewTestAuthRepo(t)

	user1 := entity.User{
		Login:        "duplicate_test",
		PasswordHash: "hash123",
	}
	err := r.CreateUser(ctx, &user1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}
	user2 := entity.User{
		Login:        "duplicate_test",
		PasswordHash: "hash456",
	}
	err = r.CreateUser(ctx, &user2)

	assert.ErrorIs(t, err, repo.ErrDuplicateEntry, "should return duplicate error")
}

func TestFindByLogin(t *testing.T) {
	ctx := context.Background()
	r := NewTestAuthRepo(t)

	testUser := entity.User{
		Login:        "findme",
		PasswordHash: "secrethash",
	}
	err := r.CreateUser(ctx, &testUser)
	require.NoError(t, err, "test user creation should be succeed")

	tests := []struct {
		name  string
		login string
		want  error
	}{
		{
			name:  "existing user",
			login: "findme",
			want:  nil,
		},
		{
			name:  "non-existing user",
			login: "notfound",
			want:  repo.ErrNotFound,
		},
		{
			name:  "empty login",
			login: "",
			want:  repo.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			user, err := r.FindByLogin(ctx, tc.login)

			if tc.want != nil {
				assert.ErrorIs(t, err, tc.want, "should return error")
				assert.ErrorIs(t, err, repo.ErrNotFound, "should return expected error")
				assert.Nil(t, user, "user should be nil")
			} else {
				require.NoError(t, err, "should not return an error")
				require.NotNil(t, user, "user should not be nil")

				assert.Equal(t, tc.login, user.Login, "login should match")
				assert.Equal(t, "secrethash", user.PasswordHash, "password hash should match")
				assert.NotZero(t, user.ID, "user ID should be set")
			}
		})
	}
}

func TestUpdatePasswordByID(t *testing.T) {
	ctx := context.Background()
	r := NewTestAuthRepo(t)

	// Создаем тестового пользователя
	testUser := entity.User{
		Login:        "updateme",
		PasswordHash: "oldhash",
	}
	err := r.CreateUser(ctx, &testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Получаем созданного пользователя для получения ID
	createdUser, err := r.FindByLogin(ctx, testUser.Login)
	if err != nil {
		t.Fatalf("failed to find created user: %v", err)
	}

	tests := []struct {
		name        string
		userID      int64
		newPassword string
		want        error
	}{
		{
			name:        "update existing user",
			userID:      createdUser.ID,
			newPassword: "newhash123",
			want:        nil,
		},
		{
			name:        "update non-existing user",
			userID:      99999,
			newPassword: "somehash",
			want:        repo.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			request := repo.UpdatePasswordRequest{
				ID:           tc.userID,
				PasswordHash: tc.newPassword,
			}

			response, err := r.UpdatePasswordByID(ctx, &request)

			if tc.want != nil {
				assert.Error(t, err, "Should return an error")
				assert.ErrorIs(t, err, tc.want, "Should return expected error type")
				assert.Nil(t, response, "Response should be nil on error")
			} else {
				require.NoError(t, err, "Should not return an error")
				require.NotNil(t, response, "Response should not be nil")

				assert.Equal(t, tc.userID, response.ID, "Response ID should match request ID")
				assert.NotZero(t, response.UpdatedAt, "UpdatedAt should be set")
			}
		})
	}
}

// AssertUserEqual помогает сравнивать пользователей в тестах
func AssertUserEqual(t *testing.T, want, got *entity.User) {
	t.Helper()

	if want.ID != got.ID {
		t.Errorf("want ID: %d, got: %d", want.ID, got.ID)
	}

	if want.Login != got.Login {
		t.Errorf("want Login: %s, got: %s", want.Login, got.Login)
	}

	if want.PasswordHash != got.PasswordHash {
		t.Errorf("want PasswordHash: %s, got: %s", want.PasswordHash, got.PasswordHash)
	}

	if !want.CreatedAt.IsZero() && !want.CreatedAt.Equal(got.CreatedAt) {
		t.Errorf("want CreatedAt: %s, got: %s", want.CreatedAt, got.CreatedAt)
	}

	if !want.UpdatedAt.IsZero() && !want.UpdatedAt.Equal(got.UpdatedAt) {
		t.Errorf("want UpdatedAt: %s, got: %s", want.UpdatedAt, got.UpdatedAt)
	}
}
