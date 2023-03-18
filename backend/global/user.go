package global

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

const (
	RoleNil UserRole = iota
	RoleSuper
	RoleFileManager
	RoleUser
)

type UserRole uint8

func (r UserRole) String() string {
	switch r {
	case RoleSuper:
		return "super"
	case RoleFileManager:
		return "filemgr"
	case RoleUser:
		return "user"
	}
	return "nil"
}

func (r UserRole) Nick() string {
	switch r {
	case RoleSuper:
		return "课程组长"
	case RoleFileManager:
		return "归档代理"
	case RoleUser:
		return "课程组员"
	}
	return "nil"
}

const (
	UserTableUser    = "user"
	UserTableMessage = "msg"
)

var (
	ErrInvalidRole       = errors.New("invalid role")
	ErrEmptyPassword     = errors.New("empty password")
	ErrEmptyName         = errors.New("empty name")
	ErrInvalidUsersCount = errors.New("invalid users count")
	ErrEmptyUserID       = errors.New("empty user ID")
	ErrEmptyContect      = errors.New("empty contact")
	ErrUsernameExists    = errors.New("username exists")
	ErrInvalidName       = errors.New("invalid name")
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
}

// User stores a user in table named UserTableUser
type User struct {
	ID   *int
	Role UserRole
	Date int64 // Date is the creating date's unix timestamp
	Pswd string
	Last int64 // Last is the last password reseting unix timestamp
	Name string
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
	_ = u.notifyUserAdded(opname, user.Name)
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.InsertUnique(UserTableUser, user)
}

// UpdateUserInfo ...
func (u *UserDatabase) UpdateUserInfo(id int, nick, avtr, desc string) error {
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	if nick != "" {
		user.Nick = nick
	}
	if avtr != "" {
		user.Avtr = avtr
	}
	if desc != "" {
		user.Desc = desc
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.Insert(UserTableUser, &user)
}

// UpdateUserRole ...
func (u *UserDatabase) UpdateUserRole(id int, nr UserRole) error {
	if nr == RoleNil || nr > RoleUser {
		return ErrInvalidRole
	}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Role = nr
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.Insert(UserTableUser, &user)
}

// UpdateUserPassword ...
func (u *UserDatabase) UpdateUserPassword(id int, npwd string) error {
	if npwd == "" {
		return ErrEmptyPassword
	}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Last = time.Now().Unix()
	user.Pswd = npwd
	_ = u.notifyPasswordChange(user.Name, npwd)
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.Insert(UserTableUser, &user)
}

// UpdateUserContact ...
func (u *UserDatabase) UpdateUserContact(id int, ncont string) error {
	if ncont == "" {
		return ErrEmptyContect
	}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Cont = ncont
	_ = u.notifyContactChange(user.Name, ncont)
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.Insert(UserTableUser, &user)
}

// GetUserByName avoids sql injection by removing ; ' " =
func (u *UserDatabase) GetUserByName(username string) (user User, err error) {
	username = strings.NewReplacer(";", "", "'", "", `"`, "", "=", "").Replace(username)
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

// IsNameExists avoids sql injection by removing ; ' " =
func (u *UserDatabase) IsNameExists(username string) bool {
	username = strings.NewReplacer(";", "", "'", "", `"`, "", "=", "").Replace(username)
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
		user.Pswd = ""
		users[i] = user
		i++
		if i >= n {
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

// Message is shown in the workbench
type Message struct {
	ID   *int
	ToID int // ToID user's ID
	Date int64
	Text string // Text is the message content
	Name string // Name is the user's name to add in register message
	Cont string // Cont is the user's phone number to add in register message or an operator's name in add user message
	Pswd string // Pswd is the user's password to add in register message
}

// SendMessage will send a message
func (u *UserDatabase) SendMessage(m *Message) error {
	m.ID = nil
	m.Date = time.Now().Unix()
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.InsertUnique(UserTableMessage, m)
}

// NotifyRegister will send register notification to all supers
func (u *UserDatabase) NotifyRegister(name, cont, pswd string) error {
	if name == "" {
		return ErrEmptyName
	}
	for _, c := range name {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') && !(c >= 'a' && c <= 'z') {
			return ErrInvalidName
		}
	}
	if pswd == "" {
		return ErrEmptyPassword
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: "收到来自 " + name + " 的注册请求, 联系方式: " + cont,
		Name: name,
		Cont: cont,
		Pswd: pswd,
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyUserAdded will send notification to all supers
func (u *UserDatabase) notifyUserAdded(opname, name string) error {
	if opname == "" || name == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: opname + "添加了用户 " + name,
		Name: name,
		Cont: opname,
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyContactChange will send notification to all supers
func (u *UserDatabase) notifyContactChange(name, cont string) error {
	if name == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: "用户 " + name + " 更改联系方式为: " + cont,
		Name: name,
		Cont: cont,
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyPasswordChange will send notification to all supers
func (u *UserDatabase) notifyPasswordChange(name, npwd string) error {
	if name == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: "用户 " + name + " 更改了密码",
		Name: name,
		Pswd: npwd,
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMessagesOfUser set Pswd field to empty
func (u *UserDatabase) GetMessagesOfUser(to int) (ms []Message, err error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	n, err := u.db.Count(UserTableMessage)
	if err != nil {
		return
	}
	ms = make([]Message, 0, n)
	m := Message{}
	err = u.db.FindFor(UserTableMessage, &m, "WHERE ToID="+strconv.Itoa(to), func() error {
		m.Pswd = ""
		ms = append(ms, m)
		return nil
	})
	return
}

// GetMessageByID ...
func (u *UserDatabase) GetMessageByID(id int) (m Message, err error) {
	u.mu.RLock()
	err = u.db.Find(UserTableMessage, &m, "WHERE ID="+strconv.Itoa(id))
	u.mu.RUnlock()
	return
}

// DelMessageByID ...
func (u *UserDatabase) DelMessageByID(id int) (err error) {
	u.mu.Lock()
	err = u.db.Del(UserTableMessage, "WHERE ID="+strconv.Itoa(id))
	u.mu.Unlock()
	return
}
