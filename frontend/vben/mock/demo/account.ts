import { MockMethod } from 'vite-plugin-mock'
import { resultSuccess, resultError } from '../_util'
import { ResultEnum } from '../../src/enums/httpEnum'

const userInfo = {
  name: 'fumiama',
  userid: '00000001',
  email: 'fumiama@demo.com',
  signature: '天何所沓，十二焉分。日月安属，列星安陈。',
  introduction: '日は山の端にかかりぬ。',
  title: '课程组长',
  group: '信息与通信工程学院-网络工程系',
  tags: [
    {
      key: '0',
      label: '很有想法的',
    },
    {
      key: '1',
      label: '专注设计',
    },
    {
      key: '2',
      label: '辣~',
    },
    {
      key: '3',
      label: '大长腿',
    },
    {
      key: '4',
      label: '川妹子',
    },
    {
      key: '5',
      label: '海纳百川',
    },
  ],
  notifyCount: 12,
  unreadCount: 11,
  country: '中国',
  address: '四川成都',
  phone: '028-61830156',
}

export default [
  {
    url: '/api/account/getAccountInfo',
    timeout: 1000,
    method: 'get',
    response: () => {
      return resultSuccess(userInfo)
    },
  },
  {
    url: '/api/user/sessionTimeout',
    method: 'post',
    statusCode: 401,
    response: () => {
      return resultError()
    },
  },
  {
    url: '/api/user/tokenExpired',
    method: 'post',
    statusCode: 200,
    response: () => {
      return resultError('Token Expired!', { code: ResultEnum.TIMEOUT as number })
    },
  },
] as MockMethod[]
