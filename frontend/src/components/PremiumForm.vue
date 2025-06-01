<template>
  <div class="premium-form" :class="formClasses">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      :label-width="labelWidth"
      :label-position="labelPosition"
      :size="size"
      :disabled="disabled"
      :validate-on-rule-change="validateOnRuleChange"
      :scroll-to-error="scrollToError"
      @validate="handleValidate"
      class="form-container"
    >
      <!-- 表单头部 -->
      <div class="form-header" v-if="title || subtitle || $slots.header">
        <slot name="header">
          <div class="form-title-section" v-if="title || subtitle">
            <h3 class="form-title" v-if="title">{{ title }}</h3>
            <p class="form-subtitle" v-if="subtitle">{{ subtitle }}</p>
          </div>
        </slot>
      </div>
      
      <!-- 步骤指示器 -->
      <div class="form-steps" v-if="steps && steps.length > 0">
        <el-steps :active="currentStep" :align-center="stepsAlignCenter" :process-status="stepStatus">
          <el-step
            v-for="(step, index) in steps"
            :key="index"
            :title="step.title"
            :description="step.description"
            :icon="step.icon"
          />
        </el-steps>
      </div>
      
      <!-- 表单字段 -->
      <div class="form-content">
        <template v-for="(field, index) in filteredFields" :key="field.prop || index">
          <!-- 分组标题 -->
          <div class="form-group" v-if="field.type === 'group'">
            <div class="group-header">
              <h4 class="group-title">
                <el-icon v-if="field.icon"><component :is="field.icon" /></el-icon>
                {{ field.title }}
              </h4>
              <p class="group-description" v-if="field.description">{{ field.description }}</p>
            </div>
            <div class="group-content">
              <el-row :gutter="field.gutter || 20">
                <template v-for="subField in field.fields" :key="subField.prop">
                  <el-col :span="subField.span || 24">
                    <premium-form-item :field="subField" :form-data="formData" />
                  </el-col>
                </template>
              </el-row>
            </div>
          </div>
          
          <!-- 分隔线 -->
          <el-divider v-else-if="field.type === 'divider'" :content-position="field.position || 'center'">
            {{ field.title }}
          </el-divider>
          
          <!-- 自定义内容 -->
          <div class="form-custom" v-else-if="field.type === 'custom'">
            <slot :name="field.slot" :form-data="formData" :field="field" />
          </div>
          
          <!-- 网格布局 -->
          <el-row v-else-if="field.type === 'row'" :gutter="field.gutter || 20">
            <template v-for="colField in field.fields" :key="colField.prop">
              <el-col :span="colField.span || 12">
                <premium-form-item :field="colField" :form-data="formData" />
              </el-col>
            </template>
          </el-row>
          
          <!-- 普通表单项 -->
          <premium-form-item v-else :field="field" :form-data="formData" />
        </template>
      </div>
      
      <!-- 表单底部 -->
      <div class="form-footer" v-if="showFooter || $slots.footer">
        <slot name="footer" :form-data="formData" :validate="validate" :resetFields="resetFields">
          <div class="form-actions" v-if="showFooter">
            <!-- 步骤导航 -->
            <div class="step-navigation" v-if="steps && steps.length > 0">
              <el-button 
                v-if="currentStep > 0"
                @click="prevStep"
                :disabled="loading"
              >
                上一步
              </el-button>
              
              <el-button 
                v-if="currentStep < steps.length - 1"
                type="primary"
                @click="nextStep"
                :loading="stepLoading"
              >
                下一步
              </el-button>
              
              <el-button 
                v-if="currentStep === steps.length - 1"
                type="primary"
                @click="handleSubmit"
                :loading="loading"
              >
                {{ submitText }}
              </el-button>
            </div>
            
            <!-- 普通操作按钮 -->
            <div class="normal-actions" v-else>
              <el-button 
                v-if="showReset"
                @click="handleReset"
                :disabled="loading"
              >
                {{ resetText }}
              </el-button>
              
              <el-button 
                v-if="showCancel"
                @click="handleCancel"
                :disabled="loading"
              >
                {{ cancelText }}
              </el-button>
              
              <el-button 
                type="primary"
                @click="handleSubmit"
                :loading="loading"
                :disabled="submitDisabled"
              >
                {{ submitText }}
              </el-button>
            </div>
          </div>
        </slot>
      </div>
    </el-form>
    
    <!-- 保存提示 -->
    <transition name="save-tip">
      <div class="save-tip" v-if="showSaveTip">
        <el-icon><SuccessFilled /></el-icon>
        <span>{{ saveTipText }}</span>
      </div>
    </transition>
  </div>
