package api

import (
	"errors"

	"github.com/fumiama/paper-manager/backend/global"
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
		Roles: []role{{RoleName: user.Role.Nick(), Value: user.Role.String()}},
	}, nil
}
