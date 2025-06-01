import request from './request'

// 获取核心指标
export function getMetrics(params) {
  return request({
    url: '/admin/statistics/metrics',
    method: 'get',
    params
  })
}

// 获取销售趋势
export function getSalesTrend(params) {
  return request({
    url: '/admin/statistics/sales-trend',
    method: 'get',
    params
  })
}

// 获取分销商排行
export function getDistributorRank(params) {
  return request({
    url: '/admin/statistics/distributor-rank',
    method: 'get',
    params
  })
}

// 获取商品分析
export function getProductAnalysis(params) {
  return request({
    url: '/admin/statistics/product-analysis',
    method: 'get',
    params
  })
}

// 获取时段分布
export function getHourDistribution(params) {
  return request({
    url: '/admin/statistics/hour-distribution',
    method: 'get',
    params
  })
}

// 获取地区分布
export function getRegionDistribution(params) {
  return request({
    url: '/admin/statistics/region-distribution',
    method: 'get',
    params
  })
}

// 获取详细数据
export function getDetailData(params) {
  return request({
    url: '/admin/statistics/details',
    method: 'get',
    params
  })
}

// 导出数据
export function exportData(params) {
  return request({
    url: '/admin/statistics/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}