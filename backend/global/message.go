package global

import (
	"strconv"
	"time"
)

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
