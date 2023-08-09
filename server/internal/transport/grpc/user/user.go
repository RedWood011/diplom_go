package user

import (
	"context"

	"RedWood011/server/internal/authorization"
	"RedWood011/server/internal/config"
	"RedWood011/server/internal/entity"
	"github.com/docker/distribution/uuid"
	"golang.org/x/exp/slog"
)

type UserService interface {
	CreateUser(ctx context.Context, user entity.User) error
	AuthUser(ctx context.Context, user entity.User) (string, error)
	DeleteUser(ctx context.Context, login string) error
}

// NewGrpcUsers функция создания обраточка запросов для пользователей
func NewGrpcUsers(userService UserService, config *config.Config, logger *slog.Logger) *GrpcUsers {
	return &GrpcUsers{
		userService: userService,
		cfg:         config,
		logger:      logger,
	}
}

// GrpcUsers структура для обраточика запросов для пользователя
type GrpcUsers struct {
	UnimplementedUsersServer
	userService UserService
	cfg         *config.Config
	logger      *slog.Logger
}

func (gh *GrpcUsers) CreateUser(ctx context.Context, in *CreateUserRequest) (*TokenResponse, error) {
	userID := uuid.Generate().String()
	user := entity.User{
		ID:       userID,
		Login:    in.Login,
		Password: in.Password,
	}
	err := gh.userService.CreateUser(ctx, user)
	if err != nil {
		gh.logger.Info("UserID:", userID, err.Error())
		return &TokenResponse{
			Status: err.Error(),
		}, nil
	}

	token, err := authorization.CreateToken(userID, gh.cfg.TokenConfig.AccessTimeLiveToken, gh.cfg.TokenConfig.AccessTimeLiveToken, gh.cfg.TokenConfig.SecretKey, gh.cfg.TokenConfig.SecretKey)
	if err != nil {
		return &TokenResponse{
			Status: err.Error(),
		}, nil
	}

	return &TokenResponse{
		Status:       "created",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// AuthUser функция авторизации пользователя
func (gh *GrpcUsers) AuthUser(ctx context.Context, in *AuthUserRequest) (*TokenResponse, error) {
	user := entity.User{
		Login:    in.Login,
		Password: in.Password,
	}

	userID, err := gh.userService.AuthUser(ctx, user)
	if err != nil {
		gh.logger.Info("UserID:", userID, err.Error())
		return &TokenResponse{
			Status: err.Error(),
		}, err
	}

	token, err := authorization.CreateToken(userID, gh.cfg.TokenConfig.AccessTimeLiveToken, gh.cfg.TokenConfig.AccessTimeLiveToken, gh.cfg.TokenConfig.SecretKey, gh.cfg.TokenConfig.SecretKey)
	if err != nil {
		gh.logger.Info("UserID:", userID, err.Error())
		return &TokenResponse{
			Status: err.Error(),
		}, nil
	}
	return &TokenResponse{
		Status:      "ok",
		AccessToken: token.AccessToken,
	}, nil
}

func (gh *GrpcUsers) DeleteUser(ctx context.Context, in *DeleteUserRequest) (*DeleteUserResponse, error) {
	err := gh.userService.DeleteUser(ctx, in.Login)
	if err != nil {
		gh.logger.Info("Delete login:", in.Login, err.Error())
		return &DeleteUserResponse{
			Status: err.Error(),
		}, err
	}

	return &DeleteUserResponse{
		Status: "ok",
	}, nil
}
