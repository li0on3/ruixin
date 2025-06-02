import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import router from '@/router/optimized'
import requestManager from '@/utils/requestManager'

// 401处理标记
let is401Processing = false

// 创建axios实例 - 直接访问后端，不使用代理
const service = axios.create({
  baseURL: window.location.protocol + '//' + window.location.hostname + ':8080/api/v1',
  timeout: 30000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    
    // 如果是getUserInfo请求，添加特殊处理
    if (config.url.includes('/admin/user/info')) {
      // 防止重复的getUserInfo请求
      config = requestManager.add(config)
    }
    
    return config
  },
  error => {
    console.error('Request error:', error)
    requestManager.remove(error.config)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    // 移除已完成的请求
    requestManager.remove(response.config)
    
    const res = response.data
    
    // 根据后端约定的响应格式处理
    if (res.code !== 200) {
      // 检查是否在登录页面且设置了跳过提示的标记
      const isLoginPage = router.currentRoute.value.path === '/login'
      const skipMessage = window.__skipLoginMessage
      
      if (!isLoginPage || !skipMessage) {
        ElMessage({
          message: res.msg || '请求失败',
          type: 'error',
          duration: 5 * 1000
        })
      }
      
      // 401: 未授权
      if (res.code === 401) {
        // 防止重复处理401错误
        if (!is401Processing) {
          is401Processing = true
          const userStore = useUserStore()
          userStore.logout()
          
          // 重置标记
          setTimeout(() => {
            is401Processing = false
          }, 2000)
        }
      }
      
      return Promise.reject(new Error(res.msg || '请求失败'))
    } else {
      return res
    }
  },
  error => {
    // 移除失败的请求
    if (error.config) {
      requestManager.remove(error.config)
    }
    
    // 如果是取消的请求，不处理
    if (axios.isCancel(error)) {
      console.log('请求被取消:', error.message)
      return Promise.reject(error)
    }
    
    console.error('Response error:', error)
    
    let message = '网络错误'
    if (error.response) {
      switch (error.response.status) {
        case 401:
          message = '未授权，请重新登录'
          // 防止重复处理401错误
          if (!is401Processing) {
            is401Processing = true
            const userStore = useUserStore()
            userStore.logout()
            
            // 重置标记
            setTimeout(() => {
              is401Processing = false
            }, 2000)
          }
          break
        case 403:
          message = '拒绝访问'
          break
        case 400:
          message = error.response.data?.msg || '请求参数错误'
          break
        case 404:
          message = '请求地址错误'
          break
        case 500:
          message = error.response.data?.msg || '服务器内部错误'
          break
        default:
          message = error.response.data?.msg || '请求失败'
      }
    } else if (error.request) {
      message = '网络连接失败'
    }
    
    // 检查是否在登录页面且设置了跳过提示的标记
    const isLoginPage = router.currentRoute.value.path === '/login'
    const skipMessage = window.__skipLoginMessage
    
    if (!isLoginPage || !skipMessage) {
      ElMessage({
        message,
        type: 'error',
        duration: 4000,
        showClose: true,
        grouping: true,
        customClass: 'custom-message-error'
      })
    }
    
    return Promise.reject(error)
  }
)

export default service