import { defHttp } from '/@/utils/http/axios'
import { getFileListModel, FilePercent, DelFile, AnalyzeFile } from './model/fileListModel'
import { DownloadFile, FileStatus } from './model/fileModel'

enum Api {
  GetFileList = '/getFileList',
  GetFilePercent = '/getFilePercent',
  DelFile = '/delFile',
  AnalyzeFile = '/analyzeFile',
  DlFile = '/dlFile',
  GetFileStatus = '/getFileStatus',
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

/**
 * @description: Download file
 */
export const downloadFile = (id: number) => {
  return defHttp.get<DownloadFile>({ url: Api.DlFile, params: { id: id } })
}

/**
 * @description: Get file status
 */
export const getFileStatus = (id: number) => {
  return defHttp.get<FileStatus>({ url: Api.GetFileStatus, params: { id: id } })
}
