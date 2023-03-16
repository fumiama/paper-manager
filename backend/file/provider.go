package file

import (
	"net/http"
	"strings"

	"github.com/fumiama/paper-manager/backend/global"
	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

// Handler serves contents in global.FileFolder
func Handler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod("GET", w, r) {
		return
	}
	i := strings.LastIndex(r.URL.Path, "/")
	fn := r.URL.Path[i+1:]
	if fn == "" {
		http.Error(w, "400 Bad Request: empty path", http.StatusBadRequest)
		return
	}
	name := global.FileFolder + fn
	logrus.Infoln("[file.Handler]\t serve", name)
	http.ServeFile(w, r, name)
}
