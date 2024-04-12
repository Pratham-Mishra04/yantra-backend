package schemas

type ResourceBucketCreateSchema struct {
	Title               string `json:"title" validate:"required,max=50"`
	Description         string `json:"description" validate:"max=500"`
	OnlyAdminViewAccess bool   `json:"onlyAdminViewAccess" validate:"required"`
	OnlyAdminEditAccess bool   `json:"onlyAdminEditAccess" validate:"required"`
}

type ResourceBucketEditSchema struct {
	Title               string  `json:"title" validate:"max=50"`
	Description         *string `json:"description" validate:"max=500"`
	OnlyAdminViewAccess *bool   `json:"onlyAdminViewAccess"`
	OnlyAdminEditAccess *bool   `json:"onlyAdminEditAccess"`
}

type ResourceFileCreateSchema struct {
	Title       string `json:"title" validate:"required,max=50"`
	Description string `json:"description" validate:"max=500"`
	Link        string `json:"link"`
}

type ResourceFileEditSchema struct {
	Title       string `json:"title" validate:"max=50"`
	Description string `json:"description" validate:"max=500"`
}
