package promptx

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/aggronmagi/promptx/internal/debug"
)

// CommandContext command context 
// 
// The command execution function uses the ~Check...~ method provided by ~CommandContext~ to
// get the parameters.
// All Check methods, if there is no input from the user, interrupt the process through panic.
// Or if you enter a value that does not meet expectations, the process will be interrupted
// by panic.
// The call to the Check function should match the ~CommandParameter~ entered when
// creating the command.
// The Index parameter of the Check function starts at 0.
// After the ~CommandParameter~ detection passes, the player input retained internally
// is of string type.
// So ~CheckString~ is valid as long as the index position is entered.
// ~CheckSelectIndex~ is only used for the parameters corresponding to ~WithArgSelect~
// to get the first few options entered by the player (starting from 0)
type CommandContext interface {
	Context

	CheckString(index int) string
	CheckInteger(index int) int64
	CheckIPPort(index int) (ip, port string)
	CheckAddrs(index int) (addrs []string)
	CheckSelectIndex(index int) int
}

// CmdContext run context
type CmdContext struct {
	Context
	// command and subcommands
	Cmds []*Cmd
	// input args
	Args []string
	// input line
	Line string
	// Root root command. use for dynamic modify command.
	Root *Cmd
	// 当前命令
	Cur *Cmd
	// default select options
	SelectCC *SelectOptions
	// default input options
	InputCC *InputOptions
	// change flag,command args changed by auto action
	ChangeFlag bool
}

var _ CommandContext = (*CmdContext)(nil)

func (ctx *CmdContext) execCommand() {
	var err error
	for k, v := range ctx.Cur.args {
		err = v.Check(ctx, k)
		if err != nil {
			ctx.Println(err)
			return
		}
	}
	ctx.fixHisotry()
	defer func() {
		if p := recover(); p != nil {
			ctx.Printf("exec %s failed %v\n", ctx.Cur.name, p)
			// 不移除失败的历史记录
			//ctx.RemoveHistory(ctx.Line)
			return
		}
	}()
	ctx.Cur.execFunc(ctx)
}

func (ctx *CmdContext) CheckString(index int) string {
	if index >= len(ctx.Args) {
		panic(fmt.Sprintf("out of range,index:%d size:%d.", index, len(ctx.Args)))
	}
	return ctx.Args[index]
}

func (ctx *CmdContext) CheckInteger(index int) int64 {
	if index >= len(ctx.Args) {
		panic(fmt.Sprintf("out of range,index:%d size:%d.", index, len(ctx.Args)))
	}
	n, err := strconv.ParseInt(ctx.Args[index], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("convert int failed,%v. index:%d value:[%s]", err, index, ctx.Args[index]))
	}
	return n
}

func (ctx *CmdContext) CheckIPPort(index int) (ip, port string) {
	if index >= len(ctx.Args) {
		panic(fmt.Sprintf("out of range,index:%d size:%d.", index, len(ctx.Args)))
	}
	ip, port, err := net.SplitHostPort(ctx.Args[index])
	if err != nil {
		panic(fmt.Sprintf("split host port failed,%v. index:%d value:[%s]", err, index, ctx.Args[index]))
	}
	return
}

func (ctx *CmdContext) CheckAddrs(index int) (addrs []string) {
	if index >= len(ctx.Args) {
		panic(fmt.Sprintf("out of range,index:%d size:%d.", index, len(ctx.Args)))
	}

	for _, v := range strings.Split(ctx.Args[index], ",") {
		_, _, err := net.SplitHostPort(v)
		if err != nil {
			panic(fmt.Sprintf("split host port failed,%v. index:%d value:[%s]", err, index, ctx.Args[index]))
		}
		addrs = append(addrs, v)
	}
	return
}

func (ctx *CmdContext) CheckSelectIndex(index int) int {
	if index >= len(ctx.Args) {
		panic(fmt.Sprintf("out of range,index:%d size:%d.", index, len(ctx.Args)))
	}
	if index >= len(ctx.Cur.args) {
		panic(fmt.Sprintf("out of range,index:%d size:%d.not register.", index, len(ctx.Args)))
	}
	iface, ok := ctx.Cur.args[index].(interface {
		SelectOptions() []string
	})
	if !ok {
		panic(fmt.Sprintf("convert SelectOption failed. index:%d", index))
	}
	for k, v := range iface.SelectOptions() {
		if v == ctx.Args[index] {
			return k
		}
	}
	panic(fmt.Sprintf("index:%d not select options. %s not in %v", index, ctx.Args[index], iface.SelectOptions()))
}

func (ctx *CmdContext) fixHisotry() {
	if !ctx.ChangeFlag {
		debug.Log("no change flag")
		return
	}
	debug.Log("remove line:" + ctx.Line)
	ctx.RemoveHistory(ctx.Line)
	newLine := ctx.Root.FixCommandLine(ctx.Line, ctx.Args)
	debug.Log("add line:" + newLine)
	ctx.AddHistory(newLine)
}
