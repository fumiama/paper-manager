import { defHttp } from '/@/utils/http/axios'
import {
  LoginParams,
  LoginResultModel,
  GetUserInfoModel,
  ResetPasswordParams,
  RegisterParams,
  ResetPasswordResultModel,
  RegisterResultModel,
} from './model/userModel'

import { ErrorMessageMode } from '/#/axios'

enum Api {
  Login = '/login',
  Logout = '/logout',
  ResetPassword = '/resetPassword',
  Register = '/register',
  GetUserInfo = '/getUserInfo',
  GetPermCode = '/getPermCode',
  TestRetry = '/testRetry',
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
 * @description: getUserInfo
 */
export function getUserInfo() {
  return defHttp.get<GetUserInfoModel>({ url: Api.GetUserInfo }, { errorMessageMode: 'none' })
}

export function getPermCode() {
  return defHttp.get<string[]>({ url: Api.GetPermCode })
}

export function doLogout() {
  return defHttp.get({ url: Api.Logout })
}

export function testRetry() {
  return defHttp.get(
    { url: Api.TestRetry },
    {
      retryRequest: {
        isOpenRetry: true,
        count: 5,
        waitTime: 1000,
      },
    },
  )
}
