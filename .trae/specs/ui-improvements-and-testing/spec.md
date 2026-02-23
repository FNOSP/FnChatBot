# UI改进与模型管理功能规格说明

## Why
当前存在以下问题：
1. 语言切换后页面没有立即更新，需要刷新才能生效
2. 现有配色方案不够美观现代，文字可读性需要提升
3. 模型管理功能不完善，无法从API获取可用模型列表
4. 聊天界面无法切换模型

## What Changes
- 修复语言切换即时生效问题
- 设计全新的深色/浅色配色方案
- 实现模型列表API获取功能
- 实现模型选择器组件
- 在聊天界面添加模型切换功能
- 制定完整的测试计划和自动化测试

## Impact
- Affected specs: 系统设置、模型配置、聊天界面
- Affected code: 
  - frontend/src/composables/useTheme.ts
  - frontend/src/style.css
  - frontend/src/components/settings/ModelServices.vue
  - frontend/src/views/ChatView.vue
  - frontend/src/store/chat.ts
  - frontend/src/i18n.ts
  - backend/internal/api/ (新增模型列表代理API)

## ADDED Requirements

### Requirement: 语言切换即时生效
系统SHALL在用户切换语言后立即更新所有界面文本，无需刷新页面。

#### Scenario: 切换语言后界面更新
- **WHEN** 用户在设置页面切换语言
- **THEN** 所有界面文本立即更新为新语言
- **AND** 语言设置保存到localStorage

### Requirement: 现代化配色方案
系统SHALL提供美观、现代、高可读性的配色方案。

#### Scenario: 浅色主题配色
- **GIVEN** 用户使用浅色主题
- **THEN** 界面使用温暖的浅色背景
- **AND** 文字对比度达到WCAG AA标准（4.5:1以上）
- **AND** 强调色醒目但不刺眼

#### Scenario: 深色主题配色
- **GIVEN** 用户使用深色主题
- **THEN** 界面使用舒适的深色背景
- **AND** 文字对比度达到WCAG AA标准
- **AND** 减少蓝光，保护眼睛

### Requirement: 模型列表API获取
系统SHALL支持从AI服务API获取可用模型列表。

#### Scenario: 获取模型列表
- **GIVEN** 用户已配置API Base URL和API Key
- **WHEN** 用户点击"获取模型列表"按钮
- **THEN** 系统调用 /v1/models API获取模型列表
- **AND** 显示可用模型供用户选择

#### Scenario: 添加模型到配置
- **GIVEN** 模型列表已加载
- **WHEN** 用户选择若干模型并确认
- **THEN** 选中的模型添加到用户模型配置列表

### Requirement: 聊天界面模型切换
系统SHALL在聊天界面提供模型切换功能。

#### Scenario: 切换聊天模型
- **GIVEN** 用户已配置多个模型
- **WHEN** 用户在聊天界面选择不同模型
- **THEN** 后续消息使用新选择的模型

### Requirement: 自动化测试
系统SHALL通过Chrome浏览器自动化测试验证所有功能。

#### Scenario: 语言切换测试
- **WHEN** 执行自动化测试
- **THEN** 验证语言切换后界面文本正确更新

#### Scenario: 主题切换测试
- **WHEN** 执行自动化测试
- **THEN** 验证主题切换后颜色正确应用

#### Scenario: 模型管理测试
- **WHEN** 执行自动化测试
- **THEN** 验证模型列表获取、添加、切换功能

#### Scenario: AI对话测试
- **WHEN** 执行自动化测试
- **THEN** 使用gpt-4o-mini-ca模型进行完整对话测试

## MODIFIED Requirements

### Requirement: 模型服务配置界面
原有的模型配置界面需要扩展，新增以下功能：
- "获取模型列表"按钮，调用外部API
- 模型选择器，支持多选
- 模型快速添加功能

## REMOVED Requirements
无移除的需求。
