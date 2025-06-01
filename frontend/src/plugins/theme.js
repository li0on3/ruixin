// 主题系统插件
import { reactive, watch } from 'vue'

// 主题配置
const themeConfig = {
  // 默认主题设置
  default: {
    name: 'light',
    displayName: '浅色主题',
    colors: {
      primary: '#f97316',
      secondary: '#64748b',
      success: '#10b981',
      warning: '#f59e0b',
      error: '#ef4444',
      info: '#3b82f6'
    }
  },
  
  // 可用主题列表
  themes: {
    light: {
      name: 'light',
      displayName: '浅色主题',
      colors: {
        primary: '#f97316',
        secondary: '#64748b',
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#3b82f6'
      }
    },
    
    dark: {
      name: 'dark',
      displayName: '深色主题',
      colors: {
        primary: '#fb923c',
        secondary: '#94a3b8',
        success: '#34d399',
        warning: '#fbbf24',
        error: '#f87171',
        info: '#60a5fa'
      }
    },
    
    auto: {
      name: 'auto',
      displayName: '跟随系统',
      colors: {
        primary: '#f97316',
        secondary: '#64748b',
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#3b82f6'
      }
    }
  }
}

// 创建主题状态
const themeState = reactive({
  current: 'light',
  isDark: false,
  isAuto: false,
  colors: { ...themeConfig.themes.light.colors }
})

// 主题管理器类
class ThemeManager {
  constructor() {
    this.state = themeState
    this.mediaQuery = null
    this.init()
  }
  
  // 初始化主题系统
  init() {
    // 从本地存储恢复主题设置
    this.restoreTheme()
    
    // 设置系统主题监听器
    this.setupMediaQuery()
    
    // 监听主题变化
    this.watchThemeChange()
    
    // 应用当前主题
    this.applyTheme(this.state.current)
  }
  
  // 从本地存储恢复主题
  restoreTheme() {
    try {
      const savedTheme = localStorage.getItem('app-theme')
      if (savedTheme && themeConfig.themes[savedTheme]) {
        this.state.current = savedTheme
      }
    } catch (error) {
      console.warn('Failed to restore theme from localStorage:', error)
    }
  }
  
  // 设置系统主题监听器
  setupMediaQuery() {
    this.mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    this.mediaQuery.addEventListener('change', this.handleSystemThemeChange.bind(this))
  }
  
  // 处理系统主题变化
  handleSystemThemeChange(e) {
    if (this.state.current === 'auto') {
      this.applySystemTheme(e.matches)
    }
  }
  
  // 监听主题状态变化
  watchThemeChange() {
    watch(
      () => this.state.current,
      (newTheme) => {
        this.applyTheme(newTheme)
        this.saveTheme(newTheme)
      }
    )
  }
  
  // 应用主题
  applyTheme(themeName) {
    const theme = themeConfig.themes[themeName]
    if (!theme) return
    
    if (themeName === 'auto') {
      this.state.isAuto = true
      this.applySystemTheme(this.mediaQuery?.matches || false)
    } else {
      this.state.isAuto = false
      this.state.isDark = themeName === 'dark'
      this.setThemeColors(theme.colors)
      this.setDocumentTheme(themeName)
    }
  }
  
  // 应用系统主题
  applySystemTheme(isDark) {
    this.state.isDark = isDark
    const theme = isDark ? themeConfig.themes.dark : themeConfig.themes.light
    this.setThemeColors(theme.colors)
    this.setDocumentTheme(isDark ? 'dark' : 'light')
  }
  
  // 设置主题颜色
  setThemeColors(colors) {
    this.state.colors = { ...colors }
    
    // 更新 CSS 自定义属性
    const root = document.documentElement
    Object.entries(colors).forEach(([key, value]) => {
      root.style.setProperty(`--theme-${key}`, value)
      
      // 转换为 RGB 值（用于透明度计算）
      const rgb = this.hexToRgb(value)
      if (rgb) {
        root.style.setProperty(`--theme-${key}-rgb`, `${rgb.r}, ${rgb.g}, ${rgb.b}`)
      }
    })
  }
  
