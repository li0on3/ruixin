<template>
  <div class="app-layout" :class="layoutClasses">
    <!-- 移动端遮罩层 -->
    <div 
      class="mobile-overlay" 
      :class="{ 'is-visible': isMobileSidebarOpen }" 
      @click="closeMobileSidebar"
    ></div>
    
    <!-- 侧边栏 -->
    <aside class="app-sidebar" :class="sidebarClasses">
      <!-- Logo 区域 -->
      <div class="sidebar-header">
        <router-link to="/dashboard" class="logo-container" @click="closeMobileSidebar">
          <div class="logo-icon">
            <el-icon><Coffee /></el-icon>
          </div>
          <div class="logo-content" v-if="!isCollapse">
            <h1 class="logo-title">瑞幸分销</h1>
            <p class="logo-subtitle">Distribution System</p>
          </div>
        </router-link>
        
        <!-- 侧边栏控制按钮（仅桌面端） -->
        <button 
          class="sidebar-toggle desktop-only" 
          @click="toggleSidebar"
          :title="isCollapse ? '展开侧边栏' : '收起侧边栏'"
        >
          <el-icon>
            <ArrowLeft v-if="!isCollapse" />
            <ArrowRight v-else />
          </el-icon>
        </button>
      </div>
      
      <!-- 导航菜单 -->
      <nav class="sidebar-nav custom-scrollbar">
        <div class="nav-section">
          <!-- 核心业务 -->
          <div class="nav-group">
            <div class="nav-group-header" v-if="!isCollapse">
              <el-icon><TrendCharts /></el-icon>
              <span>核心业务</span>
            </div>
            <div class="nav-items">
              <router-link 
                to="/dashboard" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/dashboard') }"
              >
                <el-icon><DataBoard /></el-icon>
                <span>数据仪表盘</span>
                <div class="nav-indicator"></div>
              </router-link>
              
              <router-link 
                to="/orders" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/orders') }"
              >
                <el-icon><ShoppingCart /></el-icon>
                <span>订单管理</span>
                <div class="nav-indicator"></div>
              </router-link>
              
              <router-link 
                to="/statistics" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/statistics') }"
              >
                <el-icon><TrendCharts /></el-icon>
                <span>数据统计</span>
                <div class="nav-indicator"></div>
              </router-link>
            </div>
          </div>
          
          <!-- 资源管理 -->
          <div class="nav-group">
            <div class="nav-group-header" v-if="!isCollapse">
              <el-icon><Box /></el-icon>
              <span>资源管理</span>
            </div>
            <div class="nav-items">
              <router-link 
                to="/cards" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/cards') }"
              >
                <el-icon><CreditCard /></el-icon>
                <span>卡片管理</span>
                <div class="nav-indicator"></div>
              </router-link>
              
              <div class="nav-submenu" :class="{ 'is-expanded': isSubmenuOpen('products') }">
                <div class="nav-item submenu-trigger" @click="toggleSubmenu('products')">
                  <el-icon><Goods /></el-icon>
                  <span>商品管理</span>
                  <el-icon class="submenu-arrow">
                    <ArrowDown v-if="isSubmenuOpen('products')" />
                    <ArrowRight v-else />
                  </el-icon>
                </div>
                <div class="submenu-items">
                  <router-link 
                    to="/products" 
                    class="nav-item submenu-item" 
                    @click="closeMobileSidebar"
                    :class="{ 'is-active': isRouteActive('/products') }"
                  >
                    <span>商品列表</span>
                  </router-link>
                  <router-link 
                    to="/products/available" 
                    class="nav-item submenu-item" 
                    @click="closeMobileSidebar"
                    :class="{ 'is-active': isRouteActive('/products/available') }"
                  >
                    <span>可用商品</span>
                  </router-link>
                  <router-link 
                    to="/luckin/prices" 
                    class="nav-item submenu-item" 
                    @click="closeMobileSidebar"
                    :class="{ 'is-active': isRouteActive('/luckin/prices') }"
                  >
                    <span>价格管理</span>
                  </router-link>
                </div>
              </div>
              
              <router-link 
                to="/stores" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/stores') }"
              >
                <el-icon><Shop /></el-icon>
                <span>门店管理</span>
                <div class="nav-indicator"></div>
              </router-link>
            </div>
          </div>
          
          <!-- 客户关系 -->
          <div class="nav-group">
            <div class="nav-group-header" v-if="!isCollapse">
              <el-icon><Avatar /></el-icon>
              <span>客户关系</span>
            </div>
            <div class="nav-items">
              <router-link 
                to="/distributors" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/distributors') }"
              >
                <el-icon><UserFilled /></el-icon>
                <span>分销商管理</span>
                <div class="nav-indicator"></div>
              </router-link>
            </div>
          </div>
          
          <!-- 财务中心 -->
          <div class="nav-group">
            <div class="nav-group-header" v-if="!isCollapse">
              <el-icon><Money /></el-icon>
              <span>财务中心</span>
            </div>
            <div class="nav-items">
              <div class="nav-submenu" :class="{ 'is-expanded': isSubmenuOpen('finance') }">
                <div class="nav-item submenu-trigger" @click="toggleSubmenu('finance')">
                  <el-icon><Wallet /></el-icon>
                  <span>财务管理</span>
                  <el-icon class="submenu-arrow">
                    <ArrowDown v-if="isSubmenuOpen('finance')" />
                    <ArrowRight v-else />
                  </el-icon>
                </div>
                <div class="submenu-items">
                  <router-link 
                    to="/finance/transactions" 
                    class="nav-item submenu-item" 
                    @click="closeMobileSidebar"
                    :class="{ 'is-active': isRouteActive('/finance/transactions') }"
                  >
                    <span>交易记录</span>
                  </router-link>
                  <router-link 
                    to="/finance/withdrawals" 
                    class="nav-item submenu-item" 
                    @click="closeMobileSidebar"
                    :class="{ 'is-active': isRouteActive('/finance/withdrawals') }"
                  >
                    <span>提现管理</span>
                  </router-link>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 系统设置 -->
          <div class="nav-group">
            <div class="nav-group-header" v-if="!isCollapse">
              <el-icon><Tools /></el-icon>
              <span>系统设置</span>
            </div>
            <div class="nav-items">
              <router-link 
                to="/admins" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/admins') }"
                v-if="userStore.isSuperAdmin"
              >
                <el-icon><User /></el-icon>
                <span>管理员</span>
                <div class="nav-indicator"></div>
              </router-link>
              
              <router-link 
                to="/settings" 
                class="nav-item" 
                @click="closeMobileSidebar"
                :class="{ 'is-active': isRouteActive('/settings') }"
              >
                <el-icon><Setting /></el-icon>
                <span>系统设置</span>
                <div class="nav-indicator"></div>
              </router-link>
            </div>
          </div>
        </div>
      </nav>
      
      <!-- 侧边栏底部 -->
      <div class="sidebar-footer" v-if="!isCollapse">
        <div class="user-card">
          <div class="user-avatar">
            {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
          </div>
          <div class="user-info">
            <div class="user-name">{{ userStore.userInfo?.username || '用户' }}</div>
            <div class="user-role">{{ userStore.isSuperAdmin ? '超级管理员' : '管理员' }}</div>
          </div>
          <el-dropdown trigger="click" placement="top-start">
            <button class="user-menu-btn">
              <el-icon><MoreFilled /></el-icon>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleProfile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item @click="handleSettings">
                  <el-icon><Setting /></el-icon>
                  偏好设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        
        <div class="app-version">
          <span class="version-label">Version</span>
          <span class="version-number">1.2.0</span>
        </div>
      </div>
    </aside>
    
    <!-- 主内容区 -->
    <div class="app-main">
      <!-- 顶部导航栏 -->
      <header class="app-header">
        <!-- 移动端菜单按钮 -->
        <button 
          class="mobile-menu-btn mobile-only" 
          @click="toggleMobileSidebar"
          :class="{ 'is-active': isMobileSidebarOpen }"
        >
          <span class="hamburger-line"></span>
          <span class="hamburger-line"></span>
          <span class="hamburger-line"></span>
        </button>
        
        <!-- 面包屑导航 -->
        <div class="breadcrumb-section">
          <div class="page-title">{{ currentPageTitle }}</div>
          <nav class="breadcrumb-nav">
            <ol class="breadcrumb">
              <li v-for="(item, index) in breadcrumbItems" :key="index" class="breadcrumb-item">
                <router-link 
                  v-if="item.path && index < breadcrumbItems.length - 1" 
                  :to="item.path"
                  class="breadcrumb-link"
                >
                  {{ item.title }}
                </router-link>
                <span v-else class="breadcrumb-current">{{ item.title }}</span>
                <el-icon v-if="index < breadcrumbItems.length - 1" class="breadcrumb-separator">
                  <ArrowRight />
                </el-icon>
              </li>
            </ol>
          </nav>
        </div>
        
        <!-- 全局搜索 -->
        <div class="search-section desktop-only">
          <div class="search-container">
            <el-input
              v-model="searchQuery"
              placeholder="搜索功能、订单、分销商..."
              :prefix-icon="Search"
              clearable
              @keyup.enter="handleSearch"
              @clear="clearSearch"
              class="global-search"
            />
            <div class="search-results" v-if="searchResults.length > 0">
              <div 
                class="search-item" 
                v-for="item in searchResults" 
                :key="item.id"
                @click="navigateToSearchResult(item)"
              >
                <el-icon class="search-item-icon">
                  <component :is="item.icon" />
                </el-icon>
                <div class="search-item-content">
                  <div class="search-item-title">{{ item.title }}</div>
                  <div class="search-item-type">{{ item.type }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 右侧操作区 -->
        <div class="header-actions">
          <!-- 主题切换 -->
          <button 
            class="action-btn" 
            @click="toggleTheme"
            :title="isDarkTheme ? '切换到浅色主题' : '切换到深色主题'"
          >
            <el-icon>
              <Sunny v-if="isDarkTheme" />
              <Moon v-else />
            </el-icon>
          </button>
          
          <!-- 全屏切换 -->
          <button 
            class="action-btn desktop-only" 
            @click="toggleFullscreen"
            :title="isFullscreen ? '退出全屏' : '进入全屏'"
          >
            <el-icon>
              <Aim v-if="isFullscreen" />
              <FullScreen v-else />
            </el-icon>
          </button>
          
          <!-- 通知中心 -->
          <el-popover 
            placement="bottom-end" 
            :width="400" 
            trigger="click" 
            :persistent="false"
            popper-class="notification-popover"
          >
            <template #reference>
              <button class="action-btn notification-btn" title="通知中心">
                <el-icon><Bell /></el-icon>
                <el-badge 
                  :value="unreadNotifications" 
                  :hidden="unreadNotifications === 0"
                  class="notification-badge"
                />
              </button>
            </template>
            
            <div class="notification-panel">
              <div class="notification-header">
                <h3>通知中心</h3>
                <div class="notification-actions">
                  <el-button text size="small" @click="markAllAsRead" :disabled="unreadNotifications === 0">
                    全部已读
                  </el-button>
                </div>
              </div>
              
              <div class="notification-tabs">
                <button 
                  class="notification-tab" 
                  :class="{ 'is-active': activeNotificationTab === 'all' }"
                  @click="activeNotificationTab = 'all'"
                >
                  全部 ({{ notifications.length }})
                </button>
                <button 
                  class="notification-tab" 
                  :class="{ 'is-active': activeNotificationTab === 'unread' }"
                  @click="activeNotificationTab = 'unread'"
                >
                  未读 ({{ unreadNotifications }})
                </button>
              </div>
              
              <div class="notification-list custom-scrollbar">
                <transition-group name="notification" tag="div">
                  <div 
                    class="notification-item" 
                    v-for="notification in filteredNotifications" 
                    :key="notification.id"
                    :class="{ 'is-read': notification.read }"
                    @click="handleNotificationClick(notification)"
                  >
                    <div class="notification-icon" :class="notification.type">
                      <el-icon>
                        <SuccessFilled v-if="notification.type === 'success'" />
                        <WarningFilled v-if="notification.type === 'warning'" />
                        <InfoFilled v-if="notification.type === 'info'" />
                        <CircleCloseFilled v-if="notification.type === 'error'" />
                      </el-icon>
                    </div>
                    <div class="notification-content">
                      <div class="notification-title">{{ notification.title }}</div>
                      <div class="notification-description">{{ notification.description }}</div>
                      <div class="notification-time">{{ formatNotificationTime(notification.createdAt) }}</div>
                    </div>
                    <button 
                      class="notification-close" 
                      @click.stop="removeNotification(notification.id)"
                      title="删除通知"
                    >
                      <el-icon><Close /></el-icon>
                    </button>
                  </div>
                </transition-group>
                
                <div class="notification-empty" v-if="filteredNotifications.length === 0">
                  <el-icon><Bell /></el-icon>
                  <p>{{ activeNotificationTab === 'unread' ? '暂无未读通知' : '暂无通知' }}</p>
                </div>
              </div>
            </div>
          </el-popover>
          
          <!-- 用户菜单 -->
          <el-dropdown trigger="click" placement="bottom-end">
            <div class="user-dropdown">
              <div class="user-avatar">
                {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
              </div>
              <div class="user-info desktop-only">
                <div class="user-name">{{ userStore.userInfo?.username || '用户' }}</div>
                <div class="user-role">{{ userStore.isSuperAdmin ? '超级管理员' : '管理员' }}</div>
              </div>
              <el-icon class="user-arrow desktop-only"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleProfile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item @click="handleChangePassword">
                  <el-icon><Lock /></el-icon>
                  修改密码
                </el-dropdown-item>
                <el-dropdown-item @click="handleSettings">
                  <el-icon><Setting /></el-icon>
                  偏好设置
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
      
      <!-- 页面内容 -->
      <main class="app-content">
        <!-- 页面工具栏 -->

        <!-- 路由视图 -->
        <div class="route-container">
          <router-view v-slot="{ Component }">
            <transition name="page-transition" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
    
    <!-- 快捷操作浮动按钮 -->
    <div class="quick-actions-fab mobile-only" v-if="currentQuickActions.length > 0">
      <el-button type="primary" circle size="large" @click="showQuickActionsMenu = !showQuickActionsMenu">
        <el-icon><Plus /></el-icon>
      </el-button>
      
      <transition name="quick-actions">
        <div class="quick-actions-menu" v-if="showQuickActionsMenu">
          <button 
            class="quick-action-item" 
            v-for="action in currentQuickActions" 
            :key="action.key"
            @click="handleQuickAction(action)"
          >
            <el-icon><component :is="action.icon" /></el-icon>
            <span>{{ action.title }}</span>
          </button>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Coffee, ArrowLeft, ArrowRight, TrendCharts, DataBoard, ShoppingCart,
  Box, CreditCard, Goods, Shop, Avatar, UserFilled, Money, Wallet,
  Tools, User, Setting, MoreFilled, SwitchButton, Search, Sunny, Moon,
  Aim, FullScreen, Bell, ArrowDown, ArrowUp, Close, Plus,
  SuccessFilled, WarningFilled, InfoFilled, CircleCloseFilled
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 响应式状态
const isCollapse = ref(localStorage.getItem('sidebar-collapsed') === 'true')
const isMobileSidebarOpen = ref(false)
const isDarkTheme = ref(localStorage.getItem('theme') === 'dark')
const isFullscreen = ref(false)
const expandedSubmenus = ref(new Set())
const searchQuery = ref('')
const searchResults = ref([])
const notifications = ref([
  {
    id: 1,
    type: 'success',
    title: '订单处理成功',
    description: '订单 #DO20250131001 已成功处理完成',
    createdAt: new Date(Date.now() - 5 * 60 * 1000),
    read: false
  },
  {
    id: 2,
    type: 'info',
    title: '新分销商申请',
    description: '用户"张三"提交了分销商申请，请及时审核',
    createdAt: new Date(Date.now() - 1 * 60 * 60 * 1000),
    read: false
  },
  {
    id: 3,
    type: 'warning',
    title: '库存预警',
    description: '商品"经典美式咖啡"库存不足，请及时补充',
    createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000),
    read: true
  }
])
const activeNotificationTab = ref('all')
const showQuickActionsMenu = ref(false)

// 计算属性
const layoutClasses = computed(() => ({
  'sidebar-collapsed': isCollapse.value,
  'sidebar-mobile-open': isMobileSidebarOpen.value,
  'dark-theme': isDarkTheme.value,
  'fullscreen': isFullscreen.value
}))

const sidebarClasses = computed(() => ({
  'is-collapsed': isCollapse.value,
  'is-mobile-open': isMobileSidebarOpen.value
}))

const currentPageTitle = computed(() => {
  return route.meta?.title || '瑞幸分销系统'
})

const breadcrumbItems = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  const items = matched.map(item => ({
    title: item.meta.title,
    path: item.path
  }))
  
  // 添加首页
  if (items.length > 0 && items[0].path !== '/dashboard') {
    items.unshift({ title: '首页', path: '/dashboard' })
  }
  
  return items
})

