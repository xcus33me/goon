package repo_test

import (
	"auth/internal/entity"
	"auth/internal/repo"
	"auth/internal/repo/persistent"
	"auth/pkg/postgres"
	"context"
	"errors"
	"testing"
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
			if tc.want == nil && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("want error: %v, got: %v", tc.want, err)
			}

			// Проверяем, что ID был присвоен
			if tc.want == nil && tc.user.ID == 0 {
				t.Error("expected user ID to be set after creation")
			}
		})
	}
}

// Отдельный тест для дубликата
func TestCreateUser_DuplicateLogin(t *testing.T) {
	ctx := context.Background()
	r := NewTestAuthRepo(t)

	// Создаем первого пользователя
	user1 := entity.User{
		Login:        "duplicate_test",
		PasswordHash: "hash123",
	}
	err := r.CreateUser(ctx, &user1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Пытаемся создать пользователя с тем же логином
	user2 := entity.User{
		Login:        "duplicate_test",
		PasswordHash: "hash456",
	}
	err = r.CreateUser(ctx, &user2)

	// Проверяем, что получили ошибку дубликата
	if !errors.Is(err, repo.ErrDuplicateEntry) {
		t.Fatalf("expected ErrDuplicateEntry, got: %v", err)
	}
}

func TestFindByLogin(t *testing.T) {
	ctx := context.Background()
	r := NewTestAuthRepo(t)

	// Создаем тестового пользователя
	testUser := entity.User{
		Login:        "findme",
		PasswordHash: "secrethash",
	}
	err := r.CreateUser(ctx, &testUser)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

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

			if tc.want == nil {
				// Ожидаем успех
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if user == nil {
					t.Fatal("expected user, got nil")
				}
				if user.Login != tc.login {
					t.Fatalf("want login: %s, got: %s", tc.login, user.Login)
				}
				if user.PasswordHash != "secrethash" {
					t.Fatalf("password hash mismatch")
				}
			} else {
				// Ожидаем ошибку
				if !errors.Is(err, tc.want) {
					t.Fatalf("want error: %v, got: %v", tc.want, err)
				}
				if user != nil {
					t.Fatalf("expected nil user on error, got: %+v", user)
				}
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
		name    string
		request repo.UpdatePasswordRequest
		want    error
	}{
		{
			name: "update existing user",
			request: repo.UpdatePasswordRequest{
				ID:           createdUser.ID,
				PasswordHash: "newhash123",
			},
			want: nil,
		},
		{
			name: "update non-existing user",
			request: repo.UpdatePasswordRequest{
				ID:           99999,
				PasswordHash: "somehash",
			},
			want: repo.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			response, err := r.UpdatePasswordByID(ctx, &tc.request)

			if tc.want == nil {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if response == nil {
					t.Fatal("expected response, got nil")
				}
				if response.ID != tc.request.ID {
					t.Fatalf("expected ID %d, got %d", tc.request.ID, response.ID)
				}
			} else {
				if !errors.Is(err, tc.want) {
					t.Fatalf("want error: %v, got: %v", tc.want, err)
				}
				if response != nil {
					t.Fatalf("expected nil response on error, got: %+v", response)
				}
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
