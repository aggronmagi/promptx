package promptx

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/aggronmagi/promptx/v2/blocks"
)

// Commander 自定义命令接口
// 实现此接口可以创建自定义的命令
type Commander interface {
	// Name 返回命令名称
	Name() string
	// Help 返回命令帮助信息
	Help() string
	// Exec 执行命令
	Exec(ctx blocks.Context)
}

type rootCommandConfig struct {
	// 命令前缀
	commandPrefix string
	// 非命令处理函数
	onNonCommand func(ctx blocks.Context, command string) error
	// 执行前检查函数
	preCheck func(ctx blocks.Context) error
	// Prompt 文字
	prompt string
	// History 文件路径（空字符串表示共享 history）
	history string
	// 切换时的回调函数
	onChange func(ctx blocks.Context, args ...interface{})
}

func newRootCommandConfig() *rootCommandConfig {
	return &rootCommandConfig{
		onNonCommand: nil,
		preCheck:     nil,
		prompt:       ">>> ",
		history:      "",
		onChange:     nil,
	}
}

// Command 表示一个命令
type Command struct {
	// 配置
	config *rootCommandConfig
	// 命令名称
	name string
	// 命令帮助信息
	help string
	// 命令别名
	aliases []string
	// 子命令列表
	subCommands []*Command
	// 子命令映射表（用于快速查找）
	children map[string]*Command
	// 命令执行函数
	// 统一签名：func(ctx blocks.Context, arg any)
	// 对于泛型命令，arg 是解析后的 ARG 结构体指针
	// 对于自定义命令，arg 是实现接口的结构体指针
	execFunc func(ctx blocks.Context, arg any)
	// 参数定义（用于泛型命令和自定义命令）
	argDefs []*ArgDef
	// 参数类型 新建参数结构体的反射类型, 保存解析命令行参数的值, 传递给execFunc.
	argType reflect.Type
	// 是否为 Commander 类型命令
	isCommander bool
}

// NewCommandWithFunc 创建一个新的命令（泛型方式）
// name: 命令名称
// help: 命令帮助信息
// arg: 参数结构体指针（用于解析参数定义）
// run: 命令执行函数
func NewCommandWithFunc[ARG any](
	name, help string,
	run func(ctx blocks.Context, arg *ARG),
) *Command {
	cmd := &Command{
		name:     name,
		help:     help,
		children: make(map[string]*Command),
	}
	arg := new(ARG)

	// 解析参数定义
	cmd.argDefs = parseArgDefs(arg)

	// 获取 ARG 类型
	argType := reflect.TypeOf(arg)
	if argType.Kind() == reflect.Ptr {
		argType = argType.Elem()
	}

	// 保存参数类型
	cmd.argType = argType

	// 生成闭包函数
	// 注意：这里不解析参数，参数解析在 Exec 方法中统一处理
	cmd.execFunc = func(ctx blocks.Context, arg any) {
		if argPtr, ok := arg.(*ARG); ok {
			run(ctx, argPtr)
		} else {
			panic(fmt.Sprintf("except type:%#T, got type:%#T", argType, arg))
		}
	}

	return cmd
}

// NewCommandWithFuncLegacy 创建一个新的命令（泛型方式）
// name: 命令名称
// help: 命令帮助信息
// arg: 参数结构体指针（用于解析参数定义）
// run: 命令执行函数
func NewCommandWithFuncLegacy(
	name, help string,
	run func(ctx blocks.Context),
) *Command {
	cmd := &Command{
		name:     name,
		help:     help,
		children: make(map[string]*Command),
	}
	// 解析参数定义
	cmd.argDefs = nil

	// 保存参数类型
	cmd.argType = reflect.TypeOf(struct{}{})

	// 生成闭包函数
	// 注意：这里不解析参数，参数解析在 Exec 方法中统一处理
	cmd.execFunc = func(ctx blocks.Context, arg any) {
		run(ctx)
	}

	return cmd
}

// NewCommand 创建一个自定义命令
// 如果 custom 是结构体类型，会解析其字段作为参数定义
func NewCommand(custom Commander) *Command {
	cmd := &Command{
		name:        custom.Name(),
		help:        custom.Help(),
		children:    make(map[string]*Command),
		isCommander: true, // 标记为 Commander 类型
	}

	// 保存原始类型（包括是否是指针）
	originalType := reflect.TypeOf(custom)

	// 提取结构体类型用于解析字段
	customValue := reflect.ValueOf(custom)
	if customValue.Kind() == reflect.Ptr {
		customValue = customValue.Elem()
	}
	structType := customValue.Type()

	// 检查是否是结构体类型
	if structType.Kind() == reflect.Struct {
		// 创建参数实例用于解析参数定义
		argValue := reflect.New(structType)

		// 解析参数定义
		cmd.argDefs = parseArgDefs(argValue.Interface())
		// 保存原始类型（可能是指针）
		cmd.argType = originalType
	} else {
		panic(fmt.Sprintf("except struct type, got type:%#T", structType))
	}

	// 生成执行函数
	// 有参数，需要传递解析后的参数
	cmd.execFunc = func(ctx blocks.Context, arg any) {
		if argPtr, ok := arg.(Commander); ok {
			argPtr.Exec(ctx)
		} else {
			panic(fmt.Sprintf("except type:%#T, got type:%#T", originalType, arg))
		}
	}

	return cmd
}

// Name 返回命令名称
func (c *Command) Name() string {
	return c.name
}

// Help 返回命令帮助信息
func (c *Command) Help() string {
	return c.help
}

// Aliases 设置命令别名
func (c *Command) Aliases(aliases ...string) *Command {
	c.aliases = aliases
	return c
}

// SubCommands 添加子命令
func (c *Command) SubCommands(cmds ...*Command) *Command {
	for _, cmd := range cmds {
		c.subCommands = append(c.subCommands, cmd)
	}
	c.fixChildren()
	return c
}

// fixChildren 修复子命令映射表
func (c *Command) fixChildren() {
	if len(c.subCommands) == len(c.children) {
		return
	}
	c.children = make(map[string]*Command, len(c.subCommands))

	for _, cmd := range c.subCommands {
		c.children[cmd.name] = cmd
		// 添加别名映射
		for _, alias := range cmd.aliases {
			c.children[alias] = cmd
		}
		// 递归修复子命令
		cmd.fixChildren()
	}

	// 排序子命令
	sort.Slice(c.subCommands, func(i, j int) bool {
		return c.subCommands[i].name < c.subCommands[j].name
	})
}

// Children 返回子命令列表
func (c *Command) Children() []*Command {
	c.fixChildren()
	return c.subCommands
}

// isCmd 检查名称是否匹配此命令（包括别名）
func (c *Command) isCmd(name string) bool {
	if c.name == name {
		return true
	}
	for _, alias := range c.aliases {
		if alias == name {
			return true
		}
	}
	return false
}

// findChildCmd 查找子命令
func (c *Command) findChildCmd(name string) *Command {
	c.fixChildren()
	return c.children[name]
}

// hasSubcommand 检查是否有子命令
func (c *Command) hasSubcommand() bool {
	c.fixChildren()
	return len(c.children) > 0
}

// Exec 执行命令
// arg: 已解析的参数值（由 ExecCommand 传入）
func (c *Command) Exec(ctx blocks.Context, arg any) {
	if c.execFunc != nil {
		c.execFunc(ctx, arg)
	}
}
