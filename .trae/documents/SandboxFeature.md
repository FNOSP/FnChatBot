# Sandbox 功能文档

## 目录

1. [功能说明与使用场景](#1-功能说明与使用场景)
2. [配置方法与参数说明](#2-配置方法与参数说明)
3. [权限请求流程](#3-权限请求流程)
4. [常见问题与解决方案](#4-常见问题与解决方案)
5. [使用示例与最佳实践](#5-使用示例与最佳实践)

---

## 1. 功能说明与使用场景

### 1.1 Sandbox模式的目的和作用

Sandbox（沙箱）模式是 FnChatBot 的核心安全功能，用于限制 AI 助手可访问的文件系统路径。当 Sandbox 模式启用时，AI 只能访问管理员预先配置的允许路径列表中的目录和文件，任何尝试访问未授权路径的操作都将被拦截并触发权限请求流程。

**核心功能：**

- **路径访问控制**：精确控制 AI 可访问的文件系统路径
- **命令解析**：自动解析 AI 执行的命令中包含的路径参数
- **权限请求机制**：当访问未授权路径时，向用户发起权限请求
- **路径记忆功能**：用户可选择将新路径添加到允许列表

### 1.2 安全性保障

Sandbox 模式通过以下机制保障系统安全：

| 安全机制 | 说明 |
|---------|------|
| **路径白名单** | 只有明确配置的路径才允许访问 |
| **子路径继承** | 允许路径的子目录自动获得访问权限 |
| **命令解析** | 支持解析 30+ 种常见命令中的路径参数 |
| **实时拦截** | 在命令执行前进行权限检查 |
| **审计日志** | 记录所有权限请求和响应操作 |

**支持的命令类型：**

系统可自动识别并解析以下命令中的路径参数：

- 文件操作：`cat`, `rm`, `cp`, `mv`, `touch`, `mkdir`
- 目录操作：`cd`, `ls`, `dir`, `find`
- 文本处理：`grep`, `head`, `tail`, `less`, `more`
- 编辑器：`nano`, `vim`, `vi`
- 系统命令：`chmod`, `chown`, `type`, `del`, `copy`, `move`, `xcopy`
- 重定向：`echo >> file`

### 1.3 适用场景

#### 多用户环境

在多用户共享的 FnChatBot 实例中，不同用户可能有不同的数据隔离需求。Sandbox 模式确保：

- 用户 A 无法通过 AI 访问用户 B 的私有目录
- 公共数据目录可配置为共享访问
- 敏感系统目录始终受保护

#### 生产环境

在生产服务器上部署时，Sandbox 模式是必备的安全措施：

```text
推荐配置：
├── /var/www/html          # Web 应用目录
├── /home/app/logs         # 应用日志目录
├── /opt/app/config        # 应用配置目录
└── /tmp/app               # 临时文件目录
```

#### 开发环境

开发环境中可适当放宽限制，但仍建议启用 Sandbox：

```text
推荐配置：
├── ~/Projects             # 项目目录
├── ~/Documents            # 文档目录
└── ~/.config/app          # 应用配置
```

---

## 2. 配置方法与参数说明

### 2.1 如何启用/禁用 Sandbox 模式

#### 通过 Web 界面配置

1. 打开 FnChatBot Web 界面
2. 进入 **设置** → **Sandbox 设置** 页面
3. 找到 **启用 Sandbox 模式** 开关
4. 切换开关状态即可启用或禁用

**界面截图位置说明：**
- 设置页面顶部显示 Sandbox 模式开关
- 开关右侧显示当前状态描述
- 启用后下方显示路径配置区域

#### 通过 API 配置

**获取当前配置：**

```http
GET /api/sandbox
```

响应示例：

```json
{
  "enabled": true,
  "paths": [
    {
      "id": 1,
      "path": "C:\\Users\\Documents",
      "description": "文档目录",
      "enabled": true
    }
  ]
}
```

**启用/禁用 Sandbox：**

```http
PUT /api/sandbox
Content-Type: application/json

{
  "enabled": true
}
```

### 2.2 如何添加允许访问的路径

#### 通过 Web 界面添加

1. 在 Sandbox 设置页面的 **允许的路径** 区域
2. 在路径输入框中输入路径（如：`C:\Users\Documents`）
3. 可选：在描述输入框中添加说明（如：`文档目录`）
4. 点击 **添加路径** 按钮或按 Enter 键确认

**界面截图位置说明：**
- 路径列表上方有两个输入框
- 左侧为路径输入框，右侧为描述输入框
- 最右侧为添加按钮

#### 通过 API 添加

```http
POST /api/sandbox/paths
Content-Type: application/json

{
  "path": "C:\\Users\\Documents",
  "description": "文档目录"
}
```

### 2.3 如何删除已配置的路径

#### 通过 Web 界面删除

1. 在允许的路径列表中找到要删除的路径
2. 点击路径右侧的 **删除** 按钮（垃圾桶图标）
3. 路径将从列表中移除

#### 通过 API 删除

```http
DELETE /api/sandbox/paths/{path}
```

注意：路径需要进行 URL 编码，例如：

```text
原始路径：C:\Users\Documents
编码后：C%3A%5CUsers%5CDocuments
```

### 2.4 配置项说明

#### SandboxConfig 配置

| 字段 | 类型 | 默认值 | 说明 |
|-----|------|-------|------|
| `enabled` | boolean | `false` | 是否启用 Sandbox 模式 |

#### SandboxPath 配置

| 字段 | 类型 | 默认值 | 说明 |
|-----|------|-------|------|
| `id` | uint | 自动生成 | 路径记录 ID |
| `path` | string | 必填 | 允许访问的路径（自动转换为绝对路径） |
| `description` | string | 空 | 路径描述说明 |
| `enabled` | boolean | `true` | 该路径是否启用 |

---

## 3. 权限请求流程

### 3.1 处理流程概述

当 AI 尝试访问未授权路径时，系统会触发以下流程：

```text
┌─────────────────────────────────────────────────────────────────┐
│                    权限请求处理流程                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. AI 发起命令请求                                              │
│     ↓                                                           │
│  2. Sandbox 服务解析命令中的路径                                  │
│     ↓                                                           │
│  3. 检查路径是否在允许列表中                                       │
│     ↓                                                           │
│  ┌─────────────┐                                                │
│  │ 路径已授权？ │                                                │
│  └─────────────┘                                                │
│     │                                                           │
│     ├── 是 → 执行命令                                            │
│     │                                                           │
│     └── 否 → 发起权限请求                                         │
│              ↓                                                  │
│         显示权限请求对话框                                         │
│              ↓                                                  │
│         用户选择响应                                              │
│         ┌────┴────┬──────────┐                                   │
│         ↓         ↓          ↓                                   │
│      [允许]   [允许并记住]  [拒绝]                                  │
│         │         │          │                                   │
│         ↓         ↓          ↓                                   │
│     临时允许   添加到列表    拒绝访问                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 权限请求对话框

当触发权限请求时，用户将看到如下对话框：

**对话框内容说明：**

```
┌────────────────────────────────────────────────────┐
│  AI 请求访问新路径                                   │
├────────────────────────────────────────────────────┤
│                                                    │
│  AI 尝试访问以下路径，但该路径不在允许列表中：          │
│                                                    │
│  路径：                                            │
│  ┌──────────────────────────────────────────────┐  │
│  │ C:\Users\Admin\Private                       │  │
│  └──────────────────────────────────────────────┘  │
│                                                    │
│  执行的命令：                                       │
│  ┌──────────────────────────────────────────────┐  │
│  │ cat C:\Users\Admin\Private\secret.txt        │  │
│  └──────────────────────────────────────────────┘  │
│                                                    │
│                    [拒绝] [允许] [允许并记住]        │
└────────────────────────────────────────────────────┘
```

### 3.3 用户响应选项说明

| 选项 | 行为 | 适用场景 |
|-----|------|---------|
| **允许** | 仅本次允许访问该路径，不保存到配置 | 临时性访问，不需要长期授权 |
| **允许并记住** | 允许访问并将路径添加到允许列表 | 需要长期访问该路径 |
| **拒绝** | 拒绝本次访问请求 | 不希望 AI 访问该路径 |

### 3.4 响应数据结构

用户响应通过 WebSocket 发送：

```typescript
interface PermissionResponse {
  type: 'permission_response';
  requestId: string;    // 权限请求唯一标识
  approved: boolean;    // 是否允许
  remember: boolean;    // 是否记住路径
}
```

---

## 4. 常见问题与解决方案

### Q1: 如何临时禁用 Sandbox？

**解决方案：**

1. 进入 **设置** → **Sandbox 设置**
2. 关闭 **启用 Sandbox 模式** 开关
3. Sandbox 将立即禁用，AI 可访问所有路径

**注意：** 禁用 Sandbox 会降低系统安全性，建议仅在可信环境中临时禁用。

### Q2: 路径格式有什么要求？

**Windows 系统：**

```text
支持的格式：
✓ C:\Users\Documents
✓ C:/Users/Documents
✓ D:\Projects\MyApp

不支持的格式：
✗ Users\Documents        （相对路径会自动转换为绝对路径）
✗ .\Documents            （相对路径会自动转换为绝对路径）
```

**Linux/macOS 系统：**

```text
支持的格式：
✓ /home/user/documents
✓ /var/www/html

不支持的格式：
✗ ~/documents            （会自动展开为完整路径）
✗ ./documents            （相对路径会自动转换为绝对路径）
```

**路径规范化规则：**

- 相对路径自动转换为绝对路径
- Windows 路径分隔符统一为 `\`
- Linux/macOS 路径分隔符统一为 `/`
- Windows 盘符自动大写（如 `c:` → `C:`）

### Q3: 子路径是否自动允许？

**是的，子路径自动继承访问权限。**

例如，配置了 `C:\Users\Documents` 后：

```text
✓ C:\Users\Documents\Work        （允许访问）
✓ C:\Users\Documents\Personal    （允许访问）
✓ C:\Users\Documents\Work\2024   （允许访问）
✗ C:\Users\Downloads             （不允许访问）
✗ C:\Users                       （不允许访问，父目录不在列表中）
```

**判断逻辑：**

```go
func isSubPath(parent, child string) bool {
    // 路径相同则允许
    if parent == child {
        return true
    }
    // 子路径以父路径开头则允许
    return strings.HasPrefix(child, parent + separator)
}
```

### Q4: 如何批量添加路径？

**方法一：通过 API 批量添加**

编写脚本调用 API：

```bash
#!/bin/bash
# 批量添加路径脚本

paths=(
  "/home/user/projects"
  "/home/user/documents"
  "/var/www/html"
  "/opt/app/config"
)

for path in "${paths[@]}"; do
  curl -X POST http://localhost:8080/api/sandbox/paths \
    -H "Content-Type: application/json" \
    -d "{\"path\": \"$path\", \"description\": \"Auto added\"}"
done
```

**方法二：直接操作数据库**

如果使用 SQLite 数据库：

```sql
INSERT INTO sandbox_paths (path, description, enabled) VALUES
  ('/home/user/projects', '项目目录', 1),
  ('/home/user/documents', '文档目录', 1),
  ('/var/www/html', 'Web目录', 1);
```

### Q5: 为什么某些命令没有被拦截？

**可能的原因：**

1. **Sandbox 模式未启用**：检查设置中的开关状态
2. **命令未被识别**：某些复杂命令可能无法正确解析路径
3. **路径已在允许列表中**：检查路径配置

**支持的命令列表：**

如需扩展支持的命令，可在 `sandbox_service.go` 的 `ExtractPathsFromCommand` 方法中添加新的正则表达式模式。

---

## 5. 使用示例与最佳实践

### 示例1：配置文档目录访问

**场景：** 允许 AI 访问用户的文档目录，用于文档管理和搜索。

**配置步骤：**

1. 进入 Sandbox 设置页面
2. 添加路径：
   - 路径：`C:\Users\YourName\Documents`
   - 描述：`个人文档目录`
3. 启用 Sandbox 模式

**配置后效果：**

```text
AI 可以执行：
✓ ls C:\Users\YourName\Documents
✓ cat C:\Users\YourName\Documents\notes.txt
✓ find C:\Users\YourName\Documents -name "*.pdf"

AI 无法执行：
✗ cat C:\Windows\System32\config\SAM
✗ rm C:\Program Files\SomeApp\config.ini
```

### 示例2：配置项目工作区

**场景：** 开发环境中配置多个项目目录。

**推荐配置：**

```json
{
  "enabled": true,
  "paths": [
    {
      "path": "D:\\Projects",
      "description": "项目根目录"
    },
    {
      "path": "D:\\Workspace",
      "description": "工作区目录"
    },
    {
      "path": "C:\\Users\\YourName\\.config",
      "description": "配置文件目录"
    }
  ]
}
```

**配置效果：**

- AI 可以访问所有项目和配置文件
- 系统目录和用户私有目录受保护
- 适合开发调试场景

### 最佳实践建议

#### 1. 最小权限原则

只配置必要的路径，避免添加过于宽泛的目录：

```text
❌ 不推荐：C:\                    （整个 C 盘）
❌ 不推荐：C:\Users               （所有用户目录）
✓  推荐：C:\Users\YourName\Work  （特定工作目录）
```

#### 2. 定期审计路径配置

建议定期检查允许的路径列表：

- 移除不再需要的路径
- 更新路径描述以便管理
- 检查是否有过于宽泛的配置

#### 3. 生产环境配置模板

```json
{
  "enabled": true,
  "paths": [
    {
      "path": "/var/www/html",
      "description": "Web 应用目录"
    },
    {
      "path": "/var/log/app",
      "description": "应用日志目录"
    },
    {
      "path": "/opt/app/config",
      "description": "应用配置目录"
    },
    {
      "path": "/tmp/app",
      "description": "临时文件目录"
    }
  ]
}
```

#### 4. 开发环境配置模板

```json
{
  "enabled": true,
  "paths": [
    {
      "path": "/home/dev/projects",
      "description": "项目目录"
    },
    {
      "path": "/home/dev/.config/app",
      "description": "应用配置"
    },
    {
      "path": "/tmp",
      "description": "临时目录"
    }
  ]
}
```

#### 5. 敏感路径保护

以下路径建议始终避免添加到允许列表：

| 操作系统 | 敏感路径 |
|---------|---------|
| Windows | `C:\Windows`, `C:\Program Files`, `C:\Users\Administrator` |
| Linux | `/etc`, `/root`, `/var/log`, `/proc`, `/sys` |
| macOS | `/System`, `/Library`, `/private` |

---

## 附录

### A. API 接口汇总

| 接口 | 方法 | 说明 |
|-----|------|------|
| `/api/sandbox` | GET | 获取 Sandbox 配置 |
| `/api/sandbox` | PUT | 更新 Sandbox 启用状态 |
| `/api/sandbox/paths` | POST | 添加允许路径 |
| `/api/sandbox/paths/{path}` | DELETE | 删除允许路径 |

### B. 相关源码文件

| 文件 | 说明 |
|-----|------|
| [sandbox_service.go](../backend/internal/services/sandbox_service.go) | Sandbox 核心服务实现 |
| [sandbox_handlers.go](../backend/internal/api/sandbox_handlers.go) | API 处理器 |
| [SandboxSettings.vue](../frontend/src/components/settings/SandboxSettings.vue) | 前端设置页面 |
| [PermissionRequest.vue](../frontend/src/components/PermissionRequest.vue) | 权限请求对话框 |

### C. 版本历史

| 版本 | 日期 | 变更说明 |
|-----|------|---------|
| 1.0.0 | 2024-01 | 初始版本，支持基本路径控制 |
| 1.1.0 | 2024-02 | 添加权限请求流程 |
| 1.2.0 | 2024-03 | 扩展命令解析支持 |

---

*文档最后更新：2026年2月*
