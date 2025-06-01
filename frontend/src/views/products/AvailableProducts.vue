<template>
  <div class="available-products-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <div class="header-left">
        <h2>可用商品查询</h2>
        <el-tag type="info" effect="plain">
          数据每5分钟自动更新
        </el-tag>
      </div>
      <div class="header-right">
        <el-button @click="handleRefresh" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
        <el-button type="primary" @click="handleExport">
          <el-icon><Download /></el-icon>
          导出Excel
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <el-card class="search-card">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-input
            v-model="searchText"
            placeholder="搜索商品名称、编码或别名"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="filterStatus"
            placeholder="库存状态"
            clearable
            filterable
            @change="handleFilter"
          >
            <el-option label="全部" value="" />
            <el-option label="库存充足" value="high" />
            <el-option label="库存一般" value="medium" />
            <el-option label="库存较少" value="low" />
            <el-option label="暂时缺货" value="out" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="filterCategory"
            placeholder="商品类别"
            clearable
            filterable
            @change="handleFilter"
          >
            <el-option label="全部" value="" />
            <el-option label="咖啡" value="咖啡" />
            <el-option label="茶饮" value="茶饮" />
            <el-option label="轻食" value="轻食" />
          </el-select>
        </el-col>
        <el-col :span="4" class="stats-info">
          <span>共 {{ filteredProducts.length }} 个商品</span>
        </el-col>
      </el-row>
    </el-card>

    <!-- 商品列表 -->
    <div class="products-grid" v-loading="loading">
      <el-card
        v-for="product in paginatedProducts"
        :key="product.goods_code"
        class="product-card"
        :class="{ 'out-of-stock': product.stock_status === 'out' }"
      >
        <!-- 商品头部信息 -->
        <div class="product-header">
          <div class="product-title">
            <h3>{{ product.goods_name }}</h3>
            <el-tag :type="getStockStatusType(product.stock_status)" size="small">
              {{ getStockStatusText(product.stock_status) }}
            </el-tag>
          </div>
          <div class="product-code">
            <span>编码: {{ product.goods_code }}</span>
            <el-button
              link
              size="small"
              @click="copyToClipboard(product.goods_code)"
            >
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </div>
        </div>

        <!-- 商品基本信息 -->
        <div class="product-info">
          <div class="info-item">
            <span class="label">类别:</span>
            <span class="value">{{ product.category }}</span>
          </div>
          <div class="info-item">
            <span class="label">价格:</span>
            <span class="value">{{ product.price_range }}</span>
          </div>
          <div class="info-item" v-if="product.aliases.length > 0">
            <span class="label">别名:</span>
            <span class="value">{{ product.aliases.join('、') }}</span>
          </div>
        </div>

        <!-- SKU信息 -->
        <div class="sku-info">
          <div class="sku-header">
            <span>SKU规格 ({{ product.skus.length }}种)</span>
            <el-button link size="small" @click="showProductDetail(product)">
              查看全部
            </el-button>
          </div>
          <div class="sku-preview">
            <el-tag
              v-for="(sku, index) in product.skus.slice(0, 3)"
              :key="sku.sku_code"
              size="small"
              effect="plain"
            >
              {{ sku.chinese_desc }}
            </el-tag>
            <el-tag v-if="product.skus.length > 3" size="small" type="info">
              +{{ product.skus.length - 3 }}
            </el-tag>
          </div>
        </div>

        <!-- 操作提示 -->
        <div class="product-tips" v-if="product.ordering_tips">
          <el-alert :title="product.ordering_tips" type="info" :closable="false" />
        </div>
      </el-card>
    </div>

    <!-- 分页 -->
    <div class="pagination-container" v-if="filteredProducts.length > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48, 96]"
        :total="filteredProducts.length"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 商品详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="`${currentProduct?.goods_name} - 详细信息`"
      width="800px"
    >
      <div v-if="currentProduct" class="product-detail">
        <!-- 基本信息 -->
        <el-descriptions :column="2" border>
          <el-descriptions-item label="商品编码">
            {{ currentProduct.goods_code }}
            <el-button
              link
              size="small"
              @click="copyToClipboard(currentProduct.goods_code)"
            >
              复制
            </el-button>
          </el-descriptions-item>
          <el-descriptions-item label="商品名称">
            {{ currentProduct.goods_name }}
          </el-descriptions-item>
          <el-descriptions-item label="类别">
            {{ currentProduct.category }}
          </el-descriptions-item>
          <el-descriptions-item label="价格区间">
            {{ currentProduct.price_range }}
          </el-descriptions-item>
          <el-descriptions-item label="库存状态">
            <el-tag :type="getStockStatusType(currentProduct.stock_status)">
              {{ getStockStatusText(currentProduct.stock_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="别名" :span="2">
            {{ currentProduct.aliases.length > 0 ? currentProduct.aliases.join('、') : '无' }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 可用规格选项 -->
        <div v-if="currentProduct.specs_options && Object.keys(currentProduct.specs_options).length > 0" style="margin: 20px 0;">
          <h4 style="margin-bottom: 10px;">可用规格选项</h4>
          <div v-for="(options, specType) in currentProduct.specs_options" :key="specType" class="spec-options">
            <span class="spec-label">{{ getSpecTypeLabel(specType) }}:</span>
            <el-space wrap>
              <el-tag
                v-for="option in options"
                :key="option.code"
                size="small"
                :type="getSpecTagType(specType)"
              >
                {{ option.name }} ({{ option.code }})
              </el-tag>
            </el-space>
          </div>
        </div>

        <!-- SKU列表 -->
        <h4 style="margin: 20px 0 10px;">SKU规格列表 ({{ currentProduct.skus.length }}种)</h4>
        <el-table :data="currentProduct.skus" style="width: 100%">
          <el-table-column prop="sku_code" label="SKU编码" width="120">
            <template #default="{ row }">
              {{ row.sku_code }}
              <el-button
                link
                size="small"
                @click="copyToClipboard(row.sku_code)"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </template>
          </el-table-column>
          <el-table-column prop="chinese_desc" label="中文描述" />
          <el-table-column label="规格详情">
            <template #default="{ row }">
              <el-space>
                <el-tag v-if="row.specs.size" size="small">
                  {{ row.specs.size.name }}
                </el-tag>
                <el-tag v-if="row.specs.temperature" size="small" type="warning">
                  {{ row.specs.temperature.name }}
                </el-tag>
                <el-tag v-if="row.specs.sweetness" size="small" type="success">
                  {{ row.specs.sweetness.name }}
                </el-tag>
              </el-space>
            </template>
          </el-table-column>
          <el-table-column prop="specs_code" label="规格编码" width="100" />
          <el-table-column label="默认" width="60">
            <template #default="{ row }">
              <el-tag v-if="row.is_default" type="primary" size="small">默认</el-tag>
            </template>
          </el-table-column>
        </el-table>

        <!-- 下单示例 -->
        <h4 style="margin: 20px 0 10px;">下单示例代码</h4>
        <el-tabs v-model="activeTab">
          <el-tab-pane label="简化下单" name="simplified">
            <el-input
              type="textarea"
              :value="getSimplifiedOrderExample(currentProduct)"
              :rows="10"
              readonly
            />
            <el-button
              type="primary"
              style="margin-top: 10px"
              @click="copyToClipboard(getSimplifiedOrderExample(currentProduct))"
            >
              复制代码
            </el-button>
          </el-tab-pane>
          <el-tab-pane label="标准下单" name="standard">
            <el-input
              type="textarea"
              :value="getStandardOrderExample(currentProduct)"
              :rows="10"
              readonly
            />
            <el-button
              type="primary"
              style="margin-top: 10px"
              @click="copyToClipboard(getStandardOrderExample(currentProduct))"
            >
              复制代码
            </el-button>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>

    <!-- 最后更新时间 -->
    <div class="update-time" v-if="updateTime">
      最后更新: {{ formatTime(updateTime) }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Refresh, 
  Download, 
  Search, 
  CopyDocument 
} from '@element-plus/icons-vue'
import { 
  getAvailableProducts, 
  exportAvailableProducts,
  getStockStatusText,
  getStockStatusType
} from '@/api/availableProducts'

// 数据状态
const loading = ref(false)
const products = ref([])
const updateTime = ref(null)
const searchText = ref('')
const filterStatus = ref('')
const filterCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(24)

// 详情弹窗
const detailDialogVisible = ref(false)
const currentProduct = ref(null)
const activeTab = ref('simplified')

// 自动刷新定时器
let refreshTimer = null

// 计算属性：过滤后的商品
const filteredProducts = computed(() => {
  let result = products.value

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(product => 
      product.goods_name.toLowerCase().includes(search) ||
      product.goods_code.toLowerCase().includes(search) ||
      product.aliases.some(alias => alias.toLowerCase().includes(search))
    )
  }

  // 状态过滤
  if (filterStatus.value) {
    result = result.filter(product => product.stock_status === filterStatus.value)
  }

  // 类别过滤
  if (filterCategory.value) {
    result = result.filter(product => product.category === filterCategory.value)
  }

  return result
})

// 计算属性：分页后的商品
const paginatedProducts = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredProducts.value.slice(start, end)
})

// 获取商品数据
const fetchProducts = async () => {
  loading.value = true
  try {
    const res = await getAvailableProducts()
    if (res.code === 200 && res.data) {
      products.value = res.data.products || []
      updateTime.value = res.data.updated_at
      ElMessage.success('数据加载成功')
    }
  } catch (error) {
    ElMessage.error('获取商品数据失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 刷新数据
const handleRefresh = () => {
  fetchProducts()
}

// 导出数据
const handleExport = () => {
  if (products.value.length === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }
  exportAvailableProducts({ products: filteredProducts.value })
  ElMessage.success('导出成功')
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

// 筛选处理
const handleFilter = () => {
  currentPage.value = 1
}

// 分页处理
const handleSizeChange = () => {
  currentPage.value = 1
}

const handleCurrentChange = () => {
  // 页码改变时自动处理
}

// 显示商品详情
const showProductDetail = (product) => {
  currentProduct.value = product
  detailDialogVisible.value = true
  activeTab.value = 'simplified'
}

// 复制到剪贴板
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

// 格式化时间
const formatTime = (time) => {
  return new Date(time).toLocaleString('zh-CN')
}

// 获取简化下单示例
const getSimplifiedOrderExample = (product) => {
  const sku = product.skus[0] || {}
  return `{
  "store_name": "望京SOHO店",
  "phone": "13800138000",
  "items": [
    {
      "product_name": "${product.goods_name}",
      "specs": "${sku.chinese_desc || '大杯|冰|标准糖'}",
      "quantity": 1
    }
  ],
  "callback_url": "https://your-domain.com/callback"
}`
}

// 获取标准下单示例
const getStandardOrderExample = (product) => {
  const sku = product.skus[0] || {}
  return `{
  "store_code": "10001",
  "phone": "13800138000",
  "goods": [
    {
      "goods_code": "${product.goods_code}",
      "goods_name": "${product.goods_name}",
      "sku_code": "${sku.sku_code || 'SKU001'}",
      "specs_code": "${sku.specs_code || '1_0_0'}",
      "goods_num": 1,
      "sale_price": 15.0
    }
  ],
  "callback_url": "https://your-domain.com/callback"
}`
}

// 获取规格类型标签
const getSpecTypeLabel = (specType) => {
  const labels = {
    'size': '杯型',
    'temperature': '温度',
    'sweetness': '甜度',
    'milk': '奶',
    'flavor': '口味',
    'bean': '咖啡豆',
    'other': '其他'
  }
  return labels[specType] || specType
}

// 获取规格标签类型
const getSpecTagType = (specType) => {
  const types = {
    'size': '',
    'temperature': 'warning',
    'sweetness': 'success',
    'milk': 'info',
    'flavor': 'danger'
  }
  return types[specType] || ''
}

// 设置自动刷新
const setupAutoRefresh = () => {
  // 每5分钟自动刷新
  refreshTimer = setInterval(() => {
    fetchProducts()
  }, 5 * 60 * 1000)
}

// 组件挂载
onMounted(() => {
  fetchProducts()
  setupAutoRefresh()
})

// 组件卸载
onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped lang="scss">
.available-products-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  .header-left {
    display: flex;
    align-items: center;
    gap: 10px;

    h2 {
      margin: 0;
    }
  }

  .header-right {
    display: flex;
    gap: 10px;
  }
}

.search-card {
  margin-bottom: 20px;

  .stats-info {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    color: #666;
    font-size: 14px;
  }
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.product-card {
  transition: all 0.3s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  &.out-of-stock {
    opacity: 0.8;
    
    .product-header h3 {
      color: #999;
    }
  }

  .product-header {
    margin-bottom: 12px;

    .product-title {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;

      h3 {
        margin: 0;
        font-size: 16px;
        font-weight: 500;
      }
    }

    .product-code {
      display: flex;
      align-items: center;
      gap: 5px;
      color: #666;
      font-size: 13px;
    }
  }

  .product-info {
    margin-bottom: 12px;

    .info-item {
      display: flex;
      margin-bottom: 6px;
      font-size: 14px;

      .label {
        color: #666;
        margin-right: 8px;
        min-width: 45px;
      }

      .value {
        color: #333;
        flex: 1;
      }
    }
  }

  .sku-info {
    border-top: 1px solid #f0f0f0;
    padding-top: 12px;

    .sku-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;
      font-size: 14px;
      color: #666;
    }

    .sku-preview {
      display: flex;
      flex-wrap: wrap;
      gap: 5px;
    }
  }

  .product-tips {
    margin-top: 12px;

    :deep(.el-alert) {
      padding: 8px 12px;
      
      .el-alert__title {
        font-size: 13px;
      }
    }
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.product-detail {
  h4 {
    color: #333;
    font-weight: 500;
  }

  :deep(.el-descriptions) {
    .el-descriptions__label {
      font-weight: 500;
    }
  }

  :deep(.el-tabs) {
    margin-top: 10px;
  }

  :deep(.el-textarea__inner) {
    font-family: 'Consolas', 'Monaco', monospace;
    font-size: 13px;
    line-height: 1.5;
  }
  
  .spec-options {
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    
    .spec-label {
      font-weight: 500;
      color: #666;
      margin-right: 10px;
      min-width: 60px;
    }
  }
}

.update-time {
  text-align: center;
  color: #999;
  font-size: 13px;
  margin-top: 20px;
}
</style>