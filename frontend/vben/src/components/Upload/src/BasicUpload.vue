<template>
  <div>
    <Space>
      <a-button type="primary" @click="openUploadModal" preIcon="carbon:cloud-upload">
        {{ t('component.upload.upload') }}
      </a-button>
      <Tooltip placement="bottom" v-if="showPreview">
        <template #title>
          {{ t('component.upload.uploaded') }}
          <template v-if="fileList.length">
            {{ fileList.length }}
          </template>
        </template>
      </Tooltip>
    </Space>
    <UploadModal
      v-bind="bindValue"
      :previewFileList="fileList"
      @register="registerUploadModal"
      @change="handleChange"
      @delete="handleDelete"
    />
  </div>
</template>
<script lang="ts">
  import { defineComponent, ref, watch, unref, computed } from 'vue'
  import { Tooltip, Space } from 'ant-design-vue'
  import { useModal } from '/@/components/Modal'
  import { uploadContainerProps } from './props'
  import { omit } from 'lodash-es'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { isArray } from '/@/utils/is'
  import UploadModal from './UploadModal.vue'
  export default defineComponent({
    name: 'BasicUpload',
    components: { UploadModal, Space, Tooltip },
    props: uploadContainerProps,
    emits: ['change', 'delete', 'preview-delete', 'update:value'],
    setup(props, { emit, attrs }) {
      const { t } = useI18n()
      // 上传modal
      const [registerUploadModal, { openModal: openUploadModal }] = useModal()
      //   预览modal
      const [registerPreviewModal, { openModal: openPreviewModal }] = useModal()
      const fileList = ref<string[]>([])
      const showPreview = computed(() => {
        const { emptyHidePreview } = props
        if (!emptyHidePreview) return true
        return emptyHidePreview ? fileList.value.length > 0 : true
      })
      const bindValue = computed(() => {
        const value = { ...attrs, ...props }
        return omit(value, 'onChange')
      })
      watch(
        () => props.value,
        (value = []) => {
          fileList.value = isArray(value) ? value : []
        },
        { immediate: true },
      )
      // 上传modal保存操作
      function handleChange(urls: string[]) {
        fileList.value = [...unref(fileList), ...(urls || [])]
        emit('update:value', fileList.value)
        emit('change', fileList.value)
      }
      function handleDelete(record: Recordable) {
        emit('delete', record)
      }
      function handlePreviewDelete(url: string) {
        emit('preview-delete', url)
      }
      return {
        registerUploadModal,
        openUploadModal,
        handleChange,
        registerPreviewModal,
        openPreviewModal,
        fileList,
        showPreview,
        bindValue,
        handleDelete,
        handlePreviewDelete,
        t,
      }
    },
  })
</script>
