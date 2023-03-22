<template>
  <BasicModal v-bind="$attrs" @register="registerModal" title="编辑账号" @ok="handleSubmit">
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script lang="ts">
  import { defineComponent, ref } from 'vue'
  import { BasicModal, useModalInner } from '/@/components/Modal'
  import { BasicForm, useForm } from '/@/components/Form/index'
  import { setOthersInfoApi, setRole } from '/@/api/sys/user'
  import { useUserStore } from '/@/store/modules/user'

  export default defineComponent({
    name: 'AccountModal',
    components: { BasicModal, BasicForm },
    emits: ['success', 'register'],
    setup(_, { emit }) {
      const rowId = ref('')
      const roles = [
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

      const [registerForm, { setFieldsValue, resetFields, validate }] = useForm({
        labelWidth: 100,
        baseColProps: { span: 24 },
        schemas: [
          {
            label: '角色',
            field: 'role',
            component: 'ApiSelect',
            componentProps: {
              api: () => {
                return roles
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

        rowId.value = data.record.id
        setFieldsValue({
          ...data.record,
        })
      })

      const nick2id = { 课程组长: 1, super: 1, 归档代理: 2, filemgr: 2, 课程组员: 3, user: 3 }

      async function handleSubmit() {
        try {
          const values = await validate()
          setModalProps({ confirmLoading: true })
          closeModal()
          await setOthersInfoApi({ id: Number(rowId.value), nick: values.nick, desc: values.desc })
          if (useUserStore().getUserInfo.userId != Number(rowId.value)) {
            const rid = nick2id[values.role]
            await setRole(Number(rowId.value), rid)
            values.role = roles[rid - 1].roleName
          }
          emit('success', { values: { ...values, id: rowId.value } })
        } finally {
          setModalProps({ confirmLoading: false })
        }
      }

      return { registerModal, registerForm, handleSubmit }
    },
  })
</script>
