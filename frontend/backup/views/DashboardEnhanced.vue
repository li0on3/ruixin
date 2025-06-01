<template>
  <div class="dashboard-enhanced">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <h1 class="page-title">数据概览</h1>
        <p class="page-subtitle">欢迎回来，{{ userStore.userInfo?.username || '管理员' }}！这是您的业务概览。</p>
      </div>
      <div class="page-header-actions">
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
        />
        <el-button type="primary" @click="exportReport">
          <Download />
          导出报表
        </el-button>
      </div>
    </div>

    <!-- 核心指标卡片 -->
    <div class="metrics-grid">
      <div class="stat-card" v-for="metric in metrics" :key="metric.key">
        <div class="stat-icon" :class="metric.iconClass">
          <component :is="metric.icon" />
        </div>
        <div class="stat-content">
          <h3 class="stat-title">{{ metric.title }}</h3>
          <div class="stat-value">
            <CountUp
              :end-val="metric.value"
              :duration="2"
              :prefix="metric.prefix"
              :suffix="metric.suffix"
            />
          </div>
          <div class="stat-trend" :class="metric.trend > 0 ? 'up' : 'down'">
            <TrendCharts />
            <span>{{ Math.abs(metric.trend) }}%</span>
            <span class="trend-text">较昨日</span>
          </div>
        </div>
        <div class="stat-chart">
          <div :ref="el => miniCharts[metric.key] = el" class="mini-chart"></div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-section">
      <!-- 销售趋势 -->
      <div class="chart-container main-chart">
        <div class="chart-header">
          <h2 class="chart-title">
            <TrendCharts />
            销售趋势分析
          </h2>
          <div class="chart-actions">
            <el-radio-group v-model="trendType" size="small" @change="updateTrendChart">
              <el-radio-button value="revenue">销售额</el-radio-button>
              <el-radio-button value="orders">订单量</el-radio-button>
              <el-radio-button value="both">综合</el-radio-button>
            </el-radio-group>
            <el-radio-group v-model="trendPeriod" size="small" @change="updateTrendChart">
              <el-radio-button value="7">7天</el-radio-button>
              <el-radio-button value="30">30天</el-radio-button>
              <el-radio-button value="90">90天</el-radio-button>
            </el-radio-group>
          </div>
        </div>
        <div ref="trendChartRef" class="chart-body"></div>
      </div>

      <!-- 分布图表 -->
      <div class="chart-row">
        <!-- 商品销售分布 -->
        <div class="chart-container">
          <div class="chart-header">
            <h3 class="chart-title">
              <Coffee />
              热销商品TOP10
            </h3>
            <el-dropdown trigger="click">
              <el-button text>
                <More />
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="viewAllProducts">查看全部商品</el-dropdown-item>
                  <el-dropdown-item @click="exportProductData">导出数据</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div ref="productChartRef" class="chart-body"></div>
        </div>

        <!-- 时段分布 -->
        <div class="chart-container">
          <div class="chart-header">
            <h3 class="chart-title">
              <Clock />
              订单时段分布
            </h3>
          </div>
          <div ref="hourChartRef" class="chart-body"></div>
        </div>
      </div>
    </div>

    <!-- 数据列表区域 -->
    <div class="data-section">
      <!-- 实时订单流 -->
      <div class="data-card order-stream">
        <div class="data-header">
          <h3 class="data-title">
            <ShoppingCart />
            实时订单
          </h3>
          <el-button text @click="$router.push('/orders')">
            查看全部
            <ArrowRight />
          </el-button>
        </div>
        <div class="order-list">
          <TransitionGroup name="list-slide">
            <div class="order-item" v-for="order in recentOrders" :key="order.id">
              <div class="order-avatar">
                <component :is="getOrderIcon(order.status)" />
              </div>
              <div class="order-info">
                <div class="order-header">
                  <span class="order-no">{{ order.orderNo }}</span>
                  <span class="order-time">{{ formatTime(order.createdAt) }}</span>
                </div>
                <div class="order-detail">
                  <span class="distributor">{{ order.distributorName }}</span>
                  <span class="separator">·</span>
                  <span class="store">{{ order.storeName }}</span>
                </div>
              </div>
              <div class="order-amount">
                ¥{{ order.totalAmount }}
              </div>
              <div class="order-status">
                <el-tag :type="getOrderStatusType(order.status)" size="small">
                  {{ getOrderStatusText(order.status) }}
                </el-tag>
              </div>
            </div>
          </TransitionGroup>
        </div>
      </div>

      <!-- 分销商排行 -->
      <div class="data-card distributor-rank">
        <div class="data-header">
          <h3 class="data-title">
            <Trophy />
            分销商排行榜
          </h3>
          <el-dropdown trigger="click">
            <span class="rank-filter">
              {{ rankPeriod === 'today' ? '今日' : rankPeriod === 'week' ? '本周' : '本月' }}
              <ArrowDown />
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="changeRankPeriod('today')">今日</el-dropdown-item>
                <el-dropdown-item @click="changeRankPeriod('week')">本周</el-dropdown-item>
                <el-dropdown-item @click="changeRankPeriod('month')">本月</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="rank-list">
          <div class="rank-item" v-for="(item, index) in distributorRank" :key="item.id">
            <div class="rank-no" :class="getRankClass(index + 1)">
              {{ index + 1 }}
            </div>
            <div class="rank-info">
              <div class="rank-name">{{ item.name }}</div>
              <div class="rank-stats">
                <span>{{ item.orders }} 单</span>
                <span class="separator">·</span>
                <span>¥{{ formatNumber(item.revenue) }}</span>
              </div>
            </div>
            <div class="rank-progress">
              <el-progress
                :percentage="(item.revenue / maxRevenue) * 100"
                :stroke-width="6"
                :show-text="false"
                :color="getRankColor(index + 1)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick, computed, watch } from 'vue'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { darkTheme, applyDarkTheme, getEChartsTheme } from '@/utils/echarts-dark-theme'
