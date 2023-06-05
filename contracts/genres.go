package contracts

type Genre struct {
	ID   int    `param:"id"`
	Name string `json:"name"`
}

type GetGenreByIDRequest struct {
	ID int `param:"id" validate:"nonzero"`
}

type CreateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=50"`
}

type UpdateGenreRequest struct {
	ID   int    `param:"id" validate:"nonzero"`
	Name string `json:"name" validate:"min=3,max=50"`
}

type DeleteGenreRequest struct {
	ID int `param:"id" validate:"nonzero"`
}
