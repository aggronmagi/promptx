package promptx

import (
	"fmt"
	"strings"

	"github.com/aggronmagi/promptx/v2/blocks"
	"github.com/aggronmagi/promptx/v2/buffer"
)

type Context = blocks.Context

type Promptx interface {
	blocks.Context
	blocks.Controler
	CommandGroupSwitcher
	Run() error
}

// CommandGroupSwitcher 命令组切换接口
type CommandGroupSwitcher interface {
	SwitchCommandGroup(name string) error
}

// SwitchCommandGroup 切换命令组
// 通过接口判定，如果 ctx 实现了 CommandGroupSwitcher 接口则调用，否则返回错误
func SwitchCommandGroup(ctx blocks.Context, name string) error {
	switcher, ok := ctx.(CommandGroupSwitcher)
	if !ok {
		return fmt.Errorf("context does not implement CommandGroupSwitcher interface")
	}
	return switcher.SwitchCommandGroup(name)
}

type DynamicAddCommander interface {
	AddSubCommands(cmds ...*Command)
}

// AddCommand 添加命令
func AddCommand(ctx blocks.Context, cmds ...*Command) error {
	adder, ok := ctx.(DynamicAddCommander)
	if !ok {
		return fmt.Errorf("context does not implement DynamicAddCommander interface")
	}
	adder.AddSubCommands(cmds...)
	return nil
}

// promptx 主入口结构
type promptx struct {
	// blocks application
	blocks.Application
	// 命令组管理
	groups map[string]*Command
	// 根命令（用于当前命令组）
	root *Command
}

var _ blocks.Context = &promptx{}
var _ CommandGroupSwitcher = &promptx{}
var _ DynamicAddCommander = &promptx{}

// New 创建新的 Promptx 实例
func newPromptx(c *PromptxConfigs) *promptx {
	p := &promptx{
		groups: make(map[string]*Command),
	}

	c.common = append(c.common, blocks.WithCommonOptionExec(func(ctx blocks.Context, command string) {
		p.execCommand(ctx, command)
	}))

	// 构建 blocks application
	options := []blocks.BlocksOption{
		blocks.WithInputs(c.input...),
		blocks.WithSelects(c.selects...),
		blocks.WithCommon(c.common...),
	}
	if c.manager != nil {
		options = append(options, blocks.WithManager(c.manager))
	}
	if c.inputParser != nil {
		options = append(options, blocks.WithInput(c.inputParser))
	}
	if c.outputWriter != nil {
		options = append(options, blocks.WithOutput(c.outputWriter))
	}
	if c.stderrWriter != nil {
		options = append(options, blocks.WithStderr(c.stderrWriter))
	}
	// 最后设置 context. 保证在执行命令时，ctx 是 Promptx 实例
	options = append(options, blocks.WithContext(p))
	// 构建 blocks application
	p.Application = blocks.New(options...)

	// 从配置中加载命令组
	if rootCmd, ok := c.commandGroups[""]; ok {
		p.root = rootCmd
	}
	for _, group := range c.commandGroups {
		if group.config == nil {
			group.config = newRootCommandConfig()
		}
		group.fixChildren()
		p.groups[group.name] = group
		if p.root == nil {
			p.root = group
		}
	}

	// 设置初始补全
	if p.root != nil {
		p.setupCompletion()
	}

	return p
}

// execCommand 执行命令
func (p *promptx) execCommand(ctx blocks.Context, command string) {
	if len(command) == 0 {
		return
	}

	// 执行前检查
	if p.root.config != nil && p.root.config.preCheck != nil {
		if err := p.root.config.preCheck(ctx); err != nil {
			ctx.Printf("precheck failed, %v\n", err)
			return
		}
	}

	execText := command
	isCmd := true
	if p.root.config != nil && p.root.config.commandPrefix != "" {
		if !strings.HasPrefix(command, p.root.config.commandPrefix) {
			isCmd = false
		} else {
			execText = command[len(p.root.config.commandPrefix):]
			execText = strings.TrimSpace(execText)
		}
	}

	find := false
	if isCmd {
		// 执行命令，返回是否找到命令
		find = execCommand(ctx, p.root, execText)
	}

	if !find {
		// 不是命令，调用 OnNonCommand
		if p.root.config != nil && p.root.config.onNonCommand != nil {
			if err := p.root.config.onNonCommand(ctx, command); err != nil {
				ctx.Printf("%v\n", err)
			}
		} else {
			ctx.Printf("command set deal functions. %s\n", command)
		}
	}
}

