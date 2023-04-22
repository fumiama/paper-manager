package global

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

const (
	UserTableUser            = "user"
	UserTableMessage         = "msg"
	UserTableMonthlyAPIVisit = "visit"
	UserTableRegex           = "re"
)

var (
	ErrInvalidRole       = errors.New("invalid role")
	ErrEmptyPassword     = errors.New("empty password")
	ErrEmptyName         = errors.New("empty name")
	ErrInvalidUsersCount = errors.New("invalid users count")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidAvatar     = errors.New("invalid avatar")
	ErrEmptyUserID       = errors.New("empty user ID")
	ErrEmptyContact      = errors.New("empty contact")
	ErrUsernameExists    = errors.New("username exists")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidContact    = errors.New("invalid contact")
	ErrNoSuchFieldName   = errors.New("no such field name")
	ErrEmptyRegex        = errors.New("empty regex")
)

func init() {
	isinit := utils.IsNotExist(UserDB.db.DBPath)
	err := UserDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
	_, err = UserDB.db.DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		panic(err)
	}
	err = UserDB.db.Create(UserTableUser, &User{})
	if err != nil {
		panic(err)
	}
	err = UserDB.db.Create(UserTableMessage, &Message{},
		"FOREIGN KEY(ToID) REFERENCES "+UserTableUser+"(ID)",
	)
	if err != nil {
		panic(err)
	}
	err = UserDB.db.Create(UserTableMonthlyAPIVisit, &MonthlyAPIVisit{})
	if err != nil {
		panic(err)
	}
	err = UserDB.db.Create(UserTableRegex, &Regex{})
	if err != nil {
		panic(err)
	}
	if isinit { // 添加初始账户
		UserDB.AddUser(&User{
			Role: RoleSuper,
			Pswd: "123456",
			Name: "fumiama",
			Nick: "源文雨",
			Avtr: "https://q1.qlogo.cn/g?b=qq&nk=1332524221&s=640",
			Cont: "028-61830156",
			Desc: "天何所沓，十二焉分。日月安属，列星安陈。",
		}, "系统")
		logrus.Warn("[user] 初次启动, 创建初始账户 fumiama 密码 123456")
	}
	err = UserDB.db.Close()
	if err != nil {
		panic(err)
	}
	err = os.Chmod(UserDB.db.DBPath, 0600)
	if err != nil {
		panic(err)
	}
	err = UserDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
}

// User stores a user in table named UserTableUser
type User struct {
	ID   *int
	Role UserRole
	Date int64 // Date is the creating date's unix timestamp
	Pswd string
	Last int64  // Last is the last password reseting unix timestamp
	Name string `db:"Name,UNIQUE"`
	Nick string
	Avtr string // Avtr is the user's avatar, typically a image url
	Cont string // Cont is the user's contact, ex. phone number
	Desc string
}

// AddUser but cannot customize the ID field for it is self-increasing
func (u *UserDatabase) AddUser(user *User, opname string) error {
	user.ID = nil
	if user.Role == RoleNil || user.Role > RoleUser {
		return ErrInvalidRole
	}
	if user.Pswd == "" {
		return ErrEmptyPassword
	}
	if user.Name == "" {
		return ErrEmptyName
	}
	if u.IsNameExists(user.Name) {
		return ErrUsernameExists
	}
	for _, c := range user.Name {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') && !(c >= 'a' && c <= 'z') {
			return ErrInvalidName
		}
	}
	user.Date = time.Now().Unix()
	user.Last = user.Date
	u.mu.Lock()
	err := u.db.InsertUnique(UserTableUser, user)
	u.mu.Unlock()
	if err != nil {
		return err
	}
	nu, err := u.GetUserByName(user.Name)
	if err != nil {
		return err
	}
	_ = u.notifyUserAdded(opname, user.Name, *nu.ID)
	return u.SendMessage(opname+" 创建了您的账号", opname, *nu.ID)
}

// UpdateUserInfo ...
func (u *UserDatabase) UpdateUserInfo(id int, opname, nick, avtr, desc string) error {
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	if nick != "" {
		user.Nick = nick
	}
	if avtr != "" {
		if strings.Contains(avtr, "..") {
			return ErrInvalidAvatar
		}
		user.Avtr = avtr
	}
	if desc != "" {
		user.Desc = desc
	}
	u.mu.Lock()
	err = u.db.Insert(UserTableUser, &user)
	u.mu.Unlock()
	if err != nil {
		return err
	}
	if opname != user.Name {
		return u.SendMessage(opname+" 更新了您的个人信息", opname, *user.ID)
	}
	return u.SendMessage("更新了个人信息", opname, *user.ID)
}

