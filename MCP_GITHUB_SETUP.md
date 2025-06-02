# MCP GitHub 配置指南

## 配置步骤

1. **获取 GitHub Personal Access Token**
   - 访问 https://github.com/settings/tokens/new
   - Token 名称：建议使用 "MCP Claude Integration"
   - 权限选择：
     - ✅ repo (完整的仓库访问权限)
     - ✅ workflow (如果需要操作 GitHub Actions)
   - 点击 "Generate token" 并复制生成的 token

2. **更新 MCP 配置文件**
   - 编辑文件：`~/.config/claude/mcp_servers.json`
   - 将 `YOUR_GITHUB_TOKEN_HERE` 替换为您的实际 token

3. **重启 Claude**
   - 完全退出 Claude 应用
   - 重新启动 Claude

4. **验证配置**
   - 在 Claude 中测试 MCP GitHub 功能
   - 确认可以正常创建、更新文件

## 安全提示

- Personal Access Token 相当于密码，请妥善保管
- 不要将包含 token 的配置文件提交到版本控制系统
- 定期更新 token 以保证安全性
- 如果 token 泄露，请立即在 GitHub 上撤销并生成新的

## 故障排除

如果配置后仍无法使用：
1. 检查 token 权限是否正确
2. 确认 Claude 已完全重启
3. 查看 Claude 日志以获取更多信息