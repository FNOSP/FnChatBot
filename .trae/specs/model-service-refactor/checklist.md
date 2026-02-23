# Checklist

## Skill 使用验证

- [x] 后端开发使用了 `golang-pro` skill
- [ ] 后端测试使用了 `golang-testing` skill
- [x] 前端开发使用了 `vue-expert` skill
- [ ] UI 设计使用了 `ui-ux-pro-max` skill

## 数据模型验证

- [x] Provider 模型已创建并包含所有必要字段
- [x] Model 模型已重构，通过外键关联 Provider
- [x] 数据库迁移成功，现有数据已正确迁移
- [x] 系统预定义供应商已初始化（60 个，不含 CherryIN）
- [x] CherryIN 供应商已从预定义列表中移除
- [x] 数据库初始化时调用 `InitSystemProviders()` 写入预定义供应商

## 适配器验证

- [x] ProviderAdapter 接口定义完整
- [x] OpenAI 适配器正确实现所有接口方法
- [x] Anthropic 适配器正确实现所有接口方法
- [x] Gemini 适配器正确实现所有接口方法
- [x] Ollama 适配器正确实现所有接口方法
- [x] 适配器工厂能根据类型返回正确的适配器

## API 验证

- [x] GET /api/providers 返回所有供应商
- [x] POST /api/providers 创建供应商成功
- [x] PUT /api/providers/:id 更新供应商成功
- [x] DELETE /api/providers/:id 删除供应商成功
- [x] POST /api/providers/:id/fetch-models 获取远程模型列表成功
- [x] GET /api/providers/:id/models 返回供应商下的模型
- [x] POST /api/providers/:id/models 添加模型成功
- [x] PUT /api/models/:id 更新模型成功
- [x] DELETE /api/models/:id 删除模型成功

## 前端验证

- [x] 供应商列表正确显示
- [x] 供应商切换功能正常
- [x] 供应商启用/禁用开关正常工作
- [x] BaseURL 和 API Key 配置正确保存
- [x] 获取远程模型列表功能正常
- [x] 模型添加/编辑/删除功能正常
- [x] 浅色/深色主题下 UI 显示正确

## 聊天功能验证

- [ ] OpenAI 供应商聊天功能正常（需要运行时测试）
- [ ] Anthropic 供应商聊天功能正常（需要运行时测试）
- [ ] Gemini 供应商聊天功能正常（需要运行时测试）
- [ ] Ollama 供应商聊天功能正常（需要运行时测试）
- [ ] 流式响应正确显示（需要运行时测试）
- [ ] 错误处理正确显示（需要运行时测试）

## 国际化验证

- [x] 中文翻译完整
- [x] 英文翻译完整
- [x] 日文翻译完整
- [x] 语言切换后界面正确更新

## 测试覆盖

- [ ] 适配器单元测试覆盖率 >= 80%（待编写）
- [ ] API 处理器测试覆盖率 >= 80%（待编写）
- [ ] 集成测试通过（需要运行时测试）
- [ ] 端到端测试通过（需要运行时测试）
