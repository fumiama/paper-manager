package backend

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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
	errInvalidToken = errors.New("invalid token")
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

func logout(token string) error {
	user := usertokens.Get(token)
	if user == nil {
		return errInvalidToken
	}
	loginstatus.Delete(user.Name)
	usertokens.Delete(token)
	return nil
}

func getUsersCount(token string) (int, error) {
	user := usertokens.Get(token)
	if user == nil {
		return 0, errInvalidToken
	}
	return global.UserDB.GetUsersCount()
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
	return global.UserDB.UpdateUserPassword(id, npwd)
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
	return global.UserDB.UpdateUserContact(id, ncont)
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
	return global.UserDB.UpdateUserInfo(id, n, a, d)
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
