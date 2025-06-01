// 路由工具函数
import { defineAsyncComponent } from 'vue'

/**
 * 创建懒加载路由组件（直接返回导入函数，避免 defineAsyncComponent 警告）
 * @param {Function} importFn - 动态导入函数
 * @param {Object} options - 配置选项
 * @returns {Function} 懒加载组件
 */
export function createLazyComponent(importFn, options = {}) {
  // 直接返回导入函数，Vue Router 会自动处理
  return importFn
}

/**
 * 预加载路由组件
 * @param {Array} routes - 路由配置数组
 */
export function preloadRouteComponents(routes) {
  const preloadPromises = []
  
  function collectComponents(routeList) {
    routeList.forEach(route => {
      if (route.component) {
        // 处理 defineAsyncComponent 创建的组件
        if (route.component.__asyncLoader) {
          // Vue 3 异步组件有 __asyncLoader 属性
          const preloadPromise = () => {
            try {
              return route.component.__asyncLoader()
            } catch (error) {
              console.warn(`Failed to preload component for route: ${route.path}`, error)
              return Promise.resolve()
            }
          }
          preloadPromises.push(preloadPromise)
        } else if (typeof route.component === 'function') {
          // 处理原始的动态导入函数
          const preloadPromise = () => {
            try {
              return route.component()
            } catch (error) {
              console.warn(`Failed to preload component for route: ${route.path}`, error)
              return Promise.resolve()
            }
          }
          preloadPromises.push(preloadPromise)
        }
      }
      
      if (route.children) {
        collectComponents(route.children)
      }
    })
  }
  
  collectComponents(routes)
  
  // 返回预加载函数
  return {
    // 预加载关键路由（用户最可能访问的路由）
    preloadCritical: async () => {
      const criticalPromises = preloadPromises.slice(0, 3) // 预加载前3个组件
      await Promise.allSettled(criticalPromises.map(fn => fn()))
    },
    
    // 预加载所有路由（在空闲时间执行）
    preloadAll: async () => {
      // 使用 requestIdleCallback 在浏览器空闲时预加载
      if ('requestIdleCallback' in window) {
        return new Promise(resolve => {
          window.requestIdleCallback(async () => {
            await Promise.allSettled(preloadPromises.map(fn => fn()))
            resolve()
          })
        })
      } else {
        // 降级方案：使用 setTimeout
        return new Promise(resolve => {
          setTimeout(async () => {
            await Promise.allSettled(preloadPromises.map(fn => fn()))
            resolve()
          }, 1000)
        })
      }
    }
  }
}

/**
 * 路由缓存管理
 */
export class RouteCache {
  constructor(maxSize = 10) {
    this.cache = new Map()
    this.maxSize = maxSize
  }
  
  get(key) {
    if (this.cache.has(key)) {
      // 更新访问时间
      const item = this.cache.get(key)
      this.cache.delete(key)
      this.cache.set(key, { ...item, lastAccess: Date.now() })
      return item.data
    }
    return null
  }
  
  set(key, data) {
    if (this.cache.size >= this.maxSize) {
      // 删除最旧的项目
      const oldestKey = this.cache.keys().next().value
      this.cache.delete(oldestKey)
    }
    
    this.cache.set(key, {
      data,
      lastAccess: Date.now()
    })
  }
  
  clear() {
    this.cache.clear()
  }
  
  // 获取缓存统计信息
  getStats() {
    return {
      size: this.cache.size,
      maxSize: this.maxSize,
      keys: Array.from(this.cache.keys())
    }
  }
}

/**
 * 性能监控工具
 */
export class RoutePerformanceMonitor {
  constructor() {
    this.metrics = new Map()
  }
  
  // 开始监控路由切换
  startRouteChange(from, to) {
    const key = `${from?.path || 'initial'} -> ${to.path}`
    this.metrics.set(key, {
      startTime: performance.now(),
      from: from?.path,
      to: to.path
    })
  }
  
  // 结束监控路由切换
  endRouteChange(to) {
    const entries = Array.from(this.metrics.entries())
    const currentEntry = entries.find(([key, value]) => 
      key.endsWith(`-> ${to.path}`) && !value.endTime
    )
    
    if (currentEntry) {
      const [key, value] = currentEntry
      value.endTime = performance.now()
      value.duration = value.endTime - value.startTime
      
      // 如果切换时间过长，输出警告
      if (value.duration > 1000) {
        console.warn(`Slow route change detected: ${key} took ${value.duration.toFixed(2)}ms`)
      }
    }
  }
  
