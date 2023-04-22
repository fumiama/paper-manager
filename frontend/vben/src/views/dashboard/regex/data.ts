import { FormSchema } from '/@/components/Form'
import { getUserRegex } from '/@/api/dashboard'
const colProps = {
  span: 24,
}

const userregex = await getUserRegex()

export const schemas: FormSchema[] = [
  {
    field: 'Title',
    component: 'Input',
    colProps,
    label: '试卷标题',
    defaultValue: userregex.Title,
  },
  {
    field: 'Class',
    component: 'Input',
    colProps,
    label: '课程名称',
    defaultValue: userregex.Class,
  },
  {
    field: 'OpenCl',
    component: 'Input',
    colProps,
    label: '开/闭卷',
    defaultValue: userregex.OpenCl,
  },
  {
    field: 'Date',
    component: 'Input',
    colProps,
    label: '考试日期',
    defaultValue: userregex.Date,
  },
  {
    field: 'Time',
    component: 'Input',
    colProps,
    label: '考试时长',
    defaultValue: userregex.Time,
  },
  {
    field: 'Rate',
    component: 'Input',
    colProps,
    label: '成绩占比',
    defaultValue: userregex.Rate,
  },
  {
    field: 'Major',
    component: 'Input',
    colProps,
    label: '大题题号',
    defaultValue: userregex.Major,
  },
  {
    field: 'Sub',
    component: 'Input',
    colProps,
    label: '小题题号',
    defaultValue: userregex.Sub,
  },
]
