import { MockMethod } from 'vite-plugin-mock'
import { resultError, resultSuccess, getRequestToken, requestParams } from '../_util'

export function createFileList() {
  return [
    {
      title: '数据库2020B卷',
      color: '',
      desc: '不要等待机会，而要创造机会。',
      group: '数据库与软件工程',
      date: '2020-04-01',
    },
    {
      title: '数据库2020A卷',
      color: '#3fb27f',
      desc: '现在的你决定将来的你。',
      group: '数据库与软件工程',
      date: '2020-04-01',
    },
    {
      title: 'TCP/IP2018B卷',
      color: '#e18525',
      desc: '没有什么才能比努力更重要。',
      group: 'TCP/IP',
      date: '2021-04-01',
    },
    {
      title: 'TCP/IP2018A卷',
      color: '#bf0c2c',
      desc: '热情和欲望可以突破一切难关。',
      group: 'TCP/IP',
      date: '2018-01-01',
    },
    {
      title: '接入网2016B卷',
      color: '#00d8ff',
      desc: '健康的身体是实目标的基石。',
      group: '接入网',
      date: '2016-01-01',
    },
    {
      title: '接入网2016A卷',
      color: '#4daf1bc9',
      desc: '路是走出来的，而不是空想出来的。',
      group: '接入网',
      date: '2016-01-01',
    },
  ]
}

export default [
  // mock user login
  {
    url: '/basic-api/getFileList',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const count = request.query.count
      if (!count || count <= 0) return resultError('Invalid count')
      let fl = createFileList()
      if (fl.length > count) {
        fl = fl.slice(0, count)
      }
      return resultSuccess(fl)
    },
  },
] as MockMethod[]
