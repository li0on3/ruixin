// 组件按需加载工具

/**
 * 动态加载第三方组件库
 */
export class ComponentLoader {
  constructor() {
    this.loadedLibs = new Set()
    this.loadingPromises = new Map()
  }

  /**
   * 动态加载 ECharts
   * @param {Array} modules - 需要的 ECharts 模块
   * @returns {Promise} ECharts 实例
   */
  async loadECharts(modules = []) {
    const cacheKey = `echarts-${modules.join('-')}`
    
    if (this.loadingPromises.has(cacheKey)) {
      return this.loadingPromises.get(cacheKey)
    }

    const loadPromise = this._loadEChartsModules(modules)
    this.loadingPromises.set(cacheKey, loadPromise)
    
    return loadPromise
  }

  async _loadEChartsModules(modules) {
    try {
      // 动态导入 ECharts 核心
      const echarts = await import('echarts/core')
      
      // 默认组件
      const components = [
        'GridComponent',
        'TooltipComponent',
        'LegendComponent'
      ]
      
      // 默认渲染器
      const renderers = ['CanvasRenderer']
      
      // 根据需要加载的模块动态导入
      const moduleImports = await Promise.all([
        // 组件
        ...components.map(comp => import(`echarts/components`).then(m => m[comp])),
        // 图表类型
        ...modules.map(module => {
          switch (module) {
            case 'bar':
              return import('echarts/charts').then(m => m.BarChart)
            case 'line':
              return import('echarts/charts').then(m => m.LineChart)
            case 'pie':
              return import('echarts/charts').then(m => m.PieChart)
            case 'scatter':
              return import('echarts/charts').then(m => m.ScatterChart)
            case 'radar':
              return import('echarts/charts').then(m => m.RadarChart)
            case 'map':
              return import('echarts/charts').then(m => m.MapChart)
            case 'tree':
              return import('echarts/charts').then(m => m.TreeChart)
            case 'treemap':
              return import('echarts/charts').then(m => m.TreemapChart)
            case 'graph':
              return import('echarts/charts').then(m => m.GraphChart)
            case 'gauge':
              return import('echarts/charts').then(m => m.GaugeChart)
            case 'funnel':
              return import('echarts/charts').then(m => m.FunnelChart)
            case 'parallel':
              return import('echarts/charts').then(m => m.ParallelChart)
            case 'sankey':
              return import('echarts/charts').then(m => m.SankeyChart)
            case 'boxplot':
              return import('echarts/charts').then(m => m.BoxplotChart)
            case 'candlestick':
              return import('echarts/charts').then(m => m.CandlestickChart)
            case 'heatmap':
              return import('echarts/charts').then(m => m.HeatmapChart)
            case 'lines':
              return import('echarts/charts').then(m => m.LinesChart)
            case 'pictorialBar':
              return import('echarts/charts').then(m => m.PictorialBarChart)
            case 'themeRiver':
              return import('echarts/charts').then(m => m.ThemeRiverChart)
            case 'sunburst':
              return import('echarts/charts').then(m => m.SunburstChart)
            case 'custom':
              return import('echarts/charts').then(m => m.CustomChart)
            default:
              console.warn(`Unknown ECharts module: ${module}`)
              return null
          }
        }).filter(Boolean),
        // 渲染器
        ...renderers.map(renderer => import(`echarts/renderers`).then(m => m[renderer]))
      ])

      // 注册所有组件
      echarts.use(moduleImports.filter(Boolean))
      
      return echarts
    } catch (error) {
      console.error('Failed to load ECharts modules:', error)
      throw error
    }
  }

  /**
   * 动态加载 Element Plus 组件
   * @param {Array} components - 需要的组件名称
   * @returns {Promise} 组件对象
   */
  async loadElementComponents(components = []) {
    const cacheKey = `element-${components.join('-')}`
    
    if (this.loadingPromises.has(cacheKey)) {
      return this.loadingPromises.get(cacheKey)
    }

    const loadPromise = this._loadElementComponents(components)
    this.loadingPromises.set(cacheKey, loadPromise)
    
    return loadPromise
  }

  async _loadElementComponents(components) {
    const componentMap = {}
    
    for (const componentName of components) {
      try {
        // 动态导入 Element Plus 组件 - 使用新的导入方式
        const module = await import('element-plus')
        const ElComponent = module[`El${componentName}`] || module[componentName]
        if (ElComponent) {
          componentMap[componentName] = ElComponent
        } else {
          console.warn(`Element component not found: ${componentName}`)
        }
      } catch (error) {
        console.warn(`Failed to load Element component: ${componentName}`, error)
      }
    }
    
    return componentMap
  }

  /**
   * 动态加载图标组件
   * @param {Array} icons - 需要的图标名称
   * @returns {Promise} 图标组件对象
   */
  async loadIcons(icons = []) {
    const cacheKey = `icons-${icons.join('-')}`
    
    if (this.loadingPromises.has(cacheKey)) {
      return this.loadingPromises.get(cacheKey)
    }

    const loadPromise = this._loadIcons(icons)
    this.loadingPromises.set(cacheKey, loadPromise)
    
    return loadPromise
  }

  async _loadIcons(icons) {
    try {
      const iconModule = await import('@element-plus/icons-vue')
      const selectedIcons = {}
      
      icons.forEach(iconName => {
        if (iconModule[iconName]) {
          selectedIcons[iconName] = iconModule[iconName]
        } else {
          console.warn(`Icon not found: ${iconName}`)
        }
      })
      
      return selectedIcons
    } catch (error) {
      console.error('Failed to load icons:', error)
      throw error
    }
  }

