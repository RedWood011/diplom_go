package user

import (
	"context"

	"RedWood011/client/entity"
)

func (us *Service) RegisterUser(ctx context.Context, user entity.User) (string, string, error) {
	return us.userAdapter.RegisterUser(ctx, user)
}

func (us *Service) AuthUser(ctx context.Context, user entity.User) (string, string, error) {
	return us.userAdapter.AuthUser(ctx, user)
}
