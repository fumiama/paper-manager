import { reactive } from 'vue'
import { getFileList, getFilePercent } from '/@/api/page'
import { getFileListModel } from '/@/api/page/model/fileListModel'

export const random = (min: number, max: number) =>
  Math.floor(Math.random() * (max - min + 1) + min)

export function refreshFilePercent(item: any) {
  return async () => {
    const p = await getFilePercent(item.id)
    if (p.percent) {
      item.percent = p.percent
      if (p.percent < 100) {
        setTimeout(refreshFilePercent(item), 1000)
      }
    } else item.hassettimeout = false
  }
}

export function getListOfPage(pageSize: number, page: number): any[] {
  const i = page - 1
  let lst: any[] = []
  if (i < cardList._cardList.length / pageSize)
    lst = reactive(cardList._cardList.slice(i * pageSize, page * pageSize))
  else lst = reactive(cardList._cardList.slice((cardList._cardList.length / pageSize) * pageSize))
  for (let i = 0; i < lst.length; i++) {
    if (!lst[i].hassettimeout && lst[i].percent > 0 && lst[i].percent < 100) {
      setTimeout(refreshFilePercent(lst[i]), 1000 + random(0, 1000))
      lst[i].hassettimeout = true
    }
  }
  return lst
}

async function refreshFileList() {
  const __cardList: any[] = []
  const lst = (await getFileList()) as getFileListModel
  let __totalSize = 0
  let __totalQuestions = 0
  for (let i = 0; i < lst.length; i++) {
    __cardList.push({
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
      hassettimeout: false,
      delloading: false,
    })
    __totalSize += lst[i].size
    __totalQuestions += lst[i].questions
  }
  return {
    _cardList: __cardList,
    _totalSize: __totalSize,
    _totalQuestions: __totalQuestions,
  }
}

export let cardList = reactive(await refreshFileList())

export let pagination = reactive({
  current: 1,
  total: cardList._cardList.length,
  show: true,
  pageSize: 10,
  onChange: function (page: number, pageSize: number) {
    this.current = page
    this.pageSize = pageSize
  },
})

export function refreshCardList() {
  refreshFileList().then((value) => {
    cardList = reactive(value)
    pagination = reactive({
      current: 1,
      total: cardList._cardList.length,
      show: true,
      pageSize: 10,
      onChange: function (page: number, pageSize: number) {
        this.current = page
        this.pageSize = pageSize
      },
    })
  })
}

export function deleteFileByID(id: number) {
  cardList._cardList.map((value: any, index: number) => {
    if (value.id == id) {
      cardList._cardList.splice(index, 1)
      cardList._totalSize -= value.size
      cardList._totalQuestions -= value.questions
      pagination.total = cardList._cardList.length
    }
  })
}
