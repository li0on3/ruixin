<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="sidebarWidth" class="aside modern-menu">
      <div class="logo">
        <img src="@/assets/logo.png" alt="Logo" v-if="!isCollapse">
        <span v-else>瑞</span>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        class="menu"
        :collapse="isCollapse"
        :collapse-transition="false"
        router
      >
        <el-tooltip :content="isCollapse ? '仪表盘' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/dashboard">
            <el-icon><DataLine /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-tooltip :content="isCollapse ? '卡片管理' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/cards">
            <el-icon><CreditCard /></el-icon>
            <span>卡片管理</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-tooltip :content="isCollapse ? '分销商管理' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/distributors">
            <el-icon><UserFilled /></el-icon>
            <span>分销商管理</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-tooltip :content="isCollapse ? '订单管理' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/orders">
            <el-icon><ShoppingCart /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-tooltip :content="isCollapse ? '数据统计' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/statistics">
            <el-icon><TrendCharts /></el-icon>
            <span>数据统计</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-sub-menu index="/finance">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>财务管理</span>
          </template>
          <el-menu-item index="/finance/transactions">交易记录</el-menu-item>
          <el-menu-item index="/finance/withdrawals">提现管理</el-menu-item>
        </el-sub-menu>
        
        <el-tooltip :content="isCollapse ? '价格管理' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/luckin/prices">
            <el-icon><Coffee /></el-icon>
            <span>价格管理</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-sub-menu index="/products">
          <template #title>
            <el-icon><ShoppingBag /></el-icon>
            <span>商品管理</span>
          </template>
          <el-menu-item index="/products">商品管理</el-menu-item>
          <el-menu-item index="/products/available">可用商品查询</el-menu-item>
        </el-sub-menu>
        
        <el-tooltip :content="isCollapse ? '管理员管理' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/admins" v-if="userStore.isSuperAdmin">
            <el-icon><User /></el-icon>
            <span>管理员管理</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-tooltip :content="isCollapse ? '店铺查询' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/stores">
            <el-icon><Location /></el-icon>
            <span>店铺查询</span>
          </el-menu-item>
        </el-tooltip>
        
        <el-tooltip :content="isCollapse ? '系统设置' : ''" placement="right" :disabled="!isCollapse">
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-tooltip>
      </el-menu>
    </el-aside>
    
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header modern-header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="isCollapse = !isCollapse">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
          
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute.meta?.title">
              {{ currentRoute.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- 主题切换 -->
          <el-dropdown @command="handleThemeChange" class="theme-switcher">
            <el-button text>
              <el-icon>
                <Sunny v-if="!isDarkTheme" />
                <Moon v-else />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="light">
                  <el-icon><Sunny /></el-icon>
                  浅色主题
                </el-dropdown-item>
                <el-dropdown-item command="dark">
                  <el-icon><Moon /></el-icon>
                  深色主题
                </el-dropdown-item>
                <el-dropdown-item command="auto">
                  <el-icon><Monitor /></el-icon>
                  跟随系统
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <!-- 用户信息 -->
          <el-dropdown>
            <div class="user-info">
              <el-avatar :size="32" :src="avatarUrl">
                {{ userStore.userInfo?.username?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo?.username || '用户' }}</span>
            </div>
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
      </el-header>
      
      <!-- 主内容区 -->
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox } from 'element-plus'
import { Sunny, Moon, Monitor } from '@element-plus/icons-vue'
import { useTheme } from '@/plugins/theme'

const route = useRoute()
const userStore = useUserStore()

// 主题相关
const { themeState, setTheme, isDarkTheme } = useTheme()

const isCollapse = ref(false)
const currentRoute = computed(() => route)
const activeMenu = computed(() => route.path)
const avatarUrl = computed(() => '')

// 响应式侧边栏宽度
const sidebarWidth = computed(() => {
  if (isCollapse.value) {
    return window.innerWidth < 1366 ? '80px' : '100px'
  }
  return window.innerWidth < 1366 ? '260px' : '320px'
})

const handleProfile = () => {
  // TODO: 实现个人中心功能
}

const handlePassword = () => {
  // TODO: 实现修改密码功能
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await userStore.logout()
  } catch (error) {
    // 用户取消
  }
}

// 主题切换
const handleThemeChange = (command) => {
  setTheme(command)
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
  width: 100%;
  overflow: hidden; // 防止整体滚动
  
  .aside {
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    z-index: 10;
    
    // 响应式侧边栏宽度
    @media (max-width: 1366px) {
      &:not(.el-aside--collapse) {
        width: 260px !important;
      }
      
      &.el-aside--collapse {
        width: 80px !important;
      }
    }
  }
  
  .header {
    background-color: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 32px;
    height: 100px;
    width: 100%;
    
    @media (max-width: 1366px) {
      padding: 0 16px;
      height: 80px;
    }
    
    .header-left {
      display: flex;
      align-items: center;
      flex-shrink: 0;
      
      .collapse-btn {
        font-size: 20px;
        cursor: pointer;
        margin-right: 20px;
        
        &:hover {
          color: #667eea;
        }
      }
      
      .el-breadcrumb {
        font-size: 16px;
        white-space: nowrap;
      }
    }
    
    .header-right {
      display: flex;
      align-items: center;
      gap: 16px;
      flex-shrink: 0;
      
      .theme-switcher {
        .el-button {
          &:hover {
            color: var(--el-color-primary);
          }
        }
      }
      
      .user-info {
        display: flex;
        align-items: center;
        cursor: pointer;
        
        .username {
          margin-left: 10px;
          color: #606266;
          font-size: 16px;
          white-space: nowrap;
        }
      }
    }
    
    // 响应式调整
    @media (min-width: 1440px) and (max-width: 1920px) {
      padding: 0 24px;
    }
    
    @media (min-width: 1920px) {
      padding: 0 48px;
    }
  }
  
  .main {
    background-color: #f5f7fa;
    padding: 16px;
    overflow-y: auto;
    overflow-x: hidden; // 防止内容区域出现横向滚动
    height: calc(100vh - 100px); // 减去header高度
    width: 100%;
    
    @media (max-width: 1366px) {
      height: calc(100vh - 80px); // 小屏幕header高度
      padding: 12px;
    }
    
    @media (min-width: 1440px) and (max-width: 1920px) {
      padding: 20px;
    }
    
    @media (min-width: 1920px) {
      padding: 24px 32px;
    }
    
    // 让内容充满可视区域
    :deep(.router-view) {
      height: 100%;
      min-width: 100%;
    }
  }
  
  // 侧边栏收起时的调整
  &.sidebar-collapsed {
    .main {
      width: 100%;
    }
    
    .header {
      width: 100%;
    }
  }
}

// 页面切换动画
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>