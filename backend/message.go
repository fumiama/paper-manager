package backend

import (
	"errors"
	"time"

	"github.com/fumiama/paper-manager/backend/global"
)

var (
	errInvalidMessageID   = errors.New("invalid message id")
	errNoAcceptPermission = errors.New("no accept permission")
	errNothingToDo        = errors.New("nothing to do")
)

type messageList struct {
	ID     int                `json:"id"`
	Avatar string             `json:"avatar"`
	Date   string             `json:"date"`
	Text   string             `json:"text"`
	Type   global.MessageType `json:"type"`
}

func getMessageList(token string) ([]messageList, error) {
	user := usertokens.Get(token)
	if user == nil {
		return nil, errInvalidToken
	}
	ms, err := global.UserDB.GetMessagesOfUser(*user.ID)
	if err != nil {
		return nil, nil
	}
	if len(ms) == 0 {
		return nil, nil
	}
	ml := make([]messageList, len(ms))
	am := make(map[string]string, 64)
	for i, m := range ms {
		avtr := ""
		n := m.Name
		if n == "" {
			n = m.Cont
		}
		if a, ok := am[n]; ok {
			avtr = a
		} else {
			u, err := global.UserDB.GetUserByName(n)
			if err == nil {
				avtr = u.Avtr
				am[n] = u.Avtr
			}
		}
		ml[i].ID = *m.ID
		ml[i].Avatar = avtr
		ml[i].Date = time.Unix(m.Date, 0).Format(chineseDateLayout)
		ml[i].Text = m.Text
		ml[i].Type = m.Type()
	}
	return ml, nil
}

func acceptMessage(token string, id int) error {
	user := usertokens.Get(token)
	if user == nil {
		return errInvalidToken
	}
	if !user.IsSuper() {
		return errNoAcceptPermission
	}
	m, err := global.UserDB.GetMessageByID(id)
	if err != nil {
		return err
	}
	if m.ToID != *user.ID {
		return errInvalidMessageID
	}
	switch m.Type() {
	case global.MessageRegister:
		return global.UserDB.AddUser(&global.User{
			Role: global.RoleUser,
			Pswd: m.Pswd,
			Name: m.Name,
			Cont: m.Cont,
		}, user.Name)
	case global.MessageResetPassword:
		u, err := global.UserDB.GetUserByName(m.Name)
		if err != nil {
			return err
		}
		return global.UserDB.UpdateUserPassword(*u.ID, user.Name, "123456")
	default:
		return errNothingToDo
	}
}

func delMessage(token string, id int) error {
	user := usertokens.Get(token)
	if user == nil {
		return errInvalidToken
	}
	m, err := global.UserDB.GetMessageByID(id)
	if err != nil {
		return err
	}
	if m.ToID != *user.ID {
		return errInvalidMessageID
	}
	return global.UserDB.DelMessageByID(id)
}
