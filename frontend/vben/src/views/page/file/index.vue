<template>
  <PageWrapper :title="t('routes.filelist.file')">
    <template #headerContent>
      <a-button type="primary"> 下载试卷 </a-button>
    </template>
    <div ref="chartRef" :style="{ height, width }"></div>
    <div class="docxWrap" :style="{ width }">
      <div ref="docxRef"></div>
    </div>
  </PageWrapper>
</template>
<script lang="ts">
  import { computed, defineComponent, unref, PropType, ref, Ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { PageWrapper } from '/@/components/Page'
  import { useECharts } from '/@/hooks/web/useECharts'
  import { renderAsync } from 'docx-preview'
  import { downloadFile } from '/@/api/page'
  import { DownloadFile } from '/@/api/page/model/fileModel'
  import { useGo } from '/@/hooks/web/usePage'
  import { PageEnum } from '/@/enums/pageEnum'
  import { useI18n } from '/@/hooks/web/useI18n'
  import axios from 'axios'

  const { t } = useI18n()

  let docxRef = ref(null)

  function loadDocx(file: Blob) {
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

  export default defineComponent({
    name: 'PaperAnalyzeTab',
    components: { PageWrapper },
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
      const go = useGo()

      const params = computed(() => {
        return unref(currentRoute).params
      })

      if (!params.value || !params.value.id) {
        go(PageEnum.ERROR_PAGE)
      }

      downloadFile(Number(params.value.id)).then((file: DownloadFile) => {
        if (file && file.url) {
          axios({
            method: 'get',
            responseType: 'blob',
            url: file.url,
          }).then(({ data }) => {
            loadDocx(data)
          })
        }
      })

      const chartRef = ref<HTMLDivElement | null>(null)
      const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>)
      const dataAll = [389, 259, 262, 324, 232, 176, 196, 214, 133, 370]
      const yAxisData = [
        '原因1',
        '原因2',
        '原因3',
        '原因4',
        '原因5',
        '原因6',
        '原因7',
        '原因8',
        '原因9',
        '原因10',
      ]
      onMounted(() => {
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
              text: '重复率前十',
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
              data: yAxisData.reverse(),
              axisTick: { show: false },
              axisLabel: { show: true },
              splitLine: { show: false },
              axisLine: { show: true },
            },
          ],
          series: [
            {
              name: '各渠道投诉占比',
              type: 'pie',
              radius: '30%',
              center: ['22%', '25%'],
              data: [
                { value: 335, name: '客服电话' },
                { value: 310, name: '奥迪官网' },
                { value: 234, name: '媒体曝光' },
                { value: 135, name: '质检总局' },
                { value: 105, name: '其他' },
              ],
              labelLine: { show: false },
              label: {
                show: true,
                formatter: '{b} \n ({d}%)',
              },
            },
            {
              name: '各级别投诉占比',
              type: 'pie',
              radius: '30%',
              center: ['22%', '75%'],
              labelLine: { show: false },
              data: [
                { value: 335, name: 'A级' },
                { value: 310, name: 'B级' },
                { value: 234, name: 'C级' },
                { value: 135, name: 'D级' },
              ],
              label: {
                show: true,
                formatter: '{b} \n ({d}%)',
              },
            },
            {
              name: '投诉原因TOP10',
              type: 'bar',
              xAxisIndex: 0,
              yAxisIndex: 0,
              barWidth: '45%',
              itemStyle: { color: '#86c9f4' },
              label: { show: true, position: 'right' },
              data: dataAll.sort(),
            },
          ],
        })
      })
      return {
        t,
        chartRef,
        docxRef,
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
