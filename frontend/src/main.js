import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
// import router from './router'
import router from './router/optimized.js' // 使用优化版路由
import animationsPlugin from '@/plugins/animations.js'
import performancePlugin from '@/plugins/performance.js'
import themePlugin from '@/plugins/theme.js'
import errorHandler from '@/plugins/error-handler.js'
// 导入新的增强样式系统
import '@/assets/styles/design-tokens.scss'
import '@/assets/styles/global-enhanced.scss'
import '@/assets/styles/animations-enhanced.scss'
// 保留原有样式作为兼容
import '@/assets/styles/global.scss'
import '@/assets/styles/modern-enhanced.scss'
import '@/assets/styles/mobile-optimized.scss'
import '@/assets/styles/message.scss'
import '@/assets/styles/theme-reorganized.scss'
import '@/assets/styles/button-fix-enhanced.scss'
import '@/assets/styles/dark-theme-unified.scss' // 统一的深色主题样式

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(themePlugin)
app.use(animationsPlugin)
app.use(performancePlugin, {
  enablePreload: true,
  enableLazyLoad: true,
  enablePerformanceMonitoring: process.env.NODE_ENV === 'development',
  preloadStrategy: 'auto'
})
app.use(ElementPlus, {
  locale: zhCn,
})
app.use(errorHandler) // 全局错误处理

app.mount('#app')