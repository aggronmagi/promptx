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
		"Inputs": []InputOption(nil),
		// default global select options
		"Selects": []SelectOption(nil),
		// default common options. use to create default optins
		"Common": []CommonOption(nil),
		// default manager. if it is not nil, ignore Common.
		"Manager": BlocksManager(nil),
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
//go:generate gogen option -n CommandSetOption -o gen_options_commandset.go
func promptxCommandSetOptions() interface{} {
	return map[string]interface{}{
		// comman set name
		"Name": "",
		// command set commands
		"Cmds": []*Cmd{},
		// Set the operation record of the current command set to persist the saved file. If not set, the history operation record will be cleared every time the command set is switched.
		"History": "",
		// Set the pre-detection function for all commands executed in the current command set
		"PreCheck": (func(ctx Context) error)(nil),
		// Set the prompt string when switching to the command set(custom word color)
		"Prompt": []*Word{},
		// Set the notification function when switching to the command set, args is the parameter passed by SwitchCommandSet.
		"OnChange": func(ctx Context, args []interface{}) {},
	}
}

// WithPromptStr Set the prompt string when switching to the command set(default color)
func WithPromptStr(prompt string) CommandSetOption {
	return func(cc *CommandSetOptions) CommandSetOption {
		previous := cc.Prompt
		cc.Prompt = []*Word{WordDefault(prompt)}
		return WithPrompt(previous...)
	}
}

func init() {
	InstallCommandSetOptionsWatchDog(func(cc *CommandSetOptions) {
		if len(cc.Prompt) < 1 {
			cc.Prompt = append(cc.Prompt, WordDefault(">> "))
		}
	})
}

// Terminal core terminal operations
type Terminal interface {
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

// Interaction user-facing prompt methods
type Interaction interface {
	RawInput(tip string, opts ...InputOption) (result string, err error)
	RawSelect(tip string, list []string, opts ...SelectOption) (result int)
	RawMulSel(tip string, list []string, opts ...SelectOption) (result []int)
}

// Commander command set management and execution
type Commander interface {
	// Stop stop run
	Stop()
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
}

// Context Run Command Context
type Context interface {
	Terminal
	Interaction
	Commander
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
	mgr *commandManager
}

var _ Context = &Promptx{}

// New creates a new prompt application
func New(opts ...PromptOption) PromptMain {
	cc := NewPromptOptions(opts...)
	return newPromptx(cc)
}

// newPromptx new prompt
func newPromptx(cc *PromptOptions) *Promptx {
	p := new(Promptx)
	p.mgr = newCommandManager(p)
	p.selectCC = NewSelectOptions(cc.Selects...)
	p.inputCC = NewInputOptions(cc.Inputs...)
	if cc.Manager == nil {
		cc.Manager = NewDefaultBlockManger(cc.Common...)
	}
	p.console = terminal.NewTerminalApp(cc.Input)
	cc.Manager.SetWriter(cc.Output)
	cc.Manager.SetExecContext(p)
	cc.Manager.UpdateWinSize(cc.Input.GetWinSize())
	p.cc = cc

	return p
}

// Run run prompt
func (p *Promptx) Run() (err error) {
	p.console.Run(p.cc.Manager)
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
	if iface, ok := p.cc.Manager.(interface {
		SetPrompt(prompt string)
	}); ok {
		iface.SetPrompt(prompt)
		// p.syncCh <- struct{}{}
	}
}

// SetPromptWords update prompt string. custom display.
func (p *Promptx) SetPromptWords(words ...*Word) {
	if iface, ok := p.cc.Manager.(interface {
		SetPromptWords(words ...*Word)
	}); ok {
		iface.SetPromptWords(words...)
	}
}

func (p *Promptx) ExecCommand(args []string) {
	if iface, ok := p.cc.Manager.(interface {
		ExecCommand(args []string)
	}); ok {
		iface.ExecCommand(args)
	}
}

// ResetCommands 重置命令集合
func (p *Promptx) ResetCommands(commands ...*Cmd) {
	debug.Println("reset command ", len(commands))
	if iface, ok := p.cc.Manager.(interface {
		ResetCommands(cmds ...*Cmd)
	}); ok {
		iface.ResetCommands(commands...)
	}
}

// RemoveHistory remove from history
func (p *Promptx) RemoveHistory(line string) {
	if iface, ok := p.cc.Manager.(interface {
		RemoveHistory(line string)
	}); ok {
		iface.RemoveHistory(line)
	}
}

// AddHistory add line to history
func (p *Promptx) AddHistory(line string) {
	if iface, ok := p.cc.Manager.(interface {
		AddHistory(line string)
	}); ok {
		iface.AddHistory(line)
	}
}

func (p *Promptx) ResetHistoryFile(filename string) {
	if iface, ok := p.cc.Manager.(interface {
		ResetHistoryFile(fname string)
	}); ok {
		iface.ResetHistoryFile(filename)
	}
}

func (p *Promptx) SetCommandPreCheck(f func(ctx Context) error) {
	if iface, ok := p.cc.Manager.(interface {
		SetCommandPreCheck(f func(ctx Context) error)
	}); ok {
		iface.SetCommandPreCheck(f)
	}
}

// AddCommandSet add command set,it will auto switch when first add commandset.
func (p *Promptx) AddCommandSet(name string, cmds []*Cmd, opts ...CommandSetOption) {
	p.mgr.AddCommandSet(name, cmds, opts...)
}

// SwitchCommandSet switch to specify comands set
func (p *Promptx) SwitchCommandSet(name string, args ...interface{}) {
	p.mgr.SwitchCommandSet(name, args...)
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

// getConsoleWriter use for custom command args checker
func (p *Promptx) getConsoleWriter() output.ConsoleWriter {
	return p.cc.Output
}

// getPresetOptions use for custom command args
func (p *Promptx) getPresetOptions() (*InputOptions, *SelectOptions) {
	return p.inputCC, p.selectCC
}
