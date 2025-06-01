<template>
  <div class="products-page">
    <!-- 页面标题和操作栏 -->
    <div class="page-header">
      <div>
        <h2>商品管理</h2>
        <p class="page-desc">
          <template v-if="cardCode">
            <el-tag type="primary" closable @close="clearCardFilter">卡片: {{ cardCode }}</el-tag>
            当前查看该卡片的关联产品
          </template>
          <template v-else>
            管理所有可下单的产品信息和瑞幸产品配置
          </template>
        </p>
      </div>
      <div class="actions">
        <el-button type="primary" @click="showSyncDialog = true">
          <el-icon><Refresh /></el-icon>
          同步产品
        </el-button>
        <el-button @click="showMatchLogsDialog = true">
          <el-icon><Document /></el-icon>
          匹配失败日志
        </el-button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.native.prevent="loadProducts">
        <el-form-item label="搜索商品">
          <el-input
            v-model="searchForm.keyword"
            placeholder="产品名称/代码"
            clearable
            @keyup.enter="loadProducts"
          />
        </el-form-item>
        <el-form-item label="类别">
          <el-select v-model="searchForm.category" placeholder="请选择类别" clearable filterable style="width: 150px">
            <el-option
              v-for="item in categoryList"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable filterable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadProducts">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 产品列表 -->
    <el-card>
      <el-table
        :data="products"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="product_id" label="产品ID" width="120" />
        <el-table-column prop="goods_code" label="产品代码" width="120" />
        <el-table-column prop="name" label="产品名称" min-width="200">
          <template #default="{ row }">
            <div>
              <span>{{ row.name || row.goods_name }}</span>
              <el-tag 
                v-if="row.aliases && row.aliases.length > 0"
                size="small"
                type="info"
                style="margin-left: 8px"
              >
                有{{ row.aliases.length }}个别名
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="类别" width="120" />
        <el-table-column prop="image_url" label="图片" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.image_url"
              :src="row.image_url"
              :preview-src-list="[row.image_url]"
              fit="cover"
              style="width: 50px; height: 50px"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="关联卡片" min-width="180">
          <template #default="{ row }">
            <div v-if="row.card_count && row.card_count > 0" class="card-bindings">
              <el-tag
                type="primary"
                size="small"
              >
                关联 {{ row.card_count }} 张卡片
              </el-tag>
              <el-button
                link
                size="small"
                @click="showCardBindings(row)"
                style="margin-left: 8px"
              >
                查看详情
              </el-button>
            </div>
            <span v-else class="text-gray">暂无关联</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at || row.last_sync_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              size="small"
              @click="showProductDetail(row)"
            >
              查看详情
            </el-button>
            <el-button
              link
              size="small"
              @click="editProduct(row)"
            >
              编辑
            </el-button>
            <el-button
              link
              size="small"
              @click="deleteProduct(row)"
              style="color: #f56c6c"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadProducts"
        @size-change="loadProducts"
      />
    </el-card>

    <!-- 同步产品对话框 -->
    <el-dialog
      v-model="showSyncDialog"
      title="同步产品"
      width="500px"
    >
      <el-form :model="syncForm" label-width="100px">
        <el-form-item label="卡片代码" required>
          <el-select
            v-model="syncForm.card_code"
            placeholder="选择或输入卡片代码"
            filterable
            allow-create
            :disabled="!!cardCode"
            style="width: 100%"
            @focus="loadUnusedCards"
          >
            <el-option
              v-for="card in unusedCards"
              :key="card.id"
              :label="`${card.card_code}${card.description ? ' - ' + card.description : ''}`"
              :value="card.card_code"
            >
              <div class="card-option">
                <span class="card-code">{{ card.card_code }}</span>
                <span v-if="card.description" class="card-desc">{{ card.description }}</span>
                <el-tag v-if="card.cost_price" type="info" size="small" style="margin-left: 8px">
                  成本: ¥{{ card.cost_price }}
                </el-tag>
              </div>
            </el-option>
          </el-select>
          <div v-if="cardCode" class="form-tip">
            已自动填充当前卡片代码
          </div>
          <div v-else class="form-tip">
            可选择现有未使用的卡片，或输入新的卡片代码
          </div>
        </el-form-item>
        <el-form-item label="店铺代码">
          <el-input
            v-model="syncForm.store_code"
            placeholder="留空使用系统配置"
          />
          <div class="form-tip">
            当前系统配置店铺：{{ systemStoreCode || '未配置' }}
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showSyncDialog = false">取消</el-button>
        <el-button
          type="primary"
          :loading="syncing"
          @click="handleSyncProducts"
        >
          开始同步
        </el-button>
      </template>
    </el-dialog>

    <!-- 添加/编辑产品对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="产品ID" prop="product_id">
          <el-input v-model="form.product_id" placeholder="请输入产品ID" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="产品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入产品名称" />
        </el-form-item>
        <el-form-item label="类别" prop="category">
          <el-select v-model="form.category" placeholder="请选择或输入类别" filterable allow-create style="width: 100%">
            <el-option
              v-for="item in categoryList"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="图片URL" prop="image_url">
          <el-input v-model="form.image_url" placeholder="请输入图片URL" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入产品描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="importDialog" title="批量导入产品" width="800px">
      <el-alert type="info" show-icon :closable="false" style="margin-bottom: 20px">
        <template #title>
          导入说明：请按照JSON格式导入产品数据，每个产品需要包含 product_id、name、category 等字段
        </template>
      </el-alert>
      <el-input
        v-model="importData"
        type="textarea"
        :rows="10"
        placeholder='[{"product_id":"2500","name":"标准美式","description":"经典美式咖啡","category":"美式家族","image_url":""}]'
      />
      <template #footer>
        <el-button @click="importDialog = false">取消</el-button>
        <el-button type="primary" @click="handleImport">导入</el-button>
      </template>
    </el-dialog>

    <!-- 商品详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="商品详情"
      width="800px"
    >
      <div v-if="currentProduct">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="产品ID">
            {{ currentProduct.product_id || currentProduct.goods_code }}
          </el-descriptions-item>
          <el-descriptions-item label="产品名称">
            {{ currentProduct.name || currentProduct.goods_name }}
          </el-descriptions-item>
          <el-descriptions-item label="类别">
            {{ currentProduct.category || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentProduct.status === 1 ? 'success' : 'danger'" size="small">
              {{ currentProduct.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间" :span="2">
            {{ formatDate(currentProduct.created_at || currentProduct.last_sync_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ currentProduct.description || '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 可用规格信息 -->
        <div v-if="currentProduct.available_specs" style="margin-top: 20px">
          <h4>可用规格</h4>
          <div v-for="(spec, specCode) in parseAvailableSpecs(currentProduct.available_specs)" :key="specCode" style="margin-bottom: 10px">
            <el-tag type="info">{{ spec.name }}</el-tag>
            <span style="margin-left: 10px">
              <el-tag 
                v-for="item in spec.items" 
                :key="item.code"
                :type="item.is_default ? 'primary' : ''"
                size="small"
                style="margin: 2px"
              >
                {{ item.name }}{{ item.is_default ? ' (默认)' : '' }}
              </el-tag>
            </span>
          </div>
        </div>

        <!-- SKU信息 -->
        <h4 style="margin-top: 20px">SKU信息</h4>
        <el-table
          :data="currentProduct.SKUs || currentProduct.skus || []"
          size="small"
          max-height="300"
          empty-text="暂无SKU数据"
        >
          <el-table-column prop="sku_code" label="SKU代码" width="120">
            <template #default="{ row }">
              {{ row.sku_code }}
              <el-button
                link
                size="small"
                @click="copyToClipboard(row.sku_code)"
                style="margin-left: 5px"
              >
                <el-icon><Document /></el-icon>
              </el-button>
            </template>
          </el-table-column>
          <el-table-column prop="sku_name" label="SKU名称" min-width="150">
            <template #default="{ row }">
              <el-tooltip 
                v-if="row.sku_name && row.sku_name.length > 20"
                :content="row.sku_name"
                placement="top"
              >
                <span>{{ row.sku_name }}</span>
              </el-tooltip>
              <span v-else>{{ row.sku_name || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="chinese_desc" label="中文描述" min-width="150">
            <template #default="{ row }">
              <span v-if="row.chinese_desc">{{ row.chinese_desc }}</span>
              <span v-else-if="row.mapping && row.mapping.chinese_desc">{{ row.mapping.chinese_desc }}</span>
              <span v-else class="text-gray">-</span>
            </template>
          </el-table-column>
          <el-table-column label="规格详情" min-width="180">
            <template #default="{ row }">
              <el-space v-if="row.specs" wrap>
                <el-tag v-if="row.specs.size" size="small">
                  {{ row.specs.size.name }}
                </el-tag>
                <el-tag v-if="row.specs.temperature" size="small" type="warning">
                  {{ row.specs.temperature.name }}
                </el-tag>
                <el-tag v-if="row.specs.sweetness" size="small" type="success">
                  {{ row.specs.sweetness.name }}
                </el-tag>
                <el-tag v-if="row.specs.milk" size="small" type="info">
                  {{ row.specs.milk.name }}
                </el-tag>
              </el-space>
              <span v-else class="text-gray">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="specs_code" label="规格编码" width="100">
            <template #default="{ row }">
              {{ row.specs_code || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="默认" width="60" align="center">
            <template #default="{ row }">
              <el-tag v-if="row.is_default" type="primary" size="small">默认</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <!-- 卡片绑定详情对话框 -->
    <el-dialog
      v-model="showCardBindingsDialog"
      title="关联卡片详情"
      width="600px"
    >
      <div v-if="currentProduct">
        <div class="binding-header">
          <strong>产品：</strong>{{ currentProduct.name || currentProduct.goods_name }} ({{ currentProduct.product_id || currentProduct.goods_code }})
        </div>
        <el-table
          :data="currentCardBindings"
          v-loading="bindingsLoading"
          size="small"
        >
          <el-table-column prop="card_code" label="卡片代码" width="150" />
          <el-table-column prop="price_name" label="价格配置">
            <template #default="{ row }">
              {{ row.price_name || `配置${row.card_id}` }}
            </template>
          </el-table-column>
          <el-table-column prop="price" label="面值" width="100">
            <template #default="{ row }">
              ¥{{ row.price }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 0 ? 'success' : 'danger'" size="small">
                {{ row.status === 0 ? '可用' : '已用' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <!-- 匹配失败日志对话框 -->
    <el-dialog
      v-model="showMatchLogsDialog"
      title="匹配失败日志"
      width="900px"
    >
      <el-table
        :data="matchLogs"
        v-loading="logsLoading"
        size="small"
      >
        <el-table-column prop="distributor.name" label="分销商" width="120" />
        <el-table-column prop="request_time" label="请求时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.request_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="input_product" label="输入商品" />
        <el-table-column prop="input_specs" label="输入规格">
          <template #default="{ row }">
            <div v-if="row.input_specs">
              {{ formatSpecs(row.input_specs) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="error_reason" label="错误原因" />
        <el-table-column prop="suggestions" label="建议" width="200">
          <template #default="{ row }">
            <div v-if="row.suggestions && row.suggestions.length > 0">
              <el-tag
                v-for="(suggestion, index) in row.suggestions"
                :key="index"
                size="small"
                style="margin: 2px"
              >
                {{ suggestion }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="logsPage"
        :page-size="20"
        :total="logsTotal"
        layout="total, prev, pager, next"
        @current-change="loadMatchLogs"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Document } from '@element-plus/icons-vue'
import { getProducts, syncProducts, getMatchLogs } from '@/api/product'
import { getProductList, getProductCategories } from '@/api/luckinConfig'
import { getSystemConfigs } from '@/api/system'
import { getCards, getUnusedCards } from '@/api/card'
import { formatDate } from '@/utils/date'
import request from '@/api/request'

const route = useRoute()

// 获取查询参数
const cardId = ref(route.query.card_id || '')
const cardCode = ref(route.query.card_code || '')

// 数据
const products = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const searchForm = reactive({
  keyword: '',
  category: '',
  status: '',
  card_id: cardId.value
})

// 类别列表
const categoryList = ref([])

// 未使用的卡片列表
const unusedCards = ref([])

// 同步相关
const showSyncDialog = ref(false)
const syncing = ref(false)
const syncForm = ref({
  card_code: cardCode.value || '', // 从URL获取卡片代码
  store_code: ''
})
const systemStoreCode = ref('')

// 产品管理对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const formRef = ref()
const form = reactive({
  product_id: '',
  name: '',
  category: '',
  image_url: '',
  description: '',
  status: 1
})

const rules = {
  product_id: [
    { required: true, message: '请输入产品ID', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入产品名称', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择类别', trigger: 'change' }
  ]
}

// 批量导入对话框
const importDialog = ref(false)
const importData = ref('')

// 详情相关
const showDetailDialog = ref(false)
const currentProduct = ref(null)

// 卡片绑定相关
const showCardBindingsDialog = ref(false)
const currentCardBindings = ref([])
const bindingsLoading = ref(false)

// 匹配日志相关
const showMatchLogsDialog = ref(false)
const matchLogs = ref([])
const logsLoading = ref(false)
const logsPage = ref(1)
const logsTotal = ref(0)

// 方法
const loadProducts = async () => {
  loading.value = true
  try {
    // 优先加载瑞幸产品配置数据
    const luckinParams = {
      page: currentPage.value,
      page_size: pageSize.value,
      name: searchForm.keyword,
      category: searchForm.category,
      status: searchForm.status
    }
    
    const luckinRes = await getProductList(luckinParams)
    
    if (luckinRes.data.data && luckinRes.data.data.length > 0) {
      products.value = luckinRes.data.data
      total.value = luckinRes.data.pagination.total
    } else {
      // 如果瑞幸产品为空，尝试加载旧的产品数据
      const params = {
        page: currentPage.value,
        page_size: pageSize.value,
        search: searchForm.keyword
      }
      
      if (searchForm.card_id) {
        params.card_id = searchForm.card_id
      }
      
      const res = await getProducts(params)
      products.value = res.data.list
      total.value = res.data.total
    }
  } catch (error) {
    ElMessage.error('加载产品失败')
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.category = ''
  searchForm.status = ''
  currentPage.value = 1
  loadProducts()
}

// 清除卡片筛选
const clearCardFilter = () => {
  searchForm.card_id = ''
  cardId.value = ''
  cardCode.value = ''
  currentPage.value = 1
  loadProducts()
}

const handleSyncProducts = async () => {
  if (!syncForm.value.card_code) {
    ElMessage.warning('请输入卡片代码')
    return
  }

  syncing.value = true
  try {
    const res = await syncProducts({
      card_code: syncForm.value.card_code,
      store_code: syncForm.value.store_code
    })
    
    ElMessage.success(
      `同步成功！同步${res.data.synced_count}个产品，` +
      `新增${res.data.new_count}个，更新${res.data.updated_count}个`
    )
    
    showSyncDialog.value = false
    syncForm.value = { card_code: '', store_code: '' }
    loadProducts()
  } catch (error) {
    ElMessage.error(error.response?.data?.msg || '同步失败')
  } finally {
    syncing.value = false
  }
}

// 产品管理功能
const showAddProductDialog = () => {
  dialogTitle.value = '添加产品'
  isEdit.value = false
  Object.assign(form, {
    product_id: '',
    name: '',
    category: '',
    image_url: '',
    description: '',
    status: 1
  })
  dialogVisible.value = true
}

const editProduct = (product) => {
  dialogTitle.value = '编辑产品'
  isEdit.value = true
  Object.assign(form, {
    id: product.id,
    product_id: product.product_id,
    name: product.name,
    category: product.category,
    image_url: product.image_url,
    description: product.description,
    status: product.status
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  
  try {
    if (isEdit.value) {
      await updateProduct(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createProduct(form)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    loadProducts()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || (isEdit.value ? '更新失败' : '添加失败'))
  }
}

const deleteProduct = async (product) => {
  try {
    await ElMessageBox.confirm(
      `确定删除产品【${product.name}】吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteProductApi(product.id)
    ElMessage.success('删除成功')
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const handleImport = async () => {
  if (!importData.value.trim()) {
    ElMessage.warning('请输入导入数据')
    return
  }
  
  try {
    const data = JSON.parse(importData.value)
    await batchImportProducts(data)
    ElMessage.success('导入成功')
    importDialog.value = false
    importData.value = ''
    loadProducts()
  } catch (error) {
    if (error instanceof SyntaxError) {
      ElMessage.error('JSON格式错误，请检查数据格式')
    } else {
      ElMessage.error(error.response?.data?.error || '导入失败')
    }
  }
}

const showProductDetail = (product) => {
  currentProduct.value = product
  showDetailDialog.value = true
}

// 显示卡片绑定详情
const showCardBindings = async (product) => {
  currentProduct.value = product
  showCardBindingsDialog.value = true
  bindingsLoading.value = true
  
  try {
    // 调用API获取产品的卡片绑定详情
    const res = await request({
      url: '/admin/products/cards',
      method: 'get',
      params: {
        goods_code: product.goods_code || product.product_id
      }
    })
    
    if (res.data.list && res.data.list.length > 0) {
      // 获取该产品的所有绑定卡片
      currentCardBindings.value = res.data.list[0].cards || []
    }
  } catch (error) {
    ElMessage.error('获取卡片绑定详情失败')
    currentCardBindings.value = []
  } finally {
    bindingsLoading.value = false
  }
}

const loadMatchLogs = async () => {
  logsLoading.value = true
  try {
    const res = await getMatchLogs({
      page: logsPage.value,
      page_size: 20
    })
    matchLogs.value = res.data.list
    logsTotal.value = res.data.total
  } catch (error) {
    ElMessage.error('加载日志失败')
  } finally {
    logsLoading.value = false
  }
}

const formatSpecs = (specs) => {
  if (typeof specs === 'string') {
    try {
      specs = JSON.parse(specs)
    } catch {
      return specs
    }
  }
  return Object.entries(specs).map(([k, v]) => `${k}:${v}`).join(', ')
}

const parseAvailableSpecs = (specsStr) => {
  if (!specsStr) return {}
  try {
    return JSON.parse(specsStr)
  } catch {
    return {}
  }
}

// 复制到剪贴板功能
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (err) {
    // 降级方案
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    ElMessage.success('已复制到剪贴板')
  }
}

// 获取系统配置
const loadSystemConfig = async () => {
  try {
    const res = await getSystemConfigs()
    if (res.code === 200) {
      systemStoreCode.value = res.data.sync_store_code || '390840'
    }
  } catch (error) {
    console.error('获取系统配置失败:', error)
  }
}

// 获取类别列表
const fetchCategories = async () => {
  try {
    const res = await getProductCategories()
    categoryList.value = res.data || []
  } catch (error) {
    console.error('获取类别失败:', error)
  }
}

// 加载未使用的卡片列表
const loadUnusedCards = async () => {
  if (unusedCards.value.length > 0) return // 避免重复加载
  
  try {
    const res = await getUnusedCards()
    if (res.code === 200 && res.data.list) {
      unusedCards.value = res.data.list
    }
  } catch (error) {
    console.error('获取未使用卡片失败:', error)
  }
}

// 生命周期
onMounted(() => {
  loadProducts()
  loadSystemConfig()
  fetchCategories()
})
</script>

<style lang="scss" scoped>
.products-page {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 24px;
      font-weight: 500;
    }

    .actions {
      display: flex;
      gap: 10px;
    }
  }

  .search-card {
    margin-bottom: 20px;
  }

  .text-gray {
    color: #999;
  }

  .form-tip {
    margin-top: 5px;
    font-size: 12px;
    color: #999;
  }

  .binding-header {
    margin-bottom: 15px;
    padding: 10px;
    background: #f5f5f5;
    border-radius: 4px;
  }
  
  .page-desc {
    margin-top: 8px;
    font-size: 14px;
    color: #666;
    
    .el-tag {
      margin-right: 8px;
    }
  }
  
  .card-bindings {
    display: flex;
    align-items: center;
  }
  
  // 卡片选择器样式
  .card-option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    
    .card-code {
      font-weight: 500;
      color: #303133;
    }
    
    .card-desc {
      color: #909399;
      font-size: 12px;
      margin-left: 8px;
      flex: 1;
    }
  }
}
</style>