const unreadNotifications = computed(() => {
  return notifications.value.filter(n => !n.read).length
})

const filteredNotifications = computed(() => {
  if (activeNotificationTab.value === 'unread') {
    return notifications.value.filter(n => !n.read)
  }
  return notifications.value
})

const currentPageActions = computed(() => {
  // 根据当前路由返回页面操作按钮
  const path = route.path
  const actions = []
  
  if (path === '/cards') {
    actions.push(
      { key: 'add', title: '添加卡片', type: 'primary', icon: 'Plus' },
      { key: 'import', title: '批量导入', type: 'success', icon: 'Upload' }
    )
  } else if (path === '/distributors') {
    actions.push(
      { key: 'add', title: '添加分销商', type: 'primary', icon: 'Plus' }
    )
  } else if (path === '/orders') {
    actions.push(
      { key: 'refresh', title: '刷新', icon: 'Refresh' },
      { key: 'export', title: '导出', type: 'primary', icon: 'Download' }
    )
  }
  
  return actions
})

const currentQuickActions = computed(() => {
  // 移动端快捷操作
  const path = route.path
  const actions = []
  
  if (path === '/cards') {
    actions.push(
      { key: 'add', title: '添加', icon: 'Plus' },
      { key: 'scan', title: '扫码', icon: 'Camera' },
      { key: 'import', title: '导入', icon: 'Upload' }
    )
  } else if (path === '/orders') {
    actions.push(
      { key: 'refresh', title: '刷新', icon: 'Refresh' },
      { key: 'search', title: '搜索', icon: 'Search' }
    )
  }
  
  return actions
})

