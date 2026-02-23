# 任务列表 (Tasks)

- [x] Task 1: 研究参考实现 (Study Reference)
    - [x] SubTask 1.1: 分析 `d:\github\learn-claude-code` (特别是 `agents/v4_skills_agent.py` 等) 中的 Tool 调用与 Skill 管理逻辑。
    - [x] SubTask 1.2: 确定如何在 Go 后端中适配 Python 项目中的 Agent 模式。

- [x] Task 2: 实现工具服务 (Tool Service)
    - [x] SubTask 2.1: 定义 `Tool` 和 `Function` 结构体（匹配 OpenAI API 格式）。
    - [x] SubTask 2.2: 在 `services/tool_service.go` 中实现 `GetAvailableTools` 方法，从数据库加载启用的 Skill 和 MCP，并转换为 Tool 定义。
    - [x] SubTask 2.3: 实现简单的 Skill 执行逻辑（支持通过 `Config` 定义的简单 HTTP 请求或内置函数）。
    - [x] SubTask 2.4: 实现基础的 MCP Client 逻辑（通过 HTTP POST/SSE 调用远程 MCP 工具）。

- [x] Task 3: 增强 WebSocket 对话逻辑
    - [x] SubTask 3.1: 修改 `ChatCompletionRequest` 结构体，增加 `Tools` 和 `ToolChoice` 字段。
    - [x] SubTask 3.2: 修改 `streamChatCompletion`，在请求中注入 `GetAvailableTools()` 的结果。
    - [x] SubTask 3.3: 重构 `streamChatCompletion` 或上层逻辑，以支持 "Request -> Tool Call -> Execute -> Submit Result -> Final Response" 的循环。
    - [x] SubTask 3.4: 处理 `tool_calls` 的流式响应解析（OpenAI 流式返回 Tool Call 时通常是分段的，需要拼接）。

- [x] Task 4: 编写集成测试 (Integration Tests)
    - [x] SubTask 4.1: 搭建测试基础架构（Gin, SQLite, Env Vars）。
    - [x] SubTask 4.2: 编写 `TestChatFlow_Basic`（基础对话）。
    - [x] SubTask 4.3: 编写 `TestChatFlow_WithSkill`。
        - 创建一个返回固定值的测试 Skill。
        - 验证模型是否发起调用。
    - [x] SubTask 4.4: 编写 `TestChatFlow_WithMCP`。
        - 使用 `httptest.NewServer` 模拟 MCP Server。
        - 验证后端是否正确请求 Mock Server。

- [x] Task 5: 验证与修复
    - [x] SubTask 5.1: 运行测试并修复 `streamChatCompletion` 中的 bug。
    - [x] SubTask 5.2: 确保真实 API Key 调用下，普通对话和工具调用均能正常工作。
