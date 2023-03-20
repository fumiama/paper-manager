import { isNameExist } from '/@/api/sys/user'
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

export const nameFormSchema: FormSchema[] = [
  {
    field: 'name',
    label: '用户名',
    component: 'Input',
    helpMessage: ['不能输入带有admin的用户名'],
    rules: [
      {
        required: true,
        message: '请输入用户名',
      },
      {
        validator(_, value) {
          return new Promise((resolve, reject) => {
            isNameExist(value)
              .then((v) => {
                if (v) resolve()
                else reject('用户名已存在')
              })
              .catch((err) => {
                reject(err.message || '验证失败')
              })
          })
        },
      },
    ],
  },
  {
    field: 'pwd',
    label: '密码',
    component: 'InputPassword',
    required: true,
    ifShow: false,
  },
  {
    label: '角色',
    field: 'role',
    component: 'ApiSelect',
    componentProps: {
      api: () => {
        return [
          {
            roleName: '课程组长',
            value: 'super',
          },
          {
            roleName: '归档代理',
            value: 'filemgr',
          },
          {
            roleName: '课程组员',
            value: 'user',
          },
        ]
      },
      labelField: 'roleName',
      valueField: 'value',
    },
    required: true,
  },
  {
    field: 'nick',
    label: '昵称',
    component: 'Input',
    required: true,
  },
  {
    label: '简介',
    field: 'desc',
    component: 'InputTextArea',
  },
]
