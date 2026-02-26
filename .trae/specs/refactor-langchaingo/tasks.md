# 任务列表

- [x] 任务 1: 添加依赖项
    - [ ] 运行 `go get github.com/tmc/langchaingo` 以及相关的供应商包。
- [x] 任务 2: 实现 SQLite 聊天记录存储
    - [ ] 创建 `internal/services/memory/sqlite_history.go`。
    - [ ] 使用 `gorm` 和 `models.Message` 实现 `schema.ChatMessageHistory` 接口。
- [x] 任务 3: 创建 LangChain 服务
    - [ ] 创建 `internal/services/llm/service.go`。
    - [ ] 实现工厂逻辑，根据配置（OpenAI, Anthropic 等）创建 `llms.Model`。
    - [ ] 使用 `LangChainGo` 实现 `Chat` 和 `Stream` 方法。
- [x] 任务 4: 重构 API 处理程序
    - [ ] 更新 `internal/api/model_handlers.go`（或相关处理程序）以使用新的 `LangChainService`。
    - [ ] 在聊天端点中移除旧的 `ProviderAdapter` 使用。
    - [ ] 确保 `FetchModels` 仍然工作（保留旧逻辑或适配）。
- [x] 任务 5: 清理与验证
    - [ ] 如果完全替换，移除 `internal/services/adapters/` 中过时的适配器代码。
    - [ ] 验证功能是否正常。

# 任务依赖
- 任务 3 依赖于 任务 1。
- 任务 4 依赖于 任务 2 和 任务 3。
