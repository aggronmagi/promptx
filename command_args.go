package promptx

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	completion "github.com/aggronmagi/promptx/completion"
	"github.com/aggronmagi/promptx/internal/debug"
)

// CommandParameter command args checker
type CommandParameter interface {
	Check(ctx *CmdContext, index int) (err error)
}

// InputChecker input checker
type InputChecker func(d *Document) (err error)

// InputNotEmpty Detection input must be a non-empty string
func InputNotEmpty() InputChecker {
	return func(d *Document) (err error) {
		if d.Text == "" {
			return errors.New("empty input")
		}
		return
	}
}

// InputInteger The detection input must be a non-zero value
func InputInteger() InputChecker {
	return func(d *Document) error {
		if d.Text == "" {
			return errors.New("empty input")
		}
		v, err := strconv.ParseInt(d.Text, 10, 64)
		if err != nil {
			return err
		}

		if v == 0 {
			return errors.New("zero value invalid")
		}

		return nil
	}
}

// InputNotZeroInteger The detection input must be a non-zero number
func InputNotZeroInteger() InputChecker {
	return func(d *Document) error {
		if d.Text == "" {
			return errors.New("empty input")
		}
		v, err := strconv.ParseInt(d.Text, 10, 64)
		if err != nil {
			return err
		}

		if v == 0 {
			return errors.New("zero value invalid")
		}

		return nil
	}
}

// InputNaturalNumber Detection input must be a natural number (1,2,3....)
func InputNaturalNumber() InputChecker {
	return func(d *Document) (err error) {
		if d.Text == "" {
			return errors.New("empty input")
		}
		v, err := strconv.ParseInt(d.Text, 10, 64)
		if err != nil {
			return err
		}

		if v < 1 {
			return fmt.Errorf("%d is not natural number", v)
		}

		return nil
	}
}

// InputIntegerRange The detection input must be a number in the min, max interval
func InputIntegerRange(min, max int64) InputChecker {
	if min > max {
		min, max = max, min
	}
	return func(d *Document) (err error) {
		if d.Text == "" {
			return errors.New("empty input")
		}
		var v int64
		v, err = strconv.ParseInt(d.Text, 10, 64)
		if err != nil {
			return err
		}

		if v < min || v > max {
			return fmt.Errorf("limit range(%d,%d) - %d", min, max, v)
		}
		return nil
	}
}

// InputIPPort The detection input must be an IP port(example 127.0.0.1:8080)
func InputIPPort() InputChecker {
	return func(d *Document) (err error) {
		_, _, err = net.SplitHostPort(d.Text)
		return
	}
}

// InputAddrsArray The detection input must be an array of IP ports, separated by ",".
func InputAddrsArray() InputChecker {
	return func(d *Document) (err error) {
		for _, v := range strings.Split(d.Text, ",") {
			_, _, err = net.SplitHostPort(v)
			if err != nil {
				return err
			}
		}
		return
	}
}

type inputArgsChecker struct {
	tip     string
	opts    []InputOption
	checker InputChecker
}

func (c *inputArgsChecker) Check(ctx *CmdContext, index int) (err error) {
	if index < len(ctx.Args) {
		err = c.checker(&Document{Text: ctx.Args[index]})
		if err == nil {
			// 检测通过,打印参数
			cc := *ctx.InputCC
			cc.ApplyOption(c.opts...)
			cc.SetOption(WithInputOptionPrefixText(c.tip))
			ctx.WPrintln(cc.ResultText(&cc, 1, &Document{Text: ctx.Args[index]})...)
			return nil
		}
		ctx.Printf("input %s [%s] %v", c.tip, ctx.Args[index], err)
	} else {
		ctx.Args = append(ctx.Args, "")
	}
	ctx.ChangeFlag = true
	debug.Log("input change flag")
	// 重新输入
	ret, eof := ctx.Input(c.tip, c.opts...)
	if eof != nil {
		return errors.New("user cancel")
	}
	// 替换输入
	if ctx.Args[index] != "" {
		ctx.Line = strings.Replace(ctx.Line, ctx.Args[index], ret, 1)
	}
	ctx.Args[index] = ret
	return nil
}

// WithArgsInput Manually enter a string as a parameter
func WithArgsInput(tip string, check InputChecker, opts ...InputOption) CommandParameter {
	opts = append(opts, WithInputOptionValidFunc(check))
	return &inputArgsChecker{
		tip:     tip,
		opts:    opts,
		checker: check,
	}
}

type selectArgsChecker struct {
	tip  string
	list []string
	opts []SelectOption
}

func (c *selectArgsChecker) Check(ctx *CmdContext, index int) (err error) {
	if index < len(ctx.Args) {
		arg := ctx.Args[index]
		// 检测是否手动输入的索引
		var ok bool
		i, err := strconv.ParseInt(arg, 10, 32)
		if err == nil && int(i) < len(c.list) {
			arg = c.list[i]
			ok = true
		}
		// 输入的是选项
		if !ok {
			for k, v := range c.list {
				if v == arg {
					ok = true
					i = int64(k)
					break
				}
			}
		}
		// 输入的合法
		if ok {
			cc := *ctx.SelectCC
			cc.ApplyOption(c.opts...)
			cc.SetOption(WithSelectOptionMulti(false))
			cc.SetOption(WithSelectOptionTipText(c.tip))
			cc.Options = make([]*completion.Suggest, 0, len(c.list))
			for _, v := range c.list {
				cc.Options = append(cc.Options, &Suggest{
					Text: v,
				})
			}
			ctx.WPrintln(cc.FinishText(&cc, []int{int(i)})...)
			return nil
		}
	} else {
		ctx.Args = append(ctx.Args, "")
	}
	ctx.ChangeFlag = true
	debug.Log("select change flag")
	// 重新输入
	sel := ctx.Select(c.tip, c.list, c.opts...)
	if sel < 0 {
		return errors.New("user cancel")
	}
	// 替换输入
	ctx.Args[index] = c.list[sel]
	return nil
}

func (c *selectArgsChecker) SelectOptions() []string {
	return c.list
}

// WithArgSelect Single choice parameter Choose one of many
func WithArgSelect(tip string, list []string, opts ...SelectOption) CommandParameter {
	if len(list) < 1 {
		panic("select option size is zero")
	}
	// 选项中不能有空格,否则替换历史记录会导致参数个数异常
	for k, v := range list {
		list[k] = strings.Replace(v, " ", "-", -1)
	}
	return &selectArgsChecker{
		tip:  tip,
		list: list,
		opts: opts,
	}
}
