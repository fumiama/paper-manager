import { getFileList, getFilePercent } from '/@/api/page'
import { getFileListModel } from '/@/api/page/model/fileListModel'

function refreshFilePercent(arr: any[], i: number) {
  return async () => {
    const p = await getFilePercent(arr[i].id)
    arr[i].percent = p.percent
    if (p.percent < 100) {
      setTimeout(refreshFilePercent(arr, i), 1000)
    }
  }
}

export const { cardList, totalSize, totalQuestions } = await (async () => {
  const cardList: any[] = []
  const lst = (await getFileList()) as getFileListModel
  let totalSize = 0
  let totalQuestions = 0
  for (let i = 0; i < 100; i++) {
    cardList.push({
      id: lst[i].id,
      title: lst[i].title,
      description: lst[i].description,
      size: lst[i].size,
      questions: lst[i].questions,
      datetime: lst[i].datetime,
      icon: 'bi:filetype-docx',
      color: '#1890ff',
      author: lst[i].author,
      percent: lst[i].percent,
    })
    if (cardList[i].percent < 100) {
      setTimeout(refreshFilePercent(cardList, i), 10000)
    }
    totalSize += lst[i].size
    totalQuestions += lst[i].questions
  }
  return {
    cardList,
    totalSize,
    totalQuestions,
  }
})()
