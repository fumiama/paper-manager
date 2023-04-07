package backend

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	sql "github.com/FloatTech/sqlite"
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

type filelist struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Desc  string  `json:"description"`
	Size  float64 `json:"size"`
	Ques  int     `json:"questions"`
	Auth  string  `json:"author"`
	Date  string  `json:"datetime"`
	Per   uint    `json:"percent"`
}

func init() {
	apimap["/api/getFileList"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		count := -1
		var err error
		countstr := r.URL.Query().Get("count")
		if countstr != "" {
			count, err = strconv.Atoi(countstr)
			if err != nil {
				writeresult(w, codeError, nil, err.Error(), typeError)
				return
			}
		}
		lst, err := global.FileDB.ListUploadedFile()
		if err != nil && err != sql.ErrNullResult {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		if count > 0 && len(lst) > count {
			lst = lst[:count]
		}
		result := make([]filelist, len(lst))
		for i, v := range lst {
			result[i].ID = *v.ID
			j := strings.LastIndex(v.Path, "/")
			if j <= 0 {
				result[i].Title = v.Path
			} else {
				result[i].Title = v.Path[j+1:]
			}
			result[i].Desc = v.Desc
			result[i].Size = float64(v.Size) / 1024 / 1024 // MB
			result[i].Ques = v.QuesC
			result[i].Auth = v.UpName
			result[i].Date = time.Unix(v.UpTime, 0).Format(chineseYYMMDDLayout)
			if !v.HasntAnalyzed {
				result[i].Per = 100
			} else {
				result[i].Per = analyzeper.Get(*v.ID)
			}
		}
		writeresult(w, codeSuccess, &result, messageOk, typeSuccess)
	}}
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
			M string `json:"msg"`
		}
		go func() {
			_, err = global.FileDB.AddFile(id, reg, istemp, func(u uint) { analyzeper.Set(id, u) })
			ch <- struct{}{}
			close(ch)
		}()
		select {
		case <-time.After(time.Second):
			writeresult(w, codeSuccess, &message{M: "正在分析, 请耐心等待..."}, messageOk, typeSuccess)
			return
		case <-ch:
			if err != nil {
				writeresult(w, codeError, nil, err.Error(), typeError)
				return
			}
			writeresult(w, codeSuccess, &message{M: "分析完成"}, messageOk, typeSuccess)
		}
	}}
}

// PaperHandler serves protected contents in global.FileFolder
func PaperHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod("GET", w, r) {
		return
	}
	global.UserDB.VisitAPI()
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}
	fn := r.URL.Path[6:]
	if fn == "" {
		http.Error(w, "400 Bad Request: empty path", http.StatusBadRequest)
		return
	}
	name := global.FileFolder + fn
	logrus.Infoln("[file.FileHandler] serve", name)
	http.ServeFile(w, r, name)
}
