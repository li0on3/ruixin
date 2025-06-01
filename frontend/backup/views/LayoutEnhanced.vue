<template>
  <div class="app-container">
    <!-- 移动端头部 -->
    <header class="app-header mobile-only">
      <div class="header-left">
        <button class="mobile-menu-btn" @click="toggleMobileSidebar" :class="{ 'is-active': isMobileSidebarOpen }">
          <el-icon><Menu /></el-icon>
        </button>
        <h1 class="header-title">{{ currentRoute.meta?.title || '瑞幸分销系统' }}</h1>
      </div>
      <div class="header-right">
        <button class="header-action" @click="toggleFullscreen" title="全屏">
          <el-icon><FullScreen /></el-icon>
        </button>
        <!-- 移动端通知按钮 -->
        <el-popover placement="bottom" :width="300" trigger="click" :persistent="false">
          <template #reference>
            <button class="header-action" title="通知">
              <el-icon><Bell /></el-icon>
              <span class="notification-badge" v-if="unreadCount > 0">{{ unreadCount }}</span>
            </button>
          </template>
          <div class="notification-panel mobile-notification">
            <div class="notification-header">
              <h3>通知</h3>
              <el-button text size="small" @click="markAllRead">全部已读</el-button>
            </div>
            <div class="notification-list">
              <div class="notification-item" v-for="item in notifications" :key="item.id">
                <div class="notification-icon" :class="item.type">
                  <el-icon>
                    <InfoFilled v-if="item.type === 'info'" />
                    <SuccessFilled v-if="item.type === 'success'" />
                    <WarningFilled v-if="item.type === 'warning'" />
                  </el-icon>
                </div>
                <div class="notification-content">
                  <p class="notification-title">{{ item.title }}</p>
                  <p class="notification-time">{{ item.time }}</p>
                </div>
              </div>
              <div class="notification-empty" v-if="notifications.length === 0">
                暂无新通知
              </div>
            </div>
          </div>
        </el-popover>
        <!-- 移动端用户菜单 -->
        <el-dropdown trigger="click" placement="bottom-end">
          <button class="header-action user-action">
            <div class="user-avatar-mobile">
              {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
            </div>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleProfile">
                <el-icon><User /></el-icon>
                个人中心
              </el-dropdown-item>
              <el-dropdown-item @click="handlePassword">
                <el-icon><Lock /></el-icon>
                修改密码
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <!-- 移动端遮罩 -->
    <div class="mobile-overlay" :class="{ 'is-active': isMobileSidebarOpen }" @click="closeMobileSidebar"></div>
    
    <!-- 侧边栏 -->
    <aside class="app-sidebar" :class="{ 
      'is-collapse': isCollapse, 
      'is-mobile-open': isMobileSidebarOpen 
    }">
      <!-- Logo区域 -->
      <div class="sidebar-logo">
        <div class="sidebar-logo-icon">
          <Coffee v-if="!isCollapse" />
          <span v-else>瑞</span>
        </div>
        <h1 class="sidebar-logo-text" v-if="!isCollapse">瑞幸分销系统</h1>
      </div>
      
      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          :collapse-transition="false"
          router
        >
          <!-- 仪表盘 -->
          <el-menu-item index="/dashboard">
            <el-icon><DataLine /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          
          <!-- 核心业务 -->
          <div class="menu-group" v-if="!isCollapse">
            <div class="menu-group-title">核心业务</div>
          </div>
          
          <el-menu-item index="/cards">
            <el-icon><CreditCard /></el-icon>
            <span>卡片管理</span>
          </el-menu-item>
          
          <el-menu-item index="/distributors">
            <el-icon><UserFilled /></el-icon>
            <span>分销商管理</span>
          </el-menu-item>
          
          <el-menu-item index="/orders">
            <el-icon><ShoppingCart /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          
          <!-- 数据分析 -->
          <div class="menu-group" v-if="!isCollapse">
            <div class="menu-group-title">数据分析</div>
          </div>
          
          <el-menu-item index="/statistics">
            <el-icon><TrendCharts /></el-icon>
            <span>数据统计</span>
          </el-menu-item>
          
          <el-sub-menu index="/finance">
            <template #title>
              <el-icon><Wallet /></el-icon>
              <span>财务管理</span>
            </template>
            <el-menu-item index="/finance/transactions">交易记录</el-menu-item>
            <el-menu-item index="/finance/withdrawals">提现管理</el-menu-item>
          </el-sub-menu>
          
          <!-- 商品管理 -->
          <div class="menu-group" v-if="!isCollapse">
            <div class="menu-group-title">商品管理</div>
          </div>
          
          <el-menu-item index="/luckin/prices">
            <el-icon><Coffee /></el-icon>
            <span>价格管理</span>
          </el-menu-item>
          
          <el-menu-item index="/products">
            <el-icon><ShoppingBag /></el-icon>
            <span>商品列表</span>
          </el-menu-item>
          
          <el-menu-item index="/products/available">
            <el-icon><ShoppingBag /></el-icon>
            <span>可用商品</span>
          </el-menu-item>
          
          <el-menu-item index="/stores">
            <el-icon><Location /></el-icon>
            <span>店铺查询</span>
          </el-menu-item>
          
          <!-- 系统管理 -->
          <div class="menu-group" v-if="!isCollapse">
            <div class="menu-group-title">系统管理</div>
          </div>
          
          <el-menu-item index="/admins" v-if="userStore.isSuperAdmin">
            <el-icon><User /></el-icon>
            <span>管理员管理</span>
          </el-menu-item>
          
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
      </nav>
      
      <!-- 底部版权信息 -->
      <div class="sidebar-footer" v-if="!isCollapse">
        <div class="footer-info">
          <p class="version">v1.0.0</p>
          <p class="copyright">© 2024 Ruixin System</p>
        </div>
      </div>
    </aside>
    
    <!-- 主内容区 -->
    <div class="app-content">
      <!-- 顶部导航栏（桌面版） -->
      <header class="app-header desktop-only">
        <div class="header-left">
          <!-- 菜单折叠按钮 -->
          <button class="menu-toggler" @click="toggleSidebar">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </button>
          
          <!-- 增强版面包屑导航 -->
          <EnhancedBreadcrumb
            :quick-actions="currentQuickActions"
            :page-actions="currentPageActions"
            @quick-action="handleQuickAction"
            @page-action="handlePageAction"
          />
        </div>
        
        <div class="header-right">
          <!-- 主题切换 -->
          <el-dropdown @command="handleThemeChange" trigger="click" placement="bottom-end">
            <button class="header-action" title="主题切换">
              <Sunny v-if="!isDarkTheme" />
              <Moon v-else />
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="light">
                  <Sunny />
                  浅色主题
                </el-dropdown-item>
                <el-dropdown-item command="dark">
                  <Moon />
                  深色主题
                </el-dropdown-item>
                <el-dropdown-item command="auto">
                  <Monitor />
                  跟随系统
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <!-- 全屏按钮 -->
          <button class="header-action" @click="toggleFullscreen" title="全屏">
            <FullScreen />
          </button>
          
          <!-- 通知中心 -->
          <el-popover placement="bottom" :width="380" trigger="click" :persistent="false">
            <template #reference>
              <button class="header-action" title="通知">
                <Bell />
                <span class="badge" v-if="unreadCount > 0"></span>
              </button>
            </template>
            <div class="notification-panel">
              <div class="notification-header">
                <h3>通知中心</h3>
                <el-button text @click="markAllRead">全部已读</el-button>
              </div>
              <div class="notification-list">
                <div class="notification-item" v-for="item in notifications" :key="item.id">
                  <div class="notification-icon" :class="item.type">
                    <InfoFilled v-if="item.type === 'info'" />
                    <SuccessFilled v-if="item.type === 'success'" />
                    <WarningFilled v-if="item.type === 'warning'" />
                  </div>
                  <div class="notification-content">
                    <p class="notification-title">{{ item.title }}</p>
                    <p class="notification-time">{{ item.time }}</p>
                  </div>
                </div>
                <div class="notification-empty" v-if="notifications.length === 0">
                  暂无新通知
                </div>
              </div>
            </div>
          </el-popover>
          
          <!-- 用户菜单 -->
          <el-dropdown trigger="click">
            <div class="user-dropdown">
              <div class="user-avatar">
                {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
              </div>
              <span class="user-name">{{ userStore.userInfo?.username || '用户' }}</span>
              <ArrowDown />
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleProfile">
                  <User />
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item @click="handlePassword">
                  <Lock />
                  修改密码
                </el-dropdown-item>
                <el-dropdown-item @click="handleSetting">
                  <Setting />
                  偏好设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <SwitchButton />
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      
      <!-- 页面主体 -->
      <main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" :key="route.fullPath" />
          </transition>
        </router-view>
      </main>
    </div>
    
    <!-- 快捷操作组件 -->
    <QuickActions />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  DataLine, CreditCard, UserFilled, ShoppingCart, TrendCharts,
  Wallet, Coffee, ShoppingBag, Location, User, Setting,
  Fold, Expand, House, FullScreen, Bell, ArrowDown,
  InfoFilled, SuccessFilled, WarningFilled, Lock, SwitchButton,
  Menu, Plus, Upload, Refresh, Download, Sunny, Moon, Monitor
} from '@element-plus/icons-vue'
import EnhancedBreadcrumb from '@/components/EnhancedBreadcrumb.vue'
import QuickActions from '@/components/QuickActions.vue'
import { useTheme } from '@/plugins/theme'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 主题相关
const { themeState, setTheme, isDarkTheme } = useTheme()

