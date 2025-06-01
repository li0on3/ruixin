# ECharts 深色主题集成任务

## 任务描述
在 Dashboard 和统计页面中应用 ECharts 深色主题，确保图表在主题切换时能自动更新。

## 执行时间
2025-01-06

## 任务目标
1. 在 DashboardEnhanced.vue 中集成深色主题
2. 在 StatisticsEnhanced.vue 中集成深色主题
3. 确保所有图表响应主题切换
4. 图表颜色和样式在深色主题下正确显示

## 实施方案
使用已有的 echarts-dark-theme.js 和主题管理系统，在组件中直接添加主题支持。

## 修改的文件

### 1. DashboardEnhanced.vue
- 导入深色主题相关依赖
- 注册 ECharts 深色主题
- 修改所有图表初始化，添加主题参数
- 添加主题变化监听器，动态切换图表主题

### 2. StatisticsEnhanced.vue
- 导入深色主题相关依赖
- 注册 ECharts 深色主题
- 修改所有图表初始化，添加主题参数
- 添加主题变化监听器，动态切换图表主题
- 更新颜色映射函数，支持动态主题颜色

## 技术实现要点

### 1. 主题注册
```javascript
import { darkTheme, applyDarkTheme, getEChartsTheme } from '@/utils/echarts-dark-theme'
import { useTheme } from '@/plugins/theme'

const { isDarkTheme } = useTheme()
applyDarkTheme(echarts)
```

### 2. 图表初始化
```javascript
// 初始化时传入主题
const chart = echarts.init(element, getEChartsTheme())
```

### 3. 主题切换监听
```javascript
watch(() => isDarkTheme(), () => {
  // 销毁旧图表
  chart?.dispose()
  // 重新初始化
  chart = echarts.init(element, getEChartsTheme())
  // 重新渲染
  renderChart()
})
```

## 测试要点
1. 浅色主题下图表显示正常
2. 深色主题下图表颜色和样式正确
3. 主题切换时图表自动更新
4. 没有内存泄漏（图表实例正确销毁）

## 完成状态
✅ 已完成

## 后续优化建议
1. 可以考虑创建统一的 ECharts 钩子函数来管理图表实例
2. 可以添加主题切换动画效果
3. 可以为不同的图表类型定制专门的深色主题配色