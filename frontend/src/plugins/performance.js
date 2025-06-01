// 性能优化插件

import { componentLoader, PreloadStrategy } from '@/utils/componentLoader'

/**
 * 性能监控和优化插件
 */
export default {
  install(app, options = {}) {
    const {
      enablePreload = true,
      enableLazyLoad = true,
      enablePerformanceMonitoring = true,
      preloadStrategy = 'auto'
    } = options

    // 性能监控
    if (enablePerformanceMonitoring) {
      setupPerformanceMonitoring(app)
    }

    // 预加载策略
    if (enablePreload) {
      setupPreloadStrategy(preloadStrategy)
    }

    // 懒加载支持
    if (enableLazyLoad) {
      setupLazyLoad(app)
    }

    // 全局性能工具
    setupGlobalPerformanceTools(app)
  }
}

/**
 * 设置性能监控
 */
function setupPerformanceMonitoring(app) {
  // 组件渲染性能监控
  const renderTimes = new Map()

  app.mixin({
    beforeCreate() {
      if (this.$options.name) {
        renderTimes.set(this.$options.name, performance.now())
      }
    },
    
    mounted() {
      if (this.$options.name && renderTimes.has(this.$options.name)) {
        const startTime = renderTimes.get(this.$options.name)
        const renderTime = performance.now() - startTime
        
        // 如果渲染时间超过阈值，输出警告
        if (renderTime > 100) {
          console.warn(`Component ${this.$options.name} took ${renderTime.toFixed(2)}ms to render`)
        }
        
        renderTimes.delete(this.$options.name)
      }
    }
  })

  // 内存使用监控
  if (typeof PerformanceObserver !== 'undefined') {
    try {
      const observer = new PerformanceObserver((list) => {
        const entries = list.getEntries()
        entries.forEach((entry) => {
          if (entry.entryType === 'measure') {
            console.log(`Performance measure: ${entry.name} - ${entry.duration.toFixed(2)}ms`)
          }
        })
      })
      observer.observe({ entryTypes: ['measure'] })
    } catch (error) {
      console.warn('PerformanceObserver not supported:', error)
    }
  }

  // 页面性能指标收集
  window.addEventListener('load', () => {
    setTimeout(() => {
      const perfData = performance.getEntriesByType('navigation')[0]
      if (perfData) {
        const metrics = {
          dns: perfData.domainLookupEnd - perfData.domainLookupStart,
          tcp: perfData.connectEnd - perfData.connectStart,
          ttfb: perfData.responseStart - perfData.requestStart,
          download: perfData.responseEnd - perfData.responseStart,
          domReady: perfData.domContentLoadedEventEnd - perfData.fetchStart,
          total: perfData.loadEventEnd - perfData.fetchStart
        }
        
        console.log('Page Performance Metrics:', metrics)
        
        // 可以发送到分析服务
        // analytics.track('page_performance', metrics)
      }
    }, 1000)
  })
}

/**
 * 设置预加载策略
 */
function setupPreloadStrategy(strategy) {
  const executePreload = async () => {
    try {
      // 预加载关键组件
      if (strategy === 'auto' || strategy === 'critical') {
        await Promise.allSettled(
          PreloadStrategy.critical.map(loader => loader())
        )
        console.log('Critical components preloaded')
      }

      // 在空闲时间预加载常用组件
      if (strategy === 'auto' || strategy === 'common') {
        const preloadCommon = async () => {
          await Promise.allSettled(
            PreloadStrategy.common.map(loader => loader())
          )
          console.log('Common components preloaded')
        }

        if ('requestIdleCallback' in window) {
          requestIdleCallback(preloadCommon)
        } else {
          setTimeout(preloadCommon, 2000)
        }
      }

      // 根据网络状况调整预加载策略
      if (navigator.connection) {
        const connection = navigator.connection
        const isSlowNetwork = connection.effectiveType === 'slow-2g' || 
                             connection.effectiveType === '2g'
        
        if (!isSlowNetwork && (strategy === 'auto' || strategy === 'aggressive')) {
          // 在快速网络下预加载更多资源
          setTimeout(async () => {
            await Promise.allSettled([
              PreloadStrategy.features.charts(),
              PreloadStrategy.features.export(),
              PreloadStrategy.features.datetime()
            ])
            console.log('Feature components preloaded')
          }, 3000)
        }
      }
    } catch (error) {
      console.error('Preload failed:', error)
    }
  }

  // 页面加载完成后开始预加载
  if (document.readyState === 'complete') {
    executePreload()
  } else {
    window.addEventListener('load', executePreload)
  }
}

/**
 * 设置懒加载支持
 */
