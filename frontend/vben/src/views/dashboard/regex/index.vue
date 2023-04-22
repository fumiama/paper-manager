<template>
  <PageWrapper
    :title="t('routes.dashboard.regex')"
    contentBackground
    content="设置试卷解析/查重时使用的正则表达式"
    contentClass="p-4"
  >
    <BasicForm @register="register" />
  </PageWrapper>
</template>
<script lang="ts">
  import { BasicForm, useForm } from '/@/components/Form'
  import { defineComponent } from 'vue'
  import { schemas } from './data'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { PageWrapper } from '/@/components/Page'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { setUserRegex } from '/@/api/dashboard'

  export default defineComponent({
    name: 'FormBasicPage',
    components: { BasicForm, PageWrapper },
    setup() {
      const { t } = useI18n()
      const { createMessage } = useMessage()
      const [register, { validate, setProps }] = useForm({
        labelCol: {
          span: 5,
        },
        wrapperCol: {
          span: 16,
        },
        schemas: schemas,
        actionColOptions: {
          offset: 3,
          span: 16,
        },
        submitButtonOptions: {
          text: '提交',
        },
        submitFunc: customSubmitFunc,
      })

      async function customSubmitFunc() {
        try {
          const data = await validate()
          if (!data) return
          setProps({
            submitButtonOptions: {
              loading: true,
            },
          })
          const msg = await setUserRegex(data)
          setProps({
            submitButtonOptions: {
              loading: false,
            },
          })
          createMessage.success(msg)
        } catch (error) {}
      }

      return { t, register }
    },
  })
</script>
<style lang="less" scoped>
  .form-wrap {
    padding: 24px;
    background-color: @component-background;
  }
</style>
