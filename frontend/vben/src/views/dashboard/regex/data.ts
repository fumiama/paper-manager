import { FormSchema } from '/@/components/Form'
const colProps = {
  span: 24,
}

export const schemas: FormSchema[] = [
  {
    field: 'title',
    component: 'Input',
    colProps,
    label: '试卷标题',
    defaultValue: 'more 吗',
  },
  {
    field: 'class',
    component: 'Input',
    colProps,
    label: '课程名称',
    defaultValue: 'more 吗',
  },
  {
    field: 'opencl',
    component: 'Input',
    colProps,
    label: '开/闭卷',
    defaultValue: 'more 吗',
  },
  {
    field: 'date',
    component: 'Input',
    colProps,
    label: '考试日期',
    defaultValue: 'more 吗',
  },
  {
    field: 'time',
    component: 'Input',
    colProps,
    label: '考试时长',
    defaultValue: 'more 吗',
  },
  {
    field: 'rate',
    component: 'Input',
    colProps,
    label: '成绩占比',
    defaultValue: 'more 吗',
  },
  {
    field: 'major',
    component: 'Input',
    colProps,
    label: '大题题号',
    defaultValue: 'more 吗',
  },
  {
    field: 'sub',
    component: 'Input',
    colProps,
    label: '小题题号',
    defaultValue: 'more 吗',
  },
]
