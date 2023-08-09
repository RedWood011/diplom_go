package tests

import (
	"context"
	"os"
	"testing"

	"RedWood011/server/internal/config"
	"RedWood011/server/internal/database/postgres"
	"RedWood011/server/internal/services/user"
	usergrpc "RedWood011/server/internal/transport/grpc/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)

	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)

	tests := []struct {
		name    string
		query   string
		request *usergrpc.CreateUserRequest

		want    *usergrpc.TokenResponse
		wantErr bool
	}{
		{
			name: "User created successfully",
			request: &usergrpc.CreateUserRequest{
				Login:    "test123456789",
				Password: "test1234",
			},

			want: &usergrpc.TokenResponse{
				Status: "created",
			},
		},
		{
			name: "Failed to create user",
			request: &usergrpc.CreateUserRequest{
				Login:    "test",
				Password: "test",
			},

			want: &usergrpc.TokenResponse{
				Status: "invalid login or password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := grpcUsers.CreateUser(ctx, tt.request)

			assert.Equal(t, tt.want.Status, got.Status)
		})
	}

}

func TestAuthUser(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)

	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)

	tests := []struct {
		name    string
		query   string
		request *usergrpc.AuthUserRequest
		want    *usergrpc.TokenResponse
		wantErr bool
	}{
		{
			name: "User auth successfully",
			request: &usergrpc.AuthUserRequest{
				Login:    "test123456789",
				Password: "test1234",
			},
			want: &usergrpc.TokenResponse{
				Status: "ok",
			},
		},
		{
			name: "Failed to auth user",
			request: &usergrpc.AuthUserRequest{
				Login:    "test1234",
				Password: "test",
			},
			want: &usergrpc.TokenResponse{
				Status: "invalid login or password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := grpcUsers.AuthUser(ctx, tt.request)

			assert.Equal(t, tt.want.Status, got.Status)
		})
	}

}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)

	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)
	request := &usergrpc.CreateUserRequest{
		Login:    "test1111",
		Password: "test1234",
	}
	requestDeleted := &usergrpc.DeleteUserRequest{
		Login: "test1111",
	}
	got, err := grpcUsers.CreateUser(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, "created", got.Status)
	response, err := grpcUsers.DeleteUser(ctx, requestDeleted)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)

}