// SwitchCommandGroup 切换命令组（实现 CommandGroupSwitcher 接口）
func (p *promptx) SwitchCommandGroup(name string) error {
	group, ok := p.groups[name]
	if !ok {
		return fmt.Errorf("command group %s not found", name)
	}

	// 切换 history
	if group.config != nil && group.config.history != "" {
		p.ResetHistoryFile(group.config.history)
	} else {
		// 共享history也只和默认命令组共享.要不然命令组的history就太乱了.

		// 查找默认命令组, 切换前是默认命令组，则不修改.
		if defaultGroup, ok := p.groups[""]; ok && defaultGroup != p.root {
			// 默认命令组有 history 文件，则切换到默认命令组的history文件.
			if defaultGroup.config != nil {
				p.ResetHistoryFile(defaultGroup.config.history)
			}
		}
	}

	// 设置 prompt
	if group.config != nil && group.config.prompt != "" {
		p.SetPrompt(group.config.prompt)
	}

	// 切换根命令
	p.root = group

	// 设置自动补全
	p.setupCompletion()

	// 调用切换回调
	if group.config != nil && group.config.onChange != nil {
		group.config.onChange(p, name)
	}

	return nil
}

// setupCompletion 设置自动补全
func (p *promptx) setupCompletion() {
	if p.root == nil {
		return
	}

	// 确保子命令映射已更新
	p.root.fixChildren()

	manager := p.GetManager()

	// 如果 manager 支持 ApplyOption，则设置 Completer 和 Valid
	if mgr, ok := manager.(interface {
		ApplyOption(opts ...blocks.CommonOption)
	}); ok {
		commandPrefix := ""
		if p.root.config != nil {
			commandPrefix = p.root.config.commandPrefix
		}

		// 设置 WordSeparator
		sep := " "
		if commandPrefix != "" {
			sep += commandPrefix
		}

		// 创建完成器
		completer := createCompleter(p.root, commandPrefix)

		// 设置 Valid 函数来验证命令
		validFunc := func(status int, doc *buffer.Document) error {
			// 只在 FinishStatus 时验证，NormalStatus 时不验证（除非 AlwaysCheck）
			if status == blocks.NormalStatus {
				return nil
			}
			if len(doc.Text) == 0 {
				return nil
			}
			text := doc.Text
			if commandPrefix != "" {
				if !strings.HasPrefix(text, commandPrefix) {
					return nil
				}
				text = text[len(commandPrefix):]
				text = strings.TrimSpace(text)
			}
			// 解析命令检查是否存在
			cmdCtx, err := parseCommand(p.root, text)
			if err != nil {
				return err
			}
			if cmdCtx == nil || cmdCtx.cur == nil || cmdCtx.cur == p.root {
				return fmt.Errorf("not found command[%s]", doc.Text)
			}
			return nil
		}

		// 应用选项
		mgr.ApplyOption(
			blocks.WithCommonOptionValid(validFunc),
			blocks.WithCommonOptionComplete(
				blocks.WithCompleteOptionCompleter(completer),
				blocks.WithCompleteOptionCompletionFillSpace(true),
				blocks.WithCompleteOptionWordSeparator(sep),
			),
		)
	}
}

func (p *promptx) AddSubCommands(cmds ...*Command) {
	p.root.subCommands = append(p.root.subCommands, cmds...)
	p.root.fixChildren()
}

// Print 打印（blocks.Application 没有此方法，需要自定义）
func (p *promptx) Print(v ...interface{}) {
	fmt.Fprint(p.Stdout(), v...)
}

// Printf 格式化打印（blocks.Application 没有此方法，需要自定义）
func (p *promptx) Printf(format string, v ...interface{}) {
	fmt.Fprintf(p.Stdout(), format, v...)
}

// Println 打印并换行（blocks.Application 没有此方法，需要自定义）
func (p *promptx) Println(v ...interface{}) {
	fmt.Fprintln(p.Stdout(), v...)
}
