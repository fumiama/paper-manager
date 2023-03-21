package global

import (
	"errors"
	"strconv"
	"time"

	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

const (
	RoleNil UserRole = iota
	RoleSuper
	RoleFileManager
	RoleUser
	RoleTop
)

type UserRole uint8

func (r UserRole) IsVaild() bool {
	return r > RoleNil && r < RoleTop
}

func (r UserRole) String() string {
	switch r {
	case RoleSuper:
		return "super"
	case RoleFileManager:
		return "filemgr"
	case RoleUser:
		return "user"
	}
	return "invalid"
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
	return "非法角色"
}

const (
	MessageNormal MessageType = iota
	MessageRegister
	MessageUserAdded
	MessageContactChange
	MessagePasswordChange
	MessageResetPassword
	MessageOperator
)

type MessageType uint8

const (
	UserTableUser            = "user"
	UserTableMessage         = "msg"
	UserTableMonthlyAPIVisit = "visit"
)

var (
	ErrInvalidRole       = errors.New("invalid role")
	ErrEmptyPassword     = errors.New("empty password")
	ErrEmptyName         = errors.New("empty name")
	ErrInvalidUsersCount = errors.New("invalid users count")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrEmptyUserID       = errors.New("empty user ID")
	ErrEmptyContact      = errors.New("empty contact")
	ErrUsernameExists    = errors.New("username exists")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidContact    = errors.New("invalid contact")
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

// UpdateUserRole ...
func (u *UserDatabase) UpdateUserRole(id int, nr UserRole, opname string) error {
	if nr == RoleNil || nr > RoleUser {
		return ErrInvalidRole
	}
	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	if opname == user.Name {
		return ErrInvalidName
	}
	user.Role = nr
	u.mu.Lock()
	err = u.db.Insert(UserTableUser, &user)
	u.mu.Unlock()
	if err != nil {
		return err
	}
	_ = u.SendMessage("您的权限被 "+opname+" 变更为 "+user.Role.Nick(), opname, *user.ID)
	return u.notifyUpdateUserRole(user.Name, opname, nr, *user.ID)
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

// Type decide message type by fields Name, Cont and Pswd.
func (m *Message) Type() MessageType {
	switch {
	case m.Name != "" && m.Cont != "" && m.Pswd != "":
		return MessageRegister
	case m.Name == "" && m.Cont != "" && m.Pswd == "":
		return MessageUserAdded
	case m.Name != "" && m.Cont != "" && m.Pswd == "":
		return MessageContactChange
	case m.Name != "" && m.Cont == "" && m.Pswd != "":
		return MessagePasswordChange
	case m.Name != "" && m.Cont == "" && m.Pswd == "":
		return MessageResetPassword
	case m.Name == "" && m.Cont != "" && m.Pswd != "":
		return MessageOperator
	default:
		return MessageNormal
	}
}

// SendMessage will send a normal message to id
func (u *UserDatabase) SendMessage(text, opname string, to int) error {
	if !u.IsIDExists(to) {
		return ErrInvalidUserID
	}
	m := Message{ToID: to, Date: time.Now().Unix(), Text: text, Cont: opname, Pswd: "opname"}
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.InsertUnique(UserTableMessage, &m)
}

// NotifyRegister will send register notification to all supers
func (u *UserDatabase) NotifyRegister(ip, name, cont, pswd string) error {
	if name == "" {
		return ErrEmptyName
	}
	if cont == "" {
		return ErrEmptyContact
	}
	if pswd == "" {
		return ErrEmptyPassword
	}
	for _, c := range name {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') && !(c >= 'a' && c <= 'z') {
			return ErrInvalidName
		}
	}

	if u.IsNameExists(name) {
		return ErrInvalidName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: "收到来自 " + ip + ", 用户名 " + name + " 的注册请求, 联系方式: " + cont,
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

// NotifyResetPassword will send notification to all supers
func (u *UserDatabase) NotifyResetPassword(ip, name, cont string) error {
	if name == "" {
		return ErrEmptyName
	}
	if cont == "" {
		return ErrEmptyContact
	}
	for _, c := range name {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'Z') && !(c >= 'a' && c <= 'z') {
			return ErrInvalidName
		}
	}

	user, err := u.GetUserByName(name)
	if err != nil {
		return err
	}
	if cont != user.Cont {
		return ErrInvalidContact
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	err = u.SendMessage("发送重置密码请求", user.Name, *user.ID)
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: "收到来自 " + ip + ", 用户名 " + user.Name + " 的重置密码请求, 联系方式: " + user.Cont,
		Name: user.Name,
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
func (u *UserDatabase) notifyUserAdded(opname, name string, nuid int) error {
	if opname == "" || name == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: opname + " 添加了用户 " + name,
		Cont: opname,
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		if nuid == to {
			continue
		}
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyContactChange will send notification to all supers
func (u *UserDatabase) notifyContactChange(name, cont string, id int) error {
	if name == "" {
		return ErrEmptyName
	}
	if cont == "" {
		return ErrEmptyContact
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
		if id == to {
			continue
		}
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyPasswordChange will send notification to all supers
func (u *UserDatabase) notifyPasswordChange(name, npwd, opname string, id int) error {
	if name == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: "用户 " + name + " 被 " + opname + " 更改了密码",
		Name: name,
		Pswd: npwd,
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		if id == to {
			continue
		}
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyPasswordChange will send notification to all supers
func (u *UserDatabase) notifyUpdateUserRole(name, opname string, role UserRole, id int) error {
	if name == "" || opname == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: name + " 的权限被 " + opname + " 变更为 " + role.Nick(),
		Cont: opname,
		Pswd: "opname",
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		if id == to {
			continue
		}
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// notifyPasswordChange will send notification to all supers
func (u *UserDatabase) notifyDisableUser(name, opname string, id int) error {
	if name == "" || opname == "" {
		return ErrEmptyName
	}

	tos, err := u.GetSuperIDs()
	if err != nil {
		return err
	}

	m := Message{
		Date: time.Now().Unix(),
		Text: name + " 的账户被 " + opname + " 禁用",
		Cont: opname,
		Pswd: "opname",
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, to := range tos {
		if id == to {
			continue
		}
		m.ToID = to
		err = u.db.InsertUnique(UserTableMessage, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMessagesOfUser will change non-empty Pswd field to "-"
func (u *UserDatabase) GetMessagesOfUser(to int) (ms []Message, err error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	n, err := u.db.Count(UserTableMessage)
	if err != nil {
		return
	}
	ms = make([]Message, 0, n)
	m := Message{}
	err = u.db.FindFor(UserTableMessage, &m, "WHERE ToID="+strconv.Itoa(to)+" ORDER BY Date DESC", func() error {
		if m.Pswd != "" {
			m.Pswd = "-"
		}
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

// MonthlyAPIVisit counts the api visit history
type MonthlyAPIVisit struct {
	YM    uint32 // YM is yyyymm
	Count uint32 // visit count this mounth
}

// VisitAPI increases count of this mounth by 1
func (u *UserDatabase) VisitAPI() {
	now := time.Now()
	ym := uint32(now.Year())*100 + uint32(now.Month())
	var v MonthlyAPIVisit
	u.mu.Lock()
	defer u.mu.Unlock()
	_ = u.db.Find(UserTableMonthlyAPIVisit, &v, "WHERE YM="+strconv.FormatUint(uint64(ym), 10))
	v.YM = ym
	v.Count++
	err := u.db.Insert(UserTableMonthlyAPIVisit, &v)
	if err != nil {
		logrus.Warnln("[global.user] insert visit error:", err)
	}
}

// GetAnnualAPIVisitCount get the latest 12 mounths' count
func (u *UserDatabase) GetAnnualAPIVisitCount() (cnts [12]uint32) {
	var v MonthlyAPIVisit
	var yms [12]uint32
	now := time.Now()
	y100 := uint32(now.Year()) * 100
	py100 := uint32(now.Year()-1) * 100
	nm := int(now.Month())
	for i := 0; i < nm; i++ {
		yms[i] = y100 + uint32(i+1)
	}
	for i := nm; i < 12; i++ {
		yms[i] = py100 + uint32(i+1)
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	i := 0
	for _, ym := range yms {
		_ = u.db.Find(UserTableMonthlyAPIVisit, &v, "WHERE YM="+strconv.FormatUint(uint64(ym), 10))
		cnts[i] = v.Count
		i++
	}
	return
}
