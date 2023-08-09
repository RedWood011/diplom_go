package secret

import adapter "RedWood011/client/adapter/secret"

func NewSecretService(adapter *adapter.SecretAdapter) *SecretService {
	return &SecretService{
		secretAdapter: adapter,
	}
}

// SecretHandler струкутра обработчика секретов
type SecretService struct {
	secretAdapter *adapter.SecretAdapter
}
