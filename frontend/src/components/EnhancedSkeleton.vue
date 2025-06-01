<template>
  <div class="enhanced-skeleton" :class="{ 'animate-pulse': loading }">
    <!-- 表格骨架屏 -->
    <div v-if="type === 'table'" class="skeleton-table">
      <div class="skeleton-table-header">
        <div 
          v-for="col in columns" 
          :key="`header-${col}`"
          class="skeleton-table-cell header-cell"
          :style="{ width: getColumnWidth(col) }"
        >
          <div class="skeleton-line" :style="{ width: '80%' }"></div>
        </div>
      </div>
      <div class="skeleton-table-body">
        <div 
          v-for="row in rows" 
          :key="`row-${row}`"
          class="skeleton-table-row"
          :style="{ animationDelay: `${row * 50}ms` }"
        >
          <div 
            v-for="col in columns"
            :key="`cell-${row}-${col}`"
            class="skeleton-table-cell"
            :style="{ width: getColumnWidth(col) }"
          >
            <div 
              class="skeleton-line"
              :style="{ width: getRandomWidth(60, 90) }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 卡片骨架屏 -->
    <div v-else-if="type === 'card'" class="skeleton-card">
      <div v-if="showAvatar" class="skeleton-avatar"></div>
      <div class="skeleton-content">
        <div class="skeleton-line skeleton-title" :style="{ width: '70%' }"></div>
        <div class="skeleton-line-group">
          <div 
            v-for="line in contentLines" 
            :key="`line-${line}`"
            class="skeleton-line"
            :style="{ 
              width: line === contentLines ? '60%' : '100%',
              animationDelay: `${line * 100}ms`
            }"
          ></div>
        </div>
        <div v-if="showActions" class="skeleton-actions">
          <div class="skeleton-button"></div>
          <div class="skeleton-button"></div>
        </div>
      </div>
    </div>

    <!-- 列表骨架屏 -->
    <div v-else-if="type === 'list'" class="skeleton-list">
      <div 
        v-for="item in rows" 
        :key="`list-${item}`"
        class="skeleton-list-item"
        :style="{ animationDelay: `${item * 80}ms` }"
      >
        <div class="skeleton-avatar skeleton-avatar-sm"></div>
        <div class="skeleton-list-content">
          <div class="skeleton-line" :style="{ width: '60%' }"></div>
          <div class="skeleton-line" :style="{ width: '80%' }"></div>
        </div>
      </div>
    </div>

    <!-- 统计卡片骨架屏 -->
    <div v-else-if="type === 'stats'" class="skeleton-stats">
      <div 
        v-for="stat in 4" 
        :key="`stat-${stat}`"
        class="skeleton-stat-card"
        :style="{ animationDelay: `${stat * 150}ms` }"
      >
        <div class="skeleton-stat-icon"></div>
        <div class="skeleton-stat-content">
          <div class="skeleton-line skeleton-stat-number"></div>
          <div class="skeleton-line skeleton-stat-label"></div>
        </div>
      </div>
    </div>

    <!-- 图表骨架屏 -->
    <div v-else-if="type === 'chart'" class="skeleton-chart">
      <div class="skeleton-chart-header">
        <div class="skeleton-line" :style="{ width: '40%' }"></div>
        <div class="skeleton-chart-legend">
          <div class="skeleton-legend-item" v-for="i in 3" :key="`legend-${i}`"></div>
        </div>
      </div>
      <div class="skeleton-chart-body">
        <div class="skeleton-chart-y-axis">
          <div v-for="i in 5" :key="`y-${i}`" class="skeleton-y-label"></div>
        </div>
        <div class="skeleton-chart-content">
          <div 
            v-for="bar in 8" 
            :key="`bar-${bar}`"
            class="skeleton-chart-bar"
            :style="{ 
              height: getRandomHeight(30, 80),
              animationDelay: `${bar * 100}ms`
            }"
          ></div>
        </div>
      </div>
    </div>

    <!-- 默认骨架屏 -->
    <div v-else class="skeleton-default">
      <div class="skeleton-line skeleton-title"></div>
      <div class="skeleton-line-group">
        <div v-for="line in 3" :key="`default-${line}`" class="skeleton-line"></div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EnhancedSkeleton',
  props: {
    // 是否显示加载状态
    loading: {
      type: Boolean,
      default: true
    },
    // 骨架屏类型
    type: {
      type: String,
      default: 'default',
      validator: value => ['default', 'table', 'card', 'list', 'stats', 'chart'].includes(value)
    },
    // 表格行数
    rows: {
      type: Number,
      default: 5
    },
    // 表格列数
    columns: {
      type: Number,
      default: 4
    },
    // 内容行数
    contentLines: {
      type: Number,
      default: 3
    },
    // 是否显示头像
    showAvatar: {
      type: Boolean,
      default: false
    },
    // 是否显示操作按钮
    showActions: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    // 获取随机宽度
    getRandomWidth(min, max) {
      return `${Math.floor(Math.random() * (max - min + 1)) + min}%`
    },
    
    // 获取随机高度
    getRandomHeight(min, max) {
      return `${Math.floor(Math.random() * (max - min + 1)) + min}%`
    },
    
    // 获取列宽度
    getColumnWidth(col) {
      const widths = ['20%', '25%', '30%', '25%']
      return widths[col - 1] || '25%'
    }
  }
}
</script>

