<template>
  <PageWrapper
    :title="t('routes.genfile.name')"
    contentBackground
    content="使用自定义限制条件生成试卷"
    contentClass="p-4"
  >
    <div class="step-form-form">
      <a-steps :current="current">
        <a-step title="填写信息" />
        <a-step title="确认生成" />
        <a-step title="下载" />
      </a-steps>
    </div>
    <div class="mt-5">
      <Step1 @next="handleStep1Next" v-show="current === 0" />
      <Step2
        @prev="handleStepPrev"
        @next="handleStep2Next"
        v-show="current === 1"
        v-if="initSetp2"
      />
      <Step3 v-show="current === 2" @redo="handleRedo" v-if="initSetp3" />
    </div>
  </PageWrapper>
</template>
<script lang="ts">
  import { defineComponent, ref, toRefs } from 'vue'
  import { state } from './data'
  import Step1 from './Step1.vue'
  import Step2 from './Step2.vue'
  import Step3 from './Step3.vue'
  import { PageWrapper } from '/@/components/Page'
  import { Steps } from 'ant-design-vue'
  import { useI18n } from '/@/hooks/web/useI18n'

  export default defineComponent({
    name: 'FormStepPage',
    components: {
      Step1,
      Step2,
      Step3,
      PageWrapper,
      [Steps.name]: Steps,
      [Steps.Step.name]: Steps.Step,
    },
    setup() {
      const current = ref(0)

      const { t } = useI18n()

      function handleStep1Next(step1Values: any) {
        current.value++
        state.initSetp2 = true
        state.step1Values = step1Values
      }

      function handleStepPrev() {
        current.value--
      }

      function handleStep2Next(step2Values: any) {
        current.value++
        state.initSetp3 = true
        state.step2Values = step2Values
      }

      function handleRedo() {
        current.value = 0
        state.initSetp2 = false
        state.initSetp3 = false
      }

      return {
        t,
        current,
        handleStep1Next,
        handleStep2Next,
        handleRedo,
        handleStepPrev,
        ...toRefs(state),
      }
    },
  })
</script>
<style lang="less" scoped>
  .step-form-content {
    padding: 24px;
    background-color: @component-background;
  }

  .step-form-form {
    width: 750px;
    margin: 0 auto;
  }
</style>
