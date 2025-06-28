package repo_test

import (
	"auth/internal/entity"
	"auth/internal/repo"
	"auth/internal/repo/persistent"
	"context"
	"errors"
	"testing"
	"time"
)

func TestCreateUser(ctx context.Context, t *testing.T, newRepo func() *persistent.AuthRepo) {
	r := newRepo()

	tt := []struct {
		name string
		user entity.User
		want error
	}{
		{"no data", entity.User{}, nil},
		{
			name: "minimal data",
			user: entity.User{
				Login:        "Login1",
				PasswordHash: "1234567890",
			},
			want: nil,
		},
		{
			name: "maximal data",
			user: entity.User{
				ID:           1,
				Login:        "Login2",
				PasswordHash: "1234567890",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			want: nil,
		},
		{
			name: "duplicate login",
			user: entity.User{
				Login:        "Login1",
				PasswordHash: "1234567890",
			},
			want: repo.ErrDuplicateEntry,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := r.CreateUser(ctx, &tc.user)
			if tc.want == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
				return
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Errorf("want error: %v, got: %v", tc.want, err)
				return
			}
		})
	}
}

func TestFindByLogin(ctx context.Context, t *testing.T, newRepo func() *persistent.AuthRepo) {
	r := newRepo()

	err := r.CreateUser(ctx, &entity.User{
		Login:        "Login1",
		PasswordHash: "1234567890",
	})
	if err != nil {
		t.Fatal("test - FindByLogin - create user: ", err)
	}

	tt := []struct {
		name  string
		login string
		want  error
	}{
		{
			name:  "existing user",
			login: "Login1",
			want:  nil,
		},
		{
			name:  "non-existing user",
			login: "Login2",
			want:  repo.ErrNotFound,
		},
		{
			name:  "empty login",
			login: "",
			want:  repo.ErrNotFound,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			user, err := r.FindByLogin(ctx, tc.login)
			if tc.want == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
				return
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Errorf("want error: %v, got: %v", tc.want, err)
				return
			}

			if tc.login != user.Login {
				t.Errorf("want login: %s, got: %s", tc.login, user.Login)
				return
			}
		})
	}
}

func TestUpdatePasswordByID(ctx context.Context, t *testing.T, newRepo func() *persistent.AuthRepo) {
	r := newRepo()

	tUser := entity.User{
		ID:           0,
		Login:        "Login1",
		PasswordHash: "1234567890",
	}
	err := r.CreateUser(ctx, &tUser)
	if err != nil {
		t.Fatal("test - UpdatePasswordByID - create user: ", err)
	}

	createdUser, err := r.FindByLogin(ctx, tUser.Login)
	if err != nil {
		t.Fatal("test - UpdatePasswordByID - FindByLogin: ", err)
	}

	tt := []struct {
		name    string
		request repo.UpdatePasswordRequest
		want    error
	}{
		{
			name: "update for existing user",
			request: repo.UpdatePasswordRequest{
				ID:           createdUser.ID,
				PasswordHash: "0987654321",
			},
			want: nil,
		},
		{
			name: "update for non-existing user",
			request: repo.UpdatePasswordRequest{
				ID:           9999,
				PasswordHash: "0987654321",
			},
			want: repo.ErrNotFound,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := r.UpdatePasswordByID(ctx, &tc.request)
			if tc.want == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
				return
			}
		})
	}
}

func AssertUserEqual(t *testing.T, want, got *entity.User) {
	t.Helper()

	if want, got := want.ID, got.ID; want != got {
		t.Errorf("want ID: %d, got: %d", want, got)
	}

	if want, got := want.Login, got.Login; want != got {
		t.Errorf("want Login: %s, got: %s", want, got)
	}

	if want, got := want.PasswordHash, got.PasswordHash; want != got {
		t.Errorf("want PasswordHash: %s, got: %s", want, got)
	}

	if want, got := want.CreatedAt, got.CreatedAt; !want.Equal(got) {
		t.Errorf("want CreatedAt: %s, got: %s", want, got)
	}

	if want, got := want.UpdatedAt, got.UpdatedAt; !want.Equal(got) {
		t.Errorf("want UpdatedAt: %s, got: %s", want, got)
	}
}
