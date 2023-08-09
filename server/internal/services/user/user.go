package user

import (
	"context"

	"RedWood011/server/internal/entity"
	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=user.go -package=services	-destination=user_mock.go
type UserRepo interface {
	SaveUser(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, login string) error
}

type UserService struct {
	db  UserRepo
	log *slog.Logger
}

func NewUserService(db UserRepo, log *slog.Logger) *UserService {
	return &UserService{
		db:  db,
		log: log,
	}
}
