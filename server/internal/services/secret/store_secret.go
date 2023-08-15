package secret

import (
	"context"

	"RedWood011/server/internal/entity"
)

func (sc *Service) CreateSecret(ctx context.Context, secret *entity.Secret) error {
	return sc.db.SaveSecret(ctx, secret)
}

func (sc *Service) ListSecrets(ctx context.Context, userID string) ([]entity.Secret, error) {
	return sc.db.ListSecrets(ctx, userID)
}

func (sc *Service) DeleteSecrets(ctx context.Context, secretID string, userID string) error {
	return sc.db.DeleteSecret(ctx, secretID, userID)
}

func (sc *Service) GetSecret(ctx context.Context, secret entity.Secret) (entity.Secret, error) {
	return sc.db.GetSecret(ctx, secret)
}
