<template>
  <div class="statistics-enhanced-container">
    <!-- 顶部筛选区域 -->
    <div class="filter-section">
      <el-card class="filter-card">
        <div class="filter-content">
          <div class="filter-left">
            <h2 class="filter-title">
              <el-icon><TrendCharts /></el-icon>
              数据统计
            </h2>
          </div>
          <div class="filter-right">
            <el-form :model="filterForm" inline>
              <el-form-item label="时间范围">
                <el-date-picker
                  v-model="dateRange"
                  type="daterange"
                  range-separator="至"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  format="YYYY-MM-DD"
                  value-format="YYYY-MM-DD"
                  @change="handleDateChange"
                  class="modern-date-picker"
                />
              </el-form-item>
              <el-form-item>
                <el-button-group class="quick-select-group">
                  <el-button @click="selectToday" size="small">今日</el-button>
                  <el-button @click="selectYesterday" size="small">昨日</el-button>
                  <el-button @click="selectWeek" size="small">本周</el-button>
                  <el-button @click="selectMonth" size="small">本月</el-button>
                </el-button-group>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="loadAllData" :loading="loading">
                  <el-icon><Refresh /></el-icon>
                  查询
                </el-button>
                <el-button @click="exportData" :loading="exportLoading">
                  <el-icon><Download /></el-icon>
                  导出
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 核心指标卡片 -->
    <div class="metrics-section">
      <transition-group name="metric-card" tag="div" class="metrics-grid">
        <div
          v-for="(metric, index) in metricsConfig"
          :key="metric.key"
          class="metric-card-wrapper"
          :style="{ '--delay': index * 0.1 + 's' }"
        >
          <el-card class="metric-card" :class="`metric-${metric.type}`" shadow="hover">
            <div class="metric-header">
              <div class="metric-icon-wrapper">
                <el-icon class="metric-icon" :style="{ color: metric.color }">
                  <component :is="iconMap[metric.icon]" />
                </el-icon>
              </div>
              <div class="metric-info">
                <h3 class="metric-title">{{ metric.title }}</h3>
                <div class="metric-value">
                  <AnimatedNumber
                    :value="metrics[metric.valueKey] || 0"
                    :format="metric.formatter"
                    :duration="1500"
                    :delay="index * 200"
                    easing="easeOutCubic"
                  />
                </div>
              </div>
            </div>
            <div class="metric-footer">
              <div class="metric-trend">
                <span class="trend-label">环比</span>
                <span :class="['trend-value', getTrendClass(metrics[metric.growthKey])]">
                  <el-icon v-if="metrics[metric.growthKey] > 0" class="trend-icon">
                    <ArrowUp />
                  </el-icon>
                  <el-icon v-else-if="metrics[metric.growthKey] < 0" class="trend-icon">
                    <ArrowDown />
                  </el-icon>
                  {{ formatGrowth(metrics[metric.growthKey]) }}
                </span>
              </div>
            </div>
          </el-card>
        </div>
      </transition-group>
    </div>

    <!-- 主要图表区域 -->
    <div class="charts-section">
      <!-- 销售趋势图 -->
      <el-card class="chart-card trend-chart" shadow="hover">
        <template #header>
          <div class="chart-header">
            <div class="chart-title">
              <el-icon><TrendCharts /></el-icon>
              销售趋势分析
            </div>
            <div class="chart-controls">
              <el-radio-group v-model="trendType" @change="loadSalesTrend" size="small">
                <el-radio-button value="revenue">销售额</el-radio-button>
                <el-radio-button value="orders">订单数</el-radio-button>
                <el-radio-button value="both">综合</el-radio-button>
              </el-radio-group>
            </div>
          </div>
        </template>
        <div class="chart-content">
          <LoadingAnimation
            v-if="chartLoading.trend"
            type="circle"
            :visible="chartLoading.trend"
          />
          <div
            v-else
            ref="salesTrendChart"
            class="chart-container"
            style="height: 400px"
          ></div>
        </div>
      </el-card>
    </div>

    <!-- 分析图表网格 -->
    <div class="analytics-grid">
      <!-- 分销商排行 -->
      <el-card class="chart-card distributor-rank" shadow="hover">
        <template #header>
          <div class="chart-header">
            <div class="chart-title">
              <el-icon><User /></el-icon>
              分销商排行
            </div>
            <div class="chart-controls">
              <el-radio-group v-model="distributorRankType" @change="loadDistributorRank" size="small">
                <el-radio-button value="revenue">销售额</el-radio-button>
                <el-radio-button value="orders">订单数</el-radio-button>
              </el-radio-group>
            </div>
          </div>
        </template>
        <div class="chart-content">
          <LoadingAnimation
            v-if="chartLoading.distributor"
            type="bars"
            :visible="chartLoading.distributor"
          />
          <div
            v-else
            ref="distributorRankChart"
            class="chart-container"
            style="height: 350px"
          ></div>
        </div>
      </el-card>

      <!-- 商品分析 -->
      <el-card class="chart-card product-analysis" shadow="hover">
        <template #header>
          <div class="chart-header">
            <div class="chart-title">
              <el-icon><ShoppingBag /></el-icon>
              商品销售分析
            </div>
            <div class="chart-controls">
              <el-tooltip content="切换图表类型">
                <el-button-group size="small">
                  <el-button 
                    :type="productChartType === 'pie' ? 'primary' : ''" 
                    @click="productChartType = 'pie'"
                  >
                    饼图
                  </el-button>
                  <el-button 
                    :type="productChartType === 'bar' ? 'primary' : ''" 
                    @click="productChartType = 'bar'"
                  >
                    柱状图
                  </el-button>
                </el-button-group>
              </el-tooltip>
            </div>
          </div>
        </template>
        <div class="chart-content">
          <LoadingAnimation
            v-if="chartLoading.product"
            type="dots"
            :visible="chartLoading.product"
          />
          <div
            v-else
            ref="productAnalysisChart"
            class="chart-container"
            style="height: 350px"
          ></div>
        </div>
      </el-card>

      <!-- 时段分析 -->
      <el-card class="chart-card hour-analysis" shadow="hover">
        <template #header>
          <div class="chart-header">
            <div class="chart-title">
              <el-icon><Clock /></el-icon>
              订单时段分布
            </div>
          </div>
        </template>
        <div class="chart-content">
          <LoadingAnimation
            v-if="chartLoading.hour"
            type="skeleton"
            :visible="chartLoading.hour"
          />
          <div
            v-else
            ref="hourDistributionChart"
            class="chart-container"
            style="height: 300px"
          ></div>
        </div>
      </el-card>

    </div>

    <!-- 数据详情表格 -->
    <el-card class="data-table-section" shadow="hover">
      <template #header>
        <div class="table-header">
          <div class="table-title">
            <el-icon><Grid /></el-icon>
            数据详情
          </div>
          <div class="table-controls">
            <el-select 
              v-model="tableType" 
              @change="loadTableData" 
              class="table-type-select"
              size="default"
            >
              <el-option label="按日统计" value="daily" />
              <el-option label="按分销商" value="distributor" />
              <el-option label="按商品" value="product" />
              <el-option label="按门店" value="store" />
            </el-select>
          </div>
        </div>
      </template>
      
      <div class="table-content">
        <!-- 桌面端表格 -->
        <div class="desktop-only">
          <EnhancedTable
            :data="currentTableData"
            :columns="currentTableColumns"
            :loading="tableLoading"
            :show-export="true"
            :show-column-settings="true"
            export-filename="统计数据"
          />
        </div>
        
        <!-- 移动端卡片列表 -->
        <div class="mobile-only mobile-table-cards">
          <LoadingAnimation v-if="tableLoading" type="skeleton" :visible="tableLoading" />
          <div v-else class="mobile-card-list">
            <div
              v-for="(item, index) in currentTableData"
              :key="index"
              class="mobile-card-item"
            >
              <div class="card-item-header">
                <h4 class="card-item-title">
                  {{ getCardTitle(item) }}
                </h4>
                <div class="card-item-status">
                  {{ getCardStatus(item) }}
                </div>
              </div>
              <div class="card-item-content">
                <div
                  v-for="column in currentTableColumns.slice(1)"
                  :key="column.prop"
                  class="card-item-row"
                >
                  <span class="card-item-label">{{ column.label }}</span>
                  <span class="card-item-value">
                    {{ column.formatter ? column.formatter(item[column.prop]) : item[column.prop] }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, computed, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { darkTheme, applyDarkTheme, getEChartsTheme } from '@/utils/echarts-dark-theme'
import { useTheme } from '@/plugins/theme'
import dayjs from 'dayjs'
import {
  TrendCharts,
  Money,
  ShoppingCart,
  Coin,
  User,
  ShoppingBag,
  Clock,
  Grid,
  Refresh,
  Download,
  ArrowUp,
  ArrowDown
} from '@element-plus/icons-vue'
import {
  getMetrics,
  getSalesTrend,
  getDistributorRank,
  getProductAnalysis,
  getHourDistribution,
  getDetailData,
  exportData as exportStatistics
} from '@/api/statistics'
import EnhancedTable from '@/components/EnhancedTable.vue'
import AnimatedNumber from '@/components/AnimatedNumber.vue'
import LoadingAnimation from '@/components/LoadingAnimation.vue'

// 获取主题钩子
const { isDarkTheme } = useTheme()

// 注册深色主题
applyDarkTheme(echarts)

// 图标映射对象
const iconMap = {
  Money,
  ShoppingCart,
  Coin,
  TrendCharts,
  User,
  ShoppingBag,
  Clock,
  Grid,
  Refresh,
  Download,
  ArrowUp,
  ArrowDown
}

// 颜色映射对象 - 根据主题动态获取颜色
const getColorMap = () => {
  const isDark = isDarkTheme()
  return {
    // 主题色
    '--primary': isDark ? '#fb923c' : '#3b82f6',
    '--primary-light': isDark ? '#fed7aa' : '#60a5fa',
    '--primary-dark': isDark ? '#ea580c' : '#2563eb',
    '--success': isDark ? '#34d399' : '#10b981',
    '--success-light': isDark ? '#6ee7b7' : '#34d399',
    '--success-dark': isDark ? '#10b981' : '#059669',
    '--warning': isDark ? '#fbbf24' : '#f59e0b',
    '--warning-light': isDark ? '#fde68a' : '#fbbf24',
    '--warning-dark': isDark ? '#f59e0b' : '#d97706',
    '--info': isDark ? '#60a5fa' : '#6366f1',
    '--info-light': isDark ? '#93c5fd' : '#818cf8',
    '--info-dark': isDark ? '#3b82f6' : '#4f46e5',
    '--danger': isDark ? '#f87171' : '#ef4444',
    // 边框和文本色
    '--border-color': isDark ? '#374151' : '#e5e7eb',
    '--border-light': isDark ? '#1f2937' : '#f3f4f6',
    '--text-primary': isDark ? '#f3f4f6' : '#111827',
    '--text-secondary': isDark ? '#d1d5db' : '#6b7280',
    '--text-tertiary': isDark ? '#9ca3af' : '#9ca3af'
  }
}

// 获取颜色值的辅助函数
const getColor = (cssVar) => {
  const colorMap = getColorMap()
  return colorMap[cssVar] || cssVar
}

// 响应式数据
const loading = ref(false)
const exportLoading = ref(false)
const filterForm = reactive({})
const dateRange = ref([])

// 图表加载状态
const chartLoading = reactive({
  trend: false,
  distributor: false,
  product: false,
  hour: false
})

// 核心指标配置
const metricsConfig = [
  {
    key: 'revenue',
    title: '总销售额',
    icon: 'Money',
    color: 'var(--success)',
    type: 'revenue',
    valueKey: 'totalRevenue',
    growthKey: 'revenueGrowth',
    formatter: (val) => {
      const num = Number(val) || 0
      return '¥' + num.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    }
  },
  {
    key: 'orders',
    title: '订单总数',
    icon: 'ShoppingCart',
    color: 'var(--primary)',
    type: 'orders',
    valueKey: 'totalOrders',
    growthKey: 'orderGrowth',
    formatter: (val) => {
      const num = Number(val) || 0
      return Math.round(num).toLocaleString('zh-CN')
    }
  },
  {
    key: 'avgValue',
    title: '平均客单价',
    icon: 'Coin',
    color: 'var(--warning)',
    type: 'avg-value',
    valueKey: 'avgOrderValue',
    growthKey: 'avgValueGrowth',
    formatter: (val) => '¥' + (Number(val) || 0).toFixed(2)
  },
  {
    key: 'conversion',
    title: '转化率',
    icon: 'TrendCharts',
    color: 'var(--info)',
    type: 'conversion',
    valueKey: 'conversionRate',
    growthKey: 'conversionGrowth',
    formatter: (val) => ((Number(val) || 0) * 100).toFixed(2) + '%'
  }
]

// 数据状态
const metrics = ref({
  totalRevenue: 0,
  revenueGrowth: 0,
  totalOrders: 0,
  orderGrowth: 0,
  avgOrderValue: 0,
  avgValueGrowth: 0,
  conversionRate: 0,
  conversionGrowth: 0
})

// 图表数据
const salesTrendData = ref([])
const distributorRankData = ref([])
const productAnalysisData = ref([])
const hourDistributionData = ref([])

// 图表配置
const trendType = ref('both')
const distributorRankType = ref('revenue')
const productChartType = ref('pie')

// 表格相关
const tableType = ref('daily')
const tableLoading = ref(false)
const dailyData = ref([])
const distributorData = ref([])
const productData = ref([])
const storeData = ref([])

// 图表实例
let salesTrendChartInstance = null
let distributorRankChartInstance = null
let productAnalysisChartInstance = null
let hourDistributionChartInstance = null

// 计算属性
const queryParams = computed(() => {
  const params = {}
  if (dateRange.value && dateRange.value.length === 2) {
    params.start_date = dateRange.value[0]
    params.end_date = dateRange.value[1]
  }
  return params
})

const currentTableData = computed(() => {
  switch (tableType.value) {
    case 'daily': return dailyData.value
    case 'distributor': return distributorData.value
    case 'product': return productData.value
    case 'store': return storeData.value
    default: return []
  }
})

const currentTableColumns = computed(() => {
  const columnConfigs = {
    daily: [
      { prop: 'date', label: '日期', width: 120, align: 'center' },
      { prop: 'orderCount', label: '订单数', width: 120, align: 'center' },
      { prop: 'revenue', label: '销售额', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'avgOrderValue', label: '客单价', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'newCustomers', label: '新客户数', width: 120, align: 'center' },
      { prop: 'conversionRate', label: '转化率', width: 120, align: 'center', formatter: (val) => formatPercent(val) },
      { prop: 'profit', label: '利润', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'profitRate', label: '利润率', width: 120, align: 'center', formatter: (val) => formatPercent(val) }
    ],
    distributor: [
      { prop: 'distributorName', label: '分销商', minWidth: 180, align: 'left', fixed: 'left' },
      { prop: 'orderCount', label: '订单数', width: 120, align: 'center' },
      { prop: 'revenue', label: '销售额', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'commission', label: '佣金', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'avgOrderValue', label: '客单价', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'conversionRate', label: '转化率', width: 120, align: 'center', formatter: (val) => formatPercent(val) },
      { prop: 'percentage', label: '占比', width: 120, align: 'center', formatter: (val) => (Number(val) || 0).toFixed(2) + '%' },
      { prop: 'status', label: '状态', width: 100, align: 'center', formatter: (val) => val === 1 ? '活跃' : '休眠' }
    ],
    product: [
      { prop: 'productName', label: '商品名称', minWidth: 200, align: 'left', fixed: 'left' },
      { prop: 'category', label: '分类', width: 120, align: 'center' },
      { prop: 'quantity', label: '销量', width: 120, align: 'center' },
      { prop: 'revenue', label: '销售额', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'cost', label: '成本', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'profit', label: '利润', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'avgPrice', label: '均价', width: 120, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'profitRate', label: '利润率', width: 120, align: 'center', formatter: (val) => formatPercent(val) },
      { prop: 'percentage', label: '销售占比', width: 120, align: 'center', formatter: (val) => (Number(val) || 0).toFixed(2) + '%' }
    ],
    store: [
      { prop: 'storeName', label: '门店名称', minWidth: 200, align: 'left', fixed: 'left' },
      { prop: 'storeCode', label: '门店编码', width: 120, align: 'center' },
      { prop: 'city', label: '城市', width: 120, align: 'center' },
      { prop: 'district', label: '区域', width: 120, align: 'center' },
      { prop: 'orderCount', label: '订单数', width: 120, align: 'center' },
      { prop: 'revenue', label: '销售额', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'avgOrderValue', label: '客单价', width: 140, align: 'right', formatter: (val) => formatCurrency(val) },
      { prop: 'activeRate', label: '活跃度', width: 120, align: 'center', formatter: (val) => formatPercent(val) },
      { prop: 'percentage', label: '销售占比', width: 120, align: 'center', formatter: (val) => (Number(val) || 0).toFixed(2) + '%' }
    ]
  }
  return columnConfigs[tableType.value] || []
})

// 日期快捷选择
const selectToday = () => {
  const today = dayjs().format('YYYY-MM-DD')
  dateRange.value = [today, today]
  handleDateChange()
}

const selectYesterday = () => {
  const yesterday = dayjs().subtract(1, 'day').format('YYYY-MM-DD')
  dateRange.value = [yesterday, yesterday]
  handleDateChange()
}

const selectWeek = () => {
  const start = dayjs().startOf('week').format('YYYY-MM-DD')
  const end = dayjs().endOf('week').format('YYYY-MM-DD')
  dateRange.value = [start, end]
  handleDateChange()
}

const selectMonth = () => {
  const start = dayjs().startOf('month').format('YYYY-MM-DD')
  const end = dayjs().endOf('month').format('YYYY-MM-DD')
  dateRange.value = [start, end]
  handleDateChange()
}

// 事件处理
const handleDateChange = () => {
  loadAllData()
}

// 数据加载方法
const loadAllData = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  
  loading.value = true
  try {
    await Promise.all([
      loadMetrics(),
      loadSalesTrend(),
      loadDistributorRank(),
      loadProductAnalysis(),
      loadHourDistribution(),
      loadTableData()
    ])
    
    nextTick(() => {
      if (!isUnmounted) {
        initAllCharts()
      }
    })
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const loadMetrics = async () => {
  try {
    const res = await getMetrics(queryParams.value)
    metrics.value = res.data || {}
  } catch (error) {
    console.error('加载指标失败:', error)
    ElMessage.error('加载指标数据失败，请稍后重试')
    // 初始化为空值，不使用模拟数据
    metrics.value = {
      totalRevenue: 0,
      revenueGrowth: 0,
      totalOrders: 0,
      orderGrowth: 0,
      avgOrderValue: 0,
      avgValueGrowth: 0,
      conversionRate: 0,
      conversionGrowth: 0
    }
  }
}

const loadSalesTrend = async () => {
  chartLoading.trend = true
  try {
    const res = await getSalesTrend({ ...queryParams.value, type: trendType.value })
    salesTrendData.value = res.data || []
    nextTick(() => {
      // 检查组件是否已卸载
      if (!isUnmounted) {
        initSalesTrendChart()
      }
    })
  } catch (error) {
    console.error('加载销售趋势失败:', error)
    ElMessage.error('加载销售趋势数据失败')
    salesTrendData.value = []
    nextTick(() => {
      // 检查组件是否已卸载
      if (!isUnmounted) {
        initSalesTrendChart()
      }
    })
  } finally {
    chartLoading.trend = false
  }
}

const loadDistributorRank = async () => {
  chartLoading.distributor = true
  try {
    console.log('Loading distributor rank with type:', distributorRankType.value)
    const res = await getDistributorRank({ 
      ...queryParams.value, 
      type: distributorRankType.value, 
      limit: 10 
    })
    
    if (res.code === 200 && res.data) {
      distributorRankData.value = res.data
      console.log('Distributor rank data:', res.data)
    } else {
      distributorRankData.value = []
    }
    
    // 确保在下一个渲染周期初始化图表
    await nextTick()
    if (!isUnmounted) {
      initDistributorRankChart()
    }
  } catch (error) {
    console.error('加载分销商排行失败:', error)
    ElMessage.error('加载分销商排行数据失败')
    distributorRankData.value = []
    await nextTick()
    if (!isUnmounted) {
      initDistributorRankChart()
    }
  } finally {
    chartLoading.distributor = false
  }
}

const loadProductAnalysis = async () => {
  chartLoading.product = true
  try {
    const res = await getProductAnalysis({ ...queryParams.value, limit: 10 })
    productAnalysisData.value = res.data || []
    nextTick(() => {
      if (!isUnmounted) {
        initProductAnalysisChart()
      }
    })
  } catch (error) {
    console.error('加载商品分析失败:', error)
    ElMessage.error('加载商品分析数据失败')
    productAnalysisData.value = []
    nextTick(() => {
      if (!isUnmounted) {
        initProductAnalysisChart()
      }
    })
  } finally {
    chartLoading.product = false
  }
}

const loadHourDistribution = async () => {
  chartLoading.hour = true
  try {
    const res = await getHourDistribution(queryParams.value)
    hourDistributionData.value = res.data || []
    nextTick(() => {
      if (!isUnmounted) {
        initHourDistributionChart()
      }
    })
  } catch (error) {
    console.error('加载时段分布失败:', error)
    ElMessage.error('加载时段分布数据失败')
    hourDistributionData.value = []
    nextTick(() => {
      if (!isUnmounted) {
        initHourDistributionChart()
      }
    })
  } finally {
    chartLoading.hour = false
  }
}


const loadTableData = async () => {
  tableLoading.value = true
  try {
    const res = await getDetailData({ ...queryParams.value, type: tableType.value })
    const data = res.data || []
    updateTableData(data)
  } catch (error) {
    console.error('加载详细数据失败:', error)
    ElMessage.error('加载详细数据失败')
    // 重置对应的数据
    updateTableData([])
  } finally {
    tableLoading.value = false
  }
}

// 已移除模拟数据生成函数，现在仅使用真实的 API 数据

// 更新表格数据
const updateTableData = (data) => {
  switch (tableType.value) {
    case 'daily':
      dailyData.value = data
      break
    case 'distributor':
      distributorData.value = data
      break
    case 'product':
      productData.value = data
      break
    case 'store':
      storeData.value = data
      break
  }
}

// 图表初始化
const initAllCharts = () => {
  initSalesTrendChart()
  initDistributorRankChart()
  initProductAnalysisChart()
  initHourDistributionChart()
}

const initSalesTrendChart = () => {
  if (!salesTrendData.value || !salesTrendData.value.length) return
  if (isUnmounted) return
  
  const chartDom = document.querySelector('.trend-chart .chart-container')
  if (!chartDom) return
  
  // 确保销毁旧实例
  if (salesTrendChartInstance) {
    try {
      salesTrendChartInstance.dispose()
    } catch (error) {
      console.warn('销毁旧图表实例失败:', error)
    }
    salesTrendChartInstance = null
  }
  
  try {
    salesTrendChartInstance = echarts.init(chartDom, getEChartsTheme())
  } catch (error) {
    console.error('初始化销售趋势图表失败:', error)
    return
  }
  
  const dates = salesTrendData.value.map(item => item.date)
  const revenues = salesTrendData.value.map(item => item.revenue)
  const orders = salesTrendData.value.map(item => item.orders)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: getColor('--border-color'),
      textStyle: {
        color: getColor('--text-primary')
      },
      axisPointer: {
        type: 'cross',
        crossStyle: {
          color: getColor('--primary')
        }
      }
    },
    legend: {
      data: ['销售额', '订单数'],
      textStyle: {
        color: getColor('--text-primary')
      }
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisPointer: {
        type: 'shadow'
      },
      axisLine: {
        lineStyle: {
          color: getColor('--border-color')
        }
      },
      axisLabel: {
        color: getColor('--text-secondary')
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '销售额(元)',
        position: 'left',
        axisLine: {
          show: true,
          lineStyle: {
            color: getColor('--success')
          }
        },
        axisLabel: {
          formatter: '¥{value}',
          color: getColor('--text-secondary')
        },
        splitLine: {
          lineStyle: {
            color: getColor('--border-light')
          }
        }
      },
      {
        type: 'value',
        name: '订单数',
        position: 'right',
        axisLine: {
          show: true,
          lineStyle: {
            color: getColor('--primary')
          }
        },
        axisLabel: {
          formatter: '{value}单',
          color: getColor('--text-secondary')
        }
      }
    ],
    series: [
      {
        name: '销售额',
        type: 'bar',
        data: revenues,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: getColor('--success-light') },
            { offset: 1, color: getColor('--success') }
          ])
        },
        emphasis: {
          itemStyle: {
            color: getColor('--success-dark')
          }
        }
      },
      {
        name: '订单数',
        type: 'line',
        yAxisIndex: 1,
        data: orders,
        smooth: true,
        lineStyle: {
          color: getColor('--primary'),
          width: 3
        },
        itemStyle: {
          color: getColor('--primary')
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
            { offset: 1, color: 'rgba(59, 130, 246, 0.05)' }
          ])
        }
      }
    ],
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    }
  }
  
  try {
    salesTrendChartInstance.setOption(option)
  } catch (error) {
    console.error('设置销售趋势图表选项失败:', error)
  }
}

