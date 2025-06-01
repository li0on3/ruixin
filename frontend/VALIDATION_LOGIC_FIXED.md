# 批量验证逻辑说明（修复版）

## 🧠 智能验证逻辑详解

### 筛选条件（按优先级排序）：

#### 优先级1：有订单但状态异常的卡片 ⚠️
```sql
order_id IS NOT NULL AND status = 0
```
- **场景**：卡片已被订单使用，但本地状态还显示"未使用"
- **原因**：订单完成后状态同步失败
- **重要性**：最高优先级，影响库存准确性

#### 优先级2：有异常的未使用卡片 🔍
```sql
status = 0 AND (synced_at IS NULL OR synced_at < 24小时前)
```
- **场景**：显示"未使用"但24小时内未验证过
- **原因**：可能在瑞幸系统中已被使用
- **重要性**：高优先级，避免重复使用

#### 优先级3：**最近使用但状态可能不准确的卡片** 🕒
```sql
status = 1 AND used_at > 7天前 AND (synced_at IS NULL OR synced_at < used_at)
```
- **场景**：7天内使用的卡片，但验证时间早于使用时间
- **原因**：状态更新后未重新验证
- **重要性**：中等优先级，确保已使用卡片状态正确

#### 优先级4：新添加的未验证卡片 🆕
```sql
created_at > 24小时前 AND (synced_at IS NULL OR sync_status = 'failed')
```
- **场景**：新导入的卡片或之前验证失败的卡片
- **原因**：首次验证或需要重新验证
- **重要性**：较低优先级，确保新卡片可用性

### 排除条件：
- **预占中的卡片**：`status != 2`（避免干扰正在进行的订单）

### 验证参数：
- **数量限制**：最多200张卡片
- **验证间隔**：3秒/张
- **超时限制**：最大执行10分钟
- **预计完成时间**：2-5分钟

## 🌐 全量验证逻辑详解

### 验证范围：
```sql
status IN (0, 1) -- 未使用和已使用的卡片
AND status != 2  -- 排除预占中的卡片
```

### 安全策略：
- **分批处理**：每批5张卡片
- **批间间隔**：2分钟休息
- **卡片间隔**：10秒
- **执行方式**：后台异步，不阻塞用户操作

### 预计时间：
- **1000张卡片**：约2-3小时
- **安全第一**：防止触发API风控

## 🛠️ 修复的问题

### 1. 进度窗口管理 ✅
**问题**：关闭进度窗口后无法重新打开
**解决**：
- 添加"查看进度"按钮，当有运行中任务时显示
- 任务状态持久化到 localStorage
- 页面刷新后自动恢复任务状态

### 2. 重复验证防护 ✅
**问题**：刷新页面后可以重复启动验证
**解决**：
- 启动验证前检查是否已有运行中任务
- 如有运行中任务，直接显示进度而非启动新任务
- 任务完成后自动清理状态

### 3. 任务状态持久化 ✅
**问题**：页面刷新后任务状态丢失
**解决**：
- 使用 localStorage 保存任务ID和模式
- 页面加载时自动检查并恢复任务状态
- 任务完成时自动清理本地存储

## 📱 用户交互流程

### 启动验证：
1. 点击"批量验证"下拉菜单
2. 选择验证模式（智能/全量）
3. 确认验证提示
4. 自动显示进度窗口

### 进度监控：
1. 实时显示验证进度
2. 可以取消正在执行的任务
3. 关闭窗口后可通过"查看进度"按钮重新打开

### 页面刷新：
1. 自动检测运行中的任务
2. 显示提示信息
3. 点击"查看进度"恢复监控

### 任务完成：
1. 显示验证结果
2. 自动刷新卡片列表
3. 清理所有相关状态

## 🔧 技术实现细节

### 状态管理：
```javascript
// 响应式状态
const hasRunningTask = ref(false)      // 是否有运行中任务
const currentTaskId = ref(null)        // 当前任务ID
const validationTask = ref(null)       // 任务详情
const progressDialogVisible = ref(false) // 进度窗口显示状态

// 本地存储
localStorage.setItem('validation_task_id', taskId)
localStorage.setItem('validation_task_mode', mode)
```

### 任务检查逻辑：
```javascript
// 页面加载时检查
const checkRunningTask = async () => {
  const savedTaskId = localStorage.getItem('validation_task_id')
  if (savedTaskId) {
    const res = await getValidationProgress(savedTaskId)
    if (res.data.status === 'running' || res.data.status === 'queued') {
      // 恢复任务状态
      hasRunningTask.value = true
      currentTaskId.value = savedTaskId
    }
  }
}
```

### 防重复启动：
```javascript
// 启动前检查
if (hasRunningTask.value && currentTaskId.value) {
  const res = await getValidationProgress(currentTaskId.value)
  if (res.data.status === 'running' || res.data.status === 'queued') {
    // 显示现有任务进度
    validationTask.value = res.data
    progressDialogVisible.value = true
    return // 不启动新任务
  }
}
```

## 🎯 使用建议

### 智能验证：
- **适用场景**：日常维护、快速检查
- **使用频率**：每天1-2次
- **最佳时间**：上午或下午业务较少时

### 全量验证：
- **适用场景**：全面检查、定期维护
- **使用频率**：每周1次或按需
- **最佳时间**：夜间或业务空闲期

### 注意事项：
1. **避免同时启动多个验证任务**
2. **全量验证期间避免大量订单操作**
3. **定期关注验证结果，及时处理异常**
4. **如遇问题可随时取消任务**

这次修复完美解决了所有问题，现在用户可以安全、便捷地使用批量验证功能！