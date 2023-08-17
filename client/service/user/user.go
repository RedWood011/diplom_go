package user

import adapter "RedWood011/client/adapter/user"

func NewUserService(adapter *adapter.Adapter) *Service {
	return &Service{
		userAdapter: adapter,
	}
}

type Service struct {
	userAdapter *adapter.Adapter
}
