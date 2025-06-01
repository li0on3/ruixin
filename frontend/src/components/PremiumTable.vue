<template>
  <div class="premium-table" :class="tableClasses">
    <!-- 表格工具栏 -->
    <div class="table-toolbar" v-if="showToolbar">
      <div class="toolbar-content">
        <!-- 左侧工具 -->
        <div class="toolbar-left">
          <slot name="toolbar-left">
            <!-- 表格标题 -->
            <div class="table-title" v-if="title">
              <h3 class="title-text">{{ title }}</h3>
              <p class="title-subtitle" v-if="subtitle">{{ subtitle }}</p>
            </div>
            
            <!-- 搜索框 -->
            <div class="search-container" v-if="searchable">
              <el-input
                v-model="searchText"
                :placeholder="searchPlaceholder"
                :prefix-icon="Search"
                clearable
                @clear="handleSearch"
                @keyup.enter="handleSearch"
                @input="handleSearchInput"
                class="table-search"
              >
                <template #append v-if="showAdvancedSearch">
                  <el-button 
                    :icon="Filter" 
                    @click="showAdvancedFilter = !showAdvancedFilter"
                    :class="{ 'is-active': showAdvancedFilter }"
                    title="高级筛选"
                  />
                </template>
              </el-input>
            </div>
          </slot>
        </div>
        
        <!-- 右侧工具 -->
        <div class="toolbar-right">
          <slot name="toolbar-right">
            <!-- 表格配置 -->
            <div class="table-config">
              <!-- 密度设置 -->
              <el-dropdown v-if="densitySettable" trigger="click" placement="bottom-end">
                <el-button 
                  circle 
                  size="small" 
                  :icon="Grid" 
                  title="表格密度"
                  class="config-btn"
                />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item 
                      v-for="density in densityOptions" 
                      :key="density.value"
                      @click="currentDensity = density.value"
                      :class="{ 'is-active': currentDensity === density.value }"
                    >
                      <el-icon v-if="currentDensity === density.value"><Check /></el-icon>
                      <span>{{ density.label }}</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              
              <!-- 列设置 -->
              <el-popover v-if="columnSettable" placement="bottom-end" :width="280" trigger="click">
                <template #reference>
                  <el-button 
                    circle 
                    size="small" 
                    :icon="Setting" 
                    title="列设置"
                    class="config-btn"
                  />
                </template>
                <div class="column-setting">
                  <div class="setting-header">
                    <h4>列显示设置</h4>
                    <el-button text size="small" @click="resetColumns">重置</el-button>
                  </div>
                  <div class="setting-content">
                    <div class="setting-search">
                      <el-input
                        v-model="columnSearchText"
                        placeholder="搜索列名"
                        :prefix-icon="Search"
                        size="small"
                        clearable
                      />
                    </div>
                    <div class="column-list custom-scrollbar">
                      <el-checkbox-group v-model="visibleColumns" @change="handleColumnChange">
                        <div 
                          class="column-item" 
                          v-for="col in filteredSettableColumns" 
                          :key="col.prop"
                        >
                          <el-checkbox 
                            :value="col.prop"
                            :disabled="col.fixed"
                          >
                            <span class="column-label">{{ col.label }}</span>
                            <el-tag v-if="col.fixed" size="small" type="info">固定</el-tag>
                          </el-checkbox>
                        </div>
                      </el-checkbox-group>
                    </div>
                  </div>
                </div>
              </el-popover>
              
              <!-- 刷新按钮 -->
              <el-button 
                v-if="refreshable"
                circle 
                size="small" 
                :icon="Refresh" 
                @click="handleRefresh"
                title="刷新数据"
                class="config-btn"
                :loading="refreshLoading"
              />
              
              <!-- 导出功能 -->
              <el-dropdown v-if="exportable" trigger="click" placement="bottom-end">
                <el-button 
                  circle 
                  size="small" 
                  :icon="Download" 
                  title="导出数据"
                  class="config-btn"
                />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="handleExport('excel')">
                      <el-icon><Document /></el-icon>
                      导出为 Excel
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleExport('csv')">
                      <el-icon><Document /></el-icon>
                      导出为 CSV
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleExport('pdf')">
                      <el-icon><Document /></el-icon>
                      导出为 PDF
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </slot>
        </div>
      </div>
    </div>
    
    <!-- 高级筛选面板 -->
    <transition name="filter-slide">
      <div class="advanced-filter" v-if="showAdvancedFilter && filterConfig.length > 0">
        <div class="filter-header">
          <h4>高级筛选</h4>
          <div class="filter-actions">
            <el-button size="small" @click="handleResetFilter">重置</el-button>
            <el-button type="primary" size="small" @click="handleFilter">筛选</el-button>
            <el-button text size="small" @click="showAdvancedFilter = false">收起</el-button>
          </div>
        </div>
        
        <div class="filter-content">
          <el-form :model="filterForm" label-width="100px" size="small">
            <el-row :gutter="20">
              <el-col 
                v-for="filter in filterConfig" 
                :key="filter.prop"
                :span="filter.span || 8"
              >
                <el-form-item :label="filter.label">
                  <!-- 文本输入 -->
                  <el-input
                    v-if="filter.type === 'text'"
                    v-model="filterForm[filter.prop]"
                    :placeholder="filter.placeholder"
                    clearable
                  />
                  
                  <!-- 选择器 -->
                  <el-select
                    v-else-if="filter.type === 'select'"
                    v-model="filterForm[filter.prop]"
                    :placeholder="filter.placeholder"
                    clearable
                    :multiple="filter.multiple"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="option in filter.options"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                  
                  <!-- 日期选择 -->
                  <el-date-picker
                    v-else-if="filter.type === 'date'"
                    v-model="filterForm[filter.prop]"
                    type="date"
                    :placeholder="filter.placeholder"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                  
                  <!-- 日期范围 -->
                  <el-date-picker
                    v-else-if="filter.type === 'daterange'"
                    v-model="filterForm[filter.prop]"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                  
                  <!-- 数字范围 -->
                  <div v-else-if="filter.type === 'numberRange'" class="number-range">
                    <el-input-number
                      v-model="filterForm[filter.prop + '_min']"
                      :placeholder="filter.minPlaceholder || '最小值'"
                      :controls="false"
                      size="small"
                    />
                    <span class="range-separator">-</span>
                    <el-input-number
                      v-model="filterForm[filter.prop + '_max']"
                      :placeholder="filter.maxPlaceholder || '最大值'"
                      :controls="false"
                      size="small"
                    />
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </div>
      </div>
    </transition>
    
    <!-- 批量操作栏 -->
    <transition name="batch-slide">
      <div class="batch-actions" v-if="selectable && selectedRows.length > 0">
        <div class="batch-info">
          <el-checkbox 
            :indeterminate="selectedRows.length > 0 && selectedRows.length < displayData.length"
            :model-value="selectedRows.length === displayData.length"
            @change="toggleAllSelection"
          />
          <span class="selected-count">已选择 {{ selectedRows.length }} 项</span>
        </div>
        
        <div class="batch-operations">
          <slot name="batch-actions" :selectedRows="selectedRows">
            <el-button size="small" @click="clearSelection">取消选择</el-button>
            <el-button type="danger" size="small" @click="handleBatchDelete">批量删除</el-button>
          </slot>
        </div>
      </div>
    </transition>
    
    <!-- 表格主体 -->
    <div class="table-container" :class="[`density-${currentDensity}`]">
      <el-table
        ref="tableRef"
        :data="displayData"
        :height="height"
        :max-height="maxHeight"
        :stripe="stripe"
        :border="border"
        :show-summary="showSummary"
        :summary-method="summaryMethod"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
        @row-click="handleRowClick"
        @row-dblclick="handleRowDblClick"
        :row-class-name="getRowClassName"
        :cell-class-name="getCellClassName"
        :header-cell-class-name="getHeaderCellClassName"
        :default-sort="defaultSort"
        :element-loading-text="loadingText"
        :element-loading-spinner="LoadingIcon"
        class="premium-el-table"
      >
        <!-- 选择列 -->
        <el-table-column
          v-if="selectable"
          type="selection"
          width="55"
          fixed="left"
          :selectable="selectableMethod"
        />
        
        <!-- 序号列 -->
        <el-table-column
          v-if="showIndex"
          type="index"
          :label="indexLabel"
          width="70"
          :fixed="indexFixed"
          :index="indexMethod"
          align="center"
        />
        
        <!-- 数据列 -->
        <template v-for="column in displayColumns" :key="column.prop">
          <el-table-column
            v-if="!column.hide && isColumnVisible(column.prop)"
            :prop="column.prop"
            :label="column.label"
            :width="column.width"
            :min-width="column.minWidth || 120"
            :fixed="column.fixed"
            :align="column.align || 'left'"
            :sortable="column.sortable"
            :sort-method="column.sortMethod"
            :resizable="column.resizable !== false"
            :show-overflow-tooltip="column.showOverflowTooltip !== false"
            :class-name="column.className"
            :label-class-name="column.labelClassName"
            :formatter="column.formatter"
          >
            <!-- 表头插槽 -->
            <template v-if="column.headerSlot" #header="scope">
              <slot :name="column.headerSlot" v-bind="scope" />
            </template>
            
            <!-- 单元格插槽 -->
            <template v-if="column.slot" #default="scope">
              <slot :name="column.slot" v-bind="scope" />
            </template>
            
            <!-- 渲染函数 -->
            <template v-else-if="column.render" #default="scope">
              <component :is="column.render(scope.row, scope.$index, scope.column)" />
            </template>
          </el-table-column>
        </template>
        
        <!-- 操作列 -->
        <el-table-column
          v-if="actions && actions.length > 0"
          :label="actionLabel"
          :width="actionWidth"
          :fixed="actionFixed"
          :align="actionAlign"
          class-name="table-actions-column"
        >
          <template #default="scope">
            <div class="table-actions">
              <template v-for="(action, index) in getRowActions(scope.row, scope.$index)" :key="index">
                <!-- 普通按钮 -->
                <el-button
                  v-if="!action.hide && action.type !== 'dropdown'"
                  :type="action.buttonType || 'text'"
                  :size="action.size || 'small'"
                  :icon="action.icon"
                  :loading="action.loading"
                  :disabled="action.disabled"
                  @click="handleAction(action, scope.row, scope.$index)"
                  :title="action.tooltip || action.text"
                  :class="action.class"
                >
                  {{ action.text }}
                </el-button>
                
                <!-- 下拉菜单 -->
                <el-dropdown
                  v-else-if="!action.hide && action.type === 'dropdown'"
                  trigger="click"
                  size="small"
                  placement="bottom-end"
                >
                  <el-button
                    :type="action.buttonType || 'text'"
                    :size="action.size || 'small'"
                    :icon="action.icon || MoreFilled"
                    :disabled="action.disabled"
                    :class="action.class"
                  >
                    {{ action.text || '更多' }}
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item
                        v-for="(item, itemIndex) in action.items"
                        :key="itemIndex"
                        @click="handleAction(item, scope.row, scope.$index)"
                        :disabled="item.disabled"
                        :divided="item.divided"
                      >
                        <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
                        {{ item.text }}
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </div>
          </template>
        </el-table-column>
        
        <!-- 空数据插槽 -->
        <template #empty>
          <div class="table-empty">
            <slot name="empty">
              <div class="empty-content">
                <el-icon class="empty-icon"><Box /></el-icon>
                <p class="empty-text">{{ emptyText }}</p>
                <el-button v-if="emptyAction" type="primary" @click="handleEmptyAction">
                  {{ emptyActionText }}
                </el-button>
              </div>
            </slot>
          </div>
        </template>
      </el-table>
    </div>
    
    <!-- 分页器 -->
    <div class="table-pagination" v-if="pagination && totalCount > 0">
      <div class="pagination-info">
        <span class="total-info">
          共 {{ totalCount }} 条记录，第 {{ currentPage }} / {{ totalPages }} 页
        </span>
      </div>
      
      <el-pagination
        :current-page="currentPage"
        :page-size="currentPageSize"
        :total="totalCount"
        :page-sizes="pageSizes"
        :layout="paginationLayout"
        :background="true"
        :hide-on-single-page="hideOnSinglePage"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination-component"
      />
    </div>
    
    <!-- 表格统计信息 -->
    <div class="table-footer" v-if="showFooter">
      <div class="footer-stats">
        <span class="stats-item">总计：{{ totalCount }} 条</span>
        <span class="stats-item" v-if="selectable">已选：{{ selectedRows.length }} 条</span>
        <span class="stats-item" v-if="loadTime">耗时：{{ loadTime }}ms</span>
      </div>
      
      <div class="footer-actions">
        <slot name="footer-actions"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Filter, Grid, Setting, Refresh, Download, Check, Document,
  MoreFilled, Box, Loading
} from '@element-plus/icons-vue'

