package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user *UserWithPassword) error {
	queryString := "INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) returning id, created_at, role"
	err := r.db.QueryRow(ctx, queryString, user.Username, user.Email, user.PasswordHash).Scan(&user.Id, &user.CreatedAt, &user.Role)

	return err
}

func (r *Repository) GetUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	queryString := `
	SELECT id, username, email, pass_hash, role, created_at, deleted_at, bio
	FROM users
	WHERE email = $1 and deleted_at IS NULL;`

	user := UserWithPassword{
		User:         &User{},
		PasswordHash: "",
	}

	row := r.db.QueryRow(ctx, queryString, email)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if err != nil {
		return nil, fmt.Errorf("repo scan: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id int) (*User, error) {
	queryString := `
	SELECT id, username, email, role, created_at, deleted_at, bio
	FROM users
	WHERE id = $1 and deleted_at IS NULL;`

	user := &User{}

	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if err != nil {
		return nil, fmt.Errorf("repo scan: %w", err)
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	queryString := "UPDATE users SET deleted_at = NOW() WHERE id = $1 and deleted_at IS NULL;"
	cmdTag, err := r.db.Exec(ctx, queryString, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id int, bio string) error {
	queryString := "UPDATE users SET bio = $2 WHERE id = $1 and deleted_at IS NULL;"
	cmdTag, err := r.db.Exec(ctx, queryString, id, bio)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
