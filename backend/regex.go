package backend

import (
	"encoding/json"
	"net/http"

	"github.com/fumiama/paper-manager/backend/global"
)

func getUserRegex(token string) (*global.Regex, error) {
	user := usertokens.Get(token)
	if user == nil {
		return nil, errInvalidToken
	}
	return global.UserDB.GetUserRegex(user, *user.ID)
}

func init() {
	apimap["/api/getUserRegex"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		reg, err := getUserRegex(r.Header.Get("Authorization"))
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, reg, messageOk, typeSuccess)
	}}

	apimap["/api/setUserRegex"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		defer r.Body.Close()
		reg := &global.Regex{}
		err := json.NewDecoder(r.Body).Decode(reg)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		err = global.UserDB.SetUserRegex(*user.ID, reg)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, "成功", messageOk, typeSuccess)
	}}
}
