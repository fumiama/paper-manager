<template>
  <PageWrapper>
    <template #headerContent> <WorkbenchHeader /> </template>
    <div class="lg:flex">
      <div class="w-full">
        <DynamicInfo :loading="loading" class="enter-y" />
      </div>
    </div>
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { ref } from 'vue'
  import { PageWrapper } from '/@/components/Page'
  import WorkbenchHeader from './components/WorkbenchHeader.vue'
  import DynamicInfo from './components/DynamicInfo.vue'

  const loading = ref(true)

  setTimeout(() => {
    loading.value = false
  }, 1500)
</script>
