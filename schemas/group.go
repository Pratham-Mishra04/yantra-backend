package schemas

type GroupCreateSchema struct {
	Title   string `json:"title" validate:"required,max=50"`
	Content string `json:"content" validate:"required,max=1000"`
}

type GroupUpdateSchema struct {
	Title   string `json:"title" validate:"max=50"`
	Content string `json:"content" validate:"max=1000"`
}
