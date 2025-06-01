<template>
  <div class="stores-page">
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="城市">
          <el-select 
            v-model="searchForm.city_id" 
            placeholder="请选择城市" 
            clearable
            filterable
            style="width: 200px"
          >
            <el-option
              v-for="city in cities"
              :key="city.city_id"
              :label="city.city_name"
              :value="city.city_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="店铺名称或地址"
            clearable
            style="width: 300px"
          />
        </el-form-item>
        <el-form-item label="经度">
          <el-input
            v-model="searchForm.longitude"
            placeholder="如：116.397428"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="纬度">
          <el-input
            v-model="searchForm.latitude"
            placeholder="如：39.90923"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch" :loading="loading">
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button type="success" @click="getLocation">
            获取当前位置
          </el-button>
          <el-button type="warning" @click="handleSyncCities" :loading="syncLoading">
            <el-icon><Refresh /></el-icon>
            同步城市列表
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="tips">
        <el-alert
          title="使用提示"
          type="info"
          :closable="false"
          show-icon
        >
          <template #default>
            <ul style="margin: 0; padding-left: 20px;">
              <li>选择城市后会自动加载该城市的所有门店</li>
              <li>输入关键词可以搜索店铺名称或地址</li>
              <li>输入经纬度可以搜索附近的门店</li>
              <li>点击"获取当前位置"可以自动填充您的经纬度</li>
              <li>店铺代码（store_code）会在搜索结果中显示，可以复制使用</li>
            </ul>
          </template>
        </el-alert>
      </div>
    </el-card>

    <!-- 搜索结果 -->
    <el-card class="result-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>搜索结果</span>
          <span v-if="stores.length > 0" class="result-count">
            共找到 {{ total }} 家门店
          </span>
        </div>
      </template>
      
      <div v-if="stores.length === 0 && !loading" class="empty-result">
        <el-empty description="暂无搜索结果" />
      </div>
      
      <div v-else class="stores-grid">
        <div 
          v-for="store in stores" 
          :key="store.store_code"
          class="store-card"
        >
          <div class="store-header">
            <h3>{{ store.store_name }}</h3>
            <el-tag type="success" v-if="store.is_open">营业中</el-tag>
            <el-tag type="info" v-else>已关闭</el-tag>
          </div>
          
          <div class="store-info">
            <div class="info-item">
              <span class="label">店铺代码：</span>
              <span class="value">
                {{ store.store_code }}
                <el-button
                  link
                  type="primary"
                  :icon="CopyDocument"
                  @click="copyStoreCode(store.store_code)"
                >
                  复制
                </el-button>
              </span>
            </div>
            <div class="info-item">
              <span class="label">地址：</span>
              <span class="value">{{ store.address }}</span>
            </div>
            <div class="info-item">
              <span class="label">电话：</span>
              <span class="value">{{ store.phone || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">营业时间：</span>
              <span class="value">{{ store.business_hours || '-' }}</span>
            </div>
            <div class="info-item" v-if="store.distance">
              <span class="label">距离：</span>
              <span class="value">{{ formatDistance(store.distance) }}</span>
            </div>
          </div>
          
          <div class="store-features" v-if="store.features && store.features.length > 0">
            <el-tag 
              v-for="feature in store.features" 
              :key="feature"
              size="small"
              style="margin-right: 5px;"
            >
              {{ feature }}
            </el-tag>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <el-pagination
        v-if="total > pageSize"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleSearch"
        @size-change="handleSearch"
      />
    </el-card>

    <!-- 热门店铺 -->
    <el-card class="hot-stores-card" v-if="hotStores.length > 0">
      <template #header>
        <span>热门店铺</span>
      </template>
      <div class="hot-stores-list">
        <el-tag
          v-for="store in hotStores"
          :key="store.store_code"
          @click="quickSearchStore(store)"
          style="margin: 5px; cursor: pointer;"
          effect="plain"
        >
          {{ store.store_name }} ({{ store.store_code }})
        </el-tag>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CopyDocument, Refresh } from '@element-plus/icons-vue'
import { searchStores, getCities, getHotStores, syncCities } from '@/api/store'

// 数据
const loading = ref(false)
const syncLoading = ref(false)
const cities = ref([])
const stores = ref([])
const hotStores = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索表单
const searchForm = ref({
  city_id: '',
  keyword: '',
  longitude: '',
  latitude: ''
})

// 加载城市列表
const loadCities = async () => {
  try {
    const res = await getCities()
    cities.value = res.data || []
  } catch (error) {
    console.error('加载城市列表失败:', error)
  }
}

// 加载热门店铺
const loadHotStores = async () => {
  try {
    const res = await getHotStores()
    hotStores.value = res.data || []
  } catch (error) {
    console.error('加载热门店铺失败:', error)
  }
}

// 搜索店铺
const handleSearch = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      ...searchForm.value
    }
    
    // 移除空值参数
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null) {
        delete params[key]
      }
    })
    
    const res = await searchStores(params)
    stores.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    ElMessage.error('搜索失败')
  } finally {
    loading.value = false
  }
}

