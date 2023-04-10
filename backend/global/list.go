package global

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	sql "github.com/FloatTech/sqlite"
)

var (
	ErrInvalidSaveName = errors.New("invalid save name")
)

// List of file path
type List struct {
	ID            *int   // ID is self-inc
	Uploader      int    // Uploader is uid
	UpName        string // UpName is uploader's name
	UpTime        int64  // UpTime is upload time (unix timestamp)
	Size          int64  // Size of the original file
	QuesC         int    // QuesC 总小题数
	HasntAnalyzed bool   // HasntAnalyzed whether file has been analyzed
	IsTemp        bool   // IsTemp whether file is temp
	Path          string // Path of file, normally unique
	Desc          string // Desc is file's description
}

// SaveFileToTemp copy file to PaperFolder/tmp/uploader/name and add record into list.
func (f *FileDatabase) SaveFileToTemp(uploader int, file io.Reader, name string) (id int, err error) {
	user, err := UserDB.GetUserByID(uploader)
	if err != nil {
		return
	}
	if strings.Contains(name, "..") || strings.Contains(name, "/") {
		err = ErrInvalidSaveName
		return
	}
	tmpdir := PaperFolder + "tmp/" + strconv.Itoa(uploader)
	err = os.MkdirAll(tmpdir, 0755)
	if err != nil {
		return
	}
	fpath := tmpdir + "/" + name
	f.mu.RLock()
	lst, _ := sql.Find[List](&f.db, FileTableList, "WHERE Path='"+fpath+"'")
	f.mu.RUnlock()
	lst.Uploader = uploader
	lst.UpName = user.Name
	lst.UpTime = time.Now().Unix()
	lst.HasntAnalyzed = true
	lst.IsTemp = true
	lst.Path = fpath
	ff, err := os.Create(lst.Path)
	if err != nil {
		return
	}
	sz, err := io.Copy(ff, file)
	_ = ff.Close()
	if err != nil {
		_ = os.Remove(lst.Path)
		return
	}
	lst.Size = sz
	f.mu.Lock()
	err = f.db.Insert(FileTableList, &lst)
	f.mu.Unlock()
	if err != nil {
		return
	}
	if lst.ID != nil {
		id = *lst.ID
		return
	}
	f.mu.RLock()
	err = f.db.Find(FileTableList, &lst, "WHERE Path='"+fpath+"'")
	f.mu.RUnlock()
	id = *lst.ID
	return
}

// ListUploadedFile will select all file that HasntAnalyzed && IsTemp or !HasntAnalyzed && !IsTemp
func (f *FileDatabase) ListUploadedFile() (lst []*List, err error) {
	f.mu.RLock()
	lst, err = sql.FindAll[List](&f.db, FileTableList, "WHERE (HasntAnalyzed AND IsTemp) OR (NOT HasntAnalyzed AND NOT IsTemp) ORDER BY UpTime DESC")
	f.mu.RUnlock()
	return
}

func (f *FileDatabase) GetFileInfo(id int) (lst List, err error) {
	f.mu.RLock()
	lst, err = sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(id))
	f.mu.RUnlock()
	return
}
