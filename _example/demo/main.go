package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aggronmagi/promptx/v2/v2"
	"github.com/aggronmagi/promptx/v2/v2/blocks"
)

// ============================================================================
// 使用 Commander 接口的命令示例
// ============================================================================

// ColorCommand 展示彩色文本的命令
type ColorCommand struct{}

func (c *ColorCommand) Name() string {
	return "colorword"
}

func (c *ColorCommand) Help() string {
	return "show word color example"
}

func (c *ColorCommand) Exec(ctx blocks.Context) {
	ctx.WPrint(blocks.AskWord, blocks.WordRed(" Red "),
		blocks.WordCyan(" Cyan "),
		blocks.WordBlue(" Blue "),
		blocks.WordGreen(" Green "),
		blocks.WordPurple(" Purple "),
		blocks.WordTurquoise(" Turquoise "),
		blocks.WordWhite(" White "),
		blocks.WordYellow(" Yellow "),
		blocks.NewLineWord,
	)
}

// LoginCommand 登录命令示例
type LoginCommand struct {
	Server  string `arg:"选择登录的游戏服务器" select:"开发服,测试服,体验服"`
	Account string `arg:"账号" check:"NotEmpty"`
}

func (c *LoginCommand) Name() string {
	return "login"
}

func (c *LoginCommand) Help() string {
	return "登录游戏"
}

func (c *LoginCommand) Exec(ctx blocks.Context) {
	ctx.Println(c.Account, "login success on", c.Server)
}

// EditCommand 编辑模式命令（带子命令）
type EditCommand struct{}

func (c *EditCommand) Name() string {
	return "edit"
}

func (c *EditCommand) Help() string {
	return "show edit mode"
}

func (c *EditCommand) Exec(ctx blocks.Context) {
	log.Println("current is emacs mode")
}

// SayCommand 说话命令（带子命令）
type SayCommand struct{}

func (c *SayCommand) Name() string {
	return "say"
}

func (c *SayCommand) Help() string {
	return "say some words"
}

func (c *SayCommand) Exec(ctx blocks.Context) {
	// 不做任何事，只有子命令执行
}

// SayHelloCommand 说 hello
type SayHelloCommand struct{}

func (c *SayHelloCommand) Name() string {
	return "hello"
}

func (c *SayHelloCommand) Help() string {
	return "say hello"
}

func (c *SayHelloCommand) Exec(ctx blocks.Context) {
	log.Println("hello!")
}

// SayByeCommand 说再见并退出
type SayByeCommand struct{}

func (c *SayByeCommand) Name() string {
	return "bye"
}

func (c *SayByeCommand) Help() string {
	return "say bye and exit"
}

func (c *SayByeCommand) Exec(ctx blocks.Context) {
	ctx.Println("bye bye")
	ctx.Stop()
}

// SetPromptCommand 设置提示符命令（带子命令）
type SetPromptCommand struct {
	Color  string `arg:"color" select:"red,green"`
	Prompt string `arg:"prompt:" check:"NotEmpty"`
}

func (c *SetPromptCommand) Name() string {
	return "setprompt"
}

func (c *SetPromptCommand) Help() string {
	return "set prompt string"
}

func (c *SetPromptCommand) Exec(ctx blocks.Context) {
	if c.Color == "red" {
		ctx.SetPromptWords(blocks.AskWord, blocks.WordRed(" "+c.Prompt))
		return
	}
	ctx.SetPromptWords(blocks.AskWord, blocks.WordGreen(" "+c.Prompt))
}

// SetPromptSimpleCommand 不用颜色设置提示符
type SetPromptSimpleCommand struct {
	Prompt string `arg:"prompt:" check:"NotEmpty"`
}

func (c *SetPromptSimpleCommand) Name() string {
	return "simple"
}

func (c *SetPromptSimpleCommand) Help() string {
	return "no color"
}

func (c *SetPromptSimpleCommand) Exec(ctx blocks.Context) {
	ctx.SetPrompt(c.Prompt)
}

