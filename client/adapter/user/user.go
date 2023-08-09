package user

import (
	"context"
	"errors"

	"RedWood011/client/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserAdapter(address string) *UserAdapter {
	return &UserAdapter{
		address: address,
	}
}

type UserClient struct {
	UsersClient
	closeFunc func() error
}

type UserAdapter struct {
	address string
}

// Register функция регистрации пользователя
func (u *UserAdapter) RegisterUser(ctx context.Context, user entity.User) (string, string, error) {
	client, err := u.getConn()
	if err != nil {
		return "", "", err
	}
	message := CreateUserRequest{
		Login:    user.Login,
		Password: user.Password,
	}
	response, err := client.CreateUser(ctx, &message)
	if err != nil {
		return "", "", err
	}

	if response.Status == "created" {
		return response.AccessToken, response.RefreshToken, nil
	}

	return "", "", errors.New(response.Status)
}

// Auth функция авторизации пользователя
func (u *UserAdapter) AuthUser(ctx context.Context, user entity.User) (string, string, error) {
	client, err := u.getConn()
	if err != nil {
		return "", "", err
	}
	message := AuthUserRequest{
		Login:    user.Login,
		Password: user.Password,
	}
	response, err := client.AuthUser(ctx, &message)
	if err != nil {
		return "", "", err
	}
	if response.Status == "ok" {
		return response.AccessToken, response.RefreshToken, nil
	}
	return "", "", errors.New(response.Status)
}

func (u *UserAdapter) getConn() (*UserClient, error) {
	conn, err := grpc.Dial(u.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := NewUsersClient(conn)

	return &UserClient{cl, conn.Close}, nil
}
