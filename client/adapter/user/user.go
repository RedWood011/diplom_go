package user

import (
	"context"
	"errors"

	"RedWood011/client/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserAdapter(address string) *Adapter {
	return &Adapter{
		address: address,
	}
}

type Client struct {
	UsersClient
	closeFunc func() error
}

type Adapter struct {
	address string
}

func (u *Adapter) RegisterUser(ctx context.Context, user entity.User) (string, string, error) {
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

func (u *Adapter) AuthUser(ctx context.Context, user entity.User) (string, string, error) {
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

func (u *Adapter) getConn() (*Client, error) {
	conn, err := grpc.Dial(u.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := NewUsersClient(conn)

	return &Client{cl, conn.Close}, nil
}
