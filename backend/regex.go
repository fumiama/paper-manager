package backend

import (
	"net/http"

	"github.com/fumiama/paper-manager/backend/global"
)

func getUserRegex(token string) (*global.Regex, error) {
	user := usertokens.Get(token)
	if user == nil {
		return nil, errInvalidToken
	}
	return global.UserDB.GetUserRegex(*user.ID)
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
}
