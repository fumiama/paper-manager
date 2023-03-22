package backend

import (
	"net/http"

	"github.com/fumiama/paper-manager/backend/global"
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
	apimap["/api/getAnnualVisits"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		if !user.IsSuper() {
			writeresult(w, codeError, nil, errNoSetRolePermission.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, global.UserDB.GetAnnualAPIVisitCount(), messageOk, typeSuccess)
	}}
}

// APIHandler serves all backend /api call
func APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}

	if h, ok := apimap[r.URL.Path]; ok {
		global.UserDB.VisitAPI()
		h.handle(w, r)
		return
	}

	http.Error(w, "404 Not Found", http.StatusNotFound)
}
