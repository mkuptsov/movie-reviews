package movies

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"
	"github.com/cloudmachinery/movie-reviews/internal/slices"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db               *pgxpool.Pool
	genresRepository *genres.Repository
}

func NewRepository(db *pgxpool.Pool, genresRepository *genres.Repository) *Repository {
	return &Repository{
		db:               db,
		genresRepository: genresRepository,
	}
}

func (r *Repository) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queryString := `
	INSERT INTO movies 
	(title, release_date, description) 
	VALUES 
	($1, $2, $3)
	RETURNING
	id, created_at;
	`
		row := tx.QueryRow(ctx, queryString,
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

		next := slices.MapIndex(movie.Genres, func(i int, genre *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: movie.ID,
				GenreID: genre.ID,
				OrderNo: i,
			}
		})

		return r.updateGenres(ctx, []*genres.MovieGenreRelation{}, next)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
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
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queryString := `
	UPDATE movies 
	SET 
		title = $2,
		release_date = $3,
		description = $4,
		version = version + 1
	WHERE id = $1 and deleted_at IS NULL and version = $5`

		cmdTag, err := tx.Exec(ctx, queryString,
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
			_, err = r.GetMovieByID(ctx, id)
			if err != nil {
				return err
			}
			return apperrors.VersionMismatch("movie", "id", id, movie.Version)
		}

		next := slices.MapIndex(movie.Genres, func(i int, genre *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: id,
				GenreID: genre.ID,
				OrderNo: i,
			}
		})

		current, err := r.genresRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}

		return r.updateGenres(ctx, current, next)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) DeleteMovie(ctx context.Context, id int) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queryString := "UPDATE movies SET deleted_at = NOW() WHERE id = $1 and deleted_at IS NULL;"
		cmdTag, err := tx.Exec(ctx, queryString, id)
		if err != nil {
			return apperrors.Internal(err)
		}

		if cmdTag.RowsAffected() == 0 {
			return apperrors.NotFound("movie", "id", id)
		}
		current, err := r.genresRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}
		return r.updateGenres(ctx, current, []*genres.MovieGenreRelation{})
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) updateGenres(ctx context.Context, current, next []*genres.MovieGenreRelation) error {
	q := dbx.FromContext(ctx, r.db)
	addFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx,
			"INSERT INTO movie_genres (movie_id, genre_id, order_no) VALUES ($1, $2, $3)",
			mgo.MovieID, mgo.GenreID, mgo.OrderNo)
		return err
	}

	removeFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx,
			"DELETE FROM movie_genres WHERE movie_id = $1 and genre_id = $2",
			mgo.MovieID, mgo.GenreID)
		return err
	}

	return dbx.AdjustRelations(current, next, addFunc, removeFunc)
}
