# 瑞鑫电商管理系统 - 宝塔面板快速部署教程

本教程适用于已安装宝塔面板、Nginx、Redis、MySQL、Golang的服务器环境。

## 一、环境确认

确保以下软件已安装：
- 宝塔面板 7.0+
- Nginx 1.20+
- MySQL 5.7+
- Redis 5.0+
- Go 1.19+
- Node.js 16+ （用于前端打包）
- PM2 （进程管理）

```bash
# 验证安装
go version
node -v
pm2 -v
```

## 二、项目准备

### 1. 上传项目代码
将项目代码上传到服务器目录：
```bash
/www/wwwroot/ruixin
```

### 2. 创建数据库
在宝塔面板中创建数据库：
- 数据库名：`ruixin_platform`
- 用户名：`ruixin`
- 密码：设置强密码

或使用命令行：
```bash
mysql -uroot -p
CREATE DATABASE ruixin_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'ruixin'@'localhost' IDENTIFIED BY '你的密码';
GRANT ALL PRIVILEGES ON ruixin_platform.* TO 'ruixin'@'localhost';
FLUSH PRIVILEGES;
exit;
```

## 三、后端部署

### 1. 进入后端目录
```bash
cd /www/wwwroot/ruixin/backend
```

### 2. 配置文件设置
编辑 `configs/config.yaml`：
```yaml
server:
  host: "0.0.0.0"
  port: 8080  # API端口
  mode: "release"  # 生产环境

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  database: "ruixin_platform"
  username: "ruixin"
  password: "你的数据库密码"
  charset: "utf8mb4"

redis:
  host: "localhost"
  port: 6379
  password: ""  # 如果Redis设置了密码，填写这里
  db: 0

jwt:
  secret: "生成一个32位以上的随机密钥"
  expiration: 86400

log:
  level: "info"
  filename: "logs/ruixin.log"
```

### 3. 初始化数据库
```bash
# 如果有数据库文件
mysql -u ruixin -p ruixin_platform < database.sql

# 或运行初始化脚本
cd scripts
go run init_db.go
cd ..
```

### 4. 编译后端程序
```bash
# 安装依赖
go mod download

# 编译
go build -o ruixin-api cmd/api/main.go

# 添加执行权限
chmod +x ruixin-api
```

### 5. 使用PM2管理后端服务
创建PM2配置文件 `ecosystem.config.js`：
```javascript
module.exports = {
  apps: [{
    name: 'ruixin-api',
    script: './ruixin-api',
    cwd: '/www/wwwroot/ruixin/backend',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '1G',
    env: {
      NODE_ENV: 'production'
    }
  }]
}
```

启动服务：
```bash
pm2 start ecosystem.config.js
pm2 save
pm2 startup  # 设置开机自启
```

## 四、前端部署

### 1. 进入前端目录
```bash
cd /www/wwwroot/ruixin/frontend
```

### 2. 创建生产环境配置
创建 `.env.production` 文件：
```env
VITE_API_BASE_URL=http://你的域名:8081/api
VITE_APP_TITLE=瑞鑫电商管理系统
```

### 3. 安装依赖并打包
```bash
# 安装依赖
npm install

# 打包生产版本
npm run build
```

打包完成后，`dist` 目录包含所有静态文件。

## 五、Nginx配置

### 1. 在宝塔面板创建网站
- 域名：你的域名
- 端口：8081（避免与其他项目的80端口冲突）
- 根目录：`/www/wwwroot/ruixin/frontend/dist`
- PHP版本：纯静态

### 2. 修改Nginx配置
在网站设置中，修改配置文件：

```nginx
server {
    listen 8081;
    server_name 你的域名;
    
    # 前端静态文件
    root /www/wwwroot/ruixin/frontend/dist;
    index index.html;
    
    # GZIP压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
    
    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API代理
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 600;
        proxy_send_timeout 600;
        proxy_read_timeout 600;
        
        # 文件上传限制
        client_max_body_size 50M;
    }
    
    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # 日志
    access_log /www/wwwlogs/ruixin.log;
    error_log /www/wwwlogs/ruixin.error.log;
}
```

### 3. 重载Nginx
```bash
nginx -t  # 测试配置
nginx -s reload  # 重载
```

## 六、部署验证

### 1. 检查服务状态
```bash
# 查看后端服务
pm2 status
pm2 logs ruixin-api

# 测试API
curl http://localhost:8080/api/health
```

### 2. 访问系统
- 前端地址：`http://你的域名:8081`
- 默认账号：admin / admin123

### 3. 防火墙设置
如果无法访问，检查防火墙：
```bash
# 开放端口
firewall-cmd --permanent --add-port=8081/tcp
firewall-cmd --reload

# 或在宝塔面板的安全中放行端口
```

## 七、常见问题

### 1. 后端服务启动失败
```bash
# 查看详细日志
pm2 logs ruixin-api --lines 100

# 检查端口占用
netstat -tlnp | grep 8080
```

### 2. 前端页面空白
- 检查API地址配置是否正确
- 查看浏览器控制台错误
- 检查Nginx错误日志

### 3. 数据库连接失败
- 确认数据库用户名密码正确
- 检查MySQL服务状态
- 验证数据库权限

## 八、日常维护

### 1. 查看日志
```bash
# PM2日志
pm2 logs ruixin-api

# 应用日志
tail -f /www/wwwroot/ruixin/backend/logs/ruixin.log
```

### 2. 重启服务
```bash
pm2 restart ruixin-api
```

### 3. 更新部署
```bash
# 后端更新
cd /www/wwwroot/ruixin/backend
go build -o ruixin-api cmd/api/main.go
pm2 restart ruixin-api

# 前端更新
cd /www/wwwroot/ruixin/frontend
npm run build
```

### 4. 备份建议
定期备份：
- 数据库：每日备份
- 配置文件：有修改时备份
- 上传文件：每周备份

---

部署完成！如有问题，请查看日志或联系技术支持。