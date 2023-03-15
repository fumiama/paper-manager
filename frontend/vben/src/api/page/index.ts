import { defHttp } from '/@/utils/http/axios'
import { getFileListModel, FilePercent } from './model/fileListModel'

enum Api {
  GetFileList = '/getFileList',
  GetFilePercent = '/getFilePercent',
}

/**
 * @description: Get file list
 */
export const getFileList = (count?: number) => {
  return defHttp.get<getFileListModel>({ url: Api.GetFileList, params: { count: count } })
}

/**
 * @description: Get file percant
 */
export const getFilePercent = (id: number) => {
  return defHttp.get<FilePercent>({ url: Api.GetFilePercent, params: { id: id } })
}
