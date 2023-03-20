import { BasicColumn } from '/@/components/Table'
import { FormSchema } from '/@/components/Table'

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
    title: '创建时间',
    dataIndex: 'date',
    width: 180,
  },
  {
    title: '角色',
    dataIndex: 'role',
    width: 200,
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
