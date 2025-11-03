package dto

type CreateCategory struct {
	Name string `json:"name" validate:"required"`
}

type GetCategory struct {
	Name string `json:"name"`
}

type UpdateCategory struct {
	Name string `json:"name" validate:"required"`
}

type Categories struct {
	Categories []GetCategory `json:"categories"`
}
