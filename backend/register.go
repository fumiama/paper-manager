package backend

import (
	"errors"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/fumiama/paper-manager/backend/global"
)

var registerlimit = ttl.NewCache[string, bool](time.Minute * 10)

var (
	errRegisterTooFast = errors.New("register too fast")
	errInvalidIP       = errors.New("invalid IP")
)

func register(ip, name, mobile, npwd string) error {
	if registerlimit.Get(ip) {
		return errRegisterTooFast
	}
	if ip == "" {
		return errInvalidIP
	}
	registerlimit.Set(ip, true)
	return global.UserDB.NotifyRegister(ip, name, mobile, npwd)
}
