<template>
  <div class="login-enhanced-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
    
    <!-- 全屏加载 -->
    <LoadingAnimation
      v-if="initialLoading"
      type="logo"
      text="系统初始化中..."
      :visible="initialLoading"
      fullscreen
    />
    
    <!-- 登录主体 -->
    <div class="login-content" v-show="!initialLoading">
      <div class="login-card">
        <!-- 头部 -->
        <div class="login-header">
          <div class="logo">
            <div class="logo-icon">
              <el-icon><Coffee /></el-icon>
            </div>
            <div class="logo-text">
              <h1>瑞幸分销系统</h1>
              <p>Luckin Distribution Management</p>
            </div>
          </div>
        </div>
        
        <!-- 表单 -->
        <div class="login-form-container">
          <h2 class="form-title">管理员登录</h2>
          <p class="form-subtitle">请输入您的账号信息</p>
          
          <!-- 内嵌式提示卡片 -->
          <transition name="alert-fade">
            <div v-if="showAlert" class="login-alert" :class="`login-alert--${alertType}`">
              <div class="alert-icon">
                <el-icon v-if="alertType === 'error'"><CircleCloseFilled /></el-icon>
                <el-icon v-else-if="alertType === 'success'"><CircleCheckFilled /></el-icon>
                <el-icon v-else-if="alertType === 'warning'"><WarningFilled /></el-icon>
                <el-icon v-else><InfoFilled /></el-icon>
              </div>
              <div class="alert-content">
                <div class="alert-title">{{ alertTitle }}</div>
                <div class="alert-description">{{ alertDescription }}</div>
              </div>
              <div class="alert-close" @click="closeAlert">
                <el-icon><Close /></el-icon>
              </div>
            </div>
          </transition>
          
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            class="login-form"
            @keyup.enter="handleLogin"
            size="large"
          >
            <el-form-item prop="username" class="form-item" :class="{ 'is-error': hasUsernameError }">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
                :disabled="loading"
                class="form-input"
                @focus="clearFieldError('username')"
                @blur="validateField('username')"
                @input="handleUsernameInput"
                clearable
              >

              </el-input>
            </el-form-item>
            
            <el-form-item prop="password" class="form-item" :class="{ 'is-error': hasPasswordError }">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                :disabled="loading"
                show-password
                class="form-input"
                @focus="clearFieldError('password')"
                @blur="validateField('password')"
                @input="handlePasswordInput"
              >
                <template #prefix>
                  <el-icon class="input-icon"><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item class="form-item remember-me-item">
              <el-checkbox v-model="rememberMe" @keyup.enter="handleLogin">
                记住我
              </el-checkbox>
            </el-form-item>
            
            <el-form-item class="form-item login-actions">
              <el-button
                type="primary"
                :loading="loading"
                @click="handleLogin"
                class="login-button"
                :class="{ 'is-loading': loading }"
              >
                <span v-if="!loading">登 录</span>
                <span v-else>登录中...</span>
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <!-- 底部信息 -->
        <div class="login-footer">
          <div class="features">
            <div class="feature-item">
              <el-icon><Lock /></el-icon>
              <span>安全可靠</span>
            </div>
            <div class="feature-item">
              <el-icon><Monitor /></el-icon>
              <span>实时监控</span>
            </div>
            <div class="feature-item">
              <el-icon><DataLine /></el-icon>
              <span>数据分析</span>
            </div>
          </div>
          <p class="copyright">© 2025 瑞幸咖啡分销商自动化系统. All rights reserved.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import LoadingAnimation from '@/components/LoadingAnimation.vue'