const initDistributorRankChart = () => {
  if (!distributorRankData.value || !distributorRankData.value.length) return
  if (isUnmounted) return
  
  const chartDom = document.querySelector('.distributor-rank .chart-container')
  if (!chartDom) return
  
  // 确保销毁旧实例
  if (distributorRankChartInstance) {
    try {
      distributorRankChartInstance.dispose()
    } catch (error) {
      console.warn('销毁旧图表实例失败:', error)
    }
    distributorRankChartInstance = null
  }
  
  try {
    distributorRankChartInstance = echarts.init(chartDom, getEChartsTheme())
  } catch (error) {
    console.error('初始化分销商排行图表失败:', error)
    return
  }
  
  const names = distributorRankData.value.map(item => item.distributorName || '未知分销商')
  const values = distributorRankData.value.map(item => 
    distributorRankType.value === 'revenue' ? (item.revenue || 0) : (item.orders || 0)
  )
  
  console.log('Chart data - names:', names)
  console.log('Chart data - values:', values)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: getColor('--border-color'),
      formatter: function(params) {
        const value = distributorRankType.value === 'revenue' 
          ? `¥${params[0].value.toLocaleString()}` 
          : `${params[0].value}单`
        return `${params[0].name}<br/>${params[0].seriesName}: ${value}`
      }
    },
    xAxis: {
      type: 'value',
      axisLine: {
        lineStyle: {
          color: getColor('--border-color')
        }
      },
      axisLabel: {
        color: getColor('--text-secondary'),
        formatter: distributorRankType.value === 'revenue' ? '¥{value}' : '{value}单'
      },
      splitLine: {
        lineStyle: {
          color: getColor('--border-light')
        }
      }
    },
    yAxis: {
      type: 'category',
      data: names,
      axisLine: {
        lineStyle: {
          color: getColor('--border-color')
        }
      },
      axisLabel: {
        color: getColor('--text-secondary')
      }
    },
    series: [
      {
        name: distributorRankType.value === 'revenue' ? '销售额' : '订单数',
        type: 'bar',
        data: values,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: getColor('--primary-light') },
            { offset: 1, color: getColor('--primary') }
          ])
        },
        emphasis: {
          itemStyle: {
            color: getColor('--primary-dark')
          }
        }
      }
    ],
    grid: {
      left: '20%',
      right: '4%',
      bottom: '3%',
      top: '10%'
    }
  }
  
  try {
    distributorRankChartInstance.setOption(option)
  } catch (error) {
    console.error('设置分销商排行图表选项失败:', error)
  }
}

