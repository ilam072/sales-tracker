package dto

type CreateCategory struct {
	Name string `json:"name" validate:"required"`
}

type GetCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateCategory struct {
	Name string `json:"name" validate:"required"`
}

type Categories struct {
	Categories []GetCategory `json:"categories"`
}