</template>

<script>
// 表单项组件
const PremiumFormItem = {
  name: 'PremiumFormItem',
  props: {
    field: {
      type: Object,
      required: true
    },
    formData: {
      type: Object,
      required: true
    }
  },
  setup(props) {
    const renderFormItem = () => {
      const { field, formData } = props
      
      return (
        <el-form-item
          prop={field.prop}
          label={field.label}
          rules={field.rules}
          required={field.required}
          error={field.error}
          class={[
            'premium-form-item',
            field.class,
            {
              'is-required': field.required,
              'has-help': field.help,
              [`type-${field.type}`]: field.type
            }
          ]}
        >
          {field.help && (
            <template slot="label">
              <span class="form-label">
                {field.label}
                <el-tooltip content={field.help} placement="top">
                  <el-icon class="help-icon"><QuestionFilled /></el-icon>
                </el-tooltip>
              </span>
            </template>
          )}
          
          {renderControl(field, formData)}
          
          {field.description && (
            <div class="field-description">{field.description}</div>
          )}
        </el-form-item>
      )
    }
    
    const renderControl = (field, formData) => {
      const commonProps = {
        modelValue: formData[field.prop],
        'onUpdate:modelValue': (value) => {
          formData[field.prop] = value
          field.onChange && field.onChange(value, formData)
        },
        placeholder: field.placeholder,
        disabled: field.disabled,
        size: field.size,
        clearable: field.clearable !== false,
        ...field.props
      }
      
      switch (field.type) {
        case 'input':
          return (
            <el-input
              {...commonProps}
              type={field.inputType || 'text'}
              maxlength={field.maxlength}
              show-word-limit={field.showWordLimit}
              prefix-icon={field.prefixIcon}
              suffix-icon={field.suffixIcon}
            />
          )
          
        case 'textarea':
          return (
            <el-input
              {...commonProps}
              type="textarea"
              rows={field.rows || 4}
              maxlength={field.maxlength}
              show-word-limit={field.showWordLimit}
              autosize={field.autosize}
            />
          )
          
        case 'password':
          return (
            <el-input
              {...commonProps}
              type="password"
              show-password={field.showPassword !== false}
            />
          )
          
        case 'number':
          return (
            <el-input-number
              {...commonProps}
              min={field.min}
              max={field.max}
              step={field.step}
              precision={field.precision}
              controls-position={field.controlsPosition}
            />
          )
          
        case 'select':
          return (
            <el-select
              {...commonProps}
              multiple={field.multiple}
              filterable={field.filterable}
              remote={field.remote}
              remote-method={field.remoteMethod}
              loading={field.loading}
              style="width: 100%"
            >
              {field.options?.map(option => (
                <el-option
                  key={option.value}
                  label={option.label}
                  value={option.value}
                  disabled={option.disabled}
                />
              ))}
            </el-select>
          )
          
        case 'radio':
          return (
            <el-radio-group {...commonProps}>
              {field.options?.map(option => (
                <el-radio key={option.value} label={option.value}>
                  {option.label}
                </el-radio>
              ))}
            </el-radio-group>
          )
          
        case 'checkbox':
          if (field.options) {
            return (
              <el-checkbox-group {...commonProps}>
                {field.options?.map(option => (
                  <el-checkbox key={option.value} value={option.value}>
                    {option.label}
                  </el-checkbox>
                ))}
              </el-checkbox-group>
            )
          } else {
            return (
              <el-checkbox {...commonProps}>
                {field.text || field.label}
              </el-checkbox>
            )
          }
          
        case 'switch':
          return (
            <el-switch
              {...commonProps}
              active-text={field.activeText}
              inactive-text={field.inactiveText}
              active-value={field.activeValue || true}
              inactive-value={field.inactiveValue || false}
            />
          )
          
        case 'date':
          return (
            <el-date-picker
              {...commonProps}
              type={field.dateType || 'date'}
              format={field.format}
              value-format={field.valueFormat}
              shortcuts={field.shortcuts}
              style="width: 100%"
            />
          )
          
        case 'time':
          return (
            <el-time-picker
              {...commonProps}
              format={field.format}
              value-format={field.valueFormat}
              style="width: 100%"
            />
          )
          
        case 'upload':
          return (
            <el-upload
              {...field.uploadProps}
              file-list={formData[field.prop] || []}
              on-change={(file, fileList) => {
                formData[field.prop] = fileList
                field.onChange && field.onChange(fileList, formData)
              }}
            >
              <el-button type="primary">{field.uploadText || '选择文件'}</el-button>
              {field.tip && <div class="upload-tip">{field.tip}</div>}
            </el-upload>
          )
          
        case 'slider':
          return (
            <el-slider
              {...commonProps}
              min={field.min || 0}
              max={field.max || 100}
              step={field.step}
              show-stops={field.showStops}
              show-tooltip={field.showTooltip !== false}
              range={field.range}
            />
          )
          
        case 'rate':
          return (
            <el-rate
              {...commonProps}
              max={field.max || 5}
              allow-half={field.allowHalf}
              show-text={field.showText}
              texts={field.texts}
            />
          )
          
        case 'cascader':
          return (
            <el-cascader
              {...commonProps}
              options={field.options}
              props={field.cascaderProps}
              filterable={field.filterable}
              style="width: 100%"
            />
          )
          
        case 'custom':
          return field.render ? field.render(formData, field) : null
          
        default:
          return (
            <el-input {...commonProps} />
          )
      }
    }
    
    return renderFormItem
  }
}
</script>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { SuccessFilled, QuestionFilled } from '@element-plus/icons-vue'