// 方法
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
  localStorage.setItem('sidebar-collapsed', isCollapse.value.toString())
}

const toggleMobileSidebar = () => {
  isMobileSidebarOpen.value = !isMobileSidebarOpen.value
}

const closeMobileSidebar = () => {
  isMobileSidebarOpen.value = false
  showQuickActionsMenu.value = false
}

const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value
  localStorage.setItem('theme', isDarkTheme.value ? 'dark' : 'light')
  document.documentElement.setAttribute('data-theme', isDarkTheme.value ? 'dark' : 'light')
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const isRouteActive = (path) => {
  return route.path.startsWith(path)
}

const isSubmenuOpen = (key) => {
  return expandedSubmenus.value.has(key)
}

const toggleSubmenu = (key) => {
  if (expandedSubmenus.value.has(key)) {
    expandedSubmenus.value.delete(key)
  } else {
    expandedSubmenus.value.add(key)
  }
}

const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  
  // 模拟搜索结果
  searchResults.value = [
    { id: 1, title: '卡片管理', type: '功能模块', icon: 'CreditCard', path: '/cards' },
    { id: 2, title: '订单 #DO20250131001', type: '订单', icon: 'ShoppingCart', path: '/orders' },
    { id: 3, title: '分销商：张三', type: '分销商', icon: 'UserFilled', path: '/distributors' }
  ].filter(item => 
    item.title.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
}

