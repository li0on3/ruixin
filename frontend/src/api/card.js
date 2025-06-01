import request from './request'

// 获取卡片列表
export function getCards(params) {
  return request({
    url: '/admin/cards',
    method: 'get',
    params
  })
}

// 获取卡片详情
export function getCard(id) {
  return request({
    url: `/admin/cards/${id}`,
    method: 'get'
  })
}

// 创建卡片
export function createCard(data) {
  return request({
    url: '/admin/cards',
    method: 'post',
    data
  })
}

// 更新卡片
export function updateCard(id, data) {
  return request({
    url: `/admin/cards/${id}`,
    method: 'put',
    data
  })
}

// 删除卡片
export function deleteCard(id) {
  return request({
    url: `/admin/cards/${id}`,
    method: 'delete'
  })
}

// 获取卡片使用日志
export function getCardUsageLogs(id, params) {
  return request({
    url: `/admin/cards/${id}/usage-logs`,
    method: 'get',
    params
  })
}

// 批量导入卡片
export function batchImportCards(data) {
  return request({
    url: '/admin/cards/batch-import',
    method: 'post',
    data
  })
}

// 获取卡片统计信息
export function getCardStats(params) {
  return request({
    url: '/admin/cards/stats',
    method: 'get',
    params
  })
}

// 获取批次列表
export function getBatches(params) {
  return request({
    url: '/admin/card-batches',
    method: 'get',
    params
  })
}

// 获取批次详情
export function getBatch(id) {
  return request({
    url: `/admin/card-batches/${id}`,
    method: 'get'
  })
}

// 获取批次下的卡片
export function getBatchCards(id, params) {
  return request({
    url: `/admin/card-batches/${id}/cards`,
    method: 'get',
    params
  })
}

// 验证卡片是否可用
export function validateCard(data) {
  return request({
    url: '/admin/cards/validate',
    method: 'post',
    data: { card_code: data.cardCode }
  })
}

// 批量验证卡片
export function batchValidateCards(data) {
  return request({
    url: '/admin/cards/batch-validate',
    method: 'post',
    data: { card_codes: data.cardCodes }
  })
}

// 启动双模式批量验证
export function startBatchValidation(mode) {
  return request({
    url: '/admin/cards/batch-validation/start',
    method: 'post',
    data: { mode }
  })
}

// 获取验证进度
export function getValidationProgress(taskId) {
  return request({
    url: `/admin/cards/batch-validation/${taskId}`,
    method: 'get'
  })
}

// 取消验证任务
export function cancelValidation(taskId) {
  return request({
    url: `/admin/cards/batch-validation/${taskId}`,
    method: 'delete'
  })
}

// 获取验证统计信息
export function getValidationStatistics() {
  return request({
    url: '/admin/cards/validation-stats',
    method: 'get'
  })
}

// 获取未使用的卡片列表（用于下拉选择）
export function getUnusedCards(params = {}) {
  return request({
    url: '/admin/cards',
    method: 'get',
    params: {
      status: 0, // 只获取未使用的卡片
      page_size: 100, // 获取更多数据用于下拉框
      ...params
    }
  })
}