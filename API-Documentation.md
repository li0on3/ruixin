# 瑞幸咖啡管理系统 API 文档

## 目录

- [概述](#概述)
- [通用说明](#通用说明)
  - [基础信息](#基础信息)
  - [认证方式](#认证方式)
  - [响应格式](#响应格式)
  - [错误码说明](#错误码说明)
- [管理员端 API](#管理员端-api)
  - [认证管理](#认证管理)
  - [仪表盘](#仪表盘)
  - [卡片管理](#卡片管理)
  - [分销商管理](#分销商管理)
  - [订单管理](#订单管理)
  - [财务管理](#财务管理)
  - [商品管理](#商品管理)
  - [统计分析](#统计分析)
  - [系统配置](#系统配置)
- [分销商端 API](#分销商端-api)
  - [认证管理](#分销商认证管理)
  - [订单管理](#分销商订单管理)
  - [店铺和商品查询](#店铺和商品查询)
  - [财务管理](#分销商财务管理)

## 概述

瑞幸咖啡管理系统提供了完整的 RESTful API 接口，支持管理员和分销商两种角色的操作。本文档详细说明了所有 API 的使用方法、参数说明和示例。

## 通用说明

### 基础信息

- **基础URL**: `http://localhost:8080`
- **API版本**: `v1`
- **完整路径格式**: `{base_url}/api/{version}/{endpoint}`

### 认证方式

系统支持两种认证方式：

#### 1. JWT Token 认证（用于用户登录后的操作）
```
Authorization: Bearer {token}
```

#### 2. API Key 认证（用于分销商 API 调用）
```
X-API-Key: {api_key}
X-API-Secret: {api_secret}
```

### 响应格式

所有 API 响应采用统一的 JSON 格式：

```json
{
    "code": 200,        // 状态码
    "msg": "成功",      // 提示信息
    "data": {}          // 响应数据
}
```

### 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 管理员端 API

### 认证管理

#### 1. 管理员登录

**接口地址**: `POST /api/v1/admin/login`

**接口说明**: 管理员登录接口，成功后返回 JWT token

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:
```json
{
    "username": "admin",
    "password": "admin123"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@example.com",
            "real_name": "系统管理员",
            "role": "super_admin",
            "created_at": "2024-01-01T00:00:00Z"
        }
    }
}
```

#### 2. 管理员登出

**接口地址**: `POST /api/v1/admin/logout`

**接口说明**: 管理员登出接口

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "登出成功",
    "data": null
}
```

#### 3. 获取当前用户信息

**接口地址**: `GET /api/v1/admin/user/info`

**接口说明**: 获取当前登录管理员的详细信息

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "phone": "13800138000",
        "real_name": "系统管理员",
        "role": "super_admin",
        "status": 1,
        "last_login_at": "2024-01-15T10:30:00Z",
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 4. 修改密码

**接口地址**: `PUT /api/v1/admin/user/password`

**接口说明**: 修改管理员密码

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码，最少6位 |

**请求示例**:
```json
{
    "old_password": "admin123",
    "new_password": "newPassword123"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "Password changed successfully",
    "data": null
}
```

### 仪表盘

#### 1. 获取统计数据

**接口地址**: `GET /api/v1/admin/dashboard/statistics`

**接口说明**: 获取仪表盘总体统计数据，包括订单数、销售额、用户数等

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "today_orders": 156,
        "today_revenue": 3890.50,
        "total_orders": 12580,
        "total_revenue": 314520.00,
        "active_distributors": 45,
        "available_cards": 1250
    }
}
```

#### 2. 获取订单趋势

**接口地址**: `GET /api/v1/admin/dashboard/order-trend`

**接口说明**: 获取指定时间范围的订单趋势数据

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| dateRange | string | 否 | 时间范围：7/30/90/week/month，默认7天 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": [
        {
            "date": "2024-01-10",
            "orders": 145,
            "revenue": 3625.00
        },
        {
            "date": "2024-01-11",
            "orders": 168,
            "revenue": 4200.00
        }
    ]
}
```

#### 3. 获取热门商品

**接口地址**: `GET /api/v1/admin/dashboard/hot-goods`

**接口说明**: 获取前10个热门商品

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": [
        {
            "goods_code": "GOODS001",
            "goods_name": "拿铁咖啡",
            "sales_count": 1250,
            "sales_amount": 31250.00
        },
        {
            "goods_code": "GOODS002",
            "goods_name": "美式咖啡",
            "sales_count": 980,
            "sales_amount": 21560.00
        }
    ]
}
```

#### 4. 获取最新订单

**接口地址**: `GET /api/v1/admin/dashboard/recent-orders`

**接口说明**: 获取最新的10个订单

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": [
        {
            "order_no": "ORD202401150001",
            "distributor_name": "测试分销商",
            "store_name": "瑞幸咖啡(朝阳门店)",
            "total_amount": 25.00,
            "status": "completed",
            "created_at": "2024-01-15T14:30:00Z"
        }
    ]
}
```

### 卡片管理

#### 1. 获取卡片列表

**接口地址**: `GET /api/v1/admin/cards`

**接口说明**: 分页获取卡片列表，支持多条件筛选

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| status | int | 否 | 状态：0-未激活，1-正常，2-已使用，3-已过期，4-已作废 |
| price_id | int | 否 | 价格ID |
| search | string | 否 | 搜索关键词（卡号） |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "list": [
            {
                "id": 1,
                "card_code": "CARD001",
                "price_id": 1,
                "price_value": 25.00,
                "cost_price": 20.00,
                "sell_price": 25.00,
                "status": 1,
                "expired_at": "2024-12-31T23:59:59Z",
                "created_at": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 100
    }
}
```

#### 2. 获取卡片详情

**接口地址**: `GET /api/v1/admin/cards/:id`

**接口说明**: 根据ID获取单个卡片详情

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 卡片ID |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 1,
        "card_code": "CARD001",
        "price_id": 1,
        "price_value": 25.00,
        "cost_price": 20.00,
        "sell_price": 25.00,
        "status": 1,
        "description": "25元通用券",
        "expired_at": "2024-12-31T23:59:59Z",
        "created_at": "2024-01-01T00:00:00Z",
        "usage_count": 5,
        "last_used_at": "2024-01-15T10:00:00Z"
    }
}
```

#### 3. 创建卡片

**接口地址**: `POST /api/v1/admin/cards`

**接口说明**: 创建新的卡片

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card_code | string | 是 | 卡片代码，必须唯一 |
| price_id | int | 是 | 价格ID |
| cost_price | float | 是 | 成本价 |
| sell_price | float | 是 | 销售价 |
| expired_at | string | 是 | 过期时间，ISO 8601格式 |
| description | string | 否 | 描述信息 |

**请求示例**:
```json
{
    "card_code": "TEST123456",
    "price_id": 1,
    "cost_price": 10.0,
    "sell_price": 15.0,
    "expired_at": "2024-12-31T23:59:59Z",
    "description": "测试卡片"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "卡片创建成功",
    "data": {
        "id": 100,
        "card_code": "TEST123456",
        "status": 1
    }
}
```

#### 4. 更新卡片

**接口地址**: `PUT /api/v1/admin/cards/:id`

**接口说明**: 更新卡片信息

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 卡片ID |

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | int | 否 | 状态 |
| cost_price | float | 否 | 成本价 |
| sell_price | float | 否 | 销售价 |
| expired_at | string | 否 | 过期时间 |
| description | string | 否 | 描述信息 |

**请求示例**:
```json
{
    "status": 1,
    "cost_price": 12.0,
    "sell_price": 18.0,
    "expired_at": "2025-12-31T23:59:59Z",
    "description": "更新后的描述"
}
```

#### 5. 删除卡片

**接口地址**: `DELETE /api/v1/admin/cards/:id`

**接口说明**: 删除指定卡片

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 卡片ID |

**响应示例**:
```json
{
    "code": 200,
    "msg": "卡片删除成功",
    "data": null
}
```

#### 6. 批量导入卡片

**接口地址**: `POST /api/v1/admin/cards/batch-import`

**接口说明**: 批量导入卡片

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| price_id | int | 是 | 价格ID |
| cards | array | 是 | 卡片列表 |
| cards[].card_code | string | 是 | 卡片代码 |
| cards[].cost_price | float | 是 | 成本价 |
| cards[].sell_price | float | 是 | 销售价 |
| expired_at | string | 是 | 统一过期时间 |
| description | string | 否 | 批次描述 |

**请求示例**:
```json
{
    "price_id": 1,
    "cards": [
        {
            "card_code": "BATCH001",
            "cost_price": 10.0,
            "sell_price": 15.0
        },
        {
            "card_code": "BATCH002",
            "cost_price": 10.0,
            "sell_price": 15.0
        }
    ],
    "expired_at": "2024-12-31T23:59:59Z",
    "description": "批量导入测试"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "批量导入成功",
    "data": {
        "batch_id": 10,
        "success_count": 2,
        "failed_count": 0
    }
}
```

#### 7. 验证单个卡片

**接口地址**: `POST /api/v1/admin/cards/validate`

**接口说明**: 验证单个卡片的有效性并更新状态

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card_code | string | 是 | 卡片代码 |

**请求示例**:
```json
{
    "card_code": "TEST123456"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "验证完成",
    "data": {
        "is_valid": true,
        "message": "卡片有效",
        "updated": true
    }
}
```

#### 8. 启动批量验证任务

**接口地址**: `POST /api/v1/admin/cards/batch-validation/start`

**接口说明**: 启动异步批量验证任务

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mode | string | 是 | 验证模式：all-全量验证，smart-智能验证 |

**请求示例**:
```json
{
    "mode": "smart"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "验证任务已启动",
    "data": {
        "task_id": "task_20240115_143000",
        "mode": "smart",
        "total_cards": 1000,
        "status": "running"
    }
}
```

#### 9. 获取验证进度

**接口地址**: `GET /api/v1/admin/cards/batch-validation/:taskId`

**接口说明**: 获取批量验证任务的进度

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| taskId | string | 是 | 任务ID |

**响应示例**:
```json
{
    "code": 200,
    "msg": "获取成功",
    "data": {
        "task_id": "task_20240115_143000",
        "status": "running",
        "progress": 45.5,
        "validated_count": 455,
        "total_count": 1000,
        "success_count": 400,
        "failed_count": 55,
        "elapsed_time": 120
    }
}
```

### 分销商管理

#### 1. 获取分销商列表

**接口地址**: `GET /api/v1/admin/distributors`

**接口说明**: 分页获取分销商列表

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| status | int | 否 | 状态：1-正常，0-禁用 |
| name | string | 否 | 分销商名称搜索 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "list": [
            {
                "id": 1,
                "name": "测试分销商",
                "company_name": "测试公司",
                "contact_name": "张三",
                "phone": "13800138000",
                "email": "test@example.com",
                "status": 1,
                "balance": 5000.00,
                "credit_limit": 10000.00,
                "created_at": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 45,
        "page": 1,
        "page_size": 20
    }
}
```

#### 2. 获取分销商详情

**接口地址**: `GET /api/v1/admin/distributors/:id`

**接口说明**: 获取分销商详细信息，包含API密钥

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分销商ID |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 1,
        "name": "测试分销商",
        "company_name": "测试公司",
        "contact_name": "张三",
        "phone": "13800138000",
        "email": "test@example.com",
        "api_key": "dist_123456789",
        "api_secret": "secret_abcdefghijk",
        "status": 1,
        "balance": 5000.00,
        "credit_limit": 10000.00,
        "callback_url": "https://example.com/callback",
        "daily_order_limit": 100,
        "monthly_order_limit": 3000,
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 3. 创建分销商

**接口地址**: `POST /api/v1/admin/distributors`

**接口说明**: 创建新的分销商账号

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 分销商名称 |
| company_name | string | 是 | 公司名称 |
| contact_name | string | 是 | 联系人姓名 |
| phone | string | 是 | 联系电话 |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 登录密码 |
| callback_url | string | 否 | 回调地址 |
| credit_limit | float | 否 | 信用额度，默认0 |
| daily_order_limit | int | 否 | 日订单限制，默认100 |
| monthly_order_limit | int | 否 | 月订单限制，默认3000 |

**请求示例**:
```json
{
    "name": "测试分销商",
    "company_name": "测试公司",
    "contact_name": "张三",
    "phone": "13800138000",
    "email": "test@example.com",
    "password": "password123",
    "callback_url": "https://example.com/callback",
    "credit_limit": 10000.0,
    "daily_order_limit": 100,
    "monthly_order_limit": 3000
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "分销商创建成功",
    "data": {
        "id": 10,
        "api_key": "dist_987654321",
        "api_secret": "secret_zyxwvutsrq",
        "default_password": "password123"
    }
}
```

#### 4. 更新分销商

**接口地址**: `PUT /api/v1/admin/distributors/:id`

**接口说明**: 更新分销商信息

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分销商ID |

**请求参数**: 所有字段都是可选的
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 分销商名称 |
| company_name | string | 否 | 公司名称 |
| status | int | 否 | 状态：1-正常，0-禁用 |
| credit_limit | float | 否 | 信用额度 |

**请求示例**:
```json
{
    "name": "更新后的名称",
    "company_name": "更新后的公司名",
    "status": 1,
    "credit_limit": 20000.0
}
```

#### 5. 重置API密钥

**接口地址**: `POST /api/v1/admin/distributors/:id/reset-api-key`

**接口说明**: 重新生成分销商的API密钥

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 分销商ID |

**响应示例**:
```json
{
    "code": 200,
    "msg": "API密钥重置成功",
    "data": {
        "api_key": "dist_new123456",
        "api_secret": "secret_new789xyz"
    }
}
```

### 订单管理

#### 1. 获取订单列表

**接口地址**: `GET /api/v1/admin/orders`

**接口说明**: 分页获取订单列表，支持多条件筛选

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| distributor_id | int | 否 | 分销商ID |
| status | string | 否 | 订单状态 |
| store_code | string | 否 | 店铺代码 |
| card_code | string | 否 | 卡片代码 |
| start_date | string | 否 | 开始日期 |
| end_date | string | 否 | 结束日期 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "list": [
            {
                "id": 1,
                "order_no": "ORD202401150001",
                "outer_order_no": "DIST_ORDER_001",
                "distributor_id": 1,
                "distributor_name": "测试分销商",
                "card_code": "CARD001",
                "store_code": "STORE001",
                "store_name": "瑞幸咖啡(朝阳门店)",
                "store_address": "北京市朝阳区朝阳门外大街1号",
                "mobile": "13800138000",
                "total_amount": 25.00,
                "status": "completed",
                "luckin_order_no": "LK202401150001",
                "take_code": "A001",
                "created_at": "2024-01-15T10:00:00Z",
                "completed_at": "2024-01-15T10:05:00Z"
            }
        ],
        "total": 1000,
        "page": 1,
        "page_size": 20
    }
}
```

#### 2. 获取订单详情

**接口地址**: `GET /api/v1/admin/orders/:orderNo`

**接口说明**: 获取指定订单的详细信息

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| orderNo | string | 是 | 订单号 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 1,
        "order_no": "ORD202401150001",
        "outer_order_no": "DIST_ORDER_001",
        "distributor": {
            "id": 1,
            "name": "测试分销商",
            "company_name": "测试公司"
        },
        "card": {
            "card_code": "CARD001",
            "price_value": 25.00
        },
        "store": {
            "store_code": "STORE001",
            "store_name": "瑞幸咖啡(朝阳门店)",
            "store_address": "北京市朝阳区朝阳门外大街1号"
        },
        "goods": [
            {
                "goods_code": "GOODS001",
                "goods_name": "拿铁咖啡",
                "specs_code": "HOT_LARGE",
                "price": 25.00,
                "num": 1
            }
        ],
        "mobile": "13800138000",
        "take_time": "立即自取",
        "total_amount": 25.00,
        "status": "completed",
        "luckin_order_no": "LK202401150001",
        "take_code": "A001",
        "qr_code": "https://example.com/qr/xxx",
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:05:00Z",
        "completed_at": "2024-01-15T10:05:00Z"
    }
}
```

#### 3. 刷新订单状态

**接口地址**: `POST /api/v1/admin/orders/:orderNo/refresh`

**接口说明**: 从瑞幸系统刷新订单状态

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| orderNo | string | 是 | 订单号 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "订单状态已刷新",
    "data": {
        "order_no": "ORD202401150001",
        "status": "completed",
        "luckin_status": "已完成",
        "updated": true
    }
}
```

#### 4. 生成订单二维码

**接口地址**: `POST /api/v1/admin/orders/:orderNo/qrcode`

**接口说明**: 生成订单的取货二维码

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| orderNo | string | 是 | 订单号 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "二维码生成成功",
    "data": {
        "qr_code_image": "data:image/png;base64,iVBORw0KGgoAAAANS...",
        "qr_data": "https://luckin.com/pickup/xxx",
        "take_code": "A001"
    }
}
```

### 财务管理

#### 1. 充值

**接口地址**: `POST /api/v1/admin/finance/recharge`

**接口说明**: 为分销商充值

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| distributor_id | int | 是 | 分销商ID |
| amount | float | 是 | 充值金额，必须大于0 |
| remark | string | 否 | 备注信息 |

**请求示例**:
```json
{
    "distributor_id": 1,
    "amount": 1000.0,
    "remark": "管理员充值"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "充值成功",
    "data": {
        "transaction_id": 100,
        "balance_after": 6000.00
    }
}
```

#### 2. 余额调整

**接口地址**: `POST /api/v1/admin/finance/adjust`

**接口说明**: 调整分销商余额（正负均可）

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| distributor_id | int | 是 | 分销商ID |
| amount | float | 是 | 调整金额，正数增加，负数减少 |
| remark | string | 是 | 备注信息 |

**请求示例**:
```json
{
    "distributor_id": 1,
    "amount": -50.0,
    "remark": "余额调整：退款"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "余额调整成功",
    "data": {
        "transaction_id": 101,
        "balance_after": 5950.00
    }
}
```

#### 3. 获取交易记录列表

**接口地址**: `GET /api/v1/admin/finance/transactions`

**接口说明**: 获取所有交易记录

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| distributor_id | int | 否 | 分销商ID |
| type | string | 否 | 交易类型：recharge/order/adjust/withdraw |
| start_date | string | 否 | 开始日期 |
| end_date | string | 否 | 结束日期 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "data": [
            {
                "id": 100,
                "distributor_id": 1,
                "distributor_name": "测试分销商",
                "type": "recharge",
                "amount": 1000.00,
                "balance_before": 5000.00,
                "balance_after": 6000.00,
                "order_no": null,
                "remark": "管理员充值",
                "created_at": "2024-01-15T10:00:00Z"
            }
        ],
        "pagination": {
            "page": 1,
            "page_size": 20,
            "total": 500,
            "total_pages": 25
        }
    }
}
```

#### 4. 处理提现申请

**接口地址**: `POST /api/v1/admin/finance/withdrawals/:id/process`

**接口说明**: 批准或拒绝提现申请

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 提现申请ID |

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| approved | bool | 是 | 是否批准 |
| reject_reason | string | 否 | 拒绝原因（拒绝时必填） |

**请求示例**:
```json
{
    "approved": true,
    "reject_reason": ""
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "提现申请处理成功",
    "data": {
        "withdrawal_id": 10,
        "status": "approved",
        "processed_at": "2024-01-15T15:00:00Z"
    }
}
```

### 商品管理

#### 1. 获取商品列表

**接口地址**: `GET /api/v1/admin/products`

**接口说明**: 获取商品列表或指定卡片绑定的产品

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| search | string | 否 | 搜索关键词 |
| card_id | int | 否 | 卡片ID，获取该卡片绑定的产品 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "获取成功",
    "data": {
        "list": [
            {
                "id": 1,
                "goods_code": "GOODS001",
                "goods_name": "拿铁咖啡",
                "goods_pic": "https://example.com/latte.jpg",
                "category": "咖啡",
                "specs": [
                    {
                        "sku_code": "SKU001",
                        "sku_name": "大杯",
                        "price": 25.00
                    }
                ],
                "created_at": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 100
    }
}
```

#### 2. 搜索商品

**接口地址**: `GET /api/v1/admin/products/search`

**接口说明**: 根据关键词搜索商品

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 是 | 搜索关键词 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "搜索成功",
    "data": [
        {
            "goods_code": "GOODS001",
            "goods_name": "拿铁咖啡",
            "category": "咖啡"
        },
        {
            "goods_code": "GOODS010",
            "goods_name": "榛果拿铁",
            "category": "咖啡"
        }
    ]
}
```

#### 3. 同步商品信息

**接口地址**: `POST /api/v1/admin/products/sync`

**接口说明**: 从瑞幸系统同步商品信息

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| store_code | string | 否 | 店铺代码，不传则使用系统配置 |
| card_code | string | 是 | 卡片代码 |

**请求示例**:
```json
{
    "store_code": "STORE001",
    "card_code": "CARD001"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "同步成功",
    "data": {
        "synced_count": 50,
        "new_count": 5,
        "updated_count": 10,
        "failed_count": 0,
        "alias_count": 3,
        "duration": 2.5
    }
}
```

#### 4. 解析规格代码

**接口地址**: `POST /api/v1/admin/products/parse-specs`

**接口说明**: 解析商品规格代码为具体规格项

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| goods_code | string | 是 | 商品代码 |
| sku_code | string | 是 | SKU代码 |
| specs_code | string | 是 | 规格代码 |

**请求示例**:
```json
{
    "goods_code": "GOODS001",
    "sku_code": "SKU001",
    "specs_code": "HOT_LARGE_SUGAR"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "解析成功",
    "data": [
        {
            "spec_name": "温度",
            "spec_value": "热"
        },
        {
            "spec_name": "规格",
            "spec_value": "大杯"
        },
        {
            "spec_name": "糖度",
            "spec_value": "标准糖"
        }
    ]
}
```

### 统计分析

#### 1. 获取核心指标

**接口地址**: `GET /api/v1/admin/statistics/metrics`

**接口说明**: 获取指定时间范围的核心业务指标

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| start_date | string | 是 | 开始日期，格式：YYYY-MM-DD |
| end_date | string | 是 | 结束日期，格式：YYYY-MM-DD |

**响应示例**:
```json
{
    "code": 200,
    "msg": "Success",
    "data": {
        "total_orders": 5000,
        "total_revenue": 125000.00,
        "average_order_value": 25.00,
        "total_cost": 100000.00,
        "gross_profit": 25000.00,
        "profit_margin": 20.0,
        "active_distributors": 45,
        "new_distributors": 5
    }
}
```

#### 2. 获取销售趋势

**接口地址**: `GET /api/v1/admin/statistics/sales-trend`

**接口说明**: 获取销售趋势数据

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| start_date | string | 是 | 开始日期 |
| end_date | string | 是 | 结束日期 |
| type | string | 否 | 数据类型：both/revenue/quantity，默认both |

**响应示例**:
```json
{
    "code": 200,
    "msg": "Success",
    "data": [
        {
            "date": "2024-01-15",
            "revenue": 5250.00,
            "quantity": 210,
            "orders": 180
        },
        {
            "date": "2024-01-16",
            "revenue": 6000.00,
            "quantity": 240,
            "orders": 200
        }
    ]
}
```

#### 3. 获取分销商排行

**接口地址**: `GET /api/v1/admin/statistics/distributor-rank`

**接口说明**: 获取分销商销售排行榜

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| start_date | string | 是 | 开始日期 |
| end_date | string | 是 | 结束日期 |
| type | string | 否 | 排行类型：revenue/quantity，默认revenue |
| limit | int | 否 | 返回数量，默认10 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "Success",
    "data": [
        {
            "rank": 1,
            "distributor_id": 1,
            "distributor_name": "顶级分销商",
            "total_revenue": 50000.00,
            "total_orders": 2000,
            "average_order_value": 25.00
        },
        {
            "rank": 2,
            "distributor_id": 2,
            "distributor_name": "优质分销商",
            "total_revenue": 35000.00,
            "total_orders": 1400,
            "average_order_value": 25.00
        }
    ]
}
```

#### 4. 获取商品分析

**接口地址**: `GET /api/v1/admin/statistics/product-analysis`

**接口说明**: 获取商品销售分析数据

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| start_date | string | 是 | 开始日期 |
| end_date | string | 是 | 结束日期 |
| limit | int | 否 | 返回数量，默认10 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "Success",
    "data": [
        {
            "goods_code": "GOODS001",
            "goods_name": "拿铁咖啡",
            "sales_count": 1500,
            "sales_amount": 37500.00,
            "percentage": 30.0
        },
        {
            "goods_code": "GOODS002",
            "goods_name": "美式咖啡",
            "sales_count": 1200,
            "sales_amount": 26400.00,
            "percentage": 21.12
        }
    ]
}
```

### 系统配置

#### 1. 获取系统配置

**接口地址**: `GET /api/v1/admin/system/configs`

**接口说明**: 获取所有系统配置

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "max_validation_workers": "10",
        "validation_timeout": "30",
        "batch_validation_enabled": "true",
        "smart_validation_interval": "3600",
        "order_callback_timeout": "30",
        "order_callback_retry": "3"
    }
}
```

#### 2. 更新系统配置

**接口地址**: `PUT /api/v1/admin/system/configs`

**接口说明**: 更新系统配置项

**请求参数**: 键值对形式的配置项
```json
{
    "max_validation_workers": "10",
    "validation_timeout": "30",
    "batch_validation_enabled": "true"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "配置更新成功",
    "data": null
}
```

## 分销商端 API

### 分销商认证管理

#### 1. 分销商登录

**接口地址**: `POST /api/v1/distributor/login`

**接口说明**: 分销商登录接口，成功后返回JWT token和API密钥

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 登录密码 |

**请求示例**:
```json
{
    "email": "distributor@example.com",
    "password": "password123"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "distributor": {
            "id": 1,
            "name": "测试分销商",
            "email": "distributor@example.com",
            "company_name": "测试公司",
            "api_key": "dist_123456789",
            "api_secret": "secret_abcdefghijk"
        }
    }
}
```

#### 2. 获取分销商信息

**接口地址**: `GET /api/v1/distributor/profile`

**接口说明**: 获取当前分销商的详细信息

**请求头**:
```
Authorization: Bearer {token}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "id": 1,
        "name": "测试分销商",
        "company_name": "测试公司",
        "contact_name": "张三",
        "phone": "13800138000",
        "email": "distributor@example.com",
        "api_key": "dist_123456789",
        "api_secret": "secret_abcdefghijk",
        "status": 1,
        "balance": 5000.00,
        "credit_limit": 10000.00,
        "callback_url": "https://example.com/callback",
        "daily_order_limit": 100,
        "monthly_order_limit": 3000
    }
}
```

### 分销商订单管理

#### 1. 创建订单

**接口地址**: `POST /api/v1/distributor/order`

**接口说明**: 创建新订单

**请求头**:
```
X-API-Key: {api_key}
X-API-Secret: {api_secret}
```

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card | string | 是 | 卡片代码 |
| goods | array | 是 | 商品列表 |
| goods[].good_code | string | 是 | 商品代码 |
| goods[].good_name | string | 是 | 商品名称 |
| goods[].specs_code | string | 是 | 规格代码 |
| goods[].price | float | 是 | 商品价格 |
| goods[].num | int | 是 | 数量 |
| store_code | string | 是 | 店铺代码 |
| mobile | string | 是 | 手机号 |
| take_time | string | 是 | 取餐时间，如"立即自取" |
| outer_order_no | string | 是 | 外部订单号 |

**请求示例**:
```json
{
    "card": "CARD001",
    "goods": [
        {
            "good_code": "GOODS001",
            "good_name": "拿铁咖啡",
            "specs_code": "HOT_LARGE",
            "price": 25.0,
            "num": 1
        }
    ],
    "store_code": "STORE001",
    "mobile": "13800138000",
    "take_time": "立即自取",
    "outer_order_no": "TEST_ORDER_001"
}
```

**响应示例**:
```json
{
    "code": 200,
    "msg": "订单创建成功",
    "data": {
        "order_no": "ORD202401150001",
        "status": "pending",
        "total_amount": 25.00,
        "store_name": "瑞幸咖啡(朝阳门店)",
        "store_address": "北京市朝阳区朝阳门外大街1号"
    }
}
```

#### 2. 创建简化订单

**接口地址**: `POST /api/v1/distributor/order/simplified`

**接口说明**: 使用简化参数创建订单

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card | string | 是 | 卡片代码 |
| goods_list | array | 是 | 商品列表 |
| goods_list[].goods_code | string | 是 | 商品代码 |
| goods_list[].sku_code | string | 是 | SKU代码 |
| goods_list[].specs_code | string | 是 | 规格代码 |
| goods_list[].num | int | 是 | 数量 |
| store_code | string | 是 | 店铺代码 |
| mobile | string | 是 | 手机号 |
| outer_order_no | string | 是 | 外部订单号 |

**请求示例**:
```json
{
    "card": "CARD001",
    "goods_list": [
        {
            "goods_code": "GOODS001",
            "sku_code": "SKU001",
            "specs_code": "HOT_LARGE",
            "num": 1
        }
    ],
    "store_code": "STORE001",
    "mobile": "13800138000",
    "outer_order_no": "SIMPLE_ORDER_001"
}
```

#### 3. 批量创建订单

**接口地址**: `POST /api/v1/distributor/orders/batch`

**接口说明**: 批量创建订单，最多支持10个订单

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| orders | array | 是 | 订单列表，最多10个 |

每个订单的结构与单个创建订单相同。

**响应示例**:
```json
{
    "code": 200,
    "msg": "Batch order creation completed",
    "data": {
        "total": 2,
        "success": 2,
        "failed": 0,
        "results": [
            {
                "outer_order_no": "BATCH_ORDER_001",
                "success": true,
                "order_no": "ORD202401150001",
                "message": "订单创建成功"
            },
            {
                "outer_order_no": "BATCH_ORDER_002",
                "success": true,
                "order_no": "ORD202401150002",
                "message": "订单创建成功"
            }
        ]
    }
}
```

#### 4. 查询订单状态

**接口地址**: `GET /api/v1/distributor/order/:orderNo`

**接口说明**: 查询指定订单的状态

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| orderNo | string | 是 | 订单号 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "order_no": "ORD202401150001",
        "outer_order_no": "TEST_ORDER_001",
        "status": "completed",
        "luckin_status": "已完成",
        "luckin_order_no": "LK202401150001",
        "take_code": "A001",
        "total_amount": 25.00,
        "created_at": "2024-01-15T10:00:00Z",
        "completed_at": "2024-01-15T10:05:00Z"
    }
}
```

### 店铺和商品查询

#### 1. 搜索门店

**接口地址**: `GET /api/v1/distributor/stores`

**接口说明**: 根据卡片和条件搜索门店

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card | string | 是 | 卡片代码 |
| city_id | int | 否 | 城市ID |
| city_name | string | 否 | 城市名称 |
| keywords | string | 否 | 搜索关键词 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": [
        {
            "store_code": "STORE001",
            "store_name": "瑞幸咖啡(朝阳门店)",
            "store_address": "北京市朝阳区朝阳门外大街1号",
            "city_id": 1,
            "city_name": "北京",
            "longitude": 116.434,
            "latitude": 39.924,
            "distance": 1.2
        }
    ]
}
```

#### 2. 获取门店菜单

**接口地址**: `GET /api/v1/distributor/menu`

**接口说明**: 获取指定门店的商品菜单

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card | string | 是 | 卡片代码 |
| store_code | string | 是 | 店铺代码 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "categories": [
            {
                "category_name": "咖啡",
                "goods": [
                    {
                        "goods_code": "GOODS001",
                        "goods_name": "拿铁咖啡",
                        "goods_pic": "https://example.com/latte.jpg",
                        "price": 25.00,
                        "specs": [
                            {
                                "sku_code": "SKU001",
                                "sku_name": "大杯",
                                "specs": [
                                    {
                                        "spec_code": "HOT",
                                        "spec_name": "热"
                                    },
                                    {
                                        "spec_code": "ICE",
                                        "spec_name": "冰"
                                    }
                                ]
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
```

#### 3. 获取商品详情

**接口地址**: `GET /api/v1/distributor/goods`

**接口说明**: 获取指定商品的详细信息

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| card | string | 是 | 卡片代码 |
| store_code | string | 是 | 店铺代码 |
| goods_code | string | 是 | 商品代码 |

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "goods_code": "GOODS001",
        "goods_name": "拿铁咖啡",
        "goods_pic": "https://example.com/latte.jpg",
        "description": "经典意式拿铁，醇厚浓郁",
        "price": 25.00,
        "skus": [
            {
                "sku_code": "SKU001",
                "sku_name": "大杯",
                "price": 25.00,
                "specs": [
                    {
                        "spec_type": "温度",
                        "options": [
                            {
                                "spec_code": "HOT",
                                "spec_name": "热"
                            },
                            {
                                "spec_code": "ICE",
                                "spec_name": "冰"
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
```

#### 4. 获取城市列表

**接口地址**: `GET /api/v1/distributor/cities`

**接口说明**: 获取所有支持的城市列表

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": [
        {
            "city_id": 1,
            "city_name": "北京"
        },
        {
            "city_id": 2,
            "city_name": "上海"
        },
        {
            "city_id": 3,
            "city_name": "广州"
        }
    ]
}
```

#### 5. 获取可用商品列表

**接口地址**: `GET /api/v1/distributor/available-products`

**接口说明**: 获取当前分销商可下单的商品信息

**响应示例**:
```json
{
    "code": 200,
    "msg": "成功",
    "data": {
        "categories": [
            {
                "category": "咖啡",
                "products": [
                    {
                        "goods_code": "GOODS001",
                        "goods_name": "拿铁咖啡"
                    },
                    {
                        "goods_code": "GOODS002",
                        "goods_name": "美式咖啡"
                    }
                ]
            }
        ],
        "total_count": 50
    }
}
```

### 分销商财务管理

#### 1. 获取余额信息

**接口地址**: `GET /api/v1/distributor/balance`

**接口说明**: 获取当前分销商的余额信息

**响应示例**:
```json
{
    "balance": 5000.00,
    "credit_limit": 10000.00,
    "available_balance": 15000.00,
    "frozen_amount": 200.00,
    "warning_balance": 1000.00,
    "warning_enabled": true
}
```

#### 2. 获取交易记录

**接口地址**: `GET /api/v1/distributor/transactions`

**接口说明**: 获取分销商自己的交易记录

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |

**响应示例**:
```json
{
    "data": [
        {
            "id": 100,
            "type": "order",
            "amount": -25.00,
            "balance_before": 5025.00,
            "balance_after": 5000.00,
            "order_no": "ORD202401150001",
            "remark": "订单消费",
            "created_at": "2024-01-15T10:00:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "page_size": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

#### 3. 创建提现申请

**接口地址**: `POST /api/v1/distributor/withdrawal`

**接口说明**: 提交提现申请

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| amount | float | 是 | 提现金额 |
| account_info | object | 是 | 账户信息 |
| account_info.type | string | 是 | 账户类型：bank/alipay/wechat |
| account_info.account_name | string | 是 | 账户名称 |
| account_info.account_no | string | 是 | 账号 |
| account_info.bank_name | string | 否 | 银行名称（银行卡必填） |
| account_info.bank_branch | string | 否 | 开户支行（银行卡必填） |
| remark | string | 否 | 备注 |

**请求示例**:
```json
{
    "amount": 1000.0,
    "account_info": {
        "type": "bank",
        "account_name": "张三",
        "account_no": "6222021234567890123",
        "bank_name": "工商银行",
        "bank_branch": "北京朝阳支行"
    },
    "remark": "月度提现"
}
```

**响应示例**:
```json
{
    "message": "提现申请已提交",
    "data": {
        "withdrawal_id": 10,
        "amount": 1000.00,
        "status": "pending",
        "created_at": "2024-01-15T15:00:00Z"
    }
}
```

#### 4. 更新预警设置

**接口地址**: `PUT /api/v1/distributor/warning-settings`

**接口说明**: 更新余额预警设置

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| warning_balance | float | 是 | 预警余额阈值 |
| warning_enabled | bool | 是 | 是否启用预警 |
| warning_email | string | 否 | 预警邮箱 |
| warning_webhook | string | 否 | 预警Webhook地址 |

**请求示例**:
```json
{
    "warning_balance": 1000.0,
    "warning_enabled": true,
    "warning_email": "alert@example.com",
    "warning_webhook": "https://example.com/webhook/balance-warning"
}
```

**响应示例**:
```json
{
    "message": "预警设置更新成功"
}
```

## 错误响应示例

当API调用失败时，会返回相应的错误信息：

```json
{
    "code": 400,
    "msg": "请求参数无效：缺少必填字段 'card'",
    "data": null
}
```

```json
{
    "code": 401,
    "msg": "未授权：无效的API密钥",
    "data": null
}
```

```json
{
    "code": 500,
    "msg": "服务器内部错误：数据库连接失败",
    "data": null
}
```

## 注意事项

1. **认证要求**：除了登录接口外，其他所有接口都需要提供有效的认证信息
2. **频率限制**：分销商API有频率限制，基础限制为每分钟300次请求
3. **数据格式**：所有请求和响应都使用JSON格式，请确保Content-Type设置为application/json
4. **时间格式**：所有时间字段都使用ISO 8601格式（如：2024-01-15T10:00:00Z）
5. **金额格式**：所有金额字段都使用浮点数，精确到小数点后两位
6. **状态码**：请根据HTTP状态码和响应中的code字段判断请求是否成功
7. **幂等性**：创建订单时请使用唯一的outer_order_no，系统会检查重复订单
8. **异步处理**：批量操作可能采用异步处理，请根据返回的任务ID查询进度

## 开发建议

1. **错误处理**：始终检查响应的code字段，不要仅依赖HTTP状态码
2. **重试机制**：对于网络错误或5xx错误，建议实现指数退避的重试机制
3. **日志记录**：记录所有API调用的请求和响应，便于问题排查
4. **安全存储**：妥善保管API密钥，不要在客户端代码中硬编码
5. **版本兼容**：关注API版本变化，及时更新集成代码
6. **测试环境**：在生产环境部署前，充分使用测试环境进行验证
7. **监控告警**：监控API调用的成功率和响应时间，设置合理的告警阈值