# Tasks

## Skill 使用说明

实施各阶段任务时，必须使用以下 skill：

| 阶段 | 任务范围 | 必须使用的 Skill |
|------|---------|-----------------|
| Phase 1-3 | 后端数据模型、适配器、API | `golang-pro` |
| Phase 5 | 后端测试 | `golang-testing` |
| Phase 4 | 前端重构 | `vue-expert` |
| Phase 4 | UI 设计 | `ui-ux-pro-max` |

## Phase 1: 后端数据模型重构

- [x] Task 1: 创建 Provider 数据模型
  - [x] SubTask 1.1: 在 `backend/internal/models/` 创建 `provider.go`，定义 Provider 结构体
  - [x] SubTask 1.2: 定义 ProviderType 枚举常量（12 种类型）
  - [x] SubTask 1.3: 定义 ProviderApiOptions API 兼容性选项
  - [x] SubTask 1.4: 添加 GORM 模型方法和数据库迁移

- [x] Task 2: 重构 Model 数据模型
  - [x] SubTask 2.1: 创建 `backend/internal/models/model.go`，定义新 Model 结构体
  - [x] SubTask 2.2: 添加 ProviderID 外键关联
  - [x] SubTask 2.3: 添加 Capabilities 能力标签字段
  - [x] SubTask 2.4: 定义 ModelCapability 类型

- [x] Task 3: 创建预定义供应商配置
  - [x] SubTask 3.1: 创建 `backend/internal/config/providers.go`，定义 60 个预定义供应商（已移除 CherryIN）
  - [x] SubTask 3.2: 创建 `backend/internal/config/models.go`，定义预定义模型列表
  - [x] SubTask 3.3: 实现数据库初始化和迁移逻辑

## Phase 2: 供应商适配器实现

- [x] Task 4: 定义适配器接口
  - [x] SubTask 4.1: 在 `backend/internal/services/` 创建 `provider_adapter.go`
  - [x] SubTask 4.2: 定义 ProviderAdapter 接口
  - [x] SubTask 4.3: 定义 ModelInfo、ChatOptions、StreamChunk 等辅助类型

- [x] Task 5: 实现 OpenAI 适配器
  - [x] SubTask 5.1: 创建 `backend/internal/services/adapters/openai_adapter.go`
  - [x] SubTask 5.2: 实现 FetchModels 方法（调用 /v1/models）
  - [x] SubTask 5.3: 实现 BuildChatRequest 方法
  - [x] SubTask 5.4: 实现 ParseStreamResponse 方法

- [x] Task 6: 实现 Anthropic 适配器
  - [x] SubTask 6.1: 创建 `backend/internal/services/adapters/anthropic_adapter.go`
  - [x] SubTask 6.2: 实现 Anthropic API 格式的请求构建
  - [x] SubTask 6.3: 实现流式响应解析（SSE 格式）

- [x] Task 7: 实现 Gemini 适配器
  - [x] SubTask 7.1: 创建 `backend/internal/services/adapters/gemini_adapter.go`
  - [x] SubTask 7.2: 实现 Gemini API 格式的请求构建
  - [x] SubTask 7.3: 实现模型列表获取（listModels API）

- [x] Task 8: 实现 Ollama 适配器
  - [x] SubTask 8.1: 创建 `backend/internal/services/adapters/ollama_adapter.go`
  - [x] SubTask 8.2: 实现 Ollama API 格式（/api/tags、/api/chat）
  - [x] SubTask 8.3: 实现流式响应解析

- [x] Task 9: 实现 Azure OpenAI 适配器
  - [x] SubTask 9.1: 创建 `backend/internal/services/adapters/azure_adapter.go`
  - [x] SubTask 9.2: 实现 Azure 特有的认证和 API 版本处理
  - [x] SubTask 9.3: 实现模型列表获取

- [x] Task 10: 创建适配器工厂
  - [x] SubTask 10.1: 创建 `backend/internal/services/adapters/factory.go`
  - [x] SubTask 10.2: 实现 GetAdapter(providerType) 方法
  - [x] SubTask 10.3: 实现适配器注册机制