// Props
const props = defineProps({
  // 表单数据
  modelValue: {
    type: Object,
    required: true
  },
  
  // 表单字段配置
  fields: {
    type: Array,
    required: true
  },
  
  // 表单规则
  rules: {
    type: Object,
    default: () => ({})
  },
  
  // 表单样式
  title: String,
  subtitle: String,
  labelWidth: {
    type: String,
    default: '120px'
  },
  labelPosition: {
    type: String,
    default: 'right',
    validator: (value) => ['left', 'right', 'top'].includes(value)
  },
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['large', 'default', 'small'].includes(value)
  },
  
  // 功能配置
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  validateOnRuleChange: {
    type: Boolean,
    default: true
  },
  scrollToError: {
    type: Boolean,
    default: true
  },
  
  // 步骤配置
  steps: Array,
  currentStep: {
    type: Number,
    default: 0
  },
  stepStatus: {
    type: String,
    default: 'process'
  },
  stepsAlignCenter: {
    type: Boolean,
    default: true
  },
  
  // 底部操作
  showFooter: {
    type: Boolean,
    default: true
  },
  showReset: {
    type: Boolean,
    default: true
  },
  showCancel: {
    type: Boolean,
    default: false
  },
  submitText: {
    type: String,
    default: '提交'
  },
  resetText: {
    type: String,
    default: '重置'
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  submitDisabled: {
    type: Boolean,
    default: false
  },
  
  // 自动保存
  autoSave: {
    type: Boolean,
    default: false
  },
  autoSaveDelay: {
    type: Number,
    default: 2000
  },
  saveTipText: {
    type: String,
    default: '已自动保存'
  },
  
  // 样式类名
  formClass: [String, Array, Object],
  
  // 字段过滤
  visibleFields: Array
})

// Emits
const emit = defineEmits([
  'update:modelValue',
  'submit',
  'cancel',
  'reset',
  'validate',
  'step-change',
  'auto-save'
])

// 响应式数据
const formRef = ref()
const stepLoading = ref(false)
const showSaveTip = ref(false)
const saveTipTimer = ref(null)
const autoSaveTimer = ref(null)

// 计算属性
const formData = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formRules = computed(() => {
  const rules = { ...props.rules }
  
  // 从字段配置中提取规则
  props.fields.forEach(field => {
    if (field.rules && field.prop) {
      rules[field.prop] = field.rules
    }
  })
  
  return rules
})

const formClasses = computed(() => [
  'premium-form',
  props.formClass,
  {
    'with-steps': props.steps && props.steps.length > 0,
    'with-title': props.title,
    'is-loading': props.loading,
    'is-disabled': props.disabled,
    [`size-${props.size}`]: props.size,
    [`label-${props.labelPosition}`]: props.labelPosition
  }
])

