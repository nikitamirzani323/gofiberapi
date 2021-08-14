package model

type Response struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	Totalrecord int         `json:"totalrecord"`
	Record      interface{} `json:"record"`
}
