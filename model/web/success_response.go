package web

import "net/http"

type SuccessResponseData interface {
	Code() int
	Status() string
	Message() string
	Data() any
}

type successResponseData struct {
	ResCode    int    `json:"code"`
	ResStatus  string `json:"status"`
	ResMessage string `json:"message"`
	ResData    any    `json:"data"`
}

func NewStatusOKData(message string, data any) SuccessResponseData {
	return &successResponseData{
		ResCode:    http.StatusOK,
		ResStatus:  "status ok",
		ResMessage: message,
		ResData:    data,
	}
}

func (s *successResponseData) Code() int {
	return s.ResCode
}

func (s *successResponseData) Status() string {
	return s.ResStatus
}

func (s *successResponseData) Message() string {
	return s.ResMessage
}

func (s *successResponseData) Data() any {
	return s.ResData
}

type SuccessResponseMessage interface {
	Code() int
	Status() string
	Message() string
}

type successResponseMessage struct {
	ResCode    int    `json:"code"`
	ResStatus  string `json:"status"`
	ResMessage string `json:"message"`
}

func NewStatusCreated(message string) SuccessResponseMessage {
	return &successResponseMessage{
		ResCode:    http.StatusCreated,
		ResStatus:  "status created",
		ResMessage: message,
	}
}

func NewStatusOKMessage(message string) SuccessResponseMessage {
	return &successResponseData{
		ResCode:    http.StatusOK,
		ResStatus:  "status ok",
		ResMessage: message,
	}
}

func (s *successResponseMessage) Code() int {
	return s.ResCode
}

func (s *successResponseMessage) Status() string {
	return s.ResStatus
}

func (s *successResponseMessage) Message() string {
	return s.ResMessage
}
