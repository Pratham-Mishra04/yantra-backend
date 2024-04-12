package schemas

import "github.com/lib/pq"

type PostCreateSchema struct { // from request
	Content string         `json:"content" validate:"required,max=2000"`
	Tags    pq.StringArray `json:"tags" validate:"dive,alphanum"`
}

type PostUpdateSchema struct {
	Content string          `json:"content" validate:"max=2000"`
	Tags    *pq.StringArray `json:"tags" validate:"dive,alphanum"`
}
