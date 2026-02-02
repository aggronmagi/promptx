package main

import (
	"errors"
	"log"
	"strings"

	"github.com/aggronmagi/promptx/v2"
	"github.com/aggronmagi/promptx/v2/v2/blocks"
)

// 0:未登录 1:area1 2:area2
var state int

const (
	SetCommon = ""
	SetArea1  = "area1"
	SetArea2  = "area2"
)

func loginCommand() *promptx.Command {
	type LoginArgs struct {
		Server  string `arg:"选择登录的游戏服务器" select:"开发服,测试服,体验服"`
		Account string `arg:"账号" check:"NotEmpty"`
	}
	return promptx.NewCommandWithFunc("login", "登录游戏", func(ctx blocks.Context, arg *LoginArgs) {
		ctx.Println(arg.Account, "login success on", arg.Server)
		state = 1
		// 切换命令组
		promptx.SwitchCommandGroup(ctx, SetArea1)
	})
}

func logoutCommand() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("logout", "退出游戏", func(ctx blocks.Context) {
		state = 0
		promptx.SwitchCommandGroup(ctx, SetCommon)
	})
}

func resetCommand() *promptx.Command {
	type ResetLogArgs struct {
		Level string `arg:"level" select:"debug,info"`
	}
	cmd := promptx.NewCommandWithFunc("reset", "reset options", func(ctx blocks.Context, arg *struct{}) {})
	cmd.SubCommands(
		promptx.NewCommandWithFunc("log", "reset log level", func(ctx blocks.Context, arg *ResetLogArgs) {
			log.Println("set log level", arg.Level)
		}),
	)
	return cmd
}

func playCommand() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("play", "play games", func(ctx blocks.Context) {
		ctx.Println("play game in area1...")
	})
}

func play2Command() *promptx.Command {
	return promptx.NewCommandWithFuncLegacy("play", "play games", func(ctx blocks.Context) {
		ctx.Println("play game in area2...")
	})
}

func gotoCommand() *promptx.Command {
	type GotoArgs struct {
		Area string `arg:"area" select:"area1,area2"`
	}
	return promptx.NewCommandWithFunc("goto", "goto area", func(ctx blocks.Context, arg *GotoArgs) {
		state = 1
		if arg.Area == SetArea2 {
			state = 2
		}
		promptx.SwitchCommandGroup(ctx, arg.Area)
	})
}

func main() {
	// 使用新的配置 API
	config := promptx.NewConfig()
	
	// 注意：OnNonCommand 和 CommandPrefix 现在在 CommandGroupConfig 中配置
	
	// 默认命令组
	config.DefaultCommandGroup().
		AddCommand(
			resetCommand(),
			loginCommand(),
		).
		PreCheck(func(ctx blocks.Context) error {
			if state != 0 {
				return errors.New("already login")
			}
			return nil
		}).
		CommandPrompt("not login >> ").
		OnNonCommand(func(ctx blocks.Context, command string) error {
			if strings.Contains(strings.ToLower(command), "err") {
				return errors.New("error command")
			}
			ctx.Println("non command", command)
			return nil
		})
	
	// area1 命令组
	config.AddCommandGroup(SetArea1).
		AddCommand(
			gotoCommand(),
			playCommand(),
			resetCommand(),
			logoutCommand(),
		).
		PreCheck(func(ctx blocks.Context) error {
			if state == 0 {
				return errors.New("not login")
			}
			if state != 1 {
				return errors.New("not in " + SetArea1)
			}
			return nil
		}).
		CommandPrompt("area1 >> ")
	
	// area2 命令组
	config.AddCommandGroup(SetArea2).
		AddCommand(
			gotoCommand(),
			play2Command(),
			resetCommand(),
			logoutCommand(),
		).
		PreCheck(func(ctx blocks.Context) error {
			if state == 0 {
				return errors.New("not login")
			}
			if state != 2 {
				return errors.New("not in " + SetArea2)
			}
			return nil
		}).
		CommandPrompt("area2 >> ")
	
	p := config.Build()
	log.SetOutput(p.Stdout())
	p.Run()
}
