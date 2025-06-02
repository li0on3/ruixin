<template>
  <div class="dashboard-premium">
    <!-- 仪表盘头部 -->
    <div class="dashboard-header">
      <div class="header-content">
        <div class="welcome-section">
          <h1 class="dashboard-title">{{ getGreeting() }}，{{ userStore.userInfo?.username || '管理员' }}！</h1>
          <p class="dashboard-subtitle">
            今天是 {{ formatDate(new Date()) }}，这是您的业务概览
          </p>
        </div>
        
        <div class="header-actions">
          <div class="date-selector">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              @change="handleDateChange"
              :shortcuts="dateShortcuts"
              size="large"
              class="date-picker"
            />
          </div>
          
          <div class="action-buttons">
            <el-button 
              type="primary" 
              :icon="Refresh" 
              @click="refreshDashboard"
              :loading="refreshing"
              size="large"
            >
              刷新数据
            </el-button>
            
            <el-dropdown trigger="click" placement="bottom-end">
              <el-button 
                :icon="Download" 
                size="large"
                class="export-btn"
              >
                导出报表
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="exportReport('pdf')">
                    <el-icon><Document /></el-icon>
                    导出 PDF 报表
                  </el-dropdown-item>
                  <el-dropdown-item @click="exportReport('excel')">
                    <el-icon><Document /></el-icon>
                    导出 Excel 报表
                  </el-dropdown-item>
                  <el-dropdown-item @click="exportReport('image')">
                    <el-icon><Picture /></el-icon>
                    导出图片
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </div>

    <!-- 核心指标区域 -->
    <div class="metrics-section">
      <div class="metrics-grid">
        <div 
          class="metric-card"
          v-for="(metric, index) in coreMetrics"
          :key="metric.key"
          :class="metric.type"
          @click="handleMetricClick(metric)"
        >
          <div class="metric-header">
            <div class="metric-icon" :class="metric.iconClass">
              <el-icon><component :is="metric.icon" /></el-icon>
            </div>
            <div class="metric-trend" :class="{ 'positive': metric.trend > 0, 'negative': metric.trend < 0 }">
              <el-icon>
                <TrendCharts v-if="metric.trend > 0" />
                <Bottom v-else />
              </el-icon>
              <span>{{ Math.abs(metric.trend) }}%</span>
            </div>
          </div>
          
          <div class="metric-content">
            <div class="metric-value">
              <AnimatedNumber 
                :value="metric.value"
                :duration="1500"
                :prefix="metric.prefix"
                :suffix="metric.suffix"
                :decimals="metric.decimals || 0"
              />
            </div>
            <div class="metric-title">{{ metric.title }}</div>
            <div class="metric-subtitle">{{ metric.subtitle }}</div>
          </div>
          
          <div class="metric-chart">
            <div 
              :ref="el => setChartRef(metric.key, el)" 
              class="mini-chart"
              :data-chart="metric.key"
            ></div>
          </div>
          
          <div class="metric-overlay">
            <el-icon><View /></el-icon>
            <span>查看详情</span>
          </div>
        </div>
      </div>
    </div>


    <!-- 图表分析区域 -->
    <div class="charts-section">
      <div class="charts-grid">
        <!-- 销售趋势图 -->
        <div class="chart-card main-chart">
          <div class="chart-header">
            <div class="chart-title-section">
              <h3 class="chart-title">
                <el-icon><TrendCharts /></el-icon>
                销售趋势分析
              </h3>
              <p class="chart-subtitle">过去30天的销售数据变化趋势</p>
            </div>
            
            <div class="chart-controls">
              <el-radio-group 
                v-model="trendChartType" 
                size="small" 
                @change="updateTrendChart"
                class="chart-type-selector"
              >
                <el-radio-button value="revenue">销售额</el-radio-button>
                <el-radio-button value="orders">订单量</el-radio-button>
                <el-radio-button value="profit">利润</el-radio-button>
              </el-radio-group>
              
              <el-radio-group 
                v-model="trendPeriod" 
                size="small" 
                @change="updateTrendChart"
                class="period-selector"
              >
                <el-radio-button value="7">7天</el-radio-button>
                <el-radio-button value="30">30天</el-radio-button>
                <el-radio-button value="90">90天</el-radio-button>
              </el-radio-group>
              
              <el-button 
                :icon="FullScreen" 
                size="small" 
                circle 
                @click="expandChart('trend')"
                title="全屏查看"
              />
            </div>
          </div>
          
          <div class="chart-body">
            <div 
              ref="trendChartRef" 
              class="chart-container"
              v-loading="chartsLoading.trend"
            ></div>
          </div>
          
          <div class="chart-footer">
            <div class="chart-stats">
              <div class="stat-item">
                <span class="stat-label">最高值</span>
                <span class="stat-value">{{ formatTrendStatValue(trendStats.max) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">平均值</span>
                <span class="stat-value">{{ formatTrendStatValue(trendStats.avg) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">增长率</span>
                <span class="stat-value" :class="{ 'positive': trendStats.growth > 0 }">
                  {{ trendStats.growth > 0 ? '+' : '' }}{{ trendStats.growth }}%
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 商品分布图 -->
        <div class="chart-card">
          <div class="chart-header">
            <div class="chart-title-section">
              <h3 class="chart-title">
                <el-icon><Coffee /></el-icon>
                热销商品分布
              </h3>
              <p class="chart-subtitle">Top 10 畅销商品占比</p>
            </div>
            
            <div class="chart-controls">
              <el-button 
                :icon="FullScreen" 
                size="small" 
                circle 
                @click="expandChart('product')"
                title="全屏查看"
              />
            </div>
          </div>
          
          <div class="chart-body">
            <div 
              ref="productChartRef" 
              class="chart-container"
              v-loading="chartsLoading.product"
            ></div>
          </div>
          
          <div class="chart-footer">
            <div class="product-legend">
              <div 
                class="legend-item"
                v-for="(product, index) in topProducts.slice(0, 5)"
                :key="product.name"
              >
                <span class="legend-color" :style="{ backgroundColor: productColors[index] }"></span>
                <span class="legend-name">{{ product.name }}</span>
                <span class="legend-value">{{ product.percentage }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 地域分布图已移除 -->

        <!-- 时段分析图 -->
        <div class="chart-card">
          <div class="chart-header">
            <div class="chart-title-section">
              <h3 class="chart-title">
                <el-icon><Clock /></el-icon>
                时段分析
              </h3>
              <p class="chart-subtitle">24小时订单分布</p>
            </div>
            
            <div class="chart-controls">
              <el-button 
                :icon="FullScreen" 
                size="small" 
                circle 
                @click="expandChart('time')"
                title="全屏查看"
              />
            </div>
          </div>
          
          <div class="chart-body">
            <div 
              ref="timeChartRef" 
              class="chart-container"
              v-loading="chartsLoading.time"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近活动区域 -->
    <div class="activities-section">
      <div class="section-header">
        <h2 class="section-title">
          <el-icon><Bell /></el-icon>
          最近活动
        </h2>
        <el-button text @click="viewAllActivities">
          查看全部
          <el-icon class="el-icon--right"><ArrowRight /></el-icon>
        </el-button>
      </div>
      
      <div class="activities-grid">
        <!-- 最新订单 -->
        <div class="activity-card">
          <div class="activity-header">
            <h4 class="activity-title">
              <el-icon><ShoppingCart /></el-icon>
              最新订单
            </h4>
            <el-link type="primary" @click="viewAllOrders">查看全部</el-link>
          </div>
          
          <div class="activity-list">
            <div 
              class="activity-item"
              v-for="order in recentOrders"
              :key="order.orderNo"
              @click="viewOrderDetail(order)"
            >
              <div class="activity-avatar">
                <el-icon><User /></el-icon>
              </div>
              <div class="activity-content">
                <div class="activity-main">
                  <span class="activity-user">{{ order.distributorName || '未知分销商' }}</span>
                  <span class="activity-action">下单</span>
                  <span class="activity-target">{{ order.storeName || '未知门店' }}</span>
                </div>
                <div class="activity-meta">
                  <span class="activity-amount">¥{{ order.totalAmount || 0 }}</span>
                  <span class="activity-time">{{ formatRelativeTime(order.createdAt) }}</span>
                </div>
              </div>
              <div class="activity-status" :class="order.status">
                <el-tag :type="getOrderStatusType(order.status)" size="small">
                  {{ getOrderStatusText(order.status) }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>

        <!-- 系统通知已移除 -->
      </div>
    </div>

    <!-- 全屏图表弹窗 -->
    <el-dialog
      v-model="chartDialogVisible"
      :title="chartDialogTitle"
      width="90%"
      :append-to-body="true"
      :destroy-on-close="true"
      class="chart-dialog"
    >
      <div ref="expandedChartRef" class="expanded-chart"></div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage, ElLoading } from 'element-plus'
import * as echarts from 'echarts'
import {
  TrendCharts, Coffee, ShoppingCart, User, Location, Clock, Lightning,
  Download, Document, Picture, Refresh, ArrowDown, ArrowRight, FullScreen,
  Bell, View, Bottom, SuccessFilled, WarningFilled, InfoFilled, CircleCloseFilled,
  CreditCard, Money
} from '@element-plus/icons-vue'
import AnimatedNumber from '@/components/AnimatedNumber.vue'
import { 
  getDashboardStatistics, 
  getOrderTrend, 
  getHotGoods, 
  getRecentOrders,
  getDistributorRank
} from '@/api/dashboard'
import {
  getSalesTrend,
  getRegionDistribution,
  getHourDistribution
} from '@/api/statistics'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const dateRange = ref([])
const refreshing = ref(false)
const chartsLoading = reactive({
  trend: false,
  product: false,
  // region: false, // 移除地域分布图表
  time: false
})

const trendChartType = ref('revenue')
const trendPeriod = ref('30')
const chartDialogVisible = ref(false)
const chartDialogTitle = ref('')

// 图表引用
const chartRefs = reactive({})
const trendChartRef = ref()
const productChartRef = ref()
// const regionChartRef = ref() // 移除地域分布图表引用
const timeChartRef = ref()
const expandedChartRef = ref()

// 图表实例
const chartInstances = reactive({})

// resize 监听器引用
let resizeHandler = null

// 安全调用图表 resize 的辅助函数
const safeResizeChart = (chart) => {
  try {
    if (chart && typeof chart.resize === 'function' && !chart.isDisposed()) {
      chart.resize()
    }
  } catch (error) {
    console.warn('Chart resize error:', error)
  }
}

// 统一的 resize 处理函数
const handleResize = () => {
  Object.values(chartInstances).forEach(chart => {
    safeResizeChart(chart)
  })
}

// 数据
const coreMetrics = ref([
  {
    key: 'totalRevenue',
    title: '总销售额',
    subtitle: '本月累计',
    value: 0,
    prefix: '¥',
    trend: 12.5,
    icon: 'Money',
    iconClass: 'revenue',
    type: 'primary'
  },
  {
    key: 'totalOrders',
    title: '总订单数',
    subtitle: '本月累计',
    value: 0,
    trend: -3.2,
    icon: 'ShoppingCart',
    iconClass: 'orders',
    type: 'success'
  },
  {
    key: 'activeUsers',
    title: '活跃用户',
    subtitle: '本月活跃',
    value: 0,
    trend: 8.7,
    icon: 'User',
    iconClass: 'users',
    type: 'info'
  },
  {
    key: 'conversionRate',
    title: '转化率',
    subtitle: '本月平均',
    value: 0,
    suffix: '%',
    decimals: 2,
    trend: 5.3,
    icon: 'TrendCharts',
    iconClass: 'conversion',
    type: 'warning'
  }
])


const recentOrders = ref([])
// const recentNotifications = ref([]) // 移除系统通知
const trendStats = reactive({ max: 0, avg: 0, growth: 0 })
const topProducts = ref([])
const productColors = ['#FF6B35', '#F7931E', '#FFD23F', '#06D6A0', '#118AB2', '#073B4C', '#E63946', '#A8DADC', '#457B9D', '#1D3557']

// 日期快捷选项
const dateShortcuts = [
  {
    text: '最近一周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    }
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    }
  },
  {
    text: '最近三个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    }
  }
]

// 计算属性
const chartTheme = computed(() => ({
  color: [
    '#FF6B35', '#F7931E', '#FFD23F', '#06D6A0', '#118AB2',
    '#073B4C', '#E63946', '#A8DADC', '#457B9D', '#1D3557'
  ],
  backgroundColor: 'transparent',
  textStyle: {
    fontFamily: 'PingFang SC, Microsoft YaHei, sans-serif'
  },
  grid: {
    top: 60,
    left: 60,
    right: 40,
    bottom: 60,
    containLabel: true
  },
  legend: {
    textStyle: {
      color: '#6B7280'
    }
  },
  categoryAxis: {
    axisLine: {
      lineStyle: {
        color: '#E5E7EB'
      }
    },
    axisTick: {
      lineStyle: {
        color: '#E5E7EB'
      }
    },
    axisLabel: {
      color: '#6B7280'
    }
  },
  valueAxis: {
    axisLine: {
      lineStyle: {
        color: '#E5E7EB'
      }
    },
    axisTick: {
      lineStyle: {
        color: '#E5E7EB'
      }
    },
    axisLabel: {
      color: '#6B7280'
    },
    splitLine: {
      lineStyle: {
        color: '#F3F4F6'
      }
    }
  }
}))

// 方法
const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
}

const formatDate = (date) => {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  }).format(date)
}

const formatNumber = (num) => {
  return new Intl.NumberFormat('zh-CN').format(num)
}

const formatTrendStatValue = (value) => {
  if (trendChartType.value === 'orders') {
    return Math.round(value) + ' 单'
  } else {
    return '¥' + formatNumber(value.toFixed(2))
  }
}

const formatRelativeTime = (date) => {
  const now = new Date()
  const diffMs = now - new Date(date)
  const diffMins = Math.floor(diffMs / 60000)
  
  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours}小时前`
  
  const diffDays = Math.floor(diffHours / 24)
  if (diffDays < 7) return `${diffDays}天前`
  
  return new Intl.DateTimeFormat('zh-CN').format(new Date(date))
}

const setChartRef = (key, el) => {
  if (el) {
    chartRefs[key] = el
  }
}

const handleDateChange = (dates) => {
  console.log('Date range changed:', dates)
  console.log('Current dateRange value:', dateRange.value)
  
  // 当用户使用顶部日期选择器时，清空销售趋势的时间周期选择
  trendPeriod.value = ''
  
  // 确保日期范围已更新
  if (dates && dates.length === 2) {
    loadDashboardData()
  } else {
    console.warn('Invalid date range:', dates)
  }
}

const handleMetricClick = (metric) => {
  console.log('Metric clicked:', metric)
  // 根据指标类型跳转到相应页面
  switch (metric.key) {
    case 'totalRevenue':
      router.push('/finance/transactions')
      break
    case 'totalOrders':
      router.push('/orders')
      break
    case 'activeUsers':
      router.push('/distributors')
      break
    default:
      router.push('/statistics')
  }
}


const refreshDashboard = async () => {
  refreshing.value = true
  try {
    await loadDashboardData()
    ElMessage.success('数据刷新成功')
  } catch (error) {
    ElMessage.error('数据刷新失败')
  } finally {
    refreshing.value = false
  }
}

const exportReport = (type) => {
  console.log('Export report:', type)
  ElMessage.success(`正在导出${type.toUpperCase()}报表...`)
}

const expandChart = (chartType) => {
  chartDialogTitle.value = getChartTitle(chartType)
  chartDialogVisible.value = true
  
  nextTick(() => {
    if (expandedChartRef.value) {
      const expandedChart = echarts.init(expandedChartRef.value)
      const sourceChart = chartInstances[chartType]
      if (sourceChart) {
        const option = sourceChart.getOption()
        expandedChart.setOption(option)
        
        // 扩展视图的图表会随主 resize 处理器一起处理
      }
    }
  })
}

const getChartTitle = (chartType) => {
  const titles = {
    trend: '销售趋势分析',
    product: '热销商品分布',
    region: '地域分布',
    time: '时段分析'
  }
  return titles[chartType] || '图表'
}

const updateTrendChart = () => {
  initTrendChart()
}

const loadDashboardData = async () => {
  try {
    console.log('Loading dashboard data with date range:', dateRange.value)
    
    // 准备日期参数
    const dateParams = dateRange.value && dateRange.value.length === 2 ? {
      start_date: dateRange.value[0],
      end_date: dateRange.value[1]
    } : {}
    
    // 加载核心指标 - 传递日期参数
    const statsResponse = await getDashboardStatistics(dateParams)
    if (statsResponse.code === 200) {
      updateCoreMetrics(statsResponse.data)
    } else {
      // API 调用失败时显示错误信息，不使用模拟数据
      ElMessage.error('加载统计数据失败')
      // 设置空的核心指标
      coreMetrics.value = []
    }
    
    // 加载最新订单 - 传递日期参数
    const ordersResponse = await getRecentOrders(dateParams)
    if (ordersResponse.code === 200) {
      recentOrders.value = ordersResponse.data
    } else {
      ElMessage.error('加载最新订单失败')
      recentOrders.value = []
    }
    
    // 加载热销商品 - 传递日期参数
    const productsResponse = await getHotGoods(dateParams)
    if (productsResponse.code === 200) {
      topProducts.value = productsResponse.data
    } else {
      ElMessage.error('加载热销商品失败')
      topProducts.value = []
    }
    
    // 初始化图表
    await nextTick()
    initAllCharts()
    
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    ElMessage.error('加载仪表盘数据失败，请检查网络连接后重试')
    // 设置空值，不使用模拟数据
    coreMetrics.value = []
    recentOrders.value = []
    topProducts.value = []
  }
}

const updateCoreMetrics = (data) => {
  console.log('Updating core metrics with data:', data)
  
  // 根据日期范围调整标题
  const dateRangeLabel = data.dateRange?.label || '今日'
  const isToday = dateRangeLabel === '今日'
  
  // 初始化核心指标结构（根据后端实际返回的数据字段）
  coreMetrics.value = [
    {
      key: 'todayAmount',
      title: `${dateRangeLabel}销售额`,
      subtitle: isToday ? '今日累计' : dateRangeLabel,
      value: data.todayAmount || 0,
      prefix: '¥',
      trend: data.amountGrowth || 0,
      icon: 'Money',
      iconClass: 'revenue',
      type: 'primary',
      format: (val) => val.toFixed(2),
      miniChartData: []
    },
    {
      key: 'todayOrders',
      title: `${dateRangeLabel}订单数`,
      subtitle: isToday ? '今日累计' : dateRangeLabel,
      value: data.todayOrders || 0,
      trend: data.orderGrowth || 0,
      icon: 'ShoppingCart',
      iconClass: 'orders',
      type: 'success',
      miniChartData: []
    },
    {
      key: 'activeDistributors',
      title: '活跃分销商',
      subtitle: `总计 ${data.totalDistributors || 0} 个`,
      value: data.activeDistributors || 0,
      trend: 0,
      icon: 'User',
      iconClass: 'users',
      type: 'info',
      miniChartData: []
    },
    {
      key: 'activeCards',
      title: '活跃卡片',
      subtitle: `总计 ${data.totalCards || 0} 张`,
      value: data.activeCards || 0,
      trend: 0,
      icon: 'CreditCard',
      iconClass: 'cards',
      type: 'warning',
      miniChartData: []
    }
  ]
}

// 已删除 loadMockData 函数，现在仅使用真实 API 数据

const initAllCharts = () => {
  initTrendChart()
  initProductChart()
  // initRegionChart() // 移除地域分布图表
  initTimeChart()
  initMiniCharts()
}

const initTrendChart = async () => {
  if (!trendChartRef.value) return
  
  chartsLoading.trend = true
  
  try {
    // 根据 trendPeriod 计算日期范围
    let startDate, endDate
    
    if (trendPeriod.value && trendPeriod.value !== '') {
      // 如果选择了时间周期选择器，使用它的值
      const end = new Date()
      const start = new Date()
      const days = parseInt(trendPeriod.value)
      start.setTime(start.getTime() - 3600 * 1000 * 24 * (days - 1)) // 减1天以包含今天
      
      // 格式化日期为 YYYY-MM-DD 格式
      const formatDate = (date) => {
        const year = date.getFullYear()
        const month = String(date.getMonth() + 1).padStart(2, '0')
        const day = String(date.getDate()).padStart(2, '0')
        return `${year}-${month}-${day}`
      }
      
      startDate = formatDate(start)
      endDate = formatDate(end)
    } else if (dateRange.value && dateRange.value.length === 2) {
      // 否则使用顶部日期选择器的值
      startDate = dateRange.value[0]
      endDate = dateRange.value[1]
    } else {
      // 如果都没有值，使用默认的30天
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 29) // 30天
      
      const formatDate = (date) => {
        const year = date.getFullYear()
        const month = String(date.getMonth() + 1).padStart(2, '0')
        const day = String(date.getDate()).padStart(2, '0')
        return `${year}-${month}-${day}`
      }
      
      startDate = formatDate(start)
      endDate = formatDate(end)
    }
    
    console.log('Trend chart loading with date range:', startDate, 'to', endDate)
    
    // 从真实 API 获取趋势数据
    const response = await getSalesTrend({
      start_date: startDate,
      end_date: endDate
    })
    
    console.log('Sales trend API response:', response)
    
    let trendData = []
    if (response.code === 200 && response.data) {
      trendData = response.data
      console.log('Trend data:', trendData)
    } else {
      console.error('Failed to get trend data:', response)
    }
    
    const chart = echarts.init(trendChartRef.value, chartTheme.value)
    
    // 如果没有数据，显示空状态
    if (!trendData.length) {
      const option = {
        title: {
          text: '暂无数据',
          left: 'center',
          top: 'middle',
          textStyle: {
            color: '#999',
            fontSize: 16
          }
        }
      }
      chart.setOption(option)
      chartInstances.trend = chart
      return
    }
    
    const dates = trendData.map(item => item.date)
    const revenueData = trendData.map(item => item.revenue || 0)
    const ordersData = trendData.map(item => item.orders || 0)
    const avgValueData = trendData.map(item => item.avgValue || 0)
    
    // 计算利润数据（假设利润率为30%）
    const profitData = revenueData.map(revenue => revenue * 0.3)
    
    // 根据选择的图表类型设置数据
    let selectedData = []
    let selectedLabel = ''
    let selectedUnit = ''
    let selectedColor = ''
    
    switch (trendChartType.value) {
      case 'revenue':
        selectedData = revenueData
        selectedLabel = '销售额'
        selectedUnit = '元'
        selectedColor = '#FF6B35'
        break
      case 'orders':
        selectedData = ordersData
        selectedLabel = '订单量'
        selectedUnit = '单'
        selectedColor = '#06D6A0'
        break
      case 'profit':
        selectedData = profitData
        selectedLabel = '利润'
        selectedUnit = '元'
        selectedColor = '#F7931E'
        break
    }
    
    // 计算统计数据
    if (selectedData.length > 0) {
      trendStats.max = Math.max(...selectedData)
      trendStats.avg = selectedData.reduce((sum, val) => sum + val, 0) / selectedData.length
      
      // 计算增长率（最后一天与第一天比较）
      if (selectedData.length > 1 && selectedData[0] > 0) {
        const firstValue = selectedData[0]
        const lastValue = selectedData[selectedData.length - 1]
        trendStats.growth = ((lastValue - firstValue) / firstValue * 100).toFixed(1)
      } else {
        trendStats.growth = 0
      }
    } else {
      trendStats.max = 0
      trendStats.avg = 0
      trendStats.growth = 0
    }
    
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        },
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#E5E7EB',
        textStyle: {
          color: '#374151'
        },
        formatter: function(params) {
          const date = params[0].name
          const value = params[0].value
          return `${date}<br/>${selectedLabel}: ${value.toFixed(2)} ${selectedUnit}`
        }
      },
      grid: {
        top: 60,
        left: 60,
        right: 40,
        bottom: 60,
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: dates,
        boundaryGap: false,
        axisLine: {
          lineStyle: {
            color: '#E5E7EB'
          }
        },
        axisTick: {
          lineStyle: {
            color: '#E5E7EB'
          }
        },
        axisLabel: {
          color: '#6B7280'
        }
      },
      yAxis: {
        type: 'value',
        name: `${selectedLabel}(${selectedUnit})`,
        axisLine: {
          lineStyle: {
            color: '#E5E7EB'
          }
        },
        axisTick: {
          lineStyle: {
            color: '#E5E7EB'
          }
        },
        axisLabel: {
          color: '#6B7280',
          formatter: function(value) {
            if (trendChartType.value === 'orders') {
              return value
            }
            return value.toFixed(0)
          }
        },
        splitLine: {
          lineStyle: {
            color: '#F3F4F6'
          }
        }
      },
      series: [
        {
          name: selectedLabel,
          type: 'line',
          smooth: true,
          data: selectedData,
          symbol: 'circle',
          symbolSize: 6,
          itemStyle: {
            color: selectedColor
          },
          lineStyle: {
            width: 3
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: `${selectedColor}33` }, // 使用十六进制颜色加透明度
                { offset: 1, color: `${selectedColor}0D` }
              ]
            }
          }
        }
      ]
    }
    
    chart.setOption(option)
    chartInstances.trend = chart
    
  } catch (error) {
    console.error('加载趋势图表数据失败:', error)
    // 显示错误状态
    const chart = echarts.init(trendChartRef.value, chartTheme.value)
    const option = {
      title: {
        text: '数据加载失败',
        left: 'center',
        top: 'middle',
        textStyle: {
          color: '#f56565',
          fontSize: 16
        }
      }
    }
    chart.setOption(option)
    chartInstances.trend = chart
  } finally {
    chartsLoading.trend = false
  }
}

const initProductChart = () => {
  if (!productChartRef.value) return
  
  chartsLoading.product = true
  
  setTimeout(() => {
    const chart = echarts.init(productChartRef.value, chartTheme.value)
    
    // 如果没有数据，显示空状态
    if (!topProducts.value.length) {
      const option = {
        title: {
          text: '暂无数据',
          left: 'center',
          top: 'middle',
          textStyle: {
            color: '#999',
            fontSize: 16
          }
        }
      }
      chart.setOption(option)
      chartInstances.product = chart
      chartsLoading.product = false
      return
    }
    
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      series: [
        {
          name: '商品销量',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          label: {
            show: false,
            position: 'center'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '18',
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: topProducts.value.map((product, index) => ({
            name: product.name,
            value: product.percentage,
            itemStyle: {
              color: productColors[index]
            }
          }))
        }
      ]
    }
    
    chart.setOption(option)
    chartInstances.product = chart
    
    // resize 事件由统一的处理器管理
    
    chartsLoading.product = false
  }, 800)
}

const initRegionChart = async () => {
  if (!regionChartRef.value) return
  
  chartsLoading.region = true
  
  try {
    // 从真实 API 获取地区分布数据
    const response = await getRegionDistribution({
      start_date: dateRange.value?.[0],
      end_date: dateRange.value?.[1]
    })
    
    let regionData = []
    if (response.code === 200 && response.data) {
      regionData = response.data
    }
    
    const chart = echarts.init(regionChartRef.value, chartTheme.value)
    
    // 如果没有数据，显示空状态
    if (!regionData.length) {
      const option = {
        title: {
          text: '暂无数据',
          left: 'center',
          top: 'middle',
          textStyle: {
            color: '#999',
            fontSize: 16
          }
        }
      }
      chart.setOption(option)
      chartInstances.region = chart
      return
    }
    
    const cities = regionData.map(item => item.region || item.city)
    const salesData = regionData.map(item => item.revenue || 0)
    
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      xAxis: {
        type: 'category',
        data: cities
      },
      yAxis: {
        type: 'value',
        name: '销售额(元)'
      },
      series: [
        {
          name: '销售额',
          type: 'bar',
          data: salesData,
          itemStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: '#FF6B35' },
                { offset: 1, color: '#F7931E' }
              ]
            }
          },
          emphasis: {
            itemStyle: {
              color: '#E55A2B'
            }
          }
        }
      ]
    }
    
    chart.setOption(option)
    chartInstances.region = chart
    
  } catch (error) {
    console.error('加载地区分布图表数据失败:', error)
    // 显示错误状态
    const chart = echarts.init(regionChartRef.value, chartTheme.value)
    const option = {
      title: {
        text: '数据加载失败',
        left: 'center',
        top: 'middle',
        textStyle: {
          color: '#f56565',
          fontSize: 16
        }
      }
    }
    chart.setOption(option)
    chartInstances.region = chart
  } finally {
    chartsLoading.region = false
  }
}

const initTimeChart = async () => {
  if (!timeChartRef.value) return
  
  chartsLoading.time = true
  
  try {
    // 从真实 API 获取时段分布数据
    const response = await getHourDistribution({
      start_date: dateRange.value?.[0],
      end_date: dateRange.value?.[1]
    })
    
    let timeData = []
    if (response.code === 200 && response.data) {
      timeData = response.data
    }
    
    const chart = echarts.init(timeChartRef.value, chartTheme.value)
    
    // 如果没有数据，显示空状态
    if (!timeData.length) {
      const option = {
        title: {
          text: '暂无数据',
          left: 'center',
          top: 'middle',
          textStyle: {
            color: '#999',
            fontSize: 16
          }
        }
      }
      chart.setOption(option)
      chartInstances.time = chart
      return
    }
    
    const hours = timeData.map(item => `${item.hour}:00`)
    const orderData = timeData.map(item => item.orders || 0)
    
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: hours
      },
      yAxis: {
        type: 'value',
        name: '订单数'
      },
      series: [
        {
          name: '订单数',
          type: 'line',
          smooth: true,
          data: orderData,
          itemStyle: {
            color: '#118AB2'
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(17, 138, 178, 0.3)' },
                { offset: 1, color: 'rgba(17, 138, 178, 0.05)' }
              ]
            }
          }
        }
      ]
    }
    
    chart.setOption(option)
    chartInstances.time = chart
    
  } catch (error) {
    console.error('加载时段分布图表数据失败:', error)
    // 显示错误状态
    const chart = echarts.init(timeChartRef.value, chartTheme.value)
    const option = {
      title: {
        text: '数据加载失败',
        left: 'center',
        top: 'middle',
        textStyle: {
          color: '#f56565',
          fontSize: 16
        }
      }
    }
    chart.setOption(option)
    chartInstances.time = chart
  } finally {
    chartsLoading.time = false
  }
}

const initMiniCharts = () => {
  // 延迟执行以确保 DOM 已经渲染
  nextTick(() => {
    Object.keys(chartRefs).forEach(key => {
      const element = chartRefs[key]
      if (element && element.clientWidth > 0 && element.clientHeight > 0) {
        const chart = echarts.init(element)
      
      // 查找对应的指标数据中的迷你图表数据
      const metric = coreMetrics.value.find(m => m.key === key)
      let data = []
      
      // 如果指标中有 miniChartData，使用真实数据；否则显示空图表
      if (metric && metric.miniChartData && metric.miniChartData.length > 0) {
        data = metric.miniChartData
      } else {
        // 显示空图表（7个0值）
        data = [0, 0, 0, 0, 0, 0, 0]
      }
      
      const option = {
        grid: {
          top: 5,
          left: 5,
          right: 5,
          bottom: 5
        },
        xAxis: {
          type: 'category',
          show: false,
          data: ['', '', '', '', '', '', '']
        },
        yAxis: {
          type: 'value',
          show: false
        },
        series: [
          {
            type: 'line',
            smooth: true,
            symbol: 'none',
            lineStyle: {
              color: '#FF6B35',
              width: 2
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: 'rgba(255, 107, 53, 0.3)' },
                  { offset: 1, color: 'rgba(255, 107, 53, 0.05)' }
                ]
              }
            },
            data
          }
        ]
      }
      
      chart.setOption(option)
      }
    })
  })
}

const getOrderStatusType = (status) => {
  const types = {
    0: 'warning',  // pending
    1: 'info',     // doing
    2: 'success',  // success
    3: 'danger',   // failed
    4: 'warning',  // refunded
    5: 'info'      // cancelled
  }
  return types[status] || 'info'
}

const getOrderStatusText = (status) => {
  const texts = {
    0: '待处理',
    1: '处理中',
    2: '已完成',
    3: '失败',
    4: '已退款',
    5: '已取消'
  }
  return texts[status] || '未知'
}

const viewOrderDetail = (order) => {
  router.push(`/orders?orderNo=${order.orderNo}`)
}

const viewAllOrders = () => {
  router.push('/orders')
}

const viewAllActivities = () => {
  router.push('/activities')
}

// const viewAllNotifications = () => {
//   router.push('/notifications')
// }

// const markAsRead = (notification) => {
//   notification.read = true
// }

// 生命周期
onMounted(() => {
  // 设置默认日期范围为最近30天
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
  
  // 格式化日期为 YYYY-MM-DD 格式
  const formatDate = (date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
  
  dateRange.value = [formatDate(start), formatDate(end)]
  
  loadDashboardData()
  
  // 添加统一的 resize 监听器
  resizeHandler = handleResize
  window.addEventListener('resize', resizeHandler)
})

onBeforeUnmount(() => {
  // 移除 resize 监听器
  if (resizeHandler) {
    window.removeEventListener('resize', resizeHandler)
  }
  
  // 销毁图表实例
  Object.values(chartInstances).forEach(chart => {
    if (chart && typeof chart.dispose === 'function') {
      chart.dispose()
    }
  })
})
</script>

<style lang="scss" scoped>
@import '@/assets/styles/design-tokens.scss';
@import '@/assets/styles/animations-enhanced.scss';

.dashboard-premium {
  padding: var(--spacing-6);
  background: var(--bg-secondary);
  min-height: 100vh;
  
  // 仪表盘头部
  .dashboard-header {
    margin-bottom: var(--spacing-8);
    
    .header-content {
      display: flex;
      align-items: flex-end;
      justify-content: space-between;
      gap: var(--spacing-6);
      flex-wrap: wrap;
      
      .welcome-section {
        flex: 1;
        min-width: 300px;
        
        .dashboard-title {
          font-size: var(--text-3xl);
          font-weight: var(--font-bold);
          color: var(--text-primary);
          margin: 0 0 var(--spacing-2) 0;
          background: linear-gradient(135deg, var(--primary-600), var(--primary-500));
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }
        
        .dashboard-subtitle {
          font-size: var(--text-lg);
          color: var(--text-secondary);
          margin: 0;
          font-weight: var(--font-normal);
        }
      }
      
      .header-actions {
        display: flex;
        align-items: center;
        gap: var(--spacing-4);
        
        .date-selector {
          .date-picker {
            .el-input__wrapper {
              border-radius: var(--radius-xl);
              border: 1px solid var(--border-light);
              box-shadow: var(--shadow-sm);
              transition: all var(--transition-base);
              
              &:hover {
                border-color: var(--primary-500);
                box-shadow: var(--shadow-md);
              }
            }
          }
        }
        
        .action-buttons {
          display: flex;
          align-items: center;
          gap: var(--spacing-3);
          
          .el-button {
            border-radius: var(--radius-xl);
            font-weight: var(--font-medium);
            padding: var(--spacing-3) var(--spacing-6);
            transition: all var(--transition-base);
            
            &:hover {
              transform: translateY(-2px);
              box-shadow: var(--shadow-lg);
            }
          }
          
          .export-btn {
            background: linear-gradient(135deg, var(--success-500), var(--success-600));
            border: none;
            color: white;
            
            &:hover {
              background: linear-gradient(135deg, var(--success-600), var(--success-700));
            }
          }
        }
      }
    }
  }
  
  // 核心指标区域
  .metrics-section {
    margin-bottom: var(--spacing-8);
    
    .metrics-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: var(--spacing-6);
      
      .metric-card {
        background: var(--bg-primary);
        border-radius: var(--radius-2xl);
        padding: var(--spacing-6);
        box-shadow: var(--shadow-md);
        border: 1px solid var(--border-light);
        position: relative;
        overflow: hidden;
        cursor: pointer;
        transition: all var(--transition-base);
        
        &:hover {
          transform: translateY(-4px);
          box-shadow: var(--shadow-xl);
          
          .metric-overlay {
            opacity: 1;
          }
        }
        
        &.primary .metric-icon {
          background: linear-gradient(135deg, var(--primary-500), var(--primary-600));
        }
        
        &.success .metric-icon {
          background: linear-gradient(135deg, var(--success-500), var(--success-600));
        }
        
        &.info .metric-icon {
          background: linear-gradient(135deg, #3B82F6, #1D4ED8);
        }
        
        &.warning .metric-icon {
          background: linear-gradient(135deg, var(--warning-500), var(--warning-600));
        }
        
        .metric-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: var(--spacing-4);
          
          .metric-icon {
            width: 64px;
            height: 64px;
            border-radius: var(--radius-xl);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 24px;
            box-shadow: var(--shadow-lg);
          }
          
          .metric-trend {
            display: flex;
            align-items: center;
            gap: var(--spacing-1);
            font-size: var(--text-sm);
            font-weight: var(--font-semibold);
            padding: var(--spacing-1) var(--spacing-2);
            border-radius: var(--radius-full);
            
            &.positive {
              color: var(--success-600);
              background: var(--success-50);
            }
            
            &.negative {
              color: var(--error-600);
              background: var(--error-50);
            }
          }
        }
        
        .metric-content {
          margin-bottom: var(--spacing-4);
          
          .metric-value {
            font-size: var(--text-3xl);
            font-weight: var(--font-bold);
            color: var(--text-primary);
            margin-bottom: var(--spacing-2);
            line-height: 1.2;
          }
          
          .metric-title {
            font-size: var(--text-lg);
            font-weight: var(--font-semibold);
            color: var(--text-primary);
            margin-bottom: var(--spacing-1);
          }
          
          .metric-subtitle {
            font-size: var(--text-sm);
            color: var(--text-tertiary);
          }
        }
        
        .metric-chart {
          height: 60px;
          
          .mini-chart {
            width: 100%;
            height: 100%;
          }
        }
        
        .metric-overlay {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.8);
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          gap: var(--spacing-2);
          color: white;
          opacity: 0;
          transition: all var(--transition-base);
          
          .el-icon {
            font-size: 32px;
          }
          
          span {
            font-size: var(--text-lg);
            font-weight: var(--font-medium);
          }
        }
      }
    }
  }
  
  // 图表分析区域
  .charts-section {
    margin-bottom: var(--spacing-8);
    
    .charts-grid {
      display: grid;
      grid-template-columns: 1fr;
      gap: var(--spacing-6);
      
      .chart-card {
        background: var(--bg-primary);
        border-radius: var(--radius-2xl);
        box-shadow: var(--shadow-md);
        border: 1px solid var(--border-light);
        overflow: hidden;
        transition: all var(--transition-base);
        
        &:hover {
          box-shadow: var(--shadow-lg);
        }
        
        &.main-chart {
          grid-column: 1 / -1;
        }
        
        .chart-header {
          padding: var(--spacing-6);
          border-bottom: 1px solid var(--border-light);
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: var(--spacing-4);
          flex-wrap: wrap;
          
          .chart-title-section {
            .chart-title {
              display: flex;
              align-items: center;
              gap: var(--spacing-2);
              font-size: var(--text-lg);
              font-weight: var(--font-semibold);
              color: var(--text-primary);
              margin: 0 0 var(--spacing-1) 0;
              
              .el-icon {
                color: var(--primary-500);
              }
            }
            
            .chart-subtitle {
              font-size: var(--text-sm);
              color: var(--text-tertiary);
              margin: 0;
            }
          }
          
          .chart-controls {
            display: flex;
            align-items: center;
            gap: var(--spacing-3);
            
            .chart-type-selector,
            .period-selector {
              .el-radio-button {
                .el-radio-button__inner {
                  border-radius: var(--radius-md);
                  border: 1px solid var(--border-light);
                  transition: all var(--transition-fast);
                  
                  &:hover {
                    border-color: var(--primary-500);
                    color: var(--primary-500);
                  }
                }
                
                &.is-active {
                  .el-radio-button__inner {
                    background: var(--primary-500);
                    border-color: var(--primary-500);
                    color: white;
                  }
                }
              }
            }
          }
        }
        
        .chart-body {
          padding: var(--spacing-4);
          
          .chart-container {
            height: 400px;
            width: 100%;
          }
        }
        
        .chart-footer {
          padding: var(--spacing-4) var(--spacing-6);
          border-top: 1px solid var(--border-light);
          background: var(--bg-secondary);
          
          .chart-stats {
            display: flex;
            align-items: center;
            gap: var(--spacing-6);
            
            .stat-item {
              display: flex;
              flex-direction: column;
              align-items: center;
              gap: var(--spacing-1);
              
              .stat-label {
                font-size: var(--text-xs);
                color: var(--text-tertiary);
                text-transform: uppercase;
                letter-spacing: 0.5px;
              }
              
              .stat-value {
                font-size: var(--text-lg);
                font-weight: var(--font-semibold);
                color: var(--text-primary);
                
                &.positive {
                  color: var(--success-600);
                }
              }
            }
          }
          
          .product-legend {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
            gap: var(--spacing-2);
            
            .legend-item {
              display: flex;
              align-items: center;
              gap: var(--spacing-2);
              padding: var(--spacing-2);
              border-radius: var(--radius-md);
              transition: all var(--transition-fast);
              
              &:hover {
                background: var(--bg-hover);
              }
              
              .legend-color {
                width: 12px;
                height: 12px;
                border-radius: var(--radius-full);
                flex-shrink: 0;
              }
              
              .legend-name {
                flex: 1;
                font-size: var(--text-sm);
                color: var(--text-secondary);
              }
              
              .legend-value {
                font-size: var(--text-sm);
                font-weight: var(--font-medium);
                color: var(--text-primary);
              }
            }
          }
        }
      }
      
      // 响应式网格布局
      @media (min-width: 1024px) {
        grid-template-columns: 1fr 1fr;
        
        &:has(.main-chart) {
          grid-template-rows: auto auto auto;
          
          .chart-card:not(.main-chart) {
            .chart-body .chart-container {
              height: 300px;
            }
          }
        }
      }
    }
  }
  
  // 最近活动区域
  .activities-section {
    .section-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: var(--spacing-6);
      
      .section-title {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        font-size: var(--text-xl);
        font-weight: var(--font-semibold);
        color: var(--text-primary);
        margin: 0;
        
        .el-icon {
          color: var(--primary-500);
        }
      }
    }
    
    .activities-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
      gap: var(--spacing-6);
      
      .activity-card {
        background: var(--bg-primary);
        border-radius: var(--radius-xl);
        box-shadow: var(--shadow-md);
        border: 1px solid var(--border-light);
        overflow: hidden;
        
        .activity-header {
          padding: var(--spacing-5) var(--spacing-6);
          border-bottom: 1px solid var(--border-light);
          display: flex;
          align-items: center;
          justify-content: space-between;
          
          .activity-title {
            display: flex;
            align-items: center;
            gap: var(--spacing-2);
            font-size: var(--text-base);
            font-weight: var(--font-semibold);
            color: var(--text-primary);
            margin: 0;
            
            .el-icon {
              color: var(--primary-500);
            }
          }
        }
        
        .activity-list {
          max-height: 320px;
          overflow-y: auto;
          
          .activity-item {
            display: flex;
            align-items: center;
            gap: var(--spacing-3);
            padding: var(--spacing-4) var(--spacing-6);
            border-bottom: 1px solid var(--border-light);
            cursor: pointer;
            transition: all var(--transition-fast);
            
            &:last-child {
              border-bottom: none;
            }
            
            &:hover {
              background: var(--bg-hover);
            }
            
            &.notification-item.unread {
              background: var(--primary-50);
              border-left: 3px solid var(--primary-500);
            }
            
            .activity-avatar {
              width: 40px;
              height: 40px;
              border-radius: var(--radius-lg);
              display: flex;
              align-items: center;
              justify-content: center;
              color: white;
              font-size: 16px;
              flex-shrink: 0;
              background: var(--gray-500);
              
              &.success {
                background: var(--success-500);
              }
              
              &.warning {
                background: var(--warning-500);
              }
              
              &.info {
                background: #3B82F6;
              }
              
              &.error {
                background: var(--error-500);
              }
            }
            
            .activity-content {
              flex: 1;
              min-width: 0;
              
              .activity-main {
                font-size: var(--text-sm);
                color: var(--text-primary);
                margin-bottom: var(--spacing-1);
                
                .activity-user {
                  font-weight: var(--font-medium);
                }
                
                .activity-action {
                  color: var(--text-secondary);
                }
                
                .activity-target {
                  color: var(--primary-600);
                  font-weight: var(--font-medium);
                }
              }
              
              .activity-meta {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: var(--spacing-2);
                
                .activity-amount {
                  font-size: var(--text-sm);
                  font-weight: var(--font-semibold);
                  color: var(--success-600);
                }
                
                .activity-description {
                  font-size: var(--text-xs);
                  color: var(--text-secondary);
                  flex: 1;
                  min-width: 0;
                  overflow: hidden;
                  text-overflow: ellipsis;
                  white-space: nowrap;
                }
                
                .activity-time {
                  font-size: var(--text-xs);
                  color: var(--text-tertiary);
                  flex-shrink: 0;
                }
              }
            }
            
            .activity-status {
              flex-shrink: 0;
            }
            
            .notification-badge {
              width: 8px;
              height: 8px;
              background: var(--primary-500);
              border-radius: var(--radius-full);
              flex-shrink: 0;
            }
          }
        }
      }
    }
  }
}

// 全屏图表弹窗
:deep(.chart-dialog) {
  .el-dialog__body {
    padding: 0;
    
    .expanded-chart {
      height: 70vh;
      width: 100%;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .dashboard-premium {
    padding: var(--spacing-4);
    
    .dashboard-header {
      .header-content {
        flex-direction: column;
        align-items: stretch;
        gap: var(--spacing-4);
        
        .header-actions {
          flex-direction: column;
          align-items: stretch;
          
          .action-buttons {
            justify-content: space-between;
          }
        }
      }
    }
    
    .metrics-section {
      .metrics-grid {
        grid-template-columns: 1fr;
        gap: var(--spacing-4);
      }
    }
    
    .charts-section {
      .charts-grid {
        grid-template-columns: 1fr;
        
        .chart-card {
          .chart-header {
            flex-direction: column;
            align-items: stretch;
            gap: var(--spacing-3);
            
            .chart-controls {
              justify-content: space-between;
              flex-wrap: wrap;
            }
          }
          
          .chart-body {
            .chart-container {
              height: 300px;
            }
          }
          
          .chart-footer {
            .chart-stats {
              flex-wrap: wrap;
              gap: var(--spacing-3);
            }
          }
        }
      }
    }
    
    .activities-section {
      .activities-grid {
        grid-template-columns: 1fr;
        gap: var(--spacing-4);
      }
    }
  }
}

// 动画效果
.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 自定义滚动条
.custom-scrollbar {
  &::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: var(--bg-secondary);
    border-radius: var(--radius-full);
  }
  
  &::-webkit-scrollbar-thumb {
    background: var(--border-default);
    border-radius: var(--radius-full);
    
    &:hover {
      background: var(--border-strong);
    }
  }
}
</style>