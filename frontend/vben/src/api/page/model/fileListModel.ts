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

export interface AnalyzeFile {
  code: number
  msg: string
}

export interface GenerateConfig {
  Distribution: { [x: string]: any } // Distribution is map[majorname]subcount
  RateLimit: number // RateLimit 重复率上限
  YearStart: number // YearStart 起始年份（空则直到最旧）
  YearEnd: number // YearEnd 截止年份（空则直到最新）
  TypeMask: number // TypeMask & File.Type != 0 则匹配
}
