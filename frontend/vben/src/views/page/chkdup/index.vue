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
    <a-card title="查重报告" :bordered="false" class="!mt-5" v-if="tableRef && tableRef.length > 0">
      <div ref="chartRef" :style="{ height, width }"></div>
      <BasicTable
        title="详细信息"
        :columns="columns"
        :dataSource="tableRef"
        :canResize="canResize"
        :loading="loading"
        :striped="striped"
        :bordered="border"
        showTableSetting
        :pagination="pagination"
      >
        <template #toolbar>
          <a-button type="primary" @click="toggleCanResize">
            {{ !canResize ? '自适应高度' : '取消自适应' }}
          </a-button>
          <a-button type="primary" @click="toggleBorder">
            {{ !border ? '显示边框' : '隐藏边框' }}
          </a-button>
          <a-button type="primary" @click="toggleStriped">
            {{ !striped ? '显示斑马纹' : '隐藏斑马纹' }}
          </a-button>
        </template>
      </BasicTable>
    </a-card>

    <template #rightFooter>
      <a-button type="primary" @click="submitAll"> 提交 </a-button>
    </template>
  </PageWrapper>
</template>
<script lang="ts">
  import { BasicForm, useForm, ApiSelect } from '/@/components/Form'
  import { BasicTable } from '/@/components/Table'
  import { defineComponent, ref, unref, computed, Ref } from 'vue'
  import { useDebounceFn } from '@vueuse/core'
  import { useECharts } from '/@/hooks/web/useECharts'
  import { PageWrapper } from '/@/components/Page'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { schemas, taskSchemas, columns } from './data'
  import { Card } from 'ant-design-vue'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { getFileOptions, checkFileDup } from '/@/api/page'

  export default defineComponent({
    name: 'FormHightPage',
    components: { ApiSelect, BasicForm, BasicTable, PageWrapper, [Card.name]: Card },
    props: {
      width: {
        type: String as PropType<string>,
        default: '100%',
      },
      height: {
        type: String as PropType<string>,
        default: 'calc(100vh - 78px)',
      },
    },
    setup() {
      const { t } = useI18n()
      const keyword = ref<string>('')
      const searchParams = computed<Recordable>(() => {
        return { keyword: unref(keyword) }
      })
      const chartRef = ref<HTMLDivElement | null>(null)
      const tableRef = ref<any>(null)
      const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>)
      const { createMessage } = useMessage()

      const canResize = ref(false)
      const loading = ref(false)
      const striped = ref(true)
      const border = ref(true)
      const pagination = ref<any>(false)
      function toggleCanResize() {
        canResize.value = !canResize.value
      }
      function toggleStriped() {
        striped.value = !striped.value
      }
      function toggleBorder() {
        border.value = !border.value
      }

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
          const ret = await checkFileDup(
            Number(values.f1),
            taskValues.t1[0].$y,
            taskValues.t1[1].$y,
          )
          if (
            !(
              ret &&
              ret.duplications.length > 0 &&
              ret.questions.length > 0 &&
              ret.files.length > 0
            )
          ) {
            createMessage.warn('请先分析该试卷!')
            return
          }
          const barNames = ret.duplications.map((value) => {
            return value.name
          })
          const barData = ret.duplications.map((value) => {
            return value.percent
          })
          const queData = ret.questions.map((value) => {
            return { value: value.count, name: value.name }
          })
          const ptData = ret.questions.map((value) => {
            return { value: value.point, name: value.name }
          })
          tableRef.value = ret.files
          setOptions({
            title: [
              {
                text: '题量占比',
                left: '2%',
                top: '1%',
                textStyle: {
                  fontSize: 20,
                },
              },
              {
                text: '平均重复率: ' + ret.rate.toFixed(2) + '%, 前十如下',
                left: '40%',
                top: '1%',
                textStyle: {
                  fontSize: 20,
                },
              },
              {
                text: '分数占比',
                left: '2%',
                top: '50%',
                textStyle: {
                  fontSize: 20,
                },
              },
            ],
            grid: [{ left: '50%', top: '7%', width: '45%', height: '90%' }],
            tooltip: {
              formatter: '{b} ({c})',
            },
            xAxis: [
              {
                gridIndex: 0,
                axisTick: { show: false },
                axisLabel: { show: false },
                splitLine: { show: false },
                axisLine: { show: false },
              },
            ],
            yAxis: [
              {
                gridIndex: 0,
                interval: 0,
                data: barNames,
                axisTick: { show: false },
                axisLabel: { show: true },
                splitLine: { show: false },
                axisLine: { show: true },
              },
            ],
            series: [
              {
                name: '题量占比',
                type: 'pie',
                radius: '30%',
                center: ['22%', '25%'],
                data: queData,
                labelLine: { show: false },
                label: {
                  show: true,
                  formatter: function (d) {
                    return d.name + '(' + d.value + ')'
                  },
                },
              },
              {
                name: '分数占比',
                type: 'pie',
                radius: '30%',
                center: ['22%', '75%'],
                labelLine: { show: false },
                data: ptData,
                label: {
                  show: true,
                  formatter: '{b}\n　　({d}%)　',
                },
              },
              {
                name: '重复率前十',
                type: 'bar',
                xAxisIndex: 0,
                yAxisIndex: 0,
                barWidth: '45%',
                itemStyle: { color: '#86c9f4' },
                label: {
                  show: true,
                  position: 'right',
                  formatter: function (d) {
                    return d.data + '%'
                  },
                },
                data: barData,
              },
            ],
          })
        } catch (error) {
          createMessage.error('加载分析数据错误: ' + (error as unknown as Error).message)
        }
      }

      function onSearch(value: string) {
        keyword.value = value
      }

      return {
        t,
        columns,
        chartRef,
        tableRef,
        register,
        registerTask,
        submitAll,
        getFileOptions,
        searchParams,
        onSearch: useDebounceFn(onSearch, 300),
        canResize,
        loading,
        striped,
        border,
        toggleStriped,
        toggleCanResize,
        toggleBorder,
        pagination,
      }
    },
  })
</script>
<style lang="less" scoped>
  .high-form {
    padding-bottom: 48px;
  }
</style>