const filteredFields = computed(() => {
  let fields = props.fields
  
  // 步骤过滤
  if (props.steps && props.steps.length > 0) {
    fields = fields.filter(field => 
      !field.step || field.step === props.currentStep
    )
  }
  
  // 可见性过滤
  if (props.visibleFields) {
    fields = fields.filter(field => 
      !field.prop || props.visibleFields.includes(field.prop)
    )
  }
  
  // 条件显示过滤
  fields = fields.filter(field => {
    if (typeof field.show === 'function') {
      return field.show(formData.value)
    }
    return field.show !== false
  })
  
  return fields
})

// 方法
const validate = async (callback) => {
  try {
    const valid = await formRef.value?.validate()
    callback && callback(valid)
    emit('validate', valid)
    return valid
  } catch (error) {
    callback && callback(false, error)
    emit('validate', false, error)
    return false
  }
}

const validateField = (prop, callback) => {
  return formRef.value?.validateField(prop, callback)
}

const resetFields = () => {
  formRef.value?.resetFields()
  emit('reset')
}

const clearValidate = (props) => {
  formRef.value?.clearValidate(props)
}

const scrollToField = (prop) => {
  formRef.value?.scrollToField(prop)
}

const handleSubmit = async () => {
  try {
    const valid = await validate()
    if (valid) {
      emit('submit', formData.value)
    }
  } catch (error) {
    console.error('Form validation failed:', error)
  }
}

const handleReset = () => {
  resetFields()
}

const handleCancel = () => {
  emit('cancel')
}

const handleValidate = (prop, isValid, message) => {
  emit('validate', prop, isValid, message)
  
  // 自动保存
  if (props.autoSave && isValid) {
    triggerAutoSave()
  }
}

const nextStep = async () => {
  stepLoading.value = true
  
  try {
    // 验证当前步骤的字段
    const currentStepFields = props.fields
      .filter(field => field.step === props.currentStep && field.prop)
      .map(field => field.prop)
    
    if (currentStepFields.length > 0) {
      for (const prop of currentStepFields) {
        await validateField(prop)
      }
    }
    
    emit('step-change', props.currentStep + 1)
  } catch (error) {
    ElMessage.error('请完善当前步骤的信息')
  } finally {
    stepLoading.value = false
  }
}

const prevStep = () => {
  emit('step-change', props.currentStep - 1)
}

const triggerAutoSave = () => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
  
  autoSaveTimer.value = setTimeout(() => {
    emit('auto-save', formData.value)
    showSaveIndicator()
  }, props.autoSaveDelay)
}

const showSaveIndicator = () => {
  showSaveTip.value = true
  
  if (saveTipTimer.value) {
    clearTimeout(saveTipTimer.value)
  }
  
  saveTipTimer.value = setTimeout(() => {
    showSaveTip.value = false
  }, 2000)
}

// 监听器
watch(
  () => formData.value,
  () => {
    if (props.autoSave) {
      triggerAutoSave()
    }
  },
  { deep: true }
)

// 暴露方法
defineExpose({
  validate,
  validateField,
  resetFields,
  clearValidate,
  scrollToField,
  formRef
})
</script>

<style lang="scss" scoped>
@import '@/assets/styles/design-tokens.scss';
@import '@/assets/styles/animations-enhanced.scss';