import { useTheme } from '@/plugins/theme'
import CountUp from 'vue-countup-v3'
import {
  ShoppingCart, Money, UserFilled, CreditCard, TrendCharts,
  Download, Coffee, Clock, ArrowRight, Trophy, ArrowDown,
  More, SuccessFilled, WarningFilled, CircleCloseFilled
} from '@element-plus/icons-vue'
import {
  getDashboardStatistics,
  getOrderTrend,
  getHotGoods,
  getRecentOrders,
  getHourDistribution,
  getDistributorRank
} from '@/api/dashboard'
import dayjs from 'dayjs'

const userStore = useUserStore()
const { isDarkTheme } = useTheme()

// 注册深色主题
applyDarkTheme(echarts)

// 日期范围
const dateRange = ref([
  dayjs().subtract(7, 'day').format('YYYY-MM-DD'),
  dayjs().format('YYYY-MM-DD')
])

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

// 核心指标
const metrics = ref([
  {
    key: 'orders',
    title: '今日订单',
    value: 0,
    trend: 0,
    icon: ShoppingCart,
    iconClass: 'primary',
    prefix: '',
    suffix: '',
    miniChartData: []
  },
  {
    key: 'revenue',
    title: '今日销售额',
    value: 0,
    trend: 0,
    icon: Money,
    iconClass: 'success',
    prefix: '¥',
    suffix: '',
    miniChartData: []
  },
  {
    key: 'distributors',
    title: '活跃分销商',
    value: 0,
    trend: 0,
    icon: UserFilled,
    iconClass: 'warning',
    prefix: '',
    suffix: '',
    miniChartData: []
  },
  {
    key: 'cards',
    title: '可用卡片',
    value: 0,
    trend: 0,
    icon: CreditCard,
    iconClass: 'danger',
    prefix: '',
    suffix: '',
    miniChartData: []
  }
])

// 图表相关
const trendType = ref('both')
const trendPeriod = ref('7')
const trendChartRef = ref()
const productChartRef = ref()
const hourChartRef = ref()
const miniCharts = ref({})

// 数据
const recentOrders = ref([])
const distributorRank = ref([])
const rankPeriod = ref('today')

