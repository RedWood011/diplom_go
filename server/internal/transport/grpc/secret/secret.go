package secret

import (
	"context"

	"RedWood011/server/internal/config"
	"RedWood011/server/internal/entity"
	"github.com/docker/distribution/uuid"
	"golang.org/x/exp/slog"
)

type SecretService interface {
	CreateSecret(ctx context.Context, secret *entity.Secret) error
	ListSecrets(ctx context.Context, userID string) ([]entity.Secret, error)
	DeleteSecrets(ctx context.Context, secretID string, userID string) error
	GetSecret(ctx context.Context, secret entity.Secret) (entity.Secret, error)
}

type GrpcSecrets struct {
	UnimplementedSecretsServer
	secretService SecretService
	cfg           *config.Config
	logger        *slog.Logger
}

func NewGrpcSecrets(secret SecretService, cfg *config.Config, log *slog.Logger) *GrpcSecrets {
	return &GrpcSecrets{
		secretService: secret,
		cfg:           cfg,
		logger:        log,
	}
}

func (g *GrpcSecrets) CreateSecret(ctx context.Context, in *CreateSecretRequest) (*CreateSecretResponse, error) {
	userID := ctx.Value("userID").(string)
	secret := entity.Secret{
		ID:     uuid.Generate().String(),
		UserID: userID,
		Data:   in.Data,
		Name:   in.Name,
	}
	err := g.secretService.CreateSecret(ctx, &secret)
	if err != nil {
		g.logger.Info("Error create Secret UserID:", secret.UserID, err.Error())
		return &CreateSecretResponse{
			Status: err.Error(),
		}, nil
	}

	return &CreateSecretResponse{
		Status:   "created",
		SecretId: secret.ID,
	}, nil
}

func (g *GrpcSecrets) ListSecrets(ctx context.Context, in *ListSecretsRequest) (*ListSecretsResponse, error) {
	userID := ctx.Value("userID").(string)
	secrets, err := g.secretService.ListSecrets(ctx, userID)
	if err != nil {
		g.logger.Info("Error list Secrets UserID:", userID, err.Error())
		return &ListSecretsResponse{
			Status: err.Error(),
		}, nil
	}

	data := make([]*ListData, 0, len(secrets))
	for _, secret := range secrets {
		data = append(data, &ListData{
			SecretId: secret.ID,
			Name:     secret.Name,
		})
	}

	return &ListSecretsResponse{
		Status: "ok",
		Data:   data,
	}, nil
}

func (g *GrpcSecrets) DeleteSecret(ctx context.Context, in *DeleteSecretRequest) (*DeleteSecretResponse, error) {
	userID := ctx.Value("userID").(string)
	err := g.secretService.DeleteSecrets(ctx, in.SecretId, userID)
	if err != nil {
		g.logger.Info("Error delete Secrets UserID:", userID, err.Error())
		return &DeleteSecretResponse{
			Status: err.Error(),
		}, nil

	}
	return &DeleteSecretResponse{
		Status: "ok",
	}, nil
}

func (g *GrpcSecrets) GetSecret(ctx context.Context, in *GetSecretRequest) (*GetSecretResponse, error) {
	userID := ctx.Value("userID").(string)
	secret, err := g.secretService.GetSecret(ctx, entity.Secret{
		ID:     in.SecretId,
		UserID: userID,
	})
	if err != nil {
		g.logger.Info("Error get Secret UserID:", userID, err.Error())
		return &GetSecretResponse{
			Status: err.Error(),
		}, nil
	}

	return &GetSecretResponse{
		Status: "ok",
		Name:   secret.Name,
		Data:   secret.Data,
	}, nil
}
