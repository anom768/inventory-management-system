package web

import "net/http"

type OK struct {
	Code    uint   `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewOk(message string) *OK {
	return &OK{http.StatusOK, "status ok", message}
}

type Created struct {
	Code    uint   `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewCreated(message string) *Created {
	return &Created{http.StatusCreated, "status created", message}
}

type BadRequestResponse struct {
	Code    uint   `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewBadRequestResponse(message string) *BadRequestResponse {
	return &BadRequestResponse{http.StatusBadRequest, "status bad request", message}
}

type Unauthorized struct {
	Code    uint   `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewUnauthorized(message string) *Unauthorized {
	return &Unauthorized{http.StatusUnauthorized, "status unauthorized", message}
}

type InternalServerError struct {
	Code    uint   `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{http.StatusInternalServerError, "status internal server error", message}
}

type ResponseModel struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func NewResponseModel(data any) *ResponseModel {
	return &ResponseModel{http.StatusOK, "status ok", data}
}
