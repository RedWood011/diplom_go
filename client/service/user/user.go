package user

import adapter "RedWood011/client/adapter/user"

func NewUserService(adapter *adapter.UserAdapter) *UserServise {
	return &UserServise{
		userAdapter: adapter,
	}
}

type UserServise struct {
	userAdapter *adapter.UserAdapter
}
