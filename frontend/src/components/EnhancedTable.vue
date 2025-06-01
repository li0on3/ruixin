<template>
  <div class="enhanced-table">
    <!-- 表格工具栏 -->
    <div class="table-toolbar" v-if="showToolbar">
      <div class="toolbar-left">
        <slot name="toolbar-left">
          <!-- 搜索框 -->
          <el-input
            v-if="searchable"
            v-model="searchText"
            :placeholder="searchPlaceholder"
            prefix-icon="Search"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            class="table-search"
          />
        </slot>
      </div>
      
      <div class="toolbar-right">
        <slot name="toolbar-right">
          <!-- 刷新按钮 -->
          <el-button
            v-if="refreshable"
            :icon="Refresh"
            circle
            @click="handleRefresh"
            title="刷新数据"
          />
          
          <!-- 列设置 -->
          <el-popover
            v-if="columnSettable"
            placement="bottom"
            :width="200"
            trigger="click"
          >
            <template #reference>
              <el-button :icon="Setting" circle title="列设置" />
            </template>
            <div class="column-setting">
              <div class="column-setting-header">
                <span>列显示设置</span>
                <el-button text size="small" @click="resetColumns">重置</el-button>
              </div>
              <el-checkbox-group v-model="visibleColumns" @change="handleColumnChange">
                <el-checkbox
                  v-for="col in settableColumns"
                  :key="col.prop"
                  :value="col.prop"
                  :disabled="!!col.fixed"
                >
                  {{ col.label }}
                </el-checkbox>
              </el-checkbox-group>
            </div>
          </el-popover>
          
          <!-- 密度设置 -->
          <el-dropdown v-if="densitySettable" trigger="click">
            <el-button :icon="Grid" circle title="密度" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="density in densityOptions"
                  :key="density.value"
                  @click="currentDensity = density.value"
                  :class="{ active: currentDensity === density.value }"
                >
                  <el-icon v-if="currentDensity === density.value"><Check /></el-icon>
                  {{ density.label }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <!-- 导出按钮 -->
          <el-dropdown v-if="exportable" trigger="click">
            <el-button :icon="Download" circle title="导出" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleExport('excel')">
                  <Document />
                  导出 Excel
                </el-dropdown-item>
                <el-dropdown-item @click="handleExport('csv')">
                  <Document />
                  导出 CSV
                </el-dropdown-item>
                <el-dropdown-item @click="handleExport('pdf')">
                  <Document />
                  导出 PDF
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </slot>
      </div>
    </div>
    
    <!-- 高级筛选面板 -->
    <el-collapse-transition>
      <div class="advanced-filter" v-if="showAdvancedFilter">
        <div class="filter-content">
          <el-form :model="filterForm" label-width="100px">
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
                  />
                  
                  <!-- 数字范围 -->
                  <div v-else-if="filter.type === 'numberRange'" class="number-range">
                    <el-input-number
                      v-model="filterForm[filter.prop + '_min']"
                      :placeholder="filter.minPlaceholder || '最小值'"
                      :controls="false"
                    />
                    <span class="range-separator">-</span>
                    <el-input-number
                      v-model="filterForm[filter.prop + '_max']"
                      :placeholder="filter.maxPlaceholder || '最大值'"
                      :controls="false"
                    />
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
            
            <div class="filter-actions">
              <el-button type="primary" @click="handleFilter">筛选</el-button>
              <el-button @click="handleResetFilter">重置</el-button>
              <el-button text @click="showAdvancedFilter = false">收起</el-button>
            </div>
          </el-form>
        </div>
      </div>
    </el-collapse-transition>
    
    <!-- 表格主体 -->
    <div class="table-wrapper" :class="[`density-${currentDensity}`]">
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
        :row-class-name="rowClassName"
        :cell-class-name="cellClassName"
        :header-cell-class-name="headerCellClassName"
        :default-sort="defaultSort"
        class="enhanced-el-table"
      >
        <!-- 选择列 -->
        <el-table-column
          v-if="selectable"
          type="selection"
          width="55"
          fixed="left"
        />
        
        <!-- 序号列 -->
        <el-table-column
          v-if="showIndex"
          type="index"
          label="序号"
          width="60"
          fixed="left"
          :index="indexMethod"
        />
        
        <!-- 数据列 -->
        <template v-for="column in displayColumns" :key="column.prop">
          <el-table-column
            v-if="!column.hide && isColumnVisible(column.prop)"
            :prop="column.prop"
            :label="column.label"
            :width="column.width"
            :min-width="column.minWidth"
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
            <template v-if="column.slot" #default="scope">
              <slot :name="column.slot" v-bind="scope" />
            </template>
            
            <template v-else-if="column.render" #default="scope">
              <component :is="column.render(scope.row, scope.$index)" />
            </template>
          </el-table-column>
        </template>
        
        <!-- 操作列 -->
        <el-table-column
          v-if="actions && actions.length > 0"
          label="操作"
          :width="actionWidth"
          :fixed="actionFixed"
          :align="actionAlign"
          class-name="table-actions"
        >
          <template #default="scope">
            <div class="action-buttons">
              <template v-for="(action, index) in getRowActions(scope.row, scope.$index)" :key="index">
                <el-button
                  v-if="!action.hide"
                  :type="action.type || 'text'"
                  :size="action.size || 'small'"
                  :icon="action.icon"
                  :loading="action.loading"
                  :disabled="action.disabled"
                  @click="handleAction(action, scope.row, scope.$index)"
                  :title="action.tooltip"
                >
                  {{ action.text }}
                </el-button>
              </template>
            </div>
          </template>
        </el-table-column>
        
        <!-- 空数据插槽 -->
        <template #empty>
          <div class="table-empty">
            <el-empty :description="emptyText" />
          </div>
        </template>
      </el-table>
    </div>
    
    <!-- 分页 -->
    <div class="table-pagination" v-if="pagination">
      <el-pagination
        :current-page="currentPage"
        :page-size="currentPageSize"
        :total="totalCount"
        :page-sizes="pageSizes"
        :layout="paginationLayout"
        :background="true"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Refresh, Setting, Grid, Download, Check, Document,
  Plus, Edit, Delete, View, Share, Filter
} from '@element-plus/icons-vue'

