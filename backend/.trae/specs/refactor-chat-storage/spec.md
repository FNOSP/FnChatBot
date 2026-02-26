# 重构聊天存储结构 Spec

## 为什么 (Why)
目前的聊天存储结构（`Conversation` -> `Message`）较为简单，无法精细地存储 AI 的思维过程（Reasoning）、工具调用详情（Tool Calls）以及执行结果。参考 opencode 的设计，引入更颗粒化的 `Part` 结构，可以支持：
1.  **思维链回溯**：独立存储 AI 的推理过程。
2.  **多模态支持**：未来可扩展支持图片、文件等多种类型的消息部分。
3.  **精细化状态管理**：准确记录工具调用的参数、状态（Pending/Running/Completed）和结果。

## 变更内容 (What Changes)
-   **重命名/重构模型**：
    -   `Conversation` -> `Session`（保持概念一致，表名迁移）。
    -   `Message`：仅作为容器，不再直接存储 `Content`。
    -   新增 `Part` 模型：存储具体的内容片段。
-   **数据库变更**：
    -   新建 `sessions` 表（替代 `conversations`）。
    -   修改 `messages` 表结构。
    -   新建 `parts` 表。
-   **业务逻辑变更**：
    -   `SQLiteHistory` 服务：适配新的数据结构，负责将 `llms.ChatMessage` 拆解为 `Part` 进行存储，以及从 `Part` 重组消息。

## 影响范围 (Impact)
-   **受影响模块**：
    -   `internal/models`: 定义新的结构体。
    -   `internal/services/memory`: `SQLiteHistory` 的增删改查逻辑。
    -   `internal/db`: 数据库迁移逻辑。
    -   `internal/api`: 获取历史记录的接口可能需要适配（如果直接返回模型）。

## 新增需求 (ADDED Requirements)
### Requirement: Part 模型
系统应支持以下类型的 Part：
-   `text`: 普通文本内容。
-   `reasoning`: AI 的推理/思考内容。
-   `tool_call`: 工具调用请求（包含 `call_id`, `name`, `args`）。
-   `tool_result`: 工具执行结果（包含 `call_id`, `result`, `is_error`）。

#### Scenario: 存储工具调用
-   **WHEN** AI 返回包含工具调用的消息
-   **THEN** 系统应创建一条 `Message`，并创建多个 `Part`：
    -   一个 `text` Part（如果有）。
    -   一个或多个 `tool_call` Part。

## 修改需求 (MODIFIED Requirements)
### Requirement: 消息存储 (SQLiteHistory)
-   `AddMessage`: 需将 `llms.ChatMessage` 解析为 `Message` + `[]Part`。
-   `Messages`: 需查询 `Message` 及其关联的 `Part`，并按顺序重组为 `llms.ChatMessage`。

## 移除需求 (REMOVED Requirements)
### Requirement: Message.Content & Message.Meta
-   **Reason**: 内容和元数据现已分散存储在 `Part` 中。
-   **Migration**: 开发阶段不保留旧数据，直接重置数据库表结构。
