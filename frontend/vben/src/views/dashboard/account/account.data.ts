import { BasicColumn } from '/@/components/Table'
import { FormSchema } from '/@/components/Table'
import { h } from 'vue'
import { Switch } from 'ant-design-vue'
import { useMessage } from '/@/hooks/web/useMessage'
import { disableUser } from '/@/api/sys/user'

export const columns: BasicColumn[] = [
  {
    title: '用户名',
    dataIndex: 'name',
    width: 120,
  },
  {
    title: '昵称',
    dataIndex: 'nick',
    width: 120,
  },
  {
    title: '状态',
    dataIndex: 'stat',
    width: 120,
    customRender: ({ record }) => {
      if (!Reflect.has(record, 'pendingStatus')) {
        record.pendingStatus = false
      }
      return h(Switch, {
        checked: record.stat,
        checkedChildren: '已启用',
        unCheckedChildren: '已禁用',
        loading: record.pendingStatus,
        onChange(checked: boolean) {
          const { createMessage } = useMessage()
          if (checked) {
            record.stat = false
            createMessage.error('请让用户通过找回密码启用账户')
            return
          }
          record.pendingStatus = true
          disableUser(record.id, checked)
            .then(() => {
              record.stat = checked
              createMessage.success(`已成功禁用账户并清空密码, 如需重新启用, 请让用户找回密码`)
            })
            .catch((error) => {
              createMessage.error('禁用失败: ' + (error as unknown as Error).message)
            })
            .finally(() => {
              record.pendingStatus = false
            })
        },
      })
    },
  },
  {
    title: '创建时间',
    dataIndex: 'date',
    width: 240,
  },
  {
    title: '角色',
    dataIndex: 'role',
    width: 120,
  },
  {
    title: '简介',
    dataIndex: 'desc',
  },
]

export const searchFormSchema: FormSchema[] = [
  {
    field: 'name',
    label: '用户名',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    field: 'nick',
    label: '昵称',
    component: 'Input',
    colProps: { span: 8 },
  },
]
