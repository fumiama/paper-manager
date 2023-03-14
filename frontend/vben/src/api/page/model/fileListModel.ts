export interface FileListGroupItem {
  title: string
  icon: string
  color: string
  desc: string
  date: string
  group: string
}

/**
 * @description: Get filelist return value
 */
export type getFileListModel = FileListGroupItem[]