const initProductAnalysisChart = () => {
  if (!productAnalysisData.value || !productAnalysisData.value.length) return
  if (isUnmounted) return
  
  const chartDom = document.querySelector('.product-analysis .chart-container')
  if (!chartDom) return
  
  // 确保销毁旧实例
  if (productAnalysisChartInstance) {
    try {
      productAnalysisChartInstance.dispose()
    } catch (error) {
      console.warn('销毁旧图表实例失败:', error)
    }
    productAnalysisChartInstance = null
  }
  
  try {
    productAnalysisChartInstance = echarts.init(chartDom, getEChartsTheme())
  } catch (error) {
    console.error('初始化商品分析图表失败:', error)
    return
  }
  
  let option
  
  if (productChartType.value === 'pie') {
    const data = productAnalysisData.value.map((item, index) => ({
      name: item.productName,
      value: item.revenue,
      itemStyle: {
        color: `hsl(${(index * 45) % 360}, 70%, 60%)`
      }
    }))
    
    option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: ¥{c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        textStyle: {
          color: getColor('--text-primary')
        }
      },
      series: [
        {
          name: '商品销售额',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 5,
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
              fontSize: 16,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: false
          },
          data: data
        }
      ]
    }
  } else {
    const names = productAnalysisData.value.map(item => item.productName)
    const revenues = productAnalysisData.value.map(item => item.revenue)
    
    option = {
      tooltip: {
        trigger: 'axis',
        formatter: '{b}: ¥{c}'
      },
      xAxis: {
        type: 'category',
        data: names,
        axisLabel: {
          color: getColor('--text-secondary'),
          rotate: 45
        }
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          color: getColor('--text-secondary'),
          formatter: '¥{value}'
        }
      },
      series: [
        {
          type: 'bar',
          data: revenues,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: getColor('--warning-light') },
              { offset: 1, color: getColor('--warning') }
            ])
          }
        }
      ]
    }
  }
  
  try {
    productAnalysisChartInstance.setOption(option)
  } catch (error) {
    console.error('设置商品分析图表选项失败:', error)
  }
}

