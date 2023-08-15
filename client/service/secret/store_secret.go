package secret

import (
	"context"

	"RedWood011/client/entity"
)

func (s *Service) CreateSecret(ctx context.Context, secret *entity.Secret) (string, error) {
	key := []byte{247, 43, 127, 3, 22, 127, 69, 35, 113, 231, 20, 127, 207, 9, 109, 70}
	nonce := []byte{200, 3, 38, 94, 66, 137, 119, 105, 204, 99, 7, 14}
	err := secret.EncryptSecret(key, nonce)
	if err != nil {
		return "", err
	}

	return s.secretAdapter.CreateSecret(ctx, s.AccessToken, secret)
}

func (s *Service) ListSecrets(ctx context.Context) ([]entity.Secret, error) {
	key := []byte{247, 43, 127, 3, 22, 127, 69, 35, 113, 231, 20, 127, 207, 9, 109, 70}
	nonce := []byte{200, 3, 38, 94, 66, 137, 119, 105, 204, 99, 7, 14}
	secrets, err := s.secretAdapter.ListSecrets(ctx, s.AccessToken)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(secrets); i++ {
		err = secrets[i].DecryptSecretName(key, nonce)
		if err != nil {
			return nil, err
		}
	}
	return secrets, nil
}

func (s *Service) DeleteSecret(ctx context.Context, secretID string) error {
	return s.secretAdapter.DeleteSecret(ctx, s.AccessToken, secretID)
}

func (s *Service) GetSecret(ctx context.Context, secretID string) (*entity.Secret, error) {
	key := []byte{247, 43, 127, 3, 22, 127, 69, 35, 113, 231, 20, 127, 207, 9, 109, 70}
	nonce := []byte{200, 3, 38, 94, 66, 137, 119, 105, 204, 99, 7, 14}
	secret, err := s.secretAdapter.GetSecret(ctx, s.AccessToken, secretID)
	if err != nil {
		return nil, err
	}

	err = secret.DecryptSecret(key, nonce)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
