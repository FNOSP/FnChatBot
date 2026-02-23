# 任务列表 (Tasks)

- [x] Task 1: 项目脱敏与配置检查 (Desensitization)
    - [x] SubTask 1.2: 确保测试代码中仅使用环境变量注入 Key。

- [x] Task 2: 后端核心机制实现 (Backend Core Implementation)
    - [x] SubTask 2.1: 实现 `TodoWrite` 工具逻辑 (v2)，维护会话级任务状态。
    - [x] SubTask 2.2: 实现 `Task` 工具逻辑 (v3)，模拟子代理执行流程（隔离上下文）。
    - [x] SubTask 2.3: 优化 `Skill` 工具逻辑 (v4)，确保作为 `tool_result` 返回。
    - [x] SubTask 2.4: 更新 `WebSocket` 消息结构，包含 `tasks`, `subagent`, `sandbox` 字段。

- [x] Task 3: 前端 UI 增强 (Frontend UI Enhancement)
    - [x] SubTask 3.1: 创建 `TaskPanel.vue`，展示 `TodoWrite` 的任务列表（支持 `activeForm`）。
    - [x] SubTask 3.2: 更新 `MessageItem.vue`，展示子代理状态和技能加载通知。
    - [x] SubTask 3.3: 在 Bash 命令执行块中添加 "Sandbox Mode" 标签。

- [x] Task 4: E2E 测试脚本编写 (E2E Test Scripting)
    - [x] SubTask 4.1: 配置 Playwright 环境。
    - [x] SubTask 4.2: 编写 `test_complex_workflow.spec.ts`：
        - 模拟用户输入复杂重构任务。
        - 验证 Task Panel 更新。
        - 验证 Subagent 状态显示。
        - 验证 Skill 加载提示。
        - 验证 Sandbox 标签。

- [x] Task 5: 执行与验证 (Execution & Verification)
    - [x] SubTask 5.1: 启动前后端服务（注入测试 Key）。
    - [x] SubTask 5.2: 运行 E2E 测试并录制视频/截图。
    - [x] SubTask 5.3: 人工验证交互流畅性。