const clearSearch = () => {
  searchResults.value = []
}

const navigateToSearchResult = (item) => {
  if (item.path) {
    router.push(item.path)
  }
  searchQuery.value = ''
  searchResults.value = []
}

const formatNotificationTime = (date) => {
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  
  const days = Math.floor(hours / 24)
  return `${days}天前`
}

const handleNotificationClick = (notification) => {
  notification.read = true
  // 根据通知类型跳转到相应页面
}

const removeNotification = (id) => {
  const index = notifications.value.findIndex(n => n.id === id)
  if (index > -1) {
    notifications.value.splice(index, 1)
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
  ElMessage.success('已全部标记为已读')
}

const handlePageAction = (action) => {
  console.log('Page action:', action)
  ElMessage.info(`执行${action.title}操作`)
}

const handleQuickAction = (action) => {
  console.log('Quick action:', action)
  ElMessage.info(`执行${action.title}操作`)
  showQuickActionsMenu.value = false
}

const handleProfile = () => {
  router.push('/profile')
}

const handleChangePassword = () => {
  router.push('/change-password')
}

const handleSettings = () => {
  router.push('/settings')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '退出确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await userStore.logout()
    router.push('/login')
    ElMessage.success('退出成功')
  } catch (error) {
    // 用户取消
  }
}