// Props
const props = defineProps({
  // 表格数据
  data: {
    type: Array,
    default: () => []
  },
  // 列配置
  columns: {
    type: Array,
    required: true
  },
  // 加载状态
  loading: {
    type: Boolean,
    default: false
  },
  // 表格高度
  height: {
    type: [String, Number]
  },
  maxHeight: {
    type: [String, Number]
  },
  // 样式相关
  stripe: {
    type: Boolean,
    default: true
  },
  border: {
    type: Boolean,
    default: false
  },
  // 功能开关
  searchable: {
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
  // 工具栏
  showToolbar: {
    type: Boolean,
    default: true
  },
  // 选择
  selectable: {
    type: Boolean,
    default: false
  },
  // 序号
  showIndex: {
    type: Boolean,
    default: false
  },
  // 汇总
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
    default: 10
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
  // 操作列
  actions: {
    type: Array,
    default: () => []
  },
  actionWidth: {
    type: [String, Number],
    default: 200
  },
  actionFixed: {
    type: [Boolean, String],
    default: 'right'
  },
  actionAlign: {
    type: String,
    default: 'center'
  },
  // 高级筛选
  filterConfig: {
    type: Array,
    default: () => []
  },
  // 其他
  emptyText: {
    type: String,
    default: '暂无数据'
  },
  searchPlaceholder: {
    type: String,
    default: '请输入搜索内容'
  },
  rowClassName: Function,
  cellClassName: Function,
  headerCellClassName: Function,
  defaultSort: Object
})

// Emits
const emit = defineEmits([
  'search',
  'refresh',
  'filter',
  'sort-change',
  'selection-change',
  'page-change',
  'size-change',
  'export',
  'action-click'
])

// Refs
const tableRef = ref()
const searchText = ref('')
const currentPage = ref(props.currentPage)
const currentPageSize = ref(props.pageSize)
const currentDensity = ref('default')
const showAdvancedFilter = ref(false)
const filterForm = ref({})
const visibleColumns = ref([])
const selectedRows = ref([])

// 密度选项
const densityOptions = [
  { label: '紧凑', value: 'compact' },
  { label: '默认', value: 'default' },
  { label: '宽松', value: 'comfortable' }
]

// 计算属性
const totalCount = computed(() => props.total || props.data.length)

const displayData = computed(() => {
  if (!props.pagination) {
    return props.data
  }
  const start = (currentPage.value - 1) * currentPageSize.value
  const end = start + currentPageSize.value
  return props.data.slice(start, end)
})

const settableColumns = computed(() => {
  return props.columns.filter(col => !col.type && col.prop)
})

const displayColumns = computed(() => {
  return props.columns
})

// 方法
const handleSearch = () => {
  emit('search', searchText.value)
}

const handleRefresh = () => {
  emit('refresh')
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


const handleExport = (type) => {
  emit('export', {
    type,
    data: selectedRows.value.length > 0 ? selectedRows.value : props.data,
    columns: props.columns.filter(col => isColumnVisible(col.prop))
  })
}

const handleAction = (action, row, index) => {
  if (action.handler) {
    action.handler(row, index)
  }
  emit('action-click', { action, row, index })
}

const handleColumnChange = (checkedColumns) => {
  localStorage.setItem('table-visible-columns', JSON.stringify(checkedColumns))
}

const resetColumns = () => {
  visibleColumns.value = settableColumns.value.map(col => col.prop)
  handleColumnChange(visibleColumns.value)
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
    return action
  })
}

