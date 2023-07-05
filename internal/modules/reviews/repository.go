package reviews

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
	"github.com/mkuptsov/movie-reviews/internal/dbx"
	"github.com/mkuptsov/movie-reviews/internal/modules/movies"
)

type Repository struct {
	db         *pgxpool.Pool
	moviesRepo *movies.Repository
}

func NewRepository(db *pgxpool.Pool, moviesRepo *movies.Repository) *Repository {
	return &Repository{
		db:         db,
		moviesRepo: moviesRepo,
	}
}

func (r *Repository) Create(ctx context.Context, review *Review) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		if err := r.moviesRepo.Lock(ctx, tx, review.MovieID); err != nil {
			return err
		}

		err := tx.QueryRow(
			ctx,
			"insert into reviews (movie_id, user_id, title, content, rating) values ($1, $2, $3, $4, $5) returning id, created_at",
			review.MovieID, review.UserID, review.Title, review.Content, review.Rating).
			Scan(&review.ID, &review.CreatedAt)

		switch {
		case dbx.IsUniqueViolation(err, ""):
			return apperrors.AlreadyExists("review", "(movie_id,user_id)", fmt.Sprintf("(%d,%d)", review.MovieID, review.UserID))
		case err != nil:
			return apperrors.Internal(err)
		}

		return r.recalculateMovieRating(ctx, review.MovieID)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, reviewID int) (*Review, error) {
	q := dbx.FromContext(ctx, r.db)
	var review Review

	err := q.QueryRow(
		ctx,
		"select id, movie_id, user_id, title, content, rating, created_at from reviews where deleted_at is null and id = $1",
		reviewID).
		Scan(&review.ID, &review.MovieID, &review.UserID, &review.Title, &review.Content, &review.Rating, &review.CreatedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("review", "id", reviewID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &review, nil
}

func (r *Repository) GetPaginated(ctx context.Context, movieID, userID *int, offset int, limit int) ([]*Review, int, error) {
	selectQuery := dbx.StatementBuilder.
		Select("id", "movie_id", "user_id", "title", "content", "rating", "created_at").
		From("reviews").
		Where("deleted_at is null").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	countQuery := dbx.StatementBuilder.
		Select("count(*)").
		From("reviews").
		Where("deleted_at is null")

	if movieID != nil {
		selectQuery = selectQuery.Where("movie_id = ?", *movieID)
		countQuery = countQuery.Where("movie_id = ?", *movieID)
	}

	if userID != nil {
		selectQuery = selectQuery.Where("user_id = ?", *userID)
		countQuery = countQuery.Where("user_id = ?", *userID)
	}

	b := &pgx.Batch{}
	if err := dbx.QueueBatchSelect(b, selectQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	if err := dbx.QueueBatchSelect(b, countQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		var review Review
		if err = rows.Scan(&review.ID, &review.MovieID, &review.UserID, &review.Title, &review.Content, &review.Rating, &review.CreatedAt); err != nil {
			return nil, 0, apperrors.Internal(err)
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return reviews, total, nil
}

func (r *Repository) Update(ctx context.Context, reviewID, userID int, title, content string, rating int) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		review, err := r.GetByID(ctx, reviewID)
		if err != nil {
			return err
		}
		if err = r.moviesRepo.Lock(ctx, tx, review.MovieID); err != nil {
			return err
		}
		n, err := tx.Exec(
			ctx,
			"update reviews set title = $1, content = $2, rating = $3 where deleted_at is null and id = $4 and user_id = $5",
			title, content, rating, reviewID, userID)
		if err != nil {
			return apperrors.Internal(err)
		}

		if n.RowsAffected() == 0 {
			return r.specifyModificationError(ctx, reviewID, userID)
		}

		return r.recalculateMovieRating(ctx, review.MovieID)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, reviewID, userID int) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		review, err := r.GetByID(ctx, reviewID)
		if err != nil {
			return err
		}
		if err = r.moviesRepo.Lock(ctx, tx, review.MovieID); err != nil {
			return err
		}

		n, err := r.db.Exec(
			ctx,
			"update reviews set deleted_at = now() where deleted_at is null and id = $1 and user_id = $2",
			reviewID, userID)
		if err != nil {
			return apperrors.Internal(err)
		}

		if n.RowsAffected() == 0 {
			return r.specifyModificationError(ctx, reviewID, userID)
		}

		return r.recalculateMovieRating(ctx, review.MovieID)
	})
	if err != nil {
		return apperrors.EnsureInternal(err)
	}

	return nil
}

func (r *Repository) specifyModificationError(ctx context.Context, reviewID, userID int) error {
	// Review is not found by reviewID and userID then there are two possibilities:
	// 1. Review with reviewID does not exist
	// 2. Review with reviewID exists, but it is not owned by userID
	review, err := r.GetByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if review.UserID != userID {
		return apperrors.Forbidden(fmt.Sprintf("review with id %d is not owned by user with id %d", reviewID, userID))
	}

	// If we got here, then something is wrong
	return apperrors.Internal(fmt.Errorf("unexpected error creating/updating review with id %d", reviewID))
}

func (r *Repository) recalculateMovieRating(ctx context.Context, movieID int) error {
	q := dbx.FromContext(ctx, r.db)
	n, err := q.Exec(ctx, "UPDATE movies SET avg_rating = (SELECT avg(rating) FROM reviews WHERE deleted_at IS NULL and movie_id = $1) where id = $1", movieID)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("movie", "id", movieID)
	}
	return nil
}
