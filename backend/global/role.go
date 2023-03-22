package global

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
