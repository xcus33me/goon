package usecase_test

import (
	"auth/internal/entity"
	"auth/internal/usecase/auth"
	e "auth/pkg/errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

type test struct {
	name     string
	login    string
	password string
	mock     func()
	res      interface{}
	err      error
}

const testJwtSecret = "test-jwt-secret"

func authUsecase(t *testing.T) (*auth.UseCase, *MockAuthRepo, *MockAuthWebAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockAuthRepo(mockCtl)
	// webApi := NewMockAuthWebAPI(mockCtl) // not used

	useCase := auth.New(repo, testJwtSecret)

	return useCase, repo, nil // webapi not used
}

func TestLogin(t *testing.T) {
	t.Parallel()

	authUseCase, repo, _ := authUsecase(t)

	testUser := func(login, password string) *entity.User {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return &entity.User{
			ID:           1,
			Login:        login,
			PasswordHash: string(hash),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	}

	user := testUser("testuser", "correctpassword")

	tests := []test{
		{
			name:     "successful login",
			login:    "testuser",
			password: "correctpassword",
			mock: func() {
				repo.EXPECT().FindByLogin("testuser").Return(user, nil)
			},
			res: user,
			err: nil,
		},
		{
			name:     "user not found",
			login:    "unknownuser",
			password: "wrongpassword",
			mock: func() {
				repo.EXPECT().FindByLogin("unknownuser").Return(nil, e.UserNotFound)
			},
			res: nil,
			err: e.UserNotFound,
		},
		{
			name:     "incorrect password",
			login:    "testuser",
			password: "wrongpassword",
			mock: func() {
				repo.EXPECT().FindByLogin("testuser").Return(user, nil)
			},
			res: nil,
			err: e.WrongCredentials,
		},
	}

	for _, tc := range tests {
		ltc := tc
		t.Run(ltc.name, func(t *testing.T) {
			ltc.mock()

			user, token, err := authUseCase.Login(ltc.login, ltc.password)

			if ltc.res != nil {
				expectedUser := ltc.res.(*entity.User)

				// Verify fields
				require.NotNil(t, user, "User should not be nil on successful login")
				require.Equal(t, expectedUser.ID, user.ID, "User ID should match")
				require.Equal(t, expectedUser.Login, user.Login, "User login should match")
				require.NotEmpty(t, token, "Token should not be empty on successful login")

				// Verify JWT
				claims := &jwt.MapClaims{}
				parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(testJwtSecret), nil
				})
				require.NoError(t, err, "Token parsing should succeed")
				require.True(t, parsedToken.Valid, "Token should be valid")
				require.Equal(t, float64(expectedUser.ID), (*claims)["uid"], "Token UID should match user ID")
				require.Equal(t, expectedUser.Login, (*claims)["login"], "Token login should match user login")
			} else {
				require.Nil(t, user, "User should ne nil on failure")
				require.Empty(t, token, "Token should be empty on failure")
				require.ErrorIs(t, err, ltc.err, "Error should match expected error")
			}
		})
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	authUseCase, repo, _ := authUsecase(t)

	mockedCreationTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	existingUser := &entity.User{
		ID:        99,
		Login:     "existinguser",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	tests := []test{
		{
			name:     "successful registration",
			login:    "newuser123",
			password: "newpassword321",
			mock: func() {
				repo.EXPECT().FindByLogin("newuser123").Return(nil, e.UserNotFound)

				repo.EXPECT().CreateUser(gomock.Any()).
					DoAndReturn(func(user *entity.User) error {
						require.Equal(t, "newuser123", user.Login, "Login should be identical")
						require.NotEmpty(t, user.PasswordHash, "PasswordHash should not be empty")
						require.NotEqual(t, "newpassword321", user.PasswordHash, "Password should be hashed")

						user.ID = 123
						user.CreatedAt = mockedCreationTime
						user.UpdatedAt = mockedCreationTime
						return nil
					})
			},
			res: &entity.User{
				ID:        123,
				Login:     "newuser123",
				CreatedAt: mockedCreationTime,
				UpdatedAt: mockedCreationTime,
			},
			err: nil,
		},
		{
			name:     "user already exists",
			login:    "existinguser",
			password: "password",
			mock: func() {
				repo.EXPECT().FindByLogin("existinguser").Return(existingUser, nil)
			},
			res: nil,
			err: e.UserAlreadyExists,
		},
	}

	for _, tc := range tests {
		ltc := tc
		t.Run(ltc.name, func(t *testing.T) {
			ltc.mock()

			user, err := authUseCase.Register(ltc.login, ltc.password)
			if ltc.err != nil {
				require.ErrorIs(t, err, ltc.err, "Expected error '%v', received '%v'", ltc.err, err)
				require.Nil(t, user, "User should be nil in case of error")
			} else {
				require.NoError(t, err, "There should be no error, received: %v", err)
				require.NotNil(t, user, "User should not be nil on success.")

				expectedUser, ok := ltc.res.(*entity.User)
				require.True(t, ok, "Type of the expected result should be *entity.User")

				require.Equal(t, expectedUser.ID, user.ID, "User's ID doesn't match")
				require.Equal(t, expectedUser.Login, user.Login, "User's login does not match")

				require.NotEmpty(t, user.PasswordHash, "PasswordHash should not be empty")
				errCompare := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(ltc.password))
				require.NoError(t, errCompare, "Hash of the password does not match the original password")

				require.True(t, expectedUser.CreatedAt.Equal(user.CreatedAt), "CreatedAt does not match. Expected: %v, Received: %v", expectedUser.CreatedAt, user.CreatedAt)
				require.True(t, expectedUser.UpdatedAt.Equal(user.UpdatedAt), "UpdatedAt does not match. Expected: %v, Received: %v", expectedUser.UpdatedAt, user.UpdatedAt)
			}
		})
	}
}
