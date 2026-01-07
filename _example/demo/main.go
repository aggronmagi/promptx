package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aggronmagi/promptx"
	"github.com/aggronmagi/promptx/internal/debug"
)

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []*promptx.Suggest {
	return func(line string) []*promptx.Suggest {
		debug.Println("call dynamic completion", line)
		names := make([]*promptx.Suggest, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, &promptx.Suggest{
				Text:        f.Name(),
				Description: path,
			})
		}
		return names
	}
}

func colorCommand() *promptx.Cmd {
	return promptx.NewCommand("colorword", "show word color example").ExecFunc(func(ctx promptx.CommandContext) {
		ctx.WPrint(&promptx.AskWord, promptx.WordRed(" Red "),
			promptx.WordCyan(" Cyan "),
			promptx.WordBlue(" Blue "),
			promptx.WordGreen(" Green "),
			promptx.WordPurple(" Purple "),
			promptx.WordTurquoise(" Turquoise "),
			promptx.WordWhite(" White "),
			promptx.WordYellow(" Yellow "),
			&promptx.NewLineWord,
		)
	})
}

func loginCommand() *promptx.Cmd {
	return promptx.NewCommand("login", "登录游戏",
		promptx.WithArgSelect("选择登录的游戏服务器", []string{"开发服", "测试服", "体验服"}),
		promptx.WithArgsInput("账号:", promptx.CheckerNotEmpty()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		type LoginArgs struct {
			Server  string
			Account string
		}
		var args LoginArgs
		ctx.Bind(&args)
		ctx.Println("login success:", args.Account, "on", args.Server)
	})
}

func login2Command() *promptx.Cmd {
	return promptx.NewCommand("login2", "测试相似命令(test similar command)",
		promptx.WithArgSelect("选择登录的游戏服务器", []string{"开发服", "测试服", "体验服"}),
		promptx.WithArgsInput("账号:", promptx.CheckerNotEmpty()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		// 选择的登录服索引
		ctx.CheckSelectIndex(0)
		// 登录的服务字符串
		ctx.CheckString(0)
		// 输入的string
		ctx.CheckString(1)
	})
}

func editCommand() *promptx.Cmd {
	return promptx.NewCommand("edit", "show edit mode").ExecFunc(func(ctx promptx.CommandContext) {
		log.Println("current is emacs mode")
	}).SubCommands(
		promptx.NewCommand("vi", "修改vi模式").ExecFunc(func(ctx promptx.CommandContext) {
			log.Println("not support vim mode now!")
		}),
		promptx.NewCommand("emacs", "use emacs edit mode").ExecFunc(func(c promptx.CommandContext) {
			log.Println("use emacs mode now!")
		}),
	)
}

func sayCommand() *promptx.Cmd {
	return promptx.NewCommand("say", "say some words").SubCommands(
		promptx.NewCommand("hello", "say hello").ExecFunc(func(ctx promptx.CommandContext) {
			log.Println("hello!")
		}),
		promptx.NewCommand("bye", "say bye and exit").ExecFunc(func(ctx promptx.CommandContext) {
			ctx.Println("bye bye")
			ctx.Stop()
		}),
	)
}

//
// func dynamicTipCommand() *promptx.Cmd {
// 	return promptx.NewCommand("dynamic", "select files for operation").SubCommands(
// 		promptx.NewCommand("sub", "").DynamicTip(listFiles("./")).ExecFunc(func(ctx promptx.CommandContext) {
// 			//ctx.Println("select files", ctx.CheckString(0))
// 		}),
// 	).ExecFunc(func(ctx promptx.CommandContext) {
// 		ctx.Println("select files", ctx.CheckString(0))
// 	})
// }

