# InnerU_SSE

一个带有SSE（服务器发送事件）通知和简约前端演示的小型TODO服务器。

## 概述

本项目提供：
- 一个使用SQLite作为后端的简单TODO API（CRUD）Go HTTP服务器。
- SSE端点，用于向连接的客户端广播变更（todo.created, todo.updated, todo.completed, todo.deleted）。
- 位于 `web/index.html` 的简约前端演示，用于连接和查看实时通知。

## 环境要求

- Go 1.24+（开发和CI）
- （可选）curl或PowerShell用于测试

## 运行

在前台运行以查看日志：

./bin/todo-sse.exe

服务器默认监听8080端口。在浏览器中打开 http://localhost:8080。

## 前端演示

在浏览器中打开：

http://localhost:8080/index.html

- 输入令牌（演示使用 `demo-token`）并点击连接。
- 页面将显示连接状态和SSE消息出现的日志区域。

## API示例（PowerShell）

在PowerShell中使用 `Invoke-RestMethod`（推荐在Windows上使用）。根据需要替换ID。

# 创建
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/todos" -ContentType "application/json" -Body '{"title":"买牛奶","description":"2L"}'

# 列表
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/todos"

# 更新（PUT）
Invoke-RestMethod -Method Put -Uri "http://localhost:8080/api/todos/1" -ContentType "application/json" -Body '{"title":"买牛奶（2L）","description":"低脂"}'

# 标记完成（PATCH到/complete）
Invoke-RestMethod -Method Patch -Uri "http://localhost:8080/api/todos/1/complete" -ContentType "application/json" -Body '{"completed":true}'

# 删除
Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/api/todos/1"

## 从终端测试SSE（curl）

# 监听SSE流（在Windows上使用curl.exe）
curl.exe -N "http://localhost:8080/sse?token=demo-token"

# 然后在另一个终端中，创建一个todo以查看SSE通知
curl.exe -X POST -H "Content-Type: application/json" -d '{"title":"测试SSE"}' "http://localhost:8080/api/todos"

## 注意事项和最佳实践

- 不要提交构建的二进制文件或本地数据库文件。将它们添加到 `.gitignore`（参见仓库中的 `.gitignore`）。如果不小心提交了，请使用 `git rm --cached` 从git索引中删除。
- 生产环境使用时，将演示令牌机制替换为适当的身份验证，并考虑使用更健壮的数据库（Postgres）或共享的发布/订阅代理来实现多实例SSE。
- 服务器写入本地SQLite文件；SQLite只允许一个并发写入者——应用程序配置应设置 `db.SetMaxOpenConns(1)`。

## 开发

- 运行 `gofmt -w .`、`go vet ./...` 和（推荐）`staticcheck ./...`。
- 运行单元测试：`go test ./... -v`。

## 演示脚本

您可以创建一个小的PowerShell演示脚本来启动服务器并创建示例todo；参见 `demo.ps1`（未包含）。

## 许可证

MIT