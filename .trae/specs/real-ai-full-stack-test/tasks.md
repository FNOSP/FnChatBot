# 任务列表 (Tasks)

- [x] Task 1: 准备测试环境与配置 (Environment Setup)
    - [x] SubTask 1.1: 确认后端支持通过环境变量或启动参数覆盖 `ModelConfig` (确保不依赖数据库中的旧配置，或者在测试脚本中先初始化正确的数据库配置)。
    - [x] SubTask 1.2: 编写一个初始化脚本 `init_real_db.go` 或在测试脚本中调用 API，将用户提供的 API Key 和 Model ID 写入数据库作为默认模型。

- [x] Task 2: 优化工具反馈逻辑 (Refine Tool Feedback)
    - [x] SubTask 2.1: 修改 `tool_service.go`，让 `Skill` 工具返回更逼真的内容（例如模拟读取了一个 `SKILL.md`）。
    - [x] SubTask 2.2: 修改 `Task` 工具返回，包含 "Subagent execution summary..."，确保 LLM 能据此生成下一步回复。

- [x] Task 3: 编写真实 E2E 测试脚本 (Real E2E Script)
    - [x] SubTask 3.1: 创建 `frontend/tests/e2e/real_flow.ts`。
    - [x] SubTask 3.2: 脚本步骤：
        1. 启动后端（注入真实 Key）。
        2. 启动前端。
        3. 配置/确认模型设置。
        4. 发送消息："请制定一个计划来学习 Go 语言并发编程，请使用 TodoWrite 工具列出计划。" (验证 Todo)。
        5. 发送消息："请使用 Task 工具委派一个探索代理来查找当前的 Go 版本。" (验证 Subagent)。
        6. 发送消息："加载 'golang-best-practices' 技能。" (验证 Skill)。
        7. 发送消息："写一个 Hello World 的 Go 程序并用 bash 运行它 (打印输出即可)。" (验证 Sandbox/Bash)。

- [x] Task 4: 执行与验证 (Execute & Verify)
    - [x] SubTask 4.1: 运行测试脚本。
    - [x] SubTask 4.2: 观察控制台输出和前端截图，确认所有工具均被**真实 LLM** 调用。
