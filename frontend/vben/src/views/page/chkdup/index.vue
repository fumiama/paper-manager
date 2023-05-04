<template>
  <PageWrapper
    class="high-form"
    :title="t('routes.templist.dup')"
    content="设定好待查重试卷和查重年限后，即可生成查重报告。"
  >
    <a-card title="待查重试卷" :bordered="false">
      <BasicForm @register="register">
        <template #localSearch="{ model, field }">
          <ApiSelect
            :api="getFileOptions"
            showSearch
            v-model:value="model[field]"
            optionFilterProp="label"
            resultField="value"
            labelField="label"
            valueField="value"
            :params="searchParams"
            @search="onSearch"
          />
        </template>
      </BasicForm>
    </a-card>
    <a-card title="查重限定条件" :bordered="false" class="!mt-5">
      <BasicForm @register="registerTask" />
    </a-card>
    <a-card title="查重报告" :bordered="false" class="!mt-5">
      <p> aaaaa </p>
    </a-card>

    <template #rightFooter>
      <a-button type="primary" @click="submitAll"> 提交 </a-button>
    </template>
  </PageWrapper>
</template>
<script lang="ts">
  import { BasicForm, useForm, ApiSelect } from '/@/components/Form'
  import { defineComponent, ref, unref, computed } from 'vue'
  import { useDebounceFn } from '@vueuse/core'
  import { PageWrapper } from '/@/components/Page'
  import { schemas, taskSchemas } from './data'
  import { Card } from 'ant-design-vue'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { getFileOptions } from '/@/api/page'

  export default defineComponent({
    name: 'FormHightPage',
    components: { ApiSelect, BasicForm, PageWrapper, [Card.name]: Card },
    setup() {
      const { t } = useI18n()
      const keyword = ref<string>('')
      const searchParams = computed<Recordable>(() => {
        return { keyword: unref(keyword) }
      })

      const [register, { validate }] = useForm({
        layout: 'vertical',
        baseColProps: {
          span: 24,
        },
        schemas: schemas,
        showActionButtonGroup: false,
      })

      const [registerTask, { validate: validateTaskForm }] = useForm({
        layout: 'vertical',
        baseColProps: {
          span: 6,
        },
        schemas: taskSchemas,
        showActionButtonGroup: false,
      })

      async function submitAll() {
        try {
          const [values, taskValues] = await Promise.all([validate(), validateTaskForm()])
          console.log('form data:', values, taskValues)
        } catch (error) {}
      }

      function onSearch(value: string) {
        keyword.value = value
      }

      return {
        t,
        register,
        registerTask,
        submitAll,
        getFileOptions,
        searchParams,
        onSearch: useDebounceFn(onSearch, 300),
      }
    },
  })
</script>
<style lang="less" scoped>
  .high-form {
    padding-bottom: 48px;
  }
</style>
