package promptx

import (
	"fmt"
	"io"

	input "github.com/aggronmagi/promptx/input"
	"github.com/aggronmagi/promptx/internal/debug"
	output "github.com/aggronmagi/promptx/output"
	"github.com/aggronmagi/promptx/terminal"
)

// PromptOptionsOptionDeclareWithDefault promptx options
// generate by https://github.com/timestee/optiongen
//go:generate optionGen --option_with_struct_name=false --v=true
func PromptOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// default global input options
		"InputOptions": []InputOption(nil),
		// default global select options
		"SelectOptions": []SelectOption(nil),
		// default common options. use to create default optins
		"CommonOpions": []CommonOption(nil),
		// default manager. if it is not nil, ignore CommonOpions.
		"BlocksManager": BlocksManager(nil),
		// input
		"Input": input.ConsoleParser(input.NewStandardInputParser()),
		// output
		"Output": output.ConsoleWriter(output.NewStandardOutputWriter()),
		"Stderr": output.ConsoleWriter(output.NewStderrWriter()),
	}
}

// Promptx prompt command line
type Promptx struct {
	// config options
	cc *PromptOptions

	// default select options
	selectCC *SelectOptions
	// default input options
	inputCC *InputOptions
	//
	console *terminal.TerminalApp
}

// NewPromptx new prompt
func NewPromptx(opts ...PromptOption) *Promptx {
	cc := NewPromptOptions(opts...)
	return newPromptx(cc)
}

// NewCommandPromptx new with command
func NewCommandPromptx(cmds ...*Cmd) *Promptx {
	return NewPromptx(
		WithCommonOpions(
			WithCommonOptionCommands(cmds...),
		),
	)
}

// NewOptionCommandPromptx new with command and options
// use for replace NewCommandPromptx when you need apply other options.
// example: NewCommandPromptx(cmds...) => NewOptionCommandPromptx(NewPromptOptions(....),cmds...)
func NewOptionCommandPromptx(cc *PromptOptions, cmds ...*Cmd) *Promptx {
	if cc == nil {
		cc = NewPromptOptions()
	}
	cc.CommonOpions = append(cc.CommonOpions, WithCommonOptionCommands(cmds...))
	return newPromptx(cc)
}

// newPromptx new prompt
func newPromptx(cc *PromptOptions) *Promptx {
	p := new(Promptx)
	p.selectCC = NewSelectOptions(cc.SelectOptions...)
	p.inputCC = NewInputOptions(cc.InputOptions...)
	if cc.BlocksManager == nil {
		cc.BlocksManager = NewDefaultBlockManger(cc.CommonOpions...)
	}
	p.console = terminal.NewTerminalApp(cc.Input)
	cc.BlocksManager.SetWriter(cc.Output)
	cc.BlocksManager.SetExecContext(p)
	cc.BlocksManager.UpdateWinSize(cc.Input.GetWinSize())
	p.cc = cc

	return p
}

// // Start start run async
// func (p *Promptx) Start() (err error) {
// 	// go p.console.Run(p.cc.BlocksManager)
// 	return
// }

// Run run prompt
func (p *Promptx) Run() (err error) {
	p.console.Run(p.cc.BlocksManager)
	debug.Println("exit root run")
	return
}

func (p *Promptx) Stop() {
	p.console.Stop()
}

// EnterRawMode enter raw mode for read key press real time
func (p *Promptx) EnterRawMode() (err error) {
	return p.console.EnterRaw()
}

// ExitRawMode exit raw mode
func (p *Promptx) ExitRawMode() (err error) {
	return p.console.ExitRaw()
}

// Stdout return a wrap stdout writer. it can refersh view correct
func (p *Promptx) Stdout() io.Writer {
	return terminal.NewWrapWriter(p.cc.Output, p.console)
}

// Stderr std err
func (p *Promptx) Stderr() io.Writer {
	return terminal.NewWrapWriter(p.cc.Stderr, p.console)
}

