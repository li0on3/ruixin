<template>
  <div id="app">
    <router-view />
    
    <!-- 全局页面加载器 -->
    <PageLoader
      :visible="pageLoading"
      :text="loadingText"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import PageLoader from '@/components/PageLoader.vue'

const router = useRouter()
const userStore = useUserStore()

// 页面加载状态
const pageLoading = ref(false)
const loadingText = ref('页面加载中...')

// 路由加载状态管理
router.beforeEach((to, from, next) => {
  if (to.path !== from.path) {
    pageLoading.value = true
    loadingText.value = `正在加载${to.meta?.title || '页面'}...`
  }
  next()
})

router.afterEach(() => {
  // 延迟隐藏加载器，确保页面内容已渲染
  setTimeout(() => {
    pageLoading.value = false
  }, 300)
})

onMounted(() => {
  // 初始化用户信息
  const token = localStorage.getItem('token')
  if (token) {
    userStore.token = token
    userStore.getUserInfo()
  }
})
</script>

<style>
#app {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>