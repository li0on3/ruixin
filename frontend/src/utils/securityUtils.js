/**
 * 安全工具集
 * 用于增强应用安全性，特别是处理文件上传和Excel解析
 */

/**
 * 验证Excel文件
 * @param {File} file - 要验证的文件
 * @returns {boolean} 验证是否通过
 * @throws {Error} 验证失败时抛出错误
 */
export function validateExcelFile(file) {
  const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
  const ALLOWED_EXTENSIONS = ['.xlsx', '.xls'];
  const ALLOWED_MIME_TYPES = [
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-excel'
  ];

  // 检查文件是否存在
  if (!file) {
    throw new Error('请选择文件');
  }

  // 检查文件大小
  if (file.size > MAX_FILE_SIZE) {
    throw new Error(`文件大小不能超过 ${MAX_FILE_SIZE / 1024 / 1024}MB`);
  }

  // 检查文件扩展名
  const fileName = file.name.toLowerCase();
  const ext = fileName.slice(fileName.lastIndexOf('.'));
  if (!ALLOWED_EXTENSIONS.includes(ext)) {
    throw new Error('只支持 .xlsx 和 .xls 格式的文件');
  }

  // 检查MIME类型
  if (file.type && !ALLOWED_MIME_TYPES.includes(file.type)) {
    console.warn('文件MIME类型不匹配:', file.type);
  }

  return true;
}

/**
 * 安全读取Excel文件
 * @param {ArrayBuffer|string} data - Excel文件数据
 * @param {Object} options - 额外选项
 * @returns {Object} 解析后的工作簿对象
 */
export function safeReadExcel(data, options = {}) {
  try {
    // 使用安全的解析选项
    const safeOptions = {
      type: options.type || 'binary',
      cellFormula: false,    // 禁用公式解析，防止公式注入
      cellHTML: false,       // 禁用HTML解析，防止XSS
      cellNF: false,         // 禁用数字格式
      cellStyles: false,     // 禁用样式
      sheetStubs: false,     // 禁用空单元格
      password: undefined,   // 不支持密码保护的文件
      ...options            // 允许覆盖某些选项
    };

    // 动态导入xlsx以避免打包时的问题
    const XLSX = window.XLSX || require('xlsx');
    const workbook = XLSX.read(data, safeOptions);

    // 验证工作簿
    validateWorkbook(workbook);

    return workbook;
  } catch (error) {
    console.error('Excel解析失败:', error);
    throw new Error('无效的Excel文件或文件已损坏');
  }
}

/**
 * 验证工作簿内容
 * @param {Object} workbook - 工作簿对象
 * @throws {Error} 验证失败时抛出错误
 */
function validateWorkbook(workbook) {
  const MAX_SHEETS = 10;
  const MAX_ROWS_PER_SHEET = 10000;
  const MAX_COLS_PER_SHEET = 100;

  // 检查工作表数量
  if (!workbook.SheetNames || workbook.SheetNames.length === 0) {
    throw new Error('Excel文件中没有工作表');
  }

  if (workbook.SheetNames.length > MAX_SHEETS) {
    throw new Error(`Excel文件包含过多工作表（最多${MAX_SHEETS}个）`);
  }

  // 检查每个工作表
  for (const sheetName of workbook.SheetNames) {
    const sheet = workbook.Sheets[sheetName];
    if (!sheet) continue;

    // 获取工作表范围
    const range = sheet['!ref'];
    if (range) {
      const XLSX = window.XLSX || require('xlsx');
      const decodedRange = XLSX.utils.decode_range(range);
      
      const rowCount = decodedRange.e.r - decodedRange.s.r + 1;
      const colCount = decodedRange.e.c - decodedRange.s.c + 1;

      if (rowCount > MAX_ROWS_PER_SHEET) {
        throw new Error(`工作表"${sheetName}"包含过多行（最多${MAX_ROWS_PER_SHEET}行）`);
      }

      if (colCount > MAX_COLS_PER_SHEET) {
        throw new Error(`工作表"${sheetName}"包含过多列（最多${MAX_COLS_PER_SHEET}列）`);
      }
    }
  }
}

/**
 * 清理和净化Excel数据
 * @param {Array} data - 从Excel读取的数据数组
 * @returns {Array} 清理后的数据
 */
export function sanitizeExcelData(data) {
  if (!Array.isArray(data)) {
    return [];
  }

  return data.map(row => {
    const cleanRow = {};
    for (const [key, value] of Object.entries(row)) {
      // 清理键名（移除特殊字符）
      const cleanKey = key.replace(/[^\w\u4e00-\u9fa5]/g, '_');
      
      // 清理值
      if (typeof value === 'string') {
        // 移除潜在的脚本标签和危险内容
        cleanRow[cleanKey] = value
          .replace(/<script[^>]*>.*?<\/script>/gi, '')
          .replace(/<[^>]+>/g, '')
          .trim();
      } else if (typeof value === 'number' && !isNaN(value)) {
        cleanRow[cleanKey] = value;
      } else if (value instanceof Date) {
        cleanRow[cleanKey] = value.toISOString();
      } else {
        cleanRow[cleanKey] = String(value || '');
      }
    }
    return cleanRow;
  });
}

/**
 * 创建安全的文件下载
 * @param {Blob} blob - 要下载的文件内容
 * @param {string} filename - 文件名
 */
export function safeDownload(blob, filename) {
  // 清理文件名
  const safeFilename = filename.replace(/[^\w\u4e00-\u9fa5.-]/g, '_');
  
  // 创建下载链接
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = safeFilename;
  
  // 安全属性
  link.rel = 'noopener noreferrer';
  
  // 触发下载
  document.body.appendChild(link);
  link.click();
  
  // 清理
  setTimeout(() => {
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  }, 100);
}

/**
 * 检查是否在安全的环境中运行
 * @returns {Object} 环境安全信息
 */
export function checkEnvironmentSecurity() {
  const security = {
    isSecureContext: window.isSecureContext || false,
    protocol: window.location.protocol,
    isLocalhost: ['localhost', '127.0.0.1'].includes(window.location.hostname),
    isDevelopment: process.env.NODE_ENV === 'development'
  };

  // 开发环境警告
  if (security.isDevelopment && !security.isLocalhost) {
    console.warn('警告：开发服务器正在非本地环境运行，可能存在安全风险');
  }

  return security;
}

// 导出工具集
export default {
  validateExcelFile,
  safeReadExcel,
  sanitizeExcelData,
  safeDownload,
  checkEnvironmentSecurity
};