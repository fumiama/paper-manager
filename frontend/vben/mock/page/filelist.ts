import { MockMethod } from 'vite-plugin-mock'
import { resultError, resultSuccess, getRequestToken, requestParams } from '../_util'

const deletedIDs: number[] = []

function createFileList() {
  const lst: any[] = []
  for (let i = 100; i > 0; i--) {
    if (deletedIDs.includes(i)) continue
    lst.push({
      id: i,
      title: '接入网2020B卷',
      description: '电子科技大学接入网2020B卷',
      size: i,
      questions: 10,
      author: '课程组长',
      datetime: '2020-11-26 17:39',
      percent: 0,
    })
  }
  return lst
}

export default [
  // mock get filelist
  {
    url: '/basic-api/getFileList',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const count = request.query.count
      if (!count || count <= 0) return resultSuccess(createFileList())
      let fl = createFileList()
      if (fl.length > count) {
        fl = fl.slice(0, count)
      }
      return resultSuccess(fl)
    },
  },
  {
    url: '/basic-api/getFilePercent',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const id = request.query.id
      if (!id || id < 0) return resultError('Invalid id')
      return resultSuccess({
        percent: 100,
      })
    },
  },
  {
    url: '/basic-api/delFile',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const id = Number(request.query.id)
      if (!id || id < 0) return resultError('Invalid id')
      deletedIDs.push(id)
      return resultSuccess({
        msg: '已成功删除文件' + id + '.',
      })
    },
  },
] as MockMethod[]