const initHourDistributionChart = () => {
  if (!hourDistributionData.value || !hourDistributionData.value.length) return
  if (isUnmounted) return
  
  const chartDom = document.querySelector('.hour-analysis .chart-container')
  if (!chartDom) return
  
  // 确保销毁旧实例
  if (hourDistributionChartInstance) {
    try {
      hourDistributionChartInstance.dispose()
    } catch (error) {
      console.warn('销毁旧图表实例失败:', error)
    }
    hourDistributionChartInstance = null
  }
  
  try {
    hourDistributionChartInstance = echarts.init(chartDom, getEChartsTheme())
  } catch (error) {
    console.error('初始化时段分布图表失败:', error)
    return
  }
  
  const hours = hourDistributionData.value.map(item => `${item.hour}:00`)
  const orders = hourDistributionData.value.map(item => item.orders)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: '{b}: {c}单'
    },
    xAxis: {
      type: 'category',
      data: hours,
      axisLabel: {
        color: getColor('--text-secondary')
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: getColor('--text-secondary'),
        formatter: '{value}单'
      }
    },
    series: [
      {
        type: 'line',
        data: orders,
        smooth: true,
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: getColor('--info-light') },
            { offset: 1, color: 'transparent' }
          ])
        },
        lineStyle: {
          color: getColor('--info'),
          width: 3
        },
        itemStyle: {
          color: getColor('--info')
        }
      }
    ],
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    }
  }
  
  try {
    hourDistributionChartInstance.setOption(option)
  } catch (error) {
    console.error('设置时段分布图表选项失败:', error)
  }
}


