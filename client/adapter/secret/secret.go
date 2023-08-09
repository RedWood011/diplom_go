package secret

import (
	"context"
	"fmt"

	"RedWood011/client/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SecretAdapter struct {
	address     string
	accessToken string
}

type SecretClient struct {
	SecretsClient
	closeFunc func() error
}

func NewSecretAdapter(address string, access string) *SecretAdapter {
	return &SecretAdapter{
		address:     address,
		accessToken: access,
	}
}

func (c *SecretAdapter) GetSecret(ctx context.Context, secretID string) (entity.Secret, error) {
	client, err := c.getConn()
	if err != nil {
		return entity.Secret{}, err
	}

	message := GetSecretRequest{
		SecretId: secretID,
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	response, err := client.GetSecret(ctx, &message)
	if err != nil {
		return entity.Secret{}, err
	}
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated && statusErr.Message() == "invalid token" {
			return entity.Secret{}, fmt.Errorf("please login again\n")
		}
		return entity.Secret{}, err
	}

	var secret entity.Secret

	secret.ID = secretID
	secret.Data = response.Data
	secret.Name = response.Name

	client.closeFunc()

	return secret, nil
}

func (c *SecretAdapter) CreateSecret(ctx context.Context, secret *entity.Secret) (string, error) {
	client, err := c.getConn()
	if err != nil {
		return "", err
	}
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %v", c.accessToken)})
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	message := CreateSecretRequest{
		Name: secret.Name,
		Data: secret.Data,
	}
	response, err := client.CreateSecret(ctxNew, &message)
	if err != nil {
		return "", err
	}
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated && statusErr.Message() == "invalid token" {
			return "", fmt.Errorf("please login again\n")
		}
		return "", err
	}
	if response.Status != "created" {
		return "", fmt.Errorf("error create secret status: %v", response.Status)
	}

	client.closeFunc()
	return response.SecretId, nil
}

func (c *SecretAdapter) ListSecrets(ctx context.Context) ([]entity.Secret, error) {
	client, err := c.getConn()
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	message := ListSecretsRequest{}
	response, err := client.ListSecrets(ctx, &message)
	if err != nil {
		return nil, err
	}
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated && statusErr.Message() == "invalid token" {
			return nil, fmt.Errorf("please login again\n")
		}
		return nil, err
	}
	var secrets []entity.Secret
	for _, secret := range response.Data {
		secrets = append(secrets, entity.Secret{
			ID:   secret.SecretId,
			Name: secret.Name,
		})
	}

	client.closeFunc()
	return secrets, nil
}

func (c *SecretAdapter) DeleteSecret(ctx context.Context, secretID string) error {
	client, err := c.getConn()
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", c.accessToken))
	message := DeleteSecretRequest{SecretId: secretID}

	response, err := client.DeleteSecret(ctx, &message)
	if err != nil {
		return err
	}
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.Unauthenticated && statusErr.Message() == "invalid token" {
			return fmt.Errorf("please login again\n")
		}
		return err
	}
	if response.Status != "ok" {
		return fmt.Errorf("error delete secret status: %v", response.Status)
	}

	client.closeFunc()
	return nil
}

func (sa *SecretAdapter) getConn() (*SecretClient, error) {
	conn, err := grpc.Dial(sa.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cl := NewSecretsClient(conn)

	return &SecretClient{cl, conn.Close}, nil
}
