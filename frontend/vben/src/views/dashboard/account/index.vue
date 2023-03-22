<template>
  <PageWrapper dense contentFullHeight fixedHeight contentClass="flex">
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑用户资料',
                onClick: handleEdit.bind(null, record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <AccountModal @register="registerModal" @success="handleSuccess" />
  </PageWrapper>
</template>
<script lang="ts">
  import { defineComponent } from 'vue'

  import { BasicTable, useTable, TableAction } from '/@/components/Table'
  import { getUsersList } from '/@/api/sys/user'
  import { PageWrapper } from '/@/components/Page'

  import { useModal } from '/@/components/Modal'
  import AccountModal from './AccountModal.vue'

  import { columns } from './account.data'

  export default defineComponent({
    name: 'AccountManagement',
    components: { BasicTable, PageWrapper, AccountModal, TableAction },
    setup() {
      const [registerModal, { openModal }] = useModal()
      const [registerTable, { updateTableDataRecord }] = useTable({
        title: '账号列表',
        api: getUsersList,
        rowKey: 'id',
        columns,
        useSearchForm: false,
        showTableSetting: true,
        bordered: true,
        actionColumn: {
          width: 80,
          title: '操作',
          dataIndex: 'action',
        },
      })

      function handleEdit(record: Recordable) {
        // console.log(record)
        openModal(true, {
          record,
        })
      }

      function handleSuccess({ values }) {
        // 演示不刷新表格直接更新内部数据。
        // 注意：updateTableDataRecord要求表格的rowKey属性为string并且存在于每一行的record的keys中
        updateTableDataRecord(values.id, values)
      }

      return {
        registerTable,
        registerModal,
        handleEdit,
        handleSuccess,
      }
    },
  })
</script>
