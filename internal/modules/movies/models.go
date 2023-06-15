package movies

import (
	"time"

	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"
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
	Description string          `json:"description"`
	Version     int             `json:"version"`
	Genres      []*genres.Genre `json:"genres"`
}
