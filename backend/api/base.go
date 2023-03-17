package api

import (
	"encoding/json"
	"io"
)

const (
	codeError   = -1
	codeSuccess = 0
	codeTimeout = 401
)

const (
	typeSuccess = "success"
	typeError   = "error"
)

const (
	messageOk = "ok"
)

type base struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
	Type    string `json:"type"`
}

func writeresult(w io.Writer, c int, r any, m, t string) error {
	return json.NewEncoder(w).Encode(&base{
		Code:    c,
		Result:  r,
		Message: m,
		Type:    t,
	})
}
