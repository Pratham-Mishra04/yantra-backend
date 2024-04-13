package schemas

type CreatePollRequest struct {
	Title         string   `json:"title" validate:"max=50"`
	Content       string   `json:"content" validate:"required,max=500"`
	Options       []string `json:"options"`
	IsMultiAnswer bool     `json:"isMultiAnswer"`
}

type EditPollRequest struct {
	Content string `json:"question" validate:"max=500"`
}