// 图表实例
let trendChart = null
let productChart = null
let hourChart = null
const miniChartInstances = {}

// 计算最大营收（用于进度条）
const maxRevenue = computed(() => {
  if (distributorRank.value.length === 0) return 1
  return Math.max(...distributorRank.value.map(item => item.revenue))
})

// 加载统计数据
const loadStatistics = async () => {
  try {
    const { data } = await getDashboardStatistics()
    
    // 更新指标数据
    metrics.value[0].value = data.todayOrders || 0
    metrics.value[0].trend = data.orderGrowth || 0
    metrics.value[0].miniChartData = data.ordersTrend || []
    metrics.value[1].value = data.todayAmount || 0
    metrics.value[1].trend = data.amountGrowth || 0
    metrics.value[1].miniChartData = data.amountTrend || []
    metrics.value[2].value = data.activeDistributors || 0
    metrics.value[2].trend = data.distributorGrowth || 0
    metrics.value[2].miniChartData = data.distributorTrend || []
    metrics.value[3].value = data.activeCards || 0
    metrics.value[3].trend = data.cardGrowth || 0
    metrics.value[3].miniChartData = data.cardTrend || []
    
    // 初始化迷你图表
    nextTick(() => {
      initMiniCharts()
    })
  } catch (error) {
    console.error('Failed to load statistics:', error)
  }
}

// 初始化迷你图表
const initMiniCharts = () => {
  metrics.value.forEach(metric => {
    const el = miniCharts.value[metric.key]
    if (el && !miniChartInstances[metric.key]) {
      miniChartInstances[metric.key] = echarts.init(el, getEChartsTheme())
      
      // 如果有真实的迷你图表数据，使用真实数据，否则隐藏图表
      const hasData = metric.miniChartData && metric.miniChartData.length > 0
      
      if (hasData) {
        const option = {
          grid: {
            top: 0,
            right: 0,
            bottom: 0,
            left: 0
          },
          xAxis: {
            type: 'category',
            show: false,
            data: metric.miniChartData.map((_, index) => index + 1)
          },
          yAxis: {
            type: 'value',
            show: false
          },
          series: [{
            type: 'line',
            smooth: true,
            symbol: 'none',
            lineStyle: {
              width: 2,
              color: getMetricColor(metric.iconClass)
            },
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: getMetricColor(metric.iconClass, 0.3) },
                { offset: 1, color: getMetricColor(metric.iconClass, 0) }
              ])
            },
            data: metric.miniChartData
          }]
        }
        
        miniChartInstances[metric.key].setOption(option)
      } else {
        // 没有数据时显示空白图表
        const option = {
          grid: { top: 0, right: 0, bottom: 0, left: 0 },
          xAxis: { type: 'category', show: false, data: ['1'] },
          yAxis: { type: 'value', show: false },
          series: []
        }
        miniChartInstances[metric.key].setOption(option)
      }
    }
  })
}

