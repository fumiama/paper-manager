import { defHttp } from '/@/utils/http/axios'
import { MessageItem } from './model/workbenchModel'

enum Api {
  GetMessageList = '/getMessageList',
}

export const getMessageList = () => {
  return defHttp.get<MessageItem[]>({ url: Api.GetMessageList })
}