// Props 定义
const props = defineProps({
  // 基础数据
  data: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  loadingText: {
    type: String,
    default: '数据加载中...'
  },
  
  // 表格外观
  title: String,
  subtitle: String,
  height: [String, Number],
  maxHeight: [String, Number],
  stripe: {
    type: Boolean,
    default: true
  },
  border: {
    type: Boolean,
    default: false
  },
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['large', 'default', 'small'].includes(value)
  },
  
  // 功能开关
  showToolbar: {
    type: Boolean,
    default: true
  },
  searchable: {
    type: Boolean,
    default: true
  },
  showAdvancedSearch: {
    type: Boolean,
    default: true
  },
  refreshable: {
    type: Boolean,
    default: true
  },
  columnSettable: {
    type: Boolean,
    default: true
  },
  densitySettable: {
    type: Boolean,
    default: true
  },
  exportable: {
    type: Boolean,
    default: true
  },
  
  // 选择功能
  selectable: {
    type: Boolean,
    default: false
  },
  selectableMethod: Function,
  
  // 序号列
  showIndex: {
    type: Boolean,
    default: false
  },
  indexLabel: {
    type: String,
    default: '序号'
  },
  indexFixed: {
    type: [Boolean, String],
    default: false
  },
  
  // 汇总功能
  showSummary: {
    type: Boolean,
    default: false
  },
  summaryMethod: Function,
  
  // 分页
  pagination: {
    type: Boolean,
    default: true
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 20
  },
  total: {
    type: Number,
    default: 0
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  },
  paginationLayout: {
    type: String,
    default: 'total, sizes, prev, pager, next, jumper'
  },
  hideOnSinglePage: {
    type: Boolean,
    default: false
  },
  
  // 操作列
  actions: {
    type: Array,
    default: () => []
  },
  actionLabel: {
    type: String,
    default: '操作'
  },
  actionWidth: {
    type: [String, Number],
    default: 'auto'
  },
  actionFixed: {
    type: [Boolean, String],
    default: 'right'
  },
  actionAlign: {
    type: String,
    default: 'center'
  },
  
  // 筛选
  filterConfig: {
    type: Array,
    default: () => []
  },
  
  // 其他配置
  emptyText: {
    type: String,
    default: '暂无数据'
  },
  emptyAction: {
    type: Boolean,
    default: false
  },
  emptyActionText: {
    type: String,
    default: '添加数据'
  },
  searchPlaceholder: {
    type: String,
    default: '请输入搜索关键词'
  },
  showFooter: {
    type: Boolean,
    default: false
  },
  loadTime: Number,
  
  // 样式类名
  rowClassName: [String, Function],
  cellClassName: [String, Function],
  headerCellClassName: [String, Function],
  defaultSort: Object
})

