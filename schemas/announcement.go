package schemas

type AnnouncementCreateSchema struct {
	Title   string `json:"title" validate:"required,max=50"`
	Content string `json:"content" validate:"required,max=1000"`
}

type AnnouncementUpdateSchema struct {
	Title   string `json:"title" validate:"max=50"`
	Content string `json:"content" validate:"max=1000"`
}
