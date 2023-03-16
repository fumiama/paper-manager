import { MockMethod } from 'vite-plugin-mock'
import md5 from 'md5'
import { resultError, resultSuccess, getRequestToken, requestParams } from '../_util'

export function createFakeUserList() {
  return [
    {
      userId: '1',
      username: 'fumiama',
      realName: '源文雨',
      avatar: 'https://q1.qlogo.cn/g?b=qq&nk=1332524221&s=640',
      desc: 'manager',
      password: '123456',
      token: 'fakeToken1',
      homePath: '/dashboard/analysis',
      roles: [
        {
          roleName: '课程组长',
          value: 'super',
        },
      ],
    },
    {
      userId: '2',
      username: 'filemgr',
      password: '123456',
      realName: '归档代理',
      avatar: 'https://q1.qlogo.cn/g?b=qq&nk=468131917&s=640',
      desc: 'file manager',
      token: 'fakeToken2',
      homePath: '/dashboard/workbench',
      roles: [
        {
          roleName: '归档代理',
          value: 'filemgr',
        },
      ],
    },
    {
      userId: '3',
      username: 'user',
      password: '123456',
      realName: '课程组员',
      avatar: 'https://q1.qlogo.cn/g?b=qq&nk=468131931&s=640',
      desc: 'normal user',
      token: 'fakeToken3',
      homePath: '/dashboard/workbench',
      roles: [
        {
          roleName: '课程组员',
          value: 'user',
        },
      ],
    },
  ]
}

const fakeCodeList: any = {
  '1': ['1000', '3000', '5000'],

  '2': ['2000', '4000', '6000'],
}
export default [
  {
    url: '/api/getLoginSalt',
    timeout: 200,
    method: 'get',
    response: () => {
      return resultSuccess({ salt: 'nc8w9f82hfioq2ci9hcwehcq' })
    },
  },
  // mock user login
  {
    url: '/api/login',
    timeout: 200,
    method: 'post',
    response: ({ body }) => {
      const { username, password } = body
      const checkUser = createFakeUserList().find(
        (item) =>
          item.username === username &&
          password === md5(item.password + 'nc8w9f82hfioq2ci9hcwehcq'),
      )
      if (!checkUser) {
        return resultError('Incorrect account or password!')
      }
      const { userId, username: _username, token, realName, desc, roles } = checkUser
      return resultSuccess({
        roles,
        userId,
        username: _username,
        token,
        realName,
        desc,
      })
    },
  },
  // mock reset password
  {
    url: '/api/resetPassword',
    timeout: 200,
    method: 'post',
    response: ({ body }) => {
      const { username, mobile } = body
      return resultSuccess({
        msg: '已将用户' + username + '电话' + mobile + '的重置请求上报, 请耐心等待!',
      })
    },
  },
  // mock register
  {
    url: '/api/register',
    timeout: 200,
    method: 'post',
    response: ({ body }) => {
      const { username, mobile } = body
      return resultSuccess({
        msg: '已将用户' + username + '电话' + mobile + '的注册请求上报, 请耐心等待!',
      })
    },
  },
  {
    url: '/api/getUserInfo',
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const checkUser = createFakeUserList().find((item) => item.token === token)
      if (!checkUser) {
        return resultError('The corresponding user information was not obtained!')
      }
      return resultSuccess(checkUser)
    },
  },
  {
    url: '/api/getPermCode',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const checkUser = createFakeUserList().find((item) => item.token === token)
      if (!checkUser) {
        return resultError('Invalid token!')
      }
      const codeList = fakeCodeList[checkUser.userId]

      return resultSuccess(codeList)
    },
  },
  {
    url: '/api/logout',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const checkUser = createFakeUserList().find((item) => item.token === token)
      if (!checkUser) {
        return resultError('Invalid token!')
      }
      return resultSuccess(undefined, { message: 'Token has been destroyed' })
    },
  },
  {
    url: '/api/testRetry',
    statusCode: 405,
    method: 'get',
    response: () => {
      return resultError('Error!')
    },
  },
] as MockMethod[]