// 监听全屏状态变化
const handleFullscreenChange = () => {
  isFullscreen.value = !!document.fullscreenElement
}

// 监听点击外部关闭搜索结果
const handleClickOutside = (event) => {
  const searchContainer = event.target.closest('.search-container')
  if (!searchContainer) {
    searchResults.value = []
  }
}

// 生命周期
onMounted(() => {
  document.addEventListener('fullscreenchange', handleFullscreenChange)
  document.addEventListener('click', handleClickOutside)
  
  // 应用保存的主题
  if (isDarkTheme.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
  }
  
  // 监听搜索输入
  watch(searchQuery, (newValue) => {
    if (newValue.trim()) {
      handleSearch()
    } else {
      clearSearch()
    }
  })
})

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', handleFullscreenChange)
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style lang="scss" scoped>
@import '@/assets/styles/design-tokens.scss';
@import '@/assets/styles/animations-enhanced.scss';

.app-layout {
  display: flex;
  height: 100vh;
  background: var(--bg-primary);
  transition: all var(--transition-base);
  
  // 移动端遮罩
  .mobile-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: var(--z-modal-backdrop);
    opacity: 0;
    visibility: hidden;
    transition: all var(--transition-base);
    
    &.is-visible {
      opacity: 1;
      visibility: visible;
    }
  }
  
  // 侧边栏
  .app-sidebar {
    width: 280px;
    background: var(--bg-primary);
    border-right: 1px solid var(--border-light);
    display: flex;
    flex-direction: column;
    transition: all var(--transition-base);
    z-index: var(--z-fixed);
    box-shadow: var(--shadow-sm);
    
    &.is-collapsed {
      width: 80px;
      
      .logo-content {
        opacity: 0;
        transform: translateX(-20px);
      }
      
      .nav-group-header {
        opacity: 0;
        transform: translateX(-20px);
      }
      
      .nav-item span {
        opacity: 0;
        transform: translateX(-20px);
      }
    }
    
    // 侧边栏头部
    .sidebar-header {
      padding: var(--spacing-6) var(--spacing-4);
      border-bottom: 1px solid var(--border-light);
      display: flex;
      align-items: center;
      justify-content: space-between;
      min-height: 80px;
      
      .logo-container {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        text-decoration: none;
        flex: 1;
        
        .logo-icon {
          width: 48px;
          height: 48px;
          background: linear-gradient(135deg, var(--primary-500), var(--primary-600));
          border-radius: var(--radius-xl);
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-size: 24px;
          flex-shrink: 0;
          box-shadow: var(--shadow-md);
          
          .el-icon {
            font-size: 24px;
          }
        }
        
        .logo-content {
          transition: all var(--transition-base);
          
          .logo-title {
            font-size: var(--text-xl);
            font-weight: var(--font-bold);
            color: var(--text-primary);
            margin: 0;
            line-height: 1.2;
          }
          
          .logo-subtitle {
            font-size: var(--text-xs);
            color: var(--text-tertiary);
            margin: 0;
            text-transform: uppercase;
            letter-spacing: 0.5px;
          }
        }
      }
      
      .sidebar-toggle {
        width: 32px;
        height: 32px;
        border: none;
        background: var(--bg-secondary);
        border-radius: var(--radius-md);
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--text-secondary);
        cursor: pointer;
        transition: all var(--transition-fast);
        
        &:hover {
          background: var(--bg-hover);
          color: var(--text-primary);
        }
      }
    }
    
    // 导航区域
    .sidebar-nav {
      flex: 1;
      overflow-y: auto;
      padding: var(--spacing-4) var(--spacing-2);
      
      .nav-section {
        .nav-group {
          margin-bottom: var(--spacing-6);
          
          &:last-child {
            margin-bottom: 0;
          }
          
          .nav-group-header {
            display: flex;
            align-items: center;
            gap: var(--spacing-2);
            padding: var(--spacing-2) var(--spacing-4);
            margin-bottom: var(--spacing-2);
            font-size: var(--text-xs);
            font-weight: var(--font-semibold);
            color: var(--text-tertiary);
            text-transform: uppercase;
            letter-spacing: 0.5px;
            transition: all var(--transition-base);
            
            .el-icon {
              font-size: 14px;
            }
          }
          
          .nav-items {
            .nav-item {
              display: flex;
              align-items: center;
              gap: var(--spacing-3);
              padding: var(--spacing-3) var(--spacing-4);
              margin: var(--spacing-1) 0;
              border-radius: var(--radius-lg);
              text-decoration: none;
              color: var(--text-secondary);
              font-weight: var(--font-medium);
              transition: all var(--transition-fast);
              position: relative;
              overflow: hidden;
              cursor: pointer;
              
              .el-icon {
                font-size: 20px;
                flex-shrink: 0;
              }
              
              span {
                flex: 1;
                transition: all var(--transition-base);
              }
              
              .nav-indicator {
                position: absolute;
                left: 0;
                top: 50%;
                width: 4px;
                height: 0;
                background: var(--primary-500);
                border-radius: 0 var(--radius-base) var(--radius-base) 0;
                transform: translateY(-50%);
                transition: all var(--transition-base);
              }
              
              &:hover {
                background: var(--bg-hover);
                color: var(--text-primary);
                transform: translateX(2px);
              }
              
              &.is-active {
                background: var(--primary-50);
                color: var(--primary-600);
                
                .nav-indicator {
                  height: 20px;
                }
              }
            }
            
            .nav-submenu {
              .submenu-trigger {
                .submenu-arrow {
                  margin-left: auto;
                  font-size: 16px;
                  transition: all var(--transition-fast);
                }
              }
              
              &.is-expanded {
                .submenu-trigger .submenu-arrow {
                  transform: rotate(0deg);
                }
                
                .submenu-items {
                  max-height: 200px;
                  opacity: 1;
                }
              }
              
              .submenu-items {
                max-height: 0;
                opacity: 0;
                overflow: hidden;
                transition: all var(--transition-base);
                padding-left: var(--spacing-8);
                
                .submenu-item {
                  padding: var(--spacing-2) var(--spacing-4);
                  font-size: var(--text-sm);
                  
                  &:hover {
                    transform: translateX(4px);
                  }
                }
              }
            }
          }
        }
      }
    }
    
    // 侧边栏底部
    .sidebar-footer {
      padding: var(--spacing-4);
      border-top: 1px solid var(--border-light);
      
      .user-card {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        padding: var(--spacing-3);
        background: var(--bg-secondary);
        border-radius: var(--radius-lg);
        margin-bottom: var(--spacing-4);
        
        .user-avatar {
          width: 40px;
          height: 40px;
          background: var(--primary-500);
          color: white;
          border-radius: var(--radius-lg);
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: var(--font-semibold);
          flex-shrink: 0;
        }
        
        .user-info {
          flex: 1;
          min-width: 0;
          
          .user-name {
            font-weight: var(--font-medium);
            color: var(--text-primary);
            font-size: var(--text-sm);
            margin-bottom: 2px;
          }
          
          .user-role {
            font-size: var(--text-xs);
            color: var(--text-tertiary);
          }
        }
        
        .user-menu-btn {
          width: 24px;
          height: 24px;
          border: none;
          background: none;
          color: var(--text-tertiary);
          border-radius: var(--radius-md);
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          transition: all var(--transition-fast);
          
          &:hover {
            background: var(--bg-hover);
            color: var(--text-primary);
          }
        }
      }
      
      .app-version {
        display: flex;
        align-items: center;
        justify-content: space-between;
        font-size: var(--text-xs);
        color: var(--text-tertiary);
        
        .version-number {
          font-weight: var(--font-medium);
          color: var(--primary-500);
        }
      }
    }
  }
  
  // 主内容区
  .app-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
    
    // 顶部导航栏
    .app-header {
      height: 80px;
      background: var(--bg-primary);
      border-bottom: 1px solid var(--border-light);
      display: flex;
      align-items: center;
      gap: var(--spacing-4);
      padding: 0 var(--spacing-6);
      position: sticky;
      top: 0;
      z-index: var(--z-sticky);
      box-shadow: var(--shadow-sm);
      
      .mobile-menu-btn {
        width: 40px;
        height: 40px;
        border: none;
        background: none;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 4px;
        cursor: pointer;
        border-radius: var(--radius-md);
        transition: all var(--transition-fast);
        
        .hamburger-line {
          width: 20px;
          height: 2px;
          background: var(--text-secondary);
          transition: all var(--transition-fast);
          
          &:nth-child(1) {
            transform-origin: left center;
          }
          
          &:nth-child(3) {
            transform-origin: left center;
          }
        }
        
        &:hover {
          background: var(--bg-hover);
          
          .hamburger-line {
            background: var(--text-primary);
          }
        }
        
        &.is-active {
          .hamburger-line {
            &:nth-child(1) {
              transform: rotate(45deg);
            }
            
            &:nth-child(2) {
              opacity: 0;
            }
            
            &:nth-child(3) {
              transform: rotate(-45deg);
            }
          }
        }
      }
      
      .breadcrumb-section {
        flex: 1;
        min-width: 0;
        
        .page-title {
          font-size: var(--text-xl);
          font-weight: var(--font-semibold);
          color: var(--text-primary);
          margin-bottom: var(--spacing-1);
        }
        
        .breadcrumb-nav {
          .breadcrumb {
            display: flex;
            align-items: center;
            gap: var(--spacing-2);
            margin: 0;
            padding: 0;
            list-style: none;
            
            .breadcrumb-item {
              display: flex;
              align-items: center;
              gap: var(--spacing-2);
              font-size: var(--text-sm);
              
              .breadcrumb-link {
                color: var(--text-tertiary);
                text-decoration: none;
                transition: color var(--transition-fast);
                
                &:hover {
                  color: var(--primary-500);
                }
              }
              
              .breadcrumb-current {
                color: var(--text-secondary);
                font-weight: var(--font-medium);
              }
              
              .breadcrumb-separator {
                color: var(--text-tertiary);
                font-size: 12px;
              }
            }
          }
        }
      }
      
      .search-section {
        flex: 0 0 400px;
        
        .search-container {
          position: relative;
          
          .global-search {
            .el-input__wrapper {
              border-radius: var(--radius-full);
              height: 40px;
              background: var(--bg-secondary);
              border: 1px solid transparent;
              transition: all var(--transition-fast);
              
              &:hover {
                background: var(--bg-hover);
              }
              
              &.is-focus {
                background: var(--bg-primary);
                border-color: var(--primary-500);
                box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
              }
            }
          }
          
          .search-results {
            position: absolute;
            top: calc(100% + 8px);
            left: 0;
            right: 0;
            background: var(--bg-primary);
            border: 1px solid var(--border-light);
            border-radius: var(--radius-xl);
            box-shadow: var(--shadow-lg);
            z-index: var(--z-dropdown);
            max-height: 300px;
            overflow-y: auto;
            
            .search-item {
              display: flex;
              align-items: center;
              gap: var(--spacing-3);
              padding: var(--spacing-3) var(--spacing-4);
              cursor: pointer;
              transition: background var(--transition-fast);
              
              &:hover {
                background: var(--bg-hover);
              }
              
              .search-item-icon {
                width: 32px;
                height: 32px;
                background: var(--bg-secondary);
                border-radius: var(--radius-md);
                display: flex;
                align-items: center;
                justify-content: center;
                color: var(--primary-500);
              }
              
              .search-item-content {
                flex: 1;
                
                .search-item-title {
                  font-weight: var(--font-medium);
                  color: var(--text-primary);
                  margin-bottom: 2px;
                }
                
                .search-item-type {
                  font-size: var(--text-xs);
                  color: var(--text-tertiary);
                }
              }
            }
          }
        }
      }
      
      .header-actions {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        
        .action-btn {
          width: 40px;
          height: 40px;
          border: none;
          background: none;
          color: var(--text-secondary);
          border-radius: var(--radius-lg);
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          transition: all var(--transition-fast);
          position: relative;
          
          &:hover {
            background: var(--bg-hover);
            color: var(--text-primary);
          }
          
          .el-icon {
            font-size: 20px;
          }
        }
        
        .notification-btn {
          .notification-badge {
            position: absolute;
            top: 8px;
            right: 8px;
          }
        }
        
        .user-dropdown {
          display: flex;
          align-items: center;
          gap: var(--spacing-3);
          padding: var(--spacing-2) var(--spacing-3);
          border-radius: var(--radius-lg);
          cursor: pointer;
          transition: all var(--transition-fast);
          
          &:hover {
            background: var(--bg-hover);
          }
          
          .user-avatar {
            width: 40px;
            height: 40px;
            background: var(--primary-500);
            color: white;
            border-radius: var(--radius-lg);
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: var(--font-semibold);
          }
          
          .user-info {
            .user-name {
              font-weight: var(--font-medium);
              color: var(--text-primary);
              font-size: var(--text-sm);
              margin-bottom: 2px;
            }
            
            .user-role {
              font-size: var(--text-xs);
              color: var(--text-tertiary);
            }
          }
          
          .user-arrow {
            color: var(--text-tertiary);
            font-size: 16px;
          }
        }
      }
    }
    
    // 页面内容
    .app-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
      

      
      .route-container {
        flex: 1;
        overflow: auto;
        padding: var(--spacing-6);
      }
    }
  }
  
  // 快捷操作浮动按钮
  .quick-actions-fab {
    position: fixed;
    bottom: var(--spacing-6);
    right: var(--spacing-6);
    z-index: var(--z-popover);
    
    .quick-actions-menu {
      position: absolute;
      bottom: 100%;
      right: 0;
      margin-bottom: var(--spacing-4);
      display: flex;
      flex-direction: column;
      gap: var(--spacing-2);
      
      .quick-action-item {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        padding: var(--spacing-3) var(--spacing-4);
        background: var(--bg-primary);
        border: 1px solid var(--border-light);
        border-radius: var(--radius-full);
        box-shadow: var(--shadow-md);
        cursor: pointer;
        transition: all var(--transition-fast);
        white-space: nowrap;
        
        &:hover {
          background: var(--bg-hover);
          transform: translateX(-4px);
        }
        
        .el-icon {
          font-size: 16px;
          color: var(--primary-500);
        }
        
        span {
          font-size: var(--text-sm);
          font-weight: var(--font-medium);
          color: var(--text-primary);
        }
      }
    }
  }
}

