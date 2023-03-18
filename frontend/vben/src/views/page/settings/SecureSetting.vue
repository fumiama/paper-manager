<template>
  <CollapseContainer title="安全设置" :canExpan="false">
    <List>
      <template v-for="item in list" :key="item.key">
        <ListItem>
          <ListItemMeta>
            <template #title>
              {{ item.title }}
              <div class="extra" v-if="item.extra" @click="item.click">
                {{ item.extra }}
              </div>
            </template>
            <template #description>
              <div>{{ item.description }}</div>
            </template>
          </ListItemMeta>
        </ListItem>
      </template>
    </List>
  </CollapseContainer>
</template>
<script lang="ts">
  import { List } from 'ant-design-vue'
  import { defineComponent } from 'vue'
  import { CollapseContainer } from '/@/components/Container/index'
  import { useUserStore } from '/@/store/modules/user'
  import { PageEnum } from '/@/enums/pageEnum'
  import { useGo } from '/@/hooks/web/usePage'

  export default defineComponent({
    components: { CollapseContainer, List, ListItem: List.Item, ListItemMeta: List.Item.Meta },
    setup() {
      const userStore = useUserStore()
      const { last, contact } = userStore.getUserInfo
      const go = useGo()

      return {
        list: [
          {
            key: '1',
            title: '账户密码',
            description: '上次修改密码: ' + last,
            extra: '修改',
            click: () => {
              go(PageEnum.PAGE_PASSWORD_SETTINGS)
            },
          },
          {
            key: '2',
            title: '联系方式',
            description: '已记录的联系方式: ' + contact,
            extra: '修改',
            click: () => {
              go(PageEnum.PAGE_CONTACT_SETTINGS)
            },
          },
        ],
      }
    },
  })
</script>
<style lang="less" scoped>
  .extra {
    float: right;
    margin-top: 10px;
    margin-right: 30px;
    font-weight: normal;
    color: #1890ff;
    cursor: pointer;
  }
</style>
