# 深色主题白色背景修复任务

## 任务描述
用户反馈以下页面在深色主题下还有白色背景问题：
- statistics (统计页面)
- cards (卡片管理)
- products_available (可用商品)
- stores (店铺查询)
- admins (管理员管理)
- settings (系统设置)

## 问题分析
经过检查，发现这些页面使用了以下样式类：
1. `search-card` - 搜索卡片
2. `table-card` - 表格卡片
3. `metric-card` - 指标卡片
4. `result-card` - 结果卡片
5. `hot-stores-card` - 热门店铺卡片
6. `product-card` - 商品卡片
7. `store-card` - 店铺卡片
8. `page-container` - 页面容器

## 已完成的修复
在 `/home/li0on/project/ruixin/frontend/src/assets/styles/dark-theme-ultimate-fix.scss` 文件中添加了以下修复：

1. **通用页面容器修复**
   - `.page-container` 设置透明背景和正确的文字颜色

2. **特定页面卡片修复**
   - 统计页面的所有卡片背景
   - 结果卡片和热门店铺卡片
   - 店铺卡片的背景和边框
   - 商品卡片的背景和文字颜色

3. **全局卡片背景修复**
   - 确保所有 `el-card` 组件都有深色背景
   - 卡片头部使用稍浅的深色背景
   - 表格卡片特别处理
   - 告警组件的深色主题适配

## 修复要点
- 背景色使用 `#1e293b` (深色)
- 卡片头部使用 `#334155` (稍浅)
- 文字颜色使用 `#f1f5f9` (亮色)
- 次要文字使用 `#94a3b8` (灰色)
- 边框颜色使用 `#334155` 或 `#475569`
- 高亮色使用 `#fb923c` (橙色)

## 验证方法
1. 启动前端开发服务器
2. 切换到深色主题
3. 访问每个页面检查背景色是否正确
4. 确保文字可读性良好
5. 检查交互状态（hover、active等）

## 后续建议
如果还有其他页面出现白色背景问题，可以继续在 `dark-theme-ultimate-fix.scss` 文件中添加相应的修复规则。