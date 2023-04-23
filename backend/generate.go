package backend

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fumiama/paper-manager/backend/global"
)

func init() {
	apimap["/api/genFile"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		user := usertokens.Get(r.Header.Get("Authorization"))
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		conf := global.GenerateConfig{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&conf)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		docf, err := global.FileDB.GenerateFile(&conf)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		_, _ = io.Copy(w, docf)
	}}
}
