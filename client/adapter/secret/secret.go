package secret

import (
	"context"
	"fmt"

	"RedWood011/client/apperrors"
	"RedWood011/client/entity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	address string
}

type Client struct {
	SecretsClient
	closeFunc func() error
}

func NewSecretAdapter(address string) *Adapter {
	return &Adapter{
		address: address,
	}
}

func (sa *Adapter) GetSecret(ctx context.Context, token string, secretID string) (*entity.Secret, error) {
	client, err := sa.getConn()
	if err != nil {
		return nil, err
	}

	message := GetSecretRequest{
		SecretId: secretID,
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", token))
	response, err := client.GetSecret(ctx, &message)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated {
			return nil, apperrors.ErrAuth
		}
		return nil, err
	}

	var secret entity.Secret

	secret.ID = secretID
	secret.Data = response.Data
	secret.Name = response.Name

	client.closeFunc()

	return &secret, nil
}

func (sa *Adapter) CreateSecret(ctx context.Context, token string, secret *entity.Secret) (string, error) {
	client, err := sa.getConn()
	if err != nil {
		return "", err
	}
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %v", token)})
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	message := CreateSecretRequest{
		Name: secret.Name,
		Data: secret.Data,
	}
	response, err := client.CreateSecret(ctxNew, &message)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated {
			return "", apperrors.ErrAuth
		}
		return "", err
	}
	if response.Status != "created" {
		return "", fmt.Errorf("error create secret status: %v", response.Status)
	}

	client.closeFunc()
	return response.SecretId, nil
}

func (sa *Adapter) ListSecrets(ctx context.Context, token string) ([]entity.Secret, error) {
	client, err := sa.getConn()
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", token))
	message := ListSecretsRequest{}
	response, err := client.ListSecrets(ctx, &message)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated && statusErr.Message() == "invalid token" {
			return nil, apperrors.ErrAuth
		}
		return nil, err
	}
	secrets := make([]entity.Secret, 0, len(response.Data))
	for _, secret := range response.Data {
		secrets = append(secrets, entity.Secret{
			ID:   secret.SecretId,
			Name: secret.Name,
		})
	}

	client.closeFunc()
	return secrets, nil
}

func (sa *Adapter) DeleteSecret(ctx context.Context, token, secretID string) error {
	client, err := sa.getConn()
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", token))
	message := DeleteSecretRequest{SecretId: secretID}

	response, err := client.DeleteSecret(ctx, &message)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated && statusErr.Message() == "invalid token" {
			return apperrors.ErrAuth
		}
		return err
	}
	if response.Status != "ok" {
		return fmt.Errorf("error delete secret status: %v", response.Status)
	}

	client.closeFunc()
	return nil
}

func (sa *Adapter) getConn() (*Client, error) {
	conn, err := grpc.Dial(sa.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := NewSecretsClient(conn)

	return &Client{cl, conn.Close}, nil
}
