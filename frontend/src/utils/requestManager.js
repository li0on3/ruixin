import axios from 'axios'

// 请求管理器 - 防止重复请求
class RequestManager {
  constructor() {
    this.pending = new Map()
  }

  // 生成请求的唯一key
  generateKey(config) {
    const { method, url, params, data } = config
    return `${method}_${url}_${JSON.stringify(params || {})}_${JSON.stringify(data || {})}`
  }

  // 添加请求
  add(config) {
    const key = this.generateKey(config)
    
    // 如果已存在相同请求，取消之前的请求
    if (this.pending.has(key)) {
      const cancel = this.pending.get(key)
      cancel('取消重复请求')
    }
    
    // 创建新的取消令牌
    const source = axios.CancelToken.source()
    config.cancelToken = source.token
    
    // 保存取消函数
    this.pending.set(key, source.cancel)
    
    return config
  }

  // 移除请求
  remove(config) {
    const key = this.generateKey(config)
    this.pending.delete(key)
  }

  // 清空所有请求
  clear() {
    // 取消所有pending请求
    for (const cancel of this.pending.values()) {
      cancel('清空所有请求')
    }
    this.pending.clear()
  }
}

export default new RequestManager()