## Phase 3: API 接口重构

- [x] Task 11: 创建供应商管理 API
  - [x] SubTask 11.1: 创建 `backend/internal/api/provider_handlers.go`
  - [x] SubTask 11.2: 实现 GetProviders、CreateProvider、UpdateProvider、DeleteProvider
  - [x] SubTask 11.3: 实现 FetchModels 接口（调用适配器获取远程模型）
  - [x] SubTask 11.4: 实现 ToggleProvider 启用/禁用接口

- [x] Task 12: 重构模型管理 API
  - [x] SubTask 12.1: 创建 `backend/internal/api/model_handlers.go`
  - [x] SubTask 12.2: 实现 GetProviderModels、AddModelToProvider、UpdateModel、DeleteModel
  - [x] SubTask 12.3: 更新路由注册

- [x] Task 13: 重构 WebSocket 聊天逻辑
  - [x] SubTask 13.1: 修改 `websocket.go` 使用适配器模式
  - [x] SubTask 13.2: 通过 ProviderID 获取供应商配置
  - [x] SubTask 13.3: 使用适配器构建请求和解析响应

## Phase 4: 前端重构

- [x] Task 14: 更新前端类型定义
  - [x] SubTask 14.1: 创建 `frontend/src/types/provider.ts`
  - [x] SubTask 14.2: 定义 Provider、Model、ProviderType 类型
  - [x] SubTask 14.3: 定义 ProviderApiOptions 类型

- [x] Task 15: 重构 ModelServices.vue
  - [x] SubTask 15.1: 实现供应商列表侧边栏（显示 60 个预定义供应商）
  - [x] SubTask 15.2: 实现供应商配置表单
  - [x] SubTask 15.3: 实现模型管理表格
  - [x] SubTask 15.4: 实现获取远程模型弹窗
  - [x] SubTask 15.5: 实现供应商启用/禁用开关

- [x] Task 16: 更新国际化
  - [x] SubTask 16.1: 更新 `zh.json` 添加供应商相关翻译
  - [x] SubTask 16.2: 更新 `en.json` 添加供应商相关翻译
  - [x] SubTask 16.3: 更新 `ja.json` 添加供应商相关翻译

## Phase 5: 测试与验证

- [x] Task 17: 后端单元测试
  - [x] SubTask 17.1: 编写适配器单元测试
  - [x] SubTask 17.2: 编写 API 处理器测试
  - [x] SubTask 17.3: 编写模型迁移测试

- [x] Task 18: 集成测试
  - [x] SubTask 18.1: 测试供应商 CRUD 操作
  - [x] SubTask 18.2: 测试模型管理操作
  - [x] SubTask 18.3: 测试聊天功能（使用不同供应商）

- [x] Task 19: 端到端测试
  - [x] SubTask 19.1: 测试 OpenAI 供应商完整流程
  - [x] SubTask 19.2: 测试 Anthropic 供应商完整流程
  - [x] SubTask 19.3: 测试 Ollama 本地供应商
  - [x] SubTask 19.4: 测试国内供应商（DeepSeek、智谱等）

## Phase 6: 修复供应商列表为空问题

- [x] Task 20: 修复供应商初始化
  - [x] SubTask 20.1: 在数据库初始化时调用 `config.InitSystemProviders()`
  - [x] SubTask 20.2: 确保 60 个预定义供应商正确写入数据库
  - [x] SubTask 20.3: 测试供应商列表显示

# Task Dependencies

- [Task 2] depends on [Task 1]
- [Task 3] depends on [Task 1, Task 2]
- [Task 5, Task 6, Task 7, Task 8, Task 9] depend on [Task 4]
- [Task 10] depends on [Task 5, Task 6, Task 7, Task 8, Task 9]
- [Task 11, Task 12] depend on [Task 10]
- [Task 13] depends on [Task 10]
- [Task 15] depends on [Task 14]
- [Task 17, Task 18, Task 19] depend on [Task 11, Task 12, Task 13, Task 15]
