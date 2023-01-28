package main

import (
	"errors"
	"log"

	"github.com/aggronmagi/promptx"
)

// 0:未登录 1:area1 2:area2
var state int

const (
	SetCommon = ""
	SetArea1  = "area1"
	SetArea2  = "area2"
)

func loginCommand() *promptx.Cmd {
	return promptx.NewCommand("login", "登录游戏",
		promptx.WithArgSelect("选择登录的游戏服务器", []string{"开发服", "测试服", "体验服"}),
		promptx.WithArgsInput("账号:", promptx.InputNotEmpty()),
	).ExecFunc(func(ctx promptx.CommandContext) {
		// 登录的服务字符串
		ctx.CheckString(0)
		// 输入的string
		id := ctx.CheckString(1)
		ctx.Println(id, "login success")
		state = 1
		//
		ctx.SwitchCommandSet(SetArea1)
	})
}

func logoutCommand() *promptx.Cmd {
	return promptx.NewCommand("logout", "退出游戏").ExecFunc(func(ctx promptx.CommandContext) {
		state = 0
		ctx.SwitchCommandSet(SetCommon)
	})
}

func resetCommand() *promptx.Cmd {
	return promptx.NewCommand("reset", "reset options").SubCommands(
		promptx.NewCommand("log", "reset log level",
			promptx.WithArgSelect("level", []string{"debug", "info"}),
		).ExecFunc(func(ctx promptx.CommandContext) {
			log.Println("set log level", ctx.CheckString(0))
		}),
	)
}

func playCommand() *promptx.Cmd {
	return promptx.NewCommand("play", "play games").ExecFunc(func(ctx promptx.CommandContext) {
		ctx.Println("play game in area1...")
	})
}

func play2Command() *promptx.Cmd {
	return promptx.NewCommand("play", "play games").ExecFunc(func(ctx promptx.CommandContext) {
		ctx.Println("play game in area2...")
	})
}
func gotoCommand() *promptx.Cmd {
	return promptx.NewCommand("goto", "goto area",
		promptx.WithArgSelect("area", []string{SetArea1, SetArea2}),
	).ExecFunc(func(ctx promptx.CommandContext) {
		area := ctx.CheckString(0)
		state = 1 + ctx.CheckSelectIndex(0)
		ctx.SwitchCommandSet(area)
	})
}

func main() {
	p := promptx.NewPromptx()
	// default set
	p.AddCommandSet(SetCommon, []*promptx.Cmd{
		resetCommand(),
		loginCommand(),
	}, promptx.WithCommandSetOptionPreCheck(func(ctx promptx.Context) error {
		if state != 0 {
			return errors.New("already login")
		}
		return nil
	}), promptx.WithCommandSetOptionPrompt("not login >> "))
	//
	p.AddCommandSet(SetArea1, []*promptx.Cmd{
		gotoCommand(),
		playCommand(),
		resetCommand(),
		logoutCommand(),
	}, promptx.WithCommandSetOptionPreCheck(func(ctx promptx.Context) error {
		if state == 0 {
			return errors.New("not login")
		}
		if state != 1 {
			return errors.New("not in " + SetArea1)
		}
		return nil
	}), promptx.WithCommandSetOptionPrompt("area1 >> "))

	p.AddCommandSet(SetArea2, []*promptx.Cmd{
		gotoCommand(),
		play2Command(),
		resetCommand(),
		logoutCommand(),
	}, promptx.WithCommandSetOptionPreCheck(func(ctx promptx.Context) error {
		if state == 0 {
			return errors.New("not login")
		}
		if state != 2 {
			return errors.New("not in " + SetArea2)
		}
		return nil
	}), promptx.WithCommandSetOptionPrompt("area2 >> "))
	log.SetOutput(p.Stdout())
	p.Run()
}
