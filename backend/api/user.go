package api

import (
	"errors"
	"strings"
	"time"

	"github.com/fumiama/paper-manager/backend/global"
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
	Date     string `json:"date"`
	Last     string `json:"last"`
	Contact  string `json:"contact"`
}

func getUserInfo(token string) (*getUserInfoResult, error) {
	user := usertokens.Get(token)
	if user == nil {
		return nil, errInvalidToken
	}
	cont := user.Cont
	if len(cont) > 7 {
		sb := strings.Builder{}
		sb.WriteString(cont[:3])
		for i := 0; i < len(cont)-7; i++ {
			sb.WriteByte('*')
		}
		sb.WriteString(cont[len(cont)-4:])
		cont = sb.String()
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
		Date:    time.Unix(user.Date, 0).Format(chineseDateLayout),
		Last:    time.Unix(user.Last, 0).Format(chineseDateLayout),
		Contact: cont,
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
