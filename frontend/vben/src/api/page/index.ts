import { defHttp } from '/@/utils/http/axios'
import { getFileListModel, FilePercent, DelFile, AnalyzeFile } from './model/fileListModel'

enum Api {
  GetFileList = '/getFileList',
  GetFilePercent = '/getFilePercent',
  DelFile = '/delFile',
  AnalyzeFile = '/analyzeFile',
}

/**
 * @description: Get file list
 */
export const getFileList = (count?: number) => {
  return defHttp.get<getFileListModel>({ url: Api.GetFileList, params: { count: count } })
}

/**
 * @description: Get file percent
 */
export const getFilePercent = (id: number) => {
  return defHttp.get<FilePercent>({ url: Api.GetFilePercent, params: { id: id } })
}

/**
 * @description: Get file percent
 */
export const delFile = (id: number) => {
  return defHttp.get<DelFile>({ url: Api.DelFile, params: { id: id } })
}

/**
 * @description: Analyze file
 */
export const analyzeFile = (id: number) => {
  return defHttp.get<AnalyzeFile>({ url: Api.AnalyzeFile, params: { id: id } })
}
