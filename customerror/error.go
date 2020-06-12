package customerror

import "fmt"

type Error struct {
	Code     int `json:"code"`
	HTTPCode int
	Message  string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[Error] code: %d, httpCode: %d, message: %s", e.Code, e.HTTPCode, e.Message)
}
