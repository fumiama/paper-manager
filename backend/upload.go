package backend

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fumiama/imgsz"
	"github.com/sirupsen/logrus"

	"github.com/fumiama/paper-manager/backend/global"
	"github.com/fumiama/paper-manager/backend/utils"
)

type upload struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	URL     string `json:"url"`
}

// UploadHandler receives uploaded files
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMethod("POST", w, r) {
		return
	}
	token := r.Header.Get("Authorization")
	user := usertokens.Get(token)
	if user == nil {
		writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
		return
	}
	global.UserDB.VisitAPI()
	ff, h, err := r.FormFile("avatar")
	if err == nil {
		defer ff.Close()
		ct := h.Header.Get("Content-Type")
		un := h.Filename
		logrus.Infoln("[file.UploadHandler] receive avatar, username:", un, "& mime:", ct)
		if !strings.HasPrefix(ct, "image/") {
			writeresult(w, codeError, nil, "invalid mimetype", typeError)
			return
		}
		if un != user.Name {
			writeresult(w, codeError, nil, "username mismatch", typeError)
			return
		}
		err = os.MkdirAll(global.FileFolder+strconv.Itoa(*user.ID), 0755)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		buf := bytes.NewBuffer(make([]byte, 0, h.Size))
		_, err := io.Copy(buf, ff)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		data := buf.Bytes()
		_, format, err := imgsz.DecodeSize(buf)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		userf := global.FileFolder + strconv.Itoa(*user.ID) + "/"
		err = os.MkdirAll(userf, 0755)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		avf := userf + "avatar" + time.Now().Format("_20060102_15_04_05") + "." + format
		err = os.WriteFile(avf, data, 0644)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		/*err = global.UserDB.UpdateUserInfo(*user.ID, "", avf[6:], "")
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		user.Avtr = avf[6:]
		usertokens.Set(token, user)*/
		writeresult(w, codeSuccess, &upload{
			Message: messageOk,
			Code:    codeSuccess,
			URL:     avf[6:],
		}, messageOk, typeSuccess)
		logrus.Infoln("[file.UploadHandler] save avatar to", avf[6:])
		return
	}
	if err != http.ErrMissingFile {
		writeresult(w, codeError, nil, err.Error(), typeError)
		return
	}
	ff, h, err = r.FormFile("paper")
	if err == nil {
		defer ff.Close()
		ct := h.Header.Get("Content-Type")
		fn := h.Filename
		logrus.Infoln("[file.UploadHandler] receive paper, name:", fn)
		if ct != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
			writeresult(w, codeError, nil, "invalid mimetype: need docx", typeError)
			return
		}
		if strings.ContainsAny(fn, `/\`) || strings.Contains(fn, "..") {
			writeresult(w, codeError, nil, "invalid filename", typeError)
			return
		}
		id, err := global.FileDB.SaveFileToTemp(*user.ID, ff, fn)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, id, messageOk, typeSuccess)
		return
	}
	if err != http.ErrMissingFile {
		writeresult(w, codeError, nil, err.Error(), typeError)
		return
	}
}
