package backend

import (
	"net/http"

	"github.com/fumiama/paper-manager/backend/global"
	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

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
