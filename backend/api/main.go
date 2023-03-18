package api

import (
	"encoding/json"
	"net/http"

	"github.com/fumiama/paper-manager/backend/utils"
)

type apihandler struct {
	md string
	do func(w http.ResponseWriter, r *http.Request)
}

func (h *apihandler) handle(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod(h.md, w, r) {
		return
	}
	h.do(w, r)
}

var apimap = make(map[string]*apihandler, 512)

func init() {
	apimap["/api/getLoginSalt"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
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
	}}

	apimap["/api/login"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type loginbody struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var body loginbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		ret, err := login(body.Username, body.Password)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, ret, messageOk, typeSuccess)
	}}

	apimap["/api/getUserInfo"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		ret, err := getUserInfo(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, ret, messageOk, typeSuccess)
	}}

	apimap["/api/logout"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		err := logout(r.Header.Get("Authorization"))
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, nil, messageOk, typeSuccess)
	}}

	apimap["/api/getUsersCount"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		n, err := getUsersCount(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, n, messageOk, typeSuccess)
	}}
}

// Handler serves all backend /api call
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}

	if h, ok := apimap[r.URL.Path]; ok {
		h.handle(w, r)
		return
	}

	http.Error(w, "404 Not Found", http.StatusNotFound)
}