import { 
  User, Lock, Coffee, Monitor, DataLine, 
  CircleCloseFilled, CircleCheckFilled, WarningFilled, InfoFilled, Close 
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const loginFormRef = ref()
const loading = ref(false)
const initialLoading = ref(true)

// 提示卡片状态
const showAlert = ref(false)
const alertType = ref('error') // error, success, warning, info
const alertTitle = ref('')
const alertDescription = ref('')
const fieldErrors = reactive({
  username: false,
  password: false
})

// 计算属性
const hasUsernameError = computed(() => fieldErrors.username)
const hasPasswordError = computed(() => fieldErrors.password)

// 表单数据
const loginForm = reactive({
  username: '',
  password: ''
})

// 记住我状态
const rememberMe = ref(false)

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

// 错误消息映射
const errorMessageMap = {
  'Invalid username or password': {
    title: '登录失败',
    description: '用户名或密码不正确，请检查后重试。如果忘记密码，请联系管理员。',
    fields: ['username', 'password']
  },
  '用户名或密码错误': {
    title: '登录失败',
    description: '用户名或密码不正确，请检查后重试。如果忘记密码，请联系管理员。',
    fields: ['username', 'password']
  },
  'Account disabled': {
    title: '账号状态异常',
    description: '您的账号已被管理员禁用，如需帮助请联系系统管理员。',
    fields: []
  },
  '账号已被禁用': {
    title: '账号状态异常',
    description: '您的账号已被管理员禁用，如需帮助请联系系统管理员。',
    fields: []
  },
  '网络连接失败': {
    title: '网络错误',
    description: '无法连接到服务器，请检查您的网络设置或稍后重试。',
    fields: []
  },
  'Network Error': {
    title: '网络错误',
    description: '无法连接到服务器，请检查您的网络设置或稍后重试。',
    fields: []
  },
  'default': {
    title: '登录失败',
    description: '系统繁忙，请稍后重试。如问题持续，请联系技术支持。',
    fields: []
  }
}

// 显示提示卡片
const showAlertCard = (type, title, description, affectedFields = []) => {
  alertType.value = type
  alertTitle.value = title
  alertDescription.value = description
  showAlert.value = true
  
  // 标记错误字段
  affectedFields.forEach(field => {
    fieldErrors[field] = true
  })
  
  // 自动隐藏（成功提示2秒，错误提示5秒）
  const duration = type === 'success' ? 2000 : 5000
  setTimeout(() => {
    closeAlert()
  }, duration)
}

// 关闭提示卡片
const closeAlert = () => {
  showAlert.value = false
  // 清除所有字段错误状态
  Object.keys(fieldErrors).forEach(key => {
    fieldErrors[key] = false
  })
}

// 清除特定字段错误状态
const clearFieldError = (field) => {
  fieldErrors[field] = false
  // 如果所有字段都没有错误，隐藏提示
  if (!Object.values(fieldErrors).some(val => val)) {
    closeAlert()
  }
}

// 验证单个字段
const validateField = async (field) => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validateField(field)
  } catch (error) {
    // 验证失败，但不需要额外处理，Element Plus 会显示错误
  }
}

// 处理用户名输入
const handleUsernameInput = () => {
  // 清除错误状态
  if (fieldErrors.username) {
    clearFieldError('username')
  }
}

// 处理密码输入
const handlePasswordInput = () => {
  // 清除错误状态
  if (fieldErrors.password) {
    clearFieldError('password')
  }
}

// 登录处理
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  // 清除之前的提示
  closeAlert()
  
  try {
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 禁用登录时的全局消息提示
    window.__skipLoginMessage = true
    
    await userStore.login(loginForm)
    
    // 处理记住我功能
    if (rememberMe.value) {
      localStorage.setItem('savedUsername', loginForm.username)
      localStorage.setItem('rememberMe', 'true')
    } else {
      localStorage.removeItem('savedUsername')
      localStorage.removeItem('rememberMe')
    }
    
    // 显示成功提示
    showAlertCard('success', '登录成功', '欢迎回来！正在为您跳转到管理后台...')
    
    // 延迟跳转，让用户看到成功提示
    setTimeout(async () => {
      await router.replace('/dashboard')
    }, 1500)
  } catch (error) {
    console.error('登录失败:', error)
    
    // 根据错误类型显示不同的提示
    let errorConfig = errorMessageMap.default
    
    if (error?.response?.data?.msg) {
      const errorMsg = error.response.data.msg.toLowerCase()
      for (const [key, config] of Object.entries(errorMessageMap)) {
        if (key !== 'default' && errorMsg.includes(key.toLowerCase())) {
          errorConfig = config
          break
        }
      }
    } else if (error?.message) {
      const errorMsg = error.message.toLowerCase()
      for (const [key, config] of Object.entries(errorMessageMap)) {
        if (key !== 'default' && errorMsg.includes(key.toLowerCase())) {
          errorConfig = config
          break
        }
      }
    } else if (!navigator.onLine) {
      errorConfig = errorMessageMap['网络连接失败']
    }
    
    showAlertCard('error', errorConfig.title, errorConfig.description, errorConfig.fields)
    
    // 错误后自动聚焦到第一个错误字段
    setTimeout(() => {
      if (errorConfig.fields.includes('username')) {
        document.querySelector('input[type="text"]')?.focus()
      } else if (errorConfig.fields.includes('password')) {
        document.querySelector('input[type="password"]')?.focus()
      }
    }, 300)
  } finally {
    loading.value = false
    window.__skipLoginMessage = false
  }
}

