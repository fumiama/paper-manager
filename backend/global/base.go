package global

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	// DataFolder stores all backend data in
	DataFolder = "./data/"
	// FileFolder stores all uploaded blob files
	FileFolder = DataFolder + "file/"
	// ImageFolder stores images of questions
	ImageFolder = DataFolder + "image/"
	// PaperFolder stores all protected files
	PaperFolder = DataFolder + "paper/"
)

func init() {
	initdir(DataFolder)
	initdir(FileFolder)
	initsecuredir(ImageFolder)
	initsecuredir(PaperFolder)
}

func initdir(folder string) {
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		logrus.Errorln("[os.MkdirAll]", err)
		os.Exit(line())
	}
}

func initsecuredir(folder string) {
	err := os.MkdirAll(folder, 0700)
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
