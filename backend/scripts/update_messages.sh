#!/bin/bash

# 批量替换英文消息为中文

cd /home/li0on/project/ruixin/backend

# 通用消息替换
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Success"/"msg":  "成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "OK"/"msg":  "成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed"/"msg":  "失败"/g' {} +

# 认证相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid credentials"/"msg":  "用户名或密码错误"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Login successful"/"msg":  "登录成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Logout successful"/"msg":  "登出成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Unauthorized"/"msg":  "未授权"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Token expired"/"msg":  "令牌已过期"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid token"/"msg":  "无效的令牌"/g' {} +

# 参数验证
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid request parameters"/"msg":  "请求参数无效"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid parameters"/"msg":  "参数无效"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Missing required parameters"/"msg":  "缺少必需参数"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid ID"/"msg":  "无效的ID"/g' {} +

# 卡片相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid card ID"/"msg":  "无效的卡片ID"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Card not found"/"msg":  "卡片不存在"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Card created successfully"/"msg":  "卡片创建成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Card updated successfully"/"msg":  "卡片更新成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Card deleted successfully"/"msg":  "卡片删除成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to create card"/"msg":  "创建卡片失败"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to update card"/"msg":  "更新卡片失败"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to delete card"/"msg":  "删除卡片失败"/g' {} +

# 分销商相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Distributor not found"/"msg":  "分销商不存在"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid distributor ID"/"msg":  "无效的分销商ID"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Distributor created successfully"/"msg":  "分销商创建成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Distributor updated successfully"/"msg":  "分销商更新成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Distributor deleted successfully"/"msg":  "分销商删除成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Password reset successfully"/"msg":  "密码重置成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "API key reset successfully"/"msg":  "API密钥重置成功"/g' {} +

# 订单相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Order not found"/"msg":  "订单不存在"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Invalid order ID"/"msg":  "无效的订单ID"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Order created successfully"/"msg":  "订单创建成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Order status refreshed"/"msg":  "订单状态已刷新"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to create order"/"msg":  "创建订单失败"/g' {} +

# 财务相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Insufficient balance"/"msg":  "余额不足"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Recharge successful"/"msg":  "充值成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Adjustment successful"/"msg":  "调整成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Withdrawal created successfully"/"msg":  "提现申请创建成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Withdrawal processed successfully"/"msg":  "提现处理成功"/g' {} +

# 通用操作
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Operation successful"/"msg":  "操作成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Operation failed"/"msg":  "操作失败"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Update successful"/"msg":  "更新成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Delete successful"/"msg":  "删除成功"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Create successful"/"msg":  "创建成功"/g' {} +

# 数据获取相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to get data"/"msg":  "获取数据失败"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to get list"/"msg":  "获取列表失败"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Data not found"/"msg":  "数据不存在"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Record not found"/"msg":  "记录不存在"/g' {} +

# 城市相关
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Failed to sync cities"/"msg":  "同步城市数据失败"/g' {} +
find . -name "*.go" -type f -exec sed -i 's/"msg":  "Cities synced successfully"/"msg":  "城市数据同步成功"/g' {} +

echo "消息更新完成！"