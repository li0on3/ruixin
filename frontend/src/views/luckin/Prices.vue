<template>
  <div class="prices-container">
    <!-- 工具栏 -->
    <div class="search-card">
      <el-button type="primary" @click="showAddDialog">
        <el-icon><Plus /></el-icon>
        添加价格
      </el-button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="price_id" label="价格ID" width="120" />
        <el-table-column prop="price_value" label="价格值" width="120">
          <template #default="{ row }">
            ¥{{ row.price_value.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="关联商品" min-width="200">
          <template #default="{ row }">
            <div v-if="row.product_codes">
              <el-tag 
                v-for="code in row.product_codes.split(',').slice(0, 3)" 
                :key="code"
                size="small"
                style="margin: 2px"
              >
                {{ code }}
              </el-tag>
              <span v-if="row.product_codes.split(',').length > 3">
                等{{ row.product_codes.split(',').length }}个商品
              </span>
            </div>
            <span v-else style="color: #999">未配置</span>
          </template>
        </el-table-column>
        <el-table-column label="卡片统计" width="150">
          <template #default="{ row }">
            <div v-if="row.card_stats">
              可用: {{ row.card_stats.available_count || 0 }}<br>
              已使用: {{ row.card_stats.used_count || 0 }}
            </div>
            <el-button
              v-else
              link
              size="small"
              @click="loadCardStats(row)"
              :loading="row.loading"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="利润统计" width="150">
          <template #default="{ row }">
            <div v-if="row.card_stats">
              总利润: ￥{{ (row.card_stats.total_profit || 0).toFixed(2) }}<br>
              利润率: {{ (row.card_stats.profit_margin || 0).toFixed(1) }}%
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="creator.real_name" label="创建人" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleViewCards(row)">
              查看卡片
            </el-button>
            <el-button type="primary" size="small" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">
              删除
            </el-button>
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
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="价格ID" prop="price_id">
          <el-input v-model="form.price_id" placeholder="请输入价格ID" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="价格值" prop="price_value">
          <el-input-number
            v-model="form.price_value"
            :precision="2"
            :step="0.1"
            :min="0.01"
            :max="999.99"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
    
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { getPriceList, createPrice, updatePrice, deletePrice } from '@/api/luckinConfig'
import { getCardStats } from '@/api/card'
import { formatDate } from '@/utils/date'
import request from '@/api/request'

// 表格数据
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const formRef = ref()
const currentEditId = ref(null)

// 表单数据
const form = reactive({
  price_id: '',
  price_value: 0,
  status: 1
})


// 表单验证规则
const rules = {
  price_id: [
    { required: true, message: '请输入价格ID', trigger: 'blur' }
  ],
  price_value: [
    { required: true, message: '请输入价格值', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '价格值必须大于0', trigger: 'blur' }
  ]
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPriceList({
      page: currentPage.value,
      page_size: pageSize.value
    })
    tableData.value = res.data.list
    total.value = res.data.pagination.total
  } catch (error) {
    ElMessage.error('获取价格列表失败')
  } finally {
    loading.value = false
  }
}

// 显示添加对话框
const showAddDialog = () => {
  isEdit.value = false
  dialogTitle.value = '添加价格'
  currentEditId.value = null
  form.price_id = ''
  form.price_value = 0
  form.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑价格'
  currentEditId.value = row.id
  form.price_id = row.price_id
  form.price_value = row.price_value
  form.status = row.status
  
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  await formRef.value.validate()
  
  try {
    if (isEdit.value) {
      await updatePrice(currentEditId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createPrice(form)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('Submit error:', error.response?.data)
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除价格"${row.price_id}"吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deletePrice(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

// 加载卡片统计
const loadCardStats = async (row) => {
  row.loading = true
  try {
    const { data } = await getCardStats({ price_id: row.id })
    row.card_stats = data
  } catch (error) {
    console.error('Failed to load card stats:', error)
  } finally {
    row.loading = false
  }
}

// 查看卡片
const handleViewCards = (row) => {
  // 跳转到卡片管理页面，带上价格ID筛选
  window.location.href = `/cards?price_id=${row.id}`
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

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.prices-container {
  padding: 20px;
}


.table-card {
  background: #fff;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>