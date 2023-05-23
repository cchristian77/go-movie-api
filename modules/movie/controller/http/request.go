package http

import "github.com/google/uuid"

type storeRequest struct {
	Title    string      `json:"title" form:"title" validate:"required" `
	Duration int32       `json:"duration" form:"duration" validate:"required"`
	Year     int32       `json:"year" form:"year" validate:"required"`
	Synopsis string      `json:"synopsis" form:"synopsis" validate:"required"`
	GenreIDs []uuid.UUID `json:"genre_ids" form:"genre_ids" validate:"required,min=1"`
}

type updateRequest struct {
	Title    string      `json:"title" form:"title" validate:"omitempty"`
	Duration int32       `json:"duration" form:"duration" validate:"omitempty"`
	Year     int32       `json:"year" form:"year" validate:"omitempty"`
	Synopsis string      `json:"synopsis" form:"synopsis" validate:"omitempty"`
	GenreIDs []uuid.UUID `json:"genre_ids" form:"genre_ids" validate:"omitempty,min=1"`
}
