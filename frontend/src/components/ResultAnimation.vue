<template>
  <transition name="fade-scale" mode="out-in">
    <div v-if="visible" class="result-animation">
      <!-- 成功动画 -->
      <svg v-if="type === 'success'" class="success-checkmark" viewBox="0 0 52 52">
        <circle class="checkmark-circle" cx="26" cy="26" r="25" fill="none"/>
        <path class="checkmark-check" fill="none" d="M14.1 27.2l7.1 7.2 16.7-16.8"/>
      </svg>
      
      <!-- 错误动画 -->
      <svg v-else-if="type === 'error'" class="error-cross" viewBox="0 0 52 52">
        <circle class="cross-circle" cx="26" cy="26" r="25" fill="none"/>
        <path class="cross-x" fill="none" d="M16,16 L36,36 M36,16 L16,36"/>
      </svg>
      
      <!-- 警告动画 -->
      <div v-else-if="type === 'warning'" class="warning-icon animate-pulse">
        <WarningFilled />
      </div>
      
      <!-- 信息动画 -->
      <div v-else-if="type === 'info'" class="info-icon animate-breathe">
        <InfoFilled />
      </div>
      
      <!-- 标题和描述 -->
      <h3 v-if="title" class="result-title">{{ title }}</h3>
      <p v-if="description" class="result-description">{{ description }}</p>
      
      <!-- 操作按钮 -->
      <div v-if="$slots.actions" class="result-actions">
        <slot name="actions"></slot>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { WarningFilled, InfoFilled } from '@element-plus/icons-vue'

defineProps({
  visible: {
    type: Boolean,
    default: true
  },
  type: {
    type: String,
    default: 'success',
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  }
})
</script>

<style lang="scss" scoped>
.result-animation {
  text-align: center;
  padding: var(--spacing-8);
  
  .result-title {
    margin: var(--spacing-4) 0 var(--spacing-2);
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  .result-description {
    margin: 0 0 var(--spacing-6);
    color: var(--text-secondary);
    font-size: 0.875rem;
  }
  
  .result-actions {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--spacing-3);
  }
  
  // 警告图标
  .warning-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto;
    background: var(--warning-light);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3rem;
    color: var(--warning);
  }
  
  // 信息图标
  .info-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto;
    background: var(--info-light);
    border-radius: var(--radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3rem;
    color: var(--info);
  }
}
</style>