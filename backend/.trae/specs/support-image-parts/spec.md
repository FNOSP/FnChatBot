# 支持图片存储与枚举类型优化 Spec

## 为什么 (Why)
目前的 `Part` 类型使用硬编码的字符串，容易出错且难以管理。此外，系统目前不支持图片输入，限制了用户与多模态模型交互的能力。

## 变更内容 (What Changes)
-   **枚举管理**：在 `internal/models` 中引入 `PartType` 字符串枚举类型，统一管理 Part 类型（文本、文件、工具调用、工具结果）。
-   **图片支持**：
    -   更新 `Part` 模型：`Type` 字段使用新的 `PartType` 枚举。
    -   新增 `ImagePartMeta` 结构体，用于标准化图片元数据（mime类型, 文件名等）。
    -   更新 `SQLiteHistory` 中的 `AddMessage` 逻辑：
        -   处理包含图片的输入消息。
        -   将图片 Base64 内容存储在 `Content` 字段。
        -   将元数据（MIME 类型等）存储在 `Meta` 字段。
-   **WebSocket**：更新 WebSocket 消息处理逻辑，支持解析前端传来的图片数据。

## 影响范围 (Impact)
-   **受影响模块**：聊天功能、数据存储。
-   **受影响代码**：
    -   `internal/models/models.go`: 添加枚举定义和更新 `Part` 结构体使用。
    -   `internal/services/memory/sqlite_history.go`: 消息保存逻辑。
    -   `internal/api/ws/websocket.go`: 用户输入解析逻辑。

## 新增需求 (ADDED Requirements)
### Requirement: PartType 枚举
系统应定义 `PartType` 字符串枚举，包含以下值：
-   `PartTypeText` ("text")
-   `PartTypeFile` ("file") - 用于图片和其他文件
-   `PartTypeToolCall` ("tool_calls")
-   `PartTypeToolResult` ("tool_result")

### Requirement: 图片存储
系统应将图片输入作为 `Part` 记录存储：
-   `Type`: `PartTypeFile`
-   `Content`: 图片的 Base64 字符串。
-   `Meta`: JSON 对象，包含 `mime` (例如 "image/png"), `filename` (可选)。

#### Scenario: 用户发送图片
-   **WHEN** 用户通过 WebSocket 发送包含图片的消息
-   **THEN** 系统解析图片数据。
-   **AND** 保存一个 `Part`，其中 `Type="file"`, `Content=<base64>`, `Meta={"mime": "..."}`。

## 修改需求 (MODIFIED Requirements)
### Requirement: 消息处理
`SQLiteHistory` 中的 `AddMessage` 函数应更新为遍历消息的各个部分，并使用新的枚举常量创建对应的 `Part` 记录。

## 移除需求 (REMOVED Requirements)
无