// 通知面板样式
:deep(.notification-popover) {
  padding: 0 !important;
  border-radius: var(--radius-xl) !important;
  box-shadow: var(--shadow-xl) !important;
}

.notification-panel {
  .notification-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--spacing-4) var(--spacing-6);
    border-bottom: 1px solid var(--border-light);
    
    h3 {
      font-size: var(--text-lg);
      font-weight: var(--font-semibold);
      color: var(--text-primary);
      margin: 0;
    }
  }
  
  .notification-tabs {
    display: flex;
    padding: 0 var(--spacing-6);
    border-bottom: 1px solid var(--border-light);
    
    .notification-tab {
      padding: var(--spacing-3) var(--spacing-4);
      border: none;
      background: none;
      color: var(--text-secondary);
      font-size: var(--text-sm);
      font-weight: var(--font-medium);
      cursor: pointer;
      transition: all var(--transition-fast);
      border-bottom: 2px solid transparent;
      
      &:hover {
        color: var(--text-primary);
      }
      
      &.is-active {
        color: var(--primary-500);
        border-bottom-color: var(--primary-500);
      }
    }
  }
  
  .notification-list {
    max-height: 400px;
    overflow-y: auto;
    
    .notification-item {
      display: flex;
      align-items: flex-start;
      gap: var(--spacing-3);
      padding: var(--spacing-4) var(--spacing-6);
      border-bottom: 1px solid var(--border-light);
      cursor: pointer;
      transition: all var(--transition-fast);
      position: relative;
      
      &:hover {
        background: var(--bg-hover);
      }
      
      &.is-read {
        opacity: 0.7;
      }
      
      .notification-icon {
        width: 40px;
        height: 40px;
        border-radius: var(--radius-lg);
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        
        &.success {
          background: var(--success-100);
          color: var(--success-600);
        }
        
        &.warning {
          background: var(--warning-100);
          color: var(--warning-600);
        }
        
        &.info {
          background: var(--gray-100);
          color: var(--gray-600);
        }
        
        &.error {
          background: var(--error-100);
          color: var(--error-600);
        }
      }
      
      .notification-content {
        flex: 1;
        min-width: 0;
        
        .notification-title {
          font-weight: var(--font-medium);
          color: var(--text-primary);
          margin-bottom: var(--spacing-1);
          font-size: var(--text-sm);
        }
        
        .notification-description {
          color: var(--text-secondary);
          font-size: var(--text-sm);
          line-height: 1.4;
          margin-bottom: var(--spacing-2);
        }
        
        .notification-time {
          color: var(--text-tertiary);
          font-size: var(--text-xs);
        }
      }
      
      .notification-close {
        width: 24px;
        height: 24px;
        border: none;
        background: none;
        color: var(--text-tertiary);
        border-radius: var(--radius-md);
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all var(--transition-fast);
        opacity: 0;
        
        &:hover {
          background: var(--bg-hover);
          color: var(--text-primary);
        }
      }
      
      &:hover .notification-close {
        opacity: 1;
      }
    }
    
    .notification-empty {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: var(--spacing-8) var(--spacing-6);
      color: var(--text-tertiary);
      
      .el-icon {
        font-size: 48px;
        margin-bottom: var(--spacing-4);
        opacity: 0.5;
      }
      
      p {
        margin: 0;
        font-size: var(--text-sm);
      }
    }
  }
}

