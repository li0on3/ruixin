// 导出工具函数
import * as XLSX from 'xlsx'
import { ElMessage } from 'element-plus'

/**
 * 导出数据到Excel
 * @param {Array} data - 要导出的数据数组
 * @param {Array} headers - 表头配置 [{key: 'field', title: '标题', width: 100}]
 * @param {String} filename - 文件名（不含扩展名）
 */
export function exportToExcel(data, headers, filename = 'export') {
  try {
    // 处理表头
    const headerRow = headers.map(h => h.title)
    
    // 处理数据行
    const dataRows = data.map(item => {
      return headers.map(header => {
        let value = item[header.key]
        
        // 处理特殊数据类型
        if (value === null || value === undefined) {
          return ''
        }
        
        // 处理日期
        if (header.type === 'date' && value) {
          return formatDate(value)
        }
        
        // 处理金额
        if (header.type === 'currency' && typeof value === 'number') {
          return `¥${value.toFixed(2)}`
        }
        
        // 处理百分比
        if (header.type === 'percent' && typeof value === 'number') {
          return `${value.toFixed(2)}%`
        }
        
        // 处理状态映射
        if (header.type === 'status' && header.statusMap) {
          return header.statusMap[value] || value
        }
        
        return value
      })
    })
    
    // 创建工作表
    const wsData = [headerRow, ...dataRows]
    const ws = XLSX.utils.aoa_to_sheet(wsData)
    
    // 设置列宽
    const colWidths = headers.map(h => ({
      wch: h.width || 15
    }))
    ws['!cols'] = colWidths
    
    // 创建工作簿
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1')
    
    // 导出文件
    const timestamp = new Date().toISOString().slice(0, 10)
    XLSX.writeFile(wb, `${filename}_${timestamp}.xlsx`)
    
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}

/**
 * 导出数据到CSV
 * @param {Array} data - 要导出的数据数组
 * @param {Array} headers - 表头配置
 * @param {String} filename - 文件名（不含扩展名）
 */
export function exportToCSV(data, headers, filename = 'export') {
  try {
    // 处理表头
    const headerRow = headers.map(h => h.title).join(',')
    
    // 处理数据行
    const dataRows = data.map(item => {
      return headers.map(header => {
        let value = item[header.key]
        
        // 处理特殊字符和换行
        if (typeof value === 'string') {
          value = value.replace(/"/g, '""')
          if (value.includes(',') || value.includes('\n') || value.includes('"')) {
            value = `"${value}"`
          }
        }
        
        return value || ''
      }).join(',')
    })
    
    // 组合CSV内容
    const csvContent = [headerRow, ...dataRows].join('\n')
    
    // 添加BOM以支持中文
    const BOM = '\uFEFF'
    const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })
    
    // 创建下载链接
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    const timestamp = new Date().toISOString().slice(0, 10)
    link.setAttribute('download', `${filename}_${timestamp}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}

/**
 * 根据页面配置导出数据
 * @param {String} pageType - 页面类型 (cards, orders, distributors, etc.)
 * @param {Array} data - 数据
 * @param {Object} options - 额外选项
 */
export function exportPageData(pageType, data, options = {}) {
  const configs = {
    cards: {
      filename: '卡片数据',
      headers: [
        { key: 'id', title: 'ID', width: 8 },
        { key: 'card_code', title: '卡片代码', width: 20 },
        { key: 'batch_id', title: '批次ID', width: 12 },
        { key: 'cost_price', title: '成本价', width: 12, type: 'currency' },
        { key: 'sell_price', title: '销售价', width: 12, type: 'currency' },
        { key: 'status', title: '状态', width: 10, type: 'status', statusMap: { 0: '未使用', 1: '已使用' } },
        { key: 'sync_status', title: '同步状态', width: 12 },
        { key: 'created_at', title: '创建时间', width: 20, type: 'date' },
        { key: 'used_at', title: '使用时间', width: 20, type: 'date' },
        { key: 'expired_at', title: '过期时间', width: 20, type: 'date' }
      ]
    },
    orders: {
      filename: '订单数据',
      headers: [
        { key: 'id', title: 'ID', width: 8 },
        { key: 'order_no', title: '订单号', width: 20 },
        { key: 'distributor_name', title: '分销商', width: 15 },
        { key: 'store_name', title: '店铺名称', width: 20 },
        { key: 'total_amount', title: '订单金额', width: 12, type: 'currency' },
        { key: 'profit_amount', title: '利润金额', width: 12, type: 'currency' },
        { key: 'status', title: '状态', width: 10, type: 'status', statusMap: { 0: '待处理', 1: '处理中', 2: '已完成', 3: '失败' } },
        { key: 'created_at', title: '创建时间', width: 20, type: 'date' },
        { key: 'completed_at', title: '完成时间', width: 20, type: 'date' }
      ]
    },
    distributors: {
      filename: '分销商数据',
      headers: [
        { key: 'id', title: 'ID', width: 8 },
        { key: 'name', title: '分销商名称', width: 20 },
        { key: 'phone', title: '手机号', width: 15 },
        { key: 'email', title: '邮箱', width: 25 },
        { key: 'balance', title: '余额', width: 12, type: 'currency' },
        { key: 'total_orders', title: '总订单数', width: 12 },
        { key: 'total_amount', title: '总交易额', width: 12, type: 'currency' },
        { key: 'status', title: '状态', width: 10, type: 'status', statusMap: { 1: '正常', 0: '禁用' } },
        { key: 'created_at', title: '创建时间', width: 20, type: 'date' }
      ]
    },
    statistics: {
      filename: '统计报表',
      headers: [
        { key: 'date', title: '日期', width: 15 },
        { key: 'orders', title: '订单数', width: 10 },
        { key: 'revenue', title: '营收', width: 15, type: 'currency' },
        { key: 'profit', title: '利润', width: 15, type: 'currency' },
        { key: 'success_rate', title: '成功率', width: 10, type: 'percent' }
      ]
    }
  }
  
  const config = configs[pageType]
  if (!config) {
    ElMessage.error('不支持的导出类型')
    return
  }
  
  const filename = options.filename || config.filename
  const format = options.format || 'excel'
  
  if (format === 'csv') {
    exportToCSV(data, config.headers, filename)
  } else {
    exportToExcel(data, config.headers, filename)
  }
}

/**
 * 格式化日期
 * @param {String|Date} date - 日期
 * @returns {String} 格式化后的日期字符串
 */
function formatDate(date) {
  if (!date) return ''
  
  try {
    const d = new Date(date)
    return d.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch {
    return String(date)
  }
}

// 导出特定格式的统计报表
export function exportDashboardReport(data, options = {}) {
  const { 
    statistics = {},
    orderTrend = [],
    hotProducts = [],
    recentOrders = [],
    distributorRank = []
  } = data
  
  try {
    // 创建工作簿
    const wb = XLSX.utils.book_new()
    
    // 概览统计表
    if (Object.keys(statistics).length > 0) {
      const overviewData = [
        ['指标', '数值'],
        ['今日订单', statistics.todayOrders || 0],
        ['今日销售额', `¥${(statistics.todayAmount || 0).toFixed(2)}`],
        ['活跃分销商', statistics.activeDistributors || 0],
        ['可用卡片', statistics.activeCards || 0]
      ]
      const overviewWs = XLSX.utils.aoa_to_sheet(overviewData)
      XLSX.utils.book_append_sheet(wb, overviewWs, '概览统计')
    }
    
    // 趋势数据表
    if (orderTrend.length > 0) {
      const trendHeaders = ['日期', '订单数', '销售额', '利润']
      const trendData = orderTrend.map(item => [
        item.date,
        item.orders || 0,
        `¥${(item.revenue || 0).toFixed(2)}`,
        `¥${(item.profit || 0).toFixed(2)}`
      ])
      const trendWs = XLSX.utils.aoa_to_sheet([trendHeaders, ...trendData])
      XLSX.utils.book_append_sheet(wb, trendWs, '趋势数据')
    }
    
    // 热销产品表
    if (hotProducts.length > 0) {
      const productsHeaders = ['产品名称', '销量', '占比']
      const productsData = hotProducts.map(item => [
        item.name,
        item.count,
        `${item.percentage}%`
      ])
      const productsWs = XLSX.utils.aoa_to_sheet([productsHeaders, ...productsData])
      XLSX.utils.book_append_sheet(wb, productsWs, '热销产品')
    }
    
    // 最新订单表
    if (recentOrders.length > 0) {
      const ordersHeaders = ['订单号', '分销商', '店铺', '金额', '状态', '时间']
      const ordersData = recentOrders.map(item => [
        item.orderNo,
        item.distributorName,
        item.storeName,
        `¥${item.totalAmount}`,
        getOrderStatusText(item.status),
        formatDate(item.createdAt)
      ])
      const ordersWs = XLSX.utils.aoa_to_sheet([ordersHeaders, ...ordersData])
      XLSX.utils.book_append_sheet(wb, ordersWs, '最新订单')
    }
    
    // 分销商排行表
    if (distributorRank.length > 0) {
      const rankHeaders = ['排名', '分销商', '订单数', '营收']
      const rankData = distributorRank.map((item, index) => [
        index + 1,
        item.name,
        item.orders,
        `¥${item.revenue.toFixed(2)}`
      ])
      const rankWs = XLSX.utils.aoa_to_sheet([rankHeaders, ...rankData])
      XLSX.utils.book_append_sheet(wb, rankWs, '分销商排行')
    }
    
    // 导出文件
    const timestamp = new Date().toISOString().slice(0, 10)
    const filename = options.filename || `仪表盘报表_${timestamp}.xlsx`
    XLSX.writeFile(wb, filename)
    
    ElMessage.success('报表导出成功')
  } catch (error) {
    console.error('导出报表失败:', error)
    ElMessage.error('导出报表失败')
  }
}

// 获取订单状态文本
function getOrderStatusText(status) {
  const statusMap = {
    0: '待处理',
    1: '处理中', 
    2: '已完成',
    3: '失败'
  }
  return statusMap[status] || '未知'
}