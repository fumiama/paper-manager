/**
 * @description: Login interface parameters
 */
export interface LoginParams {
  username: string
  password: string
}

/**
 * @description: Reset password interface parameters
 */
export interface ResetPasswordParams {
  username: string
  mobile: string
}

/**
 * @description: Set password interface parameters
 */
export interface SetPasswordParams {
  token: string
  password: string
}

/**
 * @description: Set Contact interface parameters
 */
export interface SetContactParams {
  token: string
  contact: string
}

/**
 * @description: Set UserInfo interface parameters
 */
export interface SetUserInfoParams {
  nick: string
  desc: string
  avtr: string
}

/**
 * @description: Set Others' Info interface parameters
 */
export interface SetOthersInfoParams {
  id: number
  nick: string
  desc: string
}

/**
 * @description: Register interface parameters
 */
export interface RegisterParams {
  username: string
  mobile: string
  password: string
}

export interface RoleInfo {
  roleName: string
  value: string
}

/**
 * @description: Login interface return value
 */
export interface LoginResultModel {
  userId: string | number
  token: string
  role: RoleInfo
}

/**
 * @description: Reset password interface return value
 */
export interface ResetPasswordResultModel {
  msg: string
}

/**
 * @description: Register interface return value
 */
export interface RegisterResultModel {
  msg: string
}

/**
 * @description: Get user information return value
 */
export interface GetUserInfoModel {
  roles: RoleInfo[]
  // 用户id
  userId: string | number
  // 用户名
  username: string
  // 真实名字
  realName: string
  // 头像
  avatar: string
  // 介绍
  desc?: string
  // 创建日期
  date: string
  // 上次修改密码日期
  last: string
  // 电话
  contact: string
}

export interface GetLoginSaltModel {
  salt: string
}

export interface GetUsersListModel {
  id: number
  name: string
  nick: string
  role: number
  date: string
  desc: string
}
