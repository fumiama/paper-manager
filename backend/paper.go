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
	errNoAnalyzePermission  = errors.New("no analyze permission")
	errNoDeletePermission   = errors.New("no delete permission")
	errNoDownloadPermission = errors.New("no download permission")
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

type filestatus struct {
	Name         string        `json:"name"`
	Size         float64       `json:"size"`
	Questions    []question    `json:"questions"`
	Duplications []duplication `json:"duplications"` // Duplications length == 10
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
	apimap["/api/getFileInfo"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var err error
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
		lst, err := global.FileDB.GetFileInfo(id)
		if err != nil && err != sql.ErrNullResult {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		result := filelist{
			ID:   id,
			Desc: lst.Desc,
			Size: float64(lst.Size) / 1024 / 1024, // MB
			Ques: lst.QuesC,
			Auth: lst.UpName,
			Date: time.Unix(lst.UpTime, 0).Format(chineseYYMMDDLayout),
		}
		j := strings.LastIndex(lst.Path, "/")
		if j <= 0 {
			result.Title = lst.Path
		} else {
			result.Title = lst.Path[j+1:]
		}
		if !lst.HasntAnalyzed {
			result.Per = 100
		} else {
			result.Per = analyzeper.Get(id)
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
	apimap["/api/delFile"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		if !user.IsSuper() {
			writeresult(w, codeError, nil, errNoDeletePermission.Error(), typeError)
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
		err = global.FileDB.DelFile(id, *user.ID, false)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, "删除成功", messageOk, typeSuccess)
	}}
	apimap["/api/dlFile"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		if !user.IsSuper() {
			writeresult(w, codeError, nil, errNoDeletePermission.Error(), typeError)
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
		lst, err := global.FileDB.GetFileInfo(id)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		type message struct {
			URL string `json:"url"`
		}
		if strings.HasPrefix(lst.Path, global.PaperFolder+"tmp/") {
			uidstr := lst.Path[17:]
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
			writeresult(w, codeSuccess, &message{URL: lst.Path[6:]}, messageOk, typeSuccess)
			return
		}
		if strings.HasPrefix(lst.Path, global.PaperFolder) {
			writeresult(w, codeSuccess, &message{URL: lst.Path[6:]}, messageOk, typeSuccess)
			return
		}
		writeresult(w, codeError, nil, "parse filepath error", typeError)
	}}
	apimap["/api/getFileStatus"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
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
		file, sz, err := global.FileDB.GetFile(id, *user.ID)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		qs, ds, err := parseFileQuestions(file.Questions)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, &filestatus{
			Name:         file.Class + ".docx",
			Size:         float64(sz) / 1024 / 1024, // MB
			Questions:    qs,
			Duplications: ds,
		}, messageOk, typeSuccess)
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