.premium-form {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  position: relative;
  transition: all var(--transition-base);
  
  &:hover {
    box-shadow: var(--shadow-md);
  }
  
  &.is-loading {
    opacity: 0.8;
    pointer-events: none;
  }
  
  &.is-disabled {
    opacity: 0.6;
    background: var(--bg-secondary);
  }
  
  // 表单头部
  .form-header {
    padding: var(--spacing-6) var(--spacing-6) var(--spacing-4);
    border-bottom: 1px solid var(--border-light);
    
    .form-title-section {
      .form-title {
        font-size: var(--text-xl);
        font-weight: var(--font-semibold);
        color: var(--text-primary);
        margin: 0 0 var(--spacing-2) 0;
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
      }
      
      .form-subtitle {
        font-size: var(--text-base);
        color: var(--text-secondary);
        margin: 0;
        line-height: 1.5;
      }
    }
  }
  
  // 步骤指示器
  .form-steps {
    padding: var(--spacing-6) var(--spacing-6) var(--spacing-4);
    border-bottom: 1px solid var(--border-light);
    
    :deep(.el-steps) {
      .el-step__title {
        font-weight: var(--font-medium);
        color: var(--text-secondary);
        
        &.is-process {
          color: var(--primary-600);
          font-weight: var(--font-semibold);
        }
        
        &.is-finish {
          color: var(--success-600);
        }
      }
      
      .el-step__description {
        color: var(--text-tertiary);
        font-size: var(--text-sm);
      }
      
      .el-step__icon {
        border-color: var(--border-default);
        
        &.is-process {
          border-color: var(--primary-500);
          background: var(--primary-500);
        }
        
        &.is-finish {
          border-color: var(--success-500);
          background: var(--success-500);
        }
      }
      
      .el-step__line {
        background: var(--border-light);
        
        &.is-finish {
          background: var(--success-500);
        }
      }
    }
  }
  
  // 表单内容
  .form-content {
    padding: var(--spacing-6);
    
    // 表单分组
    .form-group {
      margin-bottom: var(--spacing-8);
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .group-header {
        margin-bottom: var(--spacing-4);
        padding-bottom: var(--spacing-3);
        border-bottom: 1px solid var(--border-light);
        
        .group-title {
          font-size: var(--text-lg);
          font-weight: var(--font-semibold);
          color: var(--text-primary);
          margin: 0 0 var(--spacing-1) 0;
          display: flex;
          align-items: center;
          gap: var(--spacing-2);
          
          .el-icon {
            color: var(--primary-500);
          }
        }
        
        .group-description {
          font-size: var(--text-sm);
          color: var(--text-tertiary);
          margin: 0;
        }
      }
      
      .group-content {
        padding-left: var(--spacing-4);
      }
    }
    
    // 自定义内容
    .form-custom {
      margin: var(--spacing-4) 0;
    }
    
    // 分隔线样式
    :deep(.el-divider) {
      margin: var(--spacing-6) 0;
      
      .el-divider__text {
        background: var(--bg-primary);
        color: var(--text-secondary);
        font-weight: var(--font-medium);
      }
    }
  }
  
  // 表单底部
  .form-footer {
    padding: var(--spacing-4) var(--spacing-6) var(--spacing-6);
    border-top: 1px solid var(--border-light);
    background: var(--bg-secondary);
    
    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: var(--spacing-3);
      
      .step-navigation {
        display: flex;
        gap: var(--spacing-3);
      }
      
      .normal-actions {
        display: flex;
        gap: var(--spacing-3);
      }
      
      .el-button {
        min-width: 100px;
        font-weight: var(--font-medium);
        border-radius: var(--radius-lg);
        transition: all var(--transition-base);
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: var(--shadow-md);
        }
        
        &.el-button--primary {
          background: linear-gradient(135deg, var(--primary-500), var(--primary-600));
          border: none;
          
          &:hover {
            background: linear-gradient(135deg, var(--primary-600), var(--primary-700));
          }
        }
      }
    }
  }
  
  // 保存提示
  .save-tip {
    position: absolute;
    top: var(--spacing-4);
    right: var(--spacing-4);
    background: var(--success-500);
    color: white;
    padding: var(--spacing-2) var(--spacing-4);
    border-radius: var(--radius-full);
    font-size: var(--text-sm);
    display: flex;
    align-items: center;
    gap: var(--spacing-2);
    box-shadow: var(--shadow-lg);
    z-index: var(--z-tooltip);
  }
  
  // 尺寸变体
  &.size-large {
    .form-content {
      padding: var(--spacing-8);
    }
    
    .form-header {
      padding: var(--spacing-8) var(--spacing-8) var(--spacing-6);
    }
    
    .form-footer {
      padding: var(--spacing-6) var(--spacing-8) var(--spacing-8);
    }
  }
  
  &.size-small {
    .form-content {
      padding: var(--spacing-4);
    }
    
    .form-header {
      padding: var(--spacing-4) var(--spacing-4) var(--spacing-3);
    }
    
    .form-footer {
      padding: var(--spacing-3) var(--spacing-4) var(--spacing-4);
    }
  }
  
  // 标签位置变体
  &.label-top {
    :deep(.el-form-item) {
      .el-form-item__label {
        text-align: left;
        margin-bottom: var(--spacing-2);
      }
    }
  }
}

