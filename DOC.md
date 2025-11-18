# DOC — TODO SSE 通知模块（说明文档）

1. 项目概述
   - 本模块实现基于 SSE 的实时通知，适用于 TODO 应用的提醒推送。
   - 已实现：基础功能的增删改查、SSE 订阅、广播接口、心跳、简单 token 校验示例、前端 demo。
   - TDD： 1. 重点文件hub.go store.go 都涵盖了单测，做了解耦  2.测试驱动开发，保证代码质量。
   - 单测：
   store.go 主要测试重点
   SQLite 连接测试
   验证 NewStore 函数能否正确打开 SQLite 数据库连接
   确保 init 方法正确创建 todos 表结构
   测试数据库文件创建和基本连接功能

   CRUD 功能测试
   Create: 测试 CreateTodo 方法是否能正确插入包含所有字段的记录
   Read: 验证 GetTodos 和 GetTodoByID 方法是否能正确检索数据
   Update: 确认 UpdateTodo 方法是否能准确修改现有记录
   Delete: 确保 DeleteTodo 方法是否能正确删除记录

   hub.go 主要测试重点
   核心功能测试：
   TestHub 中的客户端注册和注销
   向多个客户端的消息广播
   消息传递的正确性验证

   HTTP处理器测试：
   状态码验证（缺少token时返回401，有效请求返回200）
   SSE的Content-Type头部检查
   查询参数验证

1. 技术选型
   - 语言/框架：Golang 原生 net/http（轻量、易打包）。
   - 持久化：当前实现不需要持久化；如需断线补发，建议结合短期消息存储（Redis 列表或 DB）。
   - SSE vs WebSocket
 - 选用SSE做实时推送提醒，webSocket 适合双向通信，针对于实时通知，无需双向的场景，我选择SSE实现。 
   - 数据库选型：sqlite 
 - 原因： 快速设置：无需运行独立的数据库服务
         便携性：数据库随项目文件一起移动
         隔离性：每个开发者都可以拥有自己的数据库实例
         易于重置：只需删除SQLite文件即可重新开始

2. 架构
   - Hub（内存）管理本实例所有连接。
   - 多实例：各实例通过 Redis pub/sub 接收广播并推送到本地客户端。

3. 运行说明
   - 本地：go build -o bin/todo-sse.exe && .\bin\todo-sse.exe
   - Docker：docker build -t todo-sse . && docker run -p 8080:8080 todo-sse

4. API
   - /sse?token=... 订阅（EventSource）
   - /notify (POST JSON 或 GET ?msg=) 广播

5. 关键实现点
   - SSE headers: Content-Type: text/event-stream, keep-alive, no-cache
   - 心跳每 25s 发送注释行防止代理断开
   - 客户端 reconnect 使用浏览器自带的重连逻辑，可利用 id + Last-Event-ID 实现断线补发（需消息持久化）

6. 安全与鉴权建议
   - 不要在生产日志里暴露 token（如果用 query param，注意日志清理）。
   - 推荐使用 HttpOnly cookie 或预签名短期 URL 来携带凭证。
   - HTTPS 强制。

7. 已知限制与未来改进
   - 当前实现为单实例 in-memory hub，无法跨实例广播（可加入 Redis）。
   - 未实现用户分组、权限控制、消息持久化与补发逻辑。
   - 未实现细粒度限流与连接数限制。

8. AI 使用说明
   - （选用copilot作为辅助开发，主要用于简单功能快速实现，前端简易页面快速提供示例demo； 为什么没有用cursor/claude， 因为到期了，好贵好贵，学生党用不起 哈哈哈）

9.  部署注意
   - 代理（nginx 等）需关闭响应缓冲（proxy_buffering off）并支持长连接.


10. 交付说明
   code review / PR 检查清单（checklist）

   编译与测试：
   go test ./... 全部通过
   go vet 无 error
   staticcheck 报告没有高优问题

  代码质量：
  没有未使用的导包或变量
  函数短小、注释充足（公共 API 需 doc comment）
  错误处理明确，避免 silent ignore

  安全与配置：
  不在代码或 repo 中保存敏感信息（token/secret）
  .gitignore 排除本地文件（db、bin、IDE）
  SSE token 示例仅用于 demo，说明生产中要换更安全方法

  可运行性与部署：
  README 有清晰的运行步骤（build/run/test/demo）
  Dockerfile/CI 配置能在干净环境里构建

  API 文档：
  README 或 docs 给出 API 列表（示例请求/响应）
  若时间允许，附上 curl 命令示例
