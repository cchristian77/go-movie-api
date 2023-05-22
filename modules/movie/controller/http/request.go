package http

type storeRequest struct {
	Title    string `json:"title" form:"title" validate:"required" `
	Duration int32  `json:"duration" form:"duration" validate:"required"`
	Year     int32  `json:"year" form:"year" validate:"required"`
	Synopsis string `json:"synopsis" form:"synopsis" validate:"required"`
}

type updateRequest struct {
	Title    string `json:"title" form:"title" validate:"omitempty"`
	Duration int32  `json:"duration" form:"duration" validate:"omitempty"`
	Year     int32  `json:"year" form:"year" validate:"omitempty"`
	Synopsis string `json:"synopsis" form:"synopsis" validate:"omitempty"`
}
