<template>
  <div class="enhanced-breadcrumb">
    <div class="breadcrumb-nav">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/dashboard' }" class="breadcrumb-home">
          <el-icon><House /></el-icon>
          <span v-if="!isMobile">首页</span>
        </el-breadcrumb-item>
        
        <el-breadcrumb-item
          v-for="(item, index) in breadcrumbItems"
          :key="index"
          :to="item.path ? { path: item.path } : undefined"
          :class="{ 'is-current': index === breadcrumbItems.length - 1 }"
        >
          <el-icon v-if="item.icon">
            <component :is="item.icon" />
          </el-icon>
          <span>{{ item.title }}</span>
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    
    <!-- 快捷操作 -->
    <div class="breadcrumb-actions" v-if="quickActions.length > 0">
      <el-dropdown trigger="click" placement="bottom-end">
        <el-button text class="quick-actions-btn">
          <el-icon><MoreFilled /></el-icon>
          <span v-if="!isMobile">快捷操作</span>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="action in quickActions"
              :key="action.key"
              @click="handleQuickAction(action)"
              :divided="action.divided"
            >
              <el-icon v-if="action.icon">
                <component :is="action.icon" />
              </el-icon>
              {{ action.title }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    
    <!-- 页面操作按钮 -->
    <div class="page-actions" v-if="pageActions.length > 0">
      <el-button-group>
        <el-button
          v-for="action in pageActions"
          :key="action.key"
          :type="action.type || 'default'"
          :icon="action.icon"
          :loading="action.loading"
          @click="handlePageAction(action)"
          :size="isMobile ? 'small' : 'default'"
        >
          <span v-if="!isMobile || !action.icon">{{ action.title }}</span>
        </el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { House, MoreFilled } from '@element-plus/icons-vue'

const props = defineProps({
  quickActions: {
    type: Array,
    default: () => []
  },
  pageActions: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['quick-action', 'page-action'])

const route = useRoute()
const router = useRouter()

// 检测是否为移动端
const isMobile = inject('isMobile', computed(() => window.innerWidth <= 768))

// 路由配置映射
const routeConfig = {
  '/dashboard': { title: '仪表盘', icon: 'DataLine' },
  '/cards': { title: '卡片管理', icon: 'CreditCard' },
  '/cards/bindings': { title: '卡片绑定管理', icon: 'Link', parent: '/cards' },
  '/cards/batch-import': { title: '批量导入卡片', icon: 'Upload', parent: '/cards' },
  '/cards/batches': { title: '批次管理', icon: 'Files', parent: '/cards' },
  '/distributors': { title: '分销商管理', icon: 'UserFilled' },
  '/orders': { title: '订单管理', icon: 'ShoppingCart' },
  '/statistics': { title: '数据统计', icon: 'TrendCharts' },
  '/finance': { title: '财务管理', icon: 'Wallet' },
  '/finance/transactions': { title: '交易记录', icon: 'CreditCard', parent: '/finance' },
  '/finance/withdrawals': { title: '提现管理', icon: 'Wallet', parent: '/finance' },
  '/luckin/prices': { title: '价格管理', icon: 'Coffee' },
  '/products': { title: '商品管理', icon: 'ShoppingBag' },
  '/products/available': { title: '可用商品', icon: 'ShoppingBag', parent: '/products' },
  '/admins': { title: '管理员管理', icon: 'User' },
  '/settings': { title: '系统设置', icon: 'Setting' },
  '/stores': { title: '店铺查询', icon: 'Location' }
}

// 生成面包屑项目
const breadcrumbItems = computed(() => {
  const currentPath = route.path
  const items = []
  
  // 获取当前路由配置
  const currentConfig = routeConfig[currentPath]
  if (!currentConfig) {
    // 如果没有配置，使用路由的meta信息
    if (route.meta?.title) {
      items.push({
        title: route.meta.title,
        path: null // 当前页面不可点击
      })
    }
    return items
  }
  
  // 如果有父级路径，先添加父级
  if (currentConfig.parent) {
    const parentConfig = routeConfig[currentConfig.parent]
    if (parentConfig) {
      items.push({
        title: parentConfig.title,
        icon: parentConfig.icon,
        path: currentConfig.parent
      })
    }
  }
  
  // 添加当前页面
  items.push({
    title: currentConfig.title,
    icon: currentConfig.icon,
    path: null // 当前页面不可点击
  })
  
  return items
})

// 处理快捷操作
const handleQuickAction = (action) => {
  emit('quick-action', action)
}

// 处理页面操作
const handlePageAction = (action) => {
  emit('page-action', action)
}
</script>

<style lang="scss" scoped>
.enhanced-breadcrumb {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-primary);
  padding: var(--spacing-4) var(--spacing-6);
  border-bottom: 1px solid var(--border-light);
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
  backdrop-filter: blur(8px);
  
  .breadcrumb-nav {
    flex: 1;
    min-width: 0;
    
    .el-breadcrumb {
      font-size: 0.875rem;
      
      .breadcrumb-home {
        .el-breadcrumb__inner {
          display: flex;
          align-items: center;
          gap: var(--spacing-1);
          color: var(--text-secondary);
          transition: color var(--transition-fast);
          
          &:hover {
            color: var(--primary-500);
          }
        }
      }
      
      .el-breadcrumb__item {
        .el-breadcrumb__inner {
          display: flex;
          align-items: center;
          gap: var(--spacing-1);
          font-weight: 500;
          
          &:not(.is-link) {
            color: var(--text-primary);
          }
        }
        
        &.is-current {
          .el-breadcrumb__inner {
            color: var(--primary-500);
            font-weight: 600;
          }
        }
      }
    }
  }
  
  .breadcrumb-actions {
    margin-left: var(--spacing-4);
    
    .quick-actions-btn {
      display: flex;
      align-items: center;
      gap: var(--spacing-1);
      color: var(--text-secondary);
      transition: color var(--transition-fast);
      
      &:hover {
        color: var(--primary-500);
      }
    }
  }
  
  .page-actions {
    margin-left: var(--spacing-4);
    
    .el-button-group {
      .el-button {
        border-radius: var(--radius-md);
        
        &:first-child {
          border-top-right-radius: 0;
          border-bottom-right-radius: 0;
        }
        
        &:last-child {
          border-top-left-radius: 0;
          border-bottom-left-radius: 0;
        }
        
        &:not(:first-child):not(:last-child) {
          border-radius: 0;
        }
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .enhanced-breadcrumb {
    padding: var(--mobile-space-3) var(--mobile-space-4);
    flex-wrap: wrap;
    gap: var(--mobile-space-3);
    
    .breadcrumb-nav {
      order: 1;
      width: 100%;
      
      .el-breadcrumb {
        font-size: 0.75rem;
        
        .el-breadcrumb__item {
          .el-breadcrumb__inner {
            span {
              max-width: 100px;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }
        }
      }
    }
    
    .breadcrumb-actions {
      order: 2;
      margin-left: 0;
    }
    
    .page-actions {
      order: 3;
      margin-left: 0;
      
      .el-button-group {
        display: flex;
        flex-wrap: wrap;
        gap: var(--mobile-space-2);
        
        .el-button {
          flex: 1;
          min-width: 80px;
          border-radius: var(--mobile-border-radius) !important;
        }
      }
    }
  }
}
</style>