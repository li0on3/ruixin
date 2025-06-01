import request from './request'

// 获取系统配置
export function getSystemConfigs() {
  return request({
    url: '/admin/system/configs',
    method: 'get'
  })
}

// 更新系统配置
export function updateSystemConfigs(data) {
  return request({
    url: '/admin/system/configs',
    method: 'put',
    data
  })
}

// 获取店铺代码
export function getStoreCode() {
  return request({
    url: '/admin/system/store-code',
    method: 'get'
  })
}

// 更新店铺代码
export function updateStoreCode(storeCode) {
  return request({
    url: '/admin/system/store-code',
    method: 'put',
    data: { store_code: storeCode }
  })
}