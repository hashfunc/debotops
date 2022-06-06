package server

type Status string

const (
	StatusOK    Status = "OK"
	StatusError Status = "ERROR"
)

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}
