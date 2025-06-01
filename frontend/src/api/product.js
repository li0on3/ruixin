import request from './request'

// 获取商品列表
export function getProducts(params) {
  return request({
    url: '/admin/products',
    method: 'get',
    params
  })
}

// 搜索商品
export function searchProducts(keyword) {
  return request({
    url: '/admin/products/search',
    method: 'get',
    params: { keyword }
  })
}

// 同步商品
export function syncProducts(data) {
  return request({
    url: '/admin/products/sync',
    method: 'post',
    data
  })
}

// 获取匹配失败日志
export function getMatchLogs(params) {
  return request({
    url: '/admin/products/match-logs',
    method: 'get',
    params
  })
}

// 根据代码获取商品
export function getProductsByCodes(codes) {
  return request({
    url: '/admin/products/by-codes',
    method: 'post',
    data: { codes }
  })
}