// Emits 定义
const emit = defineEmits([
  'search',
  'refresh',
  'filter',
  'sort-change',
  'selection-change',
  'page-change',
  'size-change',
  'export',
  'action-click',
  'row-click',
  'row-dblclick',
  'empty-action',
  'batch-delete'
])

// 响应式数据
const tableRef = ref()
const searchText = ref('')
const currentPage = ref(props.currentPage)
const currentPageSize = ref(props.pageSize)
const currentDensity = ref(localStorage.getItem('table-density') || 'default')
const showAdvancedFilter = ref(false)
const filterForm = ref({})
const visibleColumns = ref([])
const selectedRows = ref([])
const columnSearchText = ref('')
const refreshLoading = ref(false)

// 密度选项
const densityOptions = [
  { label: '紧凑', value: 'compact' },
  { label: '标准', value: 'default' },
  { label: '宽松', value: 'comfortable' }
]

// 计算属性
const tableClasses = computed(() => ({
  'with-toolbar': props.showToolbar,
  'with-selection': props.selectable,
  'with-pagination': props.pagination,
  [`size-${props.size}`]: true,
  [`density-${currentDensity.value}`]: true
}))

const totalCount = computed(() => props.total || props.data.length)

const totalPages = computed(() => Math.ceil(totalCount.value / currentPageSize.value))

