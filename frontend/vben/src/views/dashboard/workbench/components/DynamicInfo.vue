<template>
  <Card title="我的消息" v-bind="$attrs">
    <List item-layout="horizontal" :data-source="dynamicInfoItems">
      <template #renderItem="{ item }">
        <ListItem>
          <ListItemMeta>
            <template #description>
              {{ item.date }}
            </template>
            <!-- eslint-disable-next-line -->
            <template #title> {{ item.name }} <span v-html="item.desc"> </span> </template>
            <template #avatar>
              <Icon :icon="item.avatar" :size="30" />
            </template>
          </ListItemMeta>
          <a-button ghost color="error">删除</a-button>
        </ListItem>
      </template>
    </List>
  </Card>
</template>
<script lang="ts" setup>
  import { Card, List } from 'ant-design-vue'
  import { dynamicInfoItems } from './data'
  import { Icon } from '/@/components/Icon'

  const ListItem = List.Item
  const ListItemMeta = List.Item.Meta
</script>
