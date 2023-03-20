<template>
  <PageWrapper dense contentFullHeight fixedHeight contentClass="flex">
    <BasicTable @register="registerTable" :searchInfo="searchInfo">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate">新增账号</a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑用户资料',
                onClick: handleEdit.bind(null, record),
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                tooltip: '删除此账号',
                popConfirm: {
                  title:
                    '为确保安全, 删除账号仅仅是禁止了登录, 必要情况下用户仍然可以通过忘记密码的方式找回账号',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
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
  import { defineComponent, reactive } from 'vue'

  import { BasicTable, useTable, TableAction } from '/@/components/Table'
  import { getUsersList } from '/@/api/sys/user'
  import { PageWrapper } from '/@/components/Page'

  import { useModal } from '/@/components/Modal'
  import AccountModal from './AccountModal.vue'

  import { columns, searchFormSchema } from './account.data'

  export default defineComponent({
    name: 'AccountManagement',
    components: { BasicTable, PageWrapper, AccountModal, TableAction },
    setup() {
      const [registerModal, { openModal }] = useModal()
      const searchInfo = reactive<Recordable>({})
      const [registerTable, { reload, updateTableDataRecord }] = useTable({
        title: '账号列表',
        api: getUsersList,
        rowKey: 'id',
        columns,
        formConfig: {
          labelWidth: 120,
          schemas: searchFormSchema,
          autoSubmitOnEnter: true,
        },
        useSearchForm: true,
        showTableSetting: true,
        bordered: true,
        handleSearchInfoFn(info) {
          console.log('handleSearchInfoFn', info)
          return info
        },
        actionColumn: {
          width: 120,
          title: '操作',
          dataIndex: 'action',
          // slots: { customRender: 'action' },
        },
      })

      function handleCreate() {
        openModal(true, {
          isUpdate: false,
        })
      }

      function handleEdit(record: Recordable) {
        console.log(record)
        openModal(true, {
          record,
          isUpdate: true,
        })
      }

      function handleDelete(record: Recordable) {
        console.log(record)
      }

      function handleSuccess({ isUpdate, values }) {
        if (isUpdate) {
          // 演示不刷新表格直接更新内部数据。
          // 注意：updateTableDataRecord要求表格的rowKey属性为string并且存在于每一行的record的keys中
          const result = updateTableDataRecord(values.id, values)
          console.log(result)
        } else {
          reload()
        }
      }

      function handleSelect(deptId = '') {
        searchInfo.deptId = deptId
        reload()
      }

      return {
        registerTable,
        registerModal,
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
        handleSelect,
        searchInfo,
      }
    },
  })
</script>
