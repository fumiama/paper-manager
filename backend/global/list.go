package global

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidSaveName = errors.New("invalid save name")
)

// List of file path
type List struct {
	ID       *int   // ID is self-inc
	Uploader int    // Uploader is uid
	UpTime   int64  // UpTime is upload time (unix timestamp)
	Size     int64  // Size of the original file
	IsTemp   bool   // IsTemp whether file is temp
	Path     string // Path of file
}

// SaveFileToTemp copy file to PaperFolder/tmp/uploader/name and add record into list.
func (f *FileDatabase) SaveFileToTemp(uploader int, file io.Reader, name string) (err error) {
	_, err = UserDB.GetUserByID(uploader)
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
	lst := List{
		Uploader: uploader,
		UpTime:   time.Now().Unix(),
		IsTemp:   true,
		Path:     tmpdir + "/" + name,
	}
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
	FileDB.mu.Lock()
	err = FileDB.db.Insert(FileTableList, &lst)
	FileDB.mu.Unlock()
	return
}
