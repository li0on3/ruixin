import request from './request'

// 搜索店铺
export function searchStores(params) {
  return request({
    url: '/admin/stores/search',
    method: 'get',
    params
  })
}

// 获取城市列表
export function getCities() {
  return request({
    url: '/admin/stores/cities',
    method: 'get'
  })
}

// 获取热门店铺
export function getHotStores() {
  return request({
    url: '/admin/stores/hot',
    method: 'get'
  })
}

// 同步城市列表
export function syncCities(cardCode) {
  return request({
    url: '/admin/stores/sync-cities',
    method: 'post',
    data: { card_code: cardCode }
  })
}