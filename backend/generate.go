package backend

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/fumiama/go-docx"
	"github.com/fumiama/paper-manager/backend/global"
)

var genfilecache = ttl.NewCache[int, *docx.Docx](time.Minute * 10)

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
		genfilecache.Set(*user.ID, docf)
		writeresult(w, codeSuccess, "请在10分钟内下载, 且不要在下载完成前关闭页面, 云端不会保存", messageOk, typeSuccess)
	}}

	apimap["/api/dlGen"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		user := usertokens.Get(r.Header.Get("Authorization"))
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		docf := genfilecache.Get(*user.ID)
		if docf == nil {
			writeresult(w, codeError, nil, os.ErrNotExist.Error(), typeError)
			return
		}
		_, _ = io.Copy(w, docf)
	}}
}
