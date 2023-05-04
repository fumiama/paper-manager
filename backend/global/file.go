package global

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	sql "github.com/FloatTech/sqlite"

	"github.com/fumiama/paper-manager/backend/utils"
)

const (
	FileTableList         = "lst"
	FileTableFile         = "file"
	FileTableTempFile     = "tmpfile"
	FileTableQuestion     = "question"
	FileTableTempQuestion = "tmpqstn"
)

var (
	ErrMajorSplitsTooShort       = errors.New("major splits too short")
	ErrEmptyClass                = errors.New("empty class")
	ErrHasntAnalyzed             = errors.New("hasn't analyzed")
	ErrNoGetFileStatusPermission = errors.New("no get file status permission")
)

func init() {
	err := FileDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableList, &List{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableFile, &File{},
		"FOREIGN KEY(ListID) REFERENCES "+FileTableList+"(ID)",
	)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableTempFile, &File{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableQuestion, &Question{},
		"FOREIGN KEY(ListID) REFERENCES "+FileTableList+"(ID)",
	)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableTempQuestion, &Question{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Close()
	if err != nil {
		panic(err)
	}
	err = os.Chmod(FileDB.db.DBPath, 0600)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
}

// File stores to paper/Class/2022-2023学年/第一学期/期末/A/xxx.docx
type File struct {
	ID        int64 // ID is the first 8 bytes of the original file's md5
	ListID    int   // ListID is the foreign key to List(ID)
	Year      StudyYear
	Type      PaperType
	Date      uint32        // Date is the yyyymmdd of 考试日期
	Time      time.Duration // Time is 考试时长
	Class     string        // Class is 考试科目
	Rate      string        // Rate is 成绩构成比例
	Questions []byte        // Questions is for []QuestionJSON
}

// StudyYear 学年
type StudyYear uint16

// String ex. 2022-2023学年
func (sy StudyYear) String() string {
	next := sy + 1
	return strconv.Itoa(int(sy)) + "-" + strconv.Itoa(int(next)) + "学年"
}

func (file *File) GetList(f *FileDatabase) (List, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(file.ListID))
}

// DelFile by listid
func (f *FileDatabase) DelFile(lstid, uid int, istemp bool) error {
	user, err := UserDB.GetUserByID(uid)
	if err != nil {
		return err
	}
	if !user.IsSuper() && !istemp {
		return ErrInvalidRole
	}
	ftable := ""
	if istemp {
		ftable = FileTableTempFile
	} else {
		ftable = FileTableFile
	}
	f.mu.RLock()
	lst, err := sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	f.mu.RUnlock()
	if err != nil {
		return err
	}
	if istemp && lst.Uploader != uid {
		return ErrInvalidRole
	}
	if lst.Path == "" || strings.Contains(lst.Path, "..") {
		return os.ErrNotExist
	}
	err = f.db.Del(FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	if err != nil {
		return err
	}
	if lst.HasntAnalyzed {
		return os.RemoveAll(lst.Path)
	}
	i := strings.LastIndex(lst.Path, "/")
	if i <= 0 {
		return os.ErrNotExist
	}
	parentfolder := lst.Path[:i]
	if utils.IsNotExist(parentfolder) {
		return os.ErrNotExist
	}
	f.mu.RLock()
	file, err := sql.Find[File](&f.db, ftable, "WHERE ListID="+strconv.Itoa(lstid))
	f.mu.RUnlock()
	if err != nil {
		return err
	}
	err = f.db.Del(ftable, "WHERE ListID="+strconv.Itoa(lstid))
	if err != nil {
		return err
	}
	ques := make([]QuestionJSON, 0, 64)
	err = json.Unmarshal(file.Questions, &ques)
	if err != nil {
		return err
	}
	for _, q := range ques {
		q.Delete(f, istemp)
	}
	return os.RemoveAll(parentfolder)
}

// GetFile get analyzed file's structure from List(ID)
func (f *FileDatabase) GetFile(lstid, uid int) (file *File, lst List, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	lst, err = sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	if err != nil {
		return
	}
	if lst.HasntAnalyzed {
		err = ErrHasntAnalyzed
		return
	}
	if lst.IsTemp && lst.Uploader != uid {
		err = ErrNoGetFileStatusPermission
		return
	}
	ftable := ""
	if lst.IsTemp {
		ftable = FileTableTempFile
	} else {
		ftable = FileTableFile
	}
	filestruct, err := sql.Find[File](&f.db, ftable, "WHERE ListID="+strconv.Itoa(lstid))
	if err != nil {
		return
	}
	return &filestruct, lst, nil
}

// GetFilesByYearRange ...
func (f *FileDatabase) GetFilesByYearRange(yearstart, yearend StudyYear) ([]*File, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return sql.FindAll[File](&f.db, FileTableFile, "WHERE Year>="+strconv.Itoa(int(yearstart))+" AND Year<="+strconv.Itoa(int(yearend)))
}
