<template>
  <div class="transactions-container">
    <!-- 搜索表单 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="分销商">
          <el-select v-model="searchForm.distributor_id" placeholder="请选择分销商" clearable filterable style="width: 200px">
            <el-option
              v-for="item in distributorList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="交易类型">
          <el-select v-model="searchForm.type" placeholder="请选择类型" clearable filterable style="width: 120px">
            <el-option label="充值" :value="1" />
            <el-option label="消费" :value="2" />
            <el-option label="提现" :value="3" />
            <el-option label="退款" :value="4" />
            <el-option label="调整" :value="5" />
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
          <el-button type="success" @click="showRechargeDialog">充值</el-button>
          <el-button type="warning" @click="showAdjustDialog">余额调整</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="tableData" stripe class="table-optimized">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="distributor.name" label="分销商" min-width="120" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">{{ TransactionType[row.type] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">
            <span :class="getAmountClass(row.type, row.amount)">
              {{ formatAmount(row.type, row.amount) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_before" label="变动前余额" width="110" align="right">
          <template #default="{ row }">
            ¥{{ row.balance_before.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="变动后余额" width="110" align="right">
          <template #default="{ row }">
            ¥{{ row.balance_after.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="related_id" label="关联单号" width="160" show-overflow-tooltip />
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="creator.real_name" label="操作人" width="90" />
        <el-table-column prop="created_at" label="创建时间" width="150">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
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

    <!-- 充值对话框 -->
    <el-dialog v-model="rechargeDialog" title="充值" width="500px">
      <el-form ref="rechargeFormRef" :model="rechargeForm" :rules="rechargeRules" label-width="100px">
        <el-form-item label="分销商" prop="distributor_id">
          <el-select v-model="rechargeForm.distributor_id" placeholder="请选择或搜索分销商" filterable style="width: 100%">
            <el-option
              v-for="item in distributorList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="充值金额" prop="amount">
          <el-input-number
            v-model="rechargeForm.amount"
            :precision="2"
            :step="100"
            :min="0.01"
            :max="100000"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="rechargeForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialog = false">取消</el-button>
        <el-button type="primary" @click="handleRecharge">确定充值</el-button>
      </template>
    </el-dialog>

    <!-- 余额调整对话框 -->
    <el-dialog v-model="adjustDialog" title="余额调整" width="500px">
      <el-form ref="adjustFormRef" :model="adjustForm" :rules="adjustRules" label-width="100px">
        <el-form-item label="分销商" prop="distributor_id">
          <el-select v-model="adjustForm.distributor_id" placeholder="请选择或搜索分销商" filterable style="width: 100%">
            <el-option
              v-for="item in distributorList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="调整金额" prop="amount">
          <el-input-number
            v-model="adjustForm.amount"
            :precision="2"
            :step="100"
            style="width: 100%"
          />
          <div class="el-form-item__error" style="position: static; margin-top: 5px;">
            正数为增加余额，负数为减少余额
          </div>
        </el-form-item>
        <el-form-item label="调整原因" prop="remark">
          <el-input
            v-model="adjustForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入调整原因（必填）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAdjust">确定调整</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTransactionList, recharge, adjustBalance, TransactionType } from '@/api/finance'
import { getDistributors } from '@/api/distributor'
import { formatDate } from '@/utils/date'

// 搜索表单
const searchForm = reactive({
  distributor_id: '',
  type: '',
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

// 充值对话框
const rechargeDialog = ref(false)
const rechargeFormRef = ref()
const rechargeForm = reactive({
  distributor_id: '',
  amount: 100,
  remark: ''
})

const rechargeRules = {
  distributor_id: [
    { required: true, message: '请选择分销商', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入充值金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '充值金额必须大于0', trigger: 'blur' }
  ]
}

// 余额调整对话框
const adjustDialog = ref(false)
const adjustFormRef = ref()
const adjustForm = reactive({
  distributor_id: '',
  amount: 0,
  remark: ''
})

const adjustRules = {
  distributor_id: [
    { required: true, message: '请选择分销商', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入调整金额', trigger: 'blur' },
    { type: 'number', message: '请输入有效的金额', trigger: 'blur' }
  ],
  remark: [
    { required: true, message: '请输入调整原因', trigger: 'blur' }
  ]
}

// 获取交易记录列表
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

    const res = await getTransactionList({
      ...searchForm,
      page: currentPage.value,
      page_size: pageSize.value
    })
    tableData.value = res.data.data
    total.value = res.data.pagination.total
  } catch (error) {
    ElMessage.error('获取交易记录失败')
  } finally {
    loading.value = false
  }
}

// 获取分销商列表
const fetchDistributors = async () => {
  try {
    console.log('正在获取分销商列表...')
    const res = await getDistributors({ page: 1, page_size: 100 })
    console.log('分销商API响应:', res)
    console.log('分销商数据:', res.data)
    distributorList.value = res.data.list  // 修复：正确的数据路径
    console.log('设置的分销商列表:', distributorList.value)
  } catch (error) {
    console.error('获取分销商列表失败:', error)
    ElMessage.error('获取分销商列表失败')
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
  searchForm.type = ''
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

// 显示充值对话框
const showRechargeDialog = () => {
  rechargeForm.distributor_id = ''
  rechargeForm.amount = 100
  rechargeForm.remark = ''
  rechargeDialog.value = true
}

// 处理充值
const handleRecharge = async () => {
  await rechargeFormRef.value.validate()
  
  try {
    await recharge(rechargeForm)
    ElMessage.success('充值成功')
    rechargeDialog.value = false
    fetchData()
    fetchDistributors() // 刷新分销商列表以更新余额
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '充值失败')
  }
}

// 显示余额调整对话框
const showAdjustDialog = () => {
  adjustForm.distributor_id = ''
  adjustForm.amount = 0
  adjustForm.remark = ''
  adjustDialog.value = true
}

// 处理余额调整
const handleAdjust = async () => {
  await adjustFormRef.value.validate()
  
  try {
    await adjustBalance(adjustForm)
    ElMessage.success('余额调整成功')
    adjustDialog.value = false
    fetchData()
    fetchDistributors() // 刷新分销商列表以更新余额
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '余额调整失败')
  }
}

// 获取类型标签颜色
const getTypeTagType = (type) => {
  const typeMap = {
    1: 'success', // 充值
    2: 'danger',  // 消费
    3: 'warning', // 提现
    4: 'info',    // 退款
    5: 'info'     // 调整
  }
  return typeMap[type] || 'info'
}

// 获取金额样式类
const getAmountClass = (type, amount) => {
  if (type === 1 || type === 4) { // 充值、退款
    return 'amount-in'
  } else if (type === 2 || type === 3) { // 消费、提现
    return 'amount-out'
  } else if (type === 5) { // 调整
    return amount >= 0 ? 'amount-in' : 'amount-out'
  }
  return ''
}

// 格式化金额显示
const formatAmount = (type, amount) => {
  let prefix = ''
  if (type === 1 || type === 4) { // 充值、退款
    prefix = '+'
  } else if (type === 2 || type === 3) { // 消费、提现
    prefix = '-'
  } else if (type === 5) { // 调整
    prefix = amount >= 0 ? '+' : ''
  }
  return `${prefix}¥${Math.abs(amount).toFixed(2)}`
}

onMounted(() => {
  console.log('财务页面已挂载，开始获取数据...')
  fetchData()
  fetchDistributors()
})
</script>

<style scoped>
.transactions-container {
  padding: 20px;
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

.amount-in {
  color: #67c23a;
  font-weight: bold;
}

.amount-out {
  color: #f56c6c;
  font-weight: bold;
}
</style>