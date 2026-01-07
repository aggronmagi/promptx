# Promptx Refactoring Documentation

## 🎯 Project Goals
The primary objective is to refactor `promptx` to reduce redundancy and improve modularity by separating terminal interaction (UI) from command management (Logic).

### Dimension Separation
1.  **Terminal / Interaction (UI Layer)**: Focused on raw terminal operations, screen management, and user-facing prompt blocks (Input, Select).
2.  **Shell / Command Management (Logic Layer)**: Focused on command parsing, lifecycle management, completion algorithms, and history persistence.

---

## 🗺️ Refactoring Plan

### Phase 1: Interface Isolation & Simplification
- [x] **Define Specialized Interfaces**: Split the monolithic `Context` into `Terminal`, `Interaction`, and `Commander`.
- [x] **Streamline Context**: Core interfaces now only contain essential raw methods.
- [x] **Extract Helper Functions**: Moved 20+ type-specific input methods (e.g., `InputInt`, `InputFloat`) from the interface to `promptx_helpers.go`.

### Phase 2: Component Decoupling
- [x] **Command Logic Migration**: Created `commandManager` to handle command sets and history, removing logic from the main `Promptx` struct.
- [x] **Structural Clarity**: `Promptx` now acts as a coordinator between the `TerminalApp` and `commandManager`.

### Phase 3: Option Naming & Constructor Unification
- [x] **Unified Entry Point**: Replaced multiple redundant constructors with a single `New(opts...)`.
- [x] **Simplify Option Names**: Renamed generated option functions for brevity (e.g., `WithCommonOpions` -> `WithCommon`).
- [x] **Fix Typos**: Corrected legacy typos like `CommonOpions`.
- [ ] **Fix Broken References**: (In Progress) Update `promptx_helpers.go` and `command_args.go` to use the new simplified option names.

### Phase 4: Enhanced Argument Handling
- [x] **Struct Binding**: Implemented `ctx.Bind(v interface{})` in `CommandContext` to allow strongly-typed argument retrieval without relying on positional indices.

---

## 📈 Current Progress (Status: 85%)

| Task | Status | Note |
| :--- | :--- | :--- |
| Interface Redefinition | ✅ Completed | Interfaces isolated in `promptx.go`. |
| Logic Decoupling | ✅ Completed | Command management moved to `command_mgr.go`. |
| Constructor Unification| ✅ Completed | `New()` is now the standard entry point. |
| Struct Binding (Args) | ✅ Completed | `Bind()` method implemented and ready for use. |
| Option Naming Cleanup | 🚧 In Progress | Simplified definitions done; fixing downstream usage. |
| Code Regeneration | ✅ Completed | `go generate` run with new naming conventions. |

---

## 🛠️ Key Changes & Examples

### Before (Positional & Coupled)
```go
p := promptx.NewCommandPromptx(cmds...)
p.Run()

func(ctx promptx.CommandContext) {
    id := ctx.CheckInteger(0) // Fragile positional access
}
```

### After (Clean & Isolated)
```go
// Clean constructor with dimension-aware options
p := promptx.New(promptx.WithCommon(promptx.WithCmds(cmds...)))

// Strongly-typed argument binding
type LoginArgs struct {
    Username string
    Retry    int
}
func(ctx promptx.CommandContext) {
    var args LoginArgs
    if err := ctx.Bind(&args); err == nil {
        // use args.Username, args.Retry
    }
}
```

---

## 🚀 Next Steps
1.  **Batch Fix Options**: Replace all occurrences of old `WithInputOption...` style calls in helper files.
2.  **Examples Update**: Sync `_example/` directory with the new decoupled API.
3.  **Final Linting**: Ensure zero linter errors across the project.

