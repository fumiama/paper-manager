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
}

export interface GetLoginSaltModel {
  salt: string
}
