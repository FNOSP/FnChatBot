# 模型服务功能重构规范

## 为什么需要变更

当前 FnChatBot 项目的模型服务功能存在以下问题：

1. **数据模型设计缺陷**：`ModelConfig` 混合了供应商配置和模型配置，导致同一供应商的多个模型需要重复存储 BaseURL 和 API Key
2. **供应商适配缺失**：仅支持 OpenAI 兼容格式，无法适配 Anthropic、Google Gemini 等不同供应商的 API 差异
3. **模型列表获取简陋**：`GetAvailableModels` 仅处理 OpenAI 格式响应，无法适配其他供应商
4. **前端 UI 问题**：Provider 和 Model 概念混淆，缺乏清晰的供应商级别管理
5. **缺乏预定义配置**：没有内置常用供应商的默认配置，用户体验差

## 变更内容

### 后端变更

- **新增 `Provider` 数据模型**：独立管理供应商配置（ID、名称、类型、BaseURL、API Key、启用状态）
- **重构 `ModelConfig` 模型**：移除冗余字段，通过外键关联 Provider
- **新增供应商适配层**：使用策略模式实现不同供应商的 API 适配
- **重构模型列表获取接口**：支持多种供应商格式的模型列表获取
- **新增供应商管理 API**：CRUD 操作供应商配置

### 前端变更

- **重构 `ModelServices.vue`**：实现供应商-模型两级管理界面
- **新增供应商类型定义**：与后端模型对齐
- **优化模型列表获取流程**：支持不同供应商的 API 调用

## 影响范围

- **受影响的代码**：
  - `backend/internal/models/models.go` - 数据模型重构
  - `backend/internal/api/handlers.go` - API 处理逻辑
  - `backend/internal/api/routes.go` - 路由注册
  - `backend/internal/services/` - 新增供应商适配服务
  - `frontend/src/components/settings/ModelServices.vue` - UI 重构
  - `frontend/src/locales/*.json` - 国际化更新

## 新增需求

### 需求：供应商数据模型

系统必须提供独立的供应商（Provider）数据模型，用于存储供应商级别的配置信息。

#### 场景：创建供应商配置

- **WHEN** 用户添加一个新的供应商配置
- **THEN** 系统创建 Provider 记录，包含 ID、名称、类型、BaseURL、API Key、启用状态

#### 场景：供应商关联模型

- **WHEN** 用户在供应商下添加模型
- **THEN** 模型通过外键关联到供应商，继承供应商的连接配置

### 支持的供应商类型

系统复刻 Cherry Studio 的供应商类型，支持 **12 种核心类型**：

| 类型 | 说明 | API 格式 |
|------|------|---------|
| `openai` | OpenAI 及兼容服务（最常用） | OpenAI 格式 |
| `openai-response` | OpenAI Responses API | OpenAI Response 格式 |
| `anthropic` | Anthropic Claude | Anthropic 格式 |
| `gemini` | Google Gemini | Gemini 格式 |
| `azure-openai` | Azure OpenAI | Azure OpenAI 格式 |
| `vertexai` | Google Vertex AI | Vertex AI 格式 |
| `mistral` | Mistral AI | OpenAI 兼容 |
| `aws-bedrock` | AWS Bedrock | Bedrock 格式 |
| `vertex-anthropic` | Vertex AI 上的 Anthropic | Anthropic 格式 |
| `new-api` | New API 兼容 | 扩展 OpenAI |
| `gateway` | Vercel AI Gateway | Gateway 格式 |
| `ollama` | Ollama 本地部署 | Ollama 格式 |

### 预定义供应商列表

系统内置 **60 个预定义供应商**（已移除 CherryIN），包括：

#### 国际主流供应商
- OpenAI (api.openai.com)
- Anthropic (api.anthropic.com)
- Google Gemini (generativelanguage.googleapis.com)
- xAI Grok (api.x.ai)
- Mistral (api.mistral.ai)
- Groq (api.groq.com)
- Perplexity (api.perplexity.ai)
- NVIDIA (integrate.api.nvidia.com)
- Together (api.together.xyz)
- Fireworks (api.fireworks.ai)
- Hyperbolic (api.hyperbolic.xyz)
- Cerebras (api.cerebras.ai)

