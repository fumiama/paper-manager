import { defHttp, paperHttp } from '/@/utils/http/axios'
import {
  getFileListModel,
  AnalyzeFile,
  FileListGroupItem,
  GenerateConfig,
} from './model/fileListModel'
import { DownloadFile, FileStatus } from './model/fileModel'

enum Api {
  GetFileList = '/getFileList',
  GetFileInfo = '/getFileInfo',
  GetFilePercent = '/getFilePercent',
  DelFile = '/delFile',
  AnalyzeFile = '/analyzeFile',
  DlFile = '/dlFile',
  GetFileStatus = '/getFileStatus',
  GetMajors = '/getMajors',
  GenFile = '/genFile',
  DlGen = '/dlGen',
}

/**
 * @description: Get file list
 */
export const getFileList = (permanent: boolean, count?: number) => {
  return defHttp.get<getFileListModel>({ url: Api.GetFileList, params: { count, permanent } })
}

/**
 * @description: Get file info
 */
export const getFileInfo = (id: number) => {
  return defHttp.get<FileListGroupItem>({ url: Api.GetFileInfo, params: { id } })
}

/**
 * @description: Get file percent
 */
export const getFilePercent = (id: number) => {
  return defHttp.get<number>({ url: Api.GetFilePercent, params: { id: id } })
}

/**
 * @description: Get file percent
 */
export const delFile = (id: number, permanent: boolean) => {
  return defHttp.get<string>({ url: Api.DelFile, params: { id, permanent } })
}

/**
 * @description: Analyze file
 */
export const analyzeFile = (id: number, permanent: boolean) => {
  return defHttp.get<AnalyzeFile>(
    { url: Api.AnalyzeFile, params: { id: id, permanent: permanent } },
    { errorMessageMode: 'none' },
  )
}

/**
 * @description: Download file
 */
export const downloadFile = (id: number) => {
  return defHttp.get<DownloadFile>({ url: Api.DlFile, params: { id: id } })
}

/**
 * @description: Download file to blob
 */
export const getFileBlob = (url: string) => {
  return paperHttp.get<any>({
    responseType: 'blob',
    url: url,
  })
}

/**
 * @description: Get file status
 */
export const getFileStatus = (id: number) => {
  return defHttp.get<FileStatus>({ url: Api.GetFileStatus, params: { id: id } })
}

/**
 * @description: Get majors
 */
export const getMajors = () => {
  return defHttp.get<string[]>({ url: Api.GetMajors })
}

/**
 * @description: Generate File
 */
export const generateFile = (config: GenerateConfig) => {
  return defHttp.post<string>({ url: Api.GenFile, params: config }, { errorMessageMode: 'none' })
}

/**
 * @description: Download generated file
 */
export const dlGeneratedFile = () => {
  return paperHttp.get<any>(
    { url: '/api' + Api.DlGen, responseType: 'blob' },
    { errorMessageMode: 'none' },
  )
}
