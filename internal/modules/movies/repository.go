package movies

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

func (r *Repository) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	queryString := `
	INSERT INTO movies 
	(title, release_date, description) 
	VALUES 
	($1, $2, $3)
	RETURNING
	id, created_at;
	`
	row := r.db.QueryRow(ctx, queryString,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Description,
	)

	err := row.Scan(
		&movie.ID,
		&movie.CreatedAt,
	)
	if err != nil {
		return apperrors.Internal(err)
	}
	return nil
}

func (r *Repository) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	queryString := `
	SELECT id, title, description, release_date, created_at, deleted_at, version
	FROM movies
	WHERE id = $1 and deleted_at IS NULL;`

	movie := MovieDetails{}

	row := r.db.QueryRow(ctx, queryString, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.CreatedAt,
		&movie.DeletedAt,
		&movie.Version,
	)

	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("movie", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &movie, nil
}

func (r *Repository) GetAllPaginated(ctx context.Context, offset, limit int) ([]*Movie, int, error) {
	queryPage := `
	SELECT id, title, release_date, created_at, deleted_at
	FROM movies
	WHERE deleted_at IS NULL
	ORDER BY id
	LIMIT $2
	OFFSET $1;`

	queryTotal := "SELECT count(*) FROM movies WHERE deleted_at IS NULL"

	b := &pgx.Batch{}
	b.Queue(queryPage, offset, limit)
	b.Queue(queryTotal)

	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.CreatedAt,
			&movie.DeletedAt,
		)
		if err != nil {
			return nil, 0, apperrors.Internal(err)
		}

		movies = append(movies, &movie)
	}

	var total int
	err = br.QueryRow().Scan(&total)
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return movies, total, nil
}

func (r *Repository) UpdateMovie(ctx context.Context, id int, movie *MovieDetails) error {
	queryString := `
	UPDATE movies 
	SET 
		title = $2,
		release_date = $3,
		description = $4,
		version = version + 1
	WHERE id = $1 and deleted_at IS NULL and version = $5`

	cmdTag, err := r.db.Exec(ctx, queryString,
		id,
		movie.Title,
		movie.ReleaseDate,
		movie.Description,
		movie.Version,
	)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		_, err := r.GetMovieByID(ctx, id)
		if err != nil {
			return err
		}
		return apperrors.VersionMismatch("movie", "id", id, movie.Version)
	}

	return nil
}

func (r *Repository) DeleteMovie(ctx context.Context, id int) error {
	queryString := "UPDATE movies SET deleted_at = NOW() WHERE id = $1 and deleted_at IS NULL;"
	cmdTag, err := r.db.Exec(ctx, queryString, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("movie", "id", id)
	}

	return nil
}
