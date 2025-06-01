import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'

// 成功消息
export const showSuccess = (message, options = {}) => {
  ElMessage({
    message,
    type: 'success',
    duration: 3000,
    showClose: true,
    grouping: true,
    customClass: 'custom-message-success',
    ...options
  })
}

// 错误消息
export const showError = (message, options = {}) => {
  ElMessage({
    message,
    type: 'error',
    duration: 4000,
    showClose: true,
    grouping: true,
    customClass: 'custom-message-error',
    ...options
  })
}

// 警告消息
export const showWarning = (message, options = {}) => {
  ElMessage({
    message,
    type: 'warning',
    duration: 3500,
    showClose: true,
    grouping: true,
    customClass: 'custom-message-warning',
    ...options
  })
}

// 信息消息
export const showInfo = (message, options = {}) => {
  ElMessage({
    message,
    type: 'info',
    duration: 3000,
    showClose: true,
    grouping: true,
    customClass: 'custom-message-info',
    ...options
  })
}

// 确认对话框
export const showConfirm = (message, title = '提示', options = {}) => {
  return ElMessageBox.confirm(message, title, {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    center: true,
    roundButton: true,
    distinguishCancelAndClose: true,
    ...options
  })
}

// 删除确认
export const showDeleteConfirm = (itemName = '此项') => {
  return showConfirm(
    `删除后数据将无法恢复，确定要删除${itemName}吗？`,
    '删除确认',
    { type: 'error' }
  )
}

// 操作成功通知
export const showSuccessNotify = (title, message = '', options = {}) => {
  ElNotification({
    title,
    message,
    type: 'success',
    duration: 4000,
    position: 'top-right',
    ...options
  })
}

// 操作失败通知
export const showErrorNotify = (title, message = '', options = {}) => {
  ElNotification({
    title,
    message,
    type: 'error',
    duration: 5000,
    position: 'top-right',
    ...options
  })
}

// 导出所有方法
export default {
  showSuccess,
  showError,
  showWarning,
  showInfo,
  showConfirm,
  showDeleteConfirm,
  showSuccessNotify,
  showErrorNotify
}