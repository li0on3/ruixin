# 瑞幸咖啡分销商自动化系统 - 后端服务

## 项目概述

瑞幸咖啡分销商自动化系统后端服务，提供分销商API接口和管理后台API。

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM (MySQL)
- Redis
- JWT Authentication
- Zap Logger

## 项目结构

```
backend/
├── cmd/api/          # 应用入口
├── internal/         # 内部包
│   ├── api/         # HTTP处理器和中间件
│   ├── config/      # 配置管理
│   ├── models/      # 数据模型
│   ├── repository/  # 数据访问层
│   └── services/    # 业务逻辑层
├── pkg/             # 公共包
│   ├── httpclient/  # HTTP客户端
│   └── logger/      # 日志工具
├── configs/         # 配置文件
└── scripts/         # 脚本文件
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+

### 安装依赖

```bash
go mod download
```

### 配置文件

复制配置文件模板并修改：

```bash
cp configs/config.yaml configs/config.local.yaml
```

### 运行服务

```bash
# 开发模式
make dev

# 生产模式
make build
./ruixin-api
```

## API文档

### 分销商API

所有分销商API需要在Header中提供认证信息：
- `X-API-Key`: API密钥
- `X-API-Secret`: API密钥

#### 创建订单

```
POST /api/v1/distributor/order
```

请求体：
```json
{
  "card_code": "4H8SB644",
  "store_code": "387705",
  "phone_number": "13800138000",
  "callback_url": "https://your-callback-url.com",
  "goods": [
    {
      "goods_code": "4929",
      "sku_code": "SP3349-00009",
      "quantity": 1,
      "specs": [
        {
          "specs_code": "64",
          "code": "0",
          "name": "大杯 16oz"
        }
      ]
    }
  ]
}
```

#### 查询订单

```
GET /api/v1/distributor/order/{orderNo}
```

#### 搜索门店

```
GET /api/v1/distributor/stores?card={card}&city_id={cityId}&keywords={keywords}
```

#### 获取菜单

```
GET /api/v1/distributor/menu?card={card}&store_code={storeCode}
```

#### 获取商品详情

```
GET /api/v1/distributor/goods?card={card}&store_code={storeCode}&goods_code={goodsCode}
```

## 数据库迁移

```bash
# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down
```

## Docker部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

## 测试

```bash
make test
```

## 代码规范

```bash
# 格式化代码
make fmt

# 运行linter
make lint
```