// 导出数据
const exportData = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  
  exportLoading.value = true
  try {
    const res = await exportStatistics({ ...queryParams.value, format: 'excel' })
    const blob = new Blob([res], { 
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' 
    })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `统计报表_${dateRange.value[0]}_${dateRange.value[1]}.xlsx`
    link.click()
    URL.revokeObjectURL(link.href)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  } finally {
    exportLoading.value = false
  }
}

// 格式化函数
const formatCurrency = (num) => {
  const value = Number(num) || 0
  return `¥${value.toFixed(2)}`
}

const formatPercent = (num) => {
  const value = Number(num) || 0
  return `${(value * 100).toFixed(2)}%`
}

const formatGrowth = (growth) => {
  const value = Number(growth) || 0
  const sign = value > 0 ? '+' : ''
  return `${sign}${(value * 100).toFixed(2)}%`
}

const getTrendClass = (growth) => {
  if (!growth) return 'trend-neutral'
  return growth > 0 ? 'trend-positive' : 'trend-negative'
}

// 移动端卡片辅助方法
const getCardTitle = (item) => {
  switch (tableType.value) {
    case 'daily':
      return item.date
    case 'distributor':
      return item.distributorName
    case 'product':
      return item.productName
    case 'store':
      return item.storeName
    default:
      return '数据项'
  }
}

