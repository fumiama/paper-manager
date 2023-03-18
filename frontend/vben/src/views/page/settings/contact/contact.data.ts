import { FormSchema } from '/@/components/Form'

export const formSchema: FormSchema[] = [
  {
    field: 'contactOld',
    label: '当前联系方式',
    component: 'Input',
    required: true,
  },
  {
    field: 'contactNew',
    label: '新联系方式',
    component: 'Input',
    required: true,
  },
]
