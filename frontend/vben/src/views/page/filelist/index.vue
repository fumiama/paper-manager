<template>
  <PageWrapper :class="prefixCls" :title="t('routes.filelist.name')">
    <template #headerContent>
      <BasicUpload
        v-if="hasPermission([RoleEnum.SUPER, RoleEnum.FILE_MANAGER])"
        :maxSize="20"
        :maxNumber="10"
        @change="handleChange"
        :api="uploadApi"
        :accept="['application/vnd.openxmlformats-officedocument.wordprocessingml.document']"
      />
    </template>
    <div :class="`${prefixCls}__top`">
      <a-row :gutter="12">
        <a-col :span="8" :class="`${prefixCls}__top-col`">
          <div>总文件数</div>
          <p> {{ pagination.total }}</p>
        </a-col>
        <a-col :span="8" :class="`${prefixCls}__top-col`">
          <div>占用空间</div>
          <p> {{ cardList._totalSize }}MB </p>
        </a-col>
        <a-col :span="8" :class="`${prefixCls}__top-col`">
          <div>总题目数</div>
          <p> {{ cardList._totalQuestions }} </p>
        </a-col>
      </a-row>
    </div>

    <div :class="`${prefixCls}__content`">
      <a-list :pagination="pagination">
        <template
          v-for="item in getListOfPage(pagination.pageSize, pagination.current)"
          :key="item.id"
        >
          <a-list-item class="list">
            <a-list-item-meta>
              <template #avatar>
                <Icon class="icon" v-if="item.icon" :icon="item.icon" :color="item.color" />
              </template>
              <template #title>
                <span>{{ item.title }}</span>
                <div class="extra">
                  <a-button color="success" :disabled="item.percent < 100"> 查阅 </a-button>
                  &nbsp;&nbsp;
                  <a-button
                    color="warning"
                    v-if="hasPermission([RoleEnum.SUPER, RoleEnum.FILE_MANAGER])"
                    :disabled="item.percent != 0"
                  >
                    解析
                  </a-button>
                  &nbsp;&nbsp;
                  <a-button
                    color="error"
                    v-if="hasPermission([RoleEnum.SUPER])"
                    :disabled="item.percent > 0 && item.percent < 100"
                    @click="deleteFileBy(item.id)"
                  >
                    删除
                  </a-button>
                </div>
              </template>
              <template #description>
                <div class="description">
                  {{ item.description }}
                </div>
                <div class="info">
                  <div><span>文件大小</span>{{ item.size }}MB</div>
                  <div><span>上传用户</span>{{ item.author }}</div>
                  <div><span>上传时间</span>{{ item.datetime }}</div>
                </div>
                <div class="progress">
                  <div><span>解析进度</span></div>
                  <Progress
                    :percent="item.percent"
                    :status='((): "normal" | "success" | "active" | "exception" | undefined => { 
                      if (item.percent < 100) return "active"
                      return "success"
                    })()'
                  />
                </div>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </div>
  </PageWrapper>
</template>
<script lang="ts">
  import { Progress, Row, Col } from 'ant-design-vue'
  import { defineComponent } from 'vue'
  import { Icon } from '/@/components/Icon'
  import { BasicUpload } from '/@/components/Upload'
  import { cardList, getListOfPage, deleteFileByID, pagination } from './data'
  import { PageWrapper } from '/@/components/Page'
  import { useMessage } from '/@/hooks/web/useMessage'
  import { usePermission } from '/@/hooks/web/usePermission'
  import { RoleEnum } from '/@/enums/roleEnum'
  import { List } from 'ant-design-vue'
  import { uploadApi } from '/@/api/sys/upload'
  import { useI18n } from '/@/hooks/web/useI18n'
  import { delFile } from '/@/api/page'
  import { DelFile } from '/@/api/page/model/fileListModel'

  const { t } = useI18n()
  const { createMessage } = useMessage()

  function deleteFileBy(id: number) {
    delFile(id).then((value: DelFile) => {
      createMessage.info(value.msg)
      deleteFileByID(id)
    })
  }

  export default defineComponent({
    components: {
      BasicUpload,
      Icon,
      Progress,
      PageWrapper,
      [List.name]: List,
      [List.Item.name]: List.Item,
      AListItemMeta: List.Item.Meta,
      [Row.name]: Row,
      [Col.name]: Col,
    },
    setup() {
      const { hasPermission } = usePermission()

      return {
        t,
        RoleEnum,
        handleChange: (list: string[]) => {
          createMessage.info(`已上传文件${JSON.stringify(list)}`)
        },
        uploadApi,
        hasPermission,
        prefixCls: 'list-basic',
        getListOfPage,
        deleteFileBy,
        cardList,
        pagination,
      }
    },
  })
</script>
<style lang="less" scoped>
  .list-basic {
    &__top {
      padding: 24px;
      text-align: center;
      background-color: @component-background;
      &-col {
        &:not(:last-child) {
          border-right: 1px dashed @border-color-base;
        }
        div {
          margin-bottom: 12px;
          font-size: 14px;
          line-height: 22px;
          color: @text-color;
        }
        p {
          margin: 0;
          font-size: 24px;
          line-height: 32px;
          color: @text-color;
        }
      }
    }
    &__content {
      padding: 24px;
      margin-top: 12px;
      background-color: @component-background;
      .list {
        position: relative;
      }
      .icon {
        font-size: 40px !important;
      }
      .extra {
        position: absolute;
        top: 38px;
        right: 8px;
      }
      .description {
        display: inline-block;
        width: 20%;
      }
      .info {
        display: inline-block;
        width: 40%;
        text-align: center;
        vertical-align: top;
        div {
          display: inline-block;
          padding: 0 20px;
          span {
            display: block;
          }
        }
      }
      .progress {
        display: inline-block;
        width: 15%;
        vertical-align: top;
      }
    }
  }
</style>
