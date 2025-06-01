<template>
  <div class="enhanced-form">
    <el-form
      ref="formRef"
      :model="modelValue"
      :rules="rules"
      :label-width="labelWidth"
      :label-position="labelPosition"
      :size="size"
      :disabled="disabled"
      :validate-on-rule-change="false"
      @submit.prevent
    >
      <div
        v-for="(section, sectionIndex) in formSections"
        :key="sectionIndex"
        class="form-section"
      >
        <!-- 分组标题 -->
        <div v-if="section.title" class="section-header">
          <h3 class="section-title">
            <component v-if="section.icon" :is="section.icon" />
            {{ section.title }}
          </h3>
          <p v-if="section.description" class="section-description">
            {{ section.description }}
          </p>
        </div>
        
        <!-- 表单项 -->
        <el-row :gutter="20">
          <el-col
            v-for="field in section.fields"
            :key="field.prop"
            :span="field.span || 24"
            :xs="field.xs || 24"
            :sm="field.sm || field.span || 24"
            :md="field.md || field.span || 12"
            :lg="field.lg || field.span || 12"
            :xl="field.xl || field.span || 12"
          >
            <el-form-item
              :prop="field.prop"
              :label="field.label"
              :label-width="field.labelWidth"
              :required="field.required"
              :rules="field.rules"
              :class="[`form-item-${field.type}`, field.className]"
            >
              <!-- 输入提示 -->
              <template v-if="field.tip" #label>
                <span class="form-item-label">
                  {{ field.label }}
                  <el-tooltip :content="field.tip" placement="top">
                    <el-icon class="tip-icon"><QuestionFilled /></el-icon>
                  </el-tooltip>
                </span>
              </template>
              
              <!-- 文本输入 -->
              <el-input
                v-if="field.type === 'text' || field.type === 'email' || field.type === 'url'"
                v-model="modelValue[field.prop]"
                :type="field.type"
                :placeholder="field.placeholder"
                :clearable="field.clearable !== false"
                :show-word-limit="field.showWordLimit"
                :maxlength="field.maxlength"
                :prefix-icon="field.prefixIcon"
                :suffix-icon="field.suffixIcon"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                @blur="handleFieldBlur(field)"
                @focus="handleFieldFocus(field)"
              >
                <template v-if="field.prepend" #prepend>
                  {{ field.prepend }}
                </template>
                <template v-if="field.append" #append>
                  {{ field.append }}
                </template>
              </el-input>
              
              <!-- 密码输入 -->
              <el-input
                v-else-if="field.type === 'password'"
                v-model="modelValue[field.prop]"
                type="password"
                :placeholder="field.placeholder"
                :show-password="true"
                :clearable="field.clearable !== false"
                :prefix-icon="field.prefixIcon || Lock"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              />
              
              <!-- 文本域 -->
              <el-input
                v-else-if="field.type === 'textarea'"
                v-model="modelValue[field.prop]"
                type="textarea"
                :placeholder="field.placeholder"
                :rows="field.rows || 3"
                :autosize="field.autosize"
                :show-word-limit="field.showWordLimit"
                :maxlength="field.maxlength"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              />
              
              <!-- 数字输入 -->
              <el-input-number
                v-else-if="field.type === 'number'"
                v-model="modelValue[field.prop]"
                :min="field.min"
                :max="field.max"
                :step="field.step || 1"
                :precision="field.precision"
                :controls="field.controls !== false"
                :controls-position="field.controlsPosition"
                :placeholder="field.placeholder"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              />
              
              <!-- 选择器 -->
              <el-select
                v-else-if="field.type === 'select'"
                v-model="modelValue[field.prop]"
                :placeholder="field.placeholder"
                :clearable="field.clearable !== false"
                :multiple="field.multiple"
                :filterable="field.filterable"
                :remote="field.remote"
                :remote-method="field.remoteMethod"
                :loading="field.loading"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              >
                <el-option
                  v-for="option in field.options"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                  :disabled="option.disabled"
                >
                  <span v-if="field.optionRender" v-html="field.optionRender(option)"></span>
                </el-option>
              </el-select>
              
              <!-- 级联选择 -->
              <el-cascader
                v-else-if="field.type === 'cascader'"
                v-model="modelValue[field.prop]"
                :options="field.options"
                :props="field.cascaderProps"
                :placeholder="field.placeholder"
                :clearable="field.clearable !== false"
                :filterable="field.filterable"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              />
              
              <!-- 日期选择 -->
              <el-date-picker
                v-else-if="field.type === 'date'"
                v-model="modelValue[field.prop]"
                type="date"
                :placeholder="field.placeholder"
                :format="field.format || 'YYYY-MM-DD'"
                :value-format="field.valueFormat || 'YYYY-MM-DD'"
                :clearable="field.clearable !== false"
                :disabled-date="field.disabledDate"
                :shortcuts="field.shortcuts"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              />
              
              <!-- 日期时间选择 -->
              <el-date-picker
                v-else-if="field.type === 'datetime'"
                v-model="modelValue[field.prop]"
                type="datetime"
                :placeholder="field.placeholder"
                :format="field.format || 'YYYY-MM-DD HH:mm:ss'"
                :value-format="field.valueFormat || 'YYYY-MM-DD HH:mm:ss'"
                :clearable="field.clearable !== false"
                :disabled-date="field.disabledDate"
                :disabled-hours="field.disabledHours"
                :disabled-minutes="field.disabledMinutes"
                :disabled-seconds="field.disabledSeconds"
                :shortcuts="field.shortcuts"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              />
              
              <!-- 日期范围 -->
              <el-date-picker
                v-else-if="field.type === 'daterange'"
                v-model="modelValue[field.prop]"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                :format="field.format || 'YYYY-MM-DD'"
                :value-format="field.valueFormat || 'YYYY-MM-DD'"
                :clearable="field.clearable !== false"
                :disabled-date="field.disabledDate"
                :shortcuts="field.shortcuts"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              />
              
              <!-- 时间选择 -->
              <el-time-picker
                v-else-if="field.type === 'time'"
                v-model="modelValue[field.prop]"
                :placeholder="field.placeholder"
                :format="field.format || 'HH:mm:ss'"
                :value-format="field.valueFormat || 'HH:mm:ss'"
                :clearable="field.clearable !== false"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
                class="w-full"
              />
              
              <!-- 开关 -->
              <el-switch
                v-else-if="field.type === 'switch'"
                v-model="modelValue[field.prop]"
                :active-text="field.activeText"
                :inactive-text="field.inactiveText"
                :active-value="field.activeValue"
                :inactive-value="field.inactiveValue"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              />
              
              <!-- 单选框组 -->
              <el-radio-group
                v-else-if="field.type === 'radio'"
                v-model="modelValue[field.prop]"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              >
                <el-radio
                  v-for="option in field.options"
                  :key="option.value"
                  :value="option.value"
                  :disabled="option.disabled"
                >
                  {{ option.label }}
                </el-radio>
              </el-radio-group>
              
              <!-- 多选框组 -->
              <el-checkbox-group
                v-else-if="field.type === 'checkbox'"
                v-model="modelValue[field.prop]"
                :min="field.min"
                :max="field.max"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              >
                <el-checkbox
                  v-for="option in field.options"
                  :key="option.value"
                  :value="option.value"
                  :disabled="option.disabled"
                >
                  {{ option.label }}
                </el-checkbox>
              </el-checkbox-group>
              
              <!-- 评分 -->
              <el-rate
                v-else-if="field.type === 'rate'"
                v-model="modelValue[field.prop]"
                :max="field.max || 5"
                :allow-half="field.allowHalf"
                :show-text="field.showText"
                :texts="field.texts"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              />
              
              <!-- 滑块 -->
              <el-slider
                v-else-if="field.type === 'slider'"
                v-model="modelValue[field.prop]"
                :min="field.min || 0"
                :max="field.max || 100"
                :step="field.step || 1"
                :show-input="field.showInput"
                :range="field.range"
                :marks="field.marks"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              />
              
              <!-- 颜色选择 -->
              <el-color-picker
                v-else-if="field.type === 'color'"
                v-model="modelValue[field.prop]"
                :show-alpha="field.showAlpha"
                :predefine="field.predefine"
                :disabled="field.disabled"
                @change="handleFieldChange(field, $event)"
              />
              
              <!-- 文件上传 -->
              <el-upload
                v-else-if="field.type === 'upload'"
                v-model:file-list="modelValue[field.prop]"
                :action="field.action"
                :multiple="field.multiple"
                :limit="field.limit"
                :accept="field.accept"
                :list-type="field.listType || 'text'"
                :auto-upload="field.autoUpload !== false"
                :disabled="field.disabled"
                :on-success="(res, file) => handleUploadSuccess(field, res, file)"
                :on-error="(err, file) => handleUploadError(field, err, file)"
                :on-remove="(file) => handleUploadRemove(field, file)"
                :before-upload="field.beforeUpload"
                class="w-full"
              >
                <el-button v-if="field.listType !== 'picture-card'" type="primary">
                  <Upload />
                  {{ field.buttonText || '点击上传' }}
                </el-button>
                <el-icon v-else class="upload-icon"><Plus /></el-icon>
                <template v-if="field.uploadTip" #tip>
                  <div class="upload-tip">{{ field.uploadTip }}</div>
                </template>
              </el-upload>
              
              <!-- 自定义插槽 -->
              <slot
                v-else-if="field.type === 'custom'"
                :name="field.slotName"
                :field="field"
                :value="modelValue[field.prop]"
                :model="modelValue"
              />
              
              <!-- 字段说明 -->
              <div v-if="field.description" class="field-description">
                {{ field.description }}
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </div>
      
      <!-- 表单操作按钮 -->
      <div class="form-actions" v-if="showActions">
        <slot name="actions" :validate="validate" :reset="resetForm">
          <el-button @click="handleCancel">
            {{ cancelText }}
          </el-button>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            {{ submitText }}
          </el-button>
        </slot>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  QuestionFilled, Lock, Upload, Plus
} from '@element-plus/icons-vue'

