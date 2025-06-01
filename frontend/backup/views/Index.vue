<template>
  <div class="page-container">
    <!-- 搜索表单 -->
    <div class="search-form search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="卡片代码">
          <el-input
            v-model="searchForm.cardCode"
            placeholder="请输入卡片代码"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择" clearable filterable>
            <el-option label="全部" value="" />
            <el-option label="未使用" :value="0" />
            <el-option label="已使用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="批次ID">
          <el-input
            v-model="searchForm.batchId"
            placeholder="请输入批次ID"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 工具栏 -->
    <div class="table-toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        添加卡片
      </el-button>
      <el-button type="success" @click="handleBatchImport">
        <el-icon><Upload /></el-icon>
        批量导入
      </el-button>
      <el-button @click="handleViewBatches">
        <el-icon><Files /></el-icon>
        批次管理
      </el-button>
      <el-button @click="showValidateDialog = true">
        <el-icon><Check /></el-icon>
        验证卡片
      </el-button>
      <el-dropdown @command="handleValidationMode" :disabled="false" trigger="click" placement="bottom-end">
        <el-button type="warning" :loading="batchValidating">
          <el-icon><Refresh /></el-icon>
          批量验证
          <el-icon><ArrowDown /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="smart" :disabled="smartTaskRunning">
              <div class="validation-mode-item">
                <el-icon color="#409eff"><Lightning /></el-icon>
                <div>
                  <div class="mode-title">
                    智能验证
                    <el-tag v-if="smartTaskRunning" type="warning" size="small">运行中</el-tag>
                  </div>
                  <div class="mode-desc">只验证重要卡片，快速完成</div>
                </div>
              </div>
            </el-dropdown-item>
            <el-dropdown-item command="all" :disabled="fullTaskRunning">
              <div class="validation-mode-item">
                <el-icon color="#67c23a"><Check /></el-icon>
                <div>
                  <div class="mode-title">
                    全量验证
                    <el-tag v-if="fullTaskRunning" type="warning" size="small">运行中</el-tag>
                  </div>
                  <div class="mode-desc">验证所有卡片，后台安全执行</div>
                </div>
              </div>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-dropdown 
        v-if="hasRunningTask && runningTasks.size > 1" 
        @command="showTaskProgress"
        size="small"
      >
        <el-button type="info" size="small">
          <el-icon><View /></el-icon>
          查看进度
          <el-icon><ArrowDown /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item 
              v-for="[taskId, task] in runningTasks" 
              :key="taskId"
              :command="taskId"
            >
              {{ task.mode === 'smart' ? '智能验证' : '全量验证' }}进度
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-button 
        v-else-if="hasRunningTask" 
        type="info" 
        @click="showCurrentProgress"
        size="small"
      >
        <el-icon><View /></el-icon>
        查看进度
      </el-button>
    </div>
    
    <!-- 表格 -->
    <div class="table-card">
    <el-table :data="tableData" v-loading="loading" stripe class="table-optimized">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="card_code" label="卡片代码" width="120">
        <template #default="{ row }">
          <el-tooltip :content="'完整卡片代码已隐藏'" placement="top">
            <span>{{ row.card_code }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="batch_id" label="批次ID" width="80" align="center" />
      <el-table-column prop="cost_price" label="成本价" width="90" align="right">
        <template #default="{ row }">
          ￥{{ row.cost_price?.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="sell_price" label="销售价" width="90" align="right">
        <template #default="{ row }">
          ￥{{ row.sell_price?.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="getCardStatusType(row.status)">
            {{ getCardStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sync_status" label="同步状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getSyncStatusType(row.sync_status)">
            {{ getSyncStatusText(row.sync_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="used_at" label="使用时间" width="160" show-overflow-tooltip />
      <el-table-column prop="expired_at" label="过期时间" width="160" show-overflow-tooltip />
      <el-table-column prop="bound_product_count" label="绑定产品" width="90" align="center">
        <template #default="{ row }">
          <el-badge :value="row.bound_product_count || 0" :max="99" type="primary">
            <el-button link type="primary" @click="handleManageBindings(row)">
              查看
            </el-button>
          </el-badge>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="操作" width="320" fixed="right">
        <template #default="{ row }">
          <div class="table-action-buttons">
            <el-button link type="primary" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button link type="primary" @click="handleQuickValidate(row)">
              验证
            </el-button>
            <el-button link type="primary" @click="handleViewLogs(row)">
              日志
            </el-button>
            <el-button link type="danger" @click="handleDelete(row)">
              删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
    </div>
    
    <!-- 分页 -->
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
    
    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="卡片代码" prop="card_code" v-if="!isEdit">
          <el-input v-model="form.card_code" placeholder="请输入卡片代码" show-password />
        </el-form-item>
        <el-form-item label="瑞幸产品ID" prop="luckin_product_id" v-if="!isEdit">
          <el-input-number v-model="form.luckin_product_id" :min="1" :max="100" style="width: 100%;" />
          <div class="form-tip">瑞幸API所需的产品ID，默认为6，一般不需要修改</div>
        </el-form-item>
        <el-form-item label="成本价" prop="cost_price" v-if="!isEdit">
          <el-input-number v-model="form.cost_price" :min="0" :precision="2" style="width: 100%;">
            <template #append>元</template>
          </el-input-number>
          <div class="form-tip">您采购该卡片的成本价格</div>
        </el-form-item>
        <el-form-item label="销售价" prop="sell_price" v-if="!isEdit">
          <el-input-number v-model="form.sell_price" :min="0" :precision="2" style="width: 100%;">
            <template #append>元</template>
          </el-input-number>
          <div class="form-tip">
            <el-alert type="warning" :closable="false" show-icon>
              <template #title>
                重要：此价格将作为实际扣款金额，不再使用瑞幸返回的价格
              </template>
            </el-alert>
          </div>
        </el-form-item>
        <el-form-item label="状态" prop="status" v-if="isEdit">
          <el-select v-model="form.status" filterable>
            <el-option label="未使用" :value="0" />
            <el-option label="已使用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="成本价" prop="cost_price" v-if="isEdit">
          <el-input-number v-model="form.cost_price" :min="0" :precision="2" style="width: 100%;">
            <template #append>元</template>
          </el-input-number>
          <div class="form-tip">您采购该卡片的成本价格</div>
        </el-form-item>
        <el-form-item label="销售价" prop="sell_price" v-if="isEdit">
          <el-input-number v-model="form.sell_price" :min="0" :precision="2" style="width: 100%;">
            <template #append>元</template>
          </el-input-number>
          <div class="form-tip">
            <el-alert type="warning" :closable="false" show-icon>
              <template #title>
                重要：此价格将作为实际扣款金额，不再使用瑞幸返回的价格
              </template>
            </el-alert>
          </div>
        </el-form-item>
        <el-form-item label="过期时间" prop="expired_at">
          <el-date-picker
            v-model="form.expired_at"
            type="datetime"
            placeholder="选择日期时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 使用日志对话框 -->
    <el-dialog
      v-model="logsDialogVisible"
      title="卡片使用日志"
      width="800px"
    >
      <el-table :data="logsData" v-loading="logsLoading">
        <el-table-column prop="created_at" label="使用时间" width="180" />
        <el-table-column prop="distributor_id" label="分销商ID" />
        <el-table-column prop="order_id" label="订单号" />
        <el-table-column prop="success" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error_message" label="错误信息" show-overflow-tooltip />
      </el-table>
      <el-pagination
        v-model:current-page="logsPage"
        :page-size="logsPageSize"
        :total="logsTotal"
        layout="total, prev, pager, next"
        @current-change="loadUsageLogs"
      />
    </el-dialog>
    
    <!-- 验证卡片对话框 -->
    <el-dialog v-model="showValidateDialog" title="验证卡片" width="500px">
      <el-form :model="validateForm" ref="validateFormRef">
        <el-form-item label="卡片代码" prop="card_code" :rules="[{ required: true, message: '请输入卡片代码', trigger: 'blur' }]">
          <el-input v-model="validateForm.card_code" placeholder="请输入要验证的卡片代码" />
        </el-form-item>
      </el-form>
      <div v-if="validateResult" class="validate-result">
        <el-alert
          :title="validateResult.is_valid ? '卡片有效' : '卡片无效'"
          :type="validateResult.is_valid ? 'success' : 'error'"
          :closable="false"
          show-icon
        >
          <template #default>
            <div>{{ validateResult.message }}</div>
            <div v-if="validateResult.updated" style="margin-top: 8px; font-weight: bold;">
              ✅ 卡片状态已自动更新
            </div>
          </template>
        </el-alert>
      </div>
      <template #footer>
        <el-button @click="closeValidateDialog">关闭</el-button>
        <el-button type="primary" @click="handleValidateCard" :loading="validateLoading">验证</el-button>
      </template>
    </el-dialog>
    
    <!-- 验证进度对话框 -->
    <el-dialog 
      v-model="progressDialogVisible" 
      title="批量验证进度" 
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <div class="validation-progress">
        <div class="progress-header">
          <el-tag :type="getStatusType(validationTask?.status)">
            {{ getStatusText(validationTask?.status) }}
          </el-tag>
          <span class="mode-info">{{ validationTask?.mode === 'smart' ? '智能验证' : '全量验证' }}</span>
        </div>
        
        <div v-if="validationTask?.progress" class="progress-content">
          <el-progress 
            :percentage="getProgressPercentage()" 
            :status="validationTask.status === 'completed' ? 'success' : 'active'"
            :stroke-width="12"
          />
          
          <div class="progress-stats">
            <div class="stat-item">
              <span class="stat-label">总计:</span>
              <span class="stat-value">{{ validationTask.progress.total }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">已处理:</span>
              <span class="stat-value">{{ validationTask.progress.processed }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">有效:</span>
              <span class="stat-value success">{{ validationTask.progress.valid }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">无效:</span>
              <span class="stat-value warning">{{ validationTask.progress.invalid }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">失败:</span>
              <span class="stat-value danger">{{ validationTask.progress.failed }}</span>
            </div>
          </div>
          
          <div v-if="validationTask.status === 'running'" class="progress-time">
            <span>预计剩余时间: {{ getEstimatedTime() }}</span>
          </div>
        </div>
        
        <div v-if="validationTask?.error" class="error-message">
          <el-alert type="error" :title="validationTask.error" :closable="false" />
        </div>
      </div>
      
      <template #footer>
        <el-button 
          v-if="validationTask?.status === 'running'" 
          type="danger" 
          @click="handleCancelValidation"
        >
          取消验证
        </el-button>
        <el-button 
          type="primary" 
          @click="progressDialogVisible = false"
          :disabled="validationTask?.status === 'running'"
        >
          {{ validationTask?.status === 'running' ? '验证中...' : '关闭' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElLoading } from 'element-plus'
import { Plus, Upload, Files, Check, Refresh, ArrowDown, Lightning, View } from '@element-plus/icons-vue'
import { getCards, createCard, updateCard, deleteCard, getCardUsageLogs, validateCard, batchValidateCards, startBatchValidation, getValidationProgress, cancelValidation } from '@/api/card'
import { formatDate } from '@/utils/date'
import { showSuccess, showError, showDeleteConfirm } from '@/utils/message'

const router = useRouter()

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  cardCode: '',
  status: '',
  batchId: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref()

const form = reactive({
  card_code: '',
  luckin_product_id: 6,  // 默认值6
  cost_price: 0,
  sell_price: 0,
  status: 0,
  expired_at: '',
  description: ''
})

const rules = {
  card_code: [
    { required: true, message: '请输入卡片代码', trigger: 'blur' }
  ],
  luckin_product_id: [
    { required: true, message: '请输入瑞幸产品ID', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '产品ID范围为1-100', trigger: 'blur' }
  ],
  cost_price: [
    { required: true, message: '请输入成本价', trigger: 'blur' }
  ],
  sell_price: [
    { required: true, message: '请输入销售价', trigger: 'blur' }
  ],
  expired_at: [
    { required: true, message: '请选择过期时间', trigger: 'change' }
  ]
}

// 使用日志相关
const logsDialogVisible = ref(false)
const logsLoading = ref(false)
const logsData = ref([])
const logsTotal = ref(0)
const logsPage = ref(1)
const logsPageSize = 20
const currentCardId = ref(null)

// 验证卡片相关
const showValidateDialog = ref(false)
const validateLoading = ref(false)
const validateFormRef = ref()
const validateForm = reactive({
  card_code: ''
})
const validateResult = ref(null)

// 批量验证相关
const batchValidating = ref(false)

// 新版验证进度相关
const progressDialogVisible = ref(false)
const validationTask = ref(null)
const progressTimer = ref(null)
const hasRunningTask = ref(false)
const currentTaskId = ref(null)

// 支持多任务并行
const runningTasks = ref(new Map()) // taskId -> task info
const smartTaskRunning = ref(false)
const fullTaskRunning = ref(false)


// 获取卡片状态类型
const getCardStatusType = (status) => {
  const types = {
    0: 'success',
    1: 'info'
  }
  return types[status] || 'info'
}

// 获取卡片状态文本
const getCardStatusText = (status) => {
  const texts = {
    0: '未使用',
    1: '已使用'
  }
  return texts[status] || '未知'
}

// 获取同步状态类型
const getSyncStatusType = (status) => {
  const types = {
    'pending': 'info',
    'syncing': 'warning',
    'synced': 'success',
    'failed': 'danger'
  }
  return types[status] || 'info'
}

// 获取同步状态文本
const getSyncStatusText = (status) => {
  const texts = {
    'pending': '待同步',
    'syncing': '同步中',
    'synced': '已同步',
    'failed': '同步失败'
  }
  return texts[status] || '未知'
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
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }
    if (searchForm.batchId) {
      params.batch_id = searchForm.batchId
    }
    if (searchForm.cardCode) {
      // cardCode 在后端处理
    }
    
    const { data } = await getCards(params)
    tableData.value = data.list.map(item => ({
      ...item,
      expired_at: formatDate(item.expired_at)
    }))
    total.value = data.total
  } catch (error) {
    console.error('Failed to load cards:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

// 重置
const handleReset = () => {
  searchForm.cardCode = ''
  searchForm.status = ''
  searchForm.batchId = ''
  handleSearch()
}

// 添加
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加卡片'
  // 设置默认过期时间为一年后
  const defaultExpiredAt = new Date()
  defaultExpiredAt.setFullYear(defaultExpiredAt.getFullYear() + 1)
  
  Object.assign(form, {
    card_code: '',
    luckin_product_id: 6,  // 默认值6
    cost_price: 0,
    sell_price: 0,
    status: 0,
    expired_at: defaultExpiredAt.toISOString().slice(0, -1) + 'Z', // 格式化为 YYYY-MM-DDTHH:mm:ssZ
    description: ''
  })
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑卡片'
  Object.assign(form, {
    id: row.id,
    status: row.status,
    cost_price: row.cost_price,
    sell_price: row.sell_price,
    expired_at: row.expired_at,
    description: row.description
  })
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return
  
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await updateCard(form.id, form)
      showSuccess('卡片更新成功')
    } else {
      await createCard(form)
      showSuccess('卡片创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error) {
    console.error('Submit failed:', error)
    // 错误已经在 request.js 中通过 ElMessage 显示了
    // 如果需要额外处理，可以在这里添加
  } finally {
    submitLoading.value = false
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await showDeleteConfirm(`卡片 "${row.card_code}"`)
    
    await deleteCard(row.id)
    showSuccess('卡片删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete failed:', error)
    }
  }
}

// 查看使用日志
const handleViewLogs = (row) => {
  currentCardId.value = row.id
  logsPage.value = 1
  logsDialogVisible.value = true
  loadUsageLogs()
}

// 加载使用日志
const loadUsageLogs = async () => {
  logsLoading.value = true
  try {
    const { data } = await getCardUsageLogs(currentCardId.value, {
      page: logsPage.value,
      page_size: logsPageSize
    })
    logsData.value = data.list.map(item => ({
      ...item,
      created_at: formatDate(item.created_at)
    }))
    logsTotal.value = data.total
  } catch (error) {
    console.error('Failed to load usage logs:', error)
  } finally {
    logsLoading.value = false
  }
}

// 管理绑定
const handleManageBindings = (row) => {
  // 跳转到产品管理页面，显示该卡片绑定的产品
  router.push({
    path: '/products',
    query: { card_id: row.id, card_code: row.card_code }
  })
}

// 批量导入
const handleBatchImport = () => {
  router.push('/cards/batch-import')
}

// 查看批次
const handleViewBatches = () => {
  router.push('/cards/batches')
}

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val
  loadData()
}

// 当前页改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  loadData()
}

// 验证卡片
const handleValidateCard = async () => {
  const valid = await validateFormRef.value.validate()
  if (!valid) return
  
  validateLoading.value = true
  validateResult.value = null
  
  try {
    const res = await validateCard(validateForm.card_code)
    validateResult.value = res.data
  } catch (error) {
    showError('验证失败')
  } finally {
    validateLoading.value = false
  }
}

// 关闭验证对话框
const closeValidateDialog = () => {
  showValidateDialog.value = false
  validateForm.card_code = ''
  
  // 如果有状态更新，刷新列表
  if (validateResult.value && validateResult.value.updated) {
    loadData()
  }
  
  validateResult.value = null
}

// 快速验证单个卡片
const handleQuickValidate = async (row) => {
  try {
    const loading = ElLoading.service({
      lock: true,
      text: '正在验证卡片...',
      background: 'rgba(0, 0, 0, 0.7)'
    })
    
    const res = await validateCard(row.card_code)
    
    loading.close()
    
    if (res.data.updated) {
      showSuccess(`验证完成：${res.data.message}`)
      // 刷新列表
      loadData()
    } else {
      ElMessage({
        message: res.data.message,
        type: res.data.is_valid ? 'success' : 'warning',
        duration: 3000
      })
    }
  } catch (error) {
    showError('验证失败')
  }
}

// 处理验证模式选择
const handleValidationMode = async (mode) => {
  try {
    // 检查该模式是否已在运行
    if (mode === 'smart' && smartTaskRunning.value) {
      // 智能验证已在运行，显示其进度
      const smartTask = Array.from(runningTasks.value.values()).find(task => task.mode === 'smart')
      if (smartTask) {
        validationTask.value = smartTask
        progressDialogVisible.value = true
        startProgressPolling()
        ElMessage.warning('智能验证正在进行中，已为您打开进度窗口')
        return
      }
    }
    
    if (mode === 'all' && fullTaskRunning.value) {
      // 全量验证已在运行，显示其进度
      const fullTask = Array.from(runningTasks.value.values()).find(task => task.mode === 'all')
      if (fullTask) {
        validationTask.value = fullTask
        progressDialogVisible.value = true
        startProgressPolling()
        ElMessage.warning('全量验证正在进行中，已为您打开进度窗口')
        return
      }
    }
    
    let confirmMessage = ''
    if (mode === 'smart') {
      confirmMessage = '智能验证只会检查重要的卡片（有异常、新添加、有订单的），大约需要2-5分钟。确定要继续吗？'
    } else {
      confirmMessage = '全量验证会检查所有卡片的真实状态，会在后台安全执行，可能需要较长时间。确定要继续吗？'
    }
    
    await showDeleteConfirm(confirmMessage)
    
    // 启动验证任务
    const res = await startBatchValidation(mode)
    
    // 设置任务信息并显示进度对话框
    validationTask.value = res.data
    progressDialogVisible.value = true
    
    // 更新运行状态
    runningTasks.value.set(res.data.id, res.data)
    if (mode === 'smart') {
      smartTaskRunning.value = true
    } else if (mode === 'all') {
      fullTaskRunning.value = true
    }
    
    // 更新全局状态（用于旧逻辑兼容）
    hasRunningTask.value = true
    currentTaskId.value = res.data.id
    batchValidating.value = smartTaskRunning.value || fullTaskRunning.value
    
    // 保存到本地存储
    const savedTasks = JSON.parse(localStorage.getItem('validation_tasks') || '{}')
    savedTasks[res.data.id] = { mode, taskId: res.data.id, startTime: Date.now() }
    localStorage.setItem('validation_tasks', JSON.stringify(savedTasks))
    
    // 开始轮询进度
    startProgressPolling()
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Start validation failed:', error)
      showError('启动验证失败')
    }
    batchValidating.value = false
  }
}

// 开始轮询进度
const startProgressPolling = () => {
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
  }
  
  progressTimer.value = setInterval(async () => {
    try {
      const res = await getValidationProgress(validationTask.value.id)
      validationTask.value = res.data
      
      // 如果任务完成，停止轮询
      if (res.data.status === 'completed' || res.data.status === 'failed' || res.data.status === 'cancelled') {
        // 清理任务状态
        const taskMode = res.data.mode
        runningTasks.value.delete(res.data.id)
        
        if (taskMode === 'smart') {
          smartTaskRunning.value = false
        } else if (taskMode === 'all') {
          fullTaskRunning.value = false
        }
        
        // 更新全局状态
        batchValidating.value = smartTaskRunning.value || fullTaskRunning.value
        
        // 如果没有运行中的任务了，清理全局状态
        if (runningTasks.value.size === 0) {
          hasRunningTask.value = false
          currentTaskId.value = null
          clearInterval(progressTimer.value)
          progressTimer.value = null
        }
        
        // 清理本地存储中的该任务
        const savedTasks = JSON.parse(localStorage.getItem('validation_tasks') || '{}')
        delete savedTasks[res.data.id]
        if (Object.keys(savedTasks).length === 0) {
          localStorage.removeItem('validation_tasks')
        } else {
          localStorage.setItem('validation_tasks', JSON.stringify(savedTasks))
        }
        
        // 显示完成消息
        if (res.data.status === 'completed') {
          const progress = res.data.progress
          const modeText = taskMode === 'smart' ? '智能验证' : '全量验证'
          let message = `${modeText}完成！\n总计：${progress.total}张\n有效：${progress.valid}张\n无效：${progress.invalid}张`
          if (progress.failed > 0) {
            message += `\n失败：${progress.failed}张`
          }
          showSuccess(message)
          
          // 刷新列表
          loadData()
        } else if (res.data.status === 'failed') {
          const modeText = taskMode === 'smart' ? '智能验证' : '全量验证'
          showError(`${modeText}失败：${res.data.error}`)
        }
      }
    } catch (error) {
      console.error('获取验证进度失败:', error)
    }
  }, 2000) // 每2秒轮询一次
}

// 取消验证任务
const handleCancelValidation = async () => {
  try {
    await cancelValidation(validationTask.value.id)
    showSuccess('验证任务已取消')
  } catch (error) {
    showError('取消验证失败')
  }
}

// 获取进度百分比
const getProgressPercentage = () => {
  if (!validationTask.value?.progress) return 0
  const { total, processed } = validationTask.value.progress
  if (total === 0) return 0
  return Math.round((processed / total) * 100)
}

// 获取状态类型
const getStatusType = (status) => {
  const types = {
    'queued': 'info',
    'running': 'warning', 
    'completed': 'success',
    'failed': 'danger',
    'cancelled': 'info'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    'queued': '排队中',
    'running': '验证中',
    'completed': '已完成',
    'failed': '验证失败',
    'cancelled': '已取消'
  }
  return texts[status] || '未知状态'
}

// 获取预计时间
const getEstimatedTime = () => {
  if (!validationTask.value?.progress) return '计算中...'
  
  const { total, processed } = validationTask.value.progress
  if (processed === 0) return '计算中...'
  
  const elapsed = Date.now() - new Date(validationTask.value.started_at).getTime()
  const avgTime = elapsed / processed
  const remaining = (total - processed) * avgTime
  
  if (remaining < 60000) {
    return `${Math.round(remaining / 1000)}秒`
  } else {
    return `${Math.round(remaining / 60000)}分钟`
  }
}

// 查看当前进度
const showCurrentProgress = async () => {
  // 如果只有一个任务，直接显示
  if (runningTasks.value.size === 1) {
    const [taskId, task] = runningTasks.value.entries().next().value
    showTaskProgress(taskId)
  } else if (currentTaskId.value) {
    showTaskProgress(currentTaskId.value)
  }
}

// 查看指定任务进度
const showTaskProgress = async (taskId) => {
  try {
    const res = await getValidationProgress(taskId)
    validationTask.value = res.data
    currentTaskId.value = taskId
    progressDialogVisible.value = true
    
    if (res.data.status === 'running' || res.data.status === 'queued') {
      startProgressPolling()
    }
  } catch (error) {
    showError('获取任务进度失败')
    // 清理无效任务
    runningTasks.value.delete(taskId)
    if (runningTasks.value.size === 0) {
      hasRunningTask.value = false
      currentTaskId.value = null
    }
  }
}

// 检查是否有运行中的验证任务
const checkRunningTask = async () => {
  const savedTasks = JSON.parse(localStorage.getItem('validation_tasks') || '{}')
  if (Object.keys(savedTasks).length === 0) return
  
  let hasValidTasks = false
  
  for (const [taskId, taskInfo] of Object.entries(savedTasks)) {
    try {
      const res = await getValidationProgress(taskId)
      if (res.data.status === 'running' || res.data.status === 'queued') {
        // 恢复任务状态
        runningTasks.value.set(taskId, res.data)
        
        if (taskInfo.mode === 'smart') {
          smartTaskRunning.value = true
        } else if (taskInfo.mode === 'all') {
          fullTaskRunning.value = true
        }
        
        hasValidTasks = true
      } else {
        // 任务已完成，从本地存储中删除
        delete savedTasks[taskId]
      }
    } catch (error) {
      // 任务不存在，从本地存储中删除
      delete savedTasks[taskId]
    }
  }
  
  // 更新本地存储
  if (Object.keys(savedTasks).length === 0) {
    localStorage.removeItem('validation_tasks')
  } else {
    localStorage.setItem('validation_tasks', JSON.stringify(savedTasks))
  }
  
  // 更新全局状态
  if (hasValidTasks) {
    hasRunningTask.value = true
    batchValidating.value = smartTaskRunning.value || fullTaskRunning.value
    
    const taskCount = runningTasks.value.size
    const smartText = smartTaskRunning.value ? '智能验证' : ''
    const fullText = fullTaskRunning.value ? '全量验证' : ''
    const runningText = [smartText, fullText].filter(Boolean).join('、')
    
    ElMessage.info(`检测到正在进行的${runningText}任务，可点击"查看进度"按钮查看详情`)
  }
}

// 批量验证卡片（保持兼容性）
const handleBatchValidate = async () => {
  // 默认调用智能验证
  handleValidationMode('smart')
}

onMounted(() => {
  loadData()
  checkRunningTask()
})

onUnmounted(() => {
  // 清理定时器
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
  }
})
</script>

<style lang="scss" scoped>
.ml-10 {
  margin-left: 10px;
}

.form-tip {
  margin-top: 5px;
  color: #909399;
  font-size: 12px;
  line-height: 1.4;
}

.table-action-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
}

.validate-result {
  margin-top: 20px;
}

// 验证模式选择样式
.validation-mode-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 0;
  
  .mode-title {
    font-weight: 500;
    font-size: 14px;
    line-height: 1.2;
  }
  
  .mode-desc {
    font-size: 12px;
    color: #909399;
    line-height: 1.2;
    margin-top: 2px;
  }
}

// 验证进度样式
.validation-progress {
  .progress-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
    
    .mode-info {
      font-size: 14px;
      color: #606266;
    }
  }
  
  .progress-content {
    .progress-stats {
      display: grid;
      grid-template-columns: repeat(5, 1fr);
      gap: 16px;
      margin: 20px 0;
      
      .stat-item {
        text-align: center;
        
        .stat-label {
          display: block;
          font-size: 12px;
          color: #909399;
          margin-bottom: 4px;
        }
        
        .stat-value {
          display: block;
          font-size: 18px;
          font-weight: 600;
          
          &.success { color: #67c23a; }
          &.warning { color: #e6a23c; }
          &.danger { color: #f56c6c; }
        }
      }
    }
    
    .progress-time {
      text-align: center;
      margin-top: 16px;
      color: #606266;
      font-size: 14px;
    }
  }
  
  .error-message {
    margin-top: 16px;
  }
}
</style>