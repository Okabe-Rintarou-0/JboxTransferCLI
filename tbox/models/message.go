package models

type ErrorMessage struct {
	Code    string `json:"code"`
	Status  int64  `json:"status"`
	Message string `json:"message"`
}
