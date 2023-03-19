<template>
  <Card title="我的消息" v-bind="$attrs">
    <List item-layout="horizontal" :data-source="dynamicInfoItemsRef">
      <template #renderItem="{ item }">
        <ListItem>
          <ListItemMeta>
            <template #description>
              {{ item.date }}
            </template>
            <!-- eslint-disable-next-line -->
            <template #title> <span v-html="item.text"> </span> </template>
            <template #avatar>
              <Avatar :src="item.avatar || headerImg" :size="36" />
            </template>
          </ListItemMeta>
          <a-button
            ghost
            color="success"
            v-if="
              [MessageTypeEnum.MessageRegister, MessageTypeEnum.MessageResetPassword].includes(
                item.type,
              )
            "
            >接受</a-button
          >
          &nbsp;&nbsp;
          <a-button ghost color="error">删除</a-button>
        </ListItem>
      </template>
    </List>
  </Card>
</template>
<script lang="ts" setup>
  import { ref } from 'vue'
  import { Card, List } from 'ant-design-vue'
  import { getMessageList } from '/@/api/dashboard/index'
  import { MessageTypeEnum, MessageItem } from '/@/api/dashboard/model/workbenchModel'
  import { Avatar } from 'ant-design-vue'
  import headerImg from '/@/assets/images/header.jpg'

  const ListItem = List.Item
  const ListItemMeta = List.Item.Meta
  const dynamicInfoItemsRef = ref([] as MessageItem[])
  getMessageList().then((value) => {
    dynamicInfoItemsRef.value = value
  })
</script>