#### 国内主流供应商
- 智谱 AI (open.bigmodel.cn)
- DeepSeek (api.deepseek.com)
- 月之暗面 Kimi (api.moonshot.cn)
- 百川 AI (api.baichuan-ai.com)
- 阿里百炼 (dashscope.aliyuncs.com)
- 豆包 (ark.cn-beijing.volces.com)
- MiniMax (api.minimaxi.com)
- 腾讯混元 (api.hunyuan.cloud.tencent.com)
- 百度云 (qianfan.baidubce.com)
- 零一万物 (api.lingyiwanwu.com)
- 阶跃星辰 (api.stepfun.com)
- 小米 MiMo (api.xiaomimimo.com)

#### 聚合/代理服务
- AiHubMix (aihubmix.com)
- Silicon (api.siliconflow.cn)
- OpenRouter (openrouter.ai)
- 302.AI (api.302.ai)
- TokenFlux (api.tokenflux.ai)

#### 本地部署
- Ollama (localhost:11434)
- LM Studio (localhost:1234)
- GPUStack (需配置)
- OpenVINO Model Server (localhost:8000)

#### 云服务
- Azure OpenAI (需配置)
- AWS Bedrock (需配置)
- Vertex AI (需配置)
- GitHub Models (models.github.ai)
- GitHub Copilot (api.githubcopilot.com)

#### 专用服务
- Jina (api.jina.ai) - Embedding/Rerank
- VoyageAI (api.voyageai.com) - Embedding/Rerank
- Hugging Face (router.huggingface.co)
- Vercel AI Gateway (ai-gateway.vercel.sh)

### 需求：供应商适配器

系统必须为每种供应商类型提供适配器，处理 API 调用差异。

#### 场景：获取模型列表

- **WHEN** 用户请求获取供应商的可用模型列表
- **THEN** 系统根据供应商类型选择对应的适配器，调用正确的 API 端点并解析响应

#### 场景：聊天补全调用

- **WHEN** 系统调用 AI 模型进行聊天补全
- **THEN** 系统根据供应商类型构建正确的请求格式

### 需求：预定义供应商配置

系统必须内置常用供应商的默认配置。

#### 预定义供应商列表

- OpenAI (api.openai.com)
- Anthropic (api.anthropic.com)
- Google Gemini (generativelanguage.googleapis.com)
- Azure OpenAI (用户需配置部署名)
- Ollama (localhost:11434)

### 需求：前端供应商管理界面

前端必须提供供应商级别的管理界面，支持：

1. 供应商列表展示和切换
2. 供应商启用/禁用
3. 供应商连接配置（BaseURL、API Key）
4. 获取可用模型列表
5. 模型管理（添加、编辑、删除）

## 修改需求

### 需求：模型配置简化

模型配置必须简化，仅保留模型特有属性：

- 模型 ID（如 gpt-4、claude-3-opus）
- 显示名称
- 模型能力标签
- 关联的供应商 ID

### 需求：API 接口重构

现有 API 接口需要重构以支持新的数据模型：

| 原接口 | 新接口 | 说明 |
|--------|--------|------|
| GET /api/models | GET /api/providers | 获取供应商列表 |
| POST /api/models | POST /api/providers | 创建供应商 |
| - | GET /api/providers/:id/models | 获取供应商下的模型 |
| POST /api/models/available | POST /api/providers/:id/fetch-models | 获取远程模型列表 |

## 数据模型设计

### Provider 模型

