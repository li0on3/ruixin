<template>
  <div class="page-container">
    <!-- 搜索表单 -->
    <div class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="分销商名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入分销商名称"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择" clearable filterable>
            <el-option label="全部" value="" />
            <el-option label="待审核" :value="0" />
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 数据表格 -->
    <div class="table-card">
      <!-- 工具栏 -->
      <div class="table-toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        添加分销商
      </el-button>
    </div>
    
    <!-- 表格 -->
    <div class="table-container">
      <el-table :data="tableData" v-loading="loading" stripe class="table-optimized responsive-table">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="分销商名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="company_name" label="公司名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="contact_name" label="联系人" width="100" />
      <el-table-column prop="phone" label="电话" width="120" />
      <el-table-column prop="email" label="邮箱" min-width="160" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="余额信息" width="150" align="right">
        <template #default="{ row }">
          <div>余额: ¥{{ row.balance.toFixed(2) }}</div>
          <div style="font-size: 12px; color: #999;">冻结: ¥{{ (row.frozen_amount || 0).toFixed(2) }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="total_orders" label="总订单" width="80" align="center" />
      <el-table-column label="操作" width="360" fixed="right" class="column-action">
        <template #default="{ row }">
          <div class="table-action-buttons">
            <el-button link type="primary" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button link type="primary" @click="handleViewAPIKey(row)">
              密钥
            </el-button>
            <el-button link type="warning" @click="handleResetAPIKey(row)">
              重置密钥
            </el-button>
            <el-button link type="warning" @click="handleResetPassword(row)">
              重置密码
            </el-button>
            <el-button link type="primary" @click="handleViewLogs(row)">
              日志
            </el-button>
            <el-button link type="primary" @click="handleViewBalance(row)">
              余额
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
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="分销商名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分销商名称" />
        </el-form-item>
        <el-form-item label="公司名称" prop="company_name">
          <el-input v-model="form.company_name" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="联系人" prop="contact_name">
          <el-input v-model="form.contact_name" placeholder="请输入联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="回调地址" prop="callback_url">
          <el-input v-model="form.callback_url" placeholder="请输入订单状态回调地址" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">待审核</el-radio>
            <el-radio :label="1">正常</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 查看API密钥对话框 -->
    <el-dialog
      v-model="apiKeyDialogVisible"
      title="API密钥信息"
      width="600px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="API Key">
          <div class="api-key-wrapper">
            <span>{{ currentAPIKey.api_key || '加载中...' }}</span>
            <el-button
              type="primary"
              size="small"
              text
              @click="copyToClipboard(currentAPIKey.api_key)"
            >
              复制
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="API Secret">
          <div class="api-key-wrapper">
            <span>{{ currentAPIKey.api_secret || '加载中...' }}</span>
            <el-button
              type="primary"
              size="small"
              text
              @click="copyToClipboard(currentAPIKey.api_secret)"
            >
              复制
            </el-button>
          </div>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="apiKeyDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
    
    <!-- API日志对话框 -->
    <el-dialog
      v-model="apiLogDialogVisible"
      title="API调用日志"
      width="80%"
    >
      <el-table :data="apiLogs" v-loading="apiLogLoading" stripe class="table-optimized">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="method" label="请求方法" width="90" align="center" />
        <el-table-column prop="api_endpoint" label="请求路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="response_code" label="响应状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.response_code < 400 ? 'success' : 'danger'" v-if="row.response_code">
              {{ row.response_code }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="ip_address" label="IP地址" width="130" />
        <el-table-column prop="created_at" label="请求时间" width="160" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewLogDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="apiLogDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
    
    <!-- API日志详情对话框 -->
    <el-dialog
      v-model="logDetailDialogVisible"
      title="日志详情"
      width="800px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="请求头">
          <pre>{{ currentLogDetail.request_headers }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="请求体">
          <pre>{{ currentLogDetail.request_body }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应体">
          <pre>{{ currentLogDetail.response_body }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="logDetailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
    </div>
    
    <!-- 创建成功对话框 - 显示密码和API密钥 -->
    <el-dialog
      v-model="createSuccessDialogVisible"
      title="分销商创建成功"
      width="650px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-alert
        title="请妥善保存以下信息"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        初始密码仅显示一次，请立即保存并告知分销商！
      </el-alert>
      
      <el-descriptions :column="1" border>
        <el-descriptions-item label="分销商名称">
          {{ createSuccessInfo.name }}
        </el-descriptions-item>
        <el-descriptions-item label="初始密码">
          <div class="password-wrapper">
            <span class="password-text">{{ createSuccessInfo.password || '生成中...' }}</span>
            <el-button
              type="primary"
              size="small"
              text
              @click="copyToClipboard(createSuccessInfo.password, '密码')"
            >
              <el-icon><CopyDocument /></el-icon>
              复制
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="API Key">
          <div class="api-key-wrapper">
            <span class="api-key-text">{{ createSuccessInfo.api_key || '生成中...' }}</span>
            <el-button
              type="primary"
              size="small"
              text
              @click="copyToClipboard(createSuccessInfo.api_key, 'API Key')"
            >
              <el-icon><CopyDocument /></el-icon>
              复制
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="API Secret">
          <div class="api-key-wrapper">
            <span class="api-key-text">{{ createSuccessInfo.api_secret || '生成中...' }}</span>
            <el-button
              type="primary"
              size="small"
              text
              @click="copyToClipboard(createSuccessInfo.api_secret, 'API Secret')"
            >
              <el-icon><CopyDocument /></el-icon>
              复制
            </el-button>
          </div>
        </el-descriptions-item>
      </el-descriptions>
      
      <div style="margin-top: 20px; padding: 15px; background: #f5f7fa; border-radius: 4px;">
        <p style="margin: 0 0 10px; font-weight: bold; color: #303133;">使用说明：</p>
        <ol style="margin: 0; padding-left: 20px; color: #606266; line-height: 1.8;">
          <li>请将以上信息安全地发送给分销商</li>
          <li>分销商使用初始密码登录后，建议立即修改密码</li>
          <li>API密钥用于接口调用认证，请妥善保管</li>
        </ol>
      </div>
      
      <template #footer>
        <el-button type="primary" @click="createSuccessDialogVisible = false">
          我已保存
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="resetPasswordDialogVisible"
      title="重置分销商密码"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form ref="resetPasswordFormRef" :model="resetPasswordForm" :rules="resetPasswordRules" label-width="100px">
        <el-form-item label="分销商">
          <el-input :value="currentResetDistributor.name" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="resetPasswordForm.password"
            type="password"
            placeholder="请输入新密码"
            show-password
          >
            <template #append>
              <el-button @click="generateRandomPassword">
                随机生成
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-text type="info" size="small">
            密码长度至少6位，建议包含字母和数字
          </el-text>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPasswordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleResetPasswordSubmit" :loading="resetPasswordLoading">
          确定重置
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 重置密码成功对话框 -->
    <el-dialog
      v-model="resetPasswordSuccessDialogVisible"
      title="密码重置成功"
      width="500px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-alert
        title="新密码仅显示一次"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        请立即保存新密码并告知分销商！
      </el-alert>
      
      <el-descriptions :column="1" border>
        <el-descriptions-item label="分销商名称">
          {{ resetPasswordSuccessInfo.name }}
        </el-descriptions-item>
        <el-descriptions-item label="新密码">
          <div class="password-wrapper">
            <span class="password-text">{{ resetPasswordSuccessInfo.password }}</span>
            <el-button
              type="primary"
              size="small"
              text
              @click="copyToClipboard(resetPasswordSuccessInfo.password, '新密码')"
            >
              <el-icon><CopyDocument /></el-icon>
              复制
            </el-button>
          </div>
        </el-descriptions-item>
      </el-descriptions>
      
      <template #footer>
        <el-button type="primary" @click="resetPasswordSuccessDialogVisible = false">
          我已保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, CopyDocument } from '@element-plus/icons-vue'
import { getDistributors, createDistributor, updateDistributor, deleteDistributor, resetDistributorAPIKey, getDistributorAPILogs, getDistributor, resetDistributorPassword } from '@/api/distributor'
import { formatDate } from '@/utils/date'

const router = useRouter()

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  name: '',
  status: ''
})

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const formRef = ref(null)
const currentEditId = ref(null)

// 表单数据
const form = reactive({
  name: '',
  company_name: '',
  contact_name: '',
  phone: '',
  email: '',
  callback_url: '',
  status: 1,
  remark: ''
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入分销商名称', trigger: 'blur' }
  ],
  company_name: [
    { required: true, message: '请输入公司名称', trigger: 'blur' }
  ],
  contact_name: [
    { required: true, message: '请输入联系人姓名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  callback_url: [
    { pattern: /^https?:\/\//, message: '请输入正确的URL地址', trigger: 'blur' }
  ]
}

// API密钥对话框
const apiKeyDialogVisible = ref(false)
const currentAPIKey = reactive({
  api_key: '',
  api_secret: ''
})

// API日志对话框
const apiLogDialogVisible = ref(false)
const apiLogLoading = ref(false)
const apiLogs = ref([])
const logDetailDialogVisible = ref(false)
const currentLogDetail = reactive({
  request_headers: '',
  request_body: '',
  response_body: ''
})

// 创建成功对话框
const createSuccessDialogVisible = ref(false)
const createSuccessInfo = reactive({
  name: '',
  password: '',
  api_key: '',
  api_secret: ''
})

// 重置密码对话框
const resetPasswordDialogVisible = ref(false)
const resetPasswordFormRef = ref()
const resetPasswordLoading = ref(false)
const currentResetDistributor = reactive({
  id: null,
  name: ''
})
const resetPasswordForm = reactive({
  password: ''
})
const resetPasswordRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

// 重置密码成功对话框
const resetPasswordSuccessDialogVisible = ref(false)
const resetPasswordSuccessInfo = reactive({
  name: '',
  password: ''
})

// 获取状态类型
const getStatusType = (status) => {
  const types = {
    0: 'warning',
    1: 'success',
    2: 'danger'
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    0: '待审核',
    1: '正常',
    2: '禁用'
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
    
    if (searchForm.name) {
      params.name = searchForm.name
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }
    
    const { data } = await getDistributors(params)
    tableData.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error('Failed to load distributors:', error)
    ElMessage.error('加载分销商列表失败')
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
  searchForm.name = ''
  searchForm.status = ''
  handleSearch()
}

// 重置表单
const resetForm = () => {
  form.name = ''
  form.company_name = ''
  form.contact_name = ''
  form.phone = ''
  form.email = ''
  form.callback_url = ''
  form.status = 1
  form.remark = ''
  currentEditId.value = null
}

// 添加分销商
const handleAdd = () => {
  resetForm()
  dialogTitle.value = '添加分销商'
  dialogVisible.value = true
}

// 编辑分销商
const handleEdit = (row) => {
  resetForm()
  currentEditId.value = row.id
  dialogTitle.value = '编辑分销商'
  
  // 填充表单数据
  form.name = row.name
  form.company_name = row.company_name
  form.contact_name = row.contact_name
  form.phone = row.phone
  form.email = row.email
  form.callback_url = row.callback_url || ''
  form.status = row.status
  form.remark = row.remark || ''
  
  dialogVisible.value = true
}

// 取消对话框
const handleCancel = () => {
  dialogVisible.value = false
  formRef.value?.resetFields()
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    
    submitLoading.value = true
    
    if (currentEditId.value) {
      // 编辑
      await updateDistributor(currentEditId.value, form)
      ElMessage.success('更新成功')
    } else {
      // 新增
      const response = await createDistributor(form)
      
      // 显示创建成功对话框，包含密码信息
      createSuccessInfo.name = form.name
      createSuccessInfo.password = response.data.default_password || '123456'
      createSuccessInfo.api_key = response.data.api_key
      createSuccessInfo.api_secret = response.data.api_secret
      createSuccessDialogVisible.value = true
      
      dialogVisible.value = false
    }
    
    loadData()
  } catch (error) {
    if (error !== false) {
      console.error('Submit failed:', error)
      ElMessage.error(currentEditId.value ? '更新失败' : '添加失败')
    }
  } finally {
    submitLoading.value = false
  }
}

// 查看API密钥
const handleViewAPIKey = async (row) => {
  try {
    const { data } = await getDistributor(row.id)
    currentAPIKey.api_key = data.api_key
    currentAPIKey.api_secret = data.api_secret
    apiKeyDialogVisible.value = true
  } catch (error) {
    console.error('Failed to get API key:', error)
    ElMessage.error('获取API密钥失败')
  }
}

// 重置API密钥
const handleResetAPIKey = async (row) => {
  try {
    await ElMessageBox.confirm(
      '重置密钥后，原密钥将立即失效，确定要重置吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await resetDistributorAPIKey(row.id)
    ElMessage.success('密钥重置成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to reset API key:', error)
      ElMessage.error('重置密钥失败')
    }
  }
}

// 查看余额详情
const handleViewBalance = (row) => {
  router.push({
    path: '/finance/transactions',
    query: { distributor_id: row.id }
  })
}

// 查看API日志
const handleViewLogs = async (row) => {
  apiLogLoading.value = true
  apiLogDialogVisible.value = true
  
  try {
    const { data } = await getDistributorAPILogs(row.id)
    apiLogs.value = data.list || []
  } catch (error) {
    console.error('Failed to get API logs:', error)
    ElMessage.error('获取API日志失败')
  } finally {
    apiLogLoading.value = false
  }
}

// 查看日志详情
const viewLogDetail = (row) => {
  currentLogDetail.request_headers = row.request_headers || ''
  currentLogDetail.request_body = row.request_body || ''
  currentLogDetail.response_body = row.response_body || ''
  logDetailDialogVisible.value = true
}

// 删除分销商
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除分销商"${row.name}"吗？删除后不可恢复。`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteDistributor(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete distributor:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 复制到剪贴板
const copyToClipboard = (text, label = '') => {
  if (!text) {
    ElMessage.warning('没有可复制的内容')
    return
  }
  
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success(label ? `${label}已复制到剪贴板` : '复制成功')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 生成随机密码
const generateRandomPassword = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789'
  let password = ''
  for (let i = 0; i < 8; i++) {
    password += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  resetPasswordForm.password = password
  ElMessage.success('已生成随机密码')
}

// 显示重置密码对话框
const handleResetPassword = (row) => {
  currentResetDistributor.id = row.id
  currentResetDistributor.name = row.name
  resetPasswordForm.password = ''
  resetPasswordDialogVisible.value = true
}

// 提交重置密码
const handleResetPasswordSubmit = async () => {
  try {
    await resetPasswordFormRef.value?.validate()
    
    await ElMessageBox.confirm(
      `确定要重置分销商"${currentResetDistributor.name}"的密码吗？`,
      '重置密码确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    resetPasswordLoading.value = true
    
    await resetDistributorPassword(currentResetDistributor.id, {
      new_password: resetPasswordForm.password
    })
    
    // 显示重置成功对话框
    resetPasswordSuccessInfo.name = currentResetDistributor.name
    resetPasswordSuccessInfo.password = resetPasswordForm.password
    resetPasswordSuccessDialogVisible.value = true
    
    resetPasswordDialogVisible.value = false
    ElMessage.success('密码重置成功')
    
  } catch (error) {
    if (error !== 'cancel' && error !== false) {
      console.error('Reset password failed:', error)
      ElMessage.error('密码重置失败')
    }
  } finally {
    resetPasswordLoading.value = false
  }
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

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.password-wrapper,
.api-key-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
  
  .password-text,
  .api-key-text {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 14px;
    color: #303133;
    background: #f5f7fa;
    padding: 4px 8px;
    border-radius: 4px;
    user-select: all;
    word-break: break-all;
  }
  
  .el-button {
    flex-shrink: 0;
  }
}

.el-descriptions {
  :deep(.el-descriptions__label) {
    width: 120px;
  }
}
.page-container {
  width: 100%;
  overflow: hidden;
}

.table-container {
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  
  .el-table {
    width: 100%;
    
    // 响应式列宽调整
    @media (max-width: 1366px) {
      font-size: 12px;
      
      .el-table__header th {
        padding: 10px 0;
      }
      
      .el-table__body td {
        padding: 12px 0;
      }
      
      // 调整列宽以适应小屏幕
      th, td {
        &:nth-child(1) { width: 50px !important; }  // ID
        &:nth-child(4) { width: 80px !important; }  // 联系人
        &:nth-child(5) { width: 100px !important; } // 电话
        &:nth-child(7) { width: 70px !important; }  // 状态
        &:nth-child(8) { width: 120px !important; } // 余额信息
        &:nth-child(9) { width: 60px !important; }  // 总订单
        &:last-child { width: 200px !important; }   // 操作
      }
    }
    
    @media (min-width: 1366px) and (max-width: 1600px) {
      th, td {
        &:last-child { width: 240px !important; } // 操作列稍微增大
      }
    }
  }
}

.table-action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: flex-start;
  
  @media (max-width: 1366px) {
    gap: 4px;
    
    .el-button {
      padding: 4px 8px;
      font-size: 12px;
      margin: 2px 0;
    }
  }
}

// 对话框响应式
.el-dialog {
  @media (max-width: 1366px) {
    width: 95% !important;
    max-width: 600px;
  }
}
.api-key-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
}

.table-action-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}
</style>