package gameerr

import "fmt"

type GameServerErr struct {
	Code    string
	Message string
}

func (e *GameServerErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code string, msg string) *GameServerErr {
	return &GameServerErr{
		Code:    code,
		Message: msg,
	}
}
