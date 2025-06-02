import request from './request'

// 获取仪表盘统计数据
export function getDashboardStatistics(params) {
  return request({
    url: '/admin/dashboard/statistics',
    method: 'get',
    params
  })
}

// 获取订单趋势
export function getOrderTrend(dateRange) {
  return request({
    url: '/admin/dashboard/order-trend',
    method: 'get',
    params: { dateRange }
  })
}

// 获取热门商品
export function getHotGoods(params) {
  return request({
    url: '/admin/dashboard/hot-goods',
    method: 'get',
    params
  })
}

// 获取最新订单
export function getRecentOrders(params) {
  return request({
    url: '/admin/dashboard/recent-orders',
    method: 'get',
    params
  })
}

// 获取分销商排名
export function getDistributorRank(period = 'today') {
  return request({
    url: '/admin/dashboard/distributor-rank',
    method: 'get',
    params: { period }
  })
}

// 获取小时分布
export function getHourDistribution() {
  return request({
    url: '/admin/dashboard/hour-distribution',
    method: 'get'
  })
}