const getCardStatus = (item) => {
  switch (tableType.value) {
    case 'daily':
      return `${item.orderCount}单`
    case 'distributor':
      return `${(Number(item.percentage) || 0).toFixed(1)}%`
    case 'product':
      return item.category
    case 'store':
      return item.city
    default:
      return ''
  }
}

// 监听图表类型变化
watch(productChartType, () => {
  nextTick(() => {
    if (!isUnmounted) {
      initProductAnalysisChart()
    }
  })
})

// 添加组件是否已卸载的标志
let isUnmounted = false

// 统一的 resize 处理函数（使用防抖）
let resizeTimer = null
const handleResize = () => {
  // 如果组件已卸载，直接返回
  if (isUnmounted) return
  
  if (resizeTimer) {
    clearTimeout(resizeTimer)
  }
  resizeTimer = setTimeout(() => {
    // 再次检查组件是否已卸载
    if (isUnmounted) return
    
    // 只对已初始化的图表进行 resize，并确保实例存在
    try {
      // 安全地检查和调整图表大小
      const safeResize = (chartInstance) => {
        if (chartInstance && typeof chartInstance.resize === 'function') {
          try {
            // 不使用 isDisposed() 方法，直接捕获 resize 的异常
            chartInstance.resize()
          } catch (e) {
            // 忽略已销毁图表的 resize 错误
          }
        }
      }
      
      safeResize(salesTrendChartInstance)
      safeResize(distributorRankChartInstance)
      safeResize(productAnalysisChartInstance)
      safeResize(hourDistributionChartInstance)
    } catch (error) {
      console.warn('Resize chart error:', error)
    }
  }, 300)
}

