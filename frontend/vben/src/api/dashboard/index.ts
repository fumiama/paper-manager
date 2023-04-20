import { defHttp } from '/@/utils/http/axios'
import { MessageItem, UserRegex } from './model/model'

enum Api {
  GetMessageList = '/getMessageList',
  AcceptMessage = '/acceptMessage',
  DeleteMessage = '/delMessage',
  GetAnnualVisits = '/getAnnualVisits',
  GetUserRegex = '/getUserRegex',
}

export const getAnnualVisits = () => {
  return defHttp.get<number[]>({ url: Api.GetAnnualVisits })
}

export const getMessageList = () => {
  return defHttp.get<MessageItem[]>({ url: Api.GetMessageList })
}

export const acceptMessage = (id: number) => {
  return defHttp.get<string>({ url: Api.AcceptMessage, params: { id } })
}

export const deleteMessage = (id: number) => {
  return defHttp.get<string>({ url: Api.DeleteMessage, params: { id } })
}

export const getUserRegex = () => {
  return defHttp.get<UserRegex>({ url: Api.GetUserRegex })
}
