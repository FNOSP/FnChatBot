# Tasks

## 后端开发

- [x] Task 1: 创建Sandbox数据模型
  - [x] SubTask 1.1: 创建 backend/internal/models/sandbox.go，定义SandboxConfig结构体
  - [x] SubTask 1.2: 定义SandboxPath结构体存储允许的路径
  - [x] SubTask 1.3: 添加数据库迁移支持

- [x] Task 2: 实现Sandbox API处理器
  - [x] SubTask 2.1: 创建 backend/internal/api/sandbox_handlers.go
  - [x] SubTask 2.2: 实现GetSandboxConfig获取配置接口
  - [x] SubTask 2.3: 实现UpdateSandboxConfig更新配置接口
  - [x] SubTask 2.4: 实现AddSandboxPath添加路径接口
  - [x] SubTask 2.5: 实现RemoveSandboxPath删除路径接口

- [x] Task 3: 实现Sandbox服务层
  - [x] SubTask 3.1: 创建 backend/internal/services/sandbox_service.go
  - [x] SubTask 3.2: 实现路径验证函数IsPathAllowed
  - [x] SubTask 3.3: 实现命令解析函数ExtractPathsFromCommand
  - [x] SubTask 3.4: 实现权限检查函数CheckCommandPermission

- [x] Task 4: 扩展WebSocket消息处理
  - [x] SubTask 4.1: 添加新的消息类型常量 (permission_request, permission_response, command_blocked)
  - [x] SubTask 4.2: 实现权限请求消息发送逻辑
  - [x] SubTask 4.3: 实现权限响应处理逻辑
  - [x] SubTask 4.4: 集成Sandbox服务到WebSocket处理流程

- [x] Task 5: 注册API路由
  - [x] SubTask 5.1: 在routes.go中注册Sandbox相关路由

## 前端开发

- [x] Task 6: 创建Sandbox设置组件
  - [x] SubTask 6.1: 创建 frontend/src/components/settings/SandboxSettings.vue
  - [x] SubTask 6.2: 实现路径列表展示功能
  - [x] SubTask 6.3: 实现添加路径功能（含路径选择器）
  - [x] SubTask 6.4: 实现删除路径功能
  - [x] SubTask 6.5: 实现Sandbox模式开关

- [x] Task 7: 创建权限请求弹窗组件
  - [x] SubTask 7.1: 创建 frontend/src/components/PermissionRequest.vue
  - [x] SubTask 7.2: 实现弹窗UI显示请求路径和原因
  - [x] SubTask 7.3: 实现允许/拒绝按钮交互
  - [x] SubTask 7.4: 实现与WebSocket的消息通信

- [x] Task 8: 集成到设置页面
  - [x] SubTask 8.1: 在SettingsView.vue中添加Sandbox菜单项
  - [x] SubTask 8.2: 导入并注册SandboxSettings组件

- [x] Task 9: 添加国际化支持
  - [x] SubTask 9.1: 在zh.json中添加Sandbox相关翻译
  - [x] SubTask 9.2: 在en.json中添加Sandbox相关翻译
  - [x] SubTask 9.3: 在ja.json中添加Sandbox相关翻译

## 测试

- [x] Task 10: 后端单元测试
  - [x] SubTask 10.1: 编写Sandbox服务路径验证测试
  - [x] SubTask 10.2: 编写命令解析测试
  - [x] SubTask 10.3: 编写API处理器测试

- [x] Task 11: 集成测试
  - [x] SubTask 11.1: 测试正常sandbox路径访问
  - [x] SubTask 11.2: 测试未授权路径访问拦截
  - [x] SubTask 11.3: 测试用户授权流程
  - [x] SubTask 11.4: 测试配置保存与加载

## 文档

- [x] Task 12: 编写功能文档
  - [x] SubTask 12.1: 编写功能说明与使用场景
  - [x] SubTask 12.2: 编写配置方法与参数说明
  - [x] SubTask 12.3: 编写常见问题与解决方案
  - [x] SubTask 12.4: 编写使用示例与最佳实践

# Task Dependencies
- Task 2 depends on Task 1
- Task 3 depends on Task 1
- Task 4 depends on Task 3
- Task 5 depends on Task 2
- Task 7 depends on Task 4
- Task 8 depends on Task 6
- Task 10 depends on Task 1, Task 2, Task 3
- Task 11 depends on Task 1-9
- Task 12 depends on Task 1-11