// ClearScreen clears the screen.
func (p *Promptx) ClearScreen() {
	out := p.cc.Output
	out.EraseScreen()
	out.CursorGoTo(0, 0)
	debug.AssertNoError(out.Flush())
}

// SetTitle set title
func (p *Promptx) SetTitle(title string) {
	if len(title) < 1 {
		return
	}
	out := p.cc.Output
	out.SetTitle(title)
	debug.AssertNoError(out.Flush())
}

// ClearTitle clear title
func (p *Promptx) ClearTitle() {
	out := p.cc.Output
	out.ClearTitle()
	debug.AssertNoError(out.Flush())
}

// SetPrompt update prompt.
func (p *Promptx) SetPrompt(prompt string) {
	if iface, ok := p.cc.BlocksManager.(interface {
		SetPrompt(prompt string)
	}); ok {
		iface.SetPrompt(prompt)
		// p.syncCh <- struct{}{}
	}
}

// SetPromptWords update prompt string. custom display.
func (p *Promptx) SetPromptWords(words ...*Word) {
	if iface, ok := p.cc.BlocksManager.(interface {
		SetPromptWords(words ...*Word)
	}); ok {
		iface.SetPromptWords(words...)
	}
}

func (p *Promptx) ExecCommand(args []string) {
	if iface, ok := p.cc.BlocksManager.(interface {
		ExecCommand(args []string)
	}); ok {
		iface.ExecCommand(args)
	}
}

// Print = fmt.Print
func (p *Promptx) Print(v ...interface{}) {
	fmt.Fprint(p.Stdout(), v...)
}

// Printf = fmt.Printf
func (p *Promptx) Printf(format string, v ...interface{}) {
	fmt.Fprintf(p.Stdout(), format, v...)
}

// Println = fmt.Println
func (p *Promptx) Println(v ...interface{}) {
	fmt.Fprintln(p.Stdout(), v...)
}

// WPrint  print words
func (p *Promptx) WPrint(words ...*Word) {
	out := p.cc.Output
	for _, v := range words {
		out.SetColor(v.TextColor, v.BGColor, v.Bold)
		out.WriteStr(v.Text)
	}
	out.SetColor(DefaultColor, DefaultColor, false)
}

// WPrintln print words and newline
func (p *Promptx) WPrintln(words ...*Word) {
	out := p.cc.Output
	for _, v := range words {
		out.SetColor(v.TextColor, v.BGColor, v.Bold)
		out.WriteStr(v.Text)
	}
	out.SetColor(DefaultColor, DefaultColor, false)
	out.WriteRawStr("\n")
}

// Input get input
func (p *Promptx) Input(tip string, opts ...InputOption) (result string, err error) {
	// copy a new config
	newCC := (*p.inputCC)

	// apply input opts
	newCC.ApplyOption(opts...)
	// set internal options
	newCC.SetOption(WithInputOptionFinishFunc(func(input string, eof error) {
		result, err = input, eof
	}))
	if tip != "" {
		newCC.SetOption(WithInputOptionPrefixText(tip))
	}
	//
	input := NewInputManager(&newCC)

	input.SetExecContext(p)
	input.SetWriter(p.cc.Output)
	input.UpdateWinSize(p.cc.Input.GetWinSize())
	debug.Println("enter input")
	p.console.Run(input)
	debug.Println("exit input")
	return
}

// Select get input
func (p *Promptx) Select(tip string, list []string, opts ...SelectOption) (result int) {
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	// reset options internal
	newCC.SetOption(WithSelectOptionFinishFunc(
		func(sels []int) {
			if len(sels) < 1 {
				result = -1
				return
			}
			result = sels[0]
		},
	))
	newCC.SetOption(WithSelectOptionMulti(false))
	newCC.SetOption(WithSelectOptionTipText(tip))

	// if opts set options. use opts values
	if len(newCC.Options) < 1 {
		if len(list) < 1 {
			return -1
		}
		for _, v := range list {
			newCC.Options = append(newCC.Options, &Suggest{
				Text: v,
			})
		}
	}
	sel := NewSelectManager(&newCC)

	sel.SetExecContext(p)
	sel.SetWriter(p.cc.Output)
	sel.UpdateWinSize(p.cc.Input.GetWinSize())
	p.console.Run(sel)
	return
}

