<template>
  <PageWrapper :class="prefixCls" title="试卷资源管理器">
    <div :class="`${prefixCls}__top`">
      <a-row :gutter="12">
        <a-col :span="8" :class="`${prefixCls}__top-col`">
          <div>总文件数</div>
          <p>8</p>
        </a-col>
        <a-col :span="8" :class="`${prefixCls}__top-col`">
          <div>占用空间</div>
          <p>32MB</p>
        </a-col>
        <a-col :span="8" :class="`${prefixCls}__top-col`">
          <div>总题目数</div>
          <p>24</p>
        </a-col>
      </a-row>
    </div>

    <div :class="`${prefixCls}__content`">
      <a-list :pagination="pagination">
        <template v-for="item in list" :key="item.id">
          <a-list-item class="list">
            <a-list-item-meta>
              <template #avatar>
                <Icon class="icon" v-if="item.icon" :icon="item.icon" :color="item.color" />
              </template>
              <template #title>
                <span>{{ item.title }}</span>
                <div class="extra" v-if="item.extra">
                  <a-button ghost color="success"> 成功 </a-button>
                  &nbsp;&nbsp;
                  <a-button ghost color="warning"> 警告 </a-button>
                  &nbsp;&nbsp;
                  <a-button ghost color="error"> 错误 </a-button>
                </div>
              </template>
              <template #description>
                <div class="description">
                  {{ item.description }}
                </div>
                <div class="info">
                  <div><span>Owner</span>{{ item.author }}</div>
                  <div><span>开始时间</span>{{ item.datetime }}</div>
                </div>
                <div class="progress">
                  <Progress :percent="item.percent" status="active" />
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
  import Icon from '/@/components/Icon/index'
  import { cardList } from './data'
  import { PageWrapper } from '/@/components/Page'
  import { List } from 'ant-design-vue'
  export default defineComponent({
    components: {
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
      return {
        prefixCls: 'list-basic',
        list: cardList,
        pagination: {
          show: true,
          pageSize: 3,
        },
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
