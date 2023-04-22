package global

import (
	"reflect"
	"regexp"
	"strconv"

	sql "github.com/FloatTech/sqlite"
)

// Regex stores user's config of splitting docx file
type Regex struct {
	ID     int    // ID is User(ID)
	Title  string // Title default `.*(\d{4})\s*-.*学年.*(\d).*([中末]).*([AB])\s*卷`
	Class  string // Class default `(考试科目|课程名称)：\s*(\S+)\s*`
	OpenCl string // OpenCl default `考试形式：\s*(\S+)\s*`
	Date   string // Date default `考试日期：\s*(\d+)\s*年\s*(\d+)\s*月\s*(\d+)\s*日`
	Time   string // Time default `考试时长：\s*(\d+)\s*分钟`
	Rate   string // Rate default `(成绩构成比例|课程成绩构成)：\s*(.*%)\s*`
	Major  string // Major default `([一二三四五六七八九十]+)、\s*(.*)\s*（.*([空题]?)\s*(\d*).*共\s*(\d+)\s*分.*）`
	Sub    string // Sub default `(\d+)、`
}

func GetDefaultRegex() (reg Regex) {
	reg.Title = `.*(\d{4})\s*-.*学年.*(\d).*([中末]).*([AB])\s*卷`
	reg.Class = `(考试科目|课程名称)：\s*(\S+)\s*`
	reg.OpenCl = `考试形式：\s*(\S+)\s*`
	reg.Date = `考试日期：\s*(\d+)\s*年\s*(\d+)\s*月\s*(\d*)\s*日`
	reg.Time = `考试时长：\s*(\d+)\s*分钟`
	reg.Rate = `(成绩构成比例|课程成绩构成)：\s*(.*%)\s*`
	reg.Major = `([一二三四五六七八九十]+)、\s*(.*)\s*（.*([空题]?)\s*(\d*).*共\s*(\d+)\s*分.*）`
	reg.Sub = `(\d+)、`
	return
}

// SetUserRegex set Regex.name = re
func (u *UserDatabase) SetUserRegex(id int, reg *Regex) error {
	if reg == nil {
		return ErrEmptyRegex
	}
	user, err := UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	if !user.IsSuper() && id != *user.ID {
		return ErrInvalidRole
	}
	defaultrf := reflect.ValueOf(GetDefaultRegex())
	rreg := reflect.ValueOf(reg).Elem()
	for i := 1; i < rreg.NumField(); i++ {
		if rreg.Field(i).Equal(defaultrf.Field(i)) {
			rreg.Field(i).SetString("")
		} else {
			_, err = regexp.Compile(rreg.Field(i).String())
			if err != nil {
				return err
			}
		}
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.db.Insert(UserTableRegex, reg)
}

// GetUserRegex default newRegex()
func (u *UserDatabase) GetUserRegex(id int) (*Regex, error) {
	user, err := UserDB.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if !user.IsSuper() || id != *user.ID {
		return nil, ErrInvalidRole
	}
	u.mu.RLock()
	reg, _ := sql.Find[Regex](&u.db, UserTableRegex, "WHERE ID="+strconv.Itoa(id))
	u.mu.RUnlock()
	reg.ID = *user.ID
	rf := reflect.ValueOf(&reg).Elem()
	defaultrf := reflect.ValueOf(GetDefaultRegex())
	for i := 1; i < rf.NumField(); i++ {
		if rf.Field(i).IsZero() {
			rf.Field(i).Set(defaultrf.Field(i))
		}
	}
	return &reg, nil
}
