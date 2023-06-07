package stars

import "time"

type Star struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"first_name"`
	MiddleName *string    `json:"middle_name"`
	LastName   string     `json:"last_name"`
	BirthDate  time.Time  `json:"birth_date"`
	BirthPlace *string    `json:"birth_place"`
	DeathDate  *time.Time `json:"death_date"`
	Bio        *string    `json:"bio"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
