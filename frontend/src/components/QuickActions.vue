<template>
  <div class="quick-actions">
    <!-- 快捷操作面板 -->
    <el-drawer
      v-model="drawerVisible"
      title="快捷操作"
      :size="isMobile ? '100%' : '400px'"
      direction="rtl"
    >
      <div class="quick-actions-content">
        <!-- 常用功能 -->
        <div class="actions-section">
          <h3 class="section-title">
            <el-icon><Star /></el-icon>
            常用功能
          </h3>
          <div class="actions-grid">
            <div
              v-for="action in commonActions"
              :key="action.key"
              class="action-item"
              @click="handleAction(action)"
            >
              <div class="action-icon" :style="{ background: action.color }">
                <el-icon>
                  <component :is="action.icon" />
                </el-icon>
              </div>
              <div class="action-content">
                <h4 class="action-title">{{ action.title }}</h4>
                <p class="action-desc">{{ action.description }}</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 最近访问 -->
        <div class="actions-section">
          <h3 class="section-title">
            <el-icon><Clock /></el-icon>
            最近访问
          </h3>
          <div class="recent-list">
            <div
              v-for="item in recentPages"
              :key="item.path"
              class="recent-item"
              @click="navigateTo(item.path)"
            >
              <div class="recent-icon">
                <el-icon>
                  <component :is="item.icon" />
                </el-icon>
              </div>
              <div class="recent-content">
                <h4 class="recent-title">{{ item.title }}</h4>
                <p class="recent-time">{{ formatTime(item.visitTime) }}</p>
              </div>
              <el-icon class="recent-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
        
        <!-- 快捷导航 -->
        <div class="actions-section">
          <h3 class="section-title">
            <el-icon><Menu /></el-icon>
            快捷导航
          </h3>
          <div class="nav-grid">
            <el-button
              v-for="nav in quickNavigation"
              :key="nav.path"
              class="nav-button"
              @click="navigateTo(nav.path)"
            >
              <el-icon>
                <component :is="nav.icon" />
              </el-icon>
              <span>{{ nav.title }}</span>
            </el-button>
          </div>
        </div>
        
        <!-- 系统信息 -->
        <div class="actions-section">
          <h3 class="section-title">
            <el-icon><InfoFilled /></el-icon>
            系统信息
          </h3>
          <div class="system-info">
            <div class="info-item">
              <span class="info-label">当前版本</span>
              <span class="info-value">v1.0.0</span>
            </div>
            <div class="info-item">
              <span class="info-label">在线时长</span>
              <span class="info-value">{{ onlineTime }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">最后同步</span>
              <span class="info-value">{{ lastSyncTime }}</span>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>
    
    <!-- 触发按钮 -->
    <el-button
      class="quick-actions-trigger"
      type="primary"
      circle
      @click="drawerVisible = true"
      :class="{ 'is-mobile': isMobile }"
    >
      <el-icon><Lightning /></el-icon>
    </el-button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, inject } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Star, Clock, Menu, InfoFilled, Lightning, ArrowRight,
  Plus, Upload, Search, TrendCharts, CreditCard, UserFilled,
  ShoppingCart, Setting, Files, DataLine
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

// 检测是否为移动端
const isMobile = inject('isMobile', computed(() => window.innerWidth <= 768))

// 响应式数据
const drawerVisible = ref(false)
const onlineStartTime = ref(Date.now())
const lastSyncTime = ref(new Date().toLocaleString())

// 常用功能配置
const commonActions = ref([
  {
    key: 'add-card',
    title: '添加卡片',
    description: '快速添加新的卡片',
    icon: 'Plus',
    color: 'var(--success)',
    action: () => router.push('/cards')
  },
  {
    key: 'batch-import',
    title: '批量导入',
    description: '批量导入卡片数据',
    icon: 'Upload',
    color: 'var(--primary)',
    action: () => router.push('/cards/batch-import')
  },
  {
    key: 'search-data',
    title: '数据查询',
    description: '查询各类业务数据',
    icon: 'Search',
    color: 'var(--info)',
    action: () => router.push('/statistics')
  },
  {
    key: 'view-statistics',
    title: '统计报表',
    description: '查看详细统计数据',
    icon: 'TrendCharts',
    color: 'var(--warning)',
    action: () => router.push('/statistics')
  }
])

// 快捷导航配置
const quickNavigation = ref([
  { title: '仪表盘', icon: 'DataLine', path: '/dashboard' },
  { title: '卡片管理', icon: 'CreditCard', path: '/cards' },
  { title: '分销商', icon: 'UserFilled', path: '/distributors' },
  { title: '订单', icon: 'ShoppingCart', path: '/orders' },
  { title: '统计', icon: 'TrendCharts', path: '/statistics' },
  { title: '设置', icon: 'Setting', path: '/settings' }
])

// 最近访问页面（从localStorage获取）
const recentPages = computed(() => {
  const recent = JSON.parse(localStorage.getItem('recentPages') || '[]')
  return recent.slice(0, 5) // 只显示最近5个
})

// 在线时长计算
const onlineTime = computed(() => {
  const diff = Date.now() - onlineStartTime.value
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  return `${hours}时${minutes}分`
})

// 方法
const handleAction = (action) => {
  if (action.action) {
    action.action()
  }
  drawerVisible.value = false
}

const navigateTo = (path) => {
  router.push(path)
  drawerVisible.value = false
}

const formatTime = (timestamp) => {
  const now = Date.now()
  const diff = now - timestamp
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  return '刚刚'
}

