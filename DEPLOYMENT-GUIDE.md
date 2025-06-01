# 瑞幸咖啡管理系统部署指南

本文档详细说明了如何在宝塔面板和Linux环境下部署瑞幸咖啡管理系统。

## 目录

- [系统要求](#系统要求)
- [部署架构](#部署架构)
- [宝塔面板部署](#宝塔面板部署)
  - [环境准备](#宝塔环境准备)
  - [MySQL数据库配置](#宝塔mysql配置)
  - [Redis配置](#宝塔redis配置)
  - [后端部署](#宝塔后端部署)
  - [前端部署](#宝塔前端部署)
  - [Nginx配置](#宝塔nginx配置)
- [Linux环境部署](#linux环境部署)
  - [环境准备](#linux环境准备)
  - [手动部署](#手动部署)
  - [Docker部署](#docker部署)
  - [系统服务配置](#系统服务配置)
- [部署后配置](#部署后配置)
- [性能优化](#性能优化)
- [安全配置](#安全配置)
- [常见问题](#常见问题)
- [故障排除](#故障排除)

## 系统要求

### 最低配置
- CPU: 2核
- 内存: 4GB
- 硬盘: 50GB
- 带宽: 5Mbps

### 推荐配置
- CPU: 4核
- 内存: 8GB
- 硬盘: 100GB SSD
- 带宽: 10Mbps

### 软件要求
- 操作系统: Ubuntu 20.04+ / CentOS 7+ / Debian 10+
- MySQL: 8.0+
- Redis: 6.0+
- Nginx: 1.18+
- Go: 1.21+
- Node.js: 18+

## 部署架构

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Nginx     │────▶│   Backend   │────▶│   MySQL     │
│   (80/443)  │     │   (8080)    │     │   (3306)    │
└─────────────┘     └─────────────┘     └─────────────┘
       │                    │
       │                    └───────────▶┌─────────────┐
       │                                 │   Redis     │
       ▼                                 │   (6379)    │
┌─────────────┐                         └─────────────┘
│  Frontend   │
│   (Static)  │
└─────────────┘
```

## 宝塔面板部署

### 宝塔环境准备

1. **安装宝塔面板**
```bash
# Centos安装脚本
yum install -y wget && wget -O install.sh https://download.bt.cn/install/install_6.0.sh && sh install.sh ed8484bec

# Ubuntu/Deepin安装脚本
wget -O install.sh https://download.bt.cn/install/install-ubuntu_6.0.sh && sudo bash install.sh ed8484bec

# Debian安装脚本
wget -O install.sh https://download.bt.cn/install/install-ubuntu_6.0.sh && bash install.sh ed8484bec
```

2. **登录宝塔面板**
   - 访问: http://服务器IP:8888
   - 使用安装时生成的用户名和密码登录

3. **安装必要软件**
   - 在软件商店安装：
     - Nginx 1.22
     - MySQL 8.0
     - Redis 7.0
     - PHP 8.1（可选，用于phpMyAdmin）
     - Go 1.21

### 宝塔MySQL配置

1. **创建数据库**
   - 进入数据库管理
   - 点击"添加数据库"
   - 数据库名：`ruixin_platform`
   - 用户名：`ruixin`
   - 密码：`RuiXin@2025!`（请修改为强密码）
   - 访问权限：`本地服务器`

2. **优化MySQL配置**
   - 点击MySQL设置 > 配置修改
   - 添加以下配置：
```ini
[mysqld]
# 基础设置
character-set-server=utf8mb4
collation-server=utf8mb4_general_ci
max_connections=500
max_connect_errors=1000

# InnoDB优化
innodb_buffer_pool_size=1G  # 根据内存调整，建议设置为总内存的50%-70%
innodb_log_file_size=256M
innodb_flush_log_at_trx_commit=2
innodb_flush_method=O_DIRECT

# 查询缓存
query_cache_type=1
query_cache_size=128M
query_cache_limit=2M

# 慢查询日志
slow_query_log=1
slow_query_log_file=/www/server/data/mysql-slow.log
long_query_time=2
```

3. **导入初始数据**
```bash
# 创建初始化SQL文件
cat > /tmp/init.sql << 'EOF'
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS ruixin_platform DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE ruixin_platform;

-- 初始管理员账号
INSERT INTO admins (username, password, email, real_name, role, status, created_at) 
VALUES ('admin', '$2a$10$YourHashedPasswordHere', 'admin@ruixin.com', '系统管理员', 'super_admin', 1, NOW())
ON DUPLICATE KEY UPDATE username=username;
EOF

# 导入数据
mysql -uroot -p ruixin_platform < /tmp/init.sql
```

### 宝塔Redis配置

1. **配置Redis**
   - 点击Redis设置 > 配置文件
   - 修改以下配置：
```conf
# 绑定地址
bind 127.0.0.1

# 端口
port 6379

# 密码设置（建议设置密码）
requirepass RuiXinRedis@2025

# 持久化配置
save 900 1
save 300 10
save 60 10000

# 最大内存设置
maxmemory 512mb
maxmemory-policy allkeys-lru

# 日志
logfile /www/server/redis/redis.log
```

2. **重启Redis服务**
```bash
systemctl restart redis
```

### 宝塔后端部署

1. **创建网站**
   - 点击"网站" > "添加站点"
   - 域名：`api.yourdomain.com`（或使用IP:8080）
   - 根目录：`/www/wwwroot/ruixin-backend`
   - 数据库：选择之前创建的数据库

2. **上传后端代码**
```bash
# 创建目录
mkdir -p /www/wwwroot/ruixin-backend
cd /www/wwwroot/ruixin-backend

# 上传代码（使用宝塔文件管理器或以下命令）
# 方式1：Git克隆
git clone https://github.com/your-repo/ruixin-backend.git .

# 方式2：上传压缩包并解压
# 通过宝塔文件管理器上传backend.zip，然后解压
```

3. **编译后端程序**
```bash
cd /www/wwwroot/ruixin-backend

# 设置Go代理（加速依赖下载）
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 下载依赖
go mod download

# 编译程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ruixin-api cmd/api/main.go

# 设置执行权限
chmod +x ruixin-api
```

4. **配置文件设置**
```bash
# 编辑配置文件
vim /www/wwwroot/ruixin-backend/configs/config.yaml
```

修改以下配置：
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # 生产环境使用release
  read_timeout: 30
  write_timeout: 30

database:
  driver: "mysql"
  host: "127.0.0.1"  # 本地数据库
  port: 3306
  database: "ruixin_platform"
  username: "ruixin"
  password: "RuiXin@2025!"  # 修改为实际密码
  charset: "utf8mb4"
  max_idle: 10
  max_open: 100

redis:
  host: "127.0.0.1"
  port: 6379
  password: "RuiXinRedis@2025"  # 修改为实际密码
  db: 0

jwt:
  secret: "your-super-secret-jwt-key-change-this-in-production"
  expiration: 86400

luckin:
  base_url: "https://lkcoffe.cn"
  timeout: 30
  max_retries: 3

log:
  level: "info"
  filename: "/www/wwwroot/ruixin-backend/logs/ruixin.log"
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
```

5. **创建系统服务**
```bash
# 创建服务文件
cat > /etc/systemd/system/ruixin-backend.service << EOF
[Unit]
Description=RuiXin Backend API Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www
Group=www
WorkingDirectory=/www/wwwroot/ruixin-backend
ExecStart=/www/wwwroot/ruixin-backend/ruixin-api
Restart=always
RestartSec=5
StandardOutput=append:/www/wwwroot/ruixin-backend/logs/service.log
StandardError=append:/www/wwwroot/ruixin-backend/logs/service-error.log

# 环境变量
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
systemctl daemon-reload
systemctl enable ruixin-backend
systemctl start ruixin-backend

# 查看服务状态
systemctl status ruixin-backend
```

### 宝塔前端部署

1. **创建前端网站**
   - 点击"网站" > "添加站点"
   - 域名：`www.yourdomain.com`
   - 根目录：`/www/wwwroot/ruixin-frontend`
   - PHP版本：纯静态

2. **构建前端代码**
```bash
# 在本地或构建服务器上
cd frontend
npm install
npm run build

# 将dist目录下的文件上传到服务器
# 使用宝塔文件管理器上传到 /www/wwwroot/ruixin-frontend
```

3. **设置文件权限**
```bash
chown -R www:www /www/wwwroot/ruixin-frontend
chmod -R 755 /www/wwwroot/ruixin-frontend
```

### 宝塔Nginx配置

1. **配置前端网站Nginx**
   - 点击网站设置 > 配置文件
   - 修改配置如下：

```nginx
server {
    listen 80;
    server_name www.yourdomain.com;
    root /www/wwwroot/ruixin-frontend;
    index index.html;
    
    # SSL配置（如果有证书）
    # listen 443 ssl http2;
    # ssl_certificate /www/server/panel/vhost/cert/yourdomain.com/fullchain.pem;
    # ssl_certificate_key /www/server/panel/vhost/cert/yourdomain.com/privkey.pem;
    # ssl_protocols TLSv1.2 TLSv1.3;
    # ssl_ciphers HIGH:!aNULL:!MD5;
    
    # 启用 gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/json image/svg+xml;
    
    # Vue Router History模式
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API代理
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 代理超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # 代理缓冲
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
        
        # WebSocket支持（如果需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # HTML文件不缓存
    location ~* \.html$ {
        expires -1;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # 访问日志
    access_log /www/wwwlogs/ruixin-frontend.log;
    error_log /www/wwwlogs/ruixin-frontend.error.log;
}

# HTTP重定向到HTTPS（如果启用SSL）
# server {
#     listen 80;
#     server_name www.yourdomain.com;
#     return 301 https://$server_name$request_uri;
# }
```

2. **设置防火墙规则**
   - 在宝塔面板 > 安全 中开放以下端口：
   - 80 (HTTP)
   - 443 (HTTPS，如果使用)
   - 8080 (后端API，仅限内网访问)

## Linux环境部署

### Linux环境准备

1. **更新系统**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# CentOS/RHEL
sudo yum update -y
```

2. **安装必要软件**
```bash
# Ubuntu/Debian
sudo apt install -y git wget curl vim build-essential

# CentOS/RHEL
sudo yum install -y git wget curl vim gcc gcc-c++ make
```

3. **安装MySQL 8.0**
```bash
# Ubuntu 20.04+
sudo apt install -y mysql-server
sudo mysql_secure_installation

# CentOS 7/8
sudo yum install -y https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
sudo yum install -y mysql-community-server
sudo systemctl start mysqld
sudo mysql_secure_installation
```

4. **安装Redis**
```bash
# Ubuntu/Debian
sudo apt install -y redis-server

# CentOS/RHEL
sudo yum install -y epel-release
sudo yum install -y redis

# 启动Redis
sudo systemctl start redis
sudo systemctl enable redis
```

5. **安装Nginx**
```bash
# Ubuntu/Debian
sudo apt install -y nginx

# CentOS/RHEL
sudo yum install -y nginx

# 启动Nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

6. **安装Go语言环境**
```bash
# 下载Go 1.21
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# 解压到/usr/local
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

7. **安装Node.js 18**
```bash
# 使用NodeSource仓库
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# 或使用nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 18
nvm use 18
```

### 手动部署

1. **创建部署目录**
```bash
sudo mkdir -p /opt/ruixin/{backend,frontend,logs,data}
sudo chown -R $USER:$USER /opt/ruixin
```

2. **部署后端**
```bash
# 克隆代码
cd /opt/ruixin/backend
git clone https://github.com/your-repo/ruixin-backend.git .

# 配置数据库
mysql -u root -p << EOF
CREATE DATABASE IF NOT EXISTS ruixin_platform DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE USER 'ruixin'@'localhost' IDENTIFIED BY 'RuiXin@2025!';
GRANT ALL PRIVILEGES ON ruixin_platform.* TO 'ruixin'@'localhost';
FLUSH PRIVILEGES;
EOF

# 修改配置文件
cp configs/config.yaml.example configs/config.yaml
vim configs/config.yaml  # 修改数据库密码等配置

# 编译后端
go mod download
go build -o ruixin-api cmd/api/main.go
chmod +x ruixin-api

# 创建systemd服务
sudo tee /etc/systemd/system/ruixin-backend.service > /dev/null << EOF
[Unit]
Description=RuiXin Backend API Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=$USER
WorkingDirectory=/opt/ruixin/backend
ExecStart=/opt/ruixin/backend/ruixin-api
Restart=always
RestartSec=5
StandardOutput=append:/opt/ruixin/logs/backend.log
StandardError=append:/opt/ruixin/logs/backend-error.log
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable ruixin-backend
sudo systemctl start ruixin-backend
```

3. **部署前端**
```bash
# 构建前端（在开发机器上）
cd frontend
npm install
npm run build

# 上传dist目录到服务器
scp -r dist/* user@server:/opt/ruixin/frontend/

# 配置Nginx
sudo tee /etc/nginx/sites-available/ruixin > /dev/null << 'EOF'
server {
    listen 80;
    server_name yourdomain.com;
    root /opt/ruixin/frontend;
    index index.html;
    
    # gzip压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    
    # Vue Router
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API代理
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

# 启用网站配置
sudo ln -s /etc/nginx/sites-available/ruixin /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Docker部署

1. **安装Docker和Docker Compose**
```bash
# 安装Docker
curl -fsSL https://get.docker.com | bash
sudo usermod -aG docker $USER

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

2. **准备部署文件**
```bash
cd /opt/ruixin
# 复制docker-compose.yml和相关文件
```

3. **配置环境变量**
```bash
# 创建.env文件
cat > .env << EOF
# MySQL配置
MYSQL_ROOT_PASSWORD=RootPassword@2025
MYSQL_DATABASE=ruixin_platform
MYSQL_USER=ruixin
MYSQL_PASSWORD=RuiXin@2025!

# Redis配置
REDIS_PASSWORD=RuiXinRedis@2025

# JWT密钥
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# 其他配置
TZ=Asia/Shanghai
EOF
```

4. **修改docker-compose.yml**
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: ruixin-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE}
      MYSQL_USER: ${MYSQL_USER}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
      TZ: ${TZ}
    ports:
      - "127.0.0.1:3306:3306"
    volumes:
      - ./data/mysql:/var/lib/mysql
      - ./backend/scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - ruixin-network
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci

  redis:
    image: redis:7-alpine
    container_name: ruixin-redis
    restart: always
    command: redis-server --requirepass ${REDIS_PASSWORD}
    ports:
      - "127.0.0.1:6379:6379"
    volumes:
      - ./data/redis:/data
    networks:
      - ruixin-network

  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    container_name: ruixin-backend
    restart: always
    depends_on:
      - mysql
      - redis
    environment:
      - DB_HOST=mysql
      - DB_PASSWORD=${MYSQL_PASSWORD}
      - REDIS_HOST=redis
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - GIN_MODE=release
    ports:
      - "127.0.0.1:8080:8080"
    volumes:
      - ./backend/configs:/root/configs
      - ./logs/backend:/root/logs
    networks:
      - ruixin-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: ruixin-frontend
    restart: always
    depends_on:
      - backend
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/ssl:/etc/nginx/ssl
      - ./logs/nginx:/var/log/nginx
    networks:
      - ruixin-network

networks:
  ruixin-network:
    driver: bridge

volumes:
  mysql_data:
  redis_data:
```

5. **启动服务**
```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据
docker-compose down -v
```

### 系统服务配置

1. **配置自动启动**
```bash
# Docker Compose服务
sudo tee /etc/systemd/system/ruixin-docker.service > /dev/null << EOF
[Unit]
Description=RuiXin Platform Docker Compose
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/ruixin
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable ruixin-docker
```

2. **配置日志轮转**
```bash
# 创建logrotate配置
sudo tee /etc/logrotate.d/ruixin > /dev/null << EOF
/opt/ruixin/logs/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload ruixin-backend > /dev/null 2>&1 || true
    endscript
}
EOF
```

## 部署后配置

### 初始化管理员账号

1. **生成密码哈希**
```bash
# 使用Go生成bcrypt密码
cd /opt/ruixin/backend
go run scripts/hash_password.go "your-admin-password"
```

2. **插入管理员数据**
```sql
-- 连接数据库
mysql -u ruixin -p ruixin_platform

-- 插入管理员
INSERT INTO admins (username, password, email, real_name, role, status, created_at) 
VALUES (
    'admin', 
    '$2a$10$YourGeneratedHashHere',  -- 替换为生成的哈希
    'admin@ruixin.com', 
    '系统管理员', 
    'super_admin', 
    1, 
    NOW()
);
```

### 配置定时任务

```bash
# 编辑crontab
crontab -e

# 添加定时任务
# 每天凌晨2点清理过期卡片
0 2 * * * curl -X POST http://localhost:8080/api/v1/internal/cleanup-expired-cards

# 每小时同步订单状态
0 * * * * curl -X POST http://localhost:8080/api/v1/internal/sync-order-status

# 每天凌晨3点备份数据库
0 3 * * * /opt/ruixin/scripts/backup.sh
```

### 配置备份脚本

```bash
# 创建备份脚本
mkdir -p /opt/ruixin/scripts
cat > /opt/ruixin/scripts/backup.sh << 'EOF'
#!/bin/bash
# 数据库备份脚本

BACKUP_DIR="/opt/ruixin/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="ruixin_platform"
DB_USER="ruixin"
DB_PASS="RuiXin@2025!"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u$DB_USER -p$DB_PASS $DB_NAME | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz

# 保留最近7天的备份
find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +7 -exec rm {} \;

# 如果使用云存储，可以上传到OSS/S3
# aws s3 cp $BACKUP_DIR/db_backup_$DATE.sql.gz s3://your-bucket/backups/
EOF

chmod +x /opt/ruixin/scripts/backup.sh
```

## 性能优化

### MySQL优化

1. **创建索引**
```sql
-- 订单表索引
ALTER TABLE orders ADD INDEX idx_distributor_created (distributor_id, created_at);
ALTER TABLE orders ADD INDEX idx_status_created (status, created_at);
ALTER TABLE orders ADD INDEX idx_outer_order_no (outer_order_no);

-- 卡片表索引
ALTER TABLE cards ADD INDEX idx_status_price (status, price_id);
ALTER TABLE cards ADD INDEX idx_card_code (card_code);

-- 交易表索引
ALTER TABLE transactions ADD INDEX idx_distributor_created (distributor_id, created_at);
```

2. **MySQL配置优化**
```ini
[mysqld]
# 连接数优化
max_connections = 500
max_connect_errors = 1000

# 缓冲池优化（根据服务器内存调整）
innodb_buffer_pool_size = 2G
innodb_buffer_pool_instances = 2

# 日志优化
innodb_log_file_size = 512M
innodb_log_buffer_size = 64M

# 查询缓存
query_cache_type = 1
query_cache_size = 256M
query_cache_limit = 4M

# 临时表优化
tmp_table_size = 256M
max_heap_table_size = 256M
```

### Redis优化

```conf
# 最大内存设置
maxmemory 1gb
maxmemory-policy allkeys-lru

# 持久化优化
save 900 1
save 300 10
save 60 10000

# 慢查询日志
slowlog-log-slower-than 10000
slowlog-max-len 128

# TCP优化
tcp-backlog 511
tcp-keepalive 300
```

### Nginx优化

```nginx
# 工作进程优化
worker_processes auto;
worker_rlimit_nofile 65535;

events {
    worker_connections 4096;
    use epoll;
    multi_accept on;
}

http {
    # 基础优化
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    
    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss application/rss+xml application/atom+xml image/svg+xml;
    
    # 缓存优化
    open_file_cache max=2000 inactive=20s;
    open_file_cache_valid 30s;
    open_file_cache_min_uses 2;
    open_file_cache_errors on;
    
    # 限流配置
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;
}
```

## 安全配置

### 防火墙配置

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Firewalld (CentOS)
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### SSL证书配置

1. **使用Let's Encrypt**
```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# 自动续期
sudo systemctl enable certbot.timer
```

2. **Nginx SSL配置**
```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # SSL优化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
}
```

### 安全加固

1. **限制SSH访问**
```bash
# 修改SSH配置
sudo vim /etc/ssh/sshd_config

# 禁用root登录
PermitRootLogin no

# 禁用密码登录，使用密钥
PasswordAuthentication no

# 修改默认端口
Port 2222

# 重启SSH服务
sudo systemctl restart sshd
```

2. **配置fail2ban**
```bash
# 安装fail2ban
sudo apt install fail2ban

# 创建配置
sudo tee /etc/fail2ban/jail.local > /dev/null << EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true

[nginx-limit-req]
enabled = true
filter = nginx-limit-req
action = iptables-multiport[name=nginx-limit-req,port="80,443"]
logpath = /var/log/nginx/error.log
EOF

sudo systemctl restart fail2ban
```

## 常见问题

### Q1: 后端启动失败，提示数据库连接错误
**解决方案**：
1. 检查MySQL服务是否启动
2. 验证数据库用户名密码是否正确
3. 确认数据库是否已创建
4. 检查防火墙是否阻止了3306端口

```bash
# 检查MySQL状态
systemctl status mysql

# 测试连接
mysql -u ruixin -p -h localhost ruixin_platform
```

### Q2: 前端页面显示空白
**解决方案**：
1. 检查Nginx配置是否正确
2. 确认静态文件是否上传完整
3. 查看浏览器控制台错误信息
4. 检查API代理是否正常工作

```bash
# 检查Nginx配置
nginx -t

# 查看Nginx错误日志
tail -f /var/log/nginx/error.log
```

### Q3: API请求返回502错误
**解决方案**：
1. 检查后端服务是否运行
2. 查看后端日志
3. 确认端口是否被占用
4. 检查代理配置

```bash
# 检查后端服务
systemctl status ruixin-backend

# 查看端口占用
netstat -tlnp | grep 8080

# 查看后端日志
journalctl -u ruixin-backend -f
```

### Q4: Redis连接失败
**解决方案**：
1. 检查Redis服务状态
2. 验证Redis密码配置
3. 确认绑定地址设置

```bash
# 测试Redis连接
redis-cli -a yourpassword ping
```

## 故障排除

### 日志位置

- **后端日志**: `/opt/ruixin/logs/backend.log`
- **Nginx访问日志**: `/var/log/nginx/access.log`
- **Nginx错误日志**: `/var/log/nginx/error.log`
- **MySQL日志**: `/var/log/mysql/error.log`
- **系统日志**: `journalctl -xe`

### 性能监控

1. **安装监控工具**
```bash
# 安装htop
sudo apt install htop

# 安装iotop
sudo apt install iotop

# 安装nethogs
sudo apt install nethogs
```

2. **监控命令**
```bash
# CPU和内存使用
htop

# 磁盘IO
iotop

# 网络流量
nethogs

# MySQL进程
mysqladmin -u root -p processlist

# Redis监控
redis-cli monitor
```

### 紧急恢复

1. **数据库恢复**
```bash
# 从备份恢复
gunzip < /opt/ruixin/backups/db_backup_20240115.sql.gz | mysql -u root -p ruixin_platform
```

2. **快速重启所有服务**
```bash
# 使用systemctl
sudo systemctl restart mysql redis nginx ruixin-backend

# 使用Docker
cd /opt/ruixin
docker-compose restart
```

3. **清理日志空间**
```bash
# 清理大日志文件
find /opt/ruixin/logs -name "*.log" -size +100M -exec truncate -s 0 {} \;

# 清理Docker日志
docker system prune -a
```

## 维护建议

1. **定期备份**
   - 每天备份数据库
   - 每周备份配置文件
   - 每月测试恢复流程

2. **监控告警**
   - 设置磁盘空间告警（剩余<20%）
   - 监控服务状态
   - 设置API响应时间告警

3. **安全更新**
   - 定期更新系统包
   - 及时修复安全漏洞
   - 定期审查日志

4. **性能优化**
   - 定期分析慢查询
   - 优化数据库索引
   - 清理过期数据

5. **容量规划**
   - 监控业务增长
   - 提前扩容
   - 优化资源使用

---

**重要提示**：
- 请根据实际情况修改所有密码和敏感信息
- 生产环境必须启用HTTPS
- 定期备份数据并测试恢复流程
- 建议使用专业运维工具进行监控和管理