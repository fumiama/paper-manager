<template>
  <div ref="chartRef" :style="{ height, width }"></div>
</template>
<script lang="ts" setup>
  import { onMounted, ref, Ref } from 'vue'
  import { useECharts } from '/@/hooks/web/useECharts'
  import { getAnnualVisits } from '/@/api/dashboard'

  defineProps({
    width: {
      type: String as PropType<string>,
      default: '100%',
    },
    height: {
      type: String as PropType<string>,
      default: '280px',
    },
  })

  const chartRef = ref<HTMLDivElement | null>(null)
  const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>)
  const visitsRef = ref([...new Array(12)])
  getAnnualVisits().then((visits) => {
    visitsRef.value = visits
  })
  onMounted(() => {
    setOptions({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          lineStyle: {
            width: 1,
            color: '#019680',
          },
        },
      },
      grid: { left: '1%', right: '1%', top: '2%', bottom: 0, containLabel: true },
      xAxis: {
        type: 'category',
        data: [...new Array(12)].map((_item, index) => `${index + 1}月`),
      },
      yAxis: {
        type: 'value',
        splitNumber: 4,
      },
      series: [
        {
          data: visitsRef as any,
          type: 'bar',
          barMaxWidth: 80,
        },
      ],
    })
  })
</script>
