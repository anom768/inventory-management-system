package web

type ResponseMessage struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseModel struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}