// DisableUser ...
func (u *UserDatabase) DisableUser(id int, opname string) error {
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	if opname == user.Name {
		return ErrInvalidName
	}
	user.Last = time.Now().Unix()
	user.Pswd = ""
	u.mu.Lock()
	err = u.db.Insert(UserTableUser, &user)
	u.mu.Unlock()
	if err != nil {
		return err
	}
	_ = u.SendMessage("您的账户被 "+opname+" 禁用", opname, *user.ID)
	return u.notifyDisableUser(user.Name, opname, *user.ID)
}

// UpdateUserPassword ...
func (u *UserDatabase) UpdateUserPassword(id int, opname, npwd string) error {
	if npwd == "" {
		return ErrEmptyPassword
	}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Last = time.Now().Unix()
	user.Pswd = npwd
	u.mu.Lock()
	err = u.db.Insert(UserTableUser, &user)
	u.mu.Unlock()
	if err != nil {
		return err
	}
	_ = u.notifyPasswordChange(user.Name, npwd, opname, *user.ID)
	if user.Name != opname {
		return u.SendMessage(opname+" 更新了您的密码", opname, *user.ID)
	}
	return u.SendMessage("更新了密码", opname, *user.ID)
}

// UpdateUserContact ...
func (u *UserDatabase) UpdateUserContact(id int, opname, ncont string) error {
	if ncont == "" {
		return ErrEmptyContact
	}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Cont = ncont
	u.mu.Lock()
	err = u.db.Insert(UserTableUser, &user)
	u.mu.Unlock()
	if err != nil {
		return err
	}
	_ = u.notifyContactChange(user.Name, ncont, *user.ID)
	if user.Name != opname {
		return u.SendMessage(opname+" 更新了您的联系方式", opname, *user.ID)
	}
	return u.SendMessage("更新了联系方式", opname, *user.ID)
}

// GetUserByName avoids sql injection by limiting username to 0-9A-Za-z
func (u *UserDatabase) GetUserByName(username string) (user User, err error) {
	for _, c := range username {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') && !(c >= 'a' && c <= 'z') {
			err = ErrInvalidName
			return
		}
	}
	u.mu.RLock()
	err = u.db.Find(UserTableUser, &user, "WHERE Name='"+username+"'")
	u.mu.RUnlock()
	return
}

// IsIDExists ...
func (u *UserDatabase) IsIDExists(id int) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.db.CanFind(UserTableUser, "WHERE ID="+strconv.Itoa(id))
}

// IsNameExists avoids sql injection by limiting username to 0-9A-Za-z
func (u *UserDatabase) IsNameExists(username string) bool {
	for _, c := range username {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') && !(c >= 'a' && c <= 'z') {
			return false
		}
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.db.CanFind(UserTableUser, "WHERE Name='"+username+"'")
}

// GetUserByID ...
func (u *UserDatabase) GetUserByID(id int) (user User, err error) {
	u.mu.RLock()
	err = u.db.Find(UserTableUser, &user, "WHERE ID="+strconv.Itoa(id))
	u.mu.RUnlock()
	return
}

// DelUserByID ...
func (u *UserDatabase) DelUserByID(id int) (err error) {
	u.mu.Lock()
	err = u.db.Del(UserTableUser, "WHERE ID="+strconv.Itoa(id))
	u.mu.Unlock()
	return
}

// GetUsers will set Pswd field to empty
func (u *UserDatabase) GetUsers() (users []User, err error) {
	var user User
	u.mu.RLock()
	defer u.mu.RUnlock()
	n, err := u.db.Count(UserTableUser)
	if err != nil {
		return
	}
	users = make([]User, n)
	i := 0
	err = u.db.FindFor(UserTableUser, &user, "", func() error {
		if user.Pswd != "" {
			user.Pswd = "-"
		}
		users[i] = user
		i++
		if i > n {
			return ErrInvalidUsersCount
		}
		return nil
	})
	return
}

// GetUsersCount ...
func (u *UserDatabase) GetUsersCount() (int, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.db.Count(UserTableUser)
}

func (u *UserDatabase) GetSuperIDs() (ids []int, err error) {
	var user User
	ids = make([]int, 0, 16)
	u.mu.RLock()
	defer u.mu.RUnlock()
	err = u.db.FindFor(UserTableUser, &user, "WHERE Role="+strconv.Itoa(int(RoleSuper)), func() error {
		ids = append(ids, *user.ID)
		return nil
	})
	return
}

// IsUser checks if token is valid for a user
func (user *User) IsUser() bool {
	return user.Role == RoleUser || user.Role == RoleFileManager || user.Role == RoleSuper
}

// IsFileManager checks if token is valid for a filemgr
func (user *User) IsFileManager() bool {
	return user.Role == RoleFileManager || user.Role == RoleSuper
}

// IsSuper checks if token is valid for a super
func (user *User) IsSuper() bool {
	return user.Role == RoleSuper
}
