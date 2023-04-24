<template>
  <div class="step3">
    <a-result status="success" title="试卷生成成功" :sub-title="msg">
      <template #extra>
        <a-button type="primary" @click="redo"> 再次生成 </a-button>
        <a-button @click="downloadDocx"> 下载试卷 </a-button>
      </template>
    </a-result>
    <div class="docxWrap" :style="{ width }">
      <div ref="docxRef"></div>
    </div>
  </div>
</template>
<script lang="ts">
  import { defineComponent, ref } from 'vue'
  import { Result, Descriptions } from 'ant-design-vue'
  import { renderAsync } from 'docx-preview'
  import { state } from './data'
  import { dlGeneratedFile } from '/@/api/page'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { downloadByData } from '/@/utils/file/download'

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
    components: {
      [Result.name]: Result,
      [Descriptions.name]: Descriptions,
      [Descriptions.Item.name]: Descriptions.Item,
    },
    props: {
      width: {
        type: String as PropType<string>,
        default: '100%',
      },
    },
    emits: ['redo'],
    setup(_, { emit }) {
      const { createMessage } = useMessage()
      ;(async () => {
        try {
          const data = await dlGeneratedFile()
          if (data) {
            loadDocx(data)
            return
          }
        } catch (error) {
          createMessage.error('加载docx错误: ' + (error as unknown as Error).message)
        }
      })()
      return {
        msg: state.step2Values,
        redo: () => {
          emit('redo')
        },
        docxRef,
        downloadDocx,
        docxNameRef,
        docxSizeRef,
      }
    },
  })
</script>
<style lang="less" scoped>
  .step3 {
    width: 600px;
    margin: 0 auto;
  }

  .docxWrap {
    padding-top: 0px;
    margin: 0 auto;
    overflow-x: auto;
    display: grid;
    align-items: center;
    justify-items: center;
  }
</style>
