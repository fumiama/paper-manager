<template>
  <div class="step2">
    <a-alert message="确认提交后，将进入下载页面，云端不保存。" show-icon />
    <a-descriptions :column="1" class="mt-5">
      <a-descriptions-item label="大题数">
        {{ state.step1Values.data.length }}
      </a-descriptions-item>
      <a-descriptions-item label="小题数">
        {{
          state.step1Values.data
            .map((record) => record.count || 0)
            .reduce((sum, count) => sum + count, 0)
        }}
      </a-descriptions-item>
      <a-descriptions-item label="重复率上限">
        {{ state.step1Values.values.RateLimit[1] / 100 }}
      </a-descriptions-item>
      <a-descriptions-item label="起止年份" v-if="state.step1Values.values['[YearStart, YearEnd]']">
        {{ state.step1Values.values['[YearStart, YearEnd]'][0].$y }} -
        {{ state.step1Values.values['[YearStart, YearEnd]'][1].$y }}
      </a-descriptions-item>
      <a-descriptions-item label="试卷类别" v-if="state.step1Values.values.AB">
        {{ state.step1Values.values.AB }}
      </a-descriptions-item>
      <a-descriptions-item label="考试阶段" v-if="state.step1Values.values.MiddleFinal">
        {{ state.step1Values.values.MiddleFinal }}
      </a-descriptions-item>
      <a-descriptions-item label="考试学期" v-if="state.step1Values.values.FirstSecond">
        {{ state.step1Values.values.FirstSecond }}
      </a-descriptions-item>
      <a-descriptions-item label="考试类型" v-if="state.step1Values.values.OpenClose">
        {{ state.step1Values.values.OpenClose }}
      </a-descriptions-item>
    </a-descriptions>
    <a-divider />
    <BasicForm @register="register" />
  </div>
</template>
<script lang="ts">
  import { defineComponent } from 'vue'
  import { BasicForm, useForm } from '/@/components/Form'
  import { state } from './data'
  import { generateFile } from '/@/api/page'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { Alert, Divider, Descriptions } from 'ant-design-vue'

  export default defineComponent({
    components: {
      BasicForm,
      [Alert.name]: Alert,
      [Divider.name]: Divider,
      [Descriptions.name]: Descriptions,
      [Descriptions.Item.name]: Descriptions.Item,
    },
    emits: ['next', 'prev'],
    setup(_, { emit }) {
      const [register, { setProps }] = useForm({
        labelWidth: 80,
        actionColOptions: {
          span: 16,
        },
        resetButtonOptions: {
          text: '上一步',
        },
        submitButtonOptions: {
          text: '提交',
        },
        resetFunc: customResetFunc,
        submitFunc: customSubmitFunc,
      })

      const { createMessage } = useMessage()

      async function customResetFunc() {
        emit('prev')
      }

      async function customSubmitFunc() {
        try {
          setProps({
            submitButtonOptions: {
              loading: true,
            },
          })
          let ys = 0
          let ye = 0
          if (state.step1Values.values['[YearStart, YearEnd]']) {
            ys = state.step1Values.values['[YearStart, YearEnd]'][0].$y
            ye = state.step1Values.values['[YearStart, YearEnd]'][1].$y
          }
          let tm = 0
          if (state.step1Values.values.AB) {
            if (state.step1Values.values.AB == 'A') tm |= 1
            else if (state.step1Values.values.AB == 'B') tm |= 2
          }
          if (state.step1Values.values.MiddleFinal) {
            if (state.step1Values.values.MiddleFinal == '中') tm |= 1 << 4
            else if (state.step1Values.values.MiddleFinal == '末') tm |= 2 << 4
          }
          if (state.step1Values.values.FirstSecond) {
            if (state.step1Values.values.FirstSecond == '1') tm |= 1 << 8
            else if (state.step1Values.values.FirstSecond == '2') tm |= 2 << 8
          }
          if (state.step1Values.values.OpenClose) {
            if (state.step1Values.values.OpenClose == '开卷') tm |= 1 << 12
            else if (state.step1Values.values.OpenClose == '一页纸开卷') tm |= 2 << 12
            else if (state.step1Values.values.OpenClose == '闭卷') tm |= 4 << 12
          }
          const data = await generateFile({
            Distribution: state.step1Values.data.reduce((acc, { major, count }) => {
              console.log(major, count)
              acc[major] = count
              return acc
            }, {}),
            RateLimit: state.step1Values.values.RateLimit[1] / 100,
            YearStart: ys,
            YearEnd: ye,
            TypeMask: tm,
          })
          setProps({
            submitButtonOptions: {
              loading: false,
            },
          })
          emit('next', data)
        } catch (error) {
          createMessage.error((error as unknown as Error).message)
          setProps({
            submitButtonOptions: {
              loading: false,
            },
          })
        }
      }

      return {
        register,
        state,
      }
    },
  })
</script>
<style lang="less" scoped>
  .step2 {
    width: 450px;
    margin: 0 auto;
  }
</style>
