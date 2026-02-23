# Sandbox模式功能规格说明

## Why
当前系统中AI生成的bash命令可以访问任意路径，存在安全风险。需要实现sandbox模式，限制AI只能访问用户授权的路径，当尝试访问未授权路径时触发权限请求机制，确保系统安全性。

## What Changes
- 新增Sandbox配置数据模型，存储允许访问的路径列表
- 新增Sandbox API端点，用于管理路径配置
- 实现命令执行拦截器，检查路径访问权限
- 新增权限请求WebSocket消息类型
- 前端新增Sandbox设置组件
- 前端新增权限请求弹窗组件
- 添加中英文国际化支持

## Impact
- Affected specs: 系统设置、命令执行、WebSocket通信
- Affected code: 
  - backend/internal/models/ (新增sandbox.go)
  - backend/internal/api/ (新增sandbox_handlers.go, 修改routes.go)
  - backend/internal/services/ (新增sandbox_service.go)
  - backend/internal/api/ws/ (修改websocket.go)
  - frontend/src/components/settings/ (新增SandboxSettings.vue)
  - frontend/src/components/ (新增PermissionRequest.vue)
  - frontend/src/locales/ (修改zh.json, en.json)

## ADDED Requirements

### Requirement: Sandbox配置管理
系统SHALL提供Sandbox配置管理功能，允许用户配置可访问的路径列表。

#### Scenario: 添加允许访问的路径
- **WHEN** 用户在设置页面添加新路径
- **THEN** 系统验证路径格式并保存到数据库

#### Scenario: 删除已配置的路径
- **WHEN** 用户删除某个已配置的路径
- **THEN** 系统从配置列表中移除该路径

#### Scenario: 启用/禁用Sandbox模式
- **WHEN** 用户切换Sandbox模式开关
- **THEN** 系统保存设置并立即生效

### Requirement: 路径访问权限控制
系统SHALL在执行bash命令前检查路径访问权限。

#### Scenario: 访问已授权路径
- **GIVEN** Sandbox模式已启用
- **AND** 路径在允许列表中
- **WHEN** AI尝试执行访问该路径的命令
- **THEN** 命令正常执行

#### Scenario: 访问未授权路径
- **GIVEN** Sandbox模式已启用
- **AND** 路径不在允许列表中
- **WHEN** AI尝试执行访问该路径的命令
- **THEN** 系统拦截命令并发送权限请求

### Requirement: 权限请求机制
系统SHALL在检测到未授权路径访问时触发权限请求流程。

#### Scenario: 显示权限请求弹窗
- **WHEN** 系统检测到未授权路径访问
- **THEN** 向用户显示权限请求弹窗
- **AND** 弹窗包含请求的路径和原因

#### Scenario: 用户批准路径访问
- **WHEN** 用户点击"允许"按钮
- **THEN** 系统将路径添加到允许列表
- **AND** 继续执行被拦截的命令

#### Scenario: 用户拒绝路径访问
- **WHEN** 用户点击"拒绝"按钮
- **THEN** 系统取消命令执行
- **AND** 向AI返回权限拒绝消息

### Requirement: WebSocket消息扩展
系统SHALL扩展WebSocket协议以支持权限请求。

#### Scenario: 发送权限请求消息
- **WHEN** 检测到未授权路径访问
- **THEN** 发送类型为"permission_request"的WebSocket消息
- **AND** 消息包含请求路径和命令详情

#### Scenario: 接收权限响应
- **WHEN** 用户响应权限请求
- **THEN** 发送类型为"permission_response"的WebSocket消息
- **AND** 消息包含用户决策

### Requirement: 前端设置界面
系统SHALL在前端设置页面提供Sandbox配置界面。

#### Scenario: 显示Sandbox设置
- **WHEN** 用户进入设置页面
- **THEN** 显示Sandbox设置选项卡
- **AND** 显示当前配置的路径列表

#### Scenario: 路径配置操作
- **WHEN** 用户在Sandbox设置中操作
- **THEN** 可以添加、删除路径
- **AND** 可以启用/禁用Sandbox模式

## MODIFIED Requirements

### Requirement: WebSocket消息类型扩展
原有的WebSocket消息类型需要扩展，新增以下类型：
- `permission_request`: 权限请求消息
- `permission_response`: 权限响应消息
- `command_blocked`: 命令被拦截消息

## REMOVED Requirements
无移除的需求。
