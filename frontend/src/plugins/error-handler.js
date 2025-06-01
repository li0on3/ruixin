// 全局错误处理插件
import { ElMessage } from 'element-plus'
import router from '@/router/optimized'

// 错误类型映射
const errorTypeMap = {
  'Network Error': '网络连接失败，请检查您的网络设置',
  'Request failed with status code 500': '服务器错误，请稍后重试',
  'Request failed with status code 404': '请求的资源不存在',
  'Request failed with status code 403': '您没有权限访问此资源',
  'Request failed with status code 401': '登录已过期，请重新登录',
  'timeout': '请求超时，请检查网络连接'
}

// 获取友好的错误消息
function getFriendlyMessage(error) {
  // 检查是否是已知的错误类型
  for (const [key, message] of Object.entries(errorTypeMap)) {
    if (error.message?.includes(key)) {
      return message
    }
  }
  
  // 检查是否有自定义消息
  if (error.response?.data?.msg) {
    return error.response.data.msg
  }
  
  // 默认错误消息
  return error.message || '系统出现异常，请稍后重试'
}

// 错误处理函数
function handleError(error, vm, info) {
  console.error('Global error:', error)
  console.error('Error info:', info)
  
  // 获取友好的错误消息
  const message = getFriendlyMessage(error)
  
  // 显示错误提示
  ElMessage({
    message,
    type: 'error',
    duration: 5000,
    showClose: true
  })
  
  // 如果是严重错误，可以记录到日志服务
  if (process.env.NODE_ENV === 'production') {
    // 这里可以集成错误上报服务，如 Sentry
    // logErrorToService(error, info)
  }
}

// 处理未捕获的 Promise 错误
function handleUnhandledRejection(event) {
  console.error('Unhandled promise rejection:', event.reason)
  
  // 阻止默认行为（防止控制台报错）
  event.preventDefault()
  
  // 显示错误提示
  const message = getFriendlyMessage(event.reason)
  ElMessage({
    message,
    type: 'error',
    duration: 5000,
    showClose: true
  })
}

export default {
  install(app) {
    // Vue 错误处理
    app.config.errorHandler = handleError
    
    // 全局未捕获错误处理
    window.addEventListener('error', (event) => {
      console.error('Global error event:', event)
      
      // 资源加载错误
      if (event.target !== window) {
        ElMessage({
          message: '资源加载失败，请刷新页面重试',
          type: 'warning',
          duration: 3000
        })
        return
      }
      
      // JavaScript 错误
      handleError(event.error, null, 'Global error')
    })
    
    // Promise 错误处理
    window.addEventListener('unhandledrejection', handleUnhandledRejection)
    
    // 路由错误处理
    router.onError((error) => {
      console.error('Router error:', error)
      
      // 组件加载失败
      if (error.message?.includes('Failed to fetch dynamically imported module')) {
        ElMessage({
          message: '页面加载失败，正在刷新...',
          type: 'warning',
          duration: 2000,
          onClose: () => {
            window.location.reload()
          }
        })
      } else {
        handleError(error, null, 'Router error')
      }
    })
    
    // 提供手动错误处理方法
    app.config.globalProperties.$handleError = handleError
    
    // 开发环境下的调试工具
    if (process.env.NODE_ENV === 'development') {
      window.__ERROR_HANDLER__ = {
        testError: () => {
          throw new Error('This is a test error')
        },
        testPromiseError: () => {
          Promise.reject(new Error('This is a test promise rejection'))
        },
        testNetworkError: () => {
          throw new Error('Network Error')
        }
      }
      
      console.log('Error handler debug tools available at window.__ERROR_HANDLER__')
    }
  }
}