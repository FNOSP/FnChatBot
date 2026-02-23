# 任务列表 (Tasks)

- [x] Task 1: UI 修复（浅色模式、API 验证、语言）
    - [x] SubTask 1.1: 修复 `ModelServices.vue` 中的 API 配置验证（响应式问题）和预览 URL 显示。
    - [x] SubTask 1.2: 修复 `ModelServices.vue` 的浅色模式背景问题（确保 `bg-white` 或 `bg-gray-50` 正确应用）。
    - [x] SubTask 1.3: 更新 `GeneralSettings.vue`，移除保存按钮并实现自动保存；将语言选项更新为母语名称。
    - [x] SubTask 1.4: 更新 `Sidebar.vue`、`SettingsView.vue` 和 `locales/ja.json`，修复缺失的翻译。

- [x] Task 2: 后端聊天功能实现（真实 API）
    - [x] SubTask 2.1: 更新 `internal/api/ws/websocket.go`，根据对话的 `model_id` 获取 `ModelConfig`。
    - [x] SubTask 2.2: 实现 `streamChatCompletion` 辅助函数，调用外部 API（OpenAI 兼容）并流式传输分块。
    - [x] SubTask 2.3: 在 `handleUserMessage` 中用真实 API 调用替换模拟逻辑。

- [x] Task 3: 验证（端到端）
    - [x] SubTask 3.1: 启动后端和前端。
    - [x] SubTask 3.2: 使用 `dev-browser` 配置 API（使用用户提供的 ChatAnywhere 凭据）。
    - [x] SubTask 3.3: 使用 `dev-browser` 验证模型列表获取功能。
    - [x] SubTask 3.4: 使用 `dev-browser` 开始聊天，并验证响应不是模拟文本。
    - [x] SubTask 3.5: 通过截图或 DOM 检查验证浅色模式 UI 和语言切换功能。

# 任务依赖
- Task 2 依赖 Task 1（或可并行）。
- Task 3 依赖 Task 1 和 2。
