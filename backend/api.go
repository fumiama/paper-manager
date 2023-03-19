package backend

import (
	"encoding/json"
	"net/http"
	"strings"

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

	apimap["/api/register"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type registerbody struct {
			Username string `json:"username"`
			Mobile   string `json:"mobile"`
			Password string `json:"password"`
		}
		if r.Header.Get("Authorization") != "" {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body registerbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		ip := r.RemoteAddr
		i := strings.LastIndex(ip, ":")
		if i >= 0 {
			ip = ip[:i]
		}
		err = register(ip, body.Username, body.Mobile, body.Password)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "成功, 请耐心等待通知"}, messageOk, typeSuccess)
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

	apimap["/api/setPassword"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setpasswordbody struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body setpasswordbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		err = setUserPassword(*user.ID, body.Token, body.Password)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "成功, 请重新登录"}, messageOk, typeSuccess)
		_ = logout(token)
	}}

	apimap["/api/setContact"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setcontactbody struct {
			Token   string `json:"token"`
			Contact string `json:"contact"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body setcontactbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		err = setUserContact(*user.ID, body.Token, body.Contact)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		user.Cont = hideContact(body.Contact)
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "成功, 已将消息报告给课程组长"}, messageOk, typeSuccess)
	}}

	apimap["/api/setUserInfo"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setuserinfobody struct {
			Nick string `json:"nick"`
			Desc string `json:"desc"`
			Avtr string `json:"avtr"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body setuserinfobody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		err = setUserInfo(*user.ID, &body.Nick, &body.Desc, &body.Avtr)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		user.Nick = body.Nick
		user.Desc = body.Desc
		user.Avtr = body.Avtr
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "成功"}, messageOk, typeSuccess)
	}}

	apimap["/api/getMessageList"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		ret, err := getMessageList(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, ret, messageOk, typeSuccess)
	}}
}

// APIHandler serves all backend /api call
func APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}

	if h, ok := apimap[r.URL.Path]; ok {
		h.handle(w, r)
		return
	}

	http.Error(w, "404 Not Found", http.StatusNotFound)
}
