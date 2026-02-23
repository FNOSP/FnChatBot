# Tasks

## 语言切换修复

- [x] Task 1: 修复语言切换即时生效问题
  - [x] SubTask 1.1: 分析i18n配置，找出语言切换不生效的原因
  - [x] SubTask 1.2: 修改i18n.ts，确保locale变化时触发响应式更新
  - [x] SubTask 1.3: 在GeneralSettings.vue中确保正确使用useI18n

## 配色方案设计

- [x] Task 2: 设计浅色主题配色方案
  - [x] SubTask 2.1: 设计温暖的浅色背景色系
  - [x] SubTask 2.2: 设计高对比度文字颜色
  - [x] SubTask 2.3: 设计强调色和辅助色
  - [x] SubTask 2.4: 更新style.css中的CSS变量

- [x] Task 3: 设计深色主题配色方案
  - [x] SubTask 3.1: 设计舒适的深色背景色系
  - [x] SubTask 3.2: 设计高对比度文字颜色
  - [x] SubTask 3.3: 设计强调色和辅助色
  - [x] SubTask 3.4: 更新style.css中的深色主题CSS变量

## 模型管理功能

- [x] Task 4: 后端模型列表代理API
  - [x] SubTask 4.1: 创建 /api/models/available 端点
  - [x] SubTask 4.2: 实现代理调用外部AI服务 /v1/models API
  - [x] SubTask 4.3: 处理CORS和错误响应

- [x] Task 5: 前端模型列表获取功能
  - [x] SubTask 5.1: 在ModelServices.vue中添加"获取模型列表"按钮
  - [x] SubTask 5.2: 实现模型列表获取API调用
  - [x] SubTask 5.3: 创建模型选择器弹窗组件
  - [x] SubTask 5.4: 实现模型多选和添加功能

- [x] Task 6: 聊天界面模型切换
  - [x] SubTask 6.1: 在chat store中添加当前模型状态
  - [x] SubTask 6.2: 在ChatView.vue中添加模型选择下拉框
  - [x] SubTask 6.3: 实现模型切换逻辑

## 自动化测试

- [x] Task 7: 准备测试环境
  - [x] SubTask 7.1: 启动后端服务
  - [x] SubTask 7.2: 启动前端开发服务器
  - [x] SubTask 7.3: 配置测试API（api.chatanywhere.tech）

- [x] Task 8: 编写自动化测试脚本
  - [x] SubTask 8.1: 编写语言切换测试脚本
  - [x] SubTask 8.2: 编写主题切换测试脚本
  - [x] SubTask 8.3: 编写模型配置测试脚本
  - [x] SubTask 8.4: 编写AI对话测试脚本（使用gpt-4o-mini-ca）

- [x] Task 9: 执行自动化测试
  - [x] SubTask 9.1: 执行所有测试脚本
  - [x] SubTask 9.2: 收集测试结果和截图
  - [x] SubTask 9.3: 修复发现的问题

# Task Dependencies
- Task 2, Task 3 可并行执行
- Task 4 和 Task 5 可并行执行
- Task 5 依赖 Task 4
- Task 6 依赖 Task 5
- Task 7 需要在 Task 1-6 完成后执行
- Task 8 依赖 Task 7
- Task 9 依赖 Task 8