// 表单项样式
:deep(.premium-form-item) {
  margin-bottom: var(--spacing-6);
  
  &:last-child {
    margin-bottom: 0;
  }
  
  .el-form-item__label {
    font-weight: var(--font-medium);
    color: var(--text-primary);
    line-height: 1.5;
    
    .form-label {
      display: flex;
      align-items: center;
      gap: var(--spacing-1);
      
      .help-icon {
        color: var(--text-tertiary);
        cursor: help;
        transition: color var(--transition-fast);
        
        &:hover {
          color: var(--primary-500);
        }
      }
    }
  }
  
  .el-form-item__content {
    line-height: 1.5;
  }
  
  .el-form-item__error {
    color: var(--error-500);
    font-size: var(--text-sm);
    margin-top: var(--spacing-1);
  }
  
  .field-description {
    font-size: var(--text-sm);
    color: var(--text-tertiary);
    margin-top: var(--spacing-1);
    line-height: 1.4;
  }
  
  // 特殊类型样式
  &.type-upload {
    .upload-tip {
      font-size: var(--text-xs);
      color: var(--text-tertiary);
      margin-top: var(--spacing-1);
    }
  }
  
  &.is-required {
    .el-form-item__label::before {
      content: '*';
      color: var(--error-500);
      margin-right: var(--spacing-1);
    }
  }
  
  &.has-help {
    .el-form-item__label {
      .help-icon {
        margin-left: var(--spacing-1);
      }
    }
  }
  
  // 控件样式增强
  .el-input__wrapper,
  .el-textarea__inner,
  .el-select__wrapper {
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-light);
    transition: all var(--transition-fast);
    
    &:hover {
      border-color: var(--border-default);
    }
    
    &.is-focus {
      border-color: var(--primary-500);
      box-shadow: 0 0 0 3px rgba(249, 115, 22, 0.1);
    }
  }
  
  .el-radio-group,
  .el-checkbox-group {
    .el-radio,
    .el-checkbox {
      margin-right: var(--spacing-4);
      margin-bottom: var(--spacing-2);
      
      .el-radio__label,
      .el-checkbox__label {
        color: var(--text-secondary);
        font-weight: var(--font-normal);
      }
      
      &.is-checked {
        .el-radio__label,
        .el-checkbox__label {
          color: var(--primary-600);
          font-weight: var(--font-medium);
        }
      }
    }
  }
  
  .el-switch {
    .el-switch__core {
      border-color: var(--border-default);
      
      &.is-checked {
        background-color: var(--primary-500);
        border-color: var(--primary-500);
      }
    }
  }
  
  .el-slider {
    .el-slider__runway {
      background-color: var(--border-light);
      
      .el-slider__bar {
        background-color: var(--primary-500);
      }
      
      .el-slider__button {
        border-color: var(--primary-500);
        
        &:hover {
          transform: scale(1.1);
        }
      }
    }
  }
  
  .el-rate {
    .el-rate__item {
      .el-rate__icon {
        color: var(--border-default);
        transition: color var(--transition-fast);
        
        &.hover {
          color: var(--warning-400);
        }
      }
      
      &.is-selected {
        .el-rate__icon {
          color: var(--warning-500);
        }
      }
    }
  }
}

// 保存提示动画
.save-tip-enter-active,
.save-tip-leave-active {
  transition: all var(--transition-base);
}

.save-tip-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.save-tip-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

// 响应式设计
@media (max-width: 768px) {
  .premium-form {
    border-radius: var(--radius-lg);
    
    .form-header,
    .form-content,
    .form-footer {
      padding-left: var(--spacing-4);
      padding-right: var(--spacing-4);
    }
    
    .form-steps {
      padding-left: var(--spacing-4);
      padding-right: var(--spacing-4);
      
      :deep(.el-steps) {
        .el-step__title {
          font-size: var(--text-sm);
        }
        
        .el-step__description {
          display: none;
        }
      }
    }
    
    .form-actions {
      flex-direction: column;
      align-items: stretch;
      
      .el-button {
        min-width: auto;
      }
    }
    
    :deep(.premium-form-item) {
      margin-bottom: var(--spacing-4);
      
      .el-form-item__label {
        margin-bottom: var(--spacing-1);
      }
    }
  }
}

// 暗色主题适配
[data-theme="dark"] {
  .premium-form {
    background: var(--bg-primary);
    border-color: var(--border-light);
    
    .form-footer {
      background: var(--bg-secondary);
    }
    
    .save-tip {
      background: var(--success-600);
    }
  }
}
</style>