  // 设置文档主题属性
  setDocumentTheme(themeName) {
    // 同时设置在 html 和 body 上，确保所有样式都能生效
    document.documentElement.setAttribute('data-theme', themeName)
    document.body.setAttribute('data-theme', themeName)
    
    // 设置 class
    document.documentElement.className = document.documentElement.className
      .replace(/theme-\w+/g, '')
      .trim() + ` theme-${themeName}`
    
    // 同时在 body 上设置 class
    document.body.className = document.body.className
      .replace(/theme-\w+/g, '')
      .trim() + ` theme-${themeName}`
  }
  
  // 保存主题到本地存储
  saveTheme(themeName) {
    try {
      localStorage.setItem('app-theme', themeName)
    } catch (error) {
      console.warn('Failed to save theme to localStorage:', error)
    }
  }
  
  // 切换主题
  setTheme(themeName) {
    if (themeConfig.themes[themeName]) {
      this.state.current = themeName
    }
  }
  
  // 切换到下一个主题
  toggleTheme() {
    const themes = Object.keys(themeConfig.themes)
    const currentIndex = themes.indexOf(this.state.current)
    const nextIndex = (currentIndex + 1) % themes.length
    this.setTheme(themes[nextIndex])
  }
  
  // 获取当前主题信息
  getCurrentTheme() {
    return {
      ...themeConfig.themes[this.state.current],
      isDark: this.state.isDark,
      isAuto: this.state.isAuto
    }
  }
  
  // 获取所有可用主题
  getAvailableThemes() {
    return Object.values(themeConfig.themes)
  }
  
  // 检查是否为深色主题
  isDarkTheme() {
    return this.state.isDark
  }
  
  // 检查是否为自动主题
  isAutoTheme() {
    return this.state.isAuto
  }
  
  // 获取主题颜色
  getThemeColor(colorName) {
    return this.state.colors[colorName]
  }
  
  // 工具方法：hex 转 RGB
  hexToRgb(hex) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
    return result ? {
      r: parseInt(result[1], 16),
      g: parseInt(result[2], 16),
      b: parseInt(result[3], 16)
    } : null
  }
  
  // 工具方法：生成颜色变体
  generateColorVariants(baseColor) {
    // 这里可以实现颜色变体生成算法
    // 例如生成不同亮度的颜色
    const variants = {}
    const weights = [50, 100, 200, 300, 400, 500, 600, 700, 800, 900]
    
    weights.forEach(weight => {
      variants[weight] = this.adjustColorBrightness(baseColor, weight)
    })
    
    return variants
  }
  
  // 工具方法：调整颜色亮度
  adjustColorBrightness(color, weight) {
    // 简化实现，实际项目中可以使用更复杂的颜色处理库
    const factor = (500 - weight) / 500
    const rgb = this.hexToRgb(color)
    
    if (!rgb) return color
    
    const adjustedRgb = {
      r: Math.round(rgb.r + (255 - rgb.r) * factor * 0.5),
      g: Math.round(rgb.g + (255 - rgb.g) * factor * 0.5),
      b: Math.round(rgb.b + (255 - rgb.b) * factor * 0.5)
    }
    
    return `rgb(${adjustedRgb.r}, ${adjustedRgb.g}, ${adjustedRgb.b})`
  }
  
  // 销毁主题管理器
  destroy() {
    if (this.mediaQuery) {
      this.mediaQuery.removeEventListener('change', this.handleSystemThemeChange)
    }
  }
}

// 创建全局主题管理器实例
const themeManager = new ThemeManager()

// Vue 插件
export default {
  install(app) {
    // 注入主题管理器
    app.config.globalProperties.$theme = themeManager
    app.provide('theme', themeManager)
    
    // 注册全局组件或指令（如果需要）
    
    // 开发环境下暴露到 window 对象
    if (process.env.NODE_ENV === 'development') {
      window.__THEME_MANAGER__ = themeManager
    }
  }
}

// 导出主题管理器和状态
export { themeManager, themeState, themeConfig }

// 组合式 API
export function useTheme() {
  return {
    themeState,
    themeManager,
    setTheme: themeManager.setTheme.bind(themeManager),
    toggleTheme: themeManager.toggleTheme.bind(themeManager),
    getCurrentTheme: themeManager.getCurrentTheme.bind(themeManager),
    getAvailableThemes: themeManager.getAvailableThemes.bind(themeManager),
    isDarkTheme: themeManager.isDarkTheme.bind(themeManager),
    isAutoTheme: themeManager.isAutoTheme.bind(themeManager),
    getThemeColor: themeManager.getThemeColor.bind(themeManager)
  }
}