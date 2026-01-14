package promptx

import (
	"runtime/debug"
	"strings"

	"github.com/aggronmagi/promptx/v2/blocks"
)

// commandContext 命令执行上下文（内部使用）
type commandContext struct {
	blocks.Context
	// 解析后的命令链
	cmds []*Command
	// 解析后的参数
	args []string
	// 原始输入行
	line string
	// 根命令
	root *Command
	// 当前命令
	cur *Command
}

// parseCommand 解析命令
func parseCommand(root *Command, line string) (*commandContext, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, nil
	}

	ctx := &commandContext{
		Context: nil, // 将在执行时设置
		line:    line,
		root:    root,
		args:    []string{},
		cmds:    []*Command{},
	}

	// 解析命令链
	father := root
	discard := -1

	for k, arg := range fields {
		if cmd := father.findChildCmd(arg); cmd != nil {
			ctx.cmds = append(ctx.cmds, cmd)
			father = cmd
			discard = k + 1
			continue
		}
		discard = k
		break
	}

	// 提取参数
	if discard >= 0 {
		ctx.args = fields[discard:]
	}

	// 设置当前命令
	if len(ctx.cmds) > 0 {
		ctx.cur = ctx.cmds[len(ctx.cmds)-1]
	} else {
		ctx.cur = root
	}

	return ctx, nil
}

// execCommand 执行命令
// 返回 true 表示找到并执行了命令，false 表示未找到命令
func execCommand(ctx blocks.Context, root *Command, line string) bool {
	cmdCtx, err := parseCommand(root, line)
	if err != nil {
		ctx.Printf("解析命令失败: %v\n", err)
		return false
	}

	if cmdCtx == nil {
		return false
	}

	// 检查是否有命令
	if cmdCtx.cur == nil || cmdCtx.cur == root {
		// 找不到命令，返回 false
		return false
	}

	// 解析参数
	var arg any
	if len(cmdCtx.cur.argDefs) > 0 && cmdCtx.cur.argType != nil {
		// 检查参数
		checkedArgs, err := checkArgs(ctx, cmdCtx.cur.argDefs, cmdCtx.args)
		if err != nil {
			ctx.Printf("参数检查失败: %v\n", err)
			return false
		}

		// 根据命令类型创建参数值
		if cmdCtx.cur.isCommander {
			// Commander 类型：使用保存的原始类型
			argValue := createArgValueForCommander(cmdCtx.cur.argDefs, cmdCtx.cur.argType, checkedArgs)
			if argValue.IsValid() {
				arg = argValue.Interface()
			}
		} else {
			// 普通泛型命令：创建指针
			argValue := createArgValueFromStrings(cmdCtx.cur.argDefs, cmdCtx.cur.argType, checkedArgs)
			if argValue.IsValid() && argValue.CanAddr() {
				arg = argValue.Addr().Interface()
			}
		}
	}

	// 执行命令
	defer func() {
		if p := recover(); p != nil {
			ctx.Printf("执行命令失败: %v\n%s", p, string(debug.Stack()))
		}
	}()

	cmdCtx.cur.Exec(ctx, arg)
	return true
}
