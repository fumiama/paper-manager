package global

import "time"

func init() {
	err := FileDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
}

type File struct {
	ID *int
}

func (f *FileDatabase) AddFile() {}
