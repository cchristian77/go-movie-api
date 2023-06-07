package http

type updateRequest struct {
	FullName string `json:"full_name" form:"full_name" validate:"omitempty"`
	Username string `json:"username" form:"username" validate:"omitempty,min=6"`
	Email    string `json:"email" form:"email" validate:"omitempty,email"`
}
