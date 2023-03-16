import { MockMethod } from 'vite-plugin-mock'
import { resultError, resultSuccess, getRequestToken, requestParams } from '../_util'

export default [
  {
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
  },
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
          { count: 4, point: 10, name: '一、填空题' },
          { count: 10, point: 20, name: '二、不定项选择题' },
          { count: 5, point: 10, name: '三、判断改错题' },
          { count: 5, point: 30, name: '四、简述题' },
          { count: 4, point: 30, name: '五、综合题' },
        ],
        duplications: [
          { percent: 10, name: '二.1' },
          { percent: 20, name: '二.2' },
          { percent: 30, name: '二.3' },
          { percent: 40, name: '二.4' },
          { percent: 50, name: '二.5' },
          { percent: 60, name: '二.6' },
          { percent: 70, name: '二.7' },
          { percent: 80, name: '二.8' },
          { percent: 90, name: '二.9' },
          { percent: 100, name: '二.10' },
        ],
      })
    },
  },
] as MockMethod[]
