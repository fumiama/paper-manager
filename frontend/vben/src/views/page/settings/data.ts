import { FormSchema } from '/@/components/Form/index'

export interface ListItem {
  key: string
  title: string
  description: string
  extra?: string
  avatar?: string
  color?: string
}

// tab的list
export const settingList = [
  {
    key: '1',
    name: '基本设置',
    component: 'BaseSetting',
  },
  {
    key: '2',
    name: '安全设置',
    component: 'SecureSetting',
  },
]

// 基础设置 form
export const baseSetschemas: FormSchema[] = [
  {
    field: 'name',
    component: 'Input',
    label: '昵称',
    colProps: { span: 18 },
  },
  {
    field: 'introduction',
    component: 'InputTextArea',
    label: '个人简介',
    colProps: { span: 18 },
  },
]

// 安全设置 list
export const secureSettingList: ListItem[] = [
  {
    key: '1',
    title: '账户密码',
    description: '上次修改密码: 2022年1月1日0时0分0秒',
    extra: '修改',
  },
  {
    key: '2',
    title: '我的手机',
    description: '已绑定手机: 138****8293',
    extra: '修改',
  },
]
