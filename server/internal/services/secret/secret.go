package secret

import (
	"context"

	"RedWood011/server/internal/entity"
	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=secret.go -package=secret	-destination=secret_mock.go
type SecretRepo interface {
	SaveSecret(ctx context.Context, secret *entity.Secret) error
	ListSecrets(ctx context.Context, userID string) ([]entity.Secret, error)
	GetSecret(ctx context.Context, secret entity.Secret) (entity.Secret, error)
	DeleteSecret(ctx context.Context, secretID string, userID string) error
}

type SecretService struct {
	db  SecretRepo
	log *slog.Logger
}

func NewSecretService(db SecretRepo, log *slog.Logger) *SecretService {
	return &SecretService{
		db:  db,
		log: log,
	}
}