// 页面切换动画
.page-transition-enter-active {
  transition: all var(--transition-base);
}

.page-transition-leave-active {
  transition: all var(--transition-fast);
}

.page-transition-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-transition-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

// 通知动画
.notification-enter-active,
.notification-leave-active {
  transition: all var(--transition-base);
}

.notification-enter-from,
.notification-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

.notification-move {
  transition: transform var(--transition-base);
}

// 快捷操作动画
.quick-actions-enter-active,
.quick-actions-leave-active {
  transition: all var(--transition-base);
}

.quick-actions-enter-from,
.quick-actions-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

// 响应式设计
@media (max-width: 768px) {
  .app-layout {
    .app-sidebar {
      position: fixed;
      top: 0;
      left: -280px;
      height: 100vh;
      z-index: var(--z-modal);
      transition: left var(--transition-base);
      
      &.is-mobile-open {
        left: 0;
      }
    }
    
    .app-main {
      .app-header {
        padding: 0 var(--spacing-4);
        
        .breadcrumb-section {
          .page-title {
            font-size: var(--text-lg);
          }
        }
        
        .search-section {
          display: none;
        }
        
        .header-actions {
          .user-info {
            display: none;
          }
          
          .user-arrow {
            display: none;
          }
        }
      }
      
      .app-content {
        .route-container {
          padding: var(--spacing-4);
        }
      }
    }
  }
  
  .desktop-only {
    display: none !important;
  }
}

@media (min-width: 769px) {
  .mobile-only {
    display: none !important;
  }
  
  .app-layout.sidebar-collapsed {
    .app-sidebar {
      width: 80px;
    }
  }
}

// 自定义滚动条
.custom-scrollbar {
  &::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  
  &::-webkit-scrollbar-thumb {
    background: var(--border-default);
    border-radius: var(--radius-full);
    
    &:hover {
      background: var(--border-strong);
    }
  }
}

// 暗色主题适配
.dark-theme {
  // 在这里可以添加特定的暗色主题样式覆盖
}
</style>