package genres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
	"github.com/mkuptsov/movie-reviews/internal/dbx"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetGenres(ctx context.Context) ([]*Genre, error) {
	queryString := "SELECT id, name FROM genres;"
	rows, err := r.db.Query(ctx, queryString)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	var allGenres []*Genre
	for rows.Next() {
		var genre Genre
		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, apperrors.Internal(err)
		}

		allGenres = append(allGenres, &genre)
	}

	return allGenres, nil
}

func (r *Repository) GetGenreByID(ctx context.Context, id int) (*Genre, error) {
	queryString := "SELECT id, name FROM genres WHERE id = $1;"
	row := r.db.QueryRow(ctx, queryString, id)

	var genre Genre
	err := row.Scan(&genre.ID, &genre.Name)
	if dbx.IsNoRows(err) {
		return nil, apperrors.NotFound("genre", "id", id)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &genre, nil
}

func (r *Repository) CreateGenre(ctx context.Context, name string) (*Genre, error) {
	queryString := "INSERT INTO genres (name) VALUES ($1) returning id, name;"
	row := r.db.QueryRow(ctx, queryString, name)

	var genre Genre
	err := row.Scan(&genre.ID, &genre.Name)
	if dbx.IsUniqueViolation(err, "name") {
		return nil, apperrors.AlreadyExists("genre", "name", name)
	}
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &genre, nil
}

func (r *Repository) UpdateGenre(ctx context.Context, id int, name string) error {
	queryString := "UPDATE genres SET name = $2 WHERE id = $1;"
	cmdTag, err := r.db.Exec(ctx, queryString, id, name)
	if err != nil {
		return apperrors.Internal(err)
	}
	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", id)
	}

	return nil
}

func (r *Repository) DeleteGenre(ctx context.Context, id int) error {
	queryString := "DELETE FROM genres WHERE id = $1;"
	cmdTag, err := r.db.Exec(ctx, queryString, id)
	if err != nil {
		return apperrors.Internal(err)
	}
	if cmdTag.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", id)
	}

	return nil
}

func (r *Repository) GetGenresByMovieID(ctx context.Context, id int) ([]*Genre, error) {
	queryString := `
	SELECT g.id, g.name
	FROM genres g
	INNER JOIN movie_genres mg on mg.genre_id = g.id	
	WHERE mg.movie_id = $1
	ORDER BY mg.order_no
	`

	rows, err := r.db.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	genres, err := scanGenres(rows)
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, id int) ([]*MovieGenreRelation, error) {
	queryString := "SELECT movie_id, genre_id, order_no FROM movie_genres WHERE movie_id = $1"
	q := dbx.FromContext(ctx, r.db)
	rows, err := q.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var relations []*MovieGenreRelation
	for rows.Next() {
		var relation MovieGenreRelation
		err := rows.Scan(
			&relation.MovieID,
			&relation.GenreID,
			&relation.OrderNo,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		relations = append(relations, &relation)
	}

	return relations, nil
}

func scanGenres(rows pgx.Rows) ([]*Genre, error) {
	var genres []*Genre
	for rows.Next() {
		var genre Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Name,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		genres = append(genres, &genre)
	}
	return genres, nil
}
