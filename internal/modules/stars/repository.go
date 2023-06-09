package stars

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
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

func (r *Repository) CreateStar(ctx context.Context, star *Star) error {
	queryString := `
	INSERT INTO stars 
	(first_name, middle_name, last_name, birth_date, birth_place, death_date, bio) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING
	id, created_at, deleted_at
	`
	row := r.db.QueryRow(ctx, queryString,
		star.FirstName,
		star.MiddleName,
		star.LastName,
		star.BirthDate,
		star.BirthPlace,
		star.DeathDate,
		star.Bio,
	)

	err := row.Scan(
		&star.ID,
		&star.CreatedAt,
		&star.DeletedAt,
	)
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetStarByID(ctx context.Context, id int) (*Star, error) {
	queryString := `
	SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at, deleted_at
	FROM stars
	WHERE id = $1 and deleted_at IS NULL;`

	star := Star{}

	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(
		&star.ID,
		&star.FirstName,
		&star.MiddleName,
		&star.LastName,
		&star.BirthDate,
		&star.BirthPlace,
		&star.DeathDate,
		&star.Bio,
		&star.CreatedAt,
		&star.DeletedAt,
	)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("star", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &star, nil
}

func (r *Repository) GetAllPaginated(ctx context.Context, offset, limit int) ([]*Star, int, error) {
	queryPage := `
	SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at, deleted_at
	FROM stars
	WHERE deleted_at IS NULL
	ORDER BY id
	LIMIT $2
	OFFSET $1;`

	queryTotal := "SELECT count(*) FROM stars WHERE deleted_at IS NULL"

	b := &pgx.Batch{}
	b.Queue(queryPage, offset, limit)
	b.Queue(queryTotal)

	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var stars []*Star
	for rows.Next() {
		var star Star
		err = rows.Scan(
			&star.ID,
			&star.FirstName,
			&star.MiddleName,
			&star.LastName,
			&star.BirthDate,
			&star.BirthPlace,
			&star.DeathDate,
			&star.Bio,
			&star.CreatedAt,
			&star.DeletedAt,
		)
		if err != nil {
			return nil, 0, apperrors.Internal(err)
		}

		stars = append(stars, &star)
	}

	var total int
	err = br.QueryRow().Scan(&total)
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return stars, total, nil
}

func (r *Repository) UpdateStar(ctx context.Context, id int, star *Star) error {
	queryString := `
	UPDATE stars 
	SET 
		first_name = $2,
		middle_name = $3,
		last_name = $4,
		birth_date = $5,
		birth_place = $6,
		death_date = $7,
		bio = $8
		
	WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, queryString,
		id,
		star.FirstName,
		star.MiddleName,
		star.LastName,
		star.BirthDate,
		star.BirthPlace,
		star.DeathDate,
		star.Bio,
	)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("star", "id", id)
	}

	return nil
}

func (r *Repository) DeleteStar(ctx context.Context, id int) error {
	queryString := "UPDATE stars SET deleted_at = NOW() WHERE id = $1 and deleted_at IS NULL;"
	cmdTag, err := r.db.Exec(ctx, queryString, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("star", "id", id)
	}

	return nil
}
