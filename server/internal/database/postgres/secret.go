package postgres

import (
	"context"

	"RedWood011/server/internal/apperrors"
	"RedWood011/server/internal/entity"
	"github.com/jackc/pgx/v4"
)

type Secret struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Name   []byte `db:"secret_name"`
	Data   []byte `db:"secret_data"`
}

func (r *Repository) SaveSecret(ctx context.Context, secret *entity.Secret) error {
	const query = `INSERT INTO secrets (id,user_id, secret_data, secret_name) VALUES ($1, $2,$3, $4)`

	_, err := r.db.Exec(ctx, query, secret.ID, secret.UserID, secret.Data, secret.Name)

	return err
}

func (r *Repository) ListSecrets(ctx context.Context, userID string) ([]entity.Secret, error) {
	const query = `SELECT id, user_id,secret_name FROM secrets WHERE user_id = $1 AND deleted_at IS NULL`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var result []entity.Secret
	for rows.Next() {
		secret := Secret{}
		err = rows.Scan(&secret.ID, &secret.UserID, &secret.Name)
		if err != nil {
			return nil, err
		}

		result = append(result, entity.Secret{ID: secret.ID, UserID: secret.UserID, Name: secret.Name})
	}

	return result, err
}

func (r *Repository) GetSecret(ctx context.Context, secret entity.Secret) (entity.Secret, error) {
	const query = `SELECT id, user_id,secret_data FROM secrets WHERE id = $1 AND user_id=$2 AND deleted_at IS NULL`

	var s entity.Secret
	res := r.db.QueryRow(ctx, query, secret.ID, secret.UserID)
	err := res.Scan(&secret.ID, &secret.UserID, &secret.Data)
	if err == pgx.ErrNoRows {
		return entity.Secret{}, apperrors.ErrNotFound
	}
	if err != nil {
		return entity.Secret{}, err
	}

	s = entity.Secret{ID: secret.ID, UserID: secret.UserID, Name: secret.Name, Data: secret.Data}
	return s, err
}

func (r *Repository) DeleteSecret(ctx context.Context, secretID string, userID string) error {
	const query = `UPDATE secrets SET deleted_at = current_timestamp WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, secretID, userID)

	return err
}
