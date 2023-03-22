package global

import (
	"reflect"
	"regexp"
	"strconv"
)

// Regex stores user's config of splitting docx file
type Regex struct {
	ID     int    // ID is User(ID)
	Title  string // Title default `.*(\d{4})\s*-.*学年.*(\d?).*([中末]?).*([AB]?)\s*卷`
	Class  string // Class default `考试科目：\s*(\S+)\s*`
	OpenCl string // OpenCl default `考试形式：\s*(\S+)\s*`
	Date   string // Date default `考试日期：\s*(\d+)\s*年\s*(\d+)\s*月\s*(\d+)\s*日`
	Time   string // Time default `考试时长：\s*(\d+)\s*分钟`
	Rate   string // Rate default `成绩构成比例：\s*(.*%)\s*`
	Major  string // Major default `([一二三四五六七八九十]+)、\s*(.*)\s*（.*([空题]?)\s*(\d*).*共\s*(\d+)\s*分.*）`
	Sub    string // Sub default `(\d+)、`
}

func newRegex() (reg Regex) {
	reg.Title = `.*(\d{4})\s*-.*学年.*(\d?).*([中末]?).*([AB]?)\s*卷`
	reg.Class = `考试科目：\s*(\S+)\s*`
	reg.OpenCl = `考试形式：\s*(\S+)\s*`
	reg.Date = `考试日期：\s*(\d+)\s*年\s*(\d+)\s*月\s*(\d+)\s*日`
	reg.Time = `考试时长：\s*(\d+)\s*分钟`
	reg.Rate = `成绩构成比例：\s*(.*%)\s*`
	reg.Major = `([一二三四五六七八九十]+)、\s*(.*)\s*（.*([空题]?)\s*(\d*).*共\s*(\d+)\s*分.*）`
	reg.Sub = `(\d+)、`
	return
}

// SetUserRegex set Regex.name = re
func (u *UserDatabase) SetUserRegex(id int, name, re string) error {
	if name == "" || name == "ID" {
		return ErrInvalidFieldName
	}
	if re == "" {
		return ErrEmptyRegex
	}
	user, err := UserDB.GetUserByID(id)
	if err != nil {
		return err
	}
	if !user.IsFileManager() {
		return ErrInvalidRole
	}
	_, err = regexp.Compile(re)
	if err != nil {
		return err
	}
	reg := newRegex()
	rreg := reflect.ValueOf(&reg).Elem()
	f := rreg.FieldByName(name)
	if !f.IsValid() {
		return ErrNoSuchFieldName
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	_ = u.db.Find(UserTableRegex, &reg, "WHERE ID="+strconv.Itoa(id))
	reg.ID = id
	f.SetString(re)
	return u.db.Insert(UserTableRegex, &reg)
}

// GetUserRegex default newRegex()
func (u *UserDatabase) GetUserRegex(id int) (*Regex, error) {
	user, err := UserDB.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if !user.IsFileManager() {
		return nil, ErrInvalidRole
	}
	reg := newRegex()
	u.mu.RLock()
	_ = u.db.Find(UserTableRegex, &reg, "WHERE ID="+strconv.Itoa(id))
	u.mu.RUnlock()
	return &reg, nil
}
