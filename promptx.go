package promptx

import (
	"fmt"
	"io"

	input "github.com/aggronmagi/promptx/input"
	"github.com/aggronmagi/promptx/internal/debug"
	output "github.com/aggronmagi/promptx/output"
	"github.com/aggronmagi/promptx/terminal"
)

// PromptOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n PromptOption -o gen_options_prompt.go
func promptxPromptOptions() interface{} {
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

// CommandSetOptions command set options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n CommandSetOption -f -o gen_options_commandset.go
func promptxCommandSetOptions() interface{} {
	return map[string]interface{}{
		// comman set name
		"name": "",
		// command set commands
		"commands": []*Cmd{},
		// Set the operation record of the current command set to persist the saved file. If not set, the history operation record will be cleared every time the command set is switched.
		"HistoryFile": "",
		// Set the pre-detection function for all commands executed in the current command set
		"PreCheck": (func(ctx Context) error)(nil),
		// Set the prompt string when switching to the command set(custom word color)
		"PromptWords": []*Word{},
		// Set the notification function when switching to the command set, args is the parameter passed by SwitchCommandSet.
		"ChangeNotify": func(ctx Context, args []interface{}) {},
	}
}

// WithCommandSetOptionPrompt Set the prompt string when switching to the command set(default color)
func WithCommandSetOptionPrompt(prompt string) CommandSetOption {
	return func(cc *CommandSetOptions) CommandSetOption {
		previous := cc.PromptWords
		cc.PromptWords = []*Word{WordDefault(prompt)}
		return WithCommandSetOptionPromptWords(previous...)
	}
}

func init() {
	InstallCommandSetOptionsWatchDog(func(cc *CommandSetOptions) {
		if len(cc.PromptWords) < 1 {
			cc.PromptWords = append(cc.PromptWords, WordDefault(">> "))
		}
	})
}

// Context Run Command Context
type Context interface {
	// Input get input string. if cancel return error io.EOF
	Input(tip string, opts ...InputOption) (result string, eof error)
	// Select select one from list. if cancel,return -1
	Select(tip string, list []string, opts ...SelectOption) (sel int)
	// MulSel like Select, but can choose list more then one. if cancel, return empty slice
	MulSel(tip string, list []string, opts ...SelectOption) (sel []int)
	// Stop stop run
	Stop()
	// EnterRawMode enter raw mode for read key press real time
	EnterRawMode() (err error)
	// ExitRawMode exit raw mode
	ExitRawMode() (err error)
	// Stdout return a wrap stdout writer. it can refersh view correct
	Stdout() io.Writer
	// Stderr std err
	Stderr() io.Writer
	// ClearScreen clears the screen.
	ClearScreen()
	// SetTitle set window title
	SetTitle(title string)
	// ClearTitle clear window title
	ClearTitle()
	// SetPrompt update prompt string. prompt will auto add space suffix.
	SetPrompt(prompt string)
	// SetPromptWords update prompt string. custom display.
	SetPromptWords(words ...*Word)
	// ResetCommands reset all command set.
	ResetCommands(commands ...*Cmd)
	// RemoveHistory remove from history
	RemoveHistory(line string)
	// AddHistory add line to history
	AddHistory(line string)
	// reset history file
	ResetHistoryFile(filename string)
	// SetCommandPreCheck check before exec *promptx.Cmd
	SetCommandPreCheck(f func(ctx Context) error)
	// AddCommandSet add command set,it will auto switch when first add commandset.
	AddCommandSet(name string, cmds []*Cmd, opts ...CommandSetOption)
	// SwitchCommandSet switch to specify commands set,args will pass to ChangeNotify func.
	SwitchCommandSet(name string, args ...interface{})
	// Print = fmt.Print
	Print(v ...interface{})
	// Printf = fmt.Printf
	Printf(fmt string, v ...interface{})
	// Println = fmt.Println
	Println(v ...interface{})

	// WPrint  print words
	WPrint(words ...*Word)
	// WPrintln print words and newline
	WPrintln(words ...*Word)
}

type PromptMain interface {
	Context
	Run() error
	ExecCommand(args []string)
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
	//
	sets map[string]*CommandSetOptions
}

var _ Context = &Promptx{}

// NewPromptx new prompt
func NewPromptx(opts ...PromptOption) PromptMain {
	cc := NewPromptOptions(opts...)
	return newPromptx(cc)
}

// NewCommandPromptx new with command
func NewCommandPromptx(cmds ...*Cmd) PromptMain {
	return NewPromptx(
		WithCommonOpions(
			WithCommonOptionCommands(cmds...),
		),
	)
}

// NewOptionCommandPromptx new with command and options
// use for replace NewCommandPromptx when you need apply other options.
// example: NewCommandPromptx(cmds...) => NewOptionCommandPromptx(NewPromptOptions(....),cmds...)
func NewOptionCommandPromptx(cc *PromptOptions, cmds ...*Cmd) PromptMain {
	if cc == nil {
		cc = NewPromptOptions()
	}
	cc.CommonOpions = append(cc.CommonOpions, WithCommonOptionCommands(cmds...))
	return newPromptx(cc)
}

// newPromptx new prompt
func newPromptx(cc *PromptOptions) *Promptx {
	p := new(Promptx)
	p.sets = make(map[string]*CommandSetOptions)
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

// ResetCommands 重置命令集合
func (p *Promptx) ResetCommands(commands ...*Cmd) {
	debug.Println("reset command ", len(commands))
	if iface, ok := p.cc.BlocksManager.(interface {
		ResetCommands(cmds ...*Cmd)
	}); ok {
		iface.ResetCommands(commands...)
	}
}

// RemoveHistory remove from history
func (p *Promptx) RemoveHistory(line string) {
	if iface, ok := p.cc.BlocksManager.(interface {
		RemoveHistory(line string)
	}); ok {
		iface.RemoveHistory(line)
	}
}

// AddHistory add line to history
func (p *Promptx) AddHistory(line string) {
	if iface, ok := p.cc.BlocksManager.(interface {
		AddHistory(line string)
	}); ok {
		iface.AddHistory(line)
	}
}

func (p *Promptx) ResetHistoryFile(filename string) {
	if iface, ok := p.cc.BlocksManager.(interface {
		ResetHistoryFile(fname string)
	}); ok {
		iface.ResetHistoryFile(filename)
	}
}

func (p *Promptx) SetCommandPreCheck(f func(ctx Context) error) {
	if iface, ok := p.cc.BlocksManager.(interface {
		SetCommandPreCheck(f func(ctx Context) error)
	}); ok {
		iface.SetCommandPreCheck(f)
	}
}

// AddCommandSet add command set,it will auto switch when first add commandset.
func (p *Promptx) AddCommandSet(name string, cmds []*Cmd, opts ...CommandSetOption) {
	if len(cmds) < 0 {
		panic(fmt.Sprintf("commandset %s do not have any commad", name))
	}
	if _, ok := p.sets[name]; ok {
		panic(fmt.Sprintf("commandset %s register repeated", name))
	}
	set := NewCommandSetOptions(opts...)
	set.ApplyOption(
		withCommandSetOptionName(name),
		withCommandSetOptionCommands(cmds...),
	)
	p.sets[set.Name] = set
	if len(p.sets) == 1 {
		p.SwitchCommandSet(name)
	}
}

// SwitchCommandSet switch to specify comands set
func (p *Promptx) SwitchCommandSet(name string, args ...interface{}) {
	set, ok := p.sets[name]
	if !ok {
		p.Printf("commandset %s not exists", name)
		return
	}
	p.ResetHistoryFile(set.HistoryFile)
	p.SetCommandPreCheck(set.PreCheck)
	p.ResetCommands(set.Commands...)
	p.SetPromptWords(set.PromptWords...)
	if set.ChangeNotify != nil {
		set.ChangeNotify(p, args)
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

// getConsoleWriter use for custom command args checker
func (p *Promptx) getConsoleWriter() output.ConsoleWriter {
	return p.cc.Output
}

// getPresetOptions use for custom command args
func (p *Promptx) getPresetOptions() (*InputOptions, *SelectOptions) {
	return p.inputCC, p.selectCC
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
