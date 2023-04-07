import { MockMethod } from 'vite-plugin-mock'
import { resultError, resultSuccess, getRequestToken, requestParams } from '../_util'

const deletedIDs: number[] = []

// const analyzingIDs: { id: number; per: number }[] = []

/*function createFileList() {
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
}*/

export default [
  // mock get filelist
  /*{
    url: '/api/getFileList',
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
  },*/
  /*{
    url: '/api/getFilePercent',
    timeout: 200,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const id = request.query.id
      if (!id || id < 0) return resultError('Invalid id')
      let p = 0
      analyzingIDs.map((value: { id: number; per: number }, index: number) => {
        if (!p && value.id == id) {
          value.per += 10
          if (value.per >= 100) {
            analyzingIDs.splice(index, 1)
            p = 100
          }
          p = value.per
        }
      })
      if (p > 0)
        return resultSuccess({
          percent: p,
        })
      return resultSuccess({
        percent: 100,
      })
    },
  },*/
  {
    url: '/api/delFile',
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
  /*{
    url: '/api/analyzeFile',
    timeout: 1000,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const id = Number(request.query.id)
      if (!id || id < 0) return resultError('Invalid id')
      analyzingIDs.push({ id: id, per: 1 })
      return resultSuccess({
        msg: '正在分析' + id + ', 请耐心等待...',
      })
    },
  },*/
] as MockMethod[]