// 重置搜索
const handleReset = () => {
  searchForm.value = {
    city_id: '',
    keyword: '',
    longitude: '',
    latitude: ''
  }
  currentPage.value = 1
  stores.value = []
  total.value = 0
}

// 获取当前位置
const getLocation = () => {
  if (!navigator.geolocation) {
    ElMessage.error('您的浏览器不支持地理位置功能')
    return
  }
  
  const loading = ElMessage({
    message: '正在获取位置...',
    duration: 0
  })
  
  navigator.geolocation.getCurrentPosition(
    (position) => {
      searchForm.value.longitude = position.coords.longitude.toFixed(6)
      searchForm.value.latitude = position.coords.latitude.toFixed(6)
      loading.close()
      ElMessage.success('位置获取成功')
    },
    (error) => {
      loading.close()
      let message = '获取位置失败'
      switch(error.code) {
        case error.PERMISSION_DENIED:
          message = '您拒绝了位置权限请求'
          break
        case error.POSITION_UNAVAILABLE:
          message = '位置信息不可用'
          break
        case error.TIMEOUT:
          message = '获取位置超时'
          break
      }
      ElMessage.error(message)
    }
  )
}

// 复制店铺代码
const copyStoreCode = async (code) => {
  try {
    await navigator.clipboard.writeText(code)
    ElMessage.success('店铺代码已复制到剪贴板')
  } catch (error) {
    // 降级方案
    const input = document.createElement('input')
    input.value = code
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    ElMessage.success('店铺代码已复制到剪贴板')
  }
}

// 快速搜索店铺
const quickSearchStore = (store) => {
  searchForm.value.keyword = store.store_name
  handleSearch()
}

// 格式化距离
const formatDistance = (distance) => {
  if (distance < 1000) {
    return `${distance}米`
  } else {
    return `${(distance / 1000).toFixed(1)}公里`
  }
}

// 同步城市列表
const handleSyncCities = async () => {
  try {
    syncLoading.value = true
    
    // 先尝试不提供卡片代码，让系统自动使用内部卡片
    try {
      const res = await syncCities()
      if (res.code === 200) {
        ElMessage.success('城市列表同步成功')
        cities.value = res.data || []
        return
      }
    } catch (autoError) {
      // 如果自动同步失败，提示用户输入卡片代码
      if (autoError.response?.data?.code === 400) {
        const { value: cardCode } = await ElMessageBox.prompt(
          autoError.response.data.msg || '请输入有效的卡片代码来同步城市列表',
          '同步城市',
          {
            confirmButtonText: '同步',
            cancelButtonText: '取消',
            inputPattern: /^[A-Za-z0-9]+$/,
            inputErrorMessage: '请输入有效的卡片代码'
          }
        )
        
        const res = await syncCities(cardCode)
        if (res.code === 200) {
          ElMessage.success('城市列表同步成功')
          cities.value = res.data || []
        }
      } else {
        throw autoError
      }
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步城市失败:', error)
    }
  } finally {
    syncLoading.value = false
  }
}

onMounted(() => {
  loadCities()
  loadHotStores()
})
</script>

<style lang="scss" scoped>
.stores-page {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
  
  .tips {
    margin-top: 20px;
  }
}

.result-card {
  margin-bottom: 20px;
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .result-count {
      font-size: 14px;
      color: #909399;
    }
  }
}

.empty-result {
  padding: 40px 0;
}

.stores-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.store-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px;
  transition: all 0.3s;
  
  &:hover {
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }
  
  .store-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    
    h3 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
    }
  }
  
  .store-info {
    .info-item {
      margin-bottom: 8px;
      font-size: 14px;
      
      .label {
        color: #909399;
        margin-right: 8px;
      }
      
      .value {
        color: #303133;
      }
    }
  }
  
  .store-features {
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid #f0f0f0;
  }
}

.hot-stores-card {
  .hot-stores-list {
    line-height: 2;
  }
}
</style>