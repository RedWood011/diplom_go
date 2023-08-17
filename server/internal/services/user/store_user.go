package user

import (
	"context"

	"RedWood011/server/internal/apperrors"
	"RedWood011/server/internal/entity"
)

func (us *Service) CreateUser(ctx context.Context, user entity.User) error {
	if !user.IsValidPassword() || !user.IsValidLogin() {
		return apperrors.ErrAuth
	}

	err := user.SaveHashPassword()
	if err != nil {
		return err
	}

	return us.db.SaveUser(ctx, user)
}

func (us *Service) AuthUser(ctx context.Context, user entity.User) (string, error) {
	if !user.IsValidPassword() || !user.IsValidLogin() {
		return "", apperrors.ErrAuth
	}

	existUser, err := us.db.GetUser(ctx, user)
	if err != nil {
		return "", err
	}

	if !existUser.IsEqual(user) {
		return "", apperrors.ErrAuth
	}

	return existUser.ID, nil
}

func (us *Service) DeleteUser(ctx context.Context, login string) error {
	return us.db.DeleteUser(ctx, login)
}
