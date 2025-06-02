import { createRouter, createWebHistory } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { useUserStore } from '@/store/user'
import { 
  createLazyComponent, 
  preloadRouteComponents, 
  RoutePerformanceMonitor,
  RoutePrefetcher 
} from '@/utils/routeUtils'

NProgress.configure({ showSpinner: false })

// 创建性能监控器
const performanceMonitor = new RoutePerformanceMonitor()
const routePrefetcher = new RoutePrefetcher()

// 创建加载组件
const LoadingComponent = {
  template: `
    <div style="display: flex; justify-content: center; align-items: center; height: 200px;">
      <div style="text-align: center;">
        <div style="width: 40px; height: 40px; border: 3px solid #f3f3f3; border-top: 3px solid #3498db; border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 16px;"></div>
        <p style="color: #666; font-size: 14px;">页面加载中...</p>
      </div>
      <style>
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      </style>
    </div>
  `
}

// 创建错误组件
const ErrorComponent = {
  template: `
    <div style="display: flex; justify-content: center; align-items: center; height: 200px;">
      <div style="text-align: center; color: #f56565;">
        <div style="font-size: 48px; margin-bottom: 16px;">⚠️</div>
        <p style="font-size: 16px; margin-bottom: 8px;">页面加载失败</p>
        <p style="font-size: 14px; color: #666;">请刷新页面重试</p>
      </div>
    </div>
  `
}

// 懒加载配置选项
const lazyLoadOptions = {
  loading: LoadingComponent,
  error: ErrorComponent,
  delay: 200,
  timeout: 10000
}

// 路由配置
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: createLazyComponent(
      () => import('@/views/LoginEnhanced.vue'),
      lazyLoadOptions
    ),
    meta: { 
      title: '登录', 
      requiresAuth: false,
      preload: true // 标记为需要预加载的关键路由
    }
  },
  {
    path: '/',
    component: createLazyComponent(
      () => import('@/views/LayoutPremium.vue'),
      lazyLoadOptions
    ),
    redirect: '/dashboard',
    meta: { 
      requiresAuth: true,
      preload: true
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: createLazyComponent(
          () => import('@/views/DashboardPremium.vue'),
          lazyLoadOptions
        ),
        meta: { 
          title: '仪表盘',
          preload: true,
          keepAlive: true // 保持组件活跃状态
        }
      },
      {
        path: 'cards',
        name: 'Cards',
        component: createLazyComponent(
          () => import('@/views/cards/CardsEnhanced.vue'),
          lazyLoadOptions
        ),
        meta: { 
          title: '卡片管理',
          preload: true,
          keepAlive: true
        }
      },
      {
        path: 'cards/bindings',
        name: 'CardBindings',
        component: createLazyComponent(
          () => import('@/views/cards/Bindings.vue'),
          lazyLoadOptions
        ),
        meta: { title: '卡片绑定管理' }
      },
      {
        path: 'cards/batch-import',
        name: 'CardBatchImport',
        component: createLazyComponent(
          () => import('@/views/cards/BatchImport.vue'),
          lazyLoadOptions
        ),
        meta: { title: '批量导入卡片' }
      },
      {
        path: 'cards/batches',
        name: 'CardBatches',
        component: createLazyComponent(
          () => import('@/views/cards/Batches.vue'),
          lazyLoadOptions
        ),
        meta: { title: '批次管理' }
      },
      {
        path: 'distributors',
        name: 'Distributors',
        component: createLazyComponent(
          () => import('@/views/distributors/Index.vue'),
          lazyLoadOptions
        ),
        meta: { 
          title: '分销商管理',
          preload: true,
          keepAlive: true
        }
      },
      {
        path: 'orders',
        name: 'Orders',
        component: createLazyComponent(
          () => import('@/views/orders/EnhancedIndex.vue'),
          lazyLoadOptions
        ),
        meta: { 
          title: '订单管理',
          preload: true,
          keepAlive: true
        }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: createLazyComponent(
          () => import('@/views/statistics/StatisticsEnhanced.vue'),
          lazyLoadOptions
        ),
        meta: { 
          title: '数据统计',
          preload: true
        }
      },
      {
        path: 'finance/transactions',
        name: 'FinanceTransactions',
        component: createLazyComponent(
          () => import('@/views/finance/Transactions.vue'),
          lazyLoadOptions
        ),
        meta: { title: '交易记录' }
      },
      {
        path: 'finance/withdrawals',
        name: 'FinanceWithdrawals',
        component: createLazyComponent(
          () => import('@/views/finance/Withdrawals.vue'),
          lazyLoadOptions
        ),
        meta: { title: '提现管理' }
      },
      {
        path: 'luckin/prices',
        name: 'LuckinPrices',
        component: createLazyComponent(
          () => import('@/views/luckin/Prices.vue'),
          lazyLoadOptions
        ),
        meta: { title: '价格管理' }
      },
      {
        path: 'products',
        name: 'Products',
        component: createLazyComponent(
          () => import('@/views/products/Index.vue'),
          lazyLoadOptions
        ),
        meta: { title: '商品管理' }
      },
      {
        path: 'products/available',
        name: 'AvailableProducts',
        component: createLazyComponent(
          () => import('@/views/products/AvailableProducts.vue'),
          lazyLoadOptions
        ),
        meta: { title: '可用商品查询' }
      },
      {
        path: 'admins',
        name: 'Admins',
        component: createLazyComponent(
          () => import('@/views/admins/Index.vue'),
          lazyLoadOptions
        ),
        meta: { title: '管理员管理' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: createLazyComponent(
          () => import('@/views/settings/Index.vue'),
          lazyLoadOptions
        ),
        meta: { title: '系统设置' }
      },
      {
        path: 'stores',
        name: 'Stores',
        component: createLazyComponent(
          () => import('@/views/stores/Index.vue'),
          lazyLoadOptions
        ),
        meta: { title: '店铺查询' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: createLazyComponent(
      () => import('@/views/404.vue'),
      lazyLoadOptions
    ),
    meta: { title: '404' }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  // 路由滚动行为
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.hash) {
      return { el: to.hash }
    } else {
      return { top: 0, behavior: 'smooth' }
    }
  }
})

