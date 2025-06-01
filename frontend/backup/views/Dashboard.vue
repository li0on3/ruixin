<template>
  <div class="dashboard page-container">
    <!-- 统计卡片 -->
    <el-row :gutter="24" class="stat-cards">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card metric-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #409eff">
              <el-icon><ShoppingCart /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.todayOrders || 0 }}</div>
              <div class="stat-label">今日订单</div>
            </div>
          </div>
          <div class="stat-footer">
            <span>较昨日</span>
            <span :class="statistics.orderGrowth >= 0 ? 'text-success' : 'text-danger'">
              {{ statistics.orderGrowth >= 0 ? '+' : '' }}{{ statistics.orderGrowth || 0 }}%
            </span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card metric-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">¥{{ statistics.todayAmount || 0 }}</div>
              <div class="stat-label">今日销售额</div>
            </div>
          </div>
          <div class="stat-footer">
            <span>较昨日</span>
            <span :class="statistics.amountGrowth >= 0 ? 'text-success' : 'text-danger'">
              {{ statistics.amountGrowth >= 0 ? '+' : '' }}{{ statistics.amountGrowth || 0 }}%
            </span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card metric-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.activeDistributors || 0 }}</div>
              <div class="stat-label">活跃分销商</div>
            </div>
          </div>
          <div class="stat-footer">
            <span>总分销商</span>
            <span>{{ statistics.totalDistributors || 0 }}</span>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card class="stat-card metric-card">
          <div class="stat-content">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon><CreditCard /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.activeCards || 0 }}</div>
              <div class="stat-label">可用卡片</div>
            </div>
          </div>
          <div class="stat-footer">
            <span>总卡片</span>
            <span>{{ statistics.totalCards || 0 }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 图表区域 -->
    <el-row :gutter="24" class="chart-area">
      <el-col :xs="24" :lg="16">
        <el-card class="table-card">
          <template #header>
            <div class="card-header">
              <span>订单趋势</span>
              <el-radio-group v-model="dateRange" size="small" @change="loadOrderTrend">
                <el-radio-button label="week">近7天</el-radio-button>
                <el-radio-button label="month">近30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="orderChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="8">
        <el-card class="table-card">
          <template #header>
            <span>热门商品</span>
          </template>
          <div ref="goodsChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 最新订单 -->
    <el-card class="recent-orders table-card">
      <template #header>
        <div class="card-header">
          <span>最新订单</span>
          <el-button type="primary" text @click="$router.push('/orders')">
            查看全部
          </el-button>
        </div>
      </template>
      
      <el-table :data="recentOrders" stripe>
        <el-table-column prop="orderNo" label="订单号" width="200" />
        <el-table-column prop="distributorName" label="分销商" />
        <el-table-column prop="storeName" label="门店" />
        <el-table-column prop="totalAmount" label="金额">
          <template #default="{ row }">
            ¥{{ row.totalAmount }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="下单时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getDashboardStatistics, getOrderTrend, getHotGoods, getRecentOrders } from '@/api/dashboard'
import { formatDate } from '@/utils/date'

const statistics = ref({})
const dateRange = ref('week')
const recentOrders = ref([])

const orderChartRef = ref()
const goodsChartRef = ref()
let orderChart = null
let goodsChart = null

// 获取统计数据
const loadStatistics = async () => {
  try {
    const response = await getDashboardStatistics()
    console.log('Dashboard response:', response) // 调试完整响应
    console.log('Dashboard statistics data:', response.data) // 调试数据
    statistics.value = response.data || response // 兼容不同的响应格式
  } catch (error) {
    console.error('Failed to load statistics:', error)
  }
}

// 获取订单趋势
const loadOrderTrend = async () => {
  try {
    const { data } = await getOrderTrend(dateRange.value)
    
    // 检查数据是否存在
    if (!data || !Array.isArray(data)) {
      console.warn('No order trend data available')
      return
    }
    
    const option = {
      tooltip: {
        trigger: 'axis'
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
        data: data.map(item => item.date)
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '订单数',
          type: 'line',
          smooth: true,
          data: data.map(item => item.total_orders || 0),
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0)' }
            ])
          },
          itemStyle: {
            color: '#409eff'
          }
        }
      ]
    }
    
    orderChart.setOption(option)
  } catch (error) {
    console.error('Failed to load order trend:', error)
  }
}

// 获取热门商品
const loadHotGoods = async () => {
  try {
    const { data } = await getHotGoods()
    
    // 检查数据是否存在
    if (!data || !Array.isArray(data)) {
      console.warn('No hot goods data available')
      return
    }
    
    const option = {
      tooltip: {
        trigger: 'item'
      },
      series: [
        {
          type: 'pie',
          radius: ['40%', '70%'],
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
          data: data.map(item => ({
            name: item.name,
            value: item.count
          }))
        }
      ]
    }
    
    goodsChart.setOption(option)
  } catch (error) {
    console.error('Failed to load hot goods:', error)
  }
}

