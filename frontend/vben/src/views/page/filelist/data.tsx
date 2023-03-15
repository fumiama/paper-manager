export const cardList = (() => {
  const result: any[] = []
  for (let i = 0; i < 100; i++) {
    result.push({
      id: i,
      title: 'Vben Admin',
      description: '基于Vue Next, TypeScript, Ant Design Vue实现的一套完整的企业级后台管理系统',
      datetime: '2020-11-26 17:39',
      icon: 'bi:filetype-docx',
      color: '#1890ff',
      author: 'Vben',
      percent: i,
    })
  }
  return result
})()
