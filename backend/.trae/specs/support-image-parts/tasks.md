# Tasks

- [ ] Task 1: 定义 PartType 枚举和 Meta 结构体
  - [ ] SubTask 1.1: 在 `internal/models/models.go` 中定义 `PartType` 字符串类型及常量 (`PartTypeText`, `PartTypeFile` 等)。
  - [ ] SubTask 1.2: 定义 `FilePartMeta` 结构体，用于文件/图片元数据的 JSON 序列化。

- [ ] Task 2: 更新存储逻辑
  - [ ] SubTask 2.1: 重构 `internal/services/memory/sqlite_history.go`，使用 `PartType` 常量替代硬编码字符串。
  - [ ] SubTask 2.2: 更新 `AddMessage` 方法以处理 `llms.ImagePart` (或等效内容部分) -> 保存为 `PartTypeFile`，内容为 Base64，元数据包含 MIME 类型。

- [ ] Task 3: 更新 WebSocket 处理
  - [ ] SubTask 3.1: 在 `internal/api/ws/websocket.go` 中更新输入解析逻辑，检测图片数据（假设前端发送特定的结构）。
  - [ ] SubTask 3.2: 构建包含 `llms.ImageURLPart` 或 `llms.BinaryPart` 的 `llms.MessageContent` 传递给 LLM 服务。

- [ ] Task 4: 验证
  - [ ] SubTask 4.1: 创建测试用例（在 `sqlite_history_test.go` 或类似文件中），模拟保存和加载图片消息。
  - [ ] SubTask 4.2: 验证 `Content` 字段存储了 Base64，且 `Meta` 字段包含正确的 MIME 类型。

# Task Dependencies
- Task 2 depends on Task 1.
- Task 3 depends on Task 2.
- Task 4 depends on Task 3.
