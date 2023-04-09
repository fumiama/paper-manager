export interface FileListGroupItem {
  id: number
  title: string
  description: string
  size: number
  questions: number
  author: string
  datetime: string
  percent: number
}

/**
 * @description: Get filelist return value
 */
export type getFileListModel = FileListGroupItem[]

export interface DelFile {
  msg: string
}

export interface AnalyzeFile {
  code: number
  msg: string
}
