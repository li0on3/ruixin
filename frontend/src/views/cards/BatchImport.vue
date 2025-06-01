<template>
  <div class="page-container view-container">
    <el-card class="batch-import-card">
      <template #header>
        <div class="card-header">
          <span>批量导入卡片</span>
          <el-button link @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        class="import-form"
      >
        <el-form-item label="成本价" prop="cost_price">
          <el-input-number
            v-model="form.cost_price"
            :min="0"
            :precision="2"
            :step="0.1"
            placeholder="每张卡片的成本价"
          />
          <span class="ml-10">元/张</span>
        </el-form-item>

        <el-form-item label="销售价" prop="sell_price">
          <el-input-number
            v-model="form.sell_price"
            :min="0"
            :precision="2"
            :step="0.1"
            placeholder="每张卡片的销售价"
          />
          <span class="ml-10">元/张</span>
          <div class="form-tip">
            <el-alert type="warning" :closable="false" show-icon>
              <template #title>
                重要：此价格将作为实际扣款金额
              </template>
            </el-alert>
          </div>
        </el-form-item>

        <el-form-item label="瑞幸产品ID" prop="luckin_product_id">
          <el-input-number v-model="form.luckin_product_id" :min="1" :max="100" style="width: 200px;" />
          <div class="form-tip">瑞幸API所需的产品ID，默认为6，一般不需要修改</div>
        </el-form-item>

        <el-form-item label="过期时间" prop="expired_at">
          <el-date-picker
            v-model="form.expired_at"
            type="datetime"
            placeholder="选择过期时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
          />
        </el-form-item>

        <el-form-item label="导入方式">
          <el-radio-group v-model="importMode">
            <el-radio value="codes">卡片代码</el-radio>
            <el-radio value="urls">卡片链接</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="importMode === 'codes'" label="卡片代码" prop="card_codes">
          <el-input
            v-model="form.card_codes"
            type="textarea"
            :rows="10"
            placeholder="请输入卡片代码，每行一个"
          />
          <div class="form-tip">
            请输入卡片代码，每行一个。支持从Excel复制粘贴。
            <br>当前已输入：{{ cardCodeCount }} 个卡片代码
          </div>
        </el-form-item>

        <el-form-item v-else label="卡片链接" prop="card_urls">
          <el-input
            v-model="form.card_urls"
            type="textarea"
            :rows="10"
            placeholder="请输入卡片链接，每行一个"
          />
          <div class="form-tip">
            请输入卡片链接，每行一个。支持格式：https://lkcoffe.cn/?card=XXXXX
            <br>当前已输入：{{ cardUrlCount }} 个卡片链接
          </div>
        </el-form-item>

        <el-form-item label="备注" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="批次备注信息"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            导入卡片
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>

      <!-- 预览信息 -->
      <div v-if="(importMode === 'codes' && cardCodeCount > 0) || (importMode === 'urls' && cardUrlCount > 0)" class="preview-info">
        <h4>导入预览</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="卡片数量">{{ importMode === 'codes' ? cardCodeCount : cardUrlCount }}</el-descriptions-item>
          <el-descriptions-item label="总成本">¥{{ (form.cost_price * (importMode === 'codes' ? cardCodeCount : cardUrlCount)).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="总售价">¥{{ (form.sell_price * (importMode === 'codes' ? cardCodeCount : cardUrlCount)).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="预计利润">¥{{ ((form.sell_price - form.cost_price) * (importMode === 'codes' ? cardCodeCount : cardUrlCount)).toFixed(2) }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { batchImportCards } from '@/api/card'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const importMode = ref('codes')

const form = reactive({
  luckin_product_id: 6,  // 默认值6
  cost_price: 0,
  sell_price: 0,
  card_codes: '',
  card_urls: '',
  expired_at: '',
  description: ''
})

const rules = computed(() => ({
  luckin_product_id: [
    { required: true, message: '请输入瑞幸产品ID', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '产品ID范围为1-100', trigger: 'blur' }
  ],
  cost_price: [
    { required: true, message: '请输入成本价', trigger: 'blur' },
    { type: 'number', min: 0, message: '成本价不能小于0', trigger: 'blur' }
  ],
  sell_price: [
    { required: true, message: '请输入销售价', trigger: 'blur' },
    { type: 'number', min: 0, message: '销售价不能小于0', trigger: 'blur' }
  ],
  card_codes: importMode.value === 'codes' ? [
    { required: true, message: '请输入卡片代码', trigger: 'blur' }
  ] : [],
  card_urls: importMode.value === 'urls' ? [
    { required: true, message: '请输入卡片链接', trigger: 'blur' }
  ] : [],
  expired_at: [
    { required: true, message: '请选择过期时间', trigger: 'change' }
  ]
}))

// 计算卡片代码数量
const cardCodeCount = computed(() => {
  if (!form.card_codes.trim()) return 0
  const codes = form.card_codes.split('\n').filter(code => code.trim())
  return codes.length
})

// 计算卡片链接数量
const cardUrlCount = computed(() => {
  if (!form.card_urls.trim()) return 0
  const urls = form.card_urls.split('\n').filter(url => url.trim())
  return urls.length
})


// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return

  let requestData = {
    luckin_product_id: form.luckin_product_id,
    cost_price: form.cost_price,
    sell_price: form.sell_price,
    expired_at: form.expired_at,
    description: form.description
  }

  if (importMode.value === 'codes') {
    // 处理卡片代码
    const codes = form.card_codes.split('\n')
      .map(code => code.trim())
      .filter(code => code)

    if (codes.length === 0) {
      ElMessage.error('请输入至少一个卡片代码')
      return
    }

    // 检查是否有重复
    const uniqueCodes = [...new Set(codes)]
    if (uniqueCodes.length < codes.length) {
      ElMessage.warning(`发现重复的卡片代码，已自动去重。实际导入数量：${uniqueCodes.length}`)
    }

    requestData.card_codes = uniqueCodes
  } else {
    // 处理卡片链接
    const urls = form.card_urls.split('\n')
      .map(url => url.trim())
      .filter(url => url)

    if (urls.length === 0) {
      ElMessage.error('请输入至少一个卡片链接')
      return
    }

    // 检查是否有重复
    const uniqueUrls = [...new Set(urls)]
    if (uniqueUrls.length < urls.length) {
      ElMessage.warning(`发现重复的卡片链接，已自动去重。实际导入数量：${uniqueUrls.length}`)
    }

    requestData.card_urls = uniqueUrls
  }

  loading.value = true
  try {
    const { data } = await batchImportCards(requestData)

    ElMessage.success(`批量导入成功！批次号：${data.batch_no}`)
    router.push('/cards')
  } catch (error) {
    console.error('Batch import failed:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // 设置默认过期时间为一年后
  const defaultExpiredAt = new Date()
  defaultExpiredAt.setFullYear(defaultExpiredAt.getFullYear() + 1)
  form.expired_at = defaultExpiredAt.toISOString().slice(0, -1) + 'Z'
})
</script>

<style lang="scss" scoped>
.page-container {
  width: 100%;
  height: 100%;
  overflow: auto;
  padding: 16px;
  box-sizing: border-box;
  
  &.view-container {
    display: flex;
    flex-direction: column;
  }
  
  @media (min-width: 1366px) {
    padding: 20px;
  }
  
  @media (min-width: 1920px) {
    padding: 24px 32px;
  }
}

.batch-import-card {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  
  .import-form {
    max-width: 800px;
    margin: 0 auto;
    width: 100%;
    
    @media (min-width: 1440px) {
      max-width: 900px;
    }
    
    @media (min-width: 1920px) {
      max-width: 1000px;
    }
    
    .el-form-item {
      margin-bottom: 24px;
      
      @media (min-width: 1920px) {
        margin-bottom: 28px;
      }
    }
    
    .el-input,
    .el-select,
    .el-date-picker {
      width: 100%;
    }
    
    .el-input-number {
      width: 200px;
      
      @media (min-width: 1920px) {
        width: 240px;
      }
    }
    
    .el-textarea {
      .el-textarea__inner {
        font-size: 16px;
        line-height: 1.6;
        font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
        
        @media (min-width: 1920px) {
          font-size: 18px;
        }
      }
    }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 20px;
  font-weight: 700;
  
  @media (min-width: 1920px) {
    font-size: 24px;
  }
}

.ml-10 {
  margin-left: 10px;
  color: #606266;
  font-size: 14px;
  
  @media (min-width: 1920px) {
    margin-left: 12px;
    font-size: 16px;
  }
}

.form-tip {
  margin-top: 8px;
  color: #909399;
  font-size: 14px;
  line-height: 1.6;
  
  @media (min-width: 1920px) {
    margin-top: 10px;
    font-size: 16px;
  }
}

.preview-info {
  margin-top: 40px;
  padding: 24px;
  background-color: #f5f7fa;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  
  @media (min-width: 1920px) {
    margin-top: 48px;
    padding: 32px;
    border-radius: 12px;
  }

  h4 {
    margin: 0 0 20px 0;
    color: #303133;
    font-size: 18px;
    font-weight: 700;
    
    @media (min-width: 1920px) {
      margin-bottom: 24px;
      font-size: 20px;
    }
  }
  
  .el-descriptions {
    .el-descriptions__label {
      font-weight: 600;
      color: #606266;
      
      @media (min-width: 1920px) {
        font-size: 16px;
      }
    }
    
    .el-descriptions__content {
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      
      @media (min-width: 1920px) {
        font-size: 18px;
      }
    }
  }
}

// 优化表单按钮组
.el-form-item:last-child {
  margin-top: 32px;
  
  @media (min-width: 1920px) {
    margin-top: 40px;
  }
  
  .el-button {
    min-width: 120px;
    
    @media (min-width: 1920px) {
      min-width: 140px;
    }
  }
}

// 优化单选按钮组
.el-radio-group {
  .el-radio {
    margin-right: 30px;
    
    @media (min-width: 1920px) {
      margin-right: 40px;
    }
    
    .el-radio__label {
      font-size: 16px;
      
      @media (min-width: 1920px) {
        font-size: 18px;
      }
    }
  }
}
</style>