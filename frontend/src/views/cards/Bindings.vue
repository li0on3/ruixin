<template>
  <div class="bindings-container">
    <!-- 页面标题 -->
    <el-page-header @back="goBack" style="margin-bottom: 20px">
      <template #content>
        <span class="text-large font-600 mr-3">
          卡片绑定管理 - {{ cardName }}
        </span>
      </template>
    </el-page-header>

    <!-- 工具栏 -->
    <el-card class="toolbar-card search-card">
      <el-button type="primary" @click="showAddDialog">
        <el-icon><Plus /></el-icon>
        添加绑定
      </el-button>
      <el-button @click="fetchOptions">刷新选项</el-button>
    </el-card>

    <!-- 绑定列表 -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="target_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.target_type === 'price' ? 'success' : 'primary'">
              {{ row.target_type === 'price' ? '价格' : '产品' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target_id" label="目标" min-width="200" />
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
            <el-input-number
              v-model="row.priority"
              :min="0"
              :max="999"
              size="small"
              @change="updatePriority(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="creator.real_name" label="创建人" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" size="small" text @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加绑定对话框 -->
    <el-dialog v-model="dialogVisible" title="添加绑定" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="绑定类型" prop="target_type">
          <el-radio-group v-model="form.target_type" @change="handleTypeChange">
            <el-radio label="price">价格</el-radio>
            <el-radio label="product">产品</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="选择目标" prop="target_id">
          <el-select
            v-if="form.target_type === 'price'"
            v-model="form.target_id"
            placeholder="请选择价格"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="item in priceOptions"
              :key="item.price_id"
              :label="`${item.price_id} - ¥${item.price_value}`"
              :value="item.price_id"
            />
          </el-select>
          <el-select
            v-else
            v-model="form.target_id"
            placeholder="请选择产品"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="item in productOptions"
              :key="item.product_id"
              :label="`${item.product_id} - ${item.name}`"
              :value="item.product_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number
            v-model="form.priority"
            :min="0"
            :max="999"
            style="width: 100%"
          />
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
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getCategoryBindings, createBinding, deleteBinding, updateBindingPriority, getActiveOptions } from '@/api/luckinConfig'
import { formatDate } from '@/utils/date'

const route = useRoute()
const router = useRouter()

const categoryId = ref(route.query.id)
const cardName = ref(route.query.name || '')

// 表格数据
const loading = ref(false)
const tableData = ref([])

// 选项数据
const priceOptions = ref([])
const productOptions = ref([])

// 对话框
const dialogVisible = ref(false)
const formRef = ref()

// 表单数据
const form = reactive({
  target_type: 'price',
  target_id: '',
  priority: 0
})

// 表单验证规则
const rules = {
  target_id: [
    { required: true, message: '请选择目标', trigger: 'change' }
  ]
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 获取绑定列表
const fetchData = async () => {
  if (!categoryId.value) return
  
  loading.value = true
  try {
    const res = await getCategoryBindings(categoryId.value)
    tableData.value = res.data
  } catch (error) {
    ElMessage.error('获取绑定列表失败')
  } finally {
    loading.value = false
  }
}

// 获取选项数据
const fetchOptions = async () => {
  try {
    const res = await getActiveOptions()
    priceOptions.value = res.data.prices || []
    productOptions.value = res.data.products || []
  } catch (error) {
    ElMessage.error('获取选项数据失败')
  }
}

// 显示添加对话框
const showAddDialog = () => {
  form.target_type = 'price'
  form.target_id = ''
  form.priority = 0
  dialogVisible.value = true
}

// 处理类型切换
const handleTypeChange = () => {
  form.target_id = ''
}

// 提交表单
const handleSubmit = async () => {
  await formRef.value.validate()
  
  try {
    await createBinding(categoryId.value, form)
    ElMessage.success('绑定创建成功')
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '创建失败')
  }
}

// 更新优先级
const updatePriority = async (row) => {
  try {
    await updateBindingPriority(row.id, row.priority)
    ElMessage.success('优先级更新成功')
  } catch (error) {
    ElMessage.error('更新失败')
    fetchData() // 重新加载数据
  }
}

// 删除绑定
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个绑定吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteBinding(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

onMounted(() => {
  if (!categoryId.value) {
    ElMessage.error('缺少卡片ID参数')
    router.back()
    return
  }
  
  fetchData()
  fetchOptions()
})
</script>

<style scoped>
.bindings-container {
  padding: 20px;
}

.toolbar-card {
  margin-bottom: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  background: #fff;
}
</style>