// 组件挂载
onMounted(() => {
  // 优化：缩短初始加载时间
  setTimeout(() => {
    initialLoading.value = false
    
    // 自动聚焦逻辑
    setTimeout(() => {
      if (loginForm.username) {
        // 如果用户名已填充，聚焦到密码框
        document.querySelector('input[type="password"]')?.focus()
      } else {
        // 否则聚焦到用户名框
        document.querySelector('input[type="text"]')?.focus()
      }
    }, 100)
  }, 800) // 从1500ms优化到800ms
  
  // 读取记住的用户名
  const savedUsername = localStorage.getItem('savedUsername')
  const savedRememberMe = localStorage.getItem('rememberMe') === 'true'
  
  if (savedRememberMe && savedUsername) {
    loginForm.username = savedUsername
    rememberMe.value = true
  }
  
  // 监听键盘事件
  const handleKeydown = (e) => {
    if (e.key === 'Escape' && showAlert.value) {
      closeAlert()
    }
  }
  window.addEventListener('keydown', handleKeydown)
  
  // 组件卸载时移除监听
  return () => {
    window.removeEventListener('keydown', handleKeydown)
  }
})
</script>

<style lang="scss" scoped>
.login-enhanced-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
  
  // 背景装饰
  .background-decoration {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    
    .decoration-circle {
      position: absolute;
      border-radius: 50%;
      background: rgba(255, 255, 255, 0.1);
      animation: float 6s ease-in-out infinite;
      
      &.circle-1 {
        width: 200px;
        height: 200px;
        top: 10%;
        left: 10%;
        animation-delay: 0s;
      }
      
      &.circle-2 {
        width: 150px;
        height: 150px;
        top: 60%;
        right: 15%;
        animation-delay: 2s;
      }
      
      &.circle-3 {
        width: 100px;
        height: 100px;
        bottom: 20%;
        left: 20%;
        animation-delay: 4s;
      }
    }
  }
  
  .login-content {
    width: 100%;
    max-width: 440px;
    padding: var(--spacing-6);
    position: relative;
    z-index: 1;
    
    .login-card {
      background: rgba(255, 255, 255, 0.95);
      backdrop-filter: blur(20px);
      border-radius: var(--radius-2xl);
      padding: var(--spacing-8);
      box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      border: 1px solid rgba(255, 255, 255, 0.2);
      animation: slideInUp 0.6s ease-out;
      
      .login-header {
        text-align: center;
        margin-bottom: var(--spacing-8);
        
        .logo {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: var(--spacing-4);
          margin-bottom: var(--spacing-6);
          
          .logo-icon {
            width: 64px;
            height: 64px;
            background: var(--gradient-primary);
            border-radius: var(--radius-xl);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 2rem;
            box-shadow: var(--shadow-lg);
          }
          
          .logo-text {
            text-align: left;
            
            h1 {
              margin: 0;
              font-size: 1.5rem;
              font-weight: 700;
              color: var(--text-primary);
              background: var(--gradient-primary);
              -webkit-background-clip: text;
              -webkit-text-fill-color: transparent;
              background-clip: text;
            }
            
            p {
              margin: var(--spacing-1) 0 0;
              font-size: 0.875rem;
              color: var(--text-tertiary);
              letter-spacing: 0.05em;
            }
          }
        }
      }
      
      .login-form-container {
        .form-title {
          margin: 0 0 var(--spacing-2);
          font-size: 1.5rem;
          font-weight: 600;
          color: var(--text-primary);
          text-align: center;
        }
        
        .form-subtitle {
          margin: 0 0 var(--spacing-6);
          color: var(--text-secondary);
          text-align: center;
          font-size: 0.875rem;
        }
        
        // 登录提示卡片
        .login-alert {
          display: flex;
          align-items: flex-start;
          gap: var(--spacing-3);
          padding: var(--spacing-4);
          margin-bottom: var(--spacing-6);
          border-radius: var(--radius-lg);
          border: 1px solid;
          position: relative;
          animation: slideInDown 0.3s ease-out;
          
          .alert-icon {
            flex-shrink: 0;
            font-size: 1.25rem;
            line-height: 1;
            margin-top: 2px;
          }
          
          .alert-content {
            flex: 1;
            
            .alert-title {
              font-size: 1rem;
              font-weight: 600;
              margin-bottom: var(--spacing-1);
              line-height: 1.4;
            }
            
            .alert-description {
              font-size: 0.875rem;
              line-height: 1.5;
              opacity: 0.9;
            }
          }
          
          .alert-close {
            flex-shrink: 0;
            width: 20px;
            height: 20px;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            border-radius: var(--radius-sm);
            transition: all 0.2s ease;
            opacity: 0.7;
            
            &:hover {
              opacity: 1;
              background-color: rgba(0, 0, 0, 0.05);
            }
            
            .el-icon {
              font-size: 0.875rem;
            }
          }
          
          // 错误样式
          &--error {
            background-color: #fef2f2;
            border-color: #fecaca;
            color: #dc2626;
            
            .alert-icon {
              color: #dc2626;
            }
          }
          
          // 成功样式
          &--success {
            background-color: #f0fdf4;
            border-color: #bbf7d0;
            color: #16a34a;
            
            .alert-icon {
              color: #16a34a;
            }
          }
          
          // 警告样式
          &--warning {
            background-color: #fffbeb;
            border-color: #fde68a;
            color: #d97706;
            
            .alert-icon {
              color: #d97706;
            }
          }
          
          // 信息样式
          &--info {
            background-color: #eff6ff;
            border-color: #bfdbfe;
            color: #2563eb;
            
            .alert-icon {
              color: #2563eb;
            }
          }
        }
        
        .login-form {
          .form-item {
            margin-bottom: var(--spacing-5);
            
            // 错误状态样式
            &.is-error {
              .form-input {
                .el-input__wrapper {
                  border-color: #dc2626;
                  background-color: #fef2f2;
                  
                  &:hover {
                    border-color: #dc2626;
                  }
                  
                  &.is-focus {
                    border-color: #dc2626;
                    box-shadow: 0 0 0 4px rgba(220, 38, 38, 0.1);
                  }
                }
              }
            }
            
            .form-input {
              .el-input__wrapper {
                padding: var(--spacing-4);
                border-radius: var(--radius-lg);
                border: 2px solid var(--border-default);
                background: var(--bg-primary);
                transition: all 0.3s ease;
                
                &:hover {
                  border-color: var(--primary-400);
                }
                
                &.is-focus {
                  border-color: var(--primary-500);
                  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
                }
              }
              
              .input-icon {
                color: var(--text-tertiary);
                font-size: 1.125rem;
              }
              
              .el-input__inner {
                font-size: 1rem;
                color: var(--text-primary);
                
                &::placeholder {
                  color: var(--text-tertiary);
                }
              }
            }
            
            &.remember-me-item {
              margin-bottom: var(--spacing-4);
              
              .el-checkbox {
                font-size: 0.875rem;
                color: var(--text-secondary);
                
                &:hover {
                  color: var(--primary-500);
                }
                
                .el-checkbox__label {
                  font-weight: normal;
                }
              }
            }
            
            &.login-actions {
              margin-bottom: 0;
              margin-top: var(--spacing-5);
              
              .login-button {
                width: 100%;
                height: 48px;
                font-size: 1rem;
                font-weight: 600;
                border-radius: var(--radius-lg);
                background: var(--gradient-primary);
                border: none;
                color: white;
                transition: all 0.3s ease;
                position: relative;
                overflow: hidden;
                
                &::before {
                  content: '';
                  position: absolute;
                  top: 0;
                  left: -100%;
                  width: 100%;
                  height: 100%;
                  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
                  transition: left 0.5s ease;
                }
                
                &:hover {
                  transform: translateY(-2px);
                  box-shadow: var(--shadow-lg);
                  
                  &::before {
                    left: 100%;
                  }
                }
                
                &:active {
                  transform: translateY(0);
                }
                
                &.is-loading {
                  background: var(--primary-400);
                  cursor: not-allowed;
                  
                  &:hover {
                    transform: none;
                    box-shadow: none;
                  }
                }
              }
            }
          }
        }
      }
      
      .login-footer {
        margin-top: var(--spacing-8);
        
        .features {
          display: flex;
          justify-content: space-around;
          margin-bottom: var(--spacing-6);
          
          .feature-item {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: var(--spacing-2);
            color: var(--text-secondary);
            font-size: 0.875rem;
            
            .el-icon {
              width: 32px;
              height: 32px;
              background: var(--bg-secondary);
              border-radius: var(--radius-lg);
              display: flex;
              align-items: center;
              justify-content: center;
              color: var(--primary-500);
            }
          }
        }
        
        .copyright {
          margin: 0;
          text-align: center;
          color: var(--text-tertiary);
          font-size: 0.75rem;
          line-height: 1.5;
        }
      }
    }
  }
}

