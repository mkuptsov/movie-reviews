package genres

import "github.com/mkuptsov/movie-reviews/internal/dbx"

type Genre struct {
	ID   int    `param:"id"`
	Name string `json:"name"`
}

var _ dbx.Keyer = MovieGenreRelation{}

type MovieGenreRelation struct {
	MovieID int
	GenreID int
	OrderNo int
}

func (m MovieGenreRelation) Key() any {
	type MovieGenreRelationKey struct {
		MovieID, GenreID int
	}

	return MovieGenreRelationKey{
		MovieID: m.MovieID,
		GenreID: m.GenreID,
	}
}
