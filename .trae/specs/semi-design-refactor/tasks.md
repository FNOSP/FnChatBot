# 任务列表

- [x] 任务 1：设置 Semi Design 和 Vue I18n
  - [x] 子任务 1.1：安装 `semi-design-vue` 和 `@douyinfe/semi-ui`。
  - [x] 子任务 1.2：安装 `vue-i18n`。
  - [x] 子任务 1.3：配置 `i18n` 实例和语言文件（`en`, `zh`, `ja`）。
  - [x] 子任务 1.4：更新 `main.ts` 以包含 Semi Design 和 I18n 插件。
  - [x] 子任务 1.5：配置深色/浅色主题的全局样式。

- [x] 任务 2：创建通用设置组件
  - [x] 子任务 2.1：创建 `src/components/settings/GeneralSettings.vue`。
  - [x] 子任务 2.2：添加语言切换器（使用 `SemiSelect`）。
  - [x] 子任务 2.3：添加主题切换器（使用 `SemiSwitch`）。
  - [x] 子任务 2.4：实现主题持久化逻辑（`localStorage`）。

- [x] 任务 3：重构设置视图布局
  - [x] 子任务 3.1：使用 Semi Design 的 `Layout` 和 `Nav` 组件替换 Tailwind 侧边栏。
  - [x] 子任务 3.2：更新 `SettingsView.vue` 以使用新布局。
  - [x] 子任务 3.3：确保移动端响应式设计。

- [x] 任务 4：重构模型服务组件
  - [x] 子任务 4.1：使用 `SemiInput` 替换现有输入框（Base URL, API Key）。
  - [x] 子任务 4.2：实现 URL 预览逻辑（将端点路径追加到 Base URL）。
  - [x] 子任务 4.3：实现安全 Token 输入（默认掩码显示，支持切换可见性）。
  - [x] 子任务 4.4：修复“保存更改”后 Token 消失的问题。
  - [x] 子任务 4.5：使用 Semi Design 的 `Table` 或 `List` 替换表格/列表。

- [x] 任务 5：重构其他组件
  - [x] 子任务 5.1：重构 `MCPServers.vue` 以使用 Semi Design 组件。
  - [x] 子任务 5.2：重构 `SkillManagement.vue` 以使用 Semi Design 组件。
  - [x] 子任务 5.3：确保文件上传对话框使用 Semi Design 的 `Modal` 和 `Upload`。

- [x] 任务 6：最终验证与 I18n 检查
  - [x] 子任务 6.1：验证所有硬编码字符串已移动到语言文件中。
  - [x] 子任务 6.2：验证主题切换在所有页面上正常工作。
  - [x] 子任务 6.3：验证 Token 安全功能。

# 任务依赖
- 任务 2, 3, 4, 5 依赖于 任务 1（环境设置）。
- 任务 6 依赖于所有其他任务。
