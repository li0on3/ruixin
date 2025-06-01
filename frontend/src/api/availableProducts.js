import request from './request'

// 获取可用商品列表（管理员）
export const getAvailableProducts = () => {
  return request({
    url: '/admin/available-products',
    method: 'get'
  })
}

// 获取可用商品列表（分销商）
export const getDistributorAvailableProducts = () => {
  return request({
    url: '/distributor/available-products',
    method: 'get'
  })
}

// 导出商品数据为Excel
export const exportAvailableProducts = (data) => {
  const headers = ['商品编码', '商品名称', '类别', '库存状态', '价格区间', 'SKU数量', '别名']
  const rows = data.products.map(product => [
    product.goods_code,
    product.goods_name,
    product.category,
    getStockStatusText(product.stock_status),
    product.price_range,
    product.skus.length,
    product.aliases.join(', ')
  ])
  
  // 创建CSV内容
  const csvContent = [
    headers.join(','),
    ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
  ].join('\n')
  
  // 添加BOM以支持中文
  const BOM = '\uFEFF'
  const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', `可用商品列表_${new Date().toLocaleDateString()}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 获取库存状态文本
export const getStockStatusText = (status) => {
  const statusMap = {
    'high': '库存充足',
    'medium': '库存一般',
    'low': '库存较少',
    'out': '暂时缺货'
  }
  return statusMap[status] || status
}

// 获取库存状态类型（用于Element UI标签）
export const getStockStatusType = (status) => {
  const typeMap = {
    'high': 'success',
    'medium': '',
    'low': 'warning',
    'out': 'danger'
  }
  return typeMap[status] || 'info'
}