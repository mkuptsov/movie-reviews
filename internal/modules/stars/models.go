package stars

import (
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/dbx"
)

type Star struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate time.Time  `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type StarDetails struct {
	Star
	MiddleName *string `json:"middle_name,omitempty"`
	BirthPlace *string `json:"birth_place,omitempty"`
	Bio        *string `json:"bio,omitempty"`
}

type MovieCredit struct {
	Star    Star    `json:"star"`
	Role    string  `json:"role"`
	Details *string `json:"details"`
}

type MovieCreditInfo struct {
	StarID  int     `json:"star_id"`
	Role    string  `json:"role"`
	Details *string `json:"details"`
}

var _ dbx.Keyer = MovieStarRelation{}

type MovieStarRelation struct {
	MovieID int
	StarID  int
	Role    string
	Details *string
	OrderNo int
}

func (m MovieStarRelation) Key() any {
	type MovieStarRelationKey struct {
		MovieID int
		StarID  int
		Role    string
	}

	return MovieStarRelationKey{
		MovieID: m.MovieID,
		StarID:  m.StarID,
		Role:    m.Role,
	}
}
