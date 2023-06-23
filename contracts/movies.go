package contracts

import (
	"strconv"
	"time"
)

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MovieDetails struct {
	Movie
	Description string         `json:"description"`
	Version     int            `json:"version"`
	Genres      []*Genre       `json:"genres"`
	Cast        []*MovieCredit `json:"cast"`
}

type GetMovieByIDRequest struct {
	ID int `param:"id" validate:"nonzero"`
}

type GetMoviesRequest struct {
	PaginatiedRequest
	StarID     *int    `query:"starID"`
	SearchTerm *string `query:"q"`
}

func (r *GetMoviesRequest) ToQueryParams() map[string]string {
	params := r.PaginatiedRequest.ToQueryParams()
	if r.StarID != nil {
		params["starID"] = strconv.Itoa(*r.StarID)
	}
	if r.SearchTerm != nil {
		params["q"] = *r.SearchTerm
	}
	return params
}

type CreateMovieRequest struct {
	Title       string             `json:"title" validate:"min=1,max=255"`
	ReleaseDate time.Time          `json:"release_date" validate:"nonzero"`
	Description string             `json:"description"`
	Genres      []int              `json:"genres"`
	Cast        []*MovieCreditInfo `json:"cast"`
}

type UpdateMovieRequest struct {
	ID          int                `param:"id" validate:"nonzero"`
	Title       string             `json:"title"`
	ReleaseDate time.Time          `json:"release_date"`
	Description string             `json:"description"`
	Version     int                `json:"version" validate:"min=0"`
	Genres      []int              `json:"genres"`
	Cast        []*MovieCreditInfo `json:"cast"`
}

type DeleteMovieRequest struct {
	ID int `param:"id" validate:"nonzero"`
}
