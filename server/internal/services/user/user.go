package user

import (
	"context"

	"RedWood011/server/internal/entity"

	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=user.go -package=services	-destination=user_mock.go
type Repo interface {
	SaveUser(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, login string) error
}

type Service struct {
	db  Repo
	log *slog.Logger
}

func NewUserService(db Repo, log *slog.Logger) *Service {
	return &Service{
		db:  db,
		log: log,
	}
}
