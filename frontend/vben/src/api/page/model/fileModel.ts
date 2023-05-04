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
  rate: number
  questions: Question[]
  duplications: Duplication[]
}

export interface FileDupStatus {
  name: string
  size: number
  rate: number
  questions: Question[]
  duplications: Duplication[]
  files: Duplication[]
}
