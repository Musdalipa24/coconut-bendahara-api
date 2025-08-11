package dto

type ListResponseOK struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ListResponseSaldo struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Saldo   int64  `json:"total"`
	Message string `json:"message"`
}

type ListResponseToken struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ListResponseError struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
