<template>
  <template v-if="getShow">
    <LoginFormTitle class="enter-x" />
    <Form class="p-4 enter-x" :model="formData" :rules="getFormRules" ref="formRef">
      <FormItem name="account" class="enter-x">
        <Input
          class="fix-auto-fill"
          size="large"
          v-model:value="formData.account"
          :placeholder="t('sys.login.userName')"
        />
      </FormItem>
      <FormItem name="mobile" class="enter-x">
        <Input
          size="large"
          v-model:value="formData.mobile"
          :placeholder="t('sys.login.mobile')"
          class="fix-auto-fill"
        />
      </FormItem>
      <!--<FormItem name="sms" class="enter-x">
        <CountdownInput
          size="large"
          class="fix-auto-fill"
          v-model:value="formData.sms"
          :placeholder="t('sys.login.smsCode')"
        />
      </FormItem>-->
      <FormItem name="password" class="enter-x">
        <StrengthMeter
          size="large"
          v-model:value="formData.password"
          :placeholder="t('sys.login.password')"
        />
      </FormItem>
      <FormItem name="confirmPassword" class="enter-x">
        <InputPassword
          size="large"
          visibilityToggle
          v-model:value="formData.confirmPassword"
          :placeholder="t('sys.login.confirmPassword')"
        />
      </FormItem>

      <FormItem class="enter-x" name="policy">
        <!-- No logic, you need to deal with it yourself -->
        <Checkbox v-model:checked="formData.policy" size="small">
          {{ t('sys.login.policy') }}
        </Checkbox>
      </FormItem>

      <Button
        type="primary"
        class="enter-x"
        size="large"
        block
        @click="handleRegister"
        :loading="loading"
        :disabled="!isFormDataFull()"
      >
        {{ t('sys.login.registerButton') }}
      </Button>
      <Button size="large" block class="mt-4 enter-x" @click="handleBackLogin">
        {{ t('sys.login.backSignIn') }}
      </Button>
    </Form>
  </template>
</template>
<script lang="ts" setup>
  import { reactive, ref, unref, computed } from 'vue'
  import LoginFormTitle from './LoginFormTitle.vue'
  import { Form, Input, Button, Checkbox } from 'ant-design-vue'
  import { StrengthMeter } from '/@/components/StrengthMeter'
  // import { CountdownInput } from '/@/components/CountDown'
  import { registerApi } from '/@/api/sys/user'
  import { RegisterParams } from '/@/api/sys/model/userModel'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { useLoginState, useFormRules, useFormValid, LoginStateEnum } from './useLogin'
  import { useMessage } from '/@/hooks/web/useMessage'

  const FormItem = Form.Item
  const InputPassword = Input.Password
  const { t } = useI18n()
  const { handleBackLogin, getLoginState } = useLoginState()

  const formRef = ref()
  const loading = ref(false)

  const formData = reactive({
    account: '',
    password: '',
    confirmPassword: '',
    mobile: '',
    // sms: '',
    policy: false,
  })

  const { getFormRules } = useFormRules(formData)
  const { validForm } = useFormValid(formRef)

  const getShow = computed(() => unref(getLoginState) === LoginStateEnum.REGISTER)

  const { notification } = useMessage()

  async function handleRegister() {
    const form = unref(formRef)
    if (!form) return
    const data = await validForm()
    if (!data) return
    try {
      loading.value = true
      const { msg } = await registerApi({
        username: data.account,
        mobile: data.mobile,
        password: data.password,
      } as RegisterParams)
      notification.info({
        message: t('sys.login.signUpFormTitle'),
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

  function isFormDataFull(): boolean {
    return (
      formData.account != '' &&
      formData.password != '' &&
      formData.confirmPassword != '' &&
      formData.mobile != '' &&
      formData.policy
    )
  }
</script>
