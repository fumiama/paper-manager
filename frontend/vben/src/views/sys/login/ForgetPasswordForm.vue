<template>
  <template v-if="getShow">
    <LoginFormTitle class="enter-x" />
    <Form class="p-4 enter-x" :model="formData" :rules="getFormRules" ref="formRef">
      <FormItem name="account" class="enter-x">
        <Input
          size="large"
          v-model:value="formData.account"
          :placeholder="t('sys.login.userName')"
        />
      </FormItem>

      <FormItem name="mobile" class="enter-x">
        <Input size="large" v-model:value="formData.mobile" :placeholder="t('sys.login.mobile')" />
      </FormItem>

      <FormItem class="enter-x">
        <Button
          type="primary"
          size="large"
          block
          @click="handleReset"
          :loading="loading"
          :disabled="!(formData.account && formData.mobile)"
        >
          {{ t('common.resetText') }}
        </Button>
        <Button size="large" block class="mt-4" @click="handleBackLogin">
          {{ t('sys.login.backSignIn') }}
        </Button>
      </FormItem>
    </Form>
  </template>
</template>
<script lang="ts" setup>
  import { reactive, ref, computed, unref } from 'vue'
  import LoginFormTitle from './LoginFormTitle.vue'
  import { Form, Input, Button } from 'ant-design-vue'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { useLoginState, useFormRules, LoginStateEnum, useFormValid } from './useLogin'
  import { resetPasswordApi } from '/@/api/sys/user'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { ResetPasswordParams } from '/@/api/sys/model/userModel'

  const FormItem = Form.Item
  const { t } = useI18n()
  const { handleBackLogin, getLoginState } = useLoginState()
  const { notification } = useMessage()
  const { getFormRules } = useFormRules()

  const formRef = ref()
  const loading = ref(false)

  const formData = reactive({
    account: '',
    mobile: '',
    // sms: '',
  })

  const { validForm } = useFormValid(formRef)

  const getShow = computed(() => unref(getLoginState) === LoginStateEnum.RESET_PASSWORD)

  async function handleReset() {
    const form = unref(formRef)
    if (!form) return
    const data = await validForm()
    if (!data) return
    try {
      loading.value = true
      const { msg } = await resetPasswordApi({
        username: data.account,
        phonenum: data.mobile,
      } as ResetPasswordParams)
      notification.info({
        message: t('sys.login.forgetFormTitle'),
        description: msg,
        duration: 10,
      })
    } catch (error) {
      notification.error({
        message: (error as Error).name,
        description: (error as Error).message,
        duration: 3,
      })
    } finally {
      loading.value = false
      form.resetFields()
      handleBackLogin()
    }
  }
</script>
