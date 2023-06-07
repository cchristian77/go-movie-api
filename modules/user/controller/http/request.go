package http

type storeRequest struct {
	FullName string `json:"full_name" form:"full_name" validate:"required"`
	Username string `json:"username" form:"username" validate:"required,min=6"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
	IsAdmin  bool   `json:"is_admin" form:"is_admin" validate:"omitempty"`
}

type updateRequest struct {
	FullName string `json:"full_name" form:"full_name" validate:"omitempty"`
	Username string `json:"username" form:"username" validate:"omitempty,min=6"`
	Email    string `json:"email" form:"email" validate:"omitempty,email"`
}

type loginRequest struct {
	Username string `json:"username" form:"username" validate:"required,min=6"`
	Password string `json:"password" form:"password" validate:"required"`
}
