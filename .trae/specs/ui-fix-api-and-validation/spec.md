# UI修复、API验证与真实对话功能实现规范

## 为什么需要变更
用户反馈了以下问题：
1.  **API配置**：即使填写了Base URL和API Key，点击获取模型列表时仍提示未填写；预览URL逻辑令人困惑或不正确。
2.  **UI显示**：在浅色模式下，“模型服务”区域仍然显示为深色背景。
3.  **国际化（I18n）**：语言选项未显示为母语名称；切换到日语时只有部分界面翻译。
4.  **用户体验（UX）**：设置需要手动保存（用户希望自动保存）。
5.  **核心功能**：聊天功能是模拟（Mock）的，并未实际调用配置的API，导致无法在ChatAnywhere后台看到调用记录。

## 变更内容

### 前端
- **`ModelServices.vue`**：
    - 修复验证逻辑，确保使用当前输入框的值（修复响应式问题）。
    - 优化预览URL逻辑，使其更清晰（显示完整URL或标准路径）。
    - 强制在浅色模式下应用正确的背景色（修复CSS泄漏或类名缺失问题）。
- **`GeneralSettings.vue`**：
    - 更新语言选项显示为母语名称（English, 中文, 日本語）。
    - 移除“保存”按钮；实现 `watch` 以自动保存主题和语言设置。
- **`Sidebar.vue` & `SettingsView.vue`**：
    - 确保所有文本对 `locale` 变更具有响应性（在 `computed` 或模板中使用 `t`）。
    - 在 `ja.json` 中添加缺失的翻译键值。

### 后端
- **`internal/api/ws/websocket.go`**：
    - **破坏性变更**：移除模拟响应逻辑。
    - 实现真实的OpenAI兼容聊天接口调用。
    - 根据对话的 `model_id` 获取 `ModelConfig`（Base URL, API Key）。
    - 将外部API的响应流式传输回WebSocket。

## 影响范围
- **受影响的Spec**：`ui-fix-and-api-config`。
- **受影响的代码**：
    - `frontend/src/components/settings/ModelServices.vue`
    - `frontend/src/components/settings/GeneralSettings.vue`
    - `frontend/src/components/layout/Sidebar.vue`
    - `frontend/src/locales/ja.json`
    - `backend/internal/api/ws/websocket.go`

## 新增需求
### 需求：真实聊天功能
系统必须使用配置的提供商（BaseURL, APIKey）通过外部API生成聊天回复。

### 需求：设置自动保存
通用设置（主题、语言）在变更时必须自动保存。

### 需求：母语语言名称
语言下拉菜单必须显示该语言的母语名称（如 "日本語"）。

## 修改需求
### 需求：API配置验证
当字段已填写时，验证必须立即通过。

### 需求：浅色模式一致性
“模型服务”区域在浅色模式下必须显示为浅色背景。
