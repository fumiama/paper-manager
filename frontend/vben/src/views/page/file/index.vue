<template>
  <PageWrapper :title="t('routes.filelist.file') + ': ' + docxNameRef">
    <template #headerContent>
      <a-button type="primary" @click="downloadDocx">
        下载试卷 ({{ docxSizeRef.toFixed(2) }}MB)
      </a-button>
    </template>

    <a-card title="分析报告" :bordered="false">
      <div ref="chartRef" :style="{ height, width }"></div>
    </a-card>

    <a-card title="原卷预览" :bordered="false" class="!mt-5">
      <div class="docxWrap" :style="{ width }">
        <div ref="docxRef"></div>
      </div>
    </a-card>
  </PageWrapper>
</template>
<script lang="ts">
  import { computed, defineComponent, unref, PropType, ref, Ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { PageWrapper } from '/@/components/Page'
  import { useECharts } from '/@/hooks/web/useECharts'
  import { renderAsync } from 'docx-preview'
  import { Card } from 'ant-design-vue'
  import { downloadFile, getFileStatus, getFileBlob } from '/@/api/page'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { useGo } from '/@/hooks/web/usePage'
  import { useTabs } from '/@/hooks/web/useTabs'
  import { PageEnum } from '/@/enums/pageEnum'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { downloadByData } from '/@/utils/file/download'

  const { t } = useI18n()

  let docxRef = ref(null)

  let docxNameRef = ref('paper.docx')
  let docxSizeRef = ref(0)

  let docxBlob: Blob | null = null

  function loadDocx(file: Blob) {
    docxBlob = file
    renderAsync(file, docxRef.value as unknown as HTMLElement, undefined, {
      className: 'docx', // 默认和文档样式类的类名/前缀
      inWrapper: false, // 启用围绕文档内容渲染包装器
      ignoreWidth: false, // 禁止页面渲染宽度
      ignoreHeight: false, // 禁止页面渲染高度
      ignoreFonts: false, // 禁止字体渲染
      breakPages: false, // 在分页符上启用分页
      ignoreLastRenderedPageBreak: true, //禁用lastRenderedPageBreak元素的分页
      experimental: true, // 启用实验性功能（制表符停止计算）
      trimXmlDeclaration: true, // 如果为真，xml声明将在解析之前从xml文档中删除
      debug: false, // 启用额外的日志记录
    })
  }

  function downloadDocx() {
    downloadByData(docxBlob as BlobPart, docxNameRef.value)
  }

  export default defineComponent({
    name: 'PaperAnalyzeTab',
    components: { PageWrapper, [Card.name]: Card },
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
      const { currentRoute } = useRouter()
      const { createMessage } = useMessage()
      const { closeCurrent } = useTabs()
      const go = useGo()

      const params = computed(() => {
        return unref(currentRoute).params
      })

      if (!params.value || !params.value.id) {
        go(PageEnum.ERROR_PAGE)
        closeCurrent()
      }

      ;(async () => {
        try {
          const ret = await downloadFile(Number(params.value.id))
          if (ret && ret.url) {
            const data = await getFileBlob(ret.url)
            if (data) {
              loadDocx(data)
              return
            }
          }
          go(PageEnum.ERROR_PAGE)
          closeCurrent()
        } catch (error) {
          createMessage.error('加载docx错误: ' + (error as unknown as Error).message)
          go(PageEnum.ERROR_PAGE)
          closeCurrent()
        }
      })()

      const chartRef = ref<HTMLDivElement | null>(null)
      const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>)

      ;(async () => {
        try {
          const ret = await getFileStatus(Number(params.value.id))
          if (ret && ret.duplications.length > 0 && ret.questions.length > 0) {
            docxNameRef.value = ret.name
            docxSizeRef.value = ret.size
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
                  text: '全库重复率: ' + ret.rate.toFixed(2) + '%, 前十如下',
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
            return
          }
          go(PageEnum.ERROR_PAGE)
          closeCurrent()
        } catch (error) {
          createMessage.error('加载分析数据错误: ' + (error as unknown as Error).message)
          go(PageEnum.ERROR_PAGE)
          closeCurrent()
        }
      })()

      return {
        t,
        chartRef,
        docxRef,
        downloadDocx,
        docxNameRef,
        docxSizeRef,
      }
    },
  })
</script>
<style lang="less" scoped>
  .docxWrap {
    padding-top: 0px;
    margin: 0 auto;
    overflow-x: auto;
    display: grid;
    align-items: center;
    justify-items: center;
  }
</style>
