import request from './request'

// 获取订单列表
export function getOrders(params) {
  return request({
    url: '/admin/orders',
    method: 'get',
    params
  })
}

// 获取订单详情
export function getOrder(orderNo) {
  return request({
    url: `/admin/orders/${orderNo}`,
    method: 'get'
  })
}

// 刷新订单状态
export function refreshOrderStatus(orderNo) {
  return request({
    url: `/admin/orders/${orderNo}/refresh`,
    method: 'post'
  })
}

// 获取订单统计
export function getOrderStatistics(params) {
  return request({
    url: '/admin/orders/statistics',
    method: 'get',
    params
  })
}