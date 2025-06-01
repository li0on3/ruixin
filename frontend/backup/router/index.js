import { createRouter, createWebHistory } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { useUserStore } from '@/store/user'

NProgress.configure({ showSpinner: false })

const routes = [
  {
    path: '/login',
    name: 'Login',
    // component: () => import('@/views/Login.vue'),
    component: () => import('@/views/LoginEnhanced.vue'), // 使用增强版登录页面
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/',
    // component: () => import('@/views/Layout.vue'),
    component: () => import('@/views/LayoutEnhanced.vue'), // 使用增强版布局
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        // component: () => import('@/views/Dashboard.vue'),
        component: () => import('@/views/DashboardEnhanced.vue'), // 使用增强版仪表盘
        meta: { title: '仪表盘' }
      },
      {
        path: 'cards',
        name: 'Cards',
        // component: () => import('@/views/cards/Index.vue'),
        component: () => import('@/views/cards/CardsEnhanced.vue'), // 使用增强版卡片管理
        meta: { title: '卡片管理' }
      },
      {
        path: 'cards/bindings',
        name: 'CardBindings',
        component: () => import('@/views/cards/Bindings.vue'),
        meta: { title: '卡片绑定管理' }
      },
      {
        path: 'cards/batch-import',
        name: 'CardBatchImport',
        component: () => import('@/views/cards/BatchImport.vue'),
        meta: { title: '批量导入卡片' }
      },
      {
        path: 'cards/batches',
        name: 'CardBatches',
        component: () => import('@/views/cards/Batches.vue'),
        meta: { title: '批次管理' }
      },
      {
        path: 'distributors',
        name: 'Distributors',
        component: () => import('@/views/distributors/Index.vue'),
        meta: { title: '分销商管理' }
      },
      {
        path: 'orders',
        name: 'Orders',
        // component: () => import('@/views/orders/Index.vue'),
        component: () => import('@/views/orders/EnhancedIndex.vue'), // 使用增强版订单管理
        meta: { title: '订单管理' }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        // component: () => import('@/views/statistics/Index.vue'),
        component: () => import('@/views/statistics/StatisticsEnhanced.vue'), // 使用增强版统计页面
        meta: { title: '数据统计' }
      },
      {
        path: 'finance/transactions',
        name: 'FinanceTransactions',
        component: () => import('@/views/finance/Transactions.vue'),
        meta: { title: '交易记录' }
      },
      {
        path: 'finance/withdrawals',
        name: 'FinanceWithdrawals',
        component: () => import('@/views/finance/Withdrawals.vue'),
        meta: { title: '提现管理' }
      },
      {
        path: 'luckin/prices',
        name: 'LuckinPrices',
        component: () => import('@/views/luckin/Prices.vue'),
        meta: { title: '价格管理' }
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('@/views/products/Index.vue'),
        meta: { title: '商品管理' }
      },
      {
        path: 'products/available',
        name: 'AvailableProducts',
        component: () => import('@/views/products/AvailableProducts.vue'),
        meta: { title: '可用商品查询' }
      },
      {
        path: 'admins',
        name: 'Admins',
        component: () => import('@/views/admins/Index.vue'),
        meta: { title: '管理员管理' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/Index.vue'),
        meta: { title: '系统设置' }
      },
      {
        path: 'stores',
        name: 'Stores',
        component: () => import('@/views/stores/Index.vue'),
        meta: { title: '店铺查询' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404.vue'),
    meta: { title: '404' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  NProgress.start()
  
  // 设置页面标题
  document.title = `${to.meta.title || '瑞幸分销系统'} - 瑞幸咖啡分销商自动化系统`
  
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.token) {
    next('/login')
  } else if (to.path === '/login' && userStore.token) {
    next('/')
  } else {
    next()
  }
})

router.afterEach(() => {
  NProgress.done()
})

export default router