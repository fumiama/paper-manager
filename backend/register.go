package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/fumiama/paper-manager/backend/global"
)

var registerlimit = ttl.NewCache[string, bool](time.Minute * 10)

var (
	errRequestTooFast = errors.New("request too fast")
	errInvalidIP      = errors.New("invalid IP")
)

func register(ip, name, mobile, npwd string) error {
	if registerlimit.Get(ip) {
		return errRequestTooFast
	}
	if ip == "" {
		return errInvalidIP
	}
	registerlimit.Set(ip, true)
	return global.UserDB.NotifyRegister(ip, name, mobile, npwd)
}

func init() {
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
		writeresult(w, codeSuccess, &message{M: "已上报, 请耐心等待通知"}, messageOk, typeSuccess)
	}}
}