// 销毁所有图表实例
const destroyAllCharts = () => {
  // 安全地销毁图表实例
  const safeDispose = (chartInstance) => {
    if (chartInstance) {
      try {
        // 检查 dispose 方法是否存在
        if (typeof chartInstance.dispose === 'function') {
          chartInstance.dispose()
        }
      } catch (error) {
        // 忽略销毁时的错误，可能图表已经被销毁
        console.warn('Chart dispose error (ignored):', error)
      }
    }
    return null
  }
  
  // 销毁所有图表实例
  salesTrendChartInstance = safeDispose(salesTrendChartInstance)
  distributorRankChartInstance = safeDispose(distributorRankChartInstance)
  productAnalysisChartInstance = safeDispose(productAnalysisChartInstance)
  hourDistributionChartInstance = safeDispose(hourDistributionChartInstance)
}

// 组件挂载
onMounted(() => {
  isUnmounted = false
  selectMonth() // 默认选择本月
  
  // 添加全局 resize 监听器
  window.addEventListener('resize', handleResize)
})

// 监听分销商排行类型变化
watch(distributorRankType, (newVal) => {
  console.log('Distributor rank type changed to:', newVal)
  loadDistributorRank()
})

// 监听主题变化
watch(() => isDarkTheme(), () => {
  if (isUnmounted) return
  
  // 重新初始化所有图表以应用新主题
  nextTick(() => {
    if (isUnmounted) return
    
    // 销毁旧图表实例
    destroyAllCharts()
    
    // 重新初始化所有图表
    initAllCharts()
  })
})

// 组件卸载前清理
onBeforeUnmount(() => {
  // 设置卸载标志
  isUnmounted = true
  
  // 清理定时器
  if (resizeTimer) {
    clearTimeout(resizeTimer)
    resizeTimer = null
  }
  
  // 移除事件监听器
  window.removeEventListener('resize', handleResize)
  
  // 销毁所有图表实例
  destroyAllCharts()
})
</script>

<style lang="scss" scoped>
.statistics-enhanced-container {
  padding: var(--spacing-6);
  background: var(--bg-page);
  min-height: 100vh;
  
  .filter-section {
    margin-bottom: var(--spacing-6);
    
    .filter-card {
      border: none;
      box-shadow: var(--shadow-sm);
      border-radius: var(--radius-lg);
      overflow: hidden;
      
      .filter-content {
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-wrap: wrap;
        gap: var(--spacing-4);
        
        .filter-left {
          .filter-title {
            margin: 0;
            font-size: 1.5rem;
            font-weight: 600;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: var(--spacing-2);
            
            .el-icon {
              color: var(--primary);
            }
          }
        }
        
        .filter-right {
          .el-form {
            margin: 0;
            
            .el-form-item {
              margin-right: var(--spacing-4);
              margin-bottom: 0;
              
              &:last-child {
                margin-right: 0;
              }
            }
          }
          
          .quick-select-group {
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
            }
          }
        }
      }
    }
  }
  
  .metrics-section {
    margin-bottom: var(--spacing-6);
    
    .metrics-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: var(--spacing-4);
      
      .metric-card-wrapper {
        animation: slideInUp 0.6s ease-out forwards;
        animation-delay: var(--delay, 0s);
        opacity: 0;
        transform: translateY(20px);
        
        .metric-card {
          height: 100%;
          border: none;
          border-radius: var(--radius-lg);
          transition: all 0.3s ease;
          background: linear-gradient(135deg, var(--bg-white) 0%, var(--bg-gray-50) 100%);
          
          &:hover {
            transform: translateY(-4px);
            box-shadow: var(--shadow-lg);
          }
          
          &.metric-revenue {
            border-left: 4px solid var(--success);
          }
          
          &.metric-orders {
            border-left: 4px solid var(--primary);
          }
          
          &.metric-avg-value {
            border-left: 4px solid var(--warning);
          }
          
          &.metric-conversion {
            border-left: 4px solid var(--info);
          }
          
          .metric-header {
            display: flex;
            align-items: flex-start;
            gap: var(--spacing-3);
            margin-bottom: var(--spacing-4);
            
            .metric-icon-wrapper {
              width: 48px;
              height: 48px;
              border-radius: var(--radius-lg);
              background: var(--bg-gray-100);
              display: flex;
              align-items: center;
              justify-content: center;
              
              .metric-icon {
                font-size: 1.5rem;
              }
            }
            
            .metric-info {
              flex: 1;
              
              .metric-title {
                margin: 0 0 var(--spacing-1);
                font-size: 0.875rem;
                color: var(--text-secondary);
                font-weight: 500;
              }
              
              .metric-value {
                font-size: 1.875rem;
                font-weight: 700;
                color: var(--text-primary);
                line-height: 1.2;
              }
            }
          }
          
          .metric-footer {
            .metric-trend {
              display: flex;
              align-items: center;
              justify-content: space-between;
              
              .trend-label {
                font-size: 0.75rem;
                color: var(--text-tertiary);
              }
              
              .trend-value {
                display: flex;
                align-items: center;
                gap: var(--spacing-1);
                font-size: 0.875rem;
                font-weight: 600;
                
                .trend-icon {
                  font-size: 0.75rem;
                }
                
                &.trend-positive {
                  color: var(--success);
                }
                
                &.trend-negative {
                  color: var(--danger);
                }
                
                &.trend-neutral {
                  color: var(--text-tertiary);
                }
              }
            }
          }
        }
      }
    }
  }
  
  .charts-section {
    margin-bottom: var(--spacing-6);
    
    .trend-chart {
      border: none;
      border-radius: var(--radius-lg);
      box-shadow: var(--shadow-sm);
      
      &:hover {
        box-shadow: var(--shadow-md);
      }
    }
  }
  
  .analytics-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: var(--spacing-4);
    margin-bottom: var(--spacing-6);
  }
  
  .chart-card {
    border: none;
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    transition: all 0.3s ease;
    
    &:hover {
      box-shadow: var(--shadow-md);
      transform: translateY(-2px);
    }
    
    .chart-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      
      .chart-title {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        font-size: 1.125rem;
        font-weight: 600;
        color: var(--text-primary);
        
        .el-icon {
          color: var(--primary);
        }
      }
      
      .chart-controls {
        .el-radio-group,
        .el-button-group {
          .el-radio-button,
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
          }
        }
      }
    }
    
    .chart-content {
      position: relative;
      
      .chart-container {
        width: 100%;
      }
    }
  }
  
  .data-table-section {
    border: none;
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    
    .table-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      
      .table-title {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        font-size: 1.125rem;
        font-weight: 600;
        color: var(--text-primary);
        
        .el-icon {
          color: var(--primary);
        }
      }
      
      .table-controls {
        .table-type-select {
          min-width: 120px;
        }
      }
    }
    
    .table-content {
      margin-top: var(--spacing-4);
      
      // 确保 EnhancedTable 组件正确显示
      :deep(.el-table) {
        width: 100%;
        
        .el-table__header-wrapper {
          .el-table__header {
            width: 100% !important;
          }
        }
        
        .el-table__body-wrapper {
          .el-table__body {
            width: 100% !important;
          }
        }
        
        // 修复列宽度问题
        .el-table__column-resize-proxy {
          display: none;
        }
      }
    }
  }
}

