package tests_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"RedWood011/server/internal/authorization"
	"RedWood011/server/internal/config"
	"RedWood011/server/internal/database/postgres"
	"RedWood011/server/internal/services/secret"
	"RedWood011/server/internal/services/user"
	secretgrpc "RedWood011/server/internal/transport/grpc/secret"
	usergrpc "RedWood011/server/internal/transport/grpc/user"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/metadata"
)

func TestCreateSecret(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)
	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)
	secretService := secret.NewSecretService(repo, logger)
	grpcSecrets := secretgrpc.NewGrpcSecrets(secretService, cfg, logger)
	requestUser := &usergrpc.CreateUserRequest{
		Login:    "test2222",
		Password: "test1234",
	}
	got, err := grpcUsers.CreateUser(ctx, requestUser)
	assert.NoError(t, err)
	assert.Equal(t, "created", got.Status)
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %v", got.AccessToken)})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	var userID string
	userID, err = authorization.TokenValid(ctx, cfg.TokenConfig.SecretKey)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, authorization.UserKey("userID"), authorization.UserGUID(userID).String())
	assert.NoError(t, err)

	requestSecret := &secretgrpc.CreateSecretRequest{
		Name: []byte("test"),
		Data: []byte("test"),
	}
	res, err := grpcSecrets.CreateSecret(ctx, requestSecret)
	assert.NoError(t, err)
	assert.Equal(t, "created", res.Status)
}

func TestListSecret(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)
	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)
	secretService := secret.NewSecretService(repo, logger)
	grpcSecrets := secretgrpc.NewGrpcSecrets(secretService, cfg, logger)
	requestUser := &usergrpc.CreateUserRequest{
		Login:    "test2223",
		Password: "test1234",
	}
	got, err := grpcUsers.CreateUser(ctx, requestUser)
	assert.NoError(t, err)
	assert.Equal(t, "created", got.Status)
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %v", got.AccessToken)})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	var userID string
	userID, err = authorization.TokenValid(ctx, cfg.TokenConfig.SecretKey)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, authorization.UserKey("userID"), authorization.UserGUID(userID).String())
	assert.NoError(t, err)

	requestSecret := &secretgrpc.CreateSecretRequest{
		Name: []byte("test"),
		Data: []byte("test"),
	}
	res, err := grpcSecrets.CreateSecret(ctx, requestSecret)
	assert.NoError(t, err)
	assert.Equal(t, "created", res.Status)

	listSecret, err := grpcSecrets.ListSecrets(ctx, &secretgrpc.ListSecretsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, "ok", listSecret.Status)
	assert.Equal(t, "test", string(listSecret.Data[0].Name))
}

func TestGetSecret(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)
	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)
	secretService := secret.NewSecretService(repo, logger)
	grpcSecrets := secretgrpc.NewGrpcSecrets(secretService, cfg, logger)
	requestUser := &usergrpc.CreateUserRequest{
		Login:    "test2224",
		Password: "test1234",
	}
	got, err := grpcUsers.CreateUser(ctx, requestUser)
	assert.NoError(t, err)
	assert.Equal(t, "created", got.Status)
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %v", got.AccessToken)})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	var userID string
	userID, err = authorization.TokenValid(ctx, cfg.TokenConfig.SecretKey)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, authorization.UserKey("userID"), authorization.UserGUID(userID).String())
	assert.NoError(t, err)

	requestSecret := &secretgrpc.CreateSecretRequest{
		Name: []byte("name"),
		Data: []byte("test data"),
	}
	res, err := grpcSecrets.CreateSecret(ctx, requestSecret)
	assert.NoError(t, err)
	assert.Equal(t, "created", res.Status)

	requestGetSecret := &secretgrpc.GetSecretRequest{
		SecretId: res.SecretId,
	}
	getSecret, err := grpcSecrets.GetSecret(ctx, requestGetSecret)
	assert.NoError(t, err)
	assert.Equal(t, "ok", getSecret.Status)
	assert.Equal(t, "test data", string(getSecret.Data))
}

func TestDeleteSecret(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	assert.NoError(t, err)
	err = repo.Ping(ctx)
	assert.NoError(t, err)
	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)
	secretService := secret.NewSecretService(repo, logger)
	grpcSecrets := secretgrpc.NewGrpcSecrets(secretService, cfg, logger)
	requestUser := &usergrpc.CreateUserRequest{
		Login:    "test2225",
		Password: "test1234",
	}
	got, err := grpcUsers.CreateUser(ctx, requestUser)
	assert.NoError(t, err)
	assert.Equal(t, "created", got.Status)
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %v", got.AccessToken)})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	var userID string
	userID, err = authorization.TokenValid(ctx, cfg.TokenConfig.SecretKey)
	assert.NoError(t, err)
	ctx = context.WithValue(ctx, authorization.UserKey("userID"), authorization.UserGUID(userID).String())
	assert.NoError(t, err)

	requestSecret := &secretgrpc.CreateSecretRequest{
		Name: []byte("name"),
		Data: []byte("test data"),
	}
	res, err := grpcSecrets.CreateSecret(ctx, requestSecret)
	assert.NoError(t, err)
	assert.Equal(t, "created", res.Status)

	requestGetSecret := &secretgrpc.DeleteSecretRequest{
		SecretId: res.SecretId,
	}
	delSecret, err := grpcSecrets.DeleteSecret(ctx, requestGetSecret)
	assert.NoError(t, err)
	assert.Equal(t, "ok", delSecret.Status)
}
