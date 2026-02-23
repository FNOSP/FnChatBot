# 后端测试与能力增强规范 (Backend Test & Capability Spec)

## 为什么 (Why)
为了验证 AI 在对话中调用 Skill 和 MCP (Model Context Protocol) 的能力，并确保后端核心对话流程的稳定性。目前后端尚未完全支持在流式对话中动态注入和执行工具，因此需要增强后端逻辑以支持 OpenAI 格式的 Tool Calling，并编写配套的集成测试。

## 变更内容 (What Changes)
- **新增文件**:
    - `backend/tests/integration/chat_flow_test.go`: 包含基础对话、Skill 调用、MCP 调用的集成测试。
    - `backend/internal/services/tool_service.go`: 用于聚合和格式化 Skill/MCP 为 OpenAI Tool 定义。
- **修改文件**:
    - `backend/internal/api/ws/websocket.go`:
        - 更新 `streamChatCompletion` 以支持 `tools` 参数。
        - 增加处理 `tool_calls` 响应的逻辑（多轮对话）。
        - 实现工具执行结果的回传。
- **配置**: 使用环境变量 (`TEST_BASE_URL`, `TEST_API_KEY`, `TEST_MODEL_ID`) 配置测试环境。
- **参考**: 遇到实现问题时，参考 `d:\github\learn-claude-code` 中的 Agent 和 Tool 实现模式。

## 新增需求 (ADDED Requirements)
### Requirement: 工具调用支持 (Tool Calling Support)
系统应当支持在对话中根据用户意图自动调用注册的 Skill 或 MCP 工具。
1.  **工具注入**: 在调用 OpenAI API 时，系统需将所有启用的 Skill 和 MCP 转换为 `tools` 数组参数。
2.  **多轮执行**: 当模型返回 `tool_calls` 时，系统需：
    - 解析工具调用请求。
    - 执行对应的本地 Skill 逻辑或远程 MCP 调用。
    - 将执行结果作为 `tool` 类型的消息追加到对话历史。
    - 再次调用模型以获取最终回复。

### Requirement: 集成测试覆盖 (Integration Test Coverage)
提供以下测试用例：
1.  **基础流式对话**: 验证普通文本回复。
2.  **Skill 调用测试**:
    - 注册一个测试 Skill (如 "get_current_time")。
    - 用户询问相关问题。
    - 验证系统是否调用了该 Skill 并返回了包含时间的结果。
3.  **MCP 调用测试**:
    - 启动一个 Mock MCP Server。
    - 注册该 MCP 配置。
    - 用户询问相关问题。
    - 验证系统是否向 Mock Server 发起了请求并正确处理响应。

## 影响 (Impact)
- **受影响的代码**: `websocket.go` 的对话循环逻辑将变得更加复杂（需处理多轮 Tool Call）。
- **数据结构**: 可能需要调整 `ChatMessage` 以支持 `tool_calls` 和 `tool_call_id` 字段。

## 迁移 (Migration)
- 无需数据迁移，仅逻辑增强。
