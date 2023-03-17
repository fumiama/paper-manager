package global

import (
	"sync"

	sql "github.com/FloatTech/sqlite"
)

const (
	userdbpath = DataFolder + "user.db"
	filedbpath = DataFolder + "file.db"
)

type UserDatabase struct {
	mu sync.RWMutex
	db sql.Sqlite
}
type FileDatabase struct {
	mu sync.RWMutex
	db sql.Sqlite
}

var (
	UserDB = UserDatabase{db: sql.Sqlite{DBPath: userdbpath}}
	FileDB = FileDatabase{db: sql.Sqlite{DBPath: filedbpath}}
)