const indexMethod = (index) => {
  return (currentPage.value - 1) * props.pageSize + index + 1
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
  tableRef.value?.doLayout()
}

// 分页事件处理
const handleCurrentChange = (page) => {
  currentPage.value = page
  emit('page-change', page)
}

const handleSizeChange = (size) => {
  currentPageSize.value = size
  currentPage.value = 1 // 重置到第一页
  emit('size-change', size)
  emit('page-change', 1)
}

// 监听props变化
watch(() => props.pageSize, (newSize) => {
  currentPageSize.value = newSize
})

// 初始化
onMounted(() => {
  // 恢复列显示设置
  const savedColumns = localStorage.getItem('table-visible-columns')
  if (savedColumns) {
    visibleColumns.value = JSON.parse(savedColumns)
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
  doLayout
})
</script>

<style lang="scss" scoped>
.enhanced-table {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  overflow: hidden;
  
  // 工具栏
  .table-toolbar {
    padding: var(--spacing-4) var(--spacing-6);
    border-bottom: 1px solid var(--border-light);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--spacing-4);
    flex-wrap: wrap;
    
    .toolbar-left {
      display: flex;
      align-items: center;
      gap: var(--spacing-3);
      flex: 1;
      
      .table-search {
        width: 300px;
      }
    }
    
    .toolbar-right {
      display: flex;
      align-items: center;
      gap: var(--spacing-2);
    }
  }
  
  // 高级筛选
  .advanced-filter {
    padding: var(--spacing-4) var(--spacing-6);
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-light);
    
    .filter-content {
      .filter-actions {
        display: flex;
        align-items: center;
        gap: var(--spacing-3);
        margin-top: var(--spacing-4);
      }
      
      .number-range {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        
        .range-separator {
          color: var(--text-tertiary);
        }
      }
    }
  }
  
  // 表格容器
  .table-wrapper {
    &.density-compact {
      :deep(.el-table) {
        font-size: 0.875rem;
        
        .el-table__header th {
          padding: var(--spacing-2) 0;
        }
        
        .el-table__body td {
          padding: var(--spacing-2) 0;
        }
      }
    }
    
    &.density-comfortable {
      :deep(.el-table) {
        .el-table__header th {
          padding: var(--spacing-5) 0;
        }
        
        .el-table__body td {
          padding: var(--spacing-5) 0;
        }
      }
    }
  }
  
  // 表格样式
  .enhanced-el-table {
    :deep(.el-table__header) {
      th {
        background: var(--bg-secondary);
        font-weight: 600;
        color: var(--text-primary);
      }
    }
    
    :deep(.el-table__body) {
      .table-actions {
        .action-buttons {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: var(--spacing-2);
        }
      }
    }
    
    :deep(.el-table__empty-block) {
      padding: var(--spacing-8) 0;
    }
  }
  
  // 分页
  .table-pagination {
    padding: var(--spacing-4) var(--spacing-6);
    border-top: 1px solid var(--border-light);
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  // 列设置弹窗
  .column-setting {
    .column-setting-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: var(--spacing-3);
      padding-bottom: var(--spacing-3);
      border-bottom: 1px solid var(--border-light);
      
      span {
        font-weight: 500;
        color: var(--text-primary);
      }
    }
    
    .el-checkbox-group {
      display: flex;
      flex-direction: column;
      gap: var(--spacing-2);
      
      .el-checkbox {
        margin-right: 0;
      }
    }
  }
}

// 下拉菜单激活状态
:deep(.el-dropdown-menu__item.active) {
  color: var(--primary-600);
  background: var(--primary-50);
}
</style>