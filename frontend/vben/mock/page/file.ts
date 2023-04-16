import { MockMethod } from 'vite-plugin-mock'
import { resultError, resultSuccess, getRequestToken, requestParams } from '../_util'

export default [
  /*{
    url: '/api/dlFile',
    timeout: 1000,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const id = Number(request.query.id)
      if (!id || id < 0) return resultError('Invalid id')
      return resultSuccess({
        url: '/file/' + id + '.docx',
      })
    },
  },*/
  {
    url: '/api/getFileStatus',
    timeout: 500,
    method: 'get',
    response: (request: requestParams) => {
      const token = getRequestToken(request)
      if (!token) return resultError('Invalid token')
      const id = Number(request.query.id)
      if (!id || id < 0) return resultError('Invalid id')
      return resultSuccess({
        name: '100.docx',
        size: 1.5,
        questions: [
          { count: 4, point: 10, name: '填空题' },
          { count: 10, point: 20, name: '不定项选择题' },
          { count: 5, point: 10, name: '判断改错题' },
          { count: 5, point: 30, name: '简述题' },
          { count: 4, point: 30, name: '综合题' },
        ],
        duplications: [
          { percent: 10, name: '不定项选择题.1' },
          { percent: 20, name: '判断改错题.2' },
          { percent: 30, name: '简述题.3' },
          { percent: 40, name: '综合题.4' },
          { percent: 50, name: '填空题.5' },
          { percent: 60, name: '不定项选择题.6' },
          { percent: 70, name: '判断改错题.7' },
          { percent: 80, name: '简述题.8' },
          { percent: 90, name: '综合题.9' },
          { percent: 100, name: '填空题.10' },
        ],
      })
    },
  },
] as MockMethod[]
