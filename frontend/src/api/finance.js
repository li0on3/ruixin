import request from './request'

// 财务管理相关API

// 充值
export function recharge(data) {
  return request({
    url: '/admin/finance/recharge',
    method: 'post',
    data
  })
}

// 余额调整
export function adjustBalance(data) {
  return request({
    url: '/admin/finance/adjust',
    method: 'post',
    data
  })
}

// 获取交易记录列表
export function getTransactionList(params) {
  return request({
    url: '/admin/finance/transactions',
    method: 'get',
    params
  })
}

// 获取提现申请列表
export function getWithdrawalList(params) {
  return request({
    url: '/admin/finance/withdrawals',
    method: 'get',
    params
  })
}

// 获取待处理的提现申请
export function getPendingWithdrawals(params) {
  return request({
    url: '/admin/finance/withdrawals/pending',
    method: 'get',
    params
  })
}

// 处理提现申请
export function processWithdrawal(id, data) {
  return request({
    url: `/admin/finance/withdrawals/${id}/process`,
    method: 'post',
    data
  })
}

// 获取财务统计
export function getFinanceStatistics() {
  return request({
    url: '/admin/finance/statistics',
    method: 'get'
  })
}

// 获取分销商余额
export function getDistributorBalance(id) {
  return request({
    url: `/admin/finance/distributors/${id}/balance`,
    method: 'get'
  })
}

// 交易类型枚举
export const TransactionType = {
  1: '充值',
  2: '消费',
  3: '提现',
  4: '退款',
  5: '调整'
}

// 提现状态枚举
export const WithdrawalStatus = {
  0: '待处理',
  1: '已处理',
  2: '已拒绝'
}