// 获取最新订单
const loadRecentOrders = async () => {
  try {
    const { data } = await getRecentOrders()
    recentOrders.value = data.map(order => ({
      ...order,
      createdAt: formatDate(order.createdAt)
    }))
  } catch (error) {
    console.error('Failed to load recent orders:', error)
  }
}

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const types = {
    0: 'info',
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'info',
    5: 'info'
  }
  return types[status] || 'info'
}

// 获取订单状态文本
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

onMounted(() => {
  // 初始化图表
  orderChart = echarts.init(orderChartRef.value)
  goodsChart = echarts.init(goodsChartRef.value)
  
  // 自适应
  window.addEventListener('resize', () => {
    orderChart?.resize()
    goodsChart?.resize()
  })
  
  // 加载数据
  loadStatistics()
  loadOrderTrend()
  loadHotGoods()
  loadRecentOrders()
})

onUnmounted(() => {
  orderChart?.dispose()
  goodsChart?.dispose()
})
</script>

<style lang="scss" scoped>
.dashboard {
  width: 100%;
  height: 100%;
  min-width: 1280px; // 设置最小宽度
  
  .stat-cards {
    margin-bottom: 24px;
    
    .stat-card {
      height: 100%;
      
      .stat-content {
        display: flex;
        align-items: center;
        margin-bottom: 24px;
        
        .stat-icon {
          width: 80px;
          height: 80px;
          border-radius: 20px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
          font-size: 40px;
          margin-right: 24px;
          box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
          background-image: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, rgba(255, 255, 255, 0) 100%);
          flex-shrink: 0;
          
          @media (min-width: 1440px) {
            width: 90px;
            height: 90px;
            font-size: 44px;
            margin-right: 28px;
          }
          
          @media (min-width: 1920px) {
            width: 100px;
            height: 100px;
            font-size: 48px;
            margin-right: 32px;
          }
        }
        
        .stat-info {
          flex: 1;
          min-width: 0; // 防止内容溢出
          
          .stat-value {
            font-size: 36px;
            font-weight: 900;
            color: #1e293b;
            line-height: 1.1;
            margin-bottom: 8px;
            letter-spacing: -1px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            
            @media (min-width: 1440px) {
              font-size: 42px;
              margin-bottom: 10px;
            }
            
            @media (min-width: 1920px) {
              font-size: 48px;
              margin-bottom: 12px;
            }
          }
          
          .stat-label {
            font-size: 16px;
            color: #64748b;
            font-weight: 600;
            letter-spacing: 0.5px;
            
            @media (min-width: 1440px) {
              font-size: 18px;
            }
            
            @media (min-width: 1920px) {
              font-size: 20px;
            }
          }
        }
      }
      
      .stat-footer {
        display: flex;
        justify-content: space-between;
        font-size: 16px;
        color: #64748b;
        padding-top: 16px;
        border-top: 1px solid #f0f2f5;
        font-weight: 600;
        
        @media (min-width: 1920px) {
          font-size: 18px;
          padding-top: 20px;
        }
        
        .text-success {
          color: #10b981;
          font-weight: 700;
          font-size: 17px;
          
          @media (min-width: 1920px) {
            font-size: 19px;
          }
        }
        
        .text-danger {
          color: #ef4444;
          font-weight: 700;
          font-size: 17px;
          
          @media (min-width: 1920px) {
            font-size: 19px;
          }
        }
      }
    }
  }
  
  .chart-area {
    margin-bottom: 32px;
    
    .el-col {
      margin-bottom: 24px;
      
      @media (min-width: 1440px) {
        margin-bottom: 0;
      }
    }
  }
  
  .chart-container {
    height: 400px;
    min-height: 300px;
    
    @media (min-width: 1440px) {
      height: 450px;
    }
    
    @media (min-width: 1920px) {
      height: 500px;
    }
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 18px;
    font-weight: 700;
    flex-wrap: wrap;
    gap: 16px;
    
    @media (min-width: 1920px) {
      font-size: 20px;
    }
  }
  
  .recent-orders {
    margin-top: 32px;
    
    .el-table {
      // 优化表格列宽
      .el-table__header-wrapper {
        .el-table__header {
          th {
            &:nth-child(1) { width: 200px; } // 订单号
            &:nth-child(2) { min-width: 150px; } // 分销商
            &:nth-child(3) { min-width: 200px; } // 门店
            &:nth-child(4) { width: 120px; } // 金额
            &:nth-child(5) { width: 100px; } // 状态
            &:nth-child(6) { width: 180px; } // 时间
          }
        }
      }
    }
  }
}

// 修复响应式网格在不同屏幕下的问题
.el-row {
  // 在小于1440px的屏幕上，让卡片每行显示2个
  @media (min-width: 1280px) and (max-width: 1440px) {
    &.stat-cards {
      .el-col-lg-6 {
        max-width: 50%;
        flex: 0 0 50%;
      }
    }
  }
}
</style>