// 更新趋势图表
const updateTrendChart = async () => {
  try {
    // 根据trendPeriod计算日期范围
    const days = parseInt(trendPeriod.value)
    const endDate = dayjs().format('YYYY-MM-DD')
    const startDate = dayjs().subtract(days - 1, 'day').format('YYYY-MM-DD')
    
    const { data } = await getOrderTrend(`${days}`) // 传递天数字符串
    
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          crossStyle: {
            color: '#999'
          }
        }
      },
      legend: {
        data: trendType.value === 'both' ? ['销售额', '订单量'] : [trendType.value === 'revenue' ? '销售额' : '订单量']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: data.map(item => dayjs(item.date).format('MM-DD'))
      },
      yAxis: trendType.value === 'both' ? [
        {
          type: 'value',
          name: '销售额',
          axisLabel: {
            formatter: '¥{value}'
          }
        },
        {
          type: 'value',
          name: '订单量',
          axisLabel: {
            formatter: '{value}单'
          }
        }
      ] : [{
        type: 'value',
        name: trendType.value === 'revenue' ? '销售额' : '订单量',
        axisLabel: {
          formatter: trendType.value === 'revenue' ? '¥{value}' : '{value}单'
        }
      }],
      series: trendType.value === 'both' ? [
        {
          name: '销售额',
          type: 'line',
          smooth: true,
          data: data.map(item => item.revenue || 0),
          itemStyle: {
            color: '#10b981'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
              { offset: 1, color: 'rgba(16, 185, 129, 0)' }
            ])
          }
        },
        {
          name: '订单量',
          type: 'line',
          smooth: true,
          yAxisIndex: 1,
          data: data.map(item => item.orders || 0),
          itemStyle: {
            color: '#3b82f6'
          }
        }
      ] : [{
        name: trendType.value === 'revenue' ? '销售额' : '订单量',
        type: 'line',
        smooth: true,
        data: data.map(item => trendType.value === 'revenue' ? (item.revenue || 0) : (item.orders || 0)),
        itemStyle: {
          color: trendType.value === 'revenue' ? '#10b981' : '#3b82f6'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: trendType.value === 'revenue' ? 'rgba(16, 185, 129, 0.3)' : 'rgba(59, 130, 246, 0.3)' },
            { offset: 1, color: trendType.value === 'revenue' ? 'rgba(16, 185, 129, 0)' : 'rgba(59, 130, 246, 0)' }
          ])
        }
      }]
    }
    
    trendChart.setOption(option)
  } catch (error) {
    console.error('Failed to update trend chart:', error)
  }
}

// 加载商品销售数据
const loadProductData = async () => {
  try {
    const { data } = await getHotGoods()
    
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        type: 'scroll',
        orient: 'vertical',
        right: 10,
        top: 20,
        bottom: 20
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['40%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data.map((item, index) => ({
          name: item.name,
          value: item.count,
          itemStyle: {
            color: getProductColor(index)
          }
        }))
      }]
    }
    
    productChart.setOption(option)
  } catch (error) {
    console.error('Failed to load product data:', error)
  }
}

// 加载时段分布数据
const loadHourData = async () => {
  try {
    const { data } = await getHourDistribution()
    
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: Array.from({ length: 24 }, (_, i) => `${i}时`)
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        type: 'bar',
        data: data,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#8b5cf6' },
            { offset: 1, color: '#6366f1' }
          ]),
          borderRadius: [4, 4, 0, 0]
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#a78bfa' },
              { offset: 1, color: '#818cf8' }
            ])
          }
        }
      }]
    }
    
    hourChart.setOption(option)
  } catch (error) {
    console.error('Failed to load hour data:', error)
  }
}

// 加载最新订单
const loadRecentOrders = async () => {
  try {
    const { data } = await getRecentOrders()
    recentOrders.value = data
  } catch (error) {
    console.error('Failed to load recent orders:', error)
  }
}

// 加载分销商排行
const loadDistributorRank = async () => {
  try {
    const { data } = await getDistributorRank(rankPeriod.value)
    distributorRank.value = data
  } catch (error) {
    console.error('Failed to load distributor rank:', error)
  }
}

// 处理日期变化
const handleDateChange = () => {
  loadStatistics()
  updateTrendChart()
}

// 改变排行榜时间周期
const changeRankPeriod = (period) => {
  rankPeriod.value = period
  loadDistributorRank()
}

// 导出报表
const exportReport = async () => {
  try {
    // 获取当前所有数据
    const dashboardData = {
      statistics: {
        todayOrders: metrics.value[0].value,
        todayAmount: metrics.value[1].value,
        activeDistributors: metrics.value[2].value,
        activeCards: metrics.value[3].value
      },
      orderTrend: [], // 趋势数据需要重新获取
      hotProducts: [], // 热销产品数据需要重新获取
      recentOrders: recentOrders.value,
      distributorRank: distributorRank.value
    }
    
    // 获取趋势数据
    try {
      const trendResponse = await getOrderTrend(trendPeriod.value)
      dashboardData.orderTrend = trendResponse.data || []
    } catch (error) {
      console.warn('获取趋势数据失败:', error)
    }
    
    // 获取热销产品数据
    try {
      const hotGoodsResponse = await getHotGoods()
      dashboardData.hotProducts = hotGoodsResponse.data || []
    } catch (error) {
      console.warn('获取热销产品数据失败:', error)
    }
    
    // 动态导入导出函数
    const { exportDashboardReport } = await import('@/utils/export')
    
    exportDashboardReport(dashboardData, {
      filename: `仪表盘报表_${dateRange.value[0]}_${dateRange.value[1]}.xlsx`
    })
  } catch (error) {
    console.error('导出报表失败:', error)
    ElMessage.error('导出报表失败')
  }
}

