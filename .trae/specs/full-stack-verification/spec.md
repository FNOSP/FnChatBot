# 全链路集成测试与能力增强规范 (Full-Stack Integration Test & Capability Enhancement Spec)

## 为什么 (Why)
为了验证 AI 系统的全链路能力，特别是基于 `mini-claude-code` 架构设计中的核心机制：**任务规划 (TodoWrite)**、**子代理协作 (Subagents)**、**技能加载 (Skills)** 以及 **安全沙箱 (Sandbox)**。我们需要确保这些机制不仅在后端逻辑中正确实现，还能在前端 UI 中得到直观展示，并通过端到端测试进行验证。同时，需要确保项目中不包含特定服务商（ChatAnywhere）的敏感信息。

## 变更内容 (What Changes)
- **新增规范**: `frontend/tests/e2e/full_stack_spec.md` (本文档)
- **前端增强**:
    - **任务面板 (Task Panel)**: 基于 v2 TodoWrite 理念，新增持久化任务列表展示，显示 `pending`, `in_progress` (含 `activeForm`), `completed` 状态。
    - **子代理状态 (Subagent Status)**: 基于 v3 Subagents 理念，展示当前活跃的子代理类型（如 `explore`, `code`, `plan`）及其进度。
    - **技能加载通知 (Skill Notification)**: 基于 v4 Skills 理念，当 AI 加载技能时显示通知。
    - **沙箱标识 (Sandbox Tag)**: Bash 命令执行时显示 "Sandbox Mode" 标签。
- **后端增强**:
    - **工具服务 (ToolService)**:
        - 实现 `TodoWrite` 工具：用于管理任务状态。
        - 实现 `Task` 工具：用于委派子任务给子代理。
        - 优化 `Skill` 工具：确保内容作为 `tool_result` 追加到对话历史（利用缓存）。
    - **WebSocket**: 推送任务更新、子代理状态和沙箱状态。
- **配置脱敏**: 扫描全项目，统一使用 `OPENAI_API_KEY` 或通用配置。

## 新增需求 (ADDED Requirements)
### Requirement: 任务规划可视化 (v2 TodoWrite)
系统应当支持 AI通过 `TodoWrite` 工具管理任务列表。
- **UI**: 在侧边栏或独立面板展示当前任务列表。
- **逻辑**: 每次更新为全量覆盖，高亮显示 `in_progress` 任务及其 `activeForm`（进行时描述）。

### Requirement: 子代理协作机制 (v3 Subagents)
系统应当支持 AI 通过 `Task` 工具委派任务。
- **Backend**: 识别 `Task` 调用，模拟子代理执行（隔离上下文）。
- **UI**: 当子代理执行时，主对话框显示 "Subagent [Type] working..."，并流式展示其进度或最终结果。

### Requirement: 技能动态加载 (v4 Skills)
系统应当支持 AI 通过 `Skill` 工具加载领域知识。
- **Backend**: 加载 `SKILL.md` 内容并作为 `tool_result` 返回，不修改 System Prompt。
- **UI**: 提示 "Loaded Skill: [Name]"。

### Requirement: 全链路 E2E 测试
使用 Playwright 模拟真实用户交互：
1.  **复杂任务**: 用户请求 "分析项目并重构"。
2.  **验证点**:
    - AI 调用 `TodoWrite` -> UI 显示任务列表。
    - AI 调用 `Task(explore)` -> UI 显示子代理状态。
    - AI 调用 `Skill(code-review)` -> UI 显示技能加载。
    - AI 执行 Bash -> UI 显示 Sandbox 标签。

## 影响 (Impact)
- **受影响的代码**: `frontend/src/components/*`, `backend/internal/services/tool_service.go`, `backend/internal/api/ws/websocket.go`。
- **数据结构**: WebSocket 消息需增加 `tasks`, `subagent_status`, `loaded_skills` 字段。

## 迁移 (Migration)
- 确保所有 API Key 配置通过环境变量注入，无硬编码。
