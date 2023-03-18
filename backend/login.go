package backend

import (
	"crypto/md5"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/RomiChan/syncx"
	base14 "github.com/fumiama/go-base16384"

	"github.com/fumiama/paper-manager/backend/global"
)

var (
	errNoSuchUser          = errors.New("invalid username or password")
	errTooManySalts        = errors.New("too many salts")
	errInvalidLoginStatus  = errors.New("invalid login status")
	errEmptySalt           = errors.New("empty salt")
	errWrongPassword       = errors.New("invalid username or password")
	errTooManyFailedLogins = errors.New("too many failed logins")
)

const (
	loginStatusYes = iota - 1
	loginStatusNo
	loginStatusFail1
	loginStatusFail2
	loginStatusFail3
	loginStatusFailLast
)

const maxSaltCount = 4

type saltinfo struct {
	Salt  string `json:"salt"`
	count *uintptr
}

var (
	loginsalts  = ttl.NewCache[string, saltinfo](time.Minute)
	loginstatus = syncx.Map[string, int]{}
)

func getLoginSalt(username string) (*saltinfo, error) {
	if !global.UserDB.IsNameExists(username) {
		return nil, errNoSuchUser
	}
	s, _ := loginstatus.Load(username)
	if s == loginStatusYes {
		return nil, errInvalidLoginStatus
	}
	if s >= loginStatusFailLast {
		return nil, errTooManyFailedLogins
	}
	salt := loginsalts.Get(username)
	if salt.count != nil {
		if atomic.AddUintptr(salt.count, 1) >= maxSaltCount {
			time.AfterFunc(time.Minute*2, func() { atomic.StoreUintptr(salt.count, 0) })
			return nil, errTooManySalts
		}
		if salt.Salt != "" {
			return &salt, nil
		}
	}
	buf := make([]byte, 7*(rand.Intn(8)+1))
	_, err := crand.Read(buf)
	if err != nil {
		return nil, err
	}
	salt.Salt = base14.EncodeToString(buf)
	salt.count = new(uintptr)
	loginsalts.Set(username, salt)
	return &salt, nil
}

type role struct {
	RoleName string `json:"roleName"`
	Value    string `json:"value"`
}

type loginResult struct {
	Roles    []role `json:"roles"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	Token    string `json:"token"`
	RealName string `json:"realName"`
	Desc     string `json:"desc"`
}

var (
	usertokens = ttl.NewCache[string, *global.User](time.Hour)
)

func login(username, challenge string) (*loginResult, error) {
	if !global.UserDB.IsNameExists(username) {
		return nil, errNoSuchUser
	}
	s, loaded := loginstatus.LoadOrStore(username, loginStatusFail1)
	if loaded {
		if s == loginStatusYes {
			return nil, errInvalidLoginStatus
		}
		if s >= loginStatusFailLast {
			return nil, errTooManyFailedLogins
		}
		loginstatus.Store(username, s+1)
	}
	salt := loginsalts.Get(username)
	if salt.count == nil || salt.Salt == "" {
		return nil, errEmptySalt
	}
	user, err := global.UserDB.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	h.Write(base14.StringToBytes(user.Pswd))
	h.Write(base14.StringToBytes(salt.Salt))
	passchlg := hex.EncodeToString(h.Sum(make([]byte, 0, md5.Size)))
	if passchlg != challenge {
		return nil, errWrongPassword
	}
	var buf [6 * 8]byte
	_, err = crand.Read(buf[:])
	if err != nil {
		return nil, err
	}
	token := base64.RawStdEncoding.EncodeToString(buf[:])
	usertokens.Set(token, &user)
	loginstatus.Store(username, loginStatusYes)
	return &loginResult{
		Roles:    []role{{RoleName: user.Role.Nick(), Value: user.Role.String()}},
		UserID:   *user.ID,
		Username: user.Name,
		Token:    token,
		RealName: user.Nick,
		Desc:     user.Desc,
	}, nil
}
