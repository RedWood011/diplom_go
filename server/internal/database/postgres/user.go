package postgres

import (
	"context"

	"RedWood011/server/internal/apperrors"
	"RedWood011/server/internal/entity"
	"github.com/jackc/pgx/v4"
)

type User struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (r *Repository) GetUser(ctx context.Context, user entity.User) (entity.User, error) {
	var res User
	const sqlCheckUser = `SELECT id,login,password FROM users WHERE login = $1;`

	query := r.db.QueryRow(ctx, sqlCheckUser, user.Login)
	err := query.Scan(&res.ID, &res.Login, &res.Password)
	if err == pgx.ErrNoRows {
		return entity.User{}, apperrors.ErrNotFound
	}

	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		ID:       res.ID,
		Login:    res.Login,
		Password: res.Password,
	}, nil
}

func (r *Repository) SaveUser(ctx context.Context, user entity.User) error {
	const query = `INSERT INTO users (id,login, password) VALUES ($1,$2,$3)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Login, user.Password)
	return err
}

// DeleteUser функция удаления пользователя
func (r *Repository) DeleteUser(ctx context.Context, login string) error {
	const query = `UPDATE users SET  deleted_at = current_timestamp WHERE login = $1`
	_, err := r.db.Exec(ctx, query, login)
	if err != nil {
		return apperrors.ErrUserNoDeleted
	}
	return nil
}
