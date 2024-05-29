package web

import "net/http"

type ErrorResponse interface {
	Code() int
	Status() string
	Message() string
}

type errorResponse struct {
	ErrCode    int    `json:"code"`
	ErrStatus  string `json:"status"`
	ErrMessage string `json:"message"`
}

func (e *errorResponse) Code() int {
	return e.ErrCode
}

func (e *errorResponse) Status() string {
	return e.ErrStatus
}

func (e *errorResponse) Message() string {
	return e.ErrMessage
}

func NewBadRequestError(message string) ErrorResponse {
	return &errorResponse{
		ErrCode:    http.StatusBadRequest,
		ErrStatus:  "status bad request",
		ErrMessage: message,
	}
}

func NewUnauthorizedError(message string) ErrorResponse {
	return &errorResponse{
		ErrCode:    http.StatusUnauthorized,
		ErrStatus:  "status unauthorized",
		ErrMessage: message,
	}
}

func NewNotFoundError(message string) ErrorResponse {
	return &errorResponse{
		ErrCode:    http.StatusNotFound,
		ErrStatus:  "status not found",
		ErrMessage: message,
	}
}

func NewInternalServerErrorError(message string) ErrorResponse {
	return &errorResponse{
		ErrCode:    http.StatusInternalServerError,
		ErrStatus:  "status internal server error",
		ErrMessage: message,
	}
}
