# TDesign 重构与智能对话组件集成规范

## 为什么
当前项目使用 Semi Design 构建界面，虽然功能完备，但用户希望采用 TDesign (腾讯设计体系) 进行替代，并利用 TDesign 提供的 Chatbot 智能对话组件（或参考其设计模式）来重构对话界面，以获得更现代化的 UI 体验和开箱即用的对话功能（如流式传输、状态管理、丰富的消息渲染等）。

## 变更内容
### 前端框架与 UI 库
- **移除 Semi Design**：移除 `package.json` 中的 `@kousum/semi-ui-vue` 及相关样式引入。
- **引入 TDesign**：安装 `tdesign-vue-next` (TDesign Vue Next) 作为主要的 UI 组件库。
- **引入 TDesign Chatbot**：尝试引入 TDesign 提供的 Chatbot 相关组件（如果作为独立包发布），或基于 TDesign 基础组件并参考官方 Chatbot 文档实现类似的组件结构。

### 界面重构
- **全局布局**：使用 TDesign 的 `Layout` (`t-layout`) 组件重构应用的主体结构和侧边栏导航。
- **对话界面 (`ChatView`)**：
  - 重构消息列表 (`ChatList`)：参考 TDesign Chatbot 的 `ChatMessage` 组件设计，实现消息气泡、头像、状态显示。
  - 重构输入区域 (`ChatInput`)：参考 TDesign Chatbot 的 `ChatSender` 组件设计，实现多行输入、快捷指令、文件上传等功能。
  - 增强功能：集成思考过程 (`ChatReasoning`)、加载状态 (`ChatLoading`) 等视觉反馈。
- **设置页面 (`SettingsView`)**：
  - 使用 TDesign 的表单组件 (`t-form`, `t-input`, `t-select`, `t-switch` 等) 重构所有设置项。
  - 使用 `t-tabs` 重构设置页面的导航结构。
  - 保持现有的功能逻辑（如 Token 安全显示、主题切换、i18n 等），但使用新组件实现。
- **其他组件**：
  - 重构 `ModelServices`, `MCPServers`, `SkillManagement` 等管理界面，使用 `t-table` 或 `t-list` 展示数据。
  - 重构弹窗和通知，使用 `t-dialog` and `t-message`。

### 样式与主题
- **主题适配**：适配 TDesign 的主题系统，确保深色/浅色模式正常工作。
- **样式迁移**：将基于 Semi Design 的样式调整为 TDesign 的样式变量或 Tailwind CSS 类名。

## 影响范围
- **受影响的代码**：
  - `package.json`: 依赖变更。
  - `src/main.ts`: 插件注册变更。
  - `src/App.vue`: 全局配置变更。
  - `src/views/**/*.vue`: 所有页面视图。
  - `src/components/**/*.vue`: 所有 UI 组件。
  - `src/assets/styles`: 全局样式文件。

## 新增需求
### 需求：TDesign 集成
系统必须使用 `tdesign-vue-next` 组件库替代原有的 `semi-design-vue`，确保无 Semi Design 残留代码。

### 需求：Chatbot 组件体验
对话界面应具备 TDesign Chatbot 文档中描述的特性，包括但不限于：
- 清晰的消息气泡区分（用户/AI）。
- 思考过程的折叠/展开展示。
- 输入框的自动高度调整和快捷操作。
- 更加现代化的视觉风格。

### 需求：功能一致性
重构后的系统必须保留原有的所有功能，包括模型配置、MCP 服务器管理、Skill 管理、文件上传、Markdown 渲染等。