// 查看全部商品
const viewAllProducts = () => {
  ElMessage.info('跳转到商品分析页面...')
}

// 导出商品数据
const exportProductData = () => {
  ElMessage.success('商品数据导出中...')
}

// 工具函数
const formatTime = (time) => {
  return dayjs(time).format('HH:mm')
}

const formatNumber = (num) => {
  return num.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getOrderIcon = (status) => {
  const icons = {
    0: WarningFilled,
    1: WarningFilled,
    2: SuccessFilled,
    3: CircleCloseFilled
  }
  return icons[status] || WarningFilled
}

const getOrderStatusType = (status) => {
  const types = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger'
  }
  return types[status] || 'info'
}

const getOrderStatusText = (status) => {
  const texts = {
    0: '待处理',
    1: '处理中',
    2: '已完成',
    3: '失败'
  }
  return texts[status] || '未知'
}

const getRankClass = (rank) => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'silver'
  if (rank === 3) return 'bronze'
  return ''
}

const getRankColor = (rank) => {
  if (rank === 1) return '#fbbf24'
  if (rank === 2) return '#94a3b8'
  if (rank === 3) return '#f97316'
  return '#3b82f6'
}

const getMetricColor = (type, opacity = 1) => {
  const colors = {
    primary: `rgba(59, 130, 246, ${opacity})`,
    success: `rgba(16, 185, 129, ${opacity})`,
    warning: `rgba(245, 158, 11, ${opacity})`,
    danger: `rgba(239, 68, 68, ${opacity})`
  }
  return colors[type] || colors.primary
}

const getProductColor = (index) => {
  const colors = [
    '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6',
    '#6366f1', '#ec4899', '#14b8a6', '#f97316', '#06b6d4'
  ]
  return colors[index % colors.length]
}

// 已移除 generateRandomData 函数，使用真实数据

// 生命周期
onMounted(() => {
  // 初始化主图表
  trendChart = echarts.init(trendChartRef.value, getEChartsTheme())
  productChart = echarts.init(productChartRef.value, getEChartsTheme())
  hourChart = echarts.init(hourChartRef.value, getEChartsTheme())
  
  // 响应式
  window.addEventListener('resize', () => {
    trendChart?.resize()
    productChart?.resize()
    hourChart?.resize()
    Object.values(miniChartInstances).forEach(chart => chart?.resize())
  })
  
  // 加载数据
  loadStatistics()
  updateTrendChart()
  loadProductData()
  loadHourData()
  loadRecentOrders()
  loadDistributorRank()
  
  // 定时刷新
  const timer = setInterval(() => {
    loadRecentOrders()
    loadDistributorRank()
  }, 30000)
  
  // 监听主题变化
  watch(() => isDarkTheme(), () => {
    // 重新初始化所有图表以应用新主题
    nextTick(() => {
      // 销毁旧图表实例
      trendChart?.dispose()
      productChart?.dispose()
      hourChart?.dispose()
      Object.values(miniChartInstances).forEach(chart => chart?.dispose())
      
      // 清空实例引用
      Object.keys(miniChartInstances).forEach(key => {
        delete miniChartInstances[key]
      })
      
      // 重新初始化图表
      trendChart = echarts.init(trendChartRef.value, getEChartsTheme())
      productChart = echarts.init(productChartRef.value, getEChartsTheme())
      hourChart = echarts.init(hourChartRef.value, getEChartsTheme())
      
      // 重新加载数据
      updateTrendChart()
      loadProductData()
      loadHourData()
      initMiniCharts()
    })
  })
  
  onUnmounted(() => {
    clearInterval(timer)
    trendChart?.dispose()
    productChart?.dispose()
    hourChart?.dispose()
    Object.values(miniChartInstances).forEach(chart => chart?.dispose())
  })
})
</script>

