# Semi Design 重构规范

## 为什么
当前前端界面较为简陋，用户希望使用 `semi-design-vue` 进行重构，以提升视觉效果和一致性，并增加主题切换、模型配置增强、API Token 安全处理以及国际化（i18n）支持等功能。

## 变更内容

### 前端框架与 UI 库
- **迁移至 Semi Design Vue**：将现有的 UI 组件（基于 Tailwind 的自定义组件）替换为 `semi-design-vue` 组件。
- **主题优化**：
  - 实现深色/浅色模式切换。
  - 深色模式主色调：`#2a2a2a`，文字 `#ffffff`。
  - 浅色模式主色调：`#faf8f5`，文字 `#1a1a1a`。
  - 使用 `semi-design` 的 `ConfigProvider` 进行主题管理。
  - 将主题偏好持久化存储在 `localStorage` 中。

### 模型配置增强
- **实时 URL 预览**：在 `baseUrl` 输入框下方添加预览区域，显示完整的 API 端点路径。
- **输入组件**：使用 `SemiInput` 组件，并利用 `addonAfter` 属性展示 URL 预览。

### API Token 安全
- **Token 脱敏**：默认将 API Token 显示为掩码字符（如 `••••••••`）。
- **可见性切换**：添加眼睛图标按钮，用于在明文和掩码之间切换。
- **持久化修复**：确保在切换可见性或保存其他设置时，Token 值不会丢失（修复“保存更改”后 Token 消失的问题）。

### 国际化 (i18n)
- **集成 Vue I18n**：安装并配置 `vue-i18n`。
- **语言支持**：英语 (en)、中文 (zh)、日语 (ja)。
- **语言切换器**：在新的“通用设置”部分添加语言选择器，使用 `SemiSelect` 组件。
- **资源管理**：将所有硬编码字符串提取到语言文件中（`locales/en.json`, `locales/zh.json`, `locales/ja.json`）。

## 影响范围
- **受影响的规范**：`PRD.md`（UI/UX 更新）。
- **受影响的代码**：
  - `package.json`: 添加 `semi-design-vue`, `vue-i18n` 依赖。
  - `src/App.vue`: 设置 Theme Provider。
  - `src/views/SettingsView.vue`: 布局和组件更新。
  - `src/components/settings/ModelServices.vue`: Token 输入和 URL 预览逻辑。
  - `src/components/settings/GeneralSettings.vue`: 新增用于语言和主题设置的组件。
  - `src/locales/*`: 新增翻译文件。
  - 全局样式：主题颜色定义。

## 新增需求
### 需求：Semi Design 迁移
系统必须使用 `semi-design-vue` 组件来实现所有 UI 元素，确保样式和行为的一致性。

### 需求：主题切换
系统必须支持在深色 (#2a2a2a) 和浅色 (#faf8f5) 主题之间切换，并持久化保存用户的选择。

### 需求：安全的 Token 显示
系统必须默认掩盖 API Token，并提供切换按钮以显示明文。保存配置时，如果 Token 未被修改，不得清除原有的 Token 值。

### 需求：国际化
系统必须支持英语、中文和日语，并在设置中提供语言切换器。
