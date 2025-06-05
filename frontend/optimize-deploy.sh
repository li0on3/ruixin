#!/bin/bash
# 前端性能优化部署脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== 瑞鑫电商系统前端性能优化部署 ===${NC}"

# 备份当前配置
backup_files() {
    echo -e "${YELLOW}1. 备份现有配置...${NC}"
    
    # 创建备份目录
    BACKUP_DIR="backup_$(date +%Y%m%d_%H%M%S)"
    mkdir -p $BACKUP_DIR
    
    # 备份文件
    [ -f "vite.config.js" ] && cp vite.config.js $BACKUP_DIR/
    [ -f "index.html" ] && cp index.html $BACKUP_DIR/
    [ -f "package.json" ] && cp package.json $BACKUP_DIR/
    
    echo -e "${GREEN}✓ 备份完成：$BACKUP_DIR${NC}"
}

# 应用优化配置
apply_optimizations() {
    echo -e "${YELLOW}2. 应用优化配置...${NC}"
    
    # 使用优化的vite配置
    if [ -f "vite.config.optimized.js" ]; then
        cp vite.config.js vite.config.original.js 2>/dev/null || true
        cp vite.config.optimized.js vite.config.js
        echo -e "${GREEN}✓ Vite配置已优化${NC}"
    fi
    
    # 使用优化的index.html
    if [ -f "index.optimized.html" ]; then
        cp index.html index.original.html 2>/dev/null || true
        cp index.optimized.html index.html
        echo -e "${GREEN}✓ index.html已优化${NC}"
    fi
}

# 安装压缩插件
install_compression() {
    echo -e "${YELLOW}3. 安装构建压缩插件...${NC}"
    
    # 检查是否已安装
    if ! grep -q "vite-plugin-compression" package.json; then
        npm install -D vite-plugin-compression
        echo -e "${GREEN}✓ 压缩插件已安装${NC}"
    else
        echo -e "${GREEN}✓ 压缩插件已存在${NC}"
    fi
}

# 构建项目
build_project() {
    echo -e "${YELLOW}4. 构建优化后的项目...${NC}"
    
    # 清理旧的构建
    rm -rf dist
    
    # 构建
    npm run build
    
    # 生成gzip文件
    echo -e "${YELLOW}5. 生成预压缩文件...${NC}"
    find dist -type f \( -name "*.js" -o -name "*.css" -o -name "*.html" \) -exec gzip -9 -k {} \;
    
    echo -e "${GREEN}✓ 构建完成${NC}"
}

# 显示优化结果
show_results() {
    echo -e "${YELLOW}6. 优化结果分析...${NC}"
    
    # 统计文件大小
    echo -e "\n${GREEN}构建文件大小统计：${NC}"
    echo "原始文件："
    du -sh dist/assets/*.js 2>/dev/null | sort -hr | head -5
    
    echo -e "\nGzip压缩后："
    du -sh dist/assets/*.js.gz 2>/dev/null | sort -hr | head -5
    
    # 计算压缩率
    ORIGINAL_SIZE=$(du -sb dist/assets/*.js 2>/dev/null | awk '{total += $1} END {print total}')
    GZIP_SIZE=$(du -sb dist/assets/*.js.gz 2>/dev/null | awk '{total += $1} END {print total}')
    
    if [ -n "$ORIGINAL_SIZE" ] && [ -n "$GZIP_SIZE" ] && [ "$ORIGINAL_SIZE" -gt 0 ]; then
        COMPRESSION_RATIO=$(echo "scale=2; (1 - $GZIP_SIZE / $ORIGINAL_SIZE) * 100" | bc)
        echo -e "\n${GREEN}压缩率: ${COMPRESSION_RATIO}%${NC}"
    fi
}

# 生成部署指南
generate_guide() {
    echo -e "${YELLOW}7. 生成部署指南...${NC}"
    
    cat > DEPLOYMENT_OPTIMIZATION_GUIDE.md << EOF
# 性能优化部署指南

## 已完成的优化

1. **CDN加速**
   - Vue、Element Plus、ECharts等大型库使用CDN
   - 减少服务器带宽压力

2. **代码分割**
   - 路由级别的懒加载
   - 第三方库独立打包
   - 按需加载组件

3. **资源压缩**
   - 启用Gzip压缩
   - 生成预压缩文件(.gz)

## Nginx配置更新

请将以下配置添加到您的Nginx配置中：

\`\`\`nginx
# 启用Gzip
gzip on;
gzip_vary on;
gzip_types text/plain text/css text/xml text/javascript application/json application/javascript;

# 静态资源缓存
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}

# 使用预压缩文件
location ~ \.(js|css)$ {
    gzip_static on;
}
\`\`\`

## 部署步骤

1. 上传 dist 目录到服务器
2. 更新 Nginx 配置（使用 nginx-performance.conf）
3. 重启 Nginx：\`nginx -s reload\`
4. 清理浏览器缓存后测试

## 验证优化效果

1. 使用 Chrome DevTools 的 Network 面板查看加载时间
2. 使用 Lighthouse 进行性能评分
3. 记录优化前后的加载时间对比

## 预期效果

- 首屏加载时间：从18秒降至3秒以内
- Gzip压缩率：60-80%
- 二次访问：利用缓存，1秒内加载

生成时间：$(date)
EOF

    echo -e "${GREEN}✓ 部署指南已生成：DEPLOYMENT_OPTIMIZATION_GUIDE.md${NC}"
}

# 主函数
main() {
    echo -e "${GREEN}开始执行性能优化...${NC}\n"
    
    # 检查是否在前端目录
    if [ ! -f "package.json" ]; then
        echo -e "${RED}错误：请在前端项目目录中运行此脚本${NC}"
        exit 1
    fi
    
    # 执行步骤
    backup_files
    apply_optimizations
    install_compression
    build_project
    show_results
    generate_guide
    
    echo -e "\n${GREEN}=== 优化完成！===${NC}"
    echo -e "${YELLOW}请查看 DEPLOYMENT_OPTIMIZATION_GUIDE.md 了解部署步骤${NC}"
    echo -e "${YELLOW}建议先在测试环境验证效果${NC}"
}

# 执行主函数
main