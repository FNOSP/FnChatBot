# Sandbox模式功能检查清单

## 后端实现检查

- [x] SandboxConfig数据模型已创建并包含必要字段（enabled, paths）
- [x] SandboxPath数据模型已创建并正确关联
- [x] 数据库迁移成功执行，表结构正确

- [x] GET /api/sandbox 接口返回当前配置
- [x] PUT /api/sandbox 接口更新配置成功
- [x] POST /api/sandbox/paths 接口添加路径成功
- [x] DELETE /api/sandbox/paths/:path 接口删除路径成功

- [x] IsPathAllowed函数正确验证路径权限
- [x] ExtractPathsFromCommand函数正确解析命令中的路径
- [x] CheckCommandPermission函数正确检查命令权限

- [x] WebSocket新增消息类型正确定义
- [x] permission_request消息正确发送
- [x] permission_response消息正确处理
- [x] command_blocked消息正确发送

## 前端实现检查

- [x] SandboxSettings组件正确显示路径列表
- [x] 添加路径功能正常工作
- [x] 删除路径功能正常工作
- [x] Sandbox模式开关功能正常

- [x] PermissionRequest弹窗正确显示请求信息
- [x] 允许按钮正确发送permission_response
- [x] 拒绝按钮正确发送permission_response
- [x] 弹窗样式与项目风格一致

- [x] SettingsView中Sandbox菜单项正确显示
- [x] SandboxSettings组件正确集成到设置页面

## 国际化检查

- [x] zh.json包含所有Sandbox相关翻译
- [x] en.json包含所有Sandbox相关翻译
- [x] ja.json包含所有Sandbox相关翻译
- [x] 界面文本正确使用i18n

## 功能测试检查

- [x] Sandbox模式禁用时命令正常执行
- [x] Sandbox模式启用时已授权路径命令正常执行
- [x] Sandbox模式启用时未授权路径触发权限请求
- [x] 用户批准后路径自动添加到允许列表
- [x] 用户批准后命令继续执行
- [x] 用户拒绝后命令取消执行
- [x] 配置保存后重启应用仍生效

## 代码质量检查

- [x] 后端代码符合Go项目规范
- [x] 前端代码符合Vue 3 + TypeScript规范
- [x] 无明显性能问题
- [x] 错误处理完善
- [x] 代码注释清晰

## 文档检查

- [x] 功能说明文档完整
- [x] 配置方法文档清晰
- [x] 常见问题文档覆盖主要场景
- [x] 使用示例可操作
