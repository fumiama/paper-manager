<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script lang="ts">
  import { defineComponent, ref, computed, unref } from 'vue'
  import { BasicModal, useModalInner } from '/@/components/Modal'
  import { BasicForm, useForm } from '/@/components/Form/index'
  import { isNameExist } from '/@/api/sys/user'

  export default defineComponent({
    name: 'AccountModal',
    components: { BasicModal, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
      const isUpdate = ref(true)
      const rowId = ref('')

      const [registerForm, { setFieldsValue, updateSchema, resetFields, validate }] = useForm({
        labelWidth: 100,
        baseColProps: { span: 24 },
        schemas: [
          {
            field: 'name',
            label: '用户名',
            component: 'Input',
            ifShow: () => {
              return !unref(isUpdate)
            },
            rules: [
              {
                required: true,
                message: '请输入用户名',
              },
              {
                validator(_, value) {
                  return new Promise((resolve, reject) => {
                    if (unref(isUpdate)) {
                      resolve()
                      return
                    }
                    isNameExist(value)
                      .then((v) => {
                        if (!v) resolve()
                        else reject('用户名已存在')
                      })
                      .catch((err) => {
                        reject(err.message || '验证失败')
                      })
                  })
                },
              },
            ],
          },
          {
            field: 'pwd',
            label: '密码',
            component: 'InputPassword',
            required: true,
            ifShow: false,
          },
          {
            label: '角色',
            field: 'role',
            component: 'ApiSelect',
            componentProps: {
              api: () => {
                return [
                  {
                    roleName: '课程组长',
                    value: 'super',
                  },
                  {
                    roleName: '归档代理',
                    value: 'filemgr',
                  },
                  {
                    roleName: '课程组员',
                    value: 'user',
                  },
                ]
              },
              labelField: 'roleName',
              valueField: 'value',
            },
            required: true,
          },
          {
            field: 'nick',
            label: '昵称',
            component: 'Input',
            required: true,
          },
          {
            label: '简介',
            field: 'desc',
            component: 'InputTextArea',
          },
        ],
        showActionButtonGroup: false,
        actionColOptions: {
          span: 23,
        },
      })

      const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
        resetFields()
        setModalProps({ confirmLoading: false })
        isUpdate.value = !!data?.isUpdate

        if (unref(isUpdate)) {
          rowId.value = data.record.id
          setFieldsValue({
            ...data.record,
          })
        }

        updateSchema([
          {
            field: 'pwd',
            show: !unref(isUpdate),
          },
        ])
      })

      const getTitle = computed(() => (!unref(isUpdate) ? '新增账号' : '编辑账号'))

      async function handleSubmit() {
        try {
          const values = await validate()
          setModalProps({ confirmLoading: true })
          // TODO custom api
          console.log(values)
          closeModal()
          emit('success', { isUpdate: unref(isUpdate), values: { ...values, id: rowId.value } })
        } finally {
          setModalProps({ confirmLoading: false })
        }
      }

      return { registerModal, registerForm, getTitle, handleSubmit }
    },
  })
</script>
