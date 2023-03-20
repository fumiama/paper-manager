import { defHttp } from '/@/utils/http/axios'
import {
  LoginParams,
  LoginResultModel,
  GetUserInfoModel,
  ResetPasswordParams,
  SetPasswordParams,
  SetContactParams,
  SetUserInfoParams,
  RegisterParams,
  ResetPasswordResultModel,
  RegisterResultModel,
  GetLoginSaltModel,
  GetUsersListModel,
} from './model/userModel'

import { ErrorMessageMode } from '/#/axios'

enum Api {
  GetLoginSalt = '/getLoginSalt',
  Login = '/login',
  Logout = '/logout',
  ResetPassword = '/resetPassword',
  SetPassword = '/setPassword',
  SetContact = '/setContact',
  SetUserInfo = '/setUserInfo',
  Register = '/register',
  GetUserInfo = '/getUserInfo',
  GetUsersCount = '/getUsersCount',
  GetUsersList = '/getUsersList',
  IsNameExist = '/isNameExist',
}

/**
 * @description: user login api
 */
export function loginApi(params: LoginParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<LoginResultModel>(
    {
      url: Api.Login,
      params,
    },
    {
      errorMessageMode: mode,
    },
  )
}

/**
 * @description: reset password api
 */
export function resetPasswordApi(params: ResetPasswordParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<ResetPasswordResultModel>(
    {
      url: Api.ResetPassword,
      params,
    },
    {
      errorMessageMode: mode,
    },
  )
}

/**
 * @description: set password api, borrowing the ResetPasswordResultModel as they're the same
 */
export function setPasswordApi(params: SetPasswordParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<ResetPasswordResultModel>(
    {
      url: Api.SetPassword,
      params,
    },
    {
      errorMessageMode: mode,
    },
  )
}

/**
 * @description: set contact api, borrowing the ResetPasswordResultModel as they're the same
 */
export function setContactApi(params: SetContactParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<ResetPasswordResultModel>(
    {
      url: Api.SetContact,
      params,
    },
    {
      errorMessageMode: mode,
    },
  )
}

/**
 * @description: set userinfo api, borrowing the ResetPasswordResultModel as they're the same
 */
export function setUserInfoApi(params: SetUserInfoParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<ResetPasswordResultModel>(
    {
      url: Api.SetUserInfo,
      params,
    },
    {
      errorMessageMode: mode,
    },
  )
}

/**
 * @description: register api
 */
export function registerApi(params: RegisterParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<RegisterResultModel>(
    {
      url: Api.Register,
      params,
    },
    {
      errorMessageMode: mode,
    },
  )
}

/**
 * @description: getLoginSalt
 */
export function getLoginSalt(username: string) {
  return defHttp.get<GetLoginSaltModel>({ url: Api.GetLoginSalt, params: { username: username } })
}

/**
 * @description: getUserInfo
 */
export function getUserInfo() {
  return defHttp.get<GetUserInfoModel>({ url: Api.GetUserInfo }, { errorMessageMode: 'none' })
}

/**
 * @description: getUsersCount
 */
export function getUsersCount() {
  return defHttp.get<number>({ url: Api.GetUsersCount }, { errorMessageMode: 'none' })
}

/**
 * @description: getUsersList
 */
export function getUsersList() {
  return defHttp.get<GetUsersListModel[]>({ url: Api.GetUsersList }, { errorMessageMode: 'none' })
}

/**
 * @description: isNameExist
 */
export function isNameExist(username: string) {
  return defHttp.get<boolean>(
    { url: Api.IsNameExist, params: { username } },
    { errorMessageMode: 'none' },
  )
}

export function doLogout() {
  return defHttp.get({ url: Api.Logout }, { errorMessageMode: 'none' })
}
