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
] as MockMethod[]