  /**
   * 动态加载第三方工具库
   * @param {string} libName - 库名称
   * @param {string} version - 版本号（可选）
   * @returns {Promise} 库对象
   */
  async loadThirdPartyLib(libName, version = 'latest') {
    const cacheKey = `${libName}-${version}`
    
    if (this.loadedLibs.has(cacheKey)) {
      return window[libName] || Promise.resolve()
    }

    if (this.loadingPromises.has(cacheKey)) {
      return this.loadingPromises.get(cacheKey)
    }

    const loadPromise = this._loadThirdPartyLib(libName, version)
    this.loadingPromises.set(cacheKey, loadPromise)
    
    return loadPromise
  }

  async _loadThirdPartyLib(libName, version) {
    return new Promise((resolve, reject) => {
      // 检查库是否已经加载
      if (window[libName]) {
        resolve(window[libName])
        return
      }

      // 创建 script 标签
      const script = document.createElement('script')
      
      // 设置 CDN 地址（可以根据需要修改）
      const cdnMap = {
        'dayjs': `https://unpkg.com/dayjs@${version}/dayjs.min.js`,
        'lodash': `https://unpkg.com/lodash@${version}/lodash.min.js`,
        'moment': `https://unpkg.com/moment@${version}/min/moment.min.js`,
        'axios': `https://unpkg.com/axios@${version}/dist/axios.min.js`,
        'xlsx': `https://unpkg.com/xlsx@${version}/dist/xlsx.full.min.js`
      }

      script.src = cdnMap[libName] || `https://unpkg.com/${libName}@${version}`
      script.async = true
      
      script.onload = () => {
        this.loadedLibs.add(`${libName}-${version}`)
        resolve(window[libName])
      }
      
      script.onerror = () => {
        reject(new Error(`Failed to load library: ${libName}`))
      }
      
      document.head.appendChild(script)
    })
  }

  /**
   * 获取加载统计信息
   */
  getStats() {
    return {
      loadedLibraries: Array.from(this.loadedLibs),
      activePromises: this.loadingPromises.size,
      memoryUsage: this._getMemoryUsage()
    }
  }

  _getMemoryUsage() {
    if (performance && performance.memory) {
      return {
        used: Math.round(performance.memory.usedJSHeapSize / 1024 / 1024) + ' MB',
        total: Math.round(performance.memory.totalJSHeapSize / 1024 / 1024) + ' MB',
        limit: Math.round(performance.memory.jsHeapSizeLimit / 1024 / 1024) + ' MB'
      }
    }
    return null
  }

  /**
   * 清除缓存
   */
  clearCache() {
    this.loadingPromises.clear()
    this.loadedLibs.clear()
  }
}

// 创建全局实例
export const componentLoader = new ComponentLoader()

/**
 * Vue 组件混入：按需加载功能
 */
export const LazyLoadMixin = {
  data() {
    return {
      componentLoader: componentLoader
    }
  },
  
  methods: {
    /**
     * 懒加载 ECharts
     * @param {Array} modules - 需要的模块
     * @returns {Promise}
     */
    async loadECharts(modules = ['bar', 'line', 'pie']) {
      try {
        this.chartLoading = true
        const echarts = await this.componentLoader.loadECharts(modules)
        return echarts
      } finally {
        this.chartLoading = false
      }
    },

    /**
     * 懒加载第三方库
     * @param {string} libName - 库名称
     * @returns {Promise}
     */
    async loadLib(libName) {
      return await this.componentLoader.loadThirdPartyLib(libName)
    }
  },

  beforeUnmount() {
    // 组件销毁时可以选择清除某些缓存
    // this.componentLoader.clearCache()
  }
}

/**
 * Vue 指令：延迟加载组件
 */
export const vLazyLoad = {
  mounted(el, binding) {
    const { value: componentLoader } = binding
    
    // 使用 Intersection Observer 监听元素是否进入视口
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            // 元素进入视口时触发加载
            if (typeof componentLoader === 'function') {
              componentLoader()
            }
            observer.unobserve(el)
          }
        })
      },
      {
        threshold: 0.1 // 10% 的元素可见时触发
      }
    )
    
    observer.observe(el)
    el._lazyLoadObserver = observer
  },
  
  unmounted(el) {
    if (el._lazyLoadObserver) {
      el._lazyLoadObserver.disconnect()
      delete el._lazyLoadObserver
    }
  }
}

/**
 * 预加载策略配置
 */
export const PreloadStrategy = {
  // 关键资源立即加载
  critical: [
    () => import('@/components/LoadingAnimation.vue'),
    () => import('@/components/ResultAnimation.vue')
  ],
  
  // 常用资源在空闲时加载
  common: [
    () => import('@/components/EnhancedTable.vue'),
    () => import('@/components/EnhancedForm.vue'),
    () => import('@/components/EnhancedBreadcrumb.vue')
  ],
  
  // 特殊功能按需加载
  features: {
    charts: () => componentLoader.loadECharts(['bar', 'line', 'pie']),
    export: () => componentLoader.loadThirdPartyLib('xlsx'),
    datetime: () => componentLoader.loadThirdPartyLib('dayjs')
  }
}

// 在开发环境下提供调试工具
if (process.env.NODE_ENV === 'development') {
  window.__COMPONENT_LOADER__ = componentLoader
  console.log('Component loader available at window.__COMPONENT_LOADER__')
}