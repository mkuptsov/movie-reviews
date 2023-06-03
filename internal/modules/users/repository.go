package users

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/modules/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/modules/dbx"
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
	queryString := "INSERT INTO users (username, email, pass_hash, role) VALUES ($1, $2, $3, $4) returning id, created_at, role"
	err := r.db.QueryRow(ctx, queryString, user.Username, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt, &user.Role)

	if dbx.IsUniqueViolation(err, "email") {
		return apperrors.AlreadyExists("user", "email", user.Email)
	}
	if dbx.IsUniqueViolation(err, "username") {
		return apperrors.AlreadyExists("user", "username", user.Username)
	}
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	queryString := `
	SELECT id, username, email, pass_hash, role, created_at, deleted_at, bio
	FROM users
	WHERE email = $1 and deleted_at IS NULL;`

	user := UserWithPassword{
		User: &User{},
	}

	row := r.db.QueryRow(ctx, queryString, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("user", "email", email)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id int) (*User, error) {
	queryString := `
	SELECT id, username, email, role, created_at, deleted_at, bio
	FROM users
	WHERE id = $1 and deleted_at IS NULL;`

	user := User{}

	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("user", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &user, nil
}

func (r *Repository) GetUserByUserName(ctx context.Context, userName string) (*User, error) {
	queryString := `
	SELECT id, username, email, role, created_at, deleted_at, bio
	FROM users
	WHERE username = $1 and deleted_at IS NULL;`

	user := User{}

	row := r.db.QueryRow(ctx, queryString, userName)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.Bio,
	)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("user", "username", userName)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	queryString := "UPDATE users SET deleted_at = NOW() WHERE id = $1 and deleted_at IS NULL;"
	cmdTag, err := r.db.Exec(ctx, queryString, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id int, bio string) error {
	queryString := "UPDATE users SET bio = $2 WHERE id = $1 and deleted_at IS NULL;"
	cmdTag, err := r.db.Exec(ctx, queryString, id, bio)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}
	return nil
}

func (r *Repository) UpdateUserRole(ctx context.Context, id int, roleName string) error {
	queryString := "UPDATE users SET role = $2 WHERE id = $1"
	cmdTag, err := r.db.Exec(ctx, queryString, id, roleName)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}
	return nil
}
