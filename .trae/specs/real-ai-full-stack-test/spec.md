# 真实环境全链路集成测试与机制增强 (Real-World Full-Stack Integration Test & Mechanism Enhancement)

## 为什么 (Why)
用户指出之前的测试未能产生真实的 AI API 调用，且需要验证 `mini-claude-code` 架构（v2 Todo, v3 Subagents, v4 Skills）在真实交互中的表现。我们需要在**不使用 Mock** 的情况下，配置真实的后端环境（连接 ChatAnywhere API），并运行前端，通过浏览器自动化脚本模拟用户操作，从而触发真实的链式调用。

## 变更内容 (What Changes)
- **配置注入**: 确保后端服务启动时接收用户提供的真实 `BASE_URL` 和 `API_KEY`。
- **工具逻辑优化 (`backend/internal/services/tool_service.go`)**:
    - **Skill**: 确保 `Skill` 工具不仅仅返回静态文本，而是模拟从文件系统加载 `SKILL.md` (或数据库) 的行为，并以 `tool_result` 格式返回，让 LLM 感知到新知识注入。
    - **Task**: 优化 `Task` 工具的返回结果，使其更像是一个子代理的执行报告，诱导主 LLM 继续对话。
- **前端测试脚本 (`frontend/tests/e2e/real_ai_test.ts`)**:
    - 编写一个新的 Playwright 脚本，用于执行真实的对话流程。
    - 流程包括：设置 API -> 开启对话 -> 触发任务规划 (Todo) -> 触发子任务 (Task) -> 触发技能加载 (Skill) -> 触发代码执行 (Sandbox)。

## 影响 (Impact)
- **后端**: 运行时配置依赖环境变量。
- **测试**: 需要网络连接访问外部 API。

## 新增需求 (ADDED Requirements)
### Requirement: 真实 API 交互
- 系统**必须**使用 `gpt-4o-mini-ca` 模型向 `https://api.chatanywhere.tech` 发起请求。
- 测试过程中产生的 Token 消耗应能在服务商后台观测到。

### Requirement: 机制验证 (v2/v3/v4)
1.  **v2 TodoWrite**: 用户请求复杂任务时，LLM **主动**调用 `TodoWrite`，前端 TaskPanel 实时更新。
2.  **v3 Subagents**: 用户请求委派任务时，LLM 调用 `Task`，前端显示子代理状态。
3.  **v4 Skills**: 用户请求特定领域知识时，LLM 调用 `Skill`，前端显示技能加载。
4.  **Sandbox**: LLM 生成 Bash 代码块时，前端显示 "Sandbox Mode" 标签。

## 迁移 (Migration)
- 无需数据库迁移，主要是配置和运行时环境的调整。
