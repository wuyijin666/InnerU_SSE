# DOC — TODO SSE 通知模块（说明文档）

1. 项目概述
   - 本模块实现基于 SSE 的实时通知，适用于 TODO 应用的提醒推送。
   - 已实现：SSE 订阅、广播接口、心跳、简单 token 校验示例、前端 demo。

2. 技术选型
   - 语言/框架：Golang 原生 net/http（轻量、易打包）。
   - 持久化：当前实现不需要持久化；如需断线补发，建议结合短期消息存储（Redis 列表或 DB）。
   - 扩展总线：Redis pub/sub 或 Kafka（推荐 Redis 用于快速集成）。

3. 架构
   - Hub（内存）管理本实例所有连接。
   - 多实例：各实例通过 Redis pub/sub 接收广播并推送到本地客户端。

4. 运行说明
   - 本地：go build -o bin/todo-sse.exe && .\bin\todo-sse.exe
   - Docker：docker build -t todo-sse . && docker run -p 8080:8080 todo-sse

5. API
   - /sse?token=... 订阅（EventSource）
   - /notify (POST JSON 或 GET ?msg=) 广播

6. 关键实现点
   - SSE headers: Content-Type: text/event-stream, keep-alive, no-cache
   - 心跳每 25s 发送注释行防止代理断开
   - 客户端 reconnect 使用浏览器自带的重连逻辑，可利用 id + Last-Event-ID 实现断线补发（需消息持久化）

7. 安全与鉴权建议
   - 不要在生产日志里暴露 token（如果用 query param，注意日志清理）。
   - 推荐使用 HttpOnly cookie 或预签名短期 URL 来携带凭证。
   - HTTPS 强制。

8. 已知限制与未来改进
   - 当前实现为单实例 in-memory hub，无法跨实例广播（可加入 Redis）。
   - 未实现用户分组、权限控制、消息持久化与补发逻辑。
   - 未实现细粒度限流与连接数限制。

9. AI 使用说明
   - （如使用 AI 辅助开发，请在此列出工具与使用位置）

10. 部署注意
   - 代理（nginx 等）需关闭响应缓冲（proxy_buffering off）并支持长连接.
