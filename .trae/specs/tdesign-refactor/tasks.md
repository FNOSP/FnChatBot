# 任务列表

- [x] 任务 1：环境准备与依赖迁移
  - [x] 子任务 1.1：卸载 `@kousum/semi-ui-vue` 及相关依赖。
  - [x] 子任务 1.2：安装 `tdesign-vue-next`。
  - [x] 子任务 1.3：配置 `main.ts` 以引入 TDesign。
  - [x] 子任务 1.4：清理项目中的 Semi Design 样式引用，配置 TDesign 的全局样式（包括暗色模式适配）。

- [x] 任务 2：全局布局重构
  - [x] 子任务 2.1：重构 `src/App.vue` 和 `src/views/SettingsView.vue` (布局部分)，使用 `t-layout`, `t-header`, `t-aside`, `t-content` 等组件。
  - [x] 子任务 2.2：重构侧边栏导航，使用 `t-menu` 组件。
  - [ ] 子任务 2.3：确保布局在移动端的响应式表现。

- [x] 任务 3：对话界面 (`ChatView`) 重构
  - [x] 子任务 3.1：研究并尝试引入 TDesign Chatbot 组件（如果存在独立包），或者基于 `tdesign-vue-next` 创建 `ChatMessage`, `ChatList`, `ChatSender` 等组件结构。
  - [x] 子任务 3.2：重构消息列表渲染逻辑，适配新的组件结构，支持 Markdown、代码高亮。
  - [x] 子任务 3.3：实现“思考过程” (`Thinking`/`Reasoning`) 的展示组件，支持折叠/展开。
  - [x] 子任务 3.4：重构输入区域，支持多行文本、发送按钮状态控制。
  - [x] 子任务 3.5：确保流式输出时的滚动和加载状态正常。

- [x] 任务 4：设置页面与功能组件重构
  - [x] 子任务 4.1：重构 `GeneralSettings.vue`（语言、主题设置），使用 `t-select`, `t-switch`。
  - [x] 子任务 4.2：重构 `ModelServices.vue`，使用 `t-input` (带 `suffix` 图标), `t-list` 或 `t-table`。
  - [x] 子任务 4.3：重构 `MCPServers.vue` 和 `SkillManagement.vue`，使用 TDesign 的数据展示组件。
  - [x] 子任务 4.4：重构 `UserManagement.vue` 和 `SandboxSettings.vue`。
  - [x] 子任务 4.5：重构所有弹窗 (`t-dialog`) 和通知 (`t-message`)。

- [x] 任务 5：验证与优化
  - [x] 子任务 5.1：验证所有页面的主题切换效果。
  - [x] 子任务 5.2：验证国际化 (i18n) 文本是否正确显示。
  - [x] 子任务 5.3：进行全流程测试（对话、设置保存、MCP 工具调用等），修复发现的 UI bug。
  - [x] 子任务 5.4：检查并优化构建体积（按需引入 TDesign 组件）。

# 任务依赖
- 任务 2, 3, 4 依赖于 任务 1。
- 任务 5 依赖于所有其他任务。
