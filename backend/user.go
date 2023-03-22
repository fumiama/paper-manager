package backend

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	base14 "github.com/fumiama/go-base16384"
	"github.com/fumiama/paper-manager/backend/global"
	"github.com/fumiama/paper-manager/backend/utils"
)

const (
	chineseDateLayout = "2006年01月02日15时04分05秒"
)

var (
	errInvalidToken          = errors.New("invalid token")
	errNoListUsersPermission = errors.New("no list users permission")
	errNoSetRolePermission   = errors.New("no set role permission")
	errInvalidRole           = errors.New("invalid role")
)

type getUserInfoResult struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Avatar   string `json:"avatar"`
	Desc     string `json:"desc"`
	HomePath string `json:"homePath"`
	Roles    []role `json:"roles"`
	Date     int64  `json:"date"`
	Last     string `json:"last"`
	Contact  string `json:"contact"`
}

func hideContact(cont string) string {
	if len(cont) > 7 {
		sb := strings.Builder{}
		sb.WriteString(cont[:3])
		for i := 0; i < len(cont)-7; i++ {
			sb.WriteByte('*')
		}
		sb.WriteString(cont[len(cont)-4:])
		return sb.String()
	}
	return cont
}

func getUserInfo(token string) (*getUserInfoResult, error) {
	user := usertokens.Get(token)
	if user == nil {
		return nil, errInvalidToken
	}
	return &getUserInfoResult{
		UserID:   *user.ID,
		Username: user.Name,
		RealName: user.Nick,
		Avatar:   user.Avtr,
		Desc:     user.Desc,
		HomePath: func() string {
			if user.Role == global.RoleSuper {
				return "/dashboard/analysis"
			}
			return "/dashboard/workbench"
		}(),
		Roles:   []role{{RoleName: user.Role.Nick(), Value: user.Role.String()}},
		Date:    user.Date,
		Last:    time.Unix(user.Last, 0).Format(chineseDateLayout),
		Contact: hideContact(user.Cont),
	}, nil
}

func getUsersCount(token string) (int, error) {
	user := usertokens.Get(token)
	if user == nil {
		return 0, errInvalidToken
	}
	return global.UserDB.GetUsersCount()
}

type getUsersListResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Nick string `json:"nick"`
	Stat bool   `json:"stat"`
	Role string `json:"role"`
	Date string `json:"date"`
	Desc string `json:"desc"`
}

func getUsersList(token string) ([]getUsersListResult, error) {
	user := usertokens.Get(token)
	if user == nil {
		return nil, errInvalidToken
	}
	if !user.IsSuper() {
		return nil, errNoListUsersPermission
	}
	us, err := global.UserDB.GetUsers()
	if err != nil {
		return nil, err
	}
	ret := make([]getUsersListResult, len(us))
	for i, u := range us {
		ret[i].ID = *u.ID
		ret[i].Name = u.Name
		ret[i].Nick = u.Nick
		ret[i].Stat = u.Pswd != ""
		ret[i].Role = u.Role.Nick()
		ret[i].Date = time.Unix(u.Date, 0).Format(chineseDateLayout)
		ret[i].Desc = u.Desc
	}
	return ret, nil
}

func isNameExist(token, name string) (bool, error) {
	user := usertokens.Get(token)
	if user == nil {
		return false, errInvalidToken
	}
	if !user.IsSuper() {
		return false, errNoListUsersPermission
	}
	return global.UserDB.IsNameExists(name), nil
}

func setUserPassword(id int, token, npwd string) error {
	user, err := global.UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	h := md5.New()
	h.Write(base14.StringToBytes(user.Pswd))
	h.Write(base14.StringToBytes(npwd))
	if token != hex.EncodeToString(h.Sum(make([]byte, 0, 16))) {
		return errInvalidToken
	}
	return global.UserDB.UpdateUserPassword(id, user.Name, npwd)
}

func setUserContact(id int, token, ncont string) error {
	user, err := global.UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	h := md5.New()
	h.Write(base14.StringToBytes(user.Cont))
	h.Write(base14.StringToBytes(ncont))
	if token != hex.EncodeToString(h.Sum(make([]byte, 0, 16))) {
		return errInvalidToken
	}
	return global.UserDB.UpdateUserContact(id, user.Name, ncont)
}

// setUserInfo may change the arguments
func setUserInfo(id int, nick, desc, avtr *string) error {
	user, err := global.UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	n, d, a := *nick, *desc, *avtr
	if n == "" {
		*nick = user.Nick
	}
	if n == user.Nick {
		n = ""
	}
	if d == "" {
		*desc = user.Desc
	}
	if d == user.Desc {
		d = ""
	}
	if a == "" {
		*avtr = user.Avtr
	} else if utils.IsNotExist(global.DataFolder + a) {
		return os.ErrNotExist
	}
	if a == user.Avtr {
		a = ""
	}
	return global.UserDB.UpdateUserInfo(id, user.Name, n, a, d)
}