const displayData = computed(() => {
  if (!props.pagination) {
    return props.data
  }
  const start = (currentPage.value - 1) * currentPageSize.value
  const end = start + currentPageSize.value
  return props.data.slice(start, end)
})

const settableColumns = computed(() => {
  return props.columns.filter(col => !col.type && col.prop && !col.hidden)
})

const filteredSettableColumns = computed(() => {
  if (!columnSearchText.value) return settableColumns.value
  return settableColumns.value.filter(col => 
    col.label.toLowerCase().includes(columnSearchText.value.toLowerCase())
  )
})

const displayColumns = computed(() => {
  return props.columns.filter(col => !col.type)
})

// 方法
const handleSearch = () => {
  emit('search', searchText.value)
}

const handleSearchInput = (value) => {
  // 实时搜索
  if (value.length === 0 || value.length >= 2) {
    emit('search', value)
  }
}

const handleRefresh = async () => {
  refreshLoading.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 500)) // 模拟异步操作
    emit('refresh')
  } finally {
    refreshLoading.value = false
  }
}

const handleFilter = () => {
  emit('filter', filterForm.value)
}

const handleResetFilter = () => {
  filterForm.value = {}
  emit('filter', {})
}

const handleSortChange = ({ column, prop, order }) => {
  emit('sort-change', { column, prop, order })
}

