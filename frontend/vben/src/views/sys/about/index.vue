<template>
  <PageWrapper title="关于">
    <template #headerContent>
      <div class="flex justify-between items-center">
        <span class="flex-1">
          <a :href="GITHUB_URL" target="_blank">{{ name }}</a>
          是
          <a href="https://github.com/fumiama" target="_blank">源文雨</a>
          的大学本科毕业设计项目。感谢
          <a href="https://www.sice.uestc.edu.cn/info/1302/5185.htm" target="_blank">马立香</a>
          老师在我毕业设计期间对本项目的悉心指导。
          <br />
          本项目前端使用
          <a href="https://github.com/vbenjs/vue-vben-admin" target="_blank">vue-vben-admin</a>
          ，后端使用 Golang 与 SQLite 数据库，最后统一编译为一个开箱即用的可执行文件。
        </span>
      </div>
    </template>
    <Description @register="infoRegister" class="enter-y" />
    <Description @register="register" class="my-4 enter-y" />
    <Description @register="registerDev" class="enter-y" />
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { h } from 'vue'
  import { Tag } from 'ant-design-vue'
  import { PageWrapper } from '/@/components/Page'
  import { Description, DescItem, useDescription } from '/@/components/Description/index'
  import { GITHUB_URL, SITE_URL, DOC_URL } from '/@/settings/siteSetting'

  const { pkg, lastBuildTime } = __APP_INFO__

  const { dependencies, devDependencies, name, version } = pkg

  const schema: DescItem[] = []
  const devSchema: DescItem[] = []

  const commonTagRender = (color: string) => (curVal) => h(Tag, { color }, () => curVal)
  // const commonLinkRender = (text: string) => (href) => h('a', { href, target: '_blank' }, text)

  const infoSchema: DescItem[] = [
    {
      label: '版本',
      field: 'version',
      render: commonTagRender('blue'),
    },
    {
      label: '最后编译时间',
      field: 'lastBuildTime',
      render: commonTagRender('blue'),
    },
  ]

  const infoData = {
    version,
    lastBuildTime,
    doc: DOC_URL,
    preview: SITE_URL,
    github: GITHUB_URL,
  }

  Object.keys(dependencies).forEach((key) => {
    schema.push({ field: key, label: key })
  })

  Object.keys(devDependencies).forEach((key) => {
    devSchema.push({ field: key, label: key })
  })

  const [register] = useDescription({
    title: '生产环境依赖',
    data: dependencies,
    schema: schema,
    column: 3,
  })

  const [registerDev] = useDescription({
    title: '开发环境依赖',
    data: devDependencies,
    schema: devSchema,
    column: 3,
  })

  const [infoRegister] = useDescription({
    title: '项目信息',
    data: infoData,
    schema: infoSchema,
    column: 2,
  })
</script>
