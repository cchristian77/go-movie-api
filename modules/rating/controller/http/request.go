package http

import "github.com/google/uuid"

type storeRequest struct {
	Rating    float32   `json:"rating" form:"rating" validate:"required"`
	Comment   string    `json:"comment" form:"comment" validate:"omitempty"`
	MovieUuid uuid.UUID `json:"movie_id" form:"movie_id" validate:"required"`
}

type updateRequest struct {
	Rating  float32 `json:"rating" form:"rating" validate:"omitempty"`
	Comment string  `json:"comment" form:"comment" validate:"omitempty"`
}
