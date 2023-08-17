package secret

import adapter "RedWood011/client/adapter/secret"

func NewSecretService(adapter *adapter.Adapter) *Service {
	return &Service{
		secretAdapter: adapter,
	}
}

type Service struct {
	secretAdapter *adapter.Adapter
	AccessToken   string
}
