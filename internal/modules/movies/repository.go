package movies

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/apperrors"
	"github.com/cloudmachinery/movie-reviews/internal/dbx"
	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"
	"github.com/cloudmachinery/movie-reviews/internal/modules/stars"
	"github.com/cloudmachinery/movie-reviews/internal/slices"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db               *pgxpool.Pool
	genresRepository *genres.Repository
	starsRepository  *stars.Repository
}

func NewRepository(db *pgxpool.Pool, genresRepo *genres.Repository, starsRepo *stars.Repository) *Repository {
	return &Repository{
		db:               db,
		genresRepository: genresRepo,
		starsRepository:  starsRepo,
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

		nextGenres := slices.MapIndex(movie.Genres, func(i int, genre *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: movie.ID,
				GenreID: genre.ID,
				OrderNo: i,
			}
		})

		err = r.updateGenres(ctx, []*genres.MovieGenreRelation{}, nextGenres)
		if err != nil {
			return err
		}

		nextCast := slices.MapIndex(movie.Cast, func(i int, cast *stars.MovieCredit) *stars.MovieStarRelation {
			return &stars.MovieStarRelation{
				MovieID: movie.ID,
				StarID:  cast.Star.ID,
				Role:    cast.Role,
				Details: cast.Details,
				OrderNo: i,
			}
		})

		return r.updateCast(ctx, []*stars.MovieStarRelation{}, nextCast)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}
	return nil
}

func (r *Repository) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	q := dbx.FromContext(ctx, r.db)
	queryString := `
	SELECT id, title, description, release_date, created_at, deleted_at, version
	FROM movies
	WHERE id = $1 and deleted_at IS NULL;`

	movie := MovieDetails{}

	row := q.QueryRow(ctx, queryString, id)
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

func (r *Repository) GetAllPaginated(ctx context.Context, starID *int, searchTerm *string, offset, limit int) ([]*Movie, int, error) {
	queryPage := dbx.StatementBuilder.
		Select("id, title, release_date, created_at, deleted_at").
		From("movies").
		Where("deleted_at IS NULL").
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	queryTotal := dbx.StatementBuilder.
		Select("count(*)").
		From("movies").
		Where("deleted_at IS NULL")

	if starID != nil {
		queryPage = queryPage.
			Join("movie_stars on movies.id = movie_stars.movie_id").
			Where("star_id = ?", starID)

		queryTotal = queryTotal.
			Join("movie_stars on movies.id = movie_stars.movie_id").
			Where("star_id = ?", starID)
	}

	if searchTerm != nil {
		queryPage = queryPage.
			Where("search_vector @@ to_tsquery('english', ?)", *searchTerm).
			OrderByClause("ts_rank_cd(search_vector, to_tsquery('english', ?)) DESC", *searchTerm)

		queryTotal = queryTotal.
			Where("search_vector @@ to_tsquery('english', ?)", *searchTerm)
	}

	b := &pgx.Batch{}

	err := dbx.QueueBatchSelect(b, queryPage)
	if err != nil {
		return nil, 0, err
	}
	err = dbx.QueueBatchSelect(b, queryTotal)
	if err != nil {
		return nil, 0, err
	}

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

		nextGenres := slices.MapIndex(movie.Genres, func(i int, genre *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: id,
				GenreID: genre.ID,
				OrderNo: i,
			}
		})

		currentGenres, err := r.genresRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}

		err = r.updateGenres(ctx, currentGenres, nextGenres)
		if err != nil {
			return err
		}

		nextCast := slices.MapIndex(movie.Cast, func(i int, mc *stars.MovieCredit) *stars.MovieStarRelation {
			return &stars.MovieStarRelation{
				MovieID: id,
				StarID:  mc.Star.ID,
				Role:    mc.Role,
				Details: mc.Details,
				OrderNo: i,
			}
		})

		currentCast, err := r.starsRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}

		return r.updateCast(ctx, currentCast, nextCast)
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
		currentGenres, err := r.genresRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}
		err = r.updateGenres(ctx, currentGenres, []*genres.MovieGenreRelation{})
		if err != nil {
			return err
		}
		currentCast, err := r.starsRepository.GetRelationsByMovieID(ctx, id)
		if err != nil {
			return err
		}
		return r.updateCast(ctx, currentCast, []*stars.MovieStarRelation{})
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
		if err != nil {
			return apperrors.Internal(err)
		}
		return nil
	}

	removeFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx,
			"DELETE FROM movie_genres WHERE movie_id = $1 and genre_id = $2",
			mgo.MovieID, mgo.GenreID)
		if err != nil {
			return apperrors.Internal(err)
		}
		return nil
	}

	return dbx.AdjustRelations(current, next, addFunc, removeFunc)
}

func (r *Repository) updateCast(ctx context.Context, current, next []*stars.MovieStarRelation) error {
	q := dbx.FromContext(ctx, r.db)
	addFunc := func(s *stars.MovieStarRelation) error {
		_, err := q.Exec(ctx,
			"INSERT INTO movie_stars (movie_id, star_id, role, details, order_no) VALUES ($1, $2, $3, $4, $5)",
			s.MovieID, s.StarID, s.Role, s.Details, s.OrderNo)
		if err != nil {
			return apperrors.Internal(err)
		}
		return nil
	}
	removeFunc := func(s *stars.MovieStarRelation) error {
		_, err := q.Exec(ctx,
			"DELETE FROM movie_stars WHERE movie_id = $1 and star_id = $2 and role = $3",
			s.MovieID, s.StarID, s.Role)
		if err != nil {
			return apperrors.Internal(err)
		}
		return nil
	}

	return dbx.AdjustRelations(current, next, addFunc, removeFunc)
}
