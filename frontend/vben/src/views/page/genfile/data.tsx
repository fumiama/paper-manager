import { FormSchema } from '/@/components/Form'

export const step1Schemas: FormSchema[] = [
  {
    field: 'RateLimit',
    component: 'Slider',
    label: '重复率上限',
    required: true,
    componentProps: {
      min: 0,
      max: 100,
      range: true,
      marks: {
        0: '0%',
        50: '50%',
        100: '100%',
      },
    },
    colProps: {
      span: 24,
    },
  },
  {
    field: '[YearStart, YearEnd]',
    component: 'RangePicker',
    label: '起止年份',
    required: false,
    componentProps: {
      format: 'YYYY',
      placeholder: ['起始年份', '结束年份'],
      showTime: { format: 'YYYY' },
    },
    colProps: {
      span: 12,
    },
  },
  {
    field: 'AB',
    component: 'RadioGroup',
    label: '试卷类别',
    componentProps: {
      options: [
        {
          label: 'A',
          value: 'A',
        },
        {
          label: 'B',
          value: 'B',
        },
      ],
    },
    colProps: {
      span: 12,
    },
  },
  {
    field: 'MiddleFinal',
    component: 'RadioGroup',
    label: '考试阶段',
    componentProps: {
      options: [
        {
          label: '期中',
          value: '中',
        },
        {
          label: '期末',
          value: '末',
        },
      ],
    },
    colProps: {
      span: 12,
    },
  },
  {
    field: 'FirstSecond',
    component: 'RadioGroup',
    label: '考试学期',
    componentProps: {
      options: [
        {
          label: '第1学期',
          value: '1',
        },
        {
          label: '第2学期',
          value: '2',
        },
      ],
    },
    colProps: {
      span: 12,
    },
  },
  {
    field: 'OpenClose',
    component: 'RadioGroup',
    label: '考试类型',
    componentProps: {
      options: [
        {
          label: '开卷',
          value: '开卷',
        },
        {
          label: '一页纸开卷',
          value: '一页纸开卷',
        },
        {
          label: '闭卷',
          value: '闭卷',
        },
      ],
    },
    colProps: {
      span: 20,
    },
  },
]

export const step2Schemas: FormSchema[] = [
  {
    field: 'pwd',
    component: 'InputPassword',
    label: '支付密码',
    required: true,
    defaultValue: '123456',
    colProps: {
      span: 24,
    },
  },
]
