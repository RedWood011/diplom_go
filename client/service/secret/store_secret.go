package secret

import (
	"context"

	"RedWood011/client/entity"
)

var (
	key   = []byte{247, 43, 127, 3, 22, 127, 69, 35, 113, 231, 20, 127, 207, 9, 109, 70}
	nonce = []byte{200, 3, 38, 94, 66, 137, 119, 105, 204, 99, 7, 14}
)

func (s *SecretService) CreateSecret(ctx context.Context, secret *entity.Secret) (string, error) {
	err := secret.EncryptSecret(key, nonce)
	if err != nil {
		return "", err
	}

	return s.secretAdapter.CreateSecret(ctx, secret)
}
