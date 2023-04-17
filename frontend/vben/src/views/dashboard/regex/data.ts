import { FormSchema } from '/@/components/Form'
const colProps = {
  span: 8,
}

export const schemas: FormSchema[] = [
  {
    field: 'title',
    component: 'Input',
    colProps,
    label: '标题',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'class',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'opencl',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'date',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'time',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'rate',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'major',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
  {
    field: 'sub',
    component: 'Input',
    colProps,
    label: '客户',
    helpMessage: '目标的服务对象',
    componentProps: {
      placeholder: '请描述你服务的客户，内部客户直接 @姓名／工号',
    },
  },
]
