package backend

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/fumiama/paper-manager/backend/global"
	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

const (
	chineseYYMMDDLayout = "2006年01月02日"
)

// analyzeper 分析进度缓存
var analyzeper = ttl.NewCache[int, uint](time.Hour)

var (
	errNoAnalyzePermission = errors.New("no analyze permission")
)

func init() {
	apimap["/api/getFilePercent"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		idstr := r.URL.Query().Get("id")
		if idstr == "" {
			writeresult(w, codeError, nil, "empty id", typeError)
			return
		}
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, analyzeper.Get(id), messageOk, typeSuccess)
	}}

	apimap["/api/analyzeFile"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		istemp := r.URL.Query().Get("permanent") != "true"
		if !user.IsFileManager() && !istemp {
			writeresult(w, codeError, nil, errNoAnalyzePermission.Error(), typeError)
			return
		}
		idstr := r.URL.Query().Get("id")
		if idstr == "" {
			writeresult(w, codeError, nil, "empty id", typeError)
			return
		}
		id, err := strconv.Atoi(idstr)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		reg, err := global.UserDB.GetUserRegex(*user.ID)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		ch := make(chan struct{}, 1)
		type message struct {
			C int    `json:"code"` // C 0 success 1 pending
			M string `json:"msg"`
		}
		if analyzeper.Get(id) > 0 {
			writeresult(w, codeError, nil, "已在分析!", typeError)
			return
		}
		go func() {
			err = global.FileDB.AddFile(id, reg, istemp, func(u uint) { analyzeper.Set(id, u) })
			ch <- struct{}{}
			close(ch)
		}()
		select {
		case <-time.After(time.Second):
			writeresult(w, codeSuccess, &message{C: 1, M: "正在分析, 请耐心等待..."}, messageOk, typeSuccess)
			return
		case <-ch:
			if err != nil {
				writeresult(w, codeError, nil, err.Error(), typeError)
				return
			}
			writeresult(w, codeSuccess, &message{C: 0, M: "分析完成"}, messageOk, typeSuccess)
		}
	}}
}

// PaperHandler serves protected contents in global.PaperFolder
func PaperHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod("GET", w, r) {
		return
	}
	token := r.Header.Get("Authorization")
	user := usertokens.Get(token)
	if user == nil {
		writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
		return
	}
	global.UserDB.VisitAPI()
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}
	fn := r.URL.Path[7:] // skip /paper/
	if fn == "" {
		http.Error(w, "400 Bad Request: empty path", http.StatusBadRequest)
		return
	}
	if strings.HasPrefix(fn, "tmp/") {
		uidstr := fn[4:]
		i := strings.Index(uidstr, "/")
		if i <= 0 {
			writeresult(w, codeError, nil, "extract uid error", typeError)
			return
		}
		uid, err := strconv.Atoi(uidstr[:i])
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		if uid != *user.ID {
			writeresult(w, codeError, nil, errNoDownloadPermission.Error(), typeError)
			return
		}
	}
	name := global.PaperFolder + fn
	logrus.Infoln("[file.PaperHandler] serve", name)
	http.ServeFile(w, r, name)
}