function setupLazyLoad(app) {
  // 注册懒加载指令
  app.directive('lazy-load', {
    mounted(el, binding) {
      const { value: loadFunction, modifiers } = binding
      const threshold = modifiers.threshold || 0.1
      
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              if (typeof loadFunction === 'function') {
                loadFunction()
              }
              observer.unobserve(el)
            }
          })
        },
        { threshold }
      )
      
      observer.observe(el)
      el._lazyObserver = observer
    },
    
    unmounted(el) {
      if (el._lazyObserver) {
        el._lazyObserver.disconnect()
        delete el._lazyObserver
      }
    }
  })

  // 懒加载图片指令
  app.directive('lazy-img', {
    mounted(el, binding) {
      const { value: src } = binding
      
      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              const img = entry.target
              
              // 创建新的图片元素来预加载
              const imageLoader = new Image()
              imageLoader.onload = () => {
                img.src = src
                img.classList.add('loaded')
              }
              imageLoader.onerror = () => {
                img.classList.add('error')
              }
              imageLoader.src = src
              
              observer.unobserve(img)
            }
          })
        },
        { threshold: 0.1 }
      )
      
      observer.observe(el)
      el._imgObserver = observer
    },
    
    unmounted(el) {
      if (el._imgObserver) {
        el._imgObserver.disconnect()
        delete el._imgObserver
      }
    }
  })
}

/**
 * 设置全局性能工具
 */
function setupGlobalPerformanceTools(app) {
  // 全局性能工具方法
  const performanceTools = {
    // 测量代码执行时间
    measure(name, fn) {
      const start = performance.now()
      const result = fn()
      const end = performance.now()
      console.log(`${name}: ${(end - start).toFixed(2)}ms`)
      return result
    },

    // 异步测量
    async measureAsync(name, fn) {
      const start = performance.now()
      const result = await fn()
      const end = performance.now()
      console.log(`${name}: ${(end - start).toFixed(2)}ms`)
      return result
    },

    // 内存使用情况
    getMemoryUsage() {
      if (performance.memory) {
        return {
          used: Math.round(performance.memory.usedJSHeapSize / 1024 / 1024) + ' MB',
          total: Math.round(performance.memory.totalJSHeapSize / 1024 / 1024) + ' MB',
          limit: Math.round(performance.memory.jsHeapSizeLimit / 1024 / 1024) + ' MB'
        }
      }
      return null
    },

    // 检查网络状况
    getNetworkInfo() {
      if (navigator.connection) {
        return {
          effectiveType: navigator.connection.effectiveType,
          downlink: navigator.connection.downlink,
          rtt: navigator.connection.rtt,
          saveData: navigator.connection.saveData
        }
      }
      return null
    },

    // 性能建议
    getPerformanceAdvice() {
      const advice = []
      
      // 检查内存使用
      const memory = this.getMemoryUsage()
      if (memory) {
        const usedMB = parseInt(memory.used)
        if (usedMB > 100) {
          advice.push('内存使用较高，建议清理不必要的数据')
        }
      }
      
      // 检查网络状况
      const network = this.getNetworkInfo()
      if (network) {
        if (network.effectiveType === 'slow-2g' || network.effectiveType === '2g') {
          advice.push('网络速度较慢，建议减少资源加载')
        }
        if (network.saveData) {
          advice.push('用户启用了数据节省模式')
        }
      }
      
      return advice
    },

    // 组件加载器统计
    getComponentStats() {
      return componentLoader.getStats()
    }
  }

  // 添加到全局属性
  app.config.globalProperties.$performance = performanceTools
  
  // 在开发环境下添加到 window 对象
  if (process.env.NODE_ENV === 'development') {
    window.__PERFORMANCE_TOOLS__ = performanceTools
    console.log('Performance tools available at window.__PERFORMANCE_TOOLS__')
  }
}

/**
 * 性能优化建议工具
 */
export class PerformanceAdvisor {
  constructor() {
    this.recommendations = []
    this.isMonitoring = false
  }

  startMonitoring() {
    if (this.isMonitoring) return
    
    this.isMonitoring = true
    
    // 监控长任务
    if ('PerformanceObserver' in window) {
      try {
        const observer = new PerformanceObserver((list) => {
          const entries = list.getEntries()
          entries.forEach((entry) => {
            if (entry.duration > 50) {
              this.recommendations.push({
                type: 'long-task',
                message: `检测到长任务: ${entry.duration.toFixed(2)}ms`,
                suggestion: '考虑将长任务分解为较小的任务',
                timestamp: Date.now()
              })
            }
          })
        })
        observer.observe({ entryTypes: ['longtask'] })
      } catch (error) {
        console.warn('Long task monitoring not supported:', error)
      }
    }

    // 监控 FPS
    this.monitorFPS()
  }

