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
          
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            class="login-form"
            @keyup.enter="handleLogin"
            size="large"
          >
            <el-form-item prop="username" class="form-item">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
                :disabled="loading"
                class="form-input"
              >
                <template #prefix>
                  <el-icon class="input-icon"><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item prop="password" class="form-item">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                :disabled="loading"
                show-password
                class="form-input"
              >
                <template #prefix>
                  <el-icon class="input-icon"><Lock /></el-icon>
                </template>
              </el-input>
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
    
    <!-- 结果反馈 -->
    <ResultAnimation
      v-if="showResult"
      :visible="showResult"
      :type="resultType"
      :title="resultTitle"
      :description="resultDescription"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import { User, Lock, Coffee, Monitor, DataLine } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const loginFormRef = ref()
const loading = ref(false)
const initialLoading = ref(true)
const showResult = ref(false)
const resultType = ref('success')
const resultTitle = ref('')
const resultDescription = ref('')

// 表单数据
const loginForm = reactive({
  username: '',
  password: ''
})

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

// 登录处理
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 模拟网络延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    await userStore.login(loginForm)
    
    // 登录成功会在 store 中处理，这里显示成功结果
    showSuccessResult()
    
    // 使用 replace 确保跳转，并等待一小段时间让动画显示
    setTimeout(async () => {
      await router.replace('/dashboard')
    }, 1000)
  } catch (error) {
    console.error('登录失败:', error)
    
    if (error?.message?.includes('用户名或密码')) {
      showErrorResult('用户名或密码错误，请重新输入')
    } else {
      showErrorResult('登录失败，请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}

// 显示成功结果
const showSuccessResult = () => {
  resultType.value = 'success'
  resultTitle.value = '登录成功'
  resultDescription.value = '正在跳转到管理后台...'
  showResult.value = true
  
  setTimeout(() => {
    showResult.value = false
  }, 2000)
}

// 显示错误结果
const showErrorResult = (message) => {
  resultType.value = 'error'
  resultTitle.value = '登录失败'
  resultDescription.value = message
  showResult.value = true
  
  setTimeout(() => {
    showResult.value = false
  }, 3000)
}

// 组件挂载
onMounted(() => {
  // 模拟系统初始化
  setTimeout(() => {
    initialLoading.value = false
  }, 1500)
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
        
        .login-form {
          .form-item {
            margin-bottom: var(--spacing-5);
            
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
            
            &.login-actions {
              margin-bottom: 0;
              margin-top: var(--spacing-6);
              
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
          
          .login-form {
            .form-item {
              margin-bottom: var(--mobile-space-4);
              
              .form-input {
                .el-input__wrapper {
                  padding: var(--mobile-space-3);
                }
              }
              
              &.login-actions {
                margin-top: var(--mobile-space-5);
                
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