// 动画定义
@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 提示卡片过渡动画
.alert-fade-enter-active,
.alert-fade-leave-active {
  transition: all 0.3s ease;
}

.alert-fade-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.alert-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

// 移动端适配
@media (max-width: 768px) {
  .login-enhanced-container {
    padding: var(--mobile-space-4);
    
    .login-content {
      max-width: 100%;
      padding: 0;
      
      .login-card {
        padding: var(--mobile-space-6);
        border-radius: var(--mobile-border-radius);
        
        .login-header {
          margin-bottom: var(--mobile-space-6);
          
          .logo {
            flex-direction: column;
            gap: var(--mobile-space-3);
            
            .logo-icon {
              width: 56px;
              height: 56px;
              font-size: 1.75rem;
            }
            
            .logo-text {
              text-align: center;
              
              h1 {
                font-size: 1.25rem;
              }
            }
          }
        }
        
        .login-form-container {
          .form-title {
            font-size: 1.25rem;
          }
          
          .login-alert {
            margin-bottom: var(--mobile-space-4);
          }
          
          .login-form {
            .form-item {
              margin-bottom: var(--mobile-space-4);
              
              .form-input {
                .el-input__wrapper {
                  padding: var(--mobile-space-3);
                }
              }
              
              &.remember-me-item {
                margin-bottom: var(--mobile-space-3);
                
                .el-checkbox {
                  font-size: 0.813rem;
                }
              }
              
              &.login-actions {
                margin-top: var(--mobile-space-4);
                
                .login-button {
                  height: var(--mobile-touch-target);
                }
              }
            }
          }
        }
        
        .login-footer {
          margin-top: var(--mobile-space-6);
          
          .features {
            .feature-item {
              .el-icon {
                width: 28px;
                height: 28px;
              }
              
              span {
                font-size: 0.75rem;
              }
            }
          }
        }
      }
    }
    
    .background-decoration {
      .decoration-circle {
        &.circle-1 {
          width: 120px;
          height: 120px;
        }
        
        &.circle-2 {
          width: 80px;
          height: 80px;
        }
        
        &.circle-3 {
          width: 60px;
          height: 60px;
        }
      }
    }
  }
}
</style>