package http

type storeRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type updateRequest struct {
	Name string `json:"name" form:"name" validate:"omitempty"`
}
