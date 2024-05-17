package web

type UserRegisterRequest struct {
	FullName string `json:"full_name" validate:"required,max=255"`
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	Role     string `json:"role" validate:"required,eq=admin|eq=user"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	FullName string `json:"full_name" validate:"required,min=1,max=255"`
	Username string `json:"username" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	Role     string `json:"role" validate:"required,eq=admin|eq=user"`
}