  // 获取性能报告
  getReport() {
    const entries = Array.from(this.metrics.values())
    const completedEntries = entries.filter(entry => entry.endTime)
    
    if (completedEntries.length === 0) {
      return { message: 'No completed route changes to report' }
    }
    
    const durations = completedEntries.map(entry => entry.duration)
    const avgDuration = durations.reduce((a, b) => a + b, 0) / durations.length
    const maxDuration = Math.max(...durations)
    const minDuration = Math.min(...durations)
    
    return {
      totalRouteChanges: completedEntries.length,
      averageDuration: avgDuration.toFixed(2) + 'ms',
      maxDuration: maxDuration.toFixed(2) + 'ms',
      minDuration: minDuration.toFixed(2) + 'ms',
      slowRoutes: completedEntries
        .filter(entry => entry.duration > 500)
        .map(entry => ({
          route: `${entry.from} -> ${entry.to}`,
          duration: entry.duration.toFixed(2) + 'ms'
        }))
    }
  }
  
  // 清除监控数据
  clear() {
    this.metrics.clear()
  }
}

/**
 * 路由预取策略
 */
export class RoutePrefetcher {
  constructor() {
    this.prefetchQueue = []
    this.prefetched = new Set()
    this.isRunning = false
  }
  
  // 添加预取任务
  addToPrefetchQueue(importFn, routePath) {
    if (!this.prefetched.has(routePath)) {
      this.prefetchQueue.push({ importFn, routePath })
      this.processPrefetchQueue()
    }
  }
  
  // 处理预取队列
  async processPrefetchQueue() {
    if (this.isRunning || this.prefetchQueue.length === 0) {
      return
    }
    
    this.isRunning = true
    
    while (this.prefetchQueue.length > 0) {
      const { importFn, routePath } = this.prefetchQueue.shift()
      
      try {
        // 检查网络状态
        if (navigator.connection && navigator.connection.effectiveType === 'slow-2g') {
          console.log('Skipping prefetch due to slow network')
          break
        }
        
        await importFn()
        this.prefetched.add(routePath)
        console.log(`Prefetched route: ${routePath}`)
        
        // 在每个预取之间添加短暂延迟，避免阻塞主线程
        await new Promise(resolve => setTimeout(resolve, 100))
        
      } catch (error) {
        console.warn(`Failed to prefetch route ${routePath}:`, error)
      }
    }
    
    this.isRunning = false
  }
  
  // 基于用户行为预测需要预取的路由
  predictAndPrefetch(currentRoute, routeConfig) {
    const predictions = this.getPredictions(currentRoute, routeConfig)
    
    predictions.forEach(({ importFn, path }) => {
      this.addToPrefetchQueue(importFn, path)
    })
  }
  
  // 获取路由预测
  getPredictions(currentRoute, routeConfig) {
    const predictions = []
    
    // 基于当前路由预测用户可能访问的路由
    const routePredictions = {
      '/dashboard': ['/cards', '/distributors', '/orders'],
      '/cards': ['/cards/batch-import', '/cards/batches', '/cards/bindings'],
      '/distributors': ['/orders', '/finance/transactions'],
      '/orders': ['/statistics', '/distributors'],
      '/statistics': ['/orders', '/cards']
    }
    
    const likelyRoutes = routePredictions[currentRoute] || []
    
    likelyRoutes.forEach(routePath => {
      const routeInfo = routeConfig.find(r => r.path === routePath)
      if (routeInfo && routeInfo.component) {
        let importFn
        
        // 处理 defineAsyncComponent 创建的组件
        if (routeInfo.component.__asyncLoader) {
          importFn = routeInfo.component.__asyncLoader
        } else if (typeof routeInfo.component === 'function') {
          // 处理原始的动态导入函数
          importFn = routeInfo.component
        }
        
        if (importFn) {
          predictions.push({
            importFn,
            path: routePath
          })
        }
      }
    })
    
    return predictions
  }
  
  // 获取预取统计
  getStats() {
    return {
      prefetched: Array.from(this.prefetched),
      queueLength: this.prefetchQueue.length,
      isRunning: this.isRunning
    }
  }
}