package api

import (
	"encoding/json"
	"net/http"

	"github.com/fumiama/paper-manager/backend/utils"
)

// Handler serves all backend /api call
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}
	if r.URL.Path == "/api/getLoginSalt" {
		if !utils.IsMethod("GET", w, r) {
			return
		}
		username := r.URL.Query().Get("username")
		if username == "" {
			writeresult(w, codeError, nil, "empty username", typeError)
			return
		}
		salt, err := getLoginSalt(username)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, salt, messageOk, typeSuccess)
		return
	}
	if r.URL.Path == "/api/login" {
		if !utils.IsMethod("POST", w, r) {
			return
		}
		defer r.Body.Close()
		var body loginbody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		r, err := login(body.Username, body.Password)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, r, messageOk, typeSuccess)
		return
	}
	if r.URL.Path == "/api/getUserInfo" {
		if !utils.IsMethod("GET", w, r) {
			return
		}
		token := r.Header.Get("Authorization")
		r, err := getUserInfo(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, r, messageOk, typeSuccess)
		return
	}
	if r.URL.Path == "/api/logout" {
		if !utils.IsMethod("GET", w, r) {
			return
		}
		err := logout(r.Header.Get("Authorization"))
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, nil, messageOk, typeSuccess)
		return
	}
	if !utils.IsMethod("GET", w, r) {
		return
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

type loginbody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
