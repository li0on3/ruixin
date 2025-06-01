<template>
  <transition name="page-loader">
    <div v-if="visible" class="page-loader-container">
      <div class="page-loader-content">
        <div class="loader-animation">
          <div class="loader-spinner">
            <div class="spinner-circle"></div>
            <div class="spinner-circle"></div>
            <div class="spinner-circle"></div>
          </div>
          <div class="loader-logo">
            <el-icon class="logo-icon"><Coffee /></el-icon>
            <h3 class="logo-text">瑞幸分销系统</h3>
          </div>
        </div>
        <p class="loader-text">{{ text }}</p>
        <div class="loader-progress">
          <div class="progress-bar" :style="{ width: progress + '%' }"></div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { Coffee } from '@element-plus/icons-vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: '页面加载中...'
  },
  duration: {
    type: Number,
    default: 1000
  }
})

const progress = ref(0)

// 监听visible变化，模拟进度条
watch(() => props.visible, (newVal) => {
  if (newVal) {
    progress.value = 0
    startProgress()
  } else {
    progress.value = 100
  }
})

const startProgress = () => {
  const interval = setInterval(() => {
    if (progress.value < 90) {
      progress.value += Math.random() * 10
    } else if (progress.value < 95) {
      progress.value += 1
    }
    
    if (!props.visible || progress.value >= 100) {
      clearInterval(interval)
      if (!props.visible) {
        progress.value = 100
      }
    }
  }, 100)
}

onMounted(() => {
  if (props.visible) {
    startProgress()
  }
})
</script>

<style lang="scss" scoped>
.page-loader-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  
  .page-loader-content {
    text-align: center;
    color: white;
    
    .loader-animation {
      margin-bottom: var(--spacing-6);
      
      .loader-spinner {
        display: flex;
        justify-content: center;
        gap: var(--spacing-2);
        margin-bottom: var(--spacing-4);
        
        .spinner-circle {
          width: 12px;
          height: 12px;
          background: rgba(255, 255, 255, 0.8);
          border-radius: 50%;
          animation: bounce 1.4s ease-in-out infinite both;
          
          &:nth-child(1) { animation-delay: -0.32s; }
          &:nth-child(2) { animation-delay: -0.16s; }
          &:nth-child(3) { animation-delay: 0s; }
        }
      }
      
      .loader-logo {
        .logo-icon {
          width: 64px;
          height: 64px;
          background: rgba(255, 255, 255, 0.15);
          border-radius: var(--radius-xl);
          display: inline-flex;
          align-items: center;
          justify-content: center;
          font-size: 2rem;
          margin-bottom: var(--spacing-3);
          animation: pulse 2s ease-in-out infinite;
        }
        
        .logo-text {
          margin: 0;
          font-size: 1.5rem;
          font-weight: 600;
          opacity: 0.9;
        }
      }
    }
    
    .loader-text {
      margin: 0 0 var(--spacing-4);
      font-size: 1rem;
      opacity: 0.8;
    }
    
    .loader-progress {
      width: 200px;
      height: 4px;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 2px;
      overflow: hidden;
      margin: 0 auto;
      
      .progress-bar {
        height: 100%;
        background: linear-gradient(90deg, rgba(255, 255, 255, 0.8), white);
        border-radius: 2px;
        transition: width 0.3s ease;
        position: relative;
        
        &::after {
          content: '';
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.6), transparent);
          animation: shimmer 1.5s infinite;
        }
      }
    }
  }
}

// 页面加载器过渡动画
.page-loader-enter-active {
  transition: all 0.3s ease-out;
}

.page-loader-leave-active {
  transition: all 0.5s ease-in;
}

.page-loader-enter-from {
  opacity: 0;
  transform: scale(1.1);
}

.page-loader-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

// 动画定义
@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

@keyframes pulse {
  0% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(255, 255, 255, 0.4);
  }
  70% {
    transform: scale(1.05);
    box-shadow: 0 0 0 10px rgba(255, 255, 255, 0);
  }
  100% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(255, 255, 255, 0);
  }
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

// 移动端适配
@media (max-width: 768px) {
  .page-loader-container {
    .page-loader-content {
      padding: var(--mobile-space-6);
      
      .loader-animation {
        .loader-logo {
          .logo-icon {
            width: 56px;
            height: 56px;
            font-size: 1.75rem;
          }
          
          .logo-text {
            font-size: 1.25rem;
          }
        }
      }
      
      .loader-progress {
        width: 160px;
      }
    }
  }
}
</style>