package schemas

type GroupCreateSchema struct {
	Title       string `json:"title" validate:"required,max=50"`
	Description string `json:"description" validate:"required,max=1000"`
}

type GroupUpdateSchema struct {
	Title       string `json:"title" validate:"max=50"`
	Description string `json:"description" validate:"max=1000"`
}