<style lang="scss" scoped>
.enhanced-skeleton {
  animation: skeletonFadeIn var(--transition-base) ease-out;
}

// 基础骨架线
.skeleton-line {
  height: 16px;
  background: linear-gradient(
    90deg,
    var(--gray-200) 25%,
    var(--gray-100) 50%,
    var(--gray-200) 75%
  );
  background-size: 200% 100%;
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-2);
  
  &:last-child {
    margin-bottom: 0;
  }
}

.skeleton-title {
  height: 20px;
  margin-bottom: var(--spacing-4);
}

.skeleton-line-group {
  .skeleton-line {
    margin-bottom: var(--spacing-3);
  }
}

// 表格骨架屏
.skeleton-table {
  border-radius: var(--radius-xl);
  overflow: hidden;
  border: 1px solid var(--border-light);
}

.skeleton-table-header {
  display: flex;
  background: var(--bg-secondary);
  padding: var(--spacing-4);
  gap: var(--spacing-4);
}

.skeleton-table-body {
  .skeleton-table-row {
    display: flex;
    padding: var(--spacing-4);
    gap: var(--spacing-4);
    border-bottom: 1px solid var(--border-light);
    animation: skeletonSlideIn var(--transition-base) ease-out both;
    
    &:last-child {
      border-bottom: none;
    }
  }
}

.skeleton-table-cell {
  &.header-cell .skeleton-line {
    height: 14px;
    background: var(--gray-300);
  }
}

// 卡片骨架屏
.skeleton-card {
  display: flex;
  gap: var(--spacing-4);
  padding: var(--spacing-6);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-xl);
  background: var(--bg-primary);
}

.skeleton-avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  background: linear-gradient(
    90deg,
    var(--gray-200) 25%,
    var(--gray-100) 50%,
    var(--gray-200) 75%
  );
  background-size: 200% 100%;
  flex-shrink: 0;
  
  &.skeleton-avatar-sm {
    width: 32px;
    height: 32px;
  }
}

.skeleton-content {
  flex: 1;
}

.skeleton-actions {
  display: flex;
  gap: var(--spacing-3);
  margin-top: var(--spacing-4);
}

.skeleton-button {
  width: 80px;
  height: 32px;
  background: linear-gradient(
    90deg,
    var(--gray-200) 25%,
    var(--gray-100) 50%,
    var(--gray-200) 75%
  );
  background-size: 200% 100%;
  border-radius: var(--radius-lg);
}

// 列表骨架屏
.skeleton-list {
  .skeleton-list-item {
    display: flex;
    align-items: center;
    gap: var(--spacing-3);
    padding: var(--spacing-4);
    border-bottom: 1px solid var(--border-light);
    animation: skeletonSlideIn var(--transition-base) ease-out both;
    
    &:last-child {
      border-bottom: none;
    }
  }
  
  .skeleton-list-content {
    flex: 1;
    
    .skeleton-line {
      margin-bottom: var(--spacing-2);
      
      &:last-child {
        margin-bottom: 0;
      }
    }
  }
}

// 统计卡片骨架屏
.skeleton-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-6);
}

