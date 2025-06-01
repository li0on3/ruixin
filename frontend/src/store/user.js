import { defineStore } from 'pinia'
import { login, logout, getUserInfo } from '@/api/auth'
import { showSuccess, showError } from '@/utils/message'
import router from '@/router/optimized'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    userInfo: null,
    roles: []
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
      try {
        const { data } = await getUserInfo()
        this.userInfo = data
        this.roles = [data.role]
      } catch (error) {
        this.logout()
      }
    },
    
    async logout() {
      try {
        await logout()
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.token = ''
        this.userInfo = null
        this.roles = []
        localStorage.removeItem('token')
        // 使用 replace 避免返回按钮问题
        router.replace('/login')
      }
    }
  }
})