// 响应式数据
const isCollapse = ref(localStorage.getItem('sidebarCollapse') === 'true')
const isFullscreen = ref(false)
const isMobileSidebarOpen = ref(false)
const notifications = ref([
  { id: 1, type: 'success', title: '订单处理成功', time: '5分钟前' },
  { id: 2, type: 'info', title: '新的分销商申请', time: '1小时前' },
  { id: 3, type: 'warning', title: '库存预警提醒', time: '2小时前' }
])
const unreadCount = computed(() => notifications.value.length)

// 计算属性
const currentRoute = computed(() => route)
const activeMenu = computed(() => route.path)

// 当前页面的快捷操作配置
const currentQuickActions = computed(() => {
  const path = route.path
  const actions = []
  
  // 根据不同页面配置不同的快捷操作
  if (path === '/cards') {
    actions.push(
      { key: 'add-card', title: '添加卡片', icon: 'Plus' },
      { key: 'batch-import', title: '批量导入', icon: 'Upload' },
      { key: 'export-data', title: '导出数据', icon: 'Download', divided: true },
      { key: 'refresh-data', title: '刷新数据', icon: 'Refresh' }
    )
  } else if (path === '/distributors') {
    actions.push(
      { key: 'add-distributor', title: '添加分销商', icon: 'Plus' },
      { key: 'export-data', title: '导出数据', icon: 'Download', divided: true },
      { key: 'refresh-data', title: '刷新数据', icon: 'Refresh' }
    )
  } else if (path === '/orders') {
    actions.push(
      { key: 'refresh-orders', title: '刷新订单', icon: 'Refresh' },
      { key: 'export-orders', title: '导出订单', icon: 'Download', divided: true }
    )
  } else if (path === '/statistics') {
    actions.push(
      { key: 'refresh-stats', title: '刷新统计', icon: 'Refresh' },
      { key: 'export-report', title: '导出报表', icon: 'Download', divided: true }
    )
  }
  
  return actions
})

