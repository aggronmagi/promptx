package promptx

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/aggronmagi/promptx/internal/debug"
)

// CommandContext command context
//
// The command execution function uses the ~Check...~ method provided by ~CommandContext~ to
// get the parameters.
//
// It also supports strongly-typed argument binding via the ~Bind~ method.
//
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

	ArgSelect(index int, tip string, list []string, defaultSelect ...int) int
	ArgSelectString(index int, tip string, list []string, defaultSelect ...int) string
	ArgInput(index int, tip string, opts ...InputOption) (result string, eof error)
	GetArgs() []string

	// Bind bind positional arguments to struct fields.
	Bind(v interface{}) error
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

func (ctx *CmdContext) GetArgs() []string {
	return ctx.Args
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

func (ctx *CmdContext) ArgSelect(index int, tip string, list []string, defaultSelect ...int) int {
	if index > 0 && index < len(ctx.Args) {
		sel := ctx.Args[index]
		for k, v := range list {
			if v == sel {
				return k
			}
		}
	}
	id := Select(ctx, tip, list, defaultSelect...)
	if id < 0 {
		return id
	}
	// ctx.Args = append(ctx.Args, list[id])
	return id
}

func (ctx *CmdContext) ArgSelectString(index int, tip string, list []string, defaultSelect ...int) string {
	if index > 0 && index < len(ctx.Args) {
		sel := ctx.Args[index]
		for _, v := range list {
			if v == sel {
				return v
			}
		}
	}
	id := Select(ctx, tip, list, defaultSelect...)
	if id < 0 {
		return ""
	}
	return list[id]
}

func (ctx *CmdContext) ArgInput(index int, tip string, opts ...InputOption) (result string, eof error) {
	if index >= len(ctx.Args) {
		return ctx.RawInput(tip, opts...)
	}
	return ctx.Args[index], nil
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

func (ctx *CmdContext) Bind(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("Bind: must be a pointer to struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanSet() {
			continue
		}

		if i >= len(ctx.Args) {
			break
		}

		val := ctx.Args[i]
		switch field.Kind() {
		case reflect.String:
			field.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("Bind field %s: %v", rt.Field(i).Name, err)
			}
			field.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("Bind field %s: %v", rt.Field(i).Name, err)
			}
			field.SetUint(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("Bind field %s: %v", rt.Field(i).Name, err)
			}
			field.SetFloat(n)
		case reflect.Bool:
			n, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("Bind field %s: %v", rt.Field(i).Name, err)
			}
			field.SetBool(n)
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				field.Set(reflect.ValueOf(strings.Split(val, ",")))
			}
			// Add more slice types if needed
		}
	}

	return nil
}
