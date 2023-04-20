import { FormSchema } from '/@/components/Form'
import { getUserRegex } from '/@/api/dashboard'
const colProps = {
  span: 24,
}

const userregex = await getUserRegex()

export const schemas: FormSchema[] = [
  {
    field: 'title',
    component: 'Input',
    colProps,
    label: '试卷标题',
    defaultValue: userregex.Title,
  },
  {
    field: 'class',
    component: 'Input',
    colProps,
    label: '课程名称',
    defaultValue: userregex.Class,
  },
  {
    field: 'opencl',
    component: 'Input',
    colProps,
    label: '开/闭卷',
    defaultValue: userregex.OpenCl,
  },
  {
    field: 'date',
    component: 'Input',
    colProps,
    label: '考试日期',
    defaultValue: userregex.Date,
  },
  {
    field: 'time',
    component: 'Input',
    colProps,
    label: '考试时长',
    defaultValue: userregex.Time,
  },
  {
    field: 'rate',
    component: 'Input',
    colProps,
    label: '成绩占比',
    defaultValue: userregex.Rate,
  },
  {
    field: 'major',
    component: 'Input',
    colProps,
    label: '大题题号',
    defaultValue: userregex.Major,
  },
  {
    field: 'sub',
    component: 'Input',
    colProps,
    label: '小题题号',
    defaultValue: userregex.Sub,
  },
]