<style lang="scss" scoped>
.dashboard-enhanced {
  padding: 0;
  
  // 页面头部
  .page-header {
    margin-bottom: var(--spacing-6);
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: var(--spacing-4);
    
    .page-header-content {
      .page-title {
        font-size: 1.75rem;
        font-weight: 700;
        color: var(--text-primary);
        margin: 0 0 var(--spacing-2);
      }
      
      .page-subtitle {
        font-size: 0.875rem;
        color: var(--text-secondary);
        margin: 0;
      }
    }
    
    .page-header-actions {
      display: flex;
      gap: var(--spacing-3);
      align-items: center;
    }
  }
  
  // 指标卡片网格
  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: var(--spacing-6);
    margin-bottom: var(--spacing-8);
    
    .stat-card {
      background: var(--bg-primary);
      border-radius: var(--radius-xl);
      padding: var(--spacing-6);
      position: relative;
      overflow: hidden;
      display: flex;
      align-items: center;
      gap: var(--spacing-4);
      transition: all var(--transition-base);
      border: 1px solid var(--border-light);
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: var(--shadow-lg);
        border-color: var(--primary-200);
      }
      
      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: var(--radius-xl);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 1.75rem;
        flex-shrink: 0;
        
        &.primary {
          background: var(--primary-100);
          color: var(--primary-600);
        }
        
        &.success {
          background: var(--success-light);
          color: var(--success);
        }
        
        &.warning {
          background: var(--warning-light);
          color: var(--warning);
        }
        
        &.danger {
          background: var(--danger-light);
          color: var(--danger);
        }
      }
      
      .stat-content {
        flex: 1;
        min-width: 0;
        
        .stat-title {
          font-size: 0.875rem;
          font-weight: 500;
          color: var(--text-secondary);
          margin: 0 0 var(--spacing-2);
        }
        
        .stat-value {
          font-size: 2rem;
          font-weight: 700;
          color: var(--text-primary);
          line-height: 1.2;
          margin-bottom: var(--spacing-2);
        }
        
        .stat-trend {
          display: flex;
          align-items: center;
          gap: var(--spacing-1);
          font-size: 0.875rem;
          font-weight: 500;
          
          &.up {
            color: var(--success);
            svg {
              transform: rotate(-90deg);
            }
          }
          
          &.down {
            color: var(--danger);
            svg {
              transform: rotate(90deg);
            }
          }
          
          .trend-text {
            color: var(--text-tertiary);
            font-weight: 400;
          }
        }
      }
      
      .stat-chart {
        position: absolute;
        right: 0;
        top: 0;
        width: 100px;
        height: 100%;
        opacity: 0.5;
        
        .mini-chart {
          width: 100%;
          height: 100%;
        }
      }
    }
  }
  
  // 图表区域
  .charts-section {
    margin-bottom: var(--spacing-8);
    
    .chart-container {
      background: var(--bg-primary);
      border-radius: var(--radius-xl);
      padding: var(--spacing-6);
      margin-bottom: var(--spacing-6);
      box-shadow: var(--shadow-sm);
      border: 1px solid var(--border-light);
      
      &.main-chart {
        .chart-body {
          height: 400px;
        }
      }
      
      .chart-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--spacing-6);
        
        .chart-title {
          font-size: 1.125rem;
          font-weight: 600;
          color: var(--text-primary);
          margin: 0;
          display: flex;
          align-items: center;
          gap: var(--spacing-2);
          
          svg {
            font-size: 1.25rem;
            color: var(--primary-500);
          }
        }
        
        .chart-actions {
          display: flex;
          gap: var(--spacing-3);
        }
      }
      
      .chart-body {
        height: 300px;
      }
    }
    
    .chart-row {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
      gap: var(--spacing-6);
    }
  }
  
  // 数据列表区域
  .data-section {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: var(--spacing-6);
    
    @media (max-width: 1440px) {
      grid-template-columns: 1fr;
    }
    
    .data-card {
      background: var(--bg-primary);
      border-radius: var(--radius-xl);
      padding: var(--spacing-6);
      box-shadow: var(--shadow-sm);
      border: 1px solid var(--border-light);
      
      .data-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--spacing-4);
        
        .data-title {
          font-size: 1.125rem;
          font-weight: 600;
          color: var(--text-primary);
          margin: 0;
          display: flex;
          align-items: center;
          gap: var(--spacing-2);
          
          svg {
            font-size: 1.25rem;
            color: var(--primary-500);
          }
        }
      }
    }
    
    // 订单流
    .order-stream {
      .order-list {
        max-height: 500px;
        overflow-y: auto;
        
        &::-webkit-scrollbar {
          width: 6px;
        }
        
        &::-webkit-scrollbar-track {
          background: var(--bg-secondary);
          border-radius: 3px;
        }
        
        &::-webkit-scrollbar-thumb {
          background: var(--gray-300);
          border-radius: 3px;
          
          &:hover {
            background: var(--gray-400);
          }
        }
      }
      
      .order-item {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        padding: var(--spacing-3) 0;
        border-bottom: 1px solid var(--border-light);
        
        &:last-child {
          border-bottom: none;
        }
        
        .order-avatar {
          width: 40px;
          height: 40px;
          border-radius: var(--radius-lg);
          display: flex;
          align-items: center;
          justify-content: center;
          background: var(--bg-secondary);
          color: var(--text-tertiary);
          flex-shrink: 0;
        }
        
        .order-info {
          flex: 1;
          min-width: 0;
          
          .order-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: var(--spacing-1);
            
            .order-no {
              font-weight: 500;
              color: var(--text-primary);
              font-size: 0.875rem;
            }
            
            .order-time {
              font-size: 0.75rem;
              color: var(--text-tertiary);
            }
          }
          
          .order-detail {
            font-size: 0.8125rem;
            color: var(--text-secondary);
            
            .separator {
              margin: 0 var(--spacing-1);
              color: var(--text-tertiary);
            }
          }
        }
        
        .order-amount {
          font-weight: 600;
          color: var(--text-primary);
          font-size: 0.9375rem;
        }
        
        .order-status {
          flex-shrink: 0;
        }
      }
    }
    
    // 分销商排行
    .distributor-rank {
      .rank-filter {
        display: flex;
        align-items: center;
        gap: var(--spacing-1);
        cursor: pointer;
        font-size: 0.875rem;
        color: var(--text-secondary);
        
        &:hover {
          color: var(--primary-500);
        }
      }
      
      .rank-list {
        .rank-item {
          display: flex;
          align-items: center;
          gap: var(--spacing-3);
          margin-bottom: var(--spacing-4);
          
          &:last-child {
            margin-bottom: 0;
          }
          
          .rank-no {
            width: 32px;
            height: 32px;
            border-radius: var(--radius-lg);
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            font-size: 0.875rem;
            background: var(--bg-secondary);
            color: var(--text-secondary);
            flex-shrink: 0;
            
            &.gold {
              background: linear-gradient(135deg, #fbbf24, #f59e0b);
              color: white;
            }
            
            &.silver {
              background: linear-gradient(135deg, #cbd5e1, #94a3b8);
              color: white;
            }
            
            &.bronze {
              background: linear-gradient(135deg, #fdba74, #f97316);
              color: white;
            }
          }
          
          .rank-info {
            flex: 1;
            min-width: 0;
            
            .rank-name {
              font-weight: 500;
              color: var(--text-primary);
              margin-bottom: var(--spacing-1);
            }
            
            .rank-stats {
              font-size: 0.8125rem;
              color: var(--text-secondary);
              
              .separator {
                margin: 0 var(--spacing-1);
                color: var(--text-tertiary);
              }
            }
          }
          
          .rank-progress {
            width: 100px;
            flex-shrink: 0;
          }
        }
      }
    }
  }
}

// 列表动画
.list-slide-move,
.list-slide-enter-active,
.list-slide-leave-active {
  transition: all 0.3s ease;
}

.list-slide-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.list-slide-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-slide-leave-active {
  position: absolute;
  width: 100%;
}
</style>