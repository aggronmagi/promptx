# Promptx 重构文档

## 🎯 项目目标
主要目标是通过将终端交互（UI）与命令管理（逻辑）分离，重构 `promptx` 以减少冗余并提高模块化。

### 维度分离
1.  **终端 / 交互（UI 层）**：专注于原始终端操作、屏幕管理和面向用户的提示块（Input、Select）。
2.  **Shell / 命令管理（逻辑层）**：专注于命令解析、生命周期管理、补全算法和历史持久化。

---

## 🗺️ 重构计划

### 阶段 1：接口隔离与简化
- [x] **定义专用接口**：将单一的 `Context` 拆分为 `Terminal`、`Interaction` 和 `Commander`。
- [x] **精简 Context**：核心接口现在仅包含必要的原始方法。
- [x] **提取辅助函数**：将 20+ 个类型特定的输入方法（如 `InputInt`、`InputFloat`）从接口移至 `promptx_helpers.go`。

### 阶段 2：组件解耦
- [x] **命令逻辑迁移**：创建 `commandManager` 来处理命令集和历史记录，从主 `Promptx` 结构体中移除逻辑。
- [x] **结构清晰化**：`Promptx` 现在充当 `TerminalApp` 和 `commandManager` 之间的协调者。

### 阶段 3：选项命名与构造函数统一
- [x] **统一入口点**：用单个 `New(opts...)` 替换多个冗余构造函数。
- [x] **简化选项名称**：重命名生成的选项函数以简化（如 `WithCommonOpions` -> `WithCommon`）。
- [x] **修复拼写错误**：纠正遗留的拼写错误，如 `CommonOpions`。
- [x] **修复损坏的引用**：更新 `promptx_helpers.go`、`command_args.go` 和内部文件以使用新的简化选项名称。

### 阶段 4：增强的参数处理
- [x] **结构体绑定**：在 `CommandContext` 中实现 `ctx.Bind(v interface{})`，允许强类型参数检索，而不依赖位置索引。

---

## 📈 当前进度（状态：100%）

| 任务 | 状态 | 备注 |
| :--- | :--- | :--- |
| 接口重定义 | ✅ 已完成 | 接口在 `promptx.go` 中隔离。 |
| 逻辑解耦 | ✅ 已完成 | 命令管理移至 `command_mgr.go`。 |
| 构造函数统一 | ✅ 已完成 | `New()` 现在是标准入口点。 |
| 结构体绑定（参数） | ✅ 已完成 | `Bind()` 方法已实现并可供使用。 |
| 选项命名清理 | ✅ 已完成 | 简化定义并修复下游使用。 |
| 代码重新生成 | ✅ 已完成 | 使用新的命名约定运行 `go generate`。 |

---

## 🛠️ 关键变更与示例

### 之前（位置访问与耦合）
```go
p := promptx.NewCommandPromptx(cmds...)
p.Run()

func(ctx promptx.CommandContext) {
    id := ctx.CheckInteger(0) // 脆弱的位置访问
}
```

### 之后（清晰与隔离）
```go
// 使用维度感知选项的清晰构造函数
p := promptx.New(promptx.WithCommon(promptx.WithCmds(cmds...)))

// 强类型参数绑定
type LoginArgs struct {
    Username string
    Retry    int
}
func(ctx promptx.CommandContext) {
    var args LoginArgs
    if err := ctx.Bind(&args); err == nil {
        // 使用 args.Username, args.Retry
    }
}
```

---

## 🚀 下一步
1.  **发布**：使用重构后的 API 标记新版本。
2.  **文档**：更新 README 和内部文档，包含新的使用模式。