const handleSelectionChange = (selection) => {
  selectedRows.value = selection
  emit('selection-change', selection)
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  emit('page-change', page)
}

const handleSizeChange = (size) => {
  currentPageSize.value = size
  currentPage.value = 1
  emit('size-change', size)
  emit('page-change', 1)
}

const handleExport = (type) => {
  const exportData = {
    type,
    data: selectedRows.value.length > 0 ? selectedRows.value : props.data,
    columns: props.columns.filter(col => isColumnVisible(col.prop)),
    filename: `export_${new Date().getTime()}.${type}`
  }
  emit('export', exportData)
}

const handleAction = (action, row, index) => {
  if (action.handler) {
    action.handler(row, index, action)
  }
  emit('action-click', { action, row, index })
}

const handleRowClick = (row, column, event) => {
  emit('row-click', { row, column, event })
}

const handleRowDblClick = (row, column, event) => {
  emit('row-dblclick', { row, column, event })
}

const handleColumnChange = () => {
  localStorage.setItem('table-visible-columns', JSON.stringify(visibleColumns.value))
}

const resetColumns = () => {
  visibleColumns.value = settableColumns.value.map(col => col.prop)
  handleColumnChange()
  ElMessage.success('列设置已重置')
}

const isColumnVisible = (prop) => {
  if (!prop || visibleColumns.value.length === 0) return true
  return visibleColumns.value.includes(prop)
}

const getRowActions = (row, index) => {
  return props.actions.map(action => {
    if (typeof action === 'function') {
      return action(row, index)
    }
    return { ...action }
  }).filter(action => !action.hide)
}

const indexMethod = (index) => {
  return (currentPage.value - 1) * currentPageSize.value + index + 1
}

const getRowClassName = ({ row, rowIndex }) => {
  let className = 'table-row'
  
  if (props.rowClassName) {
    if (typeof props.rowClassName === 'function') {
      className += ' ' + props.rowClassName({ row, rowIndex })
    } else {
      className += ' ' + props.rowClassName
    }
  }
  
  // 添加悬停效果类
  className += ' hover-effect'
  
  return className
}

const getCellClassName = ({ row, column, rowIndex, columnIndex }) => {
  let className = 'table-cell'
  
  if (props.cellClassName) {
    if (typeof props.cellClassName === 'function') {
      className += ' ' + props.cellClassName({ row, column, rowIndex, columnIndex })
    } else {
      className += ' ' + props.cellClassName
    }
  }
  
  return className
}

const getHeaderCellClassName = ({ column, columnIndex }) => {
  let className = 'table-header-cell'
  
  if (props.headerCellClassName) {
    if (typeof props.headerCellClassName === 'function') {
      className += ' ' + props.headerCellClassName({ column, columnIndex })
    } else {
      className += ' ' + props.headerCellClassName
    }
  }
  
  return className
}

