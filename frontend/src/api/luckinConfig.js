import request from './request'

// 价格管理相关API

// 获取价格列表
export function getPriceList(params) {
  return request({
    url: '/admin/luckin/prices',
    method: 'get',
    params
  })
}

// 创建价格
export function createPrice(data) {
  return request({
    url: '/admin/luckin/prices',
    method: 'post',
    data
  })
}

// 更新价格
export function updatePrice(id, data) {
  return request({
    url: `/admin/luckin/prices/${id}`,
    method: 'put',
    data
  })
}

// 删除价格
export function deletePrice(id) {
  return request({
    url: `/admin/luckin/prices/${id}`,
    method: 'delete'
  })
}

// 产品管理相关API

// 获取产品列表
export function getProductList(params) {
  return request({
    url: '/admin/luckin/products',
    method: 'get',
    params
  })
}

// 创建产品
export function createProduct(data) {
  return request({
    url: '/admin/luckin/products',
    method: 'post',
    data
  })
}

// 更新产品
export function updateProduct(id, data) {
  return request({
    url: `/admin/luckin/products/${id}`,
    method: 'put',
    data
  })
}

// 删除产品
export function deleteProduct(id) {
  return request({
    url: `/admin/luckin/products/${id}`,
    method: 'delete'
  })
}

// 获取产品类别列表
export function getProductCategories() {
  return request({
    url: '/admin/luckin/products/categories',
    method: 'get'
  })
}

// 批量导入产品
export function batchImportProducts(data) {
  return request({
    url: '/admin/luckin/products/import',
    method: 'post',
    data
  })
}

// 种类绑定相关API

// 获取种类绑定列表
export function getCategoryBindings(categoryId) {
  return request({
    url: `/admin/luckin/categories/${categoryId}/bindings`,
    method: 'get'
  })
}

// 创建绑定
export function createBinding(categoryId, data) {
  return request({
    url: `/admin/luckin/categories/${categoryId}/bindings`,
    method: 'post',
    data
  })
}

// 删除绑定
export function deleteBinding(bindingId) {
  return request({
    url: `/admin/luckin/bindings/${bindingId}`,
    method: 'delete'
  })
}

// 更新绑定优先级
export function updateBindingPriority(bindingId, priority) {
  return request({
    url: `/admin/luckin/bindings/${bindingId}/priority`,
    method: 'put',
    data: { priority }
  })
}

// 获取有效选项
export function getActiveOptions() {
  return request({
    url: '/admin/luckin/active-options',
    method: 'get'
  })
}

// 获取价格关联的商品
export function getPriceProducts(priceId) {
  return request({
    url: `/admin/luckin/prices/${priceId}/products`,
    method: 'get'
  })
}