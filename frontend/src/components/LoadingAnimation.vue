<template>
  <transition name="fade" mode="out-in">
    <div v-if="visible" class="loading-container" :class="[`loading-${type}`, { 'loading-fullscreen': fullscreen }]">
      <div class="loading-wrapper">
        <!-- 点状加载 -->
        <div v-if="type === 'dots'" class="loading-dots">
          <span></span>
          <span></span>
          <span></span>
        </div>
        
        <!-- 条形加载 -->
        <div v-else-if="type === 'bars'" class="loading-bars">
          <span></span>
          <span></span>
          <span></span>
          <span></span>
        </div>
        
        <!-- 圆形加载 -->
        <div v-else-if="type === 'circle'" class="loading-circle"></div>
        
        <!-- 自定义Logo加载 -->
        <div v-else-if="type === 'logo'" class="loading-logo">
          <div class="logo-icon animate-pulse">
            <Coffee />
          </div>
          <div class="loading-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
        
        <!-- 骨架屏 -->
        <div v-else-if="type === 'skeleton'" class="loading-skeleton">
          <div class="skeleton-header">
            <div class="skeleton-avatar shimmer"></div>
            <div class="skeleton-lines">
              <div class="skeleton-line shimmer"></div>
              <div class="skeleton-line short shimmer"></div>
            </div>
          </div>
          <div class="skeleton-content">
            <div class="skeleton-line shimmer"></div>
            <div class="skeleton-line shimmer"></div>
            <div class="skeleton-line short shimmer"></div>
          </div>
        </div>
        
        <!-- 加载文字 -->
        <p v-if="text" class="loading-text">{{ text }}</p>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { Coffee } from '@element-plus/icons-vue'

defineProps({
  visible: {
    type: Boolean,
    default: true
  },
  type: {
    type: String,
    default: 'dots',
    validator: (value) => ['dots', 'bars', 'circle', 'logo', 'skeleton'].includes(value)
  },
  text: {
    type: String,
    default: ''
  },
  fullscreen: {
    type: Boolean,
    default: false
  }
})
</script>

<style lang="scss" scoped>
.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-8);
  
  &.loading-fullscreen {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(2px);
    z-index: var(--z-modal);
  }
  
  .loading-wrapper {
    text-align: center;
  }
  
  .loading-text {
    margin-top: var(--spacing-4);
    color: var(--text-secondary);
    font-size: 0.875rem;
  }
  
  // Logo加载样式
  .loading-logo {
    .logo-icon {
      width: 64px;
      height: 64px;
      margin: 0 auto var(--spacing-4);
      background: var(--gradient-primary);
      border-radius: var(--radius-xl);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 2rem;
    }
  }
  
  // 骨架屏样式
  .loading-skeleton {
    width: 300px;
    
    .skeleton-header {
      display: flex;
      align-items: center;
      gap: var(--spacing-3);
      margin-bottom: var(--spacing-4);
      
      .skeleton-avatar {
        width: 48px;
        height: 48px;
        border-radius: var(--radius-full);
        background: var(--gray-200);
        flex-shrink: 0;
      }
      
      .skeleton-lines {
        flex: 1;
      }
    }
    
    .skeleton-line {
      height: 12px;
      background: var(--gray-200);
      border-radius: var(--radius-md);
      margin-bottom: var(--spacing-2);
      
      &.short {
        width: 60%;
      }
    }
    
    .skeleton-content {
      .skeleton-line {
        margin-bottom: var(--spacing-3);
      }
    }
  }
}
</style>