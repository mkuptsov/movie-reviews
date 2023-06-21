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

func (r *Repository) CreateStar(ctx context.Context, star *StarDetails) error {
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

func (r *Repository) GetStarByID(ctx context.Context, id int) (*StarDetails, error) {
	queryString := `
	SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at, deleted_at
	FROM stars
	WHERE id = $1 and deleted_at IS NULL;`

	star := StarDetails{}

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

func (r *Repository) GetAllPaginated(ctx context.Context, movieID *int, offset, limit int) ([]*Star, int, error) {
	queryPage := dbx.StatementBuilder.
		Select("id, first_name, last_name, birth_date, death_date, created_at, deleted_at").
		From("stars").
		Where("deleted_at IS NULL").
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	queryTotal := dbx.StatementBuilder.
		Select("count(*)").
		From("stars").
		Where("deleted_at IS NULL")

	if movieID != nil {
		queryPage = queryPage.
			Join("movie_stars on stars.id = movie_stars.star_id").
			Where("movie_stars.movie_id = ?", movieID)

		queryTotal = queryTotal.
			Join("movie_stars on stars.id = movie_stars.star_id").
			Where("movie_stars.movie_id = ?", movieID)
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

	var stars []*Star
	for rows.Next() {
		var star Star
		err = rows.Scan(
			&star.ID,
			&star.FirstName,
			&star.LastName,
			&star.BirthDate,
			&star.DeathDate,
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

func (r *Repository) UpdateStar(ctx context.Context, id int, star *StarDetails) error {
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

func (r *Repository) GetCastsByMovieID(ctx context.Context, id int) ([]*MovieCredit, error) {
	queryString := `
	SELECT s.id, s.first_name, s.last_name, s.birth_date, s.death_date, s.created_at, s.deleted_at, ms.role, ms.details
	FROM stars s
	INNER JOIN movie_stars ms ON ms.star_id = s.id
	WHERE ms.movie_id = $1
	ORDER BY ms.order_no`
	rows, err := r.db.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	var cast []*MovieCredit
	for rows.Next() {
		var mc MovieCredit
		err = rows.Scan(
			&mc.Star.ID,
			&mc.Star.FirstName,
			&mc.Star.LastName,
			&mc.Star.BirthDate,
			&mc.Star.DeathDate,
			&mc.Star.CreatedAt,
			&mc.Star.DeletedAt,
			&mc.Role,
			&mc.Details,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		cast = append(cast, &mc)
	}
	return cast, nil
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, id int) ([]*MovieStarRelation, error) {
	queryString := `
	SELECT movie_id, star_id, role, details, order_no
	FROM movie_stars
	WHERE movie_id = $1`
	q := dbx.FromContext(ctx, r.db)
	rows, err := q.Query(ctx, queryString, id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	var relations []*MovieStarRelation
	for rows.Next() {
		var r MovieStarRelation
		err = rows.Scan(
			&r.MovieID,
			&r.StarID,
			&r.Role,
			&r.Details,
			&r.OrderNo,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}

		relations = append(relations, &r)
	}
	return relations, nil
}
