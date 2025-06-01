<template>
  <div class="cards-enhanced-container">
    <!-- 页面加载骨架屏 -->
    <LoadingAnimation
      v-if="pageLoading"
      type="skeleton"
      text="正在加载卡片数据..."
      :visible="pageLoading"
    />
    
    <!-- 主要内容 -->
    <div v-else class="cards-content">
      <!-- 搜索区域 -->
      <el-card class="search-card" shadow="hover">
        <template #header>
          <div class="search-header">
            <h3 class="search-title">
              <el-icon><Search /></el-icon>
              筛选查询
            </h3>
            <el-button text @click="handleToggleSearch" class="toggle-btn">
              {{ showAdvancedSearch ? '简单搜索' : '高级搜索' }}
            </el-button>
          </div>
        </template>
        
        <el-form 
          :model="searchForm" 
          :inline="!showAdvancedSearch"
          class="search-form"
          :class="{ 'advanced-search': showAdvancedSearch }"
        >
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" :lg="6">
              <el-form-item label="卡片代码">
                <el-input
                  v-model="searchForm.cardCode"
                  placeholder="请输入卡片代码"
                  clearable
                  @keyup.enter="handleSearch"
                >
                  <template #prefix>
                    <el-icon><CreditCard /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
            
            <el-col :xs="24" :sm="12" :md="8" :lg="6">
              <el-form-item label="卡片状态">
                <el-select 
                  v-model="searchForm.status" 
                  placeholder="请选择状态" 
                  clearable 
                  filterable
                  style="width: 100%"
                >
                  <el-option label="全部" value="" />
                  <el-option label="未使用" :value="0">
                    <span class="option-with-status">
                      <el-tag type="success" size="small">未使用</el-tag>
                    </span>
                  </el-option>
                  <el-option label="已使用" :value="1">
                    <span class="option-with-status">
                      <el-tag type="info" size="small">已使用</el-tag>
                    </span>
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :xs="24" :sm="12" :md="8" :lg="6">
              <el-form-item label="批次ID">
                <el-input
                  v-model="searchForm.batchId"
                  placeholder="请输入批次ID"
                  clearable
                  @keyup.enter="handleSearch"
                >
                  <template #prefix>
                    <el-icon><Files /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
            
            <el-col v-if="showAdvancedSearch" :xs="24" :sm="12" :md="8" :lg="6">
              <el-form-item label="创建时间">
                <el-date-picker
                  v-model="searchForm.dateRange"
                  type="daterange"
                  range-separator="至"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  format="YYYY-MM-DD"
                  value-format="YYYY-MM-DD"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            
            <el-col :xs="24" :sm="12" :md="8" :lg="6">
              <el-form-item class="search-actions">
                <el-button type="primary" @click="handleSearch" :loading="searching">
                  <el-icon><Search /></el-icon>
                  搜索
                </el-button>
                <el-button @click="handleReset">
                  <el-icon><RefreshLeft /></el-icon>
                  重置
                </el-button>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-card>
      
      <!-- 操作工具栏 -->
      <el-card class="toolbar-card" shadow="hover">
        <div class="toolbar-content">
          <div class="toolbar-left">
            <h3 class="toolbar-title">
              <el-icon><CreditCard /></el-icon>
              卡片管理
            </h3>
            <div class="stats-info">
              <el-tag type="info" size="small">
                总计: {{ pagination.total }}
              </el-tag>
            </div>
          </div>
          
          <div class="toolbar-right">
            <el-button-group>
              <el-button type="primary" @click="handleAdd">
                <el-icon><Plus /></el-icon>
                添加卡片
              </el-button>
              <el-button type="success" @click="handleBatchImport">
                <el-icon><Upload /></el-icon>
                批量导入
              </el-button>
              <el-button @click="handleViewBatches">
                <el-icon><Files /></el-icon>
                批次管理
              </el-button>
              <el-button @click="handleValidateCard">
                <el-icon><CircleCheck /></el-icon>
                验证卡片
              </el-button>
              <el-dropdown @command="handleValidationMode" trigger="click" placement="bottom-end">
                <el-button type="warning">
                  <el-icon><DocumentChecked /></el-icon>
                  批量验证
                  <el-icon><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="smart">
                      <div class="validation-mode-item">
                        <el-icon color="#409eff"><Lightning /></el-icon>
                        <div>
                          <div class="mode-title">智能验证</div>
                          <div class="mode-desc">只验证重要卡片，快速完成</div>
                        </div>
                      </div>
                    </el-dropdown-item>
                    <el-dropdown-item command="all">
                      <div class="validation-mode-item">
                        <el-icon color="#67c23a"><DocumentChecked /></el-icon>
                        <div>
                          <div class="mode-title">全量验证</div>
                          <div class="mode-desc">验证所有卡片，后台安全执行</div>
                        </div>
                      </div>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-button-group>
          </div>
        </div>
      </el-card>
      
      <!-- 数据表格 -->
      <el-card class="table-card" shadow="hover">
        <!-- 表格加载状态 -->
        <div v-if="loading" class="table-loading">
          <LoadingAnimation
            type="skeleton"
            :visible="loading"
          />
        </div>
        
        <!-- 实际表格 -->
        <div v-else>
          <!-- 桌面端表格 -->
          <div class="desktop-only">
            <el-table
              :data="tableData"
              v-loading="loading"
              stripe
              class="enhanced-table"
            >
              <el-table-column prop="cardCode" label="卡片代码" min-width="200" fixed="left">
                <template #default="{ row }">
                  <div class="card-code">
                    <el-icon class="card-icon"><CreditCard /></el-icon>
                    <span class="code-text">{{ row.card_code || row.cardCode }}</span>
                    <el-button
                      type="info"
                      link
                      size="small"
                      @click="copyToClipboard(row.card_code || row.cardCode)"
                      title="复制"
                    >
                      <el-icon><DocumentCopy /></el-icon>
                    </el-button>
                  </div>
                </template>
              </el-table-column>
              
              <el-table-column prop="batchId" label="批次ID" width="120">
                <template #default="{ row }">
                  <el-tag size="small" type="info">{{ row.batch_id || row.batchId || '-' }}</el-tag>
                </template>
              </el-table-column>
              
              <el-table-column prop="costPrice" label="成本价" width="120">
                <template #default="{ row }">
                  <span class="price-text">¥{{ (row.cost_price || row.costPrice || 0).toFixed(2) }}</span>
                </template>
              </el-table-column>
              
              <el-table-column prop="sellPrice" label="销售价" width="120">
                <template #default="{ row }">
                  <span class="price-text">¥{{ (row.sell_price || row.sellPrice || 0).toFixed(2) }}</span>
                </template>
              </el-table-column>
              
              <el-table-column prop="status" label="状态" width="120">
                <template #default="{ row }">
                  <el-tag
                    :type="getStatusType(row.status)"
                    size="small"
                  >
                    {{ getStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              
              <el-table-column prop="syncStatus" label="同步状态" width="120">
                <template #default="{ row }">
                  <el-tag
                    :type="getSyncStatusType(row.sync_status || row.syncStatus)"
                    size="small"
                  >
                    {{ getSyncStatusText(row.sync_status || row.syncStatus) }}
                  </el-tag>
                </template>
              </el-table-column>
              
              <el-table-column prop="expiredAt" label="过期时间" width="200">
                <template #default="{ row }">
                  <div class="time-cell" v-if="row.expired_at || row.expiredAt">
                    <el-icon class="time-icon"><Clock /></el-icon>
                    <span :class="{ 'expired-text': isExpired(row.expired_at || row.expiredAt) }">
                      {{ formatDate(row.expired_at || row.expiredAt) }}
                    </span>
                  </div>
                  <span v-else class="no-expire-text">永久有效</span>
                </template>
              </el-table-column>
              
              <el-table-column prop="createdAt" label="创建时间" width="200">
                <template #default="{ row }">
                  <div class="time-cell">
                    <el-icon class="time-icon"><Clock /></el-icon>
                    <span>{{ formatDate(row.created_at || row.createdAt) }}</span>
                  </div>
                </template>
              </el-table-column>
              
              <el-table-column prop="usedAt" label="使用时间" width="200">
                <template #default="{ row }">
                  <div class="time-cell" v-if="row.used_at || row.usedAt">
                    <el-icon class="time-icon"><Clock /></el-icon>
                    <span>{{ formatDate(row.used_at || row.usedAt) }}</span>
                  </div>
                  <span v-else class="unused-text">未使用</span>
                </template>
              </el-table-column>
              
              <el-table-column label="操作" width="240" fixed="right">
                <template #default="{ row }">
                  <div class="table-actions">
                    <el-button
                      type="primary"
                      link
                      size="small"
                      @click="handleView(row)"
                    >
                      <el-icon><View /></el-icon>
                      查看
                    </el-button>
                    <el-button
                      type="primary"
                      link
                      size="small"
                      @click="handleEdit(row)"
                    >
                      <el-icon><Edit /></el-icon>
                      编辑
                    </el-button>
                    <el-popconfirm
                      title="确定要删除这张卡片吗？"
                      @confirm="handleDelete(row)"
                    >
                      <template #reference>
                        <el-button
                          type="danger"
                          link
                          size="small"
                        >
                          <el-icon><Delete /></el-icon>
                          删除
                        </el-button>
                      </template>
                    </el-popconfirm>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <!-- 移动端卡片列表 -->
          <div class="mobile-only mobile-table-cards">
            <div class="mobile-card-list">
              <div
                v-for="item in tableData"
                :key="item.id"
                class="mobile-card-item"
              >
                <div class="card-item-header">
                  <h4 class="card-item-title">
                    <el-icon><CreditCard /></el-icon>
                    {{ item.card_code || item.cardCode }}
                  </h4>
                  <el-tag
                    :type="getStatusType(item.status)"
                    size="small"
                  >
                    {{ getStatusText(item.status) }}
                  </el-tag>
                </div>
                <div class="card-item-content">
                  <div class="card-item-row">
                    <span class="card-item-label">批次ID</span>
                    <span class="card-item-value">{{ item.batch_id || item.batchId || '-' }}</span>
                  </div>
                  <div class="card-item-row">
                    <span class="card-item-label">成本价</span>
                    <span class="card-item-value">¥{{ (item.cost_price || item.costPrice || 0).toFixed(2) }}</span>
                  </div>
                  <div class="card-item-row">
                    <span class="card-item-label">销售价</span>
                    <span class="card-item-value">¥{{ (item.sell_price || item.sellPrice || 0).toFixed(2) }}</span>
                  </div>
                  <div class="card-item-row">
                    <span class="card-item-label">同步状态</span>
                    <el-tag
                      :type="getSyncStatusType(item.sync_status || item.syncStatus)"
                      size="small"
                    >
                      {{ getSyncStatusText(item.sync_status || item.syncStatus) }}
                    </el-tag>
                  </div>
                  <div class="card-item-row">
                    <span class="card-item-label">创建时间</span>
                    <span class="card-item-value">{{ formatDate(item.created_at || item.createdAt) }}</span>
                  </div>
                  <div class="card-item-row" v-if="item.used_at || item.usedAt">
                    <span class="card-item-label">使用时间</span>
                    <span class="card-item-value">{{ formatDate(item.used_at || item.usedAt) }}</span>
                  </div>
                  <div class="card-item-row" v-if="item.expired_at || item.expiredAt">
                    <span class="card-item-label">过期时间</span>
                    <span class="card-item-value" :class="{ 'expired-text': isExpired(item.expired_at || item.expiredAt) }">
                      {{ formatDate(item.expired_at || item.expiredAt) }}
                    </span>
                  </div>
                </div>
                <div class="card-item-actions">
                  <el-button size="small" @click="handleView(item)">查看</el-button>
                  <el-button size="small" type="warning" @click="handleEdit(item)">编辑</el-button>
                  <el-button size="small" type="danger" @click="handleDelete(item)">删除</el-button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.size"
              :page-sizes="[10, 20, 50, 100]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              background
              @size-change="handleSearch"
              @current-change="handleSearch"
            />
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 添加/编辑卡片对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑卡片' : '添加卡片'"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="卡片代码" prop="cardCode">
          <el-input
            v-model="formData.cardCode"
            placeholder="请输入卡片代码"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="批次ID" prop="batchId">
          <el-input
            v-model="formData.batchId"
            placeholder="请输入批次ID"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 验证卡片对话框 -->
    <el-dialog
      v-model="validateDialogVisible"
      title="验证卡片"
      width="500px"
      class="validate-dialog"
    >
      <div class="validate-type-selector">
        <el-radio-group v-model="validateType" size="large">
          <el-radio-button value="single">
            <el-icon><CircleCheck /></el-icon>
            单个验证
          </el-radio-button>
          <el-radio-button value="batch">
            <el-icon><DocumentChecked /></el-icon>
            批量验证
          </el-radio-button>
        </el-radio-group>
      </div>
      
      <el-form class="validate-form">
        <!-- 单个验证 -->
        <template v-if="validateType === 'single'">
          <el-form-item label="卡片代码">
            <el-input
              v-model="validateCardCode"
              placeholder="请输入要验证的卡片代码"
              size="large"
              clearable
              @keyup.enter="handleValidateSubmit"
            >
              <template #prefix>
                <el-icon><CreditCard /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </template>
        
        <!-- 批量验证 -->
        <template v-else>
          <el-form-item label="验证方式">
            <el-radio-group v-model="batchValidateMode">
              <el-radio value="file">文件上传</el-radio>
              <el-radio value="text">文本输入</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <!-- 文件上传模式 -->
          <el-form-item v-if="batchValidateMode === 'file'" label="选择文件">
            <el-upload
              ref="uploadRef"
              class="validate-upload"
              drag
              :auto-upload="false"
              :limit="1"
              accept=".txt,.csv,.xlsx,.xls"
              :on-change="handleFileChange"
              :on-exceed="handleExceed"
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                将文件拖到此处，或<em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                  支持 txt/csv/xlsx 格式，每行一个卡片代码
                </div>
              </template>
            </el-upload>
          </el-form-item>
          
          <!-- 文本输入模式 -->
          <el-form-item v-else label="卡片代码">
            <el-input
              v-model="batchValidateCodes"
              type="textarea"
              :rows="8"
              placeholder="请输入卡片代码，每行一个"
            />
            <div class="input-tip">
              每行输入一个卡片代码，最多支持100个
            </div>
          </el-form-item>
        </template>
      </el-form>
      
      <template #footer>
        <el-button @click="validateDialogVisible = false">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleValidateSubmit" 
          :loading="validating"
          :disabled="!canValidate"
        >
          {{ validating ? '验证中...' : '开始验证' }}
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 验证结果 -->
    <ResultAnimation
      v-if="showValidateResult"
      :visible="showValidateResult"
      :type="validateResultType"
      :title="validateResultTitle"
      :description="validateResultDescription"
    />
    
    <!-- 卡片详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="卡片详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <div v-if="currentCard" class="card-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="卡片ID">
            {{ currentCard.id }}
          </el-descriptions-item>
          <el-descriptions-item label="卡片代码">
            <span class="card-code-display">{{ currentCard.card_code || currentCard.cardCode }}</span>
            <el-button 
              type="info"
              link 
              size="small" 
              @click="copyToClipboard(currentCard.card_code || currentCard.cardCode)"
              style="margin-left: 8px"
            >
              <el-icon><DocumentCopy /></el-icon>
              复制
            </el-button>
          </el-descriptions-item>
          <el-descriptions-item label="批次ID">
            {{ currentCard.batch_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="瑞幸产品ID">
            {{ currentCard.luckin_product_id || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="成本价">
            <span class="price-text">¥{{ currentCard.cost_price?.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="销售价">
            <span class="price-text">¥{{ currentCard.sell_price?.toFixed(2) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="使用状态">
            <el-tag :type="currentCard.status === 0 ? 'success' : 'info'">
              {{ currentCard.status === 0 ? '未使用' : '已使用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="同步状态">
            <el-tag :type="getSyncStatusType(currentCard.sync_status)">
              {{ getSyncStatusText(currentCard.sync_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDateTime(currentCard.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="使用时间">
            {{ currentCard.used_at ? formatDateTime(currentCard.used_at) : '未使用' }}
          </el-descriptions-item>
          <el-descriptions-item label="过期时间">
            {{ currentCard.expired_at ? formatDateTime(currentCard.expired_at) : '永不过期' }}
          </el-descriptions-item>
          <el-descriptions-item label="绑定产品数">
            {{ currentCard.bound_product_count || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ currentCard.description || '无' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
    
    <!-- 验证进度对话框 -->
    <el-dialog 
      v-model="progressDialogVisible" 
      title="批量验证进度" 
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <div class="validation-progress">
        <div class="progress-header">
          <el-tag :type="getStatusType(validationTask?.status)">
            {{ getStatusText(validationTask?.status) }}
          </el-tag>
          <span class="mode-info">{{ validationTask?.mode === 'smart' ? '智能验证' : '全量验证' }}</span>
        </div>
        
        <div v-if="validationTask?.progress" class="progress-content">
          <el-progress 
            :percentage="getProgressPercentage()" 
            :status="validationTask.status === 'completed' ? 'success' : 'active'"
            :stroke-width="12"
          />
          
          <div class="progress-stats">
            <div class="stat-item">
              <span class="stat-label">总计:</span>
              <span class="stat-value">{{ validationTask.progress.total }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">已处理:</span>
              <span class="stat-value">{{ validationTask.progress.processed }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">有效:</span>
              <span class="stat-value success">{{ validationTask.progress.valid }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">无效:</span>
              <span class="stat-value warning">{{ validationTask.progress.invalid }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">失败:</span>
              <span class="stat-value danger">{{ validationTask.progress.failed }}</span>
            </div>
          </div>
        </div>
        
        <div v-if="validationTask?.error" class="error-message">
          <el-alert type="error" :title="validationTask.error" :closable="false" />
        </div>
      </div>
      
      <template #footer>
        <el-button 
          v-if="validationTask?.status === 'running'" 
          type="danger" 
          @click="handleCancelValidation"
        >
          取消验证
        </el-button>
        <el-button 
          type="primary" 
          @click="progressDialogVisible = false"
          :disabled="validationTask?.status === 'running'"
        >
          {{ validationTask?.status === 'running' ? '验证中...' : '关闭' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, RefreshLeft, Plus, Upload, Files, CircleCheck,
  CreditCard, Clock, View, Edit, Delete, DocumentCopy,
  DocumentChecked, UploadFilled, ArrowDown, Lightning
} from '@element-plus/icons-vue'
import * as cardApi from '@/api/card'
import { startBatchValidation, getValidationProgress, cancelValidation } from '@/api/card'

const router = useRouter()

// 响应式数据
const pageLoading = ref(true)
const loading = ref(false)
const searching = ref(false)
const submitting = ref(false)
const validating = ref(false)
const showAdvancedSearch = ref(false)
const dialogVisible = ref(false)
const validateDialogVisible = ref(false)
const isEdit = ref(false)
const showValidateResult = ref(false)
const validateResultType = ref('success')
const validateResultTitle = ref('')
const validateResultDescription = ref('')

// 批量验证相关状态
const batchValidating = ref(false)
const progressDialogVisible = ref(false)
const validationTask = ref(null)
const progressTimer = ref(null)
const smartTaskRunning = ref(false)
const fullTaskRunning = ref(false)

// 卡片详情相关状态
const detailDialogVisible = ref(false)
const currentCard = ref(null)

// 表单数据
const searchForm = reactive({
  cardCode: '',
  status: '',
  batchId: '',
  dateRange: []
})

const formData = reactive({
  id: null,
  cardCode: '',
  batchId: ''
})

const validateCardCode = ref('')
const validateType = ref('single')
const batchValidateMode = ref('text')
const batchValidateCodes = ref('')
const uploadRef = ref()

// 表格数据
const tableData = ref([])
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 表单验证规则
const formRules = {
  cardCode: [
    { required: true, message: '请输入卡片代码', trigger: 'blur' }
  ]
}

// 方法
const handleToggleSearch = () => {
  showAdvancedSearch.value = !showAdvancedSearch.value
}

const handleSearch = async () => {
  searching.value = true
  loading.value = true
  
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.size,
      cardCode: searchForm.cardCode || undefined,
      status: searchForm.status !== '' ? searchForm.status : undefined,
      batchId: searchForm.batchId || undefined,
      startDate: searchForm.dateRange?.[0] || undefined,
      endDate: searchForm.dateRange?.[1] || undefined
    }
    
    const response = await cardApi.getCards(params)
    
    if (response.code === 200) {
      tableData.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.msg || '获取卡片数据失败')
    }
    
  } catch (error) {
    console.error('搜索失败:', error)
    ElMessage.error('搜索失败')
  } finally {
    searching.value = false
    loading.value = false
  }
}

const handleReset = () => {
  Object.assign(searchForm, {
    cardCode: '',
    status: '',
    batchId: '',
    dateRange: []
  })
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(formData, {
    id: null,
    cardCode: '',
    batchId: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(formData, {
    id: row.id,
    cardCode: row.cardCode,
    batchId: row.batchId
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  submitting.value = true
  
  try {
    if (isEdit.value) {
      const response = await cardApi.updateCard(formData.id, formData)
      if (response.code === 200) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        handleSearch()
      } else {
        ElMessage.error(response.msg || '更新失败')
      }
    } else {
      const response = await cardApi.createCard(formData)
      if (response.code === 200) {
        ElMessage.success('创建成功')
        dialogVisible.value = false
        handleSearch()
      } else {
        ElMessage.error(response.msg || '创建失败')
      }
    }
    
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

const handleView = (row) => {
  console.log('查看卡片:', row)
  currentCard.value = row
  detailDialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这张卡片吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const response = await cardApi.deleteCard(row.id)
    
    if (response.code === 200) {
      ElMessage.success('删除成功')
      handleSearch()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handleBatchImport = () => {
  router.push('/cards/batch-import')
}

const handleViewBatches = () => {
  router.push('/cards/batches')
}

const handleValidateCard = () => {
  validateCardCode.value = ''
  validateType.value = 'single'
  batchValidateMode.value = 'text'
  batchValidateCodes.value = ''
  validateDialogVisible.value = true
}

const handleBatchValidate = () => {
  validateCardCode.value = ''
  validateType.value = 'batch'
  batchValidateMode.value = 'text'
  batchValidateCodes.value = ''
  validateDialogVisible.value = true
}

// 处理验证模式选择
const handleValidationMode = async (mode) => {
  try {
    // 检查该模式是否已在运行
    if (mode === 'smart' && smartTaskRunning.value) {
      ElMessage.warning('智能验证正在进行中')
      return
    }
    
    if (mode === 'all' && fullTaskRunning.value) {
      ElMessage.warning('全量验证正在进行中')
      return
    }
    
    let confirmMessage = ''
    if (mode === 'smart') {
      confirmMessage = '智能验证只会检查重要的卡片（有异常、新添加、有订单的），大约需要2-5分钟。确定要继续吗？'
    } else {
      confirmMessage = '全量验证会检查所有卡片的真实状态，会在后台安全执行，可能需要较长时间。确定要继续吗？'
    }
    
    await ElMessageBox.confirm(confirmMessage, '确认验证', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 启动验证任务
    const res = await startBatchValidation(mode)
    
    // 设置任务信息并显示进度对话框
    validationTask.value = res.data
    progressDialogVisible.value = true
    
    // 更新运行状态
    if (mode === 'smart') {
      smartTaskRunning.value = true
    } else if (mode === 'all') {
      fullTaskRunning.value = true
    }
    
    batchValidating.value = true
    
    // 开始轮询进度
    startProgressPolling()
    
    ElMessage.success('验证任务已启动')
    
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('启动验证失败: ' + (error.message || error))
    }
  }
}

// 开始轮询进度
const startProgressPolling = () => {
  if (progressTimer.value) {
    clearInterval(progressTimer.value)
  }
  
  progressTimer.value = setInterval(async () => {
    try {
      const res = await getValidationProgress(validationTask.value.id)
      validationTask.value = res.data
      
      // 如果任务完成，停止轮询
      if (res.data.status === 'completed' || res.data.status === 'failed' || res.data.status === 'cancelled') {
        // 清理任务状态
        const taskMode = res.data.mode
        
        if (taskMode === 'smart') {
          smartTaskRunning.value = false
        } else if (taskMode === 'all') {
          fullTaskRunning.value = false
        }
        
        batchValidating.value = false
        clearInterval(progressTimer.value)
        progressTimer.value = null
        
        // 显示完成消息
        if (res.data.status === 'completed') {
          const progress = res.data.progress
          const modeText = taskMode === 'smart' ? '智能验证' : '全量验证'
          ElMessage.success(`${modeText}完成！总计：${progress.total}张，有效：${progress.valid}张，无效：${progress.invalid}张`)
          
          // 刷新列表
          await loadData()
        } else if (res.data.status === 'failed') {
          const modeText = taskMode === 'smart' ? '智能验证' : '全量验证'
          ElMessage.error(`${modeText}失败：${res.data.error}`)
        }
      }
    } catch (error) {
      console.error('获取验证进度失败:', error)
    }
  }, 2000) // 每2秒轮询一次
}

// 取消验证任务
const handleCancelValidation = async () => {
  try {
    await cancelValidation(validationTask.value.id)
    ElMessage.success('验证任务已取消')
  } catch (error) {
    ElMessage.error('取消验证失败')
  }
}

// 获取进度百分比
const getProgressPercentage = () => {
  if (!validationTask.value?.progress) return 0
  const { total, processed } = validationTask.value.progress
  if (total === 0) return 0
  return Math.round((processed / total) * 100)
}

// 计算属性：是否可以验证
const canValidate = computed(() => {
  if (validateType.value === 'single') {
    return validateCardCode.value.trim() !== ''
  } else {
    if (batchValidateMode.value === 'text') {
      return batchValidateCodes.value.trim() !== ''
    }
    return uploadRef.value?.uploadFiles?.length > 0
  }
})

const handleFileChange = (file) => {
  // 文件变更处理
  console.log('File selected:', file)
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const handleValidateSubmit = async () => {
  if (validateType.value === 'single') {
    // 单个验证
    if (!validateCardCode.value) {
      ElMessage.warning('请输入卡片代码')
      return
    }
    
    validating.value = true
    
    try {
      const response = await cardApi.validateCard({
        cardCode: validateCardCode.value
      })
      
      validateDialogVisible.value = false
      
      if (response.code === 200) {
        if (response.data.isValid) {
          showValidateSuccess()
        } else {
          showValidateError(response.data.message || '卡片无效或已被使用')
        }
      } else {
        showValidateError(response.msg || '验证失败')
      }
      
    } catch (error) {
      console.error('验证失败:', error)
      showValidateError('网络错误')
    } finally {
      validating.value = false
    }
  } else {
    // 批量验证
    let cardCodes = []
    
    if (batchValidateMode.value === 'text') {
      // 文本输入模式
      cardCodes = batchValidateCodes.value
        .split('\n')
        .map(code => code.trim())
        .filter(code => code)
      
      if (cardCodes.length === 0) {
        ElMessage.warning('请输入卡片代码')
        return
      }
      
      if (cardCodes.length > 100) {
        ElMessage.warning('批量验证最多支持100个卡片')
        return
      }
    } else {
      // 文件上传模式
      if (!uploadRef.value?.uploadFiles?.length) {
        ElMessage.warning('请选择要上传的文件')
        return
      }
      
      // TODO: 解析文件内容获取卡片代码
      ElMessage.info('文件解析功能开发中...')
      return
    }
    
    validating.value = true
    
    try {
      const response = await cardApi.batchValidateCards({
        cardCodes: cardCodes
      })
      
      validateDialogVisible.value = false
      
      if (response.code === 200) {
        showBatchValidateResult(response.data)
      } else {
        ElMessage.error(response.msg || '批量验证失败')
      }
      
    } catch (error) {
      console.error('批量验证失败:', error)
      ElMessage.error('批量验证失败')
    } finally {
      validating.value = false
    }
  }
}

const showValidateSuccess = () => {
  validateResultType.value = 'success'
  validateResultTitle.value = '卡片验证成功'
  validateResultDescription.value = '该卡片有效且可以使用'
  showValidateResult.value = true
  
  setTimeout(() => {
    showValidateResult.value = false
  }, 3000)
}

const showValidateError = (message = '卡片无效或已被使用') => {
  validateResultType.value = 'error'
  validateResultTitle.value = '卡片验证失败'
  validateResultDescription.value = message
  showValidateResult.value = true
  
  setTimeout(() => {
    showValidateResult.value = false
  }, 3000)
}

const showBatchValidateResult = (result) => {
  const { success, failed, details } = result
  const total = success + failed
  
  ElMessageBox.alert(
    `<div class="batch-validate-result">
      <div class="result-summary">
        <div class="summary-item">
          <span class="label">验证总数：</span>
          <span class="value">${total}</span>
        </div>
        <div class="summary-item success">
          <span class="label">有效卡片：</span>
          <span class="value">${success}</span>
        </div>
        <div class="summary-item error">
          <span class="label">无效卡片：</span>
          <span class="value">${failed}</span>
        </div>
      </div>
      ${details && details.length > 0 ? `
        <div class="result-details">
          <h4>详细结果：</h4>
          <div class="details-list">
            ${details.slice(0, 10).map(item => `
              <div class="detail-item ${item.isValid ? 'valid' : 'invalid'}">
                <span class="code">${item.cardCode}</span>
                <span class="status">${item.isValid ? '有效' : item.message || '无效'}</span>
              </div>
            `).join('')}
            ${details.length > 10 ? `<div class="more-tip">...还有 ${details.length - 10} 条记录</div>` : ''}
          </div>
        </div>
      ` : ''}
    </div>`,
    '批量验证结果',
    {
      confirmButtonText: '确定',
      dangerouslyUseHTMLString: true,
      customClass: 'batch-validate-dialog'
    }
  )
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 卡片状态处理函数
const getStatusType = (status) => {
  const types = {
    0: 'success',  // 未使用
    1: 'info',     // 已使用
    2: 'warning'   // 预占中
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    0: '未使用',
    1: '已使用',
    2: '预占中'
  }
  return texts[status] || '未知'
}

// 同步状态处理函数
const getSyncStatusType = (syncStatus) => {
  const types = {
    'pending': 'warning',
    'syncing': 'info',
    'synced': 'success',
    'failed': 'danger'
  }
  return types[syncStatus] || 'info'
}

const getSyncStatusText = (syncStatus) => {
  const texts = {
    'pending': '待同步',
    'syncing': '同步中',
    'synced': '已同步',
    'failed': '同步失败'
  }
  return texts[syncStatus] || '未知'
}

// 检查是否过期
const isExpired = (expiredAt) => {
  if (!expiredAt) return false
  return new Date(expiredAt) < new Date()
}

// 格式化时间 - 显示完整的日期时间
const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'
  try {
    const date = new Date(dateTime)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    return '-'
  }
}

// 加载数据函数
const loadData = async () => {
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.size,
      cardCode: searchForm.cardCode || undefined,
      status: searchForm.status !== '' ? searchForm.status : undefined,
      batchId: searchForm.batchId || undefined,
      startDate: searchForm.dateRange?.[0] || undefined,
      endDate: searchForm.dateRange?.[1] || undefined
    }
    
    const response = await cardApi.getCards(params)
    
    if (response.code === 200) {
      tableData.value = response.data.list || []
      pagination.total = response.data.total || 0
    } else {
      ElMessage.error(response.msg || '获取卡片数据失败')
    }
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败')
  }
}

// 组件挂载
onMounted(() => {
  pageLoading.value = false
  
  // 加载数据
  handleSearch()
})
</script>

<style lang="scss" scoped>
.cards-enhanced-container {
  padding: var(--spacing-6);
  background: var(--bg-page);
  min-height: 100vh;
  
  .cards-content {
    .search-card,
    .toolbar-card,
    .table-card {
      margin-bottom: var(--spacing-4);
      border: none;
      border-radius: var(--radius-lg);
      box-shadow: var(--shadow-sm);
      transition: all 0.3s ease;
      
      &:hover {
        box-shadow: var(--shadow-md);
      }
    }
    
    .search-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      
      .search-title {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        margin: 0;
        font-size: 1.125rem;
        font-weight: 600;
        color: var(--text-primary);
        
        .el-icon {
          color: var(--primary);
        }
      }
      
      .toggle-btn {
        color: var(--primary);
      }
    }
    
    .search-form {
      &.advanced-search {
        .el-form-item {
          margin-bottom: var(--spacing-4);
        }
      }
      
      .search-actions {
        .el-button {
          margin-right: var(--spacing-2);
        }
      }
      
      .option-with-status {
        display: flex;
        align-items: center;
        width: 100%;
      }
    }
    
    .toolbar-content {
      display: flex;
      align-items: center;
      justify-content: space-between;
      flex-wrap: wrap;
      gap: var(--spacing-4);
      
      .toolbar-left {
        display: flex;
        align-items: center;
        gap: var(--spacing-4);
        
        .toolbar-title {
          display: flex;
          align-items: center;
          gap: var(--spacing-2);
          margin: 0;
          font-size: 1.125rem;
          font-weight: 600;
          color: var(--text-primary);
          
          .el-icon {
            color: var(--primary);
          }
        }
        
        .stats-info {
          display: flex;
          gap: var(--spacing-2);
        }
      }
    }
    
    .table-loading {
      padding: var(--spacing-8);
    }
    
    .enhanced-table {
      .card-code {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        
        .card-icon {
          color: var(--primary);
        }
        
        .code-text {
          font-family: 'JetBrains Mono', 'Fira Code', monospace;
          font-weight: 600;
        }
      }
      
      .time-cell {
        display: flex;
        align-items: center;
        gap: var(--spacing-1);
        
        .time-icon {
          color: var(--text-tertiary);
        }
      }
      
      .unused-text {
        color: var(--text-tertiary);
        font-style: italic;
      }
      
      .table-actions {
        display: flex;
        gap: var(--spacing-1);
        
        .el-button {
          padding: var(--spacing-1) var(--spacing-2);
        }
      }
    }
    
    .pagination-container {
      display: flex;
      justify-content: center;
      margin-top: var(--spacing-6);
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .cards-enhanced-container {
    padding: var(--mobile-space-4);
    
    .toolbar-content {
      flex-direction: column;
      align-items: stretch;
      
      .toolbar-left {
        justify-content: space-between;
      }
      
      .toolbar-right {
        .el-button-group {
          display: flex;
          flex-wrap: wrap;
          gap: var(--mobile-space-2);
          
          .el-button {
            flex: 1;
            min-width: 120px;
          }
        }
      }
    }
  }
}

// 验证对话框样式
.validate-dialog {
  .validate-type-selector {
    margin-bottom: var(--spacing-6);
    text-align: center;
    
    .el-radio-group {
      width: 100%;
      
      .el-radio-button {
        width: 50%;
        
        .el-radio-button__inner {
          width: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
          gap: var(--spacing-2);
          padding: var(--spacing-3) 0;
          font-size: 1rem;
          
          .el-icon {
            font-size: 1.2rem;
          }
        }
      }
    }
  }
  
  .validate-form {
    .el-form-item {
      margin-bottom: var(--spacing-4);
    }
    
    .validate-upload {
      width: 100%;
      
      .el-upload-dragger {
        padding: var(--spacing-6);
        
        .el-icon--upload {
          font-size: 3rem;
          margin-bottom: var(--spacing-3);
        }
      }
    }
    
    .input-tip {
      margin-top: var(--spacing-2);
      font-size: 0.875rem;
      color: var(--text-tertiary);
    }
  }
}

// 批量验证结果样式
:global(.batch-validate-dialog) {
  .batch-validate-result {
    .result-summary {
      display: flex;
      justify-content: space-around;
      padding: var(--spacing-4);
      background: var(--bg-secondary);
      border-radius: var(--radius-lg);
      margin-bottom: var(--spacing-4);
      
      .summary-item {
        text-align: center;
        
        .label {
          display: block;
          font-size: 0.875rem;
          color: var(--text-secondary);
          margin-bottom: var(--spacing-1);
        }
        
        .value {
          display: block;
          font-size: 1.5rem;
          font-weight: 700;
          color: var(--text-primary);
        }
        
        &.success .value {
          color: var(--success);
        }
        
        &.error .value {
          color: var(--danger);
        }
      }
    }
    
    .result-details {
      h4 {
        margin: 0 0 var(--spacing-3) 0;
        font-size: 1rem;
        color: var(--text-primary);
      }
      
      .details-list {
        max-height: 300px;
        overflow-y: auto;
        
        .detail-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: var(--spacing-2) var(--spacing-3);
          margin-bottom: var(--spacing-2);
          border-radius: var(--radius-md);
          background: var(--bg-secondary);
          
          &.valid {
            border-left: 3px solid var(--success);
            
            .status {
              color: var(--success);
            }
          }
          
          &.invalid {
            border-left: 3px solid var(--danger);
            
            .status {
              color: var(--danger);
            }
          }
          
          .code {
            font-family: monospace;
            font-weight: 500;
          }
          
          .status {
            font-size: 0.875rem;
          }
        }
        
        .more-tip {
          text-align: center;
          padding: var(--spacing-3);
          color: var(--text-tertiary);
          font-style: italic;
        }
      }
    }
  }
}

// 修复表格显示问题
:deep(.el-table) {
  width: 100% !important;
  
  .el-table__header-wrapper,
  .el-table__body-wrapper {
    width: 100% !important;
  }
  
  .el-table__header,
  .el-table__body {
    width: 100% !important;
  }
  
  // 确保固定列正常显示
  .el-table__fixed,
  .el-table__fixed-right {
    height: 100% !important;
  }
  
  // 优化表格单元格内容显示
  .el-table__cell {
    .cell {
      display: flex;
      align-items: center;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
  
  // 价格文本样式
  .price-text {
    font-weight: 600;
    color: var(--primary);
  }
  
  // 时间单元格样式
  .time-cell {
    .expired-text {
      color: var(--danger);
    }
  }
  
  // 未使用文本样式
  .unused-text,
  .no-expire-text {
    color: var(--text-tertiary);
    font-style: italic;
  }
}

// 卡片详情对话框样式
.card-detail {
  .card-code-display {
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-weight: 600;
    color: var(--primary);
  }
  
  .price-text {
    font-weight: 600;
    color: var(--primary);
  }
}

// 增强表格样式
.enhanced-table {
  border-radius: var(--radius-lg);
  overflow: hidden;
  
  :deep(.el-table__empty-block) {
    padding: var(--spacing-8) 0;
  }
  
  :deep(.el-loading-mask) {
    background-color: rgba(255, 255, 255, 0.9);
  }
}

// 验证模式选择样式
.validation-mode-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 0;
  
  .mode-title {
    font-weight: 500;
    font-size: 14px;
    line-height: 1.2;
    color: var(--text-primary);
  }
  
  .mode-desc {
    font-size: 12px;
    color: var(--text-secondary);
    line-height: 1.2;
    margin-top: 2px;
  }
}

// 验证进度样式
.validation-progress {
  .progress-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
    
    .mode-info {
      font-size: 14px;
      color: #606266;
    }
  }
  
  .progress-content {
    .progress-stats {
      display: grid;
      grid-template-columns: repeat(5, 1fr);
      gap: 16px;
      margin: 20px 0;
      
      .stat-item {
        text-align: center;
        
        .stat-label {
          display: block;
          font-size: 12px;
          color: #909399;
          margin-bottom: 4px;
        }
        
        .stat-value {
          display: block;
          font-size: 18px;
          font-weight: 600;
          
          &.success { color: #67c23a; }
          &.warning { color: #e6a23c; }
          &.danger { color: #f56c6c; }
        }
      }
    }
  }
  
  .error-message {
    margin-top: 16px;
  }
}
</style>