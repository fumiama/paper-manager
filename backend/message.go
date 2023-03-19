package backend

import (
	"time"

	"github.com/fumiama/paper-manager/backend/global"
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
	for i, m := range ms {
		avtr := ""
		u, err := global.UserDB.GetUserByName(m.Name)
		if err == nil {
			avtr = u.Avtr
		}
		ml[i].ID = *m.ID
		ml[i].Avatar = avtr
		ml[i].Date = time.Unix(m.Date, 0).Format(chineseDateLayout)
		ml[i].Text = m.Text
		ml[i].Type = m.Type()
	}
	return ml, nil
}