```go
type Provider struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(100);not null" json:"name"`
    Type      string    `gorm:"type:varchar(50);not null;index" json:"type"` // openai, anthropic, gemini, etc.
    BaseURL   string    `gorm:"type:varchar(500);not null" json:"base_url"`
    APIKey    string    `gorm:"type:varchar(500)" json:"api_key,omitempty"`
    Enabled   bool      `gorm:"default:false;index" json:"enabled"`
    IsSystem  bool      `gorm:"default:false" json:"is_system"` // 系统预定义
    Models    []Model   `gorm:"foreignKey:ProviderID" json:"models,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Model 模型

```go
type Model struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    ProviderID  uint      `gorm:"not null;index" json:"provider_id"`
    Provider    Provider  `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
    ModelID     string    `gorm:"type:varchar(100);not null" json:"model_id"` // 实际模型 ID
    Name        string    `gorm:"type:varchar(100);not null" json:"name"`     // 显示名称
    Capabilities []string `gorm:"type:text;serializer:json" json:"capabilities"` // 能力标签
    IsDefault   bool      `gorm:"default:false" json:"is_default"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### ProviderType 枚举

```go
type ProviderType string

const (
    ProviderTypeOpenAI          ProviderType = "openai"
    ProviderTypeAnthropic       ProviderType = "anthropic"
    ProviderTypeGemini          ProviderType = "gemini"
    ProviderTypeAzureOpenAI     ProviderType = "azure-openai"
    ProviderTypeOllama          ProviderType = "ollama"
    ProviderTypeOpenAICompatible ProviderType = "openai-compatible"
)
```

## API 设计

### 供应商管理 API

```
GET    /api/providers           - 获取所有供应商
POST   /api/providers           - 创建供应商
GET    /api/providers/:id       - 获取单个供应商
PUT    /api/providers/:id       - 更新供应商
DELETE /api/providers/:id       - 删除供应商
POST   /api/providers/:id/fetch-models - 获取远程模型列表
```

### 模型管理 API

```
GET    /api/providers/:id/models - 获取供应商下的模型
POST   /api/providers/:id/models - 添加模型到供应商
PUT    /api/models/:id           - 更新模型
DELETE /api/models/:id           - 删除模型
```

## 适配器设计

### 适配器接口

```go
type ProviderAdapter interface {
    // 获取可用模型列表
    FetchModels(baseURL, apiKey string) ([]ModelInfo, error)
    
    // 构建聊天补全请求
    BuildChatRequest(model string, messages []ChatMessage, options ChatOptions) (*http.Request, error)
    
    // 解析流式响应
    ParseStreamResponse(reader *bufio.Reader) (StreamChunk, error)
    
    // 获取适配器支持的供应商类型
    ProviderType() ProviderType
}
```

### 适配器实现

- `OpenAIAdapter` - OpenAI 及兼容服务
- `AnthropicAdapter` - Anthropic Claude
- `GeminiAdapter` - Google Gemini
- `OllamaAdapter` - Ollama 本地部署

## 迁移策略

1. **数据库迁移**：创建新表结构，迁移现有数据
2. **API 兼容**：保留旧 API 一段时间，标记为废弃
3. **前端渐进**：先实现新界面，后移除旧代码

## 技术实施要求

### 后端开发 Skill

实施后端任务时，必须使用 `golang-pro` skill，确保：

- 使用 Go 1.21+ 特性和最佳实践
- 遵循 Go 代码规范（gofmt、golangci-lint）
- 实现完整的单元测试（覆盖率 >= 80%）
- 使用 context.Context 处理阻塞操作
- 正确处理所有错误

### 后端测试 Skill

编写测试时，必须使用 `golang-testing` skill，确保：

- 使用表驱动测试（table-driven tests）
- 编写基准测试（benchmarks）
- 使用模糊测试（fuzzing）验证边界情况
- 运行竞态检测（-race flag）

### 前端开发 Skill

实施前端任务时，必须使用 `vue-expert` skill，确保：

- 使用 Vue 3 Composition API
- 使用 TypeScript 进行类型安全开发
- 遵循 Vue 3 最佳实践
- 正确使用 Pinia 进行状态管理（如需要）

### UI/UX 设计 Skill

设计界面时，必须使用 `ui-ux-pro-max` skill，确保：

- 遵循现代 UI/UX 设计原则
- 支持浅色/深色主题
- 响应式设计
- 良好的用户体验
