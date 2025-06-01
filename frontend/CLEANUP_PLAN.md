# 前端代码清理计划

## 清理时间：2024年1月

## 当前使用的路由配置
根据 `main.js`，系统使用 `router/optimized.js` 作为路由配置。

## 正在使用的页面组件（必须保留）

根据 `router/optimized.js` 的配置，以下页面正在使用：

1. **LoginEnhanced.vue** - 登录页面
2. **LayoutPremium.vue** - 布局组件
3. **DashboardPremium.vue** - 仪表盘
4. **cards/CardsEnhanced.vue** - 卡片管理
5. **cards/Bindings.vue** - 卡片绑定
6. **cards/BatchImport.vue** - 批量导入
7. **cards/Batches.vue** - 批次管理
8. **distributors/Index.vue** - 分销商管理（注意：使用的是 Index.vue，不是 EnhancedIndex.vue）
9. **orders/EnhancedIndex.vue** - 订单管理
10. **statistics/StatisticsEnhanced.vue** - 统计页面
11. **finance/Transactions.vue** - 交易记录
12. **finance/Withdrawals.vue** - 提现管理
13. **luckin/Prices.vue** - 价格管理
14. **products/Index.vue** - 商品管理
15. **products/AvailableProducts.vue** - 可用商品
16. **admins/Index.vue** - 管理员管理
17. **settings/Index.vue** - 系统设置
18. **stores/Index.vue** - 店铺查询
19. **404.vue** - 404页面

## 可以安全删除的页面组件

### 1. 登录相关
- ❌ `/src/views/Login.vue` - 旧版登录页面

### 2. 布局组件
- ❌ `/src/views/Layout.vue` - 旧版布局
- ❌ `/src/views/LayoutEnhanced.vue` - 增强版布局（已被 LayoutPremium.vue 替代）

### 3. 仪表盘
- ❌ `/src/views/Dashboard.vue` - 旧版仪表盘
- ❌ `/src/views/DashboardEnhanced.vue` - 增强版仪表盘（已被 DashboardPremium.vue 替代）

### 4. 卡片管理
- ❌ `/src/views/cards/Index.vue` - 旧版卡片管理（已被 CardsEnhanced.vue 替代）

### 5. 分销商管理
- ❌ `/src/views/distributors/EnhancedIndex.vue` - 增强版分销商管理（未使用）

### 6. 订单管理
- ❌ `/src/views/orders/Index.vue` - 旧版订单管理（已被 EnhancedIndex.vue 替代）

### 7. 统计页面
- ❌ `/src/views/statistics/Index.vue` - 旧版统计页面（已被 StatisticsEnhanced.vue 替代）

### 8. 路由配置
- ❌ `/src/router/index.js` - 旧版路由配置（已被 optimized.js 替代）

## 可以删除的样式文件

### 深色主题相关（已合并到 dark-theme-unified.scss）
- ❌ `/src/assets/styles/dark-theme-debug.scss`
- ❌ `/src/assets/styles/dark-theme-fix.scss`
- ❌ `/src/assets/styles/dark-theme-fixes.scss`
- ❌ `/src/assets/styles/dark-theme-optimization.scss`
- ❌ `/src/assets/styles/dark-theme-pages-fix.scss`
- ❌ `/src/assets/styles/dark-theme-ultimate-fix.scss`
- ❌ `/src/assets/styles/dark-theme-force-fix.scss`
- ❌ `/src/assets/styles/dark-theme-override.scss`

### 其他可能废弃的样式
- ❌ `/src/assets/styles/animations.scss` - 如果已被 animations-enhanced.scss 替代
- ❌ `/src/assets/styles/modern.scss` - 如果已被 modern-enhanced.scss 替代

## 可以删除的插件文件
- ❌ `/src/plugins/dark-theme-runtime-fix.js` - 运行时修复插件（已不再使用）

## 删除前的确认清单

1. ✅ 确认 `router/optimized.js` 是当前使用的路由配置
2. ✅ 确认没有其他地方引用这些要删除的组件
3. ✅ 确认深色主题在 `dark-theme-unified.scss` 中正常工作
4. ✅ 建议先备份这些文件，以防需要回滚

## 删除命令（请谨慎执行）

```bash
# 备份文件（可选）
mkdir -p /home/li0on/project/ruixin/frontend/backup
cp -r /home/li0on/project/ruixin/frontend/src/views/Login.vue /home/li0on/project/ruixin/frontend/backup/
# ... 其他文件类似

# 删除页面组件
rm -f /home/li0on/project/ruixin/frontend/src/views/Login.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/Layout.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/LayoutEnhanced.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/Dashboard.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/DashboardEnhanced.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/cards/Index.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/distributors/EnhancedIndex.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/orders/Index.vue
rm -f /home/li0on/project/ruixin/frontend/src/views/statistics/Index.vue

# 删除路由配置
rm -f /home/li0on/project/ruixin/frontend/src/router/index.js

# 删除深色主题样式
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-debug.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-fix.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-fixes.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-optimization.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-pages-fix.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-ultimate-fix.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-force-fix.scss
rm -f /home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-override.scss

# 删除插件
rm -f /home/li0on/project/ruixin/frontend/src/plugins/dark-theme-runtime-fix.js
```

## 注意事项

1. **distributors/Index.vue 必须保留**，因为路由中使用的是 Index.vue，不是 EnhancedIndex.vue
2. 所有正在使用的页面都已标记为"必须保留"
3. 建议在删除前先进行备份
4. 删除后需要测试系统功能是否正常