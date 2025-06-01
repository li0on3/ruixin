<template>
  <div class="enhanced-orders-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">订单管理</h1>
      <div class="page-actions">
        <el-button type="primary" :icon="Refresh" @click="loadData">
          刷新数据
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="statistics-section">
      <el-row :gutter="24">
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stat-card">
            <div class="stat-icon total">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.total_orders || 0 }}</div>
              <div class="stat-label">总订单数</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stat-card">
            <div class="stat-icon success">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ statistics.success_orders || 0 }}</div>
              <div class="stat-label">成功订单</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stat-card">
            <div class="stat-icon amount">
              <el-icon><Wallet /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">¥{{ (statistics.total_amount || 0).toFixed(2) }}</div>
              <div class="stat-label">总金额</div>
            </div>
          </div>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <div class="stat-card">
            <div class="stat-icon rate">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ (statistics.success_rate || 0).toFixed(1) }}%</div>
              <div class="stat-label">成功率</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 高级搜索 -->
    <el-card class="search-card" v-show="showSearch">
      <el-form :model="searchForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item label="订单号">
              <el-input
                v-model="searchForm.orderNo"
                placeholder="请输入订单号"
                clearable
                @keyup.enter="handleSearch"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item label="分销商">
              <el-select 
                v-model="searchForm.distributorId" 
                placeholder="请选择分销商" 
                clearable 
                filterable
                :loading="distributorsLoading"
              >
                <el-option label="全部" value="" />
                <el-option
                  v-for="distributor in distributors"
                  :key="distributor.id"
                  :label="distributor.company_name || distributor.contact_name"
                  :value="distributor.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item label="订单状态">
              <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
                <el-option label="全部" value="" />
                <el-option label="待处理" :value="0" />
                <el-option label="处理中" :value="1" />
                <el-option label="已完成" :value="2" />
                <el-option label="失败" :value="3" />
                <el-option label="已退款" :value="4" />
                <el-option label="已取消" :value="5" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item label="日期范围">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                :shortcuts="dateShortcuts"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item label="手机号">
              <el-input
                v-model="searchForm.phoneNumber"
                placeholder="请输入手机号"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8">
            <el-form-item label="取餐码">
              <el-input
                v-model="searchForm.takeCode"
                placeholder="请输入取餐码"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24" style="text-align: right;">
            <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
            <el-button :icon="Refresh" @click="handleReset">重置</el-button>
            <el-button text @click="showSearch = false">收起</el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <EnhancedTable
        ref="tableRef"
        :data="tableData"
        :columns="tableColumns"
        :loading="loading"
        :actions="tableActions"
        :pagination="false"
        :searchable="false"
        :show-toolbar="true"
        :stripe="true"
        @refresh="loadData"
        @export="handleExport"
      >
        <!-- 自定义搜索栏左侧内容 -->
        <template #toolbar-left>
          <el-button 
            v-if="!showSearch" 
            type="primary" 
            text 
            :icon="Search" 
            @click="showSearch = true"
          >
            高级搜索
          </el-button>
        </template>

        <!-- 订单号列 -->
        <template #orderNo="{ row }">
          <el-link type="primary" @click="handleViewDetail(row)">
            {{ row.order_no }}
          </el-link>
        </template>

        <!-- 状态列 -->
        <template #status="{ row }">
          <el-tag :type="getOrderStatusType(row.status)" effect="light">
            {{ getOrderStatusText(row.status) }}
          </el-tag>
        </template>

        <!-- 金额列 -->
        <template #amount="{ row }">
          <div class="amount-cell">
            <div class="main-amount">¥{{ row.total_amount.toFixed(2) }}</div>
            <div class="sub-amount" v-if="row.profit_amount">
              利润: ¥{{ row.profit_amount.toFixed(2) }}
            </div>
          </div>
        </template>

        <!-- 商品信息列 -->
        <template #goods="{ row }">
          <div class="goods-info">
            <div v-for="(item, index) in row.goods" :key="index" class="goods-item">
              <span class="goods-name">{{ item.goods_name }}</span>
              <span class="goods-sku">{{ item.sku_name }}</span>
              <span class="goods-qty">×{{ item.quantity }}</span>
            </div>
          </div>
        </template>

        <!-- 门店信息列 -->
        <template #store="{ row }">
          <div class="store-info">
            <div class="store-name">{{ row.store_name }}</div>
            <div class="store-address">{{ row.store_address }}</div>
          </div>
        </template>

        <!-- 时间列 -->
        <template #createdAt="{ row }">
          <div class="time-info">
            <div>{{ formatDate(row.created_at) }}</div>
            <div class="time-ago">{{ getTimeAgo(row.created_at) }}</div>
          </div>
        </template>
      </EnhancedTable>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 订单详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="订单详情"
      width="900px"
      class="order-detail-dialog"
    >
      <div v-if="currentOrder" class="order-detail-content">
        <!-- 基本信息 -->
        <div class="detail-section">
          <h3 class="section-title">基本信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单号">
              <span class="mono-text">{{ currentOrder.order_no }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getOrderStatusType(currentOrder.status)">
                {{ getOrderStatusText(currentOrder.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="分销商">
              {{ currentOrder.distributor_name || currentOrder.distributor_id }}
            </el-descriptions-item>
            <el-descriptions-item label="手机号">
              {{ currentOrder.phone_number }}
            </el-descriptions-item>
            <el-descriptions-item label="取餐码">
              <span class="take-code">{{ currentOrder.take_code || '-' }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="下单时间">
              {{ formatDate(currentOrder.created_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 门店信息 -->
        <div class="detail-section">
          <h3 class="section-title">门店信息</h3>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="门店名称">
              {{ currentOrder.store_name }}
            </el-descriptions-item>
            <el-descriptions-item label="门店地址">
              {{ currentOrder.store_address }}
            </el-descriptions-item>
            <el-descriptions-item label="门店ID">
              {{ currentOrder.store_id }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 商品信息 -->
        <div class="detail-section">
          <h3 class="section-title">商品信息</h3>
          <el-table :data="currentOrder.goods || []" border stripe>
            <el-table-column prop="goods_name" label="商品名称" min-width="150" />
            <el-table-column prop="sku_name" label="规格" min-width="120" />
            <el-table-column prop="quantity" label="数量" width="80" align="center" />
            <el-table-column label="原价" width="100" align="right">
              <template #default="{ row }">
                ¥{{ row.original_price.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column label="售价" width="100" align="right">
              <template #default="{ row }">
                ¥{{ row.sale_price.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column label="小计" width="100" align="right">
              <template #default="{ row }">
                ¥{{ (row.sale_price * row.quantity).toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 金额信息 -->
        <div class="detail-section">
          <h3 class="section-title">金额信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单总额">
              <span class="price-text">¥{{ currentOrder.total_amount.toFixed(2) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="成本金额">
              <span class="price-text">¥{{ currentOrder.cost_amount.toFixed(2) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="利润金额">
              <span class="price-text profit">¥{{ currentOrder.profit_amount.toFixed(2) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="瑞幸原价" v-if="currentOrder.luckin_price">
              <span class="price-text original">¥{{ currentOrder.luckin_price.toFixed(2) }}</span>
              <el-tag type="info" size="small" class="ml-2">使用自定义价格</el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 取餐二维码 -->
        <div v-if="currentOrder.qr_data" class="detail-section qr-section">
          <h3 class="section-title">取餐二维码</h3>
          <div class="qr-container">
            <el-image 
              :src="`data:image/png;base64,${currentOrder.qr_data}`" 
              class="qr-image"
              fit="contain"
            >
              <template #error>
                <div class="qr-error">
                  <el-icon><Picture /></el-icon>
                  <span>二维码加载失败</span>
                </div>
              </template>
            </el-image>
            <div class="qr-tips">
              请使用瑞幸咖啡APP扫描二维码取餐
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button 
          type="primary" 
          v-if="currentOrder?.status === 1"
          @click="handleRefreshStatus(currentOrder)"
        >
          刷新状态
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Refresh, Document, CircleCheck, Wallet, TrendCharts,
  Picture, Download, View, Timer
} from '@element-plus/icons-vue'
import EnhancedTable from '@/components/EnhancedTable.vue'
import { getOrders, getOrder, refreshOrderStatus, getOrderStatistics } from '@/api/order'
import { getDistributors } from '@/api/distributor'
import { formatDate } from '@/utils/date'
import * as XLSX from 'xlsx'

// 状态管理
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const dateRange = ref([])
const statistics = ref({})
const showSearch = ref(false)
const distributors = ref([])
const distributorsLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  orderNo: '',
  distributorId: '',
  status: '',
  phoneNumber: '',
  takeCode: ''
})

// 详情相关
const detailDialogVisible = ref(false)
const currentOrder = ref(null)

// 表格引用
const tableRef = ref()

// 日期快捷选项
const dateShortcuts = [
  {
    text: '今天',
    value: () => {
      const end = new Date()
      const start = new Date()
      return [start, end]
    }
  },
  {
    text: '昨天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setDate(start.getDate() - 1)
      end.setDate(end.getDate() - 1)
      return [start, end]
    }
  },
  {
    text: '最近一周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setDate(start.getDate() - 7)
      return [start, end]
    }
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setMonth(start.getMonth() - 1)
      return [start, end]
    }
  }
]

// 表格列配置
const tableColumns = [
  {
    prop: 'order_no',
    label: '订单号',
    width: 200,
    fixed: 'left',
    showOverflowTooltip: true,
    slot: 'orderNo'
  },
  {
    prop: 'distributor_name',
    label: '分销商',
    width: 150,
    showOverflowTooltip: true
  },
  {
    prop: 'store_name',
    label: '门店信息',
    minWidth: 200,
    showOverflowTooltip: true,
    slot: 'store'
  },
  {
    prop: 'goods',
    label: '商品信息',
    minWidth: 250,
    showOverflowTooltip: true,
    slot: 'goods'
  },
  {
    prop: 'phone_number',
    label: '手机号',
    width: 120
  },
  {
    prop: 'total_amount',
    label: '金额',
    width: 120,
    align: 'right',
    slot: 'amount',
    sortable: 'custom'
  },
  {
    prop: 'status',
    label: '状态',
    width: 100,
    align: 'center',
    slot: 'status'
  },
  {
    prop: 'take_code',
    label: '取餐码',
    width: 100,
    align: 'center'
  },
  {
    prop: 'created_at',
    label: '下单时间',
    width: 180,
    sortable: 'custom',
    slot: 'createdAt'
  }
]

// 表格操作配置
const tableActions = computed(() => [
  {
    text: '详情',
    type: 'primary',
    icon: View,
    handler: handleViewDetail
  },
  (row) => ({
    text: '刷新状态',
    type: 'warning',
    icon: Refresh,
    handler: () => handleRefreshStatus(row),
    hide: row.status !== 1
  })
])

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

// 获取相对时间
const getTimeAgo = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return ''
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    // 添加搜索条件
    if (searchForm.orderNo) {
      params.order_no = searchForm.orderNo
    }
    if (searchForm.distributorId) {
      params.distributor_id = searchForm.distributorId
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }
    if (searchForm.phoneNumber) {
      params.phone_number = searchForm.phoneNumber
    }
    if (searchForm.takeCode) {
      params.take_code = searchForm.takeCode
    }
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const { data } = await getOrders(params)
    tableData.value = data.list || []
    total.value = data.total || 0
    
    // 加载统计信息
    loadStatistics()
  } catch (error) {
    console.error('Failed to load orders:', error)
    ElMessage.error('加载订单数据失败')
  } finally {
    loading.value = false
  }
}

// 加载统计信息
const loadStatistics = async () => {
  try {
    const params = {}
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    if (searchForm.distributorId) {
      params.distributor_id = searchForm.distributorId
    }
    
    const { data } = await getOrderStatistics(params)
    statistics.value = data
  } catch (error) {
    console.error('Failed to load statistics:', error)
  }
}

// 加载分销商列表
const loadDistributors = async () => {
  distributorsLoading.value = true
  try {
    const { data } = await getDistributors({ page_size: 1000 })
    distributors.value = data.list || []
  } catch (error) {
    console.error('Failed to load distributors:', error)
  } finally {
    distributorsLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

// 重置
const handleReset = () => {
  searchForm.orderNo = ''
  searchForm.distributorId = ''
  searchForm.status = ''
  searchForm.phoneNumber = ''
  searchForm.takeCode = ''
  dateRange.value = []
  handleSearch()
}

// 查看详情
const handleViewDetail = async (row) => {
  try {
    const { data } = await getOrder(row.order_no)
    currentOrder.value = data
    detailDialogVisible.value = true
  } catch (error) {
    console.error('Failed to get order detail:', error)
    ElMessage.error('获取订单详情失败')
  }
}

// 刷新状态
const handleRefreshStatus = async (row) => {
  try {
    await refreshOrderStatus(row.order_no)
    ElMessage.success('状态已刷新')
    loadData()
    if (detailDialogVisible.value) {
      detailDialogVisible.value = false
    }
  } catch (error) {
    console.error('Failed to refresh order status:', error)
    ElMessage.error('刷新状态失败')
  }
}

// 导出数据
const handleExport = ({ type, data, columns }) => {
  try {
    const exportData = data.map(row => {
      const obj = {}
      columns.forEach(col => {
        if (col.prop === 'status') {
          obj[col.label] = getOrderStatusText(row[col.prop])
        } else if (col.prop === 'goods') {
          obj[col.label] = row.goods?.map(g => `${g.goods_name} ${g.sku_name} ×${g.quantity}`).join('; ')
        } else if (col.prop === 'created_at') {
          obj[col.label] = formatDate(row[col.prop])
        } else if (col.prop === 'total_amount') {
          obj[col.label] = `¥${row[col.prop].toFixed(2)}`
        } else {
          obj[col.label] = row[col.prop]
        }
      })
      return obj
    })

    if (type === 'excel' || type === 'csv') {
      const ws = XLSX.utils.json_to_sheet(exportData)
      const wb = XLSX.utils.book_new()
      XLSX.utils.book_append_sheet(wb, ws, '订单数据')
      
      const fileName = `订单数据_${new Date().toLocaleDateString()}.${type === 'excel' ? 'xlsx' : 'csv'}`
      XLSX.writeFile(wb, fileName)
      
      ElMessage.success('导出成功')
    } else {
      ElMessage.warning('PDF导出功能开发中')
    }
  } catch (error) {
    console.error('Export failed:', error)
    ElMessage.error('导出失败')
  }
}

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  loadData()
}

// 当前页改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  loadData()
}

// 初始化
onMounted(() => {
  loadData()
  loadDistributors()
})
</script>

<style lang="scss" scoped>
.enhanced-orders-page {
  padding: var(--spacing-4);
  background: var(--bg-secondary);
  min-height: 100vh;

  // 页面头部
  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--spacing-6);
    
    .page-title {
      font-size: 1.75rem;
      font-weight: 700;
      color: var(--text-primary);
      margin: 0;
    }
    
    .page-actions {
      display: flex;
      gap: var(--spacing-3);
    }
  }

  // 统计卡片
  .statistics-section {
    margin-bottom: var(--spacing-6);
    
    .stat-card {
      background: var(--bg-primary);
      border-radius: var(--radius-xl);
      padding: var(--spacing-6);
      display: flex;
      align-items: center;
      gap: var(--spacing-4);
      height: 100%;
      min-height: 120px;
      transition: all 0.3s ease;
      
      &:hover {
        box-shadow: var(--shadow-lg);
        transform: translateY(-2px);
      }
      
      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: var(--radius-lg);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 28px;
        
        &.total {
          background: var(--primary-100);
          color: var(--primary-600);
        }
        
        &.success {
          background: var(--success-100);
          color: var(--success-600);
        }
        
        &.amount {
          background: var(--warning-100);
          color: var(--warning-600);
        }
        
        &.rate {
          background: var(--info-100);
          color: var(--info-600);
        }
      }
      
      .stat-content {
        flex: 1;
        
        .stat-value {
          font-size: 2rem;
          font-weight: 700;
          color: var(--text-primary);
          line-height: 1.2;
        }
        
        .stat-label {
          font-size: 0.875rem;
          color: var(--text-secondary);
          margin-top: var(--spacing-1);
        }
      }
    }
  }

  // 搜索卡片
  .search-card {
    margin-bottom: var(--spacing-6);
    
    :deep(.el-card__body) {
      padding: var(--spacing-6);
    }
    
    .el-form {
      .el-form-item {
        margin-bottom: var(--spacing-4);
      }
    }
  }

  // 表格卡片
  .table-card {
    :deep(.el-card__body) {
      padding: 0;
    }
    
    .el-pagination {
      padding: var(--spacing-4) var(--spacing-6);
      background: var(--bg-primary);
      border-top: 1px solid var(--border-light);
    }
  }

  // 自定义单元格样式
  .amount-cell {
    .main-amount {
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .sub-amount {
      font-size: 0.75rem;
      color: var(--success-600);
      margin-top: 2px;
    }
  }

  .goods-info {
    .goods-item {
      display: flex;
      align-items: center;
      gap: var(--spacing-2);
      margin-bottom: var(--spacing-1);
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .goods-name {
        color: var(--text-primary);
        font-weight: 500;
      }
      
      .goods-sku {
        color: var(--text-secondary);
        font-size: 0.875rem;
      }
      
      .goods-qty {
        color: var(--text-tertiary);
        font-size: 0.875rem;
        margin-left: auto;
      }
    }
  }

  .store-info {
    .store-name {
      color: var(--text-primary);
      font-weight: 500;
      margin-bottom: var(--spacing-1);
    }
    
    .store-address {
      color: var(--text-secondary);
      font-size: 0.875rem;
      line-height: 1.4;
    }
  }

  .time-info {
    .time-ago {
      font-size: 0.75rem;
      color: var(--text-tertiary);
      margin-top: 2px;
    }
  }
}

// 订单详情对话框
.order-detail-dialog {
  :deep(.el-dialog__body) {
    padding: 0;
  }
  
  .order-detail-content {
    .detail-section {
      padding: var(--spacing-6);
      border-bottom: 1px solid var(--border-light);
      
      &:last-child {
        border-bottom: none;
      }
      
      .section-title {
        font-size: 1.125rem;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 var(--spacing-4) 0;
      }
      
      &.qr-section {
        background: var(--bg-secondary);
      }
    }
    
    .mono-text {
      font-family: 'Courier New', monospace;
      font-weight: 500;
    }
    
    .take-code {
      font-size: 1.25rem;
      font-weight: 700;
      color: var(--primary-600);
      font-family: 'Courier New', monospace;
    }
    
    .price-text {
      font-weight: 600;
      font-size: 1.125rem;
      
      &.profit {
        color: var(--success-600);
      }
      
      &.original {
        text-decoration: line-through;
        color: var(--text-tertiary);
      }
    }
    
    .qr-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: var(--spacing-4);
      
      .qr-image {
        width: 200px;
        height: 200px;
        border: 1px solid var(--border-base);
        border-radius: var(--radius-base);
        padding: var(--spacing-2);
        background: white;
      }
      
      .qr-error {
        width: 200px;
        height: 200px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: var(--spacing-2);
        color: var(--text-tertiary);
        
        .el-icon {
          font-size: 3rem;
        }
      }
      
      .qr-tips {
        color: var(--text-secondary);
        font-size: 0.875rem;
      }
    }
  }
}

// 响应式布局
@media (max-width: 768px) {
  .enhanced-orders-page {
    padding: var(--spacing-3);
    
    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: var(--spacing-3);
      
      .page-title {
        font-size: 1.5rem;
      }
    }
    
    .statistics-section {
      .stat-card {
        padding: var(--spacing-4);
        
        .stat-icon {
          width: 48px;
          height: 48px;
          font-size: 24px;
        }
        
        .stat-content {
          .stat-value {
            font-size: 1.5rem;
          }
        }
      }
    }
  }
  
  .order-detail-dialog {
    :deep(.el-dialog) {
      width: 95% !important;
    }
  }
}
</style>