// 动画定义
@keyframes slideInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.metric-card-enter-active {
  transition: all 0.5s ease;
}

.metric-card-enter-from {
  opacity: 0;
  transform: scale(0.8) translateY(20px);
}

// 响应式设计
@media (max-width: 768px) {
  .statistics-enhanced-container {
    padding: var(--spacing-4);
    
    .filter-section .filter-card .filter-content {
      flex-direction: column;
      align-items: flex-start;
      
      .filter-right {
        width: 100%;
        
        .el-form {
          flex-direction: column;
          align-items: stretch;
          
          .el-form-item {
            margin-right: 0;
            margin-bottom: var(--spacing-3);
          }
        }
      }
    }
    
    .metrics-grid {
      grid-template-columns: 1fr;
    }
    
    .analytics-grid {
      grid-template-columns: 1fr;
    }
    
    // 移动端表格卡片化
    .mobile-table-cards {
      .mobile-card-list {
        .mobile-card-item {
          background: var(--bg-primary);
          border: 1px solid var(--border-default);
          border-radius: var(--mobile-border-radius);
          padding: var(--mobile-space-4);
          margin-bottom: var(--mobile-space-3);
          
          &:last-child {
            margin-bottom: 0;
          }
          
          .card-item-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: var(--mobile-space-3);
            
            .card-item-title {
              font-size: var(--mobile-text-base);
              font-weight: 600;
              color: var(--text-primary);
              margin: 0;
              flex: 1;
              margin-right: var(--mobile-space-2);
            }
            
            .card-item-status {
              font-size: var(--mobile-text-sm);
              color: var(--text-secondary);
              background: var(--bg-secondary);
              padding: var(--mobile-space-1) var(--mobile-space-2);
              border-radius: var(--radius-md);
              font-weight: 500;
            }
          }
          
          .card-item-content {
            .card-item-row {
              display: flex;
              justify-content: space-between;
              align-items: center;
              padding: var(--mobile-space-2) 0;
              border-bottom: 1px solid var(--border-light);
              
              &:last-child {
                border-bottom: none;
                padding-bottom: 0;
              }
              
              .card-item-label {
                font-size: var(--mobile-text-sm);
                color: var(--text-secondary);
                font-weight: 500;
              }
              
              .card-item-value {
                font-size: var(--mobile-text-sm);
                color: var(--text-primary);
                font-weight: 600;
                text-align: right;
              }
            }
          }
        }
      }
    }
  }
}

// 深度样式修复 - 确保表格正常显示
:deep(.el-table) {
  width: 100% !important;
  
  .el-table__header-wrapper,
  .el-table__body-wrapper,
  .el-table__footer-wrapper {
    width: 100% !important;
  }
  
  .el-table__header,
  .el-table__body,
  .el-table__footer {
    width: 100% !important;
  }
  
  // 修复列宽度问题
  .el-table__body-wrapper {
    overflow-x: auto;
  }
  
  // 确保固定列正常工作
  .el-table__fixed,
  .el-table__fixed-right {
    height: 100% !important;
    
    .el-table__fixed-body-wrapper {
      top: 37px !important;
    }
  }
  
  // 隐藏列调整代理
  .el-table__column-resize-proxy {
    display: none !important;
  }
  
  // 对齐样式
  .el-table__cell {
    &[class*="text-center"] .cell {
      text-align: center;
    }
    
    &[class*="text-right"] .cell {
      text-align: right;
      justify-content: flex-end;
    }
    
    &[class*="text-left"] .cell {
      text-align: left;
      justify-content: flex-start;
    }
  }
  
  // 数值列样式
  .cell {
    &:has(.el-table__cell[class*="text-right"]) {
      font-variant-numeric: tabular-nums;
    }
  }
}

// EnhancedTable 组件样式覆盖
:deep(.enhanced-table-container) {
  .el-table {
    font-size: 14px;
    
    .el-table__header {
      th {
        background-color: #f5f7fa;
        color: #303133;
        font-weight: 600;
      }
    }
    
    .el-table__body {
      td {
        padding: 12px 0;
      }
    }
  }
}
</style>