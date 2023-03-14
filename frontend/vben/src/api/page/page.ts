import { defHttp } from '/@/utils/http/axios'
import { getFileListModel } from './model/fileListModel'

enum Api {
  GetFileList = '/getFileList',
}

/**
 * @description: Get file list
 */

export const getFileList = (count: number) => {
  return defHttp.get<getFileListModel>({ url: Api.GetFileList, params: { count: count } })
}
