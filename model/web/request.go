package web

import "time"

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

type CategoryAddRequest struct {
	ID            int    `json:"id" validate:"required"`
	Name          string `json:"name" validate:"required,max=255"`
	Specification string `json:"specification" validate:"required,max=255"`
}

type CategoryUpdateRequest struct {
	ID            int    `json:"id" validate:"required"`
	Name          string `json:"name" validate:"required,max=255"`
	Specification string `json:"specification" validate:"required,max=255"`
}

type ItemAddRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	CategoryID  int     `json:"category_id" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Description string  `json:"description" validate:"required,max=255"`
}

type ItemUpdateRequest struct {
	ID          int     `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"required,max=255"`
	CategoryID  int     `json:"category_id" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Description string  `json:"description" validate:"required,max=255"`
}

type ActivityAddRequest struct {
	ItemID        int       `json:"item_id" validate:"required"`
	Action        string    `json:"action" validate:"required"`
	QuantityChane int       `json:"quantity_change" validate:"required"`
	Timestamp     time.Time `json:"timestamp" validate:"required"`
	PerformedBy   int       `json:"performed_by" validate:"required,max=255"`
}
