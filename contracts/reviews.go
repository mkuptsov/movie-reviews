package contracts

import (
	"strconv"
	"time"
)

type Review struct {
	ID        int        `json:"id"`
	MovieID   int        `json:"movie_id"`
	UserID    int        `json:"user_id"`
	Rating    int        `json:"rating"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type GetReviewsRequest struct {
	PaginatedRequest
	MovieID *int `query:"movieId"`
	UserID  *int `query:"userId"`
}

func (r *GetReviewsRequest) ToQueryParams() map[string]string {
	params := r.PaginatedRequest.ToQueryParams()
	if r.MovieID != nil {
		params["movieId"] = strconv.Itoa(*r.MovieID)
	}
	if r.UserID != nil {
		params["userId"] = strconv.Itoa(*r.UserID)
	}
	return params
}

type GetReviewRequest struct {
	ReviewID int `param:"reviewId" validate:"nonzero"`
}

type CreateReviewRequest struct {
	MovieID int    `json:"movie_id" validate:"nonzero"`
	UserID  int    `json:"user_id" validate:"nonzero"`
	Rating  int    `json:"rating" validate:"min=1,max=10"`
	Title   string `json:"title" validate:"min=3,max=255"`
	Content string `json:"content" validate:"min=20,max=2000"`
}

type UpdateReviewRequest struct {
	ReviewID int    `json:"-" param:"reviewId" validate:"nonzero"`
	UserID   int    `json:"-" param:"userId" validate:"nonzero"`
	Rating   int    `json:"rating" validate:"min=1,max=10"`
	Title    string `json:"title" validate:"min=3,max=255"`
	Content  string `json:"content" validate:"min=20,max=2000"`
}

type DeleteReviewRequest struct {
	ReviewID int `param:"reviewId" validate:"nonzero"`
	UserID   int `param:"userId" validate:"nonzero"`
}
