/**
 * 安全的导出工具函数
 * 增强了原有export.js的安全性
 */
import * as XLSX from 'xlsx'
import { ElMessage } from 'element-plus'
import { safeDownload } from './securityUtils'

/**
 * 安全的导出数据到Excel
 * @param {Array} data - 要导出的数据数组
 * @param {Array} headers - 表头配置
 * @param {String} filename - 文件名
 * @param {Object} options - 额外选项
 */
export function safeExportToExcel(data, headers, filename = 'export', options = {}) {
  try {
    // 数据验证
    if (!Array.isArray(data)) {
      throw new Error('导出数据必须是数组');
    }
    
    // 限制导出数据量
    const MAX_ROWS = options.maxRows || 50000;
    if (data.length > MAX_ROWS) {
      ElMessage.warning(`数据量过大，只导出前${MAX_ROWS}条记录`);
      data = data.slice(0, MAX_ROWS);
    }
    
    // 处理表头
    const headerRow = headers.map(h => sanitizeString(h.title || ''));
    
    // 处理数据行
    const dataRows = data.map(item => {
      return headers.map(header => {
        let value = item[header.key];
        
        // 数据净化
        if (value === null || value === undefined) {
          return '';
        }
        
        // 字符串净化
        if (typeof value === 'string') {
          value = sanitizeString(value);
        }
        
        // 处理日期
        if (header.type === 'date' && value) {
          return formatDate(value);
        }
        
        // 处理金额
        if (header.type === 'currency' && typeof value === 'number') {
          return `¥${value.toFixed(2)}`;
        }
        
        // 处理百分比
        if (header.type === 'percent' && typeof value === 'number') {
          return `${value.toFixed(2)}%`;
        }
        
        // 处理状态映射
        if (header.type === 'status' && header.statusMap) {
          return sanitizeString(header.statusMap[value] || value);
        }
        
        return value;
      });
    });
    
    // 创建工作表（使用安全选项）
    const wsData = [headerRow, ...dataRows];
    const ws = XLSX.utils.aoa_to_sheet(wsData, {
      cellDates: false,  // 禁用日期格式
      sheetStubs: false  // 禁用空单元格
    });
    
    // 设置列宽
    const colWidths = headers.map(h => ({
      wch: Math.min(h.width || 15, 100)  // 限制最大列宽
    }));
    ws['!cols'] = colWidths;
    
    // 创建工作簿
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
    
    // 生成安全的文件名
    const timestamp = new Date().toISOString().slice(0, 10);
    const safeFilename = sanitizeFilename(`${filename}_${timestamp}.xlsx`);
    
    // 导出文件（使用安全的写入选项）
    const wbout = XLSX.write(wb, {
      bookType: 'xlsx',
      type: 'array',
      cellDates: false,
      cellFormula: false,  // 禁用公式
      cellHTML: false,     // 禁用HTML
      cellNF: false,       // 禁用数字格式
      cellStyles: false    // 禁用样式
    });
    
    // 创建Blob并下载
    const blob = new Blob([wbout], { type: 'application/octet-stream' });
    safeDownload(blob, safeFilename);
    
    ElMessage.success('导出成功');
  } catch (error) {
    console.error('导出失败:', error);
    ElMessage.error(`导出失败: ${error.message}`);
  }
}

/**
 * 净化字符串，移除潜在的危险内容
 * @param {string} str - 要净化的字符串
 * @returns {string} 净化后的字符串
 */
function sanitizeString(str) {
  if (typeof str !== 'string') {
    return String(str || '');
  }
  
  return str
    // 移除不可见字符和控制字符
    .replace(/[\x00-\x1F\x7F-\x9F]/g, '')
    // 移除潜在的公式注入
    .replace(/^[=+\-@\t\r]/g, "'")
    // 限制长度
    .slice(0, 32767)  // Excel单元格最大长度
    .trim();
}

/**
 * 净化文件名
 * @param {string} filename - 原始文件名
 * @returns {string} 净化后的文件名
 */
function sanitizeFilename(filename) {
  return filename
    .replace(/[^\w\u4e00-\u9fa5.-]/g, '_')
    .replace(/\.{2,}/g, '.')
    .slice(0, 255);  // 文件名最大长度
}

/**
 * 格式化日期
 * @param {String|Date} date - 日期
 * @returns {String} 格式化后的日期字符串
 */
function formatDate(date) {
  if (!date) return '';
  
  try {
    const d = new Date(date);
    if (isNaN(d.getTime())) {
      return '';
    }
    
    return d.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  } catch {
    return '';
  }
}

/**
 * 安全的CSV导出
 * @param {Array} data - 数据
 * @param {Array} headers - 表头
 * @param {String} filename - 文件名
 * @param {Object} options - 选项
 */
export function safeExportToCSV(data, headers, filename = 'export', options = {}) {
  try {
    // 数据验证
    if (!Array.isArray(data)) {
      throw new Error('导出数据必须是数组');
    }
    
    // 限制导出数据量
    const MAX_ROWS = options.maxRows || 100000;
    if (data.length > MAX_ROWS) {
      ElMessage.warning(`数据量过大，只导出前${MAX_ROWS}条记录`);
      data = data.slice(0, MAX_ROWS);
    }
    
    // 处理表头
    const headerRow = headers.map(h => escapeCSV(sanitizeString(h.title || '')));
    
    // 处理数据行
    const dataRows = data.map(item => {
      return headers.map(header => {
        let value = item[header.key];
        
        if (value === null || value === undefined) {
          return '';
        }
        
        // 根据类型处理值
        if (header.type === 'date' && value) {
          value = formatDate(value);
        } else if (header.type === 'currency' && typeof value === 'number') {
          value = `¥${value.toFixed(2)}`;
        } else if (header.type === 'percent' && typeof value === 'number') {
          value = `${value.toFixed(2)}%`;
        } else if (header.type === 'status' && header.statusMap) {
          value = header.statusMap[value] || value;
        }
        
        return escapeCSV(sanitizeString(String(value)));
      });
    });
    
    // 组合CSV内容
    const csvContent = [
      headerRow.join(','),
      ...dataRows.map(row => row.join(','))
    ].join('\n');
    
    // 添加BOM以支持中文
    const BOM = '\uFEFF';
    const blob = new Blob([BOM + csvContent], { 
      type: 'text/csv;charset=utf-8;' 
    });
    
    // 安全下载
    const timestamp = new Date().toISOString().slice(0, 10);
    const safeFilename = sanitizeFilename(`${filename}_${timestamp}.csv`);
    safeDownload(blob, safeFilename);
    
    ElMessage.success('导出成功');
  } catch (error) {
    console.error('导出失败:', error);
    ElMessage.error(`导出失败: ${error.message}`);
  }
}

/**
 * 转义CSV特殊字符
 * @param {string} str - 要转义的字符串
 * @returns {string} 转义后的字符串
 */
function escapeCSV(str) {
  if (typeof str !== 'string') {
    str = String(str || '');
  }
  
  // 如果包含特殊字符，需要用引号包裹
  if (str.includes(',') || str.includes('"') || str.includes('\n') || str.includes('\r')) {
    // 转义引号
    str = str.replace(/"/g, '""');
    return `"${str}"`;
  }
  
  return str;
}

// 导出便捷方法
export { formatDate, sanitizeString, sanitizeFilename };

// 默认导出
export default {
  safeExportToExcel,
  safeExportToCSV,
  formatDate,
  sanitizeString,
  sanitizeFilename
};