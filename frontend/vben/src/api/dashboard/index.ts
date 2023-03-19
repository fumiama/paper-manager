import { defHttp } from '/@/utils/http/axios'
import { MessageItem } from './model/workbenchModel'

enum Api {
  GetMessageList = '/getMessageList',
  AcceptMessage = '/acceptMessage',
}

export const getMessageList = () => {
  return defHttp.get<MessageItem[]>({ url: Api.GetMessageList })
}

export const acceptMessage = (id: number) => {
  return defHttp.get<string>({ url: Api.AcceptMessage, params: { id } })
}
