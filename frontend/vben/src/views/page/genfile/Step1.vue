<template>
  <div class="step1">
    <div class="step1-form">
      <BasicTable @register="registerTable">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <TableAction :actions="createActions(record, column)" />
          </template>
        </template>
      </BasicTable>
      <a-button block type="dashed" class="mt-5" @click="handleAdd"> 新增大题 </a-button>
      <a-divider />
      <BasicForm @register="register" />
    </div>
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue'
  import { BasicForm, useForm } from '/@/components/Form'
  import { getMajors } from '/@/api/page'
  import {
    BasicTable,
    useTable,
    TableAction,
    BasicColumn,
    ActionItem,
    EditRecordRow,
  } from '/@/components/Table'
  import { step1Schemas } from './data'
  import { Select, Input, Divider } from 'ant-design-vue'
  import { useMessage } from '/@/hooks/web/useMessage'

  const columns: BasicColumn[] = [
    {
      title: '大题',
      dataIndex: 'major',
      editRow: true,
      editComponent: 'ApiSelect',
      editComponentProps: {
        api: getMajors,
        resultField: 'Name',
        labelField: 'Name',
        valueField: 'Name',
      },
    },
    {
      title: '数量',
      dataIndex: 'count',
      editRow: true,
      editComponent: 'InputNumber',
      width: 120,
    },
  ]

  export default defineComponent({
    components: {
      BasicForm,
      BasicTable,
      TableAction,
      [Select.name]: Select,
      [Input.name]: Input,
      [Input.Group.name]: Input.Group,
      [Divider.name]: Divider,
    },
    emits: ['next'],
    setup(_, { emit }) {
      const [registerTable, { getDataSource }] = useTable({
        columns: columns,
        showIndexColumn: false,
        actionColumn: {
          width: 100,
          title: '操作',
          dataIndex: 'action',
          // slots: { customRender: 'action' },
        },
        scroll: { y: '100%' },
        pagination: false,
      })

      const [register, { validate }] = useForm({
        labelWidth: 100,
        schemas: step1Schemas,
        actionColOptions: {
          span: 14,
        },
        showResetButton: false,
        submitButtonOptions: {
          text: '下一步',
        },
        submitFunc: customSubmitFunc,
      })

      const { createMessage } = useMessage()

      async function customSubmitFunc() {
        try {
          const values = await validate()
          emit('next', values)
        } catch (error) {}
      }

      function handleEdit(record: EditRecordRow) {
        record.onEdit?.(true)
      }

      function handleCancel(record: EditRecordRow) {
        record.onEdit?.(false)
        if (record.isNew) {
          const data = getDataSource()
          const index = data.findIndex((item) => item.key === record.key)
          data.splice(index, 1)
        }
      }

      function handleDelete(record: EditRecordRow) {
        record.onEdit?.(false)
        if (record.isNew) {
          const data = getDataSource()
          const index = data.findIndex((item) => item.key === record.key)
          data.splice(index, 1)
        }
      }

      function handleSave(record: EditRecordRow) {
        record.onEdit?.(false, true)
      }

      function handleAdd() {
        const data = getDataSource()
        if (data.length >= 10) {
          createMessage.warn('最大只支持10项!')
          return
        }
        const addRow: EditRecordRow = {
          major: '',
          count: 1,
          editable: true,
          isNew: true,
          key: `${Date.now()}`,
        }
        data.push(addRow)
      }

      function createActions(record: EditRecordRow, column: BasicColumn): ActionItem[] {
        if (!record.editable) {
          return [
            {
              label: '编辑',
              onClick: handleEdit.bind(null, record),
            },
            {
              label: '删除',
              onClick: handleDelete.bind(null, record, column),
            },
          ]
        }
        return [
          {
            label: '保存',
            onClick: handleSave.bind(null, record, column),
          },
          {
            label: '取消',
            popConfirm: {
              title: '是否取消编辑',
              confirm: handleCancel.bind(null, record, column),
            },
          },
        ]
      }

      return {
        register,
        registerTable,
        handleEdit,
        createActions,
        handleAdd,
        getDataSource,
      }
    },
  })
</script>
<style lang="less" scoped>
  .step1 {
    &-form {
      width: 600px;
      margin: 0 auto;
    }

    h3 {
      margin: 0 0 12px;
      font-size: 16px;
      line-height: 32px;
      color: @text-color;
    }

    h4 {
      margin: 0 0 4px;
      font-size: 14px;
      line-height: 22px;
      color: @text-color;
    }

    p {
      color: @text-color;
    }
  }

  .pay-select {
    width: 20%;
  }

  .pay-input {
    width: 70%;
  }
</style>
