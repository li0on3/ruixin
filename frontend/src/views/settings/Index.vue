<template>
  <div class="page-container">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="基础设置" name="basic">
        <el-form 
          ref="formRef"
          :model="settingsForm" 
          :rules="rules"
          label-width="140px" 
          style="max-width: 600px"
          v-loading="loading"
        >
          <el-form-item label="系统名称">
            <el-input value="瑞幸咖啡分销商自动化系统" disabled />
          </el-form-item>
          <el-form-item label="系统版本">
            <el-input value="1.0.0" disabled />
          </el-form-item>
          
          <el-divider>商品同步配置</el-divider>
          
          <el-form-item label="店铺代码" prop="sync_store_code">
            <el-input 
              v-model="settingsForm.sync_store_code" 
              placeholder="请输入店铺代码，如：390840"
              maxlength="20"
            >
              <template #append>
                <el-tooltip content="用于卡片验证和商品同步的店铺代码">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </template>
            </el-input>
            <div class="el-form-item__description">
              此店铺代码用于：
              <ul style="margin: 5px 0; padding-left: 20px;">
                <li>创建卡片时验证卡片有效性</li>
                <li>同步商品时获取店铺的商品信息</li>
              </ul>
            </div>
          </el-form-item>
          
          <el-form-item label="自动同步商品" prop="sync_enabled">
            <el-switch 
              v-model="settingsForm.sync_enabled"
              active-value="true"
              inactive-value="false"
            />
            <span style="margin-left: 10px; color: #909399;">
              {{ settingsForm.sync_enabled === 'true' ? '开启后，添加卡片时会自动同步该卡片可购买的商品' : '关闭状态，需要手动同步商品' }}
            </span>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="saveSettings" :loading="saving">
              保存设置
            </el-button>
            <el-button @click="loadSettings">重置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      
      <el-tab-pane label="安全设置" name="security">
        <el-alert
          title="安全设置功能开发中"
          type="info"
          show-icon
          :closable="false"
        />
      </el-tab-pane>
      
      <el-tab-pane label="系统日志" name="logs">
        <el-alert
          title="系统日志功能开发中"
          type="info"
          show-icon
          :closable="false"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { QuestionFilled } from '@element-plus/icons-vue'
import { getSystemConfigs, updateSystemConfigs } from '@/api/system'

const activeTab = ref('basic')
const loading = ref(false)
const saving = ref(false)
const formRef = ref()

const settingsForm = ref({
  sync_store_code: '',
  sync_enabled: 'false'
})

const rules = {
  sync_store_code: [
    { required: true, message: '请输入店铺代码', trigger: 'blur' },
    { pattern: /^\d+$/, message: '店铺代码必须为数字', trigger: 'blur' }
  ]
}

// 加载系统配置
const loadSettings = async () => {
  loading.value = true
  try {
    const res = await getSystemConfigs()
    if (res.code === 200) {
      settingsForm.value = {
        sync_store_code: res.data.sync_store_code || '390840',
        sync_enabled: res.data.sync_enabled || 'false'
      }
    }
  } catch (error) {
    ElMessage.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

// 保存设置
const saveSettings = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return

  saving.value = true
  try {
    const res = await updateSystemConfigs(settingsForm.value)
    if (res.code === 200) {
      ElMessage.success('配置更新成功')
    } else {
      ElMessage.error(res.msg || '更新失败')
    }
  } catch (error) {
    ElMessage.error('保存配置失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.el-form-item__description {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 5px;
}

.el-form-item__description ul {
  list-style-type: disc;
}

.el-form-item__description li {
  margin-bottom: 3px;
}
</style>