const handleEmptyAction = () => {
  emit('empty-action')
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 条记录吗？`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    emit('batch-delete', selectedRows.value)
  } catch (error) {
    // 用户取消
  }
}

// 公开方法
const clearSelection = () => {
  tableRef.value?.clearSelection()
}

const toggleRowSelection = (row, selected) => {
  tableRef.value?.toggleRowSelection(row, selected)
}

const toggleAllSelection = () => {
  tableRef.value?.toggleAllSelection()
}

const toggleRowExpansion = (row, expanded) => {
  tableRef.value?.toggleRowExpansion(row, expanded)
}

const setCurrentRow = (row) => {
  tableRef.value?.setCurrentRow(row)
}

const clearSort = () => {
  tableRef.value?.clearSort()
}

const clearFilter = (columnKey) => {
  tableRef.value?.clearFilter(columnKey)
}

const doLayout = () => {
  nextTick(() => {
    tableRef.value?.doLayout()
  })
}

const scrollTo = (options) => {
  tableRef.value?.scrollTo(options)
}

const setScrollTop = (top) => {
  tableRef.value?.setScrollTop(top)
}

const setScrollLeft = (left) => {
  tableRef.value?.setScrollLeft(left)
}

// 监听器
watch(() => props.pageSize, (newSize) => {
  currentPageSize.value = newSize
})

watch(() => props.currentPage, (newPage) => {
  currentPage.value = newPage
})

watch(currentDensity, (newDensity) => {
  localStorage.setItem('table-density', newDensity)
  doLayout()
})

// 生命周期
onMounted(() => {
  // 恢复列显示设置
  const savedColumns = localStorage.getItem('table-visible-columns')
  if (savedColumns) {
    try {
      visibleColumns.value = JSON.parse(savedColumns)
    } catch (error) {
      visibleColumns.value = settableColumns.value.map(col => col.prop)
    }
  } else {
    visibleColumns.value = settableColumns.value.map(col => col.prop)
  }
})

// 暴露方法
defineExpose({
  clearSelection,
  toggleRowSelection,
  toggleAllSelection,
  toggleRowExpansion,
  setCurrentRow,
  clearSort,
  clearFilter,
  doLayout,
  scrollTo,
  setScrollTop,
  setScrollLeft,
  refresh: handleRefresh,
  getSelectedRows: () => selectedRows.value,
  getCurrentPage: () => currentPage.value,
  getTotalPages: () => totalPages.value
})
</script>

<style lang="scss" scoped>
@import '@/assets/styles/design-tokens.scss';
@import '@/assets/styles/animations-enhanced.scss';

.premium-table {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  overflow: hidden;
  transition: all var(--transition-base);
  
  &:hover {
    box-shadow: var(--shadow-md);
  }
  
  // 表格工具栏
  .table-toolbar {
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-light);
    padding: var(--spacing-4) var(--spacing-6);
    
    .toolbar-content {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: var(--spacing-4);
      flex-wrap: wrap;
      
      .toolbar-left {
        display: flex;
        align-items: center;
        gap: var(--spacing-4);
        flex: 1;
        min-width: 0;
        
        .table-title {
          .title-text {
            font-size: var(--text-lg);
            font-weight: var(--font-semibold);
            color: var(--text-primary);
            margin: 0 0 var(--spacing-1) 0;
          }
          
          .title-subtitle {
            font-size: var(--text-sm);
            color: var(--text-tertiary);
            margin: 0;
          }
        }
        
        .search-container {
          .table-search {
            width: 320px;
            
            .el-input__wrapper {
              border-radius: var(--radius-lg);
              transition: all var(--transition-fast);
              
              &:hover {
                border-color: var(--border-default);
              }
              
              &.is-focus {
                border-color: var(--primary-500);
                box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
              }
            }
            
            .el-input-group__append {
              .el-button {
                border-radius: 0 var(--radius-lg) var(--radius-lg) 0;
                
                &.is-active {
                  background: var(--primary-500);
                  color: white;
                  border-color: var(--primary-500);
                }
              }
            }
          }
        }
      }
      
      .toolbar-right {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        
        .table-config {
          display: flex;
          align-items: center;
          gap: var(--spacing-2);
          
          .config-btn {
            width: 32px;
            height: 32px;
            border-radius: var(--radius-md);
            color: var(--text-secondary);
            border: 1px solid var(--border-light);
            transition: all var(--transition-fast);
            
            &:hover {
              color: var(--primary-500);
              border-color: var(--primary-500);
              background: var(--primary-50);
            }
          }
        }
      }
    }
  }
  
  // 高级筛选面板
  .advanced-filter {
    background: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-light);
    
    .filter-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: var(--spacing-4) var(--spacing-6);
      border-bottom: 1px solid var(--border-light);
      
      h4 {
        font-size: var(--text-base);
        font-weight: var(--font-semibold);
        color: var(--text-primary);
        margin: 0;
      }
      
      .filter-actions {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
      }
    }
    
    .filter-content {
      padding: var(--spacing-4) var(--spacing-6);
      
      .number-range {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        
        .range-separator {
          color: var(--text-tertiary);
          font-weight: var(--font-medium);
        }
        
        .el-input-number {
          flex: 1;
        }
      }
    }
  }
  
  // 批量操作栏
  .batch-actions {
    background: var(--primary-50);
    border-bottom: 1px solid var(--primary-200);
    padding: var(--spacing-3) var(--spacing-6);
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .batch-info {
      display: flex;
      align-items: center;
      gap: var(--spacing-3);
      
      .selected-count {
        font-size: var(--text-sm);
        font-weight: var(--font-medium);
        color: var(--primary-700);
      }
    }
    
    .batch-operations {
      display: flex;
      align-items: center;
      gap: var(--spacing-2);
    }
  }
  
  // 表格容器
  .table-container {
    position: relative;
    
    &.density-compact {
      :deep(.el-table) {
        font-size: var(--text-sm);
        
        .el-table__header th,
        .el-table__body td {
          padding: var(--spacing-2) var(--spacing-3);
        }
      }
    }
    
    &.density-comfortable {
      :deep(.el-table) {
        .el-table__header th,
        .el-table__body td {
          padding: var(--spacing-5) var(--spacing-4);
        }
      }
    }
    
    // 表格主体样式
    .premium-el-table {
      border: none;
      
      :deep(.el-table__header) {
        th {
          background: var(--bg-secondary);
          border-bottom: 1px solid var(--border-light);
          font-weight: var(--font-semibold);
          color: var(--text-primary);
          font-size: var(--text-sm);
          position: relative;
          
          &:hover {
            background: var(--bg-hover);
          }
          
          .cell {
            padding: 0 var(--spacing-4);
            line-height: 1.5;
          }
        }
      }
      
      :deep(.el-table__body) {
        .table-row {
          transition: all var(--transition-fast);
          border-bottom: 1px solid var(--border-light);
          
          &.hover-effect:hover {
            background: var(--bg-hover);
            transform: translateY(-1px);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
          }
          
          td {
            border-bottom: 1px solid var(--border-light);
            
            .cell {
              padding: 0 var(--spacing-4);
              line-height: 1.5;
            }
          }
          
          // 表格操作列
          .table-actions-column {
            .table-actions {
              display: flex;
              align-items: center;
              justify-content: center;
              gap: var(--spacing-1);
              flex-wrap: wrap;
              
              .el-button {
                transition: all var(--transition-fast);
                
                &:hover {
                  transform: translateY(-1px);
                }
              }
            }
          }
        }
      }
      
      // 空状态
      :deep(.el-table__empty-block) {
        .table-empty {
          padding: var(--spacing-12) var(--spacing-6);
          
          .empty-content {
            text-align: center;
            
            .empty-icon {
              font-size: 64px;
              color: var(--text-tertiary);
              margin-bottom: var(--spacing-4);
              opacity: 0.5;
            }
            
            .empty-text {
              font-size: var(--text-base);
              color: var(--text-secondary);
              margin: 0 0 var(--spacing-4) 0;
            }
          }
        }
      }
      
      // 加载状态
      :deep(.el-loading-spinner) {
        .circular {
          width: 42px;
          height: 42px;
          color: var(--primary-500);
        }
      }
    }
  }
  
  // 分页器
  .table-pagination {
    padding: var(--spacing-4) var(--spacing-6);
    border-top: 1px solid var(--border-light);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--spacing-4);
    flex-wrap: wrap;
    
    .pagination-info {
      .total-info {
        font-size: var(--text-sm);
        color: var(--text-secondary);
      }
    }
    
    .pagination-component {
      :deep(.el-pager) {
        li {
          margin: 0 var(--spacing-1);
          border-radius: var(--radius-md);
          transition: all var(--transition-fast);
          
          &:hover {
            color: var(--primary-600);
            background: var(--primary-50);
          }
          
          &.is-active {
            background: var(--primary-500);
            color: white;
            box-shadow: var(--shadow-sm);
          }
        }
      }
      
      :deep(.btn-prev),
      :deep(.btn-next) {
        border-radius: var(--radius-md);
        transition: all var(--transition-fast);
        
        &:hover {
          color: var(--primary-600);
          background: var(--primary-50);
        }
      }
      
      :deep(.el-select) {
        .el-input__wrapper {
          border-radius: var(--radius-md);
        }
      }
    }
  }
  
  // 表格底部
  .table-footer {
    padding: var(--spacing-3) var(--spacing-6);
    background: var(--bg-secondary);
    border-top: 1px solid var(--border-light);
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .footer-stats {
      display: flex;
      align-items: center;
      gap: var(--spacing-4);
      
      .stats-item {
        font-size: var(--text-sm);
        color: var(--text-secondary);
        
        &:not(:last-child)::after {
          content: '|';
          margin-left: var(--spacing-4);
          color: var(--border-default);
        }
      }
    }
    
    .footer-actions {
      display: flex;
      align-items: center;
      gap: var(--spacing-2);
    }
  }
}

// 列设置弹窗
.column-setting {
  .setting-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--spacing-4);
    padding-bottom: var(--spacing-3);
    border-bottom: 1px solid var(--border-light);
    
    h4 {
      font-size: var(--text-base);
      font-weight: var(--font-semibold);
      color: var(--text-primary);
      margin: 0;
    }
  }
  
  .setting-content {
    .setting-search {
      margin-bottom: var(--spacing-3);
    }
    
    .column-list {
      max-height: 240px;
      overflow-y: auto;
      
      .column-item {
        padding: var(--spacing-2) 0;
        border-bottom: 1px solid var(--border-light);
        
        &:last-child {
          border-bottom: none;
        }
        
        .el-checkbox {
          width: 100%;
          
          .el-checkbox__label {
            display: flex;
            align-items: center;
            justify-content: space-between;
            width: 100%;
            
            .column-label {
              flex: 1;
              margin-right: var(--spacing-2);
            }
          }
        }
      }
    }
  }
}

// 下拉菜单活跃状态
:deep(.el-dropdown-menu__item.is-active) {
  color: var(--primary-600);
  background: var(--primary-50);
  
  .el-icon {
    color: var(--primary-600);
  }
}

// 动画效果
.filter-slide-enter-active,
.filter-slide-leave-active {
  transition: all var(--transition-base);
  overflow: hidden;
}

.filter-slide-enter-from,
.filter-slide-leave-to {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.filter-slide-enter-to,
.filter-slide-leave-from {
  max-height: 200px;
  opacity: 1;
}

.batch-slide-enter-active,
.batch-slide-leave-active {
  transition: all var(--transition-base);
}

.batch-slide-enter-from,
.batch-slide-leave-to {
  opacity: 0;
  transform: translateY(-100%);
}

// 响应式设计
@media (max-width: 768px) {
  .premium-table {
    border-radius: var(--radius-lg);
    
    .table-toolbar {
      padding: var(--spacing-3) var(--spacing-4);
      
      .toolbar-content {
        flex-direction: column;
        align-items: stretch;
        gap: var(--spacing-3);
        
        .toolbar-left {
          .search-container {
            .table-search {
              width: 100%;
            }
          }
        }
        
        .toolbar-right {
          justify-content: space-between;
        }
      }
    }
    
    .advanced-filter {
      .filter-header {
        padding: var(--spacing-3) var(--spacing-4);
        flex-direction: column;
        align-items: stretch;
        gap: var(--spacing-3);
      }
      
      .filter-content {
        padding: var(--spacing-3) var(--spacing-4);
      }
    }
    
    .table-container {
      overflow-x: auto;
      
      .premium-el-table {
        min-width: 600px;
      }
    }
    
    .table-pagination {
      padding: var(--spacing-3) var(--spacing-4);
      flex-direction: column;
      align-items: stretch;
      gap: var(--spacing-3);
      
      .pagination-info {
        text-align: center;
      }
      
      .pagination-component {
        :deep(.el-pagination) {
          justify-content: center;
        }
      }
    }
  }
}

// 自定义滚动条
.custom-scrollbar {
  &::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }
  
  &::-webkit-scrollbar-track {
    background: var(--bg-secondary);
    border-radius: var(--radius-full);
  }
  
  &::-webkit-scrollbar-thumb {
    background: var(--border-default);
    border-radius: var(--radius-full);
    
    &:hover {
      background: var(--border-strong);
    }
  }
}
</style>