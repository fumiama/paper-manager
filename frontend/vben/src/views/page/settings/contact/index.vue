<template>
  <PageWrapper title="修改当前联系方式" content="修改成功后会通知课程组长！">
    <div class="py-8 bg-white flex flex-col justify-center items-center">
      <BasicForm @register="register" />
      <div class="flex justify-center">
        <a-button @click="resetFields"> 重置 </a-button>
        <a-button class="!ml-4" type="primary" @click="handleSubmit"> 确认 </a-button>
      </div>
    </div>
  </PageWrapper>
</template>
<script lang="ts">
  import { defineComponent } from 'vue'
  import { PageWrapper } from '/@/components/Page'
  import { BasicForm, useForm } from '/@/components/Form'
  import { setContactApi } from '/@/api/sys/user'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { useUserStore } from '/@/store/modules/user'
  import md5 from 'md5'

  import { formSchema } from './contact.data'
  export default defineComponent({
    name: 'ChangeContact',
    components: { BasicForm, PageWrapper },
    setup() {
      const [register, { validate, resetFields }] = useForm({
        size: 'large',
        baseColProps: { span: 24 },
        labelWidth: 100,
        showActionButtonGroup: false,
        schemas: formSchema,
      })

      const { createMessage } = useMessage()
      const userStore = useUserStore()

      async function handleSubmit() {
        try {
          const values = await validate()
          const { contactOld, contactNew } = values
          const { msg } = await setContactApi({
            token: md5(contactOld + contactNew),
            contact: contactNew,
          })
          createMessage.success(msg)
          userStore.getUserInfoAction()
        } catch (_) {}
      }

      return { register, resetFields, handleSubmit }
    },
  })
</script>
