import { FormSchema } from '/@/components/Form'
import { BasicColumn } from '/@/components/Table/src/types/table'

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
export const columns: BasicColumn[] = [
  {
    title: '试卷',
    dataIndex: 'name',
    fixed: 'left',
    width: 200,
  },
  {
    title: '重复率(%)',
    dataIndex: 'percent',
    width: 150,
  },
]
