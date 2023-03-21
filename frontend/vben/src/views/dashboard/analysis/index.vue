<template>
  <div class="p-4">
    <SiteAnalysis :loading="loading" />
  </div>
</template>
<script lang="ts" setup>
  import { ref } from 'vue'
  import SiteAnalysis from './components/SiteAnalysis.vue'

  const loading = ref(true)

  setTimeout(() => {
    loading.value = false
  }, 1500)
</script>