// 预加载关键路由
const { preloadCritical, preloadAll } = preloadRouteComponents(routes)

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 开始性能监控
  performanceMonitor.startRouteChange(from, to)
  
  // 开始进度条
  NProgress.start()
  
  // 设置页面标题
  document.title = `${to.meta.title || '瑞幸分销系统'} - 瑞幸咖啡分销商自动化系统`
  
  // 权限检查
  const userStore = useUserStore()
  
  // 需要认证的路由但没有 token
  if (to.meta.requiresAuth && !userStore.token) {
    // 如果已经在跳转到登录页的过程中，直接返回
    if (from.path === to.path) {
      next(false)
      return
    }
    
    // 保存尝试访问的路径
    sessionStorage.setItem('redirectPath', to.fullPath)
    next({
      path: '/login',
      replace: true
    })
    return
  }
  
  // 已登录用户访问登录页
  if (to.path === '/login' && userStore.token) {
    // 检查是否有保存的重定向路径
    const redirectPath = sessionStorage.getItem('redirectPath')
    if (redirectPath) {
      sessionStorage.removeItem('redirectPath')
      next({
        path: redirectPath,
        replace: true
      })
    } else {
      next({
        path: '/',
        replace: true
      })
    }
    return
  }
  
  // 预取相关路由
  if (from.path) {
    routePrefetcher.predictAndPrefetch(to.path, routes)
  }
  
  next()
})

router.afterEach((to, from) => {
  // 结束性能监控
  performanceMonitor.endRouteChange(to)
  
  // 结束进度条
  NProgress.done()
  
  // 记录路由访问（用于预取优化）
  if (typeof gtag !== 'undefined') {
    gtag('config', 'GA_MEASUREMENT_ID', {
      page_title: to.meta.title,
      page_location: window.location.href
    })
  }
})

// 路由错误处理
router.onError((error) => {
  console.error('Router error:', error)
  NProgress.done()
  
  // 可以在这里添加错误上报
  if (typeof Sentry !== 'undefined') {
    Sentry.captureException(error)
  }
})

// 初始化性能优化
let isInitialized = false

// 创建预加载助手
const preloadHelper = preloadRouteComponents(routes)

router.isReady().then(() => {
  if (!isInitialized) {
    isInitialized = true
    
    // 预加载关键路由
    preloadHelper.preloadCritical().then(() => {
      console.log('Critical routes preloaded')
    })
    
    // 在空闲时间预加载所有路由
    preloadHelper.preloadAll().then(() => {
      console.log('All routes preloaded')
    })
    
    // 开发环境下提供性能调试工具
    if (process.env.NODE_ENV === 'development') {
      window.__ROUTE_PERFORMANCE__ = {
        getReport: () => performanceMonitor.getReport(),
        getPrefetchStats: () => routePrefetcher.getStats(),
        clearMetrics: () => {
          performanceMonitor.clear()
          console.log('Performance metrics cleared')
        }
      }
      
      console.log('Route performance tools available at window.__ROUTE_PERFORMANCE__')
    }
  }
})

export default router

// 导出性能工具（用于测试和调试）
export { performanceMonitor, routePrefetcher }