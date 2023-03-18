<template>
  <div class="lg:flex">
    <Avatar :src="userinfo.avatar || headerImg" :size="72" class="!mx-auto !block" />
    <div class="md:ml-6 flex flex-col justify-center md:mt-0 mt-2">
      <h1 class="md:text-lg text-md">
        {{
          ((): string => {
            let hour: number = new Date().getHours()
            if (hour < 10) return '早'
            else if (hour < 19) return '午'
            else return '晚'
          })()
        }}安, {{ userinfo.realName }}, 要注意劳逸结合哦!</h1
      >
      <span class="text-secondary"> {{ userinfo.desc }} </span>
    </div>
    <div class="flex flex-1 justify-end md:mt-0 mt-4">
      <div class="flex flex-col justify-center text-right md:mr-10 mr-4">
        <span class="text-secondary"> 课程组人数 </span>
        <span class="text-2xl"> {{ userscount }} </span>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
  import { ref, computed } from 'vue'
  import { Avatar } from 'ant-design-vue'
  import { useUserStore } from '/@/store/modules/user'
  import headerImg from '/@/assets/images/header.jpg'
  import { getUsersCount } from '/@/api/sys/user'

  const userStore = useUserStore()
  const userinfo = computed(() => userStore.getUserInfo)
  const userscount = ref(0)
  getUsersCount().then((value: number) => {
    userscount.value = value
  })
</script>
