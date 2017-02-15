package opc

import "fmt"

type OracleError struct {
	StatusCode int
	Message    string `json:"message"`
}

func (e OracleError) Error() string {
	return fmt.Sprintf("%s: %s", e.StatusCode, e.Message)
}