// Select get input
func (p *Promptx) MulSel(tip string, list []string, opts ...SelectOption) (result []int) {
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	newCC.SetOption(WithSelectOptionFinishFunc(
		func(sels []int) {
			result = sels
		},
	))
	newCC.SetOption(WithSelectOptionMulti(true))

	// if opts set options. use opts values
	if len(newCC.Options) < 1 {
		if len(list) < 1 {
			return
		}
		for _, v := range list {
			newCC.Options = append(newCC.Options, &Suggest{
				Text: v,
			})
		}
	}
	sel := NewSelectManager(&newCC)
	sel.SetExecContext(p)
	sel.SetWriter(p.cc.Output)
	sel.UpdateWinSize(p.cc.Input.GetWinSize())
	p.console.Run(sel)
	return
}

// // run internal
// func (p *Promptx) run() (err error) {

// 	defer func() {
// 		p.start.Store(false)
// 		p.wg.Done()
// 	}()
// 	if p.t == nil {
// 		p.t = NewTerminal(p.cc.Stdin, p.cc.Stdout, p.cc.ChanSize)
// 		p.mgr.SetWriter(output.NewConsoleWriter(p.t))
// 	}
// 	// update windows size
// 	w, h, err := term.GetSize(syscall.Stdout)
// 	if err != nil {
// 		return err
// 	}
// 	p.mgr.UpdateWinSize(w, h)
// 	// render pre
// 	p.mgr.Render(NormalStatus)
// 	p.mgr.SetExecFunc(p.Exec)
// 	p.mgr.SetExecContext(p)

// 	p.stop = make(chan struct{})
// 	p.t.Start()
// 	defer func() {
// 		close(p.stop)
// 		p.t.Close()
// 		p.stop = nil
// 		p.t = nil
// 	}()

// 	exitCh := make(chan int)
// 	winSize := make(chan *WinSize)
// 	p.execCh = make(chan func())
// 	go HandleSignals(exitCh, winSize, p.stop)

// 	go func() {
// 		for {
// 			select {
// 			case <-p.stop:
// 				return
// 			case f := <-p.execCh:
// 				go func() {
// 					f()
// 					if p.exchange.Load() {
// 						p.refreshCh <- struct{}{}
// 					}
// 					p.cond.Signal()
// 				}()
// 			}
// 		}
// 	}()

// 	// event chan
// 	for {
// 		select {
// 		case in, ok := <-p.t.InputChan():
// 			if !ok {
// 				return
// 			}
// 			key := input.GetKey(in)

// 			if p.getCurrent().Event(key, in) {
// 				if p.exchangeNext(true) {
// 					debug.Println("recv exit. but change next screen")
// 					break
// 				}
// 				debug.Println("recv exit and not change next screen")
// 				return
// 			}
// 			p.exchangeNext(true)
// 			p.exchange.Store(false)
// 		case <-p.refreshCh:
// 			p.getCurrent().Render(NormalStatus)
// 		case <-p.syncCh:
// 			p.exchangeNext(true)
// 			p.exchange.Store(false)
// 		case size := <-winSize:
// 			p.getCurrent().UpdateWinSize(size.Col, size.Row)
// 			p.getCurrent().Render(NormalStatus)
// 		case code := <-exitCh:
// 			p.getCurrent().Render(CancelStatus)
// 			fmt.Println("exit code", code)
// 			// os.Exit(code)
// 			return
// 		}
// 	}

// 	return
// }