  monitorFPS() {
    let frames = 0
    let lastTime = performance.now()
    
    const measureFPS = () => {
      frames++
      const currentTime = performance.now()
      
      if (currentTime >= lastTime + 1000) {
        const fps = Math.round((frames * 1000) / (currentTime - lastTime))
        
        if (fps < 30) {
          this.recommendations.push({
            type: 'low-fps',
            message: `FPS 较低: ${fps}`,
            suggestion: '考虑优化动画或减少 DOM 操作',
            timestamp: Date.now()
          })
        }
        
        frames = 0
        lastTime = currentTime
      }
      
      if (this.isMonitoring) {
        requestAnimationFrame(measureFPS)
      }
    }
    
    requestAnimationFrame(measureFPS)
  }

  getRecommendations() {
    return this.recommendations.slice(-10) // 返回最近10条建议
  }

  clearRecommendations() {
    this.recommendations = []
  }

  stopMonitoring() {
    this.isMonitoring = false
  }
}

// 创建全局性能顾问实例
export const performanceAdvisor = new PerformanceAdvisor()

/**
 * 高级性能优化工具
 */
export class AdvancedPerformanceManager {
  constructor() {
    this.metrics = new Map()
    this.observers = []
    this.config = {
      enableVirtualScrolling: true,
      enableImageOptimization: true,
      enableCodeSplitting: true,
      enablePrefetching: true
    }
  }

  // 虚拟滚动优化
  setupVirtualScrolling(container, itemHeight, visibleCount) {
    let isScrolling = false
    let scrollTimeout

    const handleScroll = () => {
      if (!isScrolling) {
        requestAnimationFrame(() => {
          // 计算可见范围
          const scrollTop = container.scrollTop
          const startIndex = Math.floor(scrollTop / itemHeight)
          const endIndex = Math.min(startIndex + visibleCount, container.children.length)

          // 只渲染可见区域的元素
          Array.from(container.children).forEach((child, index) => {
            const shouldShow = index >= startIndex && index <= endIndex
            child.style.display = shouldShow ? '' : 'none'
          })

          isScrolling = false
        })
      }
      isScrolling = true

      // 防抖处理
      clearTimeout(scrollTimeout)
      scrollTimeout = setTimeout(() => {
        isScrolling = false
      }, 100)
    }

    container.addEventListener('scroll', handleScroll, { passive: true })
    return () => container.removeEventListener('scroll', handleScroll)
  }

  // 图片优化
  optimizeImages() {
    const images = document.querySelectorAll('img')
    images.forEach(img => {
      // 懒加载
      if (!img.hasAttribute('loading')) {
        img.setAttribute('loading', 'lazy')
      }

      // 响应式图片
      if (!img.hasAttribute('sizes') && !img.closest('picture')) {
        img.setAttribute('sizes', '(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw')
      }

      // 错误处理
      img.addEventListener('error', () => {
        img.style.display = 'none'
      }, { once: true })
    })
  }

  // 资源预取
  prefetchResource(url, type = 'fetch') {
    const link = document.createElement('link')
    link.rel = 'prefetch'
    link.href = url
    if (type !== 'fetch') {
      link.as = type
    }
    document.head.appendChild(link)
  }

  // 关键资源预加载
  preloadCriticalResources(resources) {
    resources.forEach(({ url, type, crossorigin }) => {
      const link = document.createElement('link')
      link.rel = 'preload'
      link.href = url
      link.as = type
      if (crossorigin) {
        link.crossOrigin = crossorigin
      }
      document.head.appendChild(link)
    })
  }

  // Web Workers 支持
  createWebWorker(workerScript) {
    return new Promise((resolve, reject) => {
      try {
        const worker = new Worker(workerScript)
        resolve(worker)
      } catch (error) {
        reject(error)
      }
    })
  }

  // Service Worker 注册
  registerServiceWorker(swPath) {
    if ('serviceWorker' in navigator) {
      return navigator.serviceWorker.register(swPath)
        .then(registration => {
          console.log('Service Worker registered:', registration)
          return registration
        })
        .catch(error => {
          console.error('Service Worker registration failed:', error)
          throw error
        })
    }
    return Promise.reject(new Error('Service Worker not supported'))
  }

  // 内存泄漏检测
  detectMemoryLeaks() {
    if (!performance.memory) return null

    const initialMemory = performance.memory.usedJSHeapSize
    
    return {
      start: () => {
        this.metrics.set('memoryStart', performance.memory.usedJSHeapSize)
      },
      
      check: () => {
        const current = performance.memory.usedJSHeapSize
        const start = this.metrics.get('memoryStart') || initialMemory
        const increase = current - start
        
        return {
          current: Math.round(current / 1024 / 1024),
          increase: Math.round(increase / 1024 / 1024),
          isLeak: increase > 10 * 1024 * 1024 // 10MB 增长可能是泄漏
        }
      }
    }
  }

