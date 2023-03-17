<template>
  <CollapseContainer title="基本设置" :canExpan="false">
    <a-row :gutter="24">
      <a-col :span="14">
        <BasicForm @register="register" />
      </a-col>
      <a-col :span="10">
        <div class="change-avatar">
          <div class="mb-2">头像</div>
          <CropperAvatar
            :uploadApi="onUpload as any"
            :value="avatarRef"
            btnText="更换头像"
            :btnProps="{ preIcon: 'ant-design:cloud-upload-outlined' }"
            @change="updateAvatar"
            width="150"
          />
        </div>
      </a-col>
    </a-row>
    <Button type="primary" @click="handleSubmit"> 更新基本信息 </Button>
  </CollapseContainer>
</template>
<script lang="ts">
  import { Button, Row, Col } from 'ant-design-vue'
  import { ref, defineComponent, onMounted } from 'vue'
  import { BasicForm, useForm } from '/@/components/Form/index'
  import { CollapseContainer } from '/@/components/Container'
  import { CropperAvatar } from '/@/components/Cropper'

  import { useMessage } from '/@/hooks/web/useMessage'

  import headerImg from '/@/assets/images/header.jpg'
  import { baseSetschemas } from './data'
  import { useUserStore } from '/@/store/modules/user'
  import { uploadApi } from '/@/api/sys/upload'

  export default defineComponent({
    components: {
      BasicForm,
      CollapseContainer,
      Button,
      ARow: Row,
      ACol: Col,
      CropperAvatar,
    },
    setup() {
      const { createMessage } = useMessage()
      const userStore = useUserStore()
      const { avatar } = userStore.getUserInfo

      const [register, { setFieldsValue }] = useForm({
        labelWidth: 120,
        schemas: baseSetschemas,
        showActionButtonGroup: false,
      })

      onMounted(async () => {
        const data = userStore.getUserInfo
        setFieldsValue(data)
      })

      const avatarRef = ref(avatar || headerImg)

      function updateAvatar({ src }) {
        const userinfo = userStore.getUserInfo
        userinfo.avatar = src
        userStore.setUserInfo(userinfo)
      }

      async function onUpload(value: { file: Blob; name: string }) {
        const data = userStore.getUserInfo
        const result = await uploadApi(
          {
            name: 'avatar',
            file: value.file,
            filename: data.username,
          },
          () => {},
        )
        avatarRef.value = result.data.url
        return result
      }

      return {
        avatarRef,
        register,
        onUpload,
        updateAvatar,
        handleSubmit: () => {
          createMessage.success('更新成功！')
        },
      }
    },
  })
</script>

<style lang="less" scoped>
  .change-avatar {
    img {
      display: block;
      margin-bottom: 15px;
      border-radius: 50%;
    }
  }
</style>