// 当前页面的页面操作配置
const currentPageActions = computed(() => {
  const path = route.path
  const actions = []
  
  // 根据不同页面配置不同的页面操作
  if (path === '/cards') {
    actions.push(
      { key: 'add', title: '添加', type: 'primary', icon: 'Plus' },
      { key: 'import', title: '导入', type: 'success', icon: 'Upload' }
    )
  } else if (path === '/distributors') {
    actions.push(
      { key: 'add', title: '添加', type: 'primary', icon: 'Plus' }
    )
  } else if (path === '/statistics') {
    actions.push(
      { key: 'export', title: '导出', type: 'primary', icon: 'Download' }
    )
  }
  
  return actions
})

// 方法
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
  localStorage.setItem('sidebarCollapse', isCollapse.value)
}

const toggleMobileSidebar = () => {
  isMobileSidebarOpen.value = !isMobileSidebarOpen.value
}

const closeMobileSidebar = () => {
  isMobileSidebarOpen.value = false
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

const markAllRead = () => {
  notifications.value = []
  ElMessage.success('已全部标记为已读')
}

const handleProfile = () => {
  router.push('/profile')
}

const handlePassword = () => {
  router.push('/password')
}

const handleSetting = () => {
  router.push('/preference')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await userStore.logout()
    ElMessage.success('退出成功')
    router.push('/login')
  } catch (error) {
    // 用户取消
  }
}

