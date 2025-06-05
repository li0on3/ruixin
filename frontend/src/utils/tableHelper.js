/**
 * 表格助手工具
 * 用于管理表格的列显示设置等功能
 */

/**
 * 清除表格列显示设置
 * 当表格列显示异常时，可以调用此方法重置
 */
export function clearTableColumnSettings() {
  localStorage.removeItem('table-visible-columns')
  console.log('表格列设置已清除，刷新页面后将恢复默认显示')
}

/**
 * 获取表格列设置
 * @returns {Array|null} 返回保存的列设置或null
 */
export function getTableColumnSettings() {
  const saved = localStorage.getItem('table-visible-columns')
  if (saved) {
    try {
      return JSON.parse(saved)
    } catch (e) {
      console.error('解析表格列设置失败:', e)
      return null
    }
  }
  return null
}

/**
 * 保存表格列设置
 * @param {Array} columns - 要保存的列数组
 */
export function saveTableColumnSettings(columns) {
  if (Array.isArray(columns)) {
    localStorage.setItem('table-visible-columns', JSON.stringify(columns))
  }
}

// 导出到全局，方便在控制台调试使用
if (typeof window !== 'undefined') {
  window.clearTableColumnSettings = clearTableColumnSettings
}