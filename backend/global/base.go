package global

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	// DataFolder stores all backend data in
	DataFolder = "./data/"
	// FileFolder stores all blob files
	FileFolder = DataFolder + "file/"
)

func init() {
	initdir(DataFolder)
	initdir(FileFolder)
}

func initdir(folder string) {
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		logrus.Errorln("[os.MkdirAll]", err)
		os.Exit(line())
	}
}

func line() int {
	_, _, fileLine, ok := runtime.Caller(2)
	if ok {
		return fileLine
	}
	return -1
}