// 主题切换处理
const handleThemeChange = (command) => {
  setTheme(command)
  const themeNames = {
    light: '浅色主题',
    dark: '深色主题',
    auto: '跟随系统'
  }
  ElMessage.success(`已切换到${themeNames[command]}`)
}

// 处理快捷操作
const handleQuickAction = (action) => {
  console.log('Quick action:', action)
  
  switch (action.key) {
    case 'add-card':
    case 'add-distributor':
      ElMessage.info(`执行${action.title}操作`)
      break
    case 'batch-import':
      router.push('/cards/batch-import')
      break
    case 'export-data':
      handleExportData('cards')
      break
    case 'export-orders':
      handleExportData('orders')
      break
    case 'export-report':
      handleExportData('statistics')
      break
    case 'refresh-data':
    case 'refresh-orders':
    case 'refresh-stats':
      ElMessage.success('数据刷新成功')
      break
    default:
      ElMessage.info(`执行${action.title}操作`)
  }
}

// 处理页面操作
const handlePageAction = (action) => {
  console.log('Page action:', action)
  
  switch (action.key) {
    case 'add':
      ElMessage.info('打开添加对话框')
      break
    case 'import':
      router.push('/cards/batch-import')
      break
    case 'export':
      const path = route.path
      if (path === '/statistics') {
        handleExportData('statistics')
      } else if (path === '/cards') {
        handleExportData('cards')
      } else if (path === '/orders') {
        handleExportData('orders')
      } else {
        ElMessage.info('导出功能开发中...')
      }
      break
    default:
      ElMessage.info(`执行${action.title}操作`)
  }
}

// 处理导出数据
const handleExportData = async (dataType) => {
  try {
    // 动态导入导出函数
    const { exportPageData } = await import('@/utils/export')
    
    // 这里应该调用相应的API获取数据，这里先用示例数据
    let data = []
    
    switch (dataType) {
      case 'cards':
        // 实际项目中应该调用 getCards API
        ElMessage.info('卡片数据导出功能开发中...')
        return
      case 'orders':
        // 实际项目中应该调用 getOrders API
        ElMessage.info('订单数据导出功能开发中...')
        return
      case 'statistics':
        // 实际项目中应该调用统计相关的API
        ElMessage.info('统计数据导出功能开发中...')
        return
      default:
        ElMessage.warning('未知的导出类型')
        return
    }
    
    exportPageData(dataType, data)
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}

// 监听全屏变化
const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement
}

onMounted(() => {
  document.addEventListener('fullscreenchange', handleFullscreenChange)
})

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
})
</script>

<style lang="scss" scoped>
// 页面切换动画
.fade-slide-enter-active {
  transition: all 0.3s ease-out;
}

