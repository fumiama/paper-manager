<template>
  <Card title="项目" v-bind="$attrs">
    <template #extra>
      <a-button type="link" size="small" @click="nav2filelist">更多</a-button>
    </template>

    <CardGrid v-for="item in items" :key="item.title" class="!md:w-1/3 !w-full">
      <span class="flex">
        <Icon :icon="item.icon" :color="item.color" size="30" />
        <span class="text-lg ml-4">{{ item.title }}</span>
      </span>
      <div class="flex mt-2 h-10 text-secondary">{{ item.desc }}</div>
      <div class="flex justify-between text-secondary">
        <span>{{ item.group }}</span>
        <span>{{ item.date }}</span>
      </div>
    </CardGrid>
  </Card>
</template>
<script lang="ts">
  import { defineComponent } from 'vue'
  import { Card } from 'ant-design-vue'
  import { Icon } from '/@/components/Icon'
  import { getFileList } from '/@/api/page/page'
  import { router } from '/@/router'
  import { PageEnum } from '/@/enums/pageEnum'

  async function nav2filelist() {
    router.push(PageEnum.PAGE_FILELIST)
  }

  const fl = await getFileList(6)

  for (let i = 0; i < fl.length; i++) {
    fl[i].icon = 'ion:newspaper-outline'
  }

  export default defineComponent({
    components: { Card, CardGrid: Card.Grid, Icon },
    setup() {
      return { items: fl, nav2filelist: nav2filelist }
    },
  })
</script>