  // 网络状态适配
  adaptToNetworkConditions() {
    if (!navigator.connection) return

    const connection = navigator.connection
    const updateStrategy = () => {
      switch (connection.effectiveType) {
        case 'slow-2g':
        case '2g':
          this.config.enablePrefetching = false
          this.config.enableImageOptimization = true
          break
        case '3g':
          this.config.enablePrefetching = true
          this.config.enableImageOptimization = true
          break
        case '4g':
          this.config.enablePrefetching = true
          this.config.enableImageOptimization = false
          break
      }
    }

    updateStrategy()
    connection.addEventListener('change', updateStrategy)
  }

  // 性能预算检查
  checkPerformanceBudget(budget) {
    const navigationEntry = performance.getEntriesByType('navigation')[0]
    if (!navigationEntry) return null

    const metrics = {
      fcp: this.getFCP(),
      lcp: this.getLCP(),
      cls: this.getCLS(),
      fid: this.getFID(),
      ttfb: navigationEntry.responseStart - navigationEntry.requestStart
    }

    const violations = []
    Object.entries(budget).forEach(([metric, limit]) => {
      if (metrics[metric] > limit) {
        violations.push({
          metric,
          actual: metrics[metric],
          budget: limit,
          over: metrics[metric] - limit
        })
      }
    })

    return { metrics, violations }
  }

  // Core Web Vitals 测量
  getFCP() {
    const fcpEntry = performance.getEntriesByName('first-contentful-paint')[0]
    return fcpEntry ? fcpEntry.startTime : null
  }

  getLCP() {
    return new Promise(resolve => {
      const observer = new PerformanceObserver(list => {
        const entries = list.getEntries()
        const lastEntry = entries[entries.length - 1]
        resolve(lastEntry.startTime)
        observer.disconnect()
      })
      observer.observe({ entryTypes: ['largest-contentful-paint'] })
    })
  }

  getCLS() {
    return new Promise(resolve => {
      let clsValue = 0
      const observer = new PerformanceObserver(list => {
        for (const entry of list.getEntries()) {
          if (!entry.hadRecentInput) {
            clsValue += entry.value
          }
        }
        resolve(clsValue)
      })
      observer.observe({ entryTypes: ['layout-shift'] })
      
      // 5秒后停止观察
      setTimeout(() => {
        observer.disconnect()
        resolve(clsValue)
      }, 5000)
    })
  }

  getFID() {
    return new Promise(resolve => {
      const observer = new PerformanceObserver(list => {
        for (const entry of list.getEntries()) {
          resolve(entry.processingStart - entry.startTime)
          observer.disconnect()
          break
        }
      })
      observer.observe({ entryTypes: ['first-input'] })
    })
  }

  // 生成性能报告
  generatePerformanceReport() {
    return {
      timestamp: Date.now(),
      url: location.href,
      userAgent: navigator.userAgent,
      connection: navigator.connection ? {
        effectiveType: navigator.connection.effectiveType,
        downlink: navigator.connection.downlink,
        rtt: navigator.connection.rtt
      } : null,
      memory: performance.memory ? {
        used: Math.round(performance.memory.usedJSHeapSize / 1024 / 1024),
        total: Math.round(performance.memory.totalJSHeapSize / 1024 / 1024),
        limit: Math.round(performance.memory.jsHeapSizeLimit / 1024 / 1024)
      } : null,
      recommendations: performanceAdvisor.getRecommendations()
    }
  }

  // 性能优化建议
  getOptimizationSuggestions() {
    const suggestions = []

    // 检查资源大小
    const resources = performance.getEntriesByType('resource')
    const largeResources = resources.filter(r => r.transferSize > 1024 * 1024) // > 1MB

    if (largeResources.length > 0) {
      suggestions.push({
        type: 'resource-size',
        message: `发现 ${largeResources.length} 个大型资源`,
        suggestion: '考虑压缩或拆分大型资源',
        resources: largeResources.map(r => ({ name: r.name, size: r.transferSize }))
      })
    }

    // 检查未使用的CSS
    if (document.styleSheets.length > 5) {
      suggestions.push({
        type: 'css-optimization',
        message: '加载了较多样式表',
        suggestion: '考虑合并或按需加载CSS'
      })
    }

    // 检查JavaScript执行时间
    const scriptResources = resources.filter(r => r.name.includes('.js'))
    const slowScripts = scriptResources.filter(r => r.duration > 100)

    if (slowScripts.length > 0) {
      suggestions.push({
        type: 'script-performance',
        message: `发现 ${slowScripts.length} 个执行缓慢的脚本`,
        suggestion: '考虑代码分割或优化脚本'
      })
    }

    return suggestions
  }
}

// 创建高级性能管理器实例
export const advancedPerformanceManager = new AdvancedPerformanceManager()