import { defineStore } from 'pinia'
import { login, logout, getUserInfo } from '@/api/auth'
import { showSuccess, showError } from '@/utils/message'
import router from '@/router/optimized'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    userInfo: null,
    roles: [],
    // 添加状态标记
    isLoggingOut: false,
    isGettingUserInfo: false,
    lastGetUserInfoTime: 0
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isSuperAdmin: (state) => state.roles.includes('super_admin'),
    isAdmin: (state) => state.roles.includes('admin') || state.roles.includes('super_admin')
  },
  
  actions: {
    async login(loginForm) {
      try {
        const { data } = await login(loginForm)
        this.token = data.token
        this.userInfo = data.user
        this.roles = [data.user.role]
        
        localStorage.setItem('token', data.token)
        
        // 移除这里的成功提示，保留 LoginEnhanced.vue 中的动画效果
        // showSuccess('登录成功，欢迎回来！')
        router.push('/')
      } catch (error) {
        // 错误消息已在 request.js 中处理
        throw error
      }
    },
    
    async getUserInfo() {
      // 防止重复请求
      const now = Date.now()
      if (this.isGettingUserInfo || (now - this.lastGetUserInfoTime < 1000)) {
        return
      }
      
      // 如果没有token，直接返回
      if (!this.token) {
        return
      }
      
      this.isGettingUserInfo = true
      this.lastGetUserInfoTime = now
      
      try {
        const { data } = await getUserInfo()
        this.userInfo = data
        this.roles = [data.role]
      } catch (error) {
        // 静默处理错误，不自动logout
        console.error('获取用户信息失败:', error)
        // 只有在401错误时才清理状态，其他错误保持现状
        if (error?.response?.status === 401) {
          // 清理本地状态但不跳转
          this.token = ''
          this.userInfo = null
          this.roles = []
          localStorage.removeItem('token')
        }
      } finally {
        this.isGettingUserInfo = false
      }
    },
    
    async logout() {
      // 防止重复执行logout
      if (this.isLoggingOut) {
        return
      }
      
      this.isLoggingOut = true
      
      try {
        // 只有在有token时才调用logout接口
        if (this.token) {
          await logout()
        }
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.token = ''
        this.userInfo = null
        this.roles = []
        localStorage.removeItem('token')
        
        // 只有在不在登录页时才跳转
        if (router.currentRoute.value.path !== '/login') {
          router.replace('/login')
        }
        
        // 延迟重置标记，防止快速重复调用
        setTimeout(() => {
          this.isLoggingOut = false
        }, 1000)
      }
    }
  }
})