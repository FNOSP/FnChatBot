# Tasks

- [x] Task 1: 定义新的数据模型 (Models)
  - [x] SubTask 1.1: 在 `internal/models` 中定义 `Session`, `Message`, `Part` 结构体。
  - [x] SubTask 1.2: 配置 GORM 的关联关系（HasMany, BelongsTo）。
  - [x] SubTask 1.3: 更新 `internal/db/db.go` 中的自动迁移逻辑。

- [x] Task 2: 重构记忆服务 (Memory Service)
  - [x] SubTask 2.1: 更新 `SQLiteHistory.AddMessage` 方法，实现消息拆分为 Parts 的逻辑。
  - [x] SubTask 2.2: 更新 `SQLiteHistory.Messages` 方法，实现从 Parts 重组消息的逻辑。
  - [x] SubTask 2.3: 适配 `Clear` 和其他辅助方法。

- [x] Task 3: 适配上层调用 (Service Layer)
  - [x] SubTask 3.1: 检查 `internal/services/llm/service.go` 是否需要调整（主要是确保它通过 `SQLiteHistory` 交互）。
  - [x] SubTask 3.2: 检查 API 层是否直接依赖了旧的 `Message` 结构，如有则进行适配。

- [x] Task 4: 验证与清理
  - [x] SubTask 4.1: 清理旧表结构（Drop Tables），重新初始化数据库。
  - [x] SubTask 4.2: 运行测试确保聊天流程正常（包括工具调用场景）。

# Task Dependencies
- Task 2 depends on Task 1.
- Task 3 depends on Task 2.
- Task 4 depends on Task 3.