.skeleton-stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-4);
  padding: var(--spacing-6);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-xl);
  background: var(--bg-primary);
  animation: skeletonSlideIn var(--transition-base) ease-out both;
}

.skeleton-stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  background: linear-gradient(
    90deg,
    var(--primary-200) 25%,
    var(--primary-100) 50%,
    var(--primary-200) 75%
  );
  background-size: 200% 100%;
}

.skeleton-stat-content {
  flex: 1;
}

.skeleton-stat-number {
  height: 24px;
  width: 60%;
  margin-bottom: var(--spacing-2);
}

.skeleton-stat-label {
  height: 14px;
  width: 80%;
}

// 图表骨架屏
.skeleton-chart {
  border: 1px solid var(--border-light);
  border-radius: var(--radius-xl);
  padding: var(--spacing-6);
  background: var(--bg-primary);
}

.skeleton-chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-6);
}

.skeleton-chart-legend {
  display: flex;
  gap: var(--spacing-4);
}

.skeleton-legend-item {
  width: 60px;
  height: 12px;
  background: linear-gradient(
    90deg,
    var(--gray-200) 25%,
    var(--gray-100) 50%,
    var(--gray-200) 75%
  );
  background-size: 200% 100%;
  border-radius: var(--radius-sm);
}

.skeleton-chart-body {
  display: flex;
  gap: var(--spacing-4);
  height: 200px;
}

.skeleton-chart-y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 40px;
}

.skeleton-y-label {
  width: 100%;
  height: 12px;
  background: linear-gradient(
    90deg,
    var(--gray-200) 25%,
    var(--gray-100) 50%,
    var(--gray-200) 75%
  );
  background-size: 200% 100%;
  border-radius: var(--radius-sm);
}

.skeleton-chart-content {
  flex: 1;
  display: flex;
  align-items: end;
  gap: var(--spacing-2);
}

.skeleton-chart-bar {
  flex: 1;
  background: linear-gradient(
    180deg,
    var(--primary-200) 0%,
    var(--primary-100) 100%
  );
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  animation: skeletonBarGrow var(--transition-slow) ease-out both;
}

// 默认骨架屏
.skeleton-default {
  padding: var(--spacing-6);
}

// 动画
.animate-pulse {
  .skeleton-line,
  .skeleton-avatar,
  .skeleton-button,
  .skeleton-stat-icon,
  .skeleton-legend-item,
  .skeleton-y-label {
    animation: skeletonShimmer 1.5s ease-in-out infinite;
  }
}

@keyframes skeletonShimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

@keyframes skeletonFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes skeletonSlideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes skeletonBarGrow {
  from {
    transform: scaleY(0);
  }
  to {
    transform: scaleY(1);
  }
}

// 暗色主题适配
[data-theme="dark"] {
  .skeleton-line,
  .skeleton-avatar,
  .skeleton-button,
  .skeleton-legend-item,
  .skeleton-y-label {
    background: linear-gradient(
      90deg,
      var(--gray-700) 25%,
      var(--gray-600) 50%,
      var(--gray-700) 75%
    );
    background-size: 200% 100%;
  }
  
  .skeleton-stat-icon {
    background: linear-gradient(
      90deg,
      var(--primary-800) 25%,
      var(--primary-700) 50%,
      var(--primary-800) 75%
    );
    background-size: 200% 100%;
  }
  
  .skeleton-chart-bar {
    background: linear-gradient(
      180deg,
      var(--primary-600) 0%,
      var(--primary-700) 100%
    );
  }
}

// 响应式适配
@include mobile-only {
  .skeleton-stats {
    grid-template-columns: 1fr;
  }
  
  .skeleton-card {
    flex-direction: column;
    text-align: center;
  }
  
  .skeleton-chart-body {
    height: 150px;
  }
  
  .skeleton-table {
    .skeleton-table-header,
    .skeleton-table-row {
      flex-direction: column;
      gap: var(--spacing-2);
    }
    
    .skeleton-table-cell {
      width: 100% !important;
    }
  }
}

// 性能优化
.enhanced-skeleton {
  contain: layout style paint;
  will-change: contents;
}

// 无障碍支持
.enhanced-skeleton[aria-label]::before {
  content: attr(aria-label);
  position: absolute;
  left: -10000px;
  width: 1px;
  height: 1px;
  overflow: hidden;
}
</style>