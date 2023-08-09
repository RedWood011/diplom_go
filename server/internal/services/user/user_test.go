package user

import (
	"context"
	"testing"

	"RedWood011/server/internal/apperrors"
	"RedWood011/server/internal/entity"
	"github.com/docker/distribution/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestCreateUserOK(t *testing.T) {
	ctrl := gomock.NewController(t)

	l := slog.Logger{}
	UserRepository := NewMockUserAdapter(ctrl)
	createUser := entity.User{ID: uuid.Generate().String(), Login: "Test1234", Password: "00000000"}

	UserRepository.EXPECT().SaveUser(gomock.Any(), gomock.Any()).Return(nil)

	userService := NewUserService(UserRepository, &l)
	err := userService.CreateUser(context.Background(), createUser)
	require.NoError(t, err)
}

func TestCreateUserInvalidLoginOrPassword(t *testing.T) {
	ctrl := gomock.NewController(t)

	l := slog.Logger{}
	UserRepository := NewMockUserAdapter(ctrl)

	userService := NewUserService(UserRepository, &l)

	testTable := []struct {
		name     string
		login    string
		password string
		err      error
	}{
		{
			name:     "Failed short login",
			login:    "t",
			password: "test",
			err:      apperrors.ErrAuth,
		},
		{
			name:     "Failed short password",
			login:    "tes",
			password: "",
			err:      apperrors.ErrAuth,
		},
	}
	for _, testCases := range testTable {
		t.Run(testCases.name, func(t *testing.T) {
			createUser := entity.User{ID: uuid.Generate().String(), Login: testCases.login, Password: testCases.password}
			err := userService.CreateUser(context.Background(), createUser)
			require.Equal(t, testCases.err, err)
		})
	}
}

func TestAuthUserOK(t *testing.T) {
	ctrl := gomock.NewController(t)

	l := slog.Logger{}
	UserRepository := NewMockUserAdapter(ctrl)
	userService := NewUserService(UserRepository, &l)

	mockUser := entity.User{
		ID:       uuid.Generate().String(),
		Login:    "validLogin",
		Password: "validPassword",
	}
	user := mockUser

	err := mockUser.SaveHashPassword()
	require.NoError(t, err)

	UserRepository.EXPECT().GetUser(gomock.Any(), user).Return(mockUser, nil)

	userID, err := userService.AuthUser(context.Background(), user)

	assert.Equal(t, mockUser.ID, userID)
	assert.NoError(t, err)
}

func TestAuthUserFailed(t *testing.T) {
	ctrl := gomock.NewController(t)

	l := slog.Logger{}
	UserRepository := NewMockUserAdapter(ctrl)
	userService := NewUserService(UserRepository, &l)

	mockUser := entity.User{
		ID:       uuid.Generate().String(),
		Login:    "Login",
		Password: "Password",
	}

	UserRepository.EXPECT().GetUser(gomock.Any(), mockUser).Return(entity.User{}, nil)

	userID, err := userService.AuthUser(context.Background(), mockUser)

	assert.NotEqual(t, mockUser.ID, userID)
	assert.Error(t, err)
}
