package dbx

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
)

func QueueBatchSelect(b *pgx.Batch, sb squirrel.SelectBuilder) error {
	sql, args, err := sb.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}
	b.Queue(sql, args...)
	return nil
}
