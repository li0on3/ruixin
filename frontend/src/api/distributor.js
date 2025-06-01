import request from './request'

// 获取分销商列表
export function getDistributors(params) {
  return request({
    url: '/admin/distributors',
    method: 'get',
    params
  })
}

// 获取分销商详情
export function getDistributor(id) {
  return request({
    url: `/admin/distributors/${id}`,
    method: 'get'
  })
}

// 创建分销商
export function createDistributor(data) {
  return request({
    url: '/admin/distributors',
    method: 'post',
    data
  })
}

// 更新分销商
export function updateDistributor(id, data) {
  return request({
    url: `/admin/distributors/${id}`,
    method: 'put',
    data
  })
}

// 删除分销商
export function deleteDistributor(id) {
  return request({
    url: `/admin/distributors/${id}`,
    method: 'delete'
  })
}

// 重置API密钥
export function resetDistributorAPIKey(id) {
  return request({
    url: `/admin/distributors/${id}/reset-api-key`,
    method: 'post'
  })
}

// 获取API调用日志
export function getDistributorAPILogs(id, params) {
  return request({
    url: `/admin/distributors/${id}/api-logs`,
    method: 'get',
    params
  })
}