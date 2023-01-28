package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
		promptx.WithArgsInput("账号:", promptx.InputNotEmpty()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		// 选择的登录服索引
		ctx.CheckSelectIndex(0)
		// 登录的服务字符串
		ctx.CheckString(0)
		// 输入的string
		ctx.CheckString(1)
		ctx.Println("login success")
	})
}

func login2Command() *promptx.Cmd {
	return promptx.NewCommand("login2", "测试相似命令(test similar command)",
		promptx.WithArgSelect("选择登录的游戏服务器", []string{"开发服", "测试服", "体验服"}),
		promptx.WithArgsInput("账号:", promptx.InputNotEmpty()),
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
		promptx.WithArgsInput("prompt:", promptx.InputNotEmpty()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		color := ctx.CheckSelectIndex(0)
		if color == 0 {
			ctx.SetPromptWords(&promptx.AskWord, promptx.WordRed(" "+ctx.CheckString(1)))
			return
		}
		ctx.SetPromptWords(&promptx.AskWord, promptx.WordGreen(" "+ctx.CheckString(1)))
	}).SubCommands(
		promptx.NewCommand("simple", "no color",
			promptx.WithArgsInput("prompt:", promptx.InputNotEmpty()),
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
		promptx.WithArgsInput("sleep second:", promptx.InputInteger()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		time.Sleep(time.Second * time.Duration(ctx.CheckInteger(0)))
	})
}

func asyncPrintCommand() *promptx.Cmd {
	return promptx.NewCommand("async-print", "async print message",
		promptx.WithArgsInput("print second:", promptx.InputNaturalNumber()),
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
	// dynamicTipCommand(),
}

func main() {
	// new promptx
	p := promptx.NewCommandPromptx(cmds...)
	// set log writer
	log.SetOutput(p.Stdout())

	p.ExecCommand([]string{"edit", "vi"})

	p.Input("input xx:")
	p.Select("select xx:", []string{
		"x1",
		"x2",
		"x3",
	})

	p.Run()
}