.fade-slide-leave-active {
  transition: all 0.2s ease-in;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

// 菜单分组样式
.menu-group {
  padding: var(--spacing-4) var(--spacing-3) var(--spacing-2);
  
  .menu-group-title {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
}

// 通知面板样式
.notification-panel {
  .notification-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-bottom: var(--spacing-3);
    margin-bottom: var(--spacing-3);
    border-bottom: 1px solid var(--border-light);
    
    h3 {
      margin: 0;
      font-size: 1rem;
      font-weight: 600;
    }
  }
  
  .notification-list {
    max-height: 400px;
    overflow-y: auto;
  }
  
  .notification-item {
    display: flex;
    align-items: flex-start;
    gap: var(--spacing-3);
    padding: var(--spacing-3);
    border-radius: var(--radius-lg);
    cursor: pointer;
    transition: background var(--transition-fast);
    
    &:hover {
      background: var(--bg-secondary);
    }
    
    .notification-icon {
      width: 40px;
      height: 40px;
      border-radius: var(--radius-lg);
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
      
      &.info {
        background: var(--info-light);
        color: var(--info);
      }
      
      &.success {
        background: var(--success-light);
        color: var(--success);
      }
      
      &.warning {
        background: var(--warning-light);
        color: var(--warning);
      }
    }
    
    .notification-content {
      flex: 1;
      
      .notification-title {
        margin: 0 0 var(--spacing-1);
        font-weight: 500;
        color: var(--text-primary);
      }
      
      .notification-time {
        margin: 0;
        font-size: 0.875rem;
        color: var(--text-tertiary);
      }
    }
  }
  
  .notification-empty {
    padding: var(--spacing-8) var(--spacing-4);
    text-align: center;
    color: var(--text-tertiary);
  }
}

// 底部版权样式
.sidebar-footer {
  padding: var(--spacing-4);
  border-top: 1px solid var(--border-light);
  
  .footer-info {
    text-align: center;
    
    p {
      margin: var(--spacing-1) 0;
      font-size: 0.75rem;
      color: var(--text-tertiary);
    }
    
    .version {
      font-weight: 600;
      color: var(--primary-500);
    }
  }
}

// 移动端专用样式
@media (max-width: 768px) {
  .app-container {
    flex-direction: column;
  }
  
  .app-header.mobile-only {
    display: flex !important;
  }
  
  .app-header.desktop-only {
    display: none !important;
  }
  
  .app-sidebar {
    position: fixed !important;
    top: 0;
    left: -100%;
    width: 280px !important;
    height: 100vh;
    z-index: 1050;
    transition: left 0.3s ease;
    box-shadow: var(--shadow-xl);
    
    &.is-mobile-open {
      left: 0;
    }
  }
  
  .mobile-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1030;
    opacity: 0;
    visibility: hidden;
    transition: all 0.3s ease;
    
    &.is-active {
      opacity: 1;
      visibility: visible;
    }
  }
  
  .app-main {
    margin-left: 0 !important;
    padding-top: calc(var(--mobile-header-height) + var(--mobile-space-4));
    padding-bottom: calc(var(--mobile-space-8) + 56px); // 为快捷操作按钮留出空间
  }
  
  // 移动端头部样式
  .header-left {
    display: flex;
    align-items: center;
    gap: var(--mobile-space-3);
    
    .mobile-menu-btn {
      width: var(--mobile-touch-target);
      height: var(--mobile-touch-target);
      border: none;
      background: none;
      color: var(--text-primary);
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: var(--radius-md);
      cursor: pointer;
      transition: all 0.2s ease;
      
      &:hover {
        background: var(--gray-100);
      }
      
      &:active {
        background: var(--gray-200);
      }
    }
    
    .header-title {
      font-size: var(--mobile-text-lg);
      font-weight: 600;
      color: var(--text-primary);
      margin: 0;
    }
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: var(--mobile-space-2);
    
    .header-action {
      width: var(--mobile-touch-target);
      height: var(--mobile-touch-target);
      border: none;
      background: none;
      color: var(--text-secondary);
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: var(--radius-md);
      cursor: pointer;
      transition: all 0.2s ease;
      position: relative;
      
      &:hover {
        background: var(--gray-100);
        color: var(--text-primary);
      }
      
      .notification-badge {
        position: absolute;
        top: 4px;
        right: 4px;
        background: var(--danger);
        color: white;
        border-radius: 50%;
        min-width: 18px;
        height: 18px;
        font-size: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
      }
      
      .user-avatar-mobile {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        background: var(--primary-500);
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        font-weight: 600;
      }
    }
  }
  
  // 移动端通知面板
  .mobile-notification {
    .notification-header {
      padding-bottom: var(--mobile-space-3);
      margin-bottom: var(--mobile-space-3);
      
      h3 {
        font-size: var(--mobile-text-base);
      }
    }
    
    .notification-item {
      padding: var(--mobile-space-3);
      
      .notification-icon {
        width: 36px;
        height: 36px;
      }
      
      .notification-content {
        .notification-title {
          font-size: var(--mobile-text-sm);
        }
        
        .notification-time {
          font-size: var(--mobile-text-xs);
        }
      }
    }
  }
}
</style>