// ============================================================================
// 使用函数风格的命令示例（对于简单命令）
// ============================================================================

func lsCommand() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("ls", "linux command ls", func(ctx promptx.Context) {
		ctx.ExitRawMode()
		cmd := exec.Command("ls", "-alh")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		ctx.EnterRawMode()
	})
}

func bashCommand() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("bash", "enter linux bash", func(ctx promptx.Context) {
		err := ctx.ExitRawMode()
		ctx.Println("exit rawmode:", err)
		cmd := exec.Command("bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		ctx.Println("run result:", err)
		ctx.EnterRawMode()
	})
}

func sleepCommand() *promptx.Command {
	type SleepArgs struct {
		Seconds int64 `arg:"sleep second" check:"Integer"`
	}
	return promptx.NewCommandWithFunc("sleep", "sleep some second", func(ctx blocks.Context, arg *SleepArgs) {
		time.Sleep(time.Second * time.Duration(arg.Seconds))
	})
}

func asyncPrintCommand() *promptx.Command {
	type AsyncPrintArgs struct {
		Seconds int64 `arg:"print second" check:"NaturalNumber"`
	}
	return promptx.NewCommandWithFunc("async-print", "async print message", func(ctx blocks.Context, arg *AsyncPrintArgs) {
		c, cancel := context.WithTimeout(context.Background(), time.Duration(arg.Seconds)*time.Second)
		go func() {
			defer cancel()
			ticker := time.NewTicker(time.Millisecond * 500)
			for {
				select {
				case <-ticker.C:
					fmt.Fprintln(ctx.Stdout(), "async message "+time.Now().String())
				case <-c.Done():
					return
				}
			}
		}()
		ctx.Println("async message command")
	})
}

func panicCommand() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("panic", "test panic command", func(ctx promptx.Context) {
		panic("test panic")
	})
}

func deepCommand() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("deep1", "deep command", func(ctx promptx.Context) {}).SubCommands(
		promptx.NewCommandWithFuncLegacy("deep2-0", "deep2", func(ctx promptx.Context) {}).SubCommands(
			promptx.NewCommandWithFuncLegacy("deep3", "deep2-1-3", func(ctx promptx.Context) {}),
		),
		promptx.NewCommandWithFuncLegacy("deep2-1", "deep2", func(ctx promptx.Context) {}).SubCommands(
			promptx.NewCommandWithFuncLegacy("deep3", "deep2-1-3", func(ctx promptx.Context) {}),
		),
	)
}

func main() {
	// 使用新的配置 API
	config := promptx.NewConfig()

	// 添加 Commander 格式的命令
	config.DefaultCommandGroup().
		AddCommand(
			// Commander 格式的命令
			promptx.NewCommand(&ColorCommand{}),
			promptx.NewCommand(&LoginCommand{}),
			promptx.NewCommand(&EditCommand{}).SubCommands(
				promptx.NewCommandWithFuncLegacy("vi", "修改vi模式", func(ctx promptx.Context) {
					log.Println("not support vim mode now!")
				}),
				promptx.NewCommandWithFuncLegacy("emacs", "use emacs edit mode", func(c promptx.Context) {
					log.Println("use emacs mode now!")
				}),
			),
			promptx.NewCommand(&SayCommand{}).SubCommands(
				promptx.NewCommand(&SayHelloCommand{}),
				promptx.NewCommand(&SayByeCommand{}),
			),
			promptx.NewCommand(&SetPromptCommand{}).SubCommands(
				promptx.NewCommand(&SetPromptSimpleCommand{}),
			),
			// 函数风格的命令
			lsCommand(),
			bashCommand(),
			sleepCommand(),
			asyncPrintCommand(),
			panicCommand(),
			deepCommand(),
		).
		CommandPrefix("!").
		OnNonCommand(func(ctx blocks.Context, command string) error {
			ctx.Println("non command", command)
			return nil
		})

	p := config.Build()

	// set log writer
	log.SetOutput(p.Stdout())

	p.Run()
}
