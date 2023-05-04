package backend

import (
	"container/heap"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	sql "github.com/FloatTech/sqlite"
	"github.com/fumiama/paper-manager/backend/global"
	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

var (
	errNoDeletePermission   = errors.New("no delete permission")
	errExtractUIDError      = errors.New("extract uid error")
	errParseFilePath        = errors.New("parse filepath error")
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

func getFileList(count, uid int, istemp *bool) ([]filelist, error) {
	lst, err := global.FileDB.ListUploadedFile(istemp, uid)
	if err != nil && err != sql.ErrNullResult {
		return nil, err
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
	return result, nil
}

func getFileInfo(lstid int) (*filelist, error) {
	lst, err := global.FileDB.ListFileByID(lstid)
	if err != nil && err != sql.ErrNullResult {
		return nil, err
	}
	result := filelist{
		ID:   lstid,
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
		result.Per = analyzeper.Get(lstid)
	}
	return &result, nil
}

func dlFile(lstid int, user *global.User) (string, error) {
	lst, err := global.FileDB.ListFileByID(lstid)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(lst.Path, global.PaperFolder+"tmp/") {
		uidstr := lst.Path[17:]
		i := strings.Index(uidstr, "/")
		if i <= 0 {
			return "", errExtractUIDError
		}
		uid, err := strconv.Atoi(uidstr[:i])
		if err != nil {
			return "", err
		}
		if uid != *user.ID {
			return "", errNoDownloadPermission
		}
		return lst.Path[6:], nil
	}
	if strings.HasPrefix(lst.Path, global.PaperFolder) {
		return lst.Path[6:], nil
	}
	return "", errParseFilePath
}

type filestatus struct {
	Name         string        `json:"name"`
	Size         float64       `json:"size"`
	Rate         float64       `json:"rate"`
	Questions    []question    `json:"questions"`
	Duplications []duplication `json:"duplications"` // Duplications length == 10
}

func getFileStatus(lstid int, user *global.User) (*filestatus, error) {
	file, lst, err := global.FileDB.GetFile(lstid, *user.ID)
	if err != nil {
		return nil, err
	}
	qs, ds, filerate, err := parseFileQuestions(file.Questions, lst.IsTemp, true)
	if err != nil {
		return nil, err
	}
	return &filestatus{
		Name:         file.Class + ".docx",
		Size:         float64(lst.Size) / 1024 / 1024, // MB
		Rate:         filerate * 100,
		Questions:    qs,
		Duplications: ds,
	}, nil
}

type filestatusdup struct {
	Name         string        `json:"name"`
	Size         float64       `json:"size"`
	Rate         int           `json:"rate"`
	Questions    []question    `json:"questions"`
	Duplications []duplication `json:"duplications"` // Duplications length == 10
	Files        []duplication `json:"files"`        // Files is {dup: desc}[] in yearstart..yearend
}

func checkFileDup(lstid int, user *global.User, yearstart, yearend global.StudyYear) (*filestatusdup, error) {
	file, lst, err := global.FileDB.GetFile(lstid, *user.ID)
	if err != nil {
		return nil, err
	}
	files, err := global.FileDB.GetFilesByYearRange(yearstart, yearend)
	if err != nil {
		return nil, err
	}
	qs := make([]question, 0, 16)
	filesdups := make([]duplication, 0, 16)
	ques := make([]global.QuestionJSON, 0, 64)
	myques := make([]global.QuestionJSON, 0, 64)
	err = json.Unmarshal(file.Questions, &myques)
	if err != nil {
		return nil, err
	}
	dh := make(duplications, 0, 16)
	heap.Init(&dh)
	filerate := 0
	isfirstmy := true
	totl := 0
	for _, f := range files {
		if f.ListID == lstid {
			continue
		}
		ques = ques[:0]
		err := json.Unmarshal(f.Questions, &ques)
		if err != nil {
			return nil, err
		}
		sum := 0.0
		cnt := 0
		for _, q := range myques {
			if isfirstmy {
				qs = append(qs, question{
					Count: len(q.Sub),
					Point: q.Points,
					Name:  q.Name,
				})
			}
			for i, subq := range q.Sub {
				p, err := getQuestionDupFromPaper(subq, ques, lst.IsTemp)
				if err != nil {
					return nil, err
				}
				heap.Push(&dh, duplication{
					Percent: int(math.Round(p * 100)),
					Name:    q.Name + "." + strconv.Itoa(i+1),
				})
				sum += p
				cnt++
			}
		}
		isfirstmy = false
		pc := int(math.Round(sum * 100 / float64(cnt)))
		flst, err := f.GetList(&global.FileDB)
		name := ""
		if err != nil {
			name = err.Error()
		} else {
			name = flst.Desc
		}
		filesdups = append(filesdups, duplication{
			Percent: pc,
			Name:    name,
		})
		filerate += pc
		totl++
	}
	if totl == 0 {
		return nil, sql.ErrNullResult
	}
	i := dh.Len()
	ds := make([]duplication, 10)
	if i > 10 {
		i = 10
	} else {
		for j := i; j < 10; j++ {
			ds[j] = duplication{Name: "N/A"}
		}
	}
	for i--; i >= 0; i-- {
		ds[i] = heap.Pop(&dh).(duplication)
	}
	return &filestatusdup{
		Name:         file.Class + ".docx",
		Size:         float64(lst.Size) / 1024 / 1024, // MB
		Rate:         filerate / totl,
		Questions:    qs,
		Duplications: ds,
		Files:        filesdups,
	}, nil
}

func init() {
	apimap["/api/getFileList"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		istemp := (*bool)(nil)
		permanent := r.URL.Query().Get("permanent")
		if permanent != "" {
			b := permanent != "true"
			istemp = &b
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
		lst, err := getFileList(count, *user.ID, istemp)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, &lst, messageOk, typeSuccess)
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
		inf, err := getFileInfo(id)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, inf, messageOk, typeSuccess)
	}}

	apimap["/api/delFile"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		istemp := r.URL.Query().Get("permanent") != "true"
		if !user.IsSuper() && !istemp {
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
		err = global.FileDB.DelFile(id, *user.ID, istemp)
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
		type message struct {
			URL string `json:"url"`
		}
		u, err := dlFile(id, user)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, &message{URL: u}, messageOk, typeSuccess)
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
		fstat, err := getFileStatus(id, user)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, fstat, messageOk, typeSuccess)
	}}

	apimap["/api/checkFileDup"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
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
		ysstr := r.URL.Query().Get("ys")
		if ysstr == "" {
			writeresult(w, codeError, nil, "empty ys", typeError)
			return
		}
		ys, err := strconv.Atoi(ysstr)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		yestr := r.URL.Query().Get("ye")
		if yestr == "" {
			writeresult(w, codeError, nil, "empty ye", typeError)
			return
		}
		ye, err := strconv.Atoi(yestr)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		fstat, err := checkFileDup(id, user, global.StudyYear(ys), global.StudyYear(ye))
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, fstat, messageOk, typeSuccess)
	}}
}

// FileHandler serves contents in global.FileFolder
func FileHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod("GET", w, r) {
		return
	}
	global.UserDB.VisitAPI()
	if r.URL.Path[0] != '/' {
		r.URL.Path = "/" + r.URL.Path
	}
	fn := r.URL.Path[6:] // skip /file/
	if fn == "" {
		http.Error(w, "400 Bad Request: empty path", http.StatusBadRequest)
		return
	}
	name := global.FileFolder + fn
	logrus.Infoln("[file.FileHandler] serve", name)
	http.ServeFile(w, r, name)
}
