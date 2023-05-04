import { FormSchema } from '/@/components/Form'

export const schemas: FormSchema[] = [
  {
    field: 'f1',
    component: 'Input',
    label: '试卷名',
    slot: 'localSearch',
    required: true,
  },
]
export const taskSchemas: FormSchema[] = [
  {
    field: 't1',
    component: 'RangePicker',
    label: '查询年限范围',
    required: true,
    componentProps: {
      format: 'YYYY',
      placeholder: ['起始年份', '结束年份'],
      showTime: { format: 'YYYY' },
    },
  },
]
