# TODO SSE 通知服务（Golang）

简介
- 这是一个基于 Go 的 SSE（Server-Sent Events）通知示例服务，适合作为 TODO 应用的实时提醒/通知模块。
- 支持简单 token 认证（通过查询参数 token），支持心跳、Last-Event-ID 协议占位，支持 /notify 广播接口用于测试或后端触发。

快速运行
1. Go 环境 (1.20+)：
   go build -o bin/todo-sse.exe
   .\bin\todo-sse.exe

2. Docker：
   docker build -t todo-sse:latest .
   docker run -p 8080:8080 todo-sse:latest

接口
- GET /sse?token=your-token  -> SSE 订阅（浏览器端用 EventSource）
- POST /notify (JSON body) -> 广播通知到所有已连接客户端
- GET /notify?msg=xxx -> 简单广播用于快速测试

示例前端
- web/index.html 提供了一个演示页面，使用 EventSource 连接并显示消息。

开发说明
- hub.go 实现了连接管理与广播逻辑，单实例下工作良好。
- 多实例扩展：建议接入 Redis pub/sub 或 Kafka 作为消息总线。

文档
- DOC.md 包含设计决策、运行与部署说明、AI 使用记录等（请补充）。