// setOthersInfo may change the arguments
func setOthersInfo(id int, opname, nick, desc string) error {
	user, err := global.UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	if nick == user.Nick {
		nick = ""
	} else if nick != "" {
		user.Nick = nick
	}
	if desc == user.Desc {
		desc = ""
	} else if desc != "" {
		user.Desc = desc
	}
	return global.UserDB.UpdateUserInfo(id, opname, nick, "", desc)
}

func setUserRole(id int, role global.UserRole, opname string) error {
	if !role.IsVaild() {
		return errInvalidRole
	}
	user, err := global.UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	if role == user.Role {
		return nil
	}
	return global.UserDB.UpdateUserRole(*user.ID, role, opname)
}

func resetPassword(ip, name, mobile string) error {
	if registerlimit.Get(ip) {
		return errRequestTooFast
	}
	if ip == "" {
		return errInvalidIP
	}
	registerlimit.Set(ip, true)
	return global.UserDB.NotifyResetPassword(ip, name, mobile)
}

func init() {
	apimap["/api/getUserInfo"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		ret, err := getUserInfo(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, ret, messageOk, typeSuccess)
	}}

	apimap["/api/getUsersCount"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		n, err := getUsersCount(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, n, messageOk, typeSuccess)
	}}

	apimap["/api/getUsersList"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		ret, err := getUsersList(token)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, &ret, messageOk, typeSuccess)
	}}

	apimap["/api/isNameExist"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		name := r.URL.Query().Get("username")
		if name == "" {
			writeresult(w, codeError, nil, "empty username", typeError)
			return
		}
		yes, err := isNameExist(token, name)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, yes, messageOk, typeSuccess)
	}}

	apimap["/api/setPassword"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setpasswordbody struct {
			Token    string `json:"token"`
			Password string `json:"password"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body setpasswordbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		err = setUserPassword(*user.ID, body.Token, body.Password)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "成功, 请重新登录"}, messageOk, typeSuccess)
		_ = logout(token)
	}}

	apimap["/api/setContact"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setcontactbody struct {
			Token   string `json:"token"`
			Contact string `json:"contact"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body setcontactbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		err = setUserContact(*user.ID, body.Token, body.Contact)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		user.Cont = hideContact(body.Contact)
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "成功, 已将消息报告给课程组长"}, messageOk, typeSuccess)
	}}

	apimap["/api/setUserInfo"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setuserinfobody struct {
			ID   int    `json:"id"`
			Nick string `json:"nick"`
			Desc string `json:"desc"`
			Avtr string `json:"avtr"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body setuserinfobody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		type message struct {
			M string `json:"msg"`
		}
		if body.ID != 0 {
			if !user.IsSuper() {
				writeresult(w, codeError, nil, "no permission to set others' info", typeError)
				return
			}
			err = setOthersInfo(body.ID, user.Name, body.Nick, body.Desc)
			if err != nil {
				writeresult(w, codeError, nil, err.Error(), typeError)
				return
			}
			writeresult(w, codeSuccess, &message{M: "成功"}, messageOk, typeSuccess)
			return
		}
		err = setUserInfo(*user.ID, &body.Nick, &body.Desc, &body.Avtr)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		user.Nick = body.Nick
		user.Desc = body.Desc
		user.Avtr = body.Avtr
		writeresult(w, codeSuccess, &message{M: "成功"}, messageOk, typeSuccess)
	}}

	apimap["/api/setRole"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type setrolebody struct {
			ID   int             `json:"id"`
			Role global.UserRole `json:"role"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		if !user.IsSuper() {
			writeresult(w, codeError, nil, errNoSetRolePermission.Error(), typeError)
			return
		}
		var body setrolebody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		if body.ID == *user.ID {
			writeresult(w, codeError, nil, "cannot set self", typeError)
			return
		}
		err = setUserRole(body.ID, body.Role, user.Name)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, nil, messageOk, typeSuccess)
	}}

	apimap["/api/disableUser"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type disableuserbody struct {
			ID int `json:"id"`
		}
		token := r.Header.Get("Authorization")
		user := usertokens.Get(token)
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		if !user.IsSuper() {
			writeresult(w, codeError, nil, errNoSetRolePermission.Error(), typeError)
			return
		}
		var body disableuserbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		if body.ID == *user.ID {
			writeresult(w, codeError, nil, "cannot disbale self", typeError)
			return
		}
		err = global.UserDB.DisableUser(body.ID, user.Name)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		writeresult(w, codeSuccess, nil, messageOk, typeSuccess)
	}}

	apimap["/api/resetPassword"] = &apihandler{"POST", func(w http.ResponseWriter, r *http.Request) {
		type resetpwdbody struct {
			Username string `json:"username"`
			Mobile   string `json:"mobile"`
		}
		if r.Header.Get("Authorization") != "" {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		var body resetpwdbody
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		ip := r.RemoteAddr
		i := strings.LastIndex(ip, ":")
		if i >= 0 {
			ip = ip[:i]
		}
		err = resetPassword(ip, body.Username, body.Mobile)
		if err != nil {
			writeresult(w, codeError, nil, err.Error(), typeError)
			return
		}
		type message struct {
			M string `json:"msg"`
		}
		writeresult(w, codeSuccess, &message{M: "已上报, 请耐心等待通知"}, messageOk, typeSuccess)
	}}
}
