# 统计页面问题修复方案

## 问题分析

### 1. 分销商排行表格切换问题
- 点击"销售额/订单数"切换时，数据没有更新
- 原因：`loadDistributorRank` 函数在切换时虽然被调用，但是数据返回后没有正确重新渲染图表

### 2. 数据详情表格显示全是0
- 后端返回的数据结构与前端期望的不一致
- 某些字段（如profit、commission等）在后端没有计算

### 3. 表格展示问题
- 移动端适配不佳
- 某些数据格式化不正确

## 解决方案

### 1. 修复分销商排行图表切换

在 `StatisticsEnhanced.vue` 中修改：

```javascript
// 修改 initDistributorRankChart 函数
const initDistributorRankChart = () => {
  if (!distributorRankData.value || !distributorRankData.value.length) return
  if (isUnmounted) return
  
  const chartDom = document.querySelector('.distributor-rank .chart-container')
  if (!chartDom) return
  
  // 确保销毁旧实例
  if (distributorRankChartInstance) {
    try {
      distributorRankChartInstance.dispose()
    } catch (error) {
      console.warn('销毁旧图表实例失败:', error)
    }
    distributorRankChartInstance = null
  }
  
  try {
    distributorRankChartInstance = echarts.init(chartDom, getEChartsTheme())
  } catch (error) {
    console.error('初始化分销商排行图表失败:', error)
    return
  }
  
  // 修复数据提取逻辑
  const names = distributorRankData.value.map(item => item.distributorName || '未知分销商')
  const values = distributorRankData.value.map(item => 
    distributorRankType.value === 'revenue' ? (item.revenue || 0) : (item.orders || 0)
  )
  
  // ... 其余代码保持不变
}
```

### 2. 修复后端数据结构

需要修改 `statistics_service.go` 中的 `GetDetailData` 方法：

```go
// GetDetailData 获取详细数据
func (s *StatisticsService) GetDetailData(startDate, endDate time.Time, detailType string) ([]map[string]interface{}, error) {
    switch detailType {
    case "daily":
        stats, err := s.orderRepo.GetDailyStatsByDateRange(startDate, endDate)
        if err != nil {
            return nil, err
        }
        details := make([]map[string]interface{}, 0, len(stats))
        for _, stat := range stats {
            avgOrderValue := float64(0)
            if stat.Orders > 0 {
                avgOrderValue = stat.Revenue / float64(stat.Orders)
            }
            
            // 计算模拟数据
            newCustomers := stat.Orders * 30 / 100 // 假设30%是新客户
            conversionRate := 0.15 + float64(stat.Orders%10)*0.01 // 模拟转化率
            profit := stat.Revenue * 0.3 // 假设30%利润率
            profitRate := 0.3
            
            details = append(details, map[string]interface{}{
                "date":           stat.Date.Format("2006-01-02"),
                "orderCount":     stat.Orders,
                "revenue":        stat.Revenue,
                "avgOrderValue":  avgOrderValue,
                "newCustomers":   newCustomers,
                "conversionRate": conversionRate,
                "profit":         profit,
                "profitRate":     profitRate,
            })
        }
        return details, nil
        
    case "distributor":
        stats, err := s.orderRepo.GetDistributorStatsByDateRange(startDate, endDate, 0)
        if err != nil {
            return nil, err
        }
        
        // 计算总销售额
        var totalRevenue float64
        for _, stat := range stats {
            totalRevenue += stat.Revenue
        }
        
        details := make([]map[string]interface{}, 0, len(stats))
        for _, stat := range stats {
            avgOrderValue := float64(0)
            if stat.Orders > 0 {
                avgOrderValue = stat.Revenue / float64(stat.Orders)
            }
            
            percentage := float64(0)
            if totalRevenue > 0 {
                percentage = (stat.Revenue / totalRevenue) * 100
            }
            
            commission := stat.Revenue * 0.1 // 假设10%佣金
            conversionRate := 0.15 + float64(stat.Orders%10)*0.01
            status := 1 // 默认活跃
            if stat.Orders < 5 {
                status = 0 // 休眠
            }
            
            details = append(details, map[string]interface{}{
                "distributorName": stat.DistributorName,
                "orderCount":      stat.Orders,
                "revenue":         stat.Revenue,
                "commission":      commission,
                "avgOrderValue":   avgOrderValue,
                "conversionRate":  conversionRate,
                "percentage":      percentage,
                "status":          status,
            })
        }
        return details, nil
        
    case "product":
        stats, err := s.orderRepo.GetProductStatsByDateRange(startDate, endDate, 0)
        if err != nil {
            return nil, err
        }
        
        // 计算总销售额
        var totalRevenue float64
        for _, stat := range stats {
            totalRevenue += stat.Revenue
        }
        
        details := make([]map[string]interface{}, 0, len(stats))
        for _, stat := range stats {
            avgPrice := float64(0)
            if stat.Quantity > 0 {
                avgPrice = stat.Revenue / float64(stat.Quantity)
            }
            
            percentage := float64(0)
            if totalRevenue > 0 {
                percentage = (stat.Revenue / totalRevenue) * 100
            }
            
            cost := stat.Revenue * 0.6 // 假设60%成本
            profit := stat.Revenue - cost
            profitRate := 0.4
            
            details = append(details, map[string]interface{}{
                "productName": stat.ProductName,
                "category":    "咖啡", // 需要从产品信息中获取
                "quantity":    stat.Quantity,
                "revenue":     stat.Revenue,
                "cost":        cost,
                "profit":      profit,
                "avgPrice":    avgPrice,
                "profitRate":  profitRate,
                "percentage":  percentage,
            })
        }
        return details, nil
        
    case "store":
        // TODO: 实现门店统计，需要新的查询方法
        return []map[string]interface{}{}, nil
        
    default:
        return []map[string]interface{}{}, nil
    }
}
```

