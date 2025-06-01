import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import router from '@/router/optimized'

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
    return config
  },
  error => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 根据后端约定的响应格式处理
    if (res.code !== 200) {
      ElMessage({
        message: res.msg || '请求失败',
        type: 'error',
        duration: 5 * 1000
      })
      
      // 401: 未授权
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        // 确保跳转到登录页
        if (router.currentRoute.value.path !== '/login') {
          router.push('/login')
        }
      }
      
      return Promise.reject(new Error(res.msg || '请求失败'))
    } else {
      return res
    }
  },
  error => {
    console.error('Response error:', error)
    
    let message = '网络错误'
    if (error.response) {
      switch (error.response.status) {
        case 401:
          message = '未授权，请重新登录'
          const userStore = useUserStore()
          userStore.logout()
          // 确保跳转到登录页
          if (router.currentRoute.value.path !== '/login') {
            router.push('/login')
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
    
    ElMessage({
      message,
      type: 'error',
      duration: 4000,
      showClose: true,
      grouping: true,
      customClass: 'custom-message-error'
    })
    
    return Promise.reject(error)
  }
)

export default service