// Props
const props = defineProps({
  // 表单数据
  modelValue: {
    type: Object,
    required: true
  },
  // 表单配置
  formSections: {
    type: Array,
    required: true
  },
  // 验证规则
  rules: {
    type: Object,
    default: () => ({})
  },
  // 表单属性
  labelWidth: {
    type: String,
    default: '120px'
  },
  labelPosition: {
    type: String,
    default: 'right'
  },
  size: {
    type: String,
    default: 'default'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  // 操作按钮
  showActions: {
    type: Boolean,
    default: true
  },
  submitText: {
    type: String,
    default: '提交'
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits([
  'update:modelValue',
  'submit',
  'cancel',
  'field-change',
  'field-focus',
  'field-blur',
  'upload-success',
  'upload-error',
  'upload-remove'
])

// Refs
const formRef = ref()

// 方法
const validate = async () => {
  return await formRef.value.validate()
}

const validateField = async (prop) => {
  return await formRef.value.validateField(prop)
}

const resetFields = () => {
  formRef.value.resetFields()
}

const clearValidate = (props) => {
  formRef.value.clearValidate(props)
}

const scrollToField = (prop) => {
  formRef.value.scrollToField(prop)
}

const handleSubmit = async () => {
  try {
    const valid = await validate()
    if (valid) {
      emit('submit', props.modelValue)
    }
  } catch (error) {
    ElMessage.warning('请检查表单填写是否正确')
    console.error('Form validation error:', error)
  }
}

const handleCancel = () => {
  emit('cancel')
}

const resetForm = () => {
  resetFields()
}

const handleFieldChange = (field, value) => {
  emit('field-change', { field, value })
  
  // 如果字段有联动配置，触发联动
  if (field.onChange) {
    field.onChange(value, props.modelValue)
  }
}

const handleFieldFocus = (field) => {
  emit('field-focus', field)
}

const handleFieldBlur = (field) => {
  emit('field-blur', field)
  
  // 失焦时验证单个字段
  if (field.validateOnBlur !== false) {
    validateField(field.prop)
  }
}

const handleUploadSuccess = (field, response, file) => {
  emit('upload-success', { field, response, file })
  
  if (field.onUploadSuccess) {
    field.onUploadSuccess(response, file, props.modelValue)
  }
}

const handleUploadError = (field, error, file) => {
  emit('upload-error', { field, error, file })
  ElMessage.error('文件上传失败')
  
  if (field.onUploadError) {
    field.onUploadError(error, file)
  }
}

const handleUploadRemove = (field, file) => {
  emit('upload-remove', { field, file })
  
  if (field.onUploadRemove) {
    field.onUploadRemove(file, props.modelValue)
  }
}

// 暴露方法
defineExpose({
  validate,
  validateField,
  resetFields,
  clearValidate,
  scrollToField
})
</script>

<style lang="scss" scoped>
.enhanced-form {
  // 表单分组
  .form-section {
    margin-bottom: var(--spacing-8);
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .section-header {
      margin-bottom: var(--spacing-6);
      padding-bottom: var(--spacing-4);
      border-bottom: 1px solid var(--border-light);
      
      .section-title {
        font-size: 1.125rem;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 var(--spacing-2);
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
        
        :deep(.el-icon) {
          font-size: 1.25rem;
          color: var(--primary-500);
        }
      }
      
      .section-description {
        font-size: 0.875rem;
        color: var(--text-secondary);
        margin: 0;
      }
    }
  }
  
  // 表单项
  :deep(.el-form-item) {
    margin-bottom: var(--spacing-6);
    
    .el-form-item__label {
      font-weight: 500;
      color: var(--text-primary);
      
      .form-item-label {
        display: flex;
        align-items: center;
        gap: var(--spacing-1);
        
        .tip-icon {
          font-size: 0.875rem;
          color: var(--text-tertiary);
          cursor: help;
          
          &:hover {
            color: var(--primary-500);
          }
        }
      }
    }
    
    .el-form-item__content {
      .field-description {
        margin-top: var(--spacing-2);
        font-size: 0.875rem;
        color: var(--text-secondary);
        line-height: 1.5;
      }
      
      .upload-tip {
        margin-top: var(--spacing-2);
        font-size: 0.875rem;
        color: var(--text-tertiary);
      }
    }
    
    // 错误提示优化
    .el-form-item__error {
      padding-top: var(--spacing-1);
      font-size: 0.875rem;
    }
  }
  
  // 特定类型的表单项样式
  .form-item-switch {
    :deep(.el-form-item__content) {
      line-height: 32px;
    }
  }
  
  .form-item-radio,
  .form-item-checkbox {
    :deep(.el-radio),
    :deep(.el-checkbox) {
      margin-right: var(--spacing-6);
      margin-bottom: var(--spacing-2);
    }
  }
  
  .form-item-upload {
    :deep(.el-upload) {
      .upload-icon {
        font-size: 2rem;
        color: var(--text-tertiary);
      }
      
      &.el-upload--picture-card {
        background: var(--bg-secondary);
        border: 2px dashed var(--border-default);
        
        &:hover {
          border-color: var(--primary-500);
          
          .upload-icon {
            color: var(--primary-500);
          }
        }
      }
    }
  }
  
  // 表单操作按钮
  .form-actions {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: var(--spacing-3);
    margin-top: var(--spacing-8);
    padding-top: var(--spacing-6);
    border-top: 1px solid var(--border-light);
  }
  
  // 工具类
  .w-full {
    width: 100%;
  }
}
</style>