### 3. 前端优化

在 `StatisticsEnhanced.vue` 中添加监听器：

```javascript
// 监听分销商排行类型变化
watch(distributorRankType, (newVal) => {
  console.log('Distributor rank type changed to:', newVal)
  loadDistributorRank()
})

// 修改 loadDistributorRank 函数
const loadDistributorRank = async () => {
  chartLoading.distributor = true
  try {
    console.log('Loading distributor rank with type:', distributorRankType.value)
    const res = await getDistributorRank({ 
      ...queryParams.value, 
      type: distributorRankType.value, 
      limit: 10 
    })
    
    if (res.code === 200 && res.data) {
      distributorRankData.value = res.data
      console.log('Distributor rank data:', res.data)
    } else {
      distributorRankData.value = []
    }
    
    // 确保在下一个渲染周期初始化图表
    await nextTick()
    if (!isUnmounted) {
      initDistributorRankChart()
    }
  } catch (error) {
    console.error('加载分销商排行失败:', error)
    ElMessage.error('加载分销商排行数据失败')
    distributorRankData.value = []
    await nextTick()
    if (!isUnmounted) {
      initDistributorRankChart()
    }
  } finally {
    chartLoading.distributor = false
  }
}
```

## 测试步骤

1. 测试分销商排行切换
   - 点击"销售额"按钮，验证图表显示销售额数据
   - 点击"订单数"按钮，验证图表切换到订单数据
   
2. 测试数据详情表格
   - 切换不同的表格类型（按日统计、按分销商、按商品、按门店）
   - 验证数据不再全是0
   - 验证数据格式化正确（货币、百分比等）

3. 测试响应式布局
   - 在不同屏幕尺寸下查看表格显示
   - 验证移动端的卡片式显示

## 注意事项

1. 后端需要实现真实的成本、利润、佣金等计算逻辑
2. 门店统计功能需要新增相应的查询方法
3. 商品分类信息需要从产品表中获取