func promptCommand() *promptx.Cmd {
	return promptx.NewCommand("setprompt", "set prompt string",
		promptx.WithArgSelect("color", []string{"red", "green"}),
		promptx.WithArgsInput("prompt:", promptx.CheckerNotEmpty()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		color := ctx.CheckSelectIndex(0)
		if color == 0 {
			ctx.SetPromptWords(&promptx.AskWord, promptx.WordRed(" "+ctx.CheckString(1)))
			return
		}
		ctx.SetPromptWords(&promptx.AskWord, promptx.WordGreen(" "+ctx.CheckString(1)))
	}).SubCommands(
		promptx.NewCommand("simple", "no color",
			promptx.WithArgsInput("prompt:", promptx.CheckerNotEmpty()),
		).ExecFunc(func(ctx promptx.CommandContext) {
			ctx.SetPrompt(ctx.CheckString(0))
		}),
	)
}
func lsCommand() *promptx.Cmd {
	return promptx.NewCommand("ls", "linux command ls").ExecFunc(func(ctx promptx.CommandContext) {
		ctx.ExitRawMode()
		cmd := exec.Command("ls", "-alh")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		ctx.EnterRawMode()
	})
}

func bashCommand() *promptx.Cmd {
	return promptx.NewCommand("bash", "enter linux bash").ExecFunc(func(ctx promptx.CommandContext) {
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
func sleepCommand() *promptx.Cmd {
	return promptx.NewCommand("sleep", "sleep some second",
		promptx.WithArgsInput("sleep second:", promptx.CheckerInteger()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		type SleepArgs struct {
			Seconds int
		}
		var args SleepArgs
		ctx.Bind(&args)
		time.Sleep(time.Second * time.Duration(args.Seconds))
	})
}

func asyncPrintCommand() *promptx.Cmd {
	return promptx.NewCommand("async-print", "async print message",
		promptx.WithArgsInput("print second:", promptx.CheckerNaturalNumber()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		c, cancel := context.WithTimeout(context.Background(), time.Duration(ctx.CheckInteger(0))*time.Second)
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

func panicCommand() *promptx.Cmd {
	return promptx.NewCommand("panic", "test panic command").ExecFunc(func(ctx promptx.CommandContext) {
		ctx.CheckInteger(1)
	})
}

func deepCommand() *promptx.Cmd {
	return promptx.NewCommand("deep1", "deep command").SubCommands(
		promptx.NewCommand("deep2-0", "deep2").SubCommands(
			promptx.NewCommand("deep3", "deep2-1-3"),
			promptx.NewCommand("deep3", "deep2-1-3"),
		),
		promptx.NewCommand("deep2-1", "deep2").SubCommands(
			promptx.NewCommand("deep3", "deep2-1-3"),
			promptx.NewCommand("deep3", "deep2-1-3"),
		),
	)
}

var cmds = []*promptx.Cmd{
	colorCommand(),
	loginCommand(),
	login2Command(),
	editCommand(),
	sayCommand(),
	promptCommand(),
	lsCommand(),
	bashCommand(),
	sleepCommand(),
	asyncPrintCommand(),
	panicCommand(),
	deepCommand(),
	// dynamicTipCommand(),
}

func main() {
	fmt.Printf("[%s]\n", strings.Trim(fmt.Sprint([]float32(nil)), "[]"))
	fmt.Printf("[%s]\n", strings.Trim(fmt.Sprint([]float32{1, 2}), "[]"))
	fmt.Printf("[%s]\n", fmt.Sprint(nil))
	tf := func(v any) string {
		in := fmt.Sprint(v)
		in = strings.TrimLeft(strings.TrimRight(in, "]"), "[")
		in = strings.TrimLeft(strings.TrimRight(in, ">"), "<")
		in = strings.Trim(strings.TrimSpace(in), "nil")
		return in
	}
	fmt.Printf("[%s]\n", tf(nil))
	fmt.Printf("[%s]\n", tf([]float32{}))
	fmt.Printf("[%s]\n", tf([]string(nil)))
	// new promptx
	p := promptx.New(promptx.WithCommon(promptx.WithCommonOptionCmds(cmds...)))
	// set log writer
	log.SetOutput(p.Stdout())

	p.ExecCommand([]string{"edit", "vi"})

	input, eof := promptx.Input(p, "input xx:", promptx.CheckerNotEmpty())
	fmt.Println(input, eof)
	sel := promptx.Select(p, "select xx:", []string{
		"x1",
		"x2",
		"exit",
	}, 1)

	if sel == 2 {
		return
	}

	p.Run()
}
