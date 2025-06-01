<template>
  <div class="withdrawals-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-title">待处理提现</div>
            <div class="stat-value">{{ pendingCount }}笔</div>
            <div class="stat-sub">金额: ¥{{ pendingAmount.toFixed(2) }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 搜索表单 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="分销商">
          <el-select v-model="searchForm.distributor_id" placeholder="请选择或搜索分销商" clearable filterable style="width: 200px">
            <el-option
              v-for="item in distributorList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable filterable style="width: 120px">
            <el-option label="待处理" :value="0" />
            <el-option label="已处理" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="distributor.name" label="分销商" />
        <el-table-column prop="amount" label="提现金额" width="120">
          <template #default="{ row }">
            <span style="color: #f56c6c; font-weight: bold;">¥{{ row.amount.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">{{ WithdrawalStatus[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="account_info" label="收款账户" width="200">
          <template #default="{ row }">
            <div v-if="getAccountInfo(row.account_info)">
              <div>{{ getAccountInfo(row.account_info).type === 'alipay' ? '支付宝' : getAccountInfo(row.account_info).type === 'wechat' ? '微信' : '银行卡' }}</div>
              <div style="font-size: 12px; color: #666;">{{ getAccountInfo(row.account_info).account_name }}</div>
              <div style="font-size: 12px; color: #666;">{{ getAccountInfo(row.account_info).account_no }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="created_at" label="申请时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="processed_at" label="处理时间" width="160">
          <template #default="{ row }">
            {{ row.processed_at ? formatDate(row.processed_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="processor.real_name" label="处理人" width="100">
          <template #default="{ row }">
            {{ row.processor?.real_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              type="success"
              size="small"
              @click="handleProcess(row, true)"
            >
              批准
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="danger"
              size="small"
              @click="handleProcess(row, false)"
            >
              拒绝
            </el-button>
            <el-tag v-if="row.status === 2" type="danger" size="small">
              {{ row.reject_reason }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

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
    </el-card>

    <!-- 拒绝原因对话框 -->
    <el-dialog v-model="rejectDialog" title="拒绝提现申请" width="400px">
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="拒绝原因" required>
          <el-input
            v-model="rejectForm.reject_reason"
            type="textarea"
            :rows="3"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmReject">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getWithdrawalList, processWithdrawal, getFinanceStatistics, WithdrawalStatus } from '@/api/finance'
import { getDistributors } from '@/api/distributor'
import { formatDate } from '@/utils/date'

// 搜索表单
const searchForm = reactive({
  distributor_id: '',
  status: '',
  start_date: '',
  end_date: ''
})

const dateRange = ref([])

// 分页
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 表格数据
const tableData = ref([])
const loading = ref(false)

// 分销商列表
const distributorList = ref([])

// 统计数据
const pendingCount = ref(0)
const pendingAmount = ref(0)

// 拒绝对话框
const rejectDialog = ref(false)
const rejectForm = reactive({
  id: null,
  reject_reason: ''
})

// 获取提现申请列表
const fetchData = async () => {
  loading.value = true
  try {
    if (dateRange.value && dateRange.value.length === 2) {
      searchForm.start_date = dateRange.value[0]
      searchForm.end_date = dateRange.value[1]
    } else {
      searchForm.start_date = ''
      searchForm.end_date = ''
    }

    const res = await getWithdrawalList({
      ...searchForm,
      page: currentPage.value,
      page_size: pageSize.value
    })
    tableData.value = res.data.data
    total.value = res.data.pagination.total
  } catch (error) {
    ElMessage.error('获取提现申请列表失败')
  } finally {
    loading.value = false
  }
}

// 获取分销商列表
const fetchDistributors = async () => {
  try {
    const res = await getDistributors({ page: 1, page_size: 100 })
    distributorList.value = res.data.list  // 修复：正确的数据路径
  } catch (error) {
    ElMessage.error('获取分销商列表失败')
  }
}

// 获取统计数据
const fetchStatistics = async () => {
  try {
    const res = await getFinanceStatistics()
    pendingCount.value = res.data.pending_withdrawals.count
    pendingAmount.value = res.data.pending_withdrawals.amount
  } catch (error) {
    console.error('获取统计数据失败', error)
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

// 重置
const handleReset = () => {
  searchForm.distributor_id = ''
  searchForm.status = ''
  searchForm.start_date = ''
  searchForm.end_date = ''
  dateRange.value = []
  handleSearch()
}

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchData()
}

// 当前页改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchData()
}

// 处理提现申请
const handleProcess = async (row, approved) => {
  if (approved) {
    // 批准
    try {
      await ElMessageBox.confirm(
        `确定批准分销商【${row.distributor.name}】的提现申请吗？提现金额：¥${row.amount.toFixed(2)}`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )
      
      await processWithdrawal(row.id, { approved: true })
      ElMessage.success('提现申请已批准')
      fetchData()
      fetchStatistics()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error.response?.data?.error || '处理失败')
      }
    }
  } else {
    // 拒绝
    rejectForm.id = row.id
    rejectForm.reject_reason = ''
    rejectDialog.value = true
  }
}

// 确认拒绝
const confirmReject = async () => {
  if (!rejectForm.reject_reason) {
    ElMessage.error('请输入拒绝原因')
    return
  }
  
  try {
    await processWithdrawal(rejectForm.id, {
      approved: false,
      reject_reason: rejectForm.reject_reason
    })
    ElMessage.success('提现申请已拒绝')
    rejectDialog.value = false
    fetchData()
    fetchStatistics()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '处理失败')
  }
}

// 获取状态标签颜色
const getStatusTagType = (status) => {
  const statusMap = {
    0: 'warning',  // 待处理
    1: 'success',  // 已处理
    2: 'danger'    // 已拒绝
  }
  return statusMap[status] || ''
}

// 解析账户信息
const getAccountInfo = (accountInfoStr) => {
  try {
    return JSON.parse(accountInfoStr)
  } catch (e) {
    return null
  }
}

onMounted(() => {
  fetchData()
  fetchDistributors()
  fetchStatistics()
})
</script>

<style scoped>
.withdrawals-container {
  padding: 20px;
}

.stat-cards {
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
}

.stat-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-sub {
  font-size: 12px;
  color: #999;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  background: #fff;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>