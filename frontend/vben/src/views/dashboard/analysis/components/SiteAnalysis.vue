<template>
  <Card
    :tab-list="[
      {
        key: 'tab1',
        tab: '访问量',
      },
    ]"
    active-tab-key="tab1"
  >
    <VisitAnalysisBar />
  </Card>
</template>
<script lang="ts" setup>
  import { Card } from 'ant-design-vue'
  import VisitAnalysisBar from './VisitAnalysisBar.vue'
</script>
