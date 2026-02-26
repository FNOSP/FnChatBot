# 使用 LangChainGo 重构后端并集成 SQLite 记忆存储

## 为什么 (Why)
目前的后端使用针对每个 AI 供应商的自定义 HTTP 适配器，这种方式难以维护和扩展。切换到 `LangChainGo` 将为多个 AI 供应商提供统一、标准的接口。此外，利用现有的 SQLite 数据库存储聊天记录（而不是仅仅依赖内存或外部系统），可以确保数据的持久化和一致性。

## 变更内容 (What Changes)
- **替换自定义适配器**: 移除或重构 `internal/services/adapters/*`，改为使用 `github.com/tmc/langchaingo`。
- **新 AI 服务**: 创建一个新的服务层，封装 `LangChainGo` 以处理聊天 (Chat) 和流式 (Streaming) 请求。
- **SQLite 记忆存储**: 使用现有的 GORM/SQLite 设置 (`internal/db` 和 `internal/models`) 实现 `langchaingo/schema.ChatMessageHistory` 接口。
- **更新 API 处理程序**: 修改 `internal/api/*` 处理程序，使用新的 `LangChainGo` 服务替代旧的 `ProviderAdapter`。
- **依赖项**: 添加 `github.com/tmc/langchaingo`。

## 影响 (Impact)
- **受影响的规范**: 无直接影响，但“聊天”和“历史记录”的实现细节会发生变化。
- **受影响的代码**:
    - `internal/services/adapters/`: 将被替换或大幅修改。
    - `internal/services/provider_adapter.go`: 接口将发生变化。
    - `internal/api/handlers.go` & `model_handlers.go`: 将更新为调用新服务。
    - `go.mod`: 新增依赖项。

## 新增需求 (ADDED Requirements)
### 需求: LangChainGo 集成
系统必须使用 `LangChainGo` 与 AI 供应商（OpenAI, Anthropic 等）进行通信。
- **场景: 聊天完成**
    - **当 (WHEN)** 用户发送消息时
    - **那么 (THEN)** 系统使用 `LangChainGo` 从选定的供应商生成回复。

### 需求: SQLite 聊天记录
系统必须从 SQLite 的 `messages` 表中存储和检索聊天记录。
- **场景: 上下文聊天**
    - **当 (WHEN)** 用户继续对话时
    - **那么 (THEN)** 系统从 SQLite 加载之前的消息并通过 `LangChainGo` 传递给 LLM。

## 修改的需求 (MODIFIED Requirements)
### 需求: 供应商适配器接口
`ProviderAdapter` 接口将不再暴露低级别的 HTTP 方法（`BuildChatRequest`, `ParseStreamResponse`）。相反，它将暴露使用 `LangChainGo` 的高级 `Generate` 或 `Stream` 方法。

## 移除的需求 (REMOVED Requirements)
### 需求: 自定义 HTTP 适配器
**原因**: 被 `LangChainGo` 的内置供应商替代。
