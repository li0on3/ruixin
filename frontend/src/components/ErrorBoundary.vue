<template>
  <div v-if="hasError" class="error-boundary">
    <div class="error-content">
      <div class="error-icon">
        <el-icon :size="64" color="#f56c6c"><CircleClose /></el-icon>
      </div>
      <h2>页面出现了一些问题</h2>
      <p>{{ errorMessage }}</p>
      <div class="error-actions">
        <el-button type="primary" @click="handleReload">
          <el-icon><Refresh /></el-icon>
          刷新页面
        </el-button>
        <el-button @click="handleBack">
          <el-icon><Back /></el-icon>
          返回上一页
        </el-button>
      </div>
      <details v-if="isDev" class="error-details">
        <summary>错误详情（开发模式）</summary>
        <pre>{{ errorInfo }}</pre>
      </details>
    </div>
  </div>
  <slot v-else />
</template>

<script setup>
import { ref, onErrorCaptured, computed } from 'vue'
import { useRouter } from 'vue-router'
import { CircleClose, Refresh, Back } from '@element-plus/icons-vue'

const router = useRouter()

const hasError = ref(false)
const errorMessage = ref('')
const errorInfo = ref('')

const isDev = computed(() => process.env.NODE_ENV === 'development')

// 捕获子组件的错误
onErrorCaptured((error, instance, info) => {
  console.error('Error captured:', error)
  console.error('Error info:', info)
  
  hasError.value = true
  errorMessage.value = error.message || '未知错误'
  errorInfo.value = {
    message: error.message,
    stack: error.stack,
    info: info,
    component: instance?.$options.name || 'Unknown'
  }
  
  // 阻止错误继续向上传播
  return false
})

// 刷新页面
const handleReload = () => {
  window.location.reload()
}

// 返回上一页
const handleBack = () => {
  router.back()
}

// 重置错误状态
const reset = () => {
  hasError.value = false
  errorMessage.value = ''
  errorInfo.value = ''
}

// 暴露重置方法
defineExpose({
  reset
})
</script>

<style lang="scss" scoped>
.error-boundary {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background-color: var(--bg-page, #f5f7fa);
}

.error-content {
  max-width: 500px;
  text-align: center;
  
  .error-icon {
    margin-bottom: 24px;
  }
  
  h2 {
    margin: 0 0 16px;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary, #303133);
  }
  
  p {
    margin: 0 0 32px;
    font-size: 16px;
    color: var(--text-secondary, #606266);
  }
  
  .error-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
    
    .el-button {
      min-width: 120px;
    }
  }
  
  .error-details {
    margin-top: 32px;
    text-align: left;
    
    summary {
      cursor: pointer;
      color: var(--text-tertiary, #909399);
      font-size: 14px;
      margin-bottom: 8px;
      
      &:hover {
        color: var(--text-secondary, #606266);
      }
    }
    
    pre {
      background-color: var(--bg-primary, #ffffff);
      border: 1px solid var(--border-color, #e4e7ed);
      border-radius: 4px;
      padding: 16px;
      font-size: 12px;
      line-height: 1.5;
      overflow-x: auto;
      white-space: pre-wrap;
      word-break: break-word;
      color: var(--text-primary, #303133);
    }
  }
}

// 深色主题适配
[data-theme="dark"] {
  .error-boundary {
    background-color: #0f172a;
  }
  
  .error-content {
    h2 {
      color: #f1f5f9;
    }
    
    p {
      color: #cbd5e1;
    }
    
    .error-details {
      pre {
        background-color: #1e293b;
        border-color: #475569;
        color: #f1f5f9;
      }
    }
  }
}
</style>