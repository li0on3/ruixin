<template>
  <div class="page-container">
    <!-- 搜索表单 -->
    <div class="search-form search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="批次号">
          <el-input
            v-model="searchForm.batchNo"
            placeholder="请输入批次号"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="价格ID">
          <el-input
            v-model="searchForm.priceId"
            placeholder="请输入价格ID"
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

    <!-- 表格 -->
    <div class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="batch_no" label="批次号" width="180" show-overflow-tooltip />
        <el-table-column label="价格信息" width="120">
          <template #default="{ row }">
            <div v-if="row.price">
              ID: {{ row.price.price_id }}<br>
              ¥{{ row.price.price_value }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="cost_price" label="成本价" width="90" align="right">
          <template #default="{ row }">
            ¥{{ row.cost_price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="sell_price" label="销售价" width="90" align="right">
          <template #default="{ row }">
            ¥{{ row.sell_price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="使用情况" width="150" align="center">
          <template #default="{ row }">
            <div>
              {{ row.used_count }} / {{ row.total_count }}
              <el-progress
                :percentage="getUsageRate(row)"
                :status="getProgressStatus(row)"
                :stroke-width="6"
                style="margin-top: 5px"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column label="利润" width="120" align="right">
          <template #default="{ row }">
            <div>
              单张: ¥{{ (row.sell_price - row.cost_price).toFixed(2) }}<br>
              总计: ¥{{ ((row.sell_price - row.cost_price) * row.used_count).toFixed(2) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="imported_at" label="导入时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.imported_at) }}
          </template>
        </el-table-column>
        <el-table-column label="导入人" width="100">
          <template #default="{ row }">
            {{ row.admin?.real_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleViewCards(row)">
              查看卡片
            </el-button>
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

    <!-- 卡片列表对话框 -->
    <el-dialog
      v-model="cardsDialogVisible"
      :title="`批次卡片列表 - ${currentBatch?.batch_no}`"
      width="900px"
    >
      <el-table :data="cardsData" v-loading="cardsLoading" max-height="500">
        <el-table-column prop="card_code" label="卡片代码" width="150" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'">
              {{ row.status === 0 ? '未使用' : '已使用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="used_at" label="使用时间" width="160">
          <template #default="{ row }">
            {{ row.used_at ? formatDate(row.used_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="order_id" label="订单ID" width="100" />
        <el-table-column prop="expired_at" label="过期时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.expired_at) }}
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 15px; text-align: center;">
        <el-pagination
          v-model:current-page="cardsPage"
          :page-size="cardsPageSize"
          :total="cardsTotal"
          layout="prev, pager, next"
          @current-change="loadBatchCards"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getBatches, getBatchCards } from '@/api/card'
import { formatDate } from '@/utils/date'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  batchNo: '',
  priceId: ''
})

// 卡片列表相关
const cardsDialogVisible = ref(false)
const cardsLoading = ref(false)
const cardsData = ref([])
const cardsTotal = ref(0)
const cardsPage = ref(1)
const cardsPageSize = 20
const currentBatch = ref(null)

// 获取使用率
const getUsageRate = (row) => {
  if (row.total_count === 0) return 0
  return Math.round((row.used_count / row.total_count) * 100)
}

// 获取进度条状态
const getProgressStatus = (row) => {
  const rate = getUsageRate(row)
  if (rate >= 90) return 'danger'
  if (rate >= 70) return 'warning'
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

    if (searchForm.priceId) {
      params.price_id = searchForm.priceId
    }

    const { data } = await getBatches(params)
    tableData.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error('Failed to load batches:', error)
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
  searchForm.batchNo = ''
  searchForm.priceId = ''
  handleSearch()
}

// 查看批次卡片
const handleViewCards = (row) => {
  currentBatch.value = row
  cardsPage.value = 1
  cardsDialogVisible.value = true
  loadBatchCards()
}

// 加载批次卡片
const loadBatchCards = async () => {
  if (!currentBatch.value) return

  cardsLoading.value = true
  try {
    const { data } = await getBatchCards(currentBatch.value.id, {
      page: cardsPage.value,
      page_size: cardsPageSize
    })
    cardsData.value = data.list || []
    cardsTotal.value = data.total || 0
  } catch (error) {
    console.error('Failed to load batch cards:', error)
  } finally {
    cardsLoading.value = false
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
</style>