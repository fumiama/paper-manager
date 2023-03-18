export interface DownloadFile {
  url: string
}

export interface Question {
  count: number
  point: number
  name: string
}

export interface Duplication {
  percent: number
  name: string
}

export interface FileStatus {
  name: string
  size: number
  questions: Question[]
  duplications: Duplication[]
}