// 记录页面访问
const recordPageVisit = (route) => {
  if (!route.meta?.title) return
  
  const recentPages = JSON.parse(localStorage.getItem('recentPages') || '[]')
  const pageInfo = {
    path: route.path,
    title: route.meta.title,
    icon: getRouteIcon(route.path),
    visitTime: Date.now()
  }
  
  // 移除重复项
  const filtered = recentPages.filter(p => p.path !== route.path)
  
  // 添加到最前面
  filtered.unshift(pageInfo)
  
  // 只保留最近10个
  const limited = filtered.slice(0, 10)
  
  localStorage.setItem('recentPages', JSON.stringify(limited))
}

// 获取路由图标
const getRouteIcon = (path) => {
  const iconMap = {
    '/dashboard': 'DataLine',
    '/cards': 'CreditCard',
    '/distributors': 'UserFilled',
    '/orders': 'ShoppingCart',
    '/statistics': 'TrendCharts',
    '/finance/transactions': 'CreditCard',
    '/finance/withdrawals': 'Wallet',
    '/products': 'ShoppingBag',
    '/settings': 'Setting'
  }
  return iconMap[path] || 'Document'
}

// 定时器更新同步时间
let syncTimer = null

onMounted(() => {
  // 记录当前页面访问
  recordPageVisit(route)
  
  // 监听路由变化
  router.afterEach((to) => {
    recordPageVisit(to)
  })
  
  // 定时更新同步时间
  syncTimer = setInterval(() => {
    lastSyncTime.value = new Date().toLocaleString()
  }, 60000) // 每分钟更新一次
})

onUnmounted(() => {
  if (syncTimer) {
    clearInterval(syncTimer)
  }
})
</script>

<style lang="scss" scoped>
.quick-actions {
  position: fixed;
  bottom: var(--spacing-8);
  right: var(--spacing-8);
  z-index: var(--z-fixed);
  
  .quick-actions-trigger {
    width: 56px;
    height: 56px;
    box-shadow: var(--shadow-lg);
    transition: all 0.3s ease;
    
    &:hover {
      transform: scale(1.1);
      box-shadow: var(--shadow-xl);
    }
    
    &.is-mobile {
      width: 48px;
      height: 48px;
      bottom: var(--mobile-space-6);
      right: var(--mobile-space-6);
    }
  }
  
  .quick-actions-content {
    padding: var(--spacing-4);
    
    .actions-section {
      margin-bottom: var(--spacing-8);
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .section-title {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        margin: 0 0 var(--spacing-4);
        font-size: 1rem;
        font-weight: 600;
        color: var(--text-primary);
        
        .el-icon {
          color: var(--primary);
        }
      }
    }
    
    // 常用功能网格
    .actions-grid {
      display: grid;
      gap: var(--spacing-3);
      
      .action-item {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        padding: var(--spacing-4);
        background: var(--bg-secondary);
        border-radius: var(--radius-lg);
        cursor: pointer;
        transition: all 0.2s ease;
        
        &:hover {
          background: var(--bg-tertiary);
          transform: translateY(-2px);
        }
        
        .action-icon {
          width: 40px;
          height: 40px;
          border-radius: var(--radius-lg);
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-size: 1.25rem;
          flex-shrink: 0;
        }
        
        .action-content {
          flex: 1;
          min-width: 0;
          
          .action-title {
            margin: 0 0 var(--spacing-1);
            font-size: 0.875rem;
            font-weight: 600;
            color: var(--text-primary);
          }
          
          .action-desc {
            margin: 0;
            font-size: 0.75rem;
            color: var(--text-secondary);
            line-height: 1.4;
          }
        }
      }
    }
    
    // 最近访问列表
    .recent-list {
      .recent-item {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        padding: var(--spacing-3);
        border-radius: var(--radius-md);
        cursor: pointer;
        transition: background 0.2s ease;
        
        &:hover {
          background: var(--bg-secondary);
        }
        
        .recent-icon {
          width: 32px;
          height: 32px;
          border-radius: var(--radius-md);
          background: var(--bg-tertiary);
          display: flex;
          align-items: center;
          justify-content: center;
          color: var(--text-secondary);
          flex-shrink: 0;
        }
        
        .recent-content {
          flex: 1;
          min-width: 0;
          
          .recent-title {
            margin: 0 0 var(--spacing-1);
            font-size: 0.875rem;
            font-weight: 500;
            color: var(--text-primary);
          }
          
          .recent-time {
            margin: 0;
            font-size: 0.75rem;
            color: var(--text-tertiary);
          }
        }
        
        .recent-arrow {
          color: var(--text-tertiary);
          font-size: 0.875rem;
        }
      }
    }
    
    // 快捷导航网格
    .nav-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
      gap: var(--spacing-2);
      
      .nav-button {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: var(--spacing-2);
        padding: var(--spacing-4);
        border: 1px solid var(--border-default);
        background: var(--bg-primary);
        
        &:hover {
          border-color: var(--primary);
          color: var(--primary);
        }
        
        span {
          font-size: 0.75rem;
        }
      }
    }
    
    // 系统信息
    .system-info {
      .info-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: var(--spacing-3) 0;
        border-bottom: 1px solid var(--border-light);
        
        &:last-child {
          border-bottom: none;
        }
        
        .info-label {
          font-size: 0.875rem;
          color: var(--text-secondary);
        }
        
        .info-value {
          font-size: 0.875rem;
          font-weight: 500;
          color: var(--text-primary);
        }
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .quick-actions {
    bottom: var(--mobile-space-6);
    right: var(--mobile-space-6);
    
    .quick-actions-content {
      padding: var(--mobile-space-4);
      
      .nav-grid {
        grid-template-columns: repeat(2, 1fr);
      }
    }
  }
}
</style>