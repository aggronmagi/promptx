package blocks

import (
	"fmt"
	"io"
	"strings"

	completion "github.com/aggronmagi/promptx/v2/completion"
	"github.com/aggronmagi/promptx/v2/input"
	"github.com/aggronmagi/promptx/v2/internal/debug"
	"github.com/aggronmagi/promptx/v2/output"
	"github.com/aggronmagi/promptx/v2/terminal"
)

// PromptOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n BlocksOption -o gen_options_blocks.go
func promptxBlocksOptions() interface{} {
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
		// Context
		"Context": Context(nil),
	}
}

////////////////////////////////////////////////////////////////////////////////

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

	// SetPrompt update prompt string. prompt will auto add space suffix.
	SetPrompt(prompt string)
	// SetPromptWords update prompt string. custom display.
	SetPromptWords(words ...*Word)

	// Stop stop run
	Stop()
}

// Interaction user-facing prompt methods
type Interaction interface {
	RawInput(tip string, opts ...InputOption) (result string, err error)
	RawSelect(tip string, list []string, opts ...SelectOption) (result int)
	RawMulSel(tip string, list []string, opts ...SelectOption) (result []int)
}

// Context Run Command Context
type Context interface {
	Terminal
	Interaction
}

type Controler interface {
	// RemoveHistory remove from history
	RemoveHistory(line string)
	// AddHistory add line to history
	AddHistory(line string)
	// reset history file
	ResetHistoryFile(filename string)

	GetPresetInputOptions() *InputOptions
	GetPresetSelectOptions() *SelectOptions
}

type Application interface {
	Context
	Controler
	Run() error
	// GetManager returns the BlocksManager used by this application.
	// This allows external packages to access the manager for configuration.
	GetManager() BlocksManager
}

type application struct {
	cc       *BlocksOptions
	inputCC  *InputOptions
	selectCC *SelectOptions
	console  *terminal.TerminalApp
}

var _ Application = &application{}

// New creates a new blocks application.
// It is the recommended entry point for creating a blocks application.
//
// Example:
//
//	app := blocks.New(blocks.WithCommon(blocks.WithPrefix(">>> ")))
func New(opts ...BlocksOption) Application {
	cc := NewBlocksOptions(opts...)
	return newApplication(cc)
}

// newApplication new application
func newApplication(cc *BlocksOptions) *application {
	app := new(application)
	if cc.Context == nil {
		cc.Context = app
	}
	app.selectCC = NewSelectOptions(cc.Selects...)
	app.inputCC = NewInputOptions(cc.Inputs...)
	if cc.Manager == nil {
		cc.Manager = NewDefaultBlockManger(cc.Common...)
	}
	app.console = terminal.NewTerminalApp(cc.Input)
	cc.Manager.SetWriter(cc.Output)
	cc.Manager.SetExecContext(cc.Context)
	cc.Manager.UpdateWinSize(cc.Input.GetWinSize())
	app.cc = cc

	return app
}

// RawInput get input
func (p *application) RawInput(tip string, opts ...InputOption) (result string, err error) {
	if !strings.HasSuffix(tip, " ") {
		tip = tip + " "
	}
	// copy a new config
	newCC := (*p.inputCC)

	// apply input opts
	newCC.ApplyOption(opts...)
	// set internal options
	newCC.SetOption(WithInputOptionOnFinish(func(input string, eof error) {
		result, err = input, eof
	}))
	if tip != "" {
		newCC.SetOption(WithInputOptionPrefix(tip))
	}
	//
	input := NewInputManager(&newCC)

	input.SetExecContext(p.cc.Context)
	input.SetWriter(p.cc.Output)
	input.UpdateWinSize(p.cc.Input.GetWinSize())
	debug.Println("enter input")
	p.console.Run(input)
	debug.Println("exit input")
	return
}

// RawSelect get select value with raw option
func (p *application) RawSelect(tip string, list []string, opts ...SelectOption) (result int) {
	if !strings.HasSuffix(tip, " ") {
		tip = tip + " "
	}
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	// reset options internal
	newCC.SetOption(WithSelectOptionOnFinish(
		func(sels []int) {
			if len(sels) < 1 {
				result = -1
				return
			}
			result = sels[0]
		},
	))
	newCC.SetOption(WithSelectOptionMulti(false))
	newCC.SetOption(WithSelectOptionTip(tip))

	// if opts set options. use opts values
	if len(newCC.Options) < 1 {
		if len(list) < 1 {
			return -1
		}
		for _, v := range list {
			newCC.Options = append(newCC.Options, &completion.Suggest{
				Text: v,
			})
		}
	}
	sel := NewSelectManager(&newCC)

	sel.SetExecContext(p.cc.Context)
	sel.SetWriter(p.cc.Output)
	sel.UpdateWinSize(p.cc.Input.GetWinSize())
	p.console.Run(sel)
	return
}

// RawMulSel get multiple value with raw option
func (p *application) RawMulSel(tip string, list []string, opts ...SelectOption) (result []int) {
	if !strings.HasSuffix(tip, " ") {
		tip = tip + " "
	}
	// copy new config
	newCC := *p.selectCC
	newCC.ApplyOption(opts...)
	newCC.SetOption(WithSelectOptionOnFinish(
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
			newCC.Options = append(newCC.Options, &completion.Suggest{
				Text: v,
			})
		}
	}
	sel := NewSelectManager(&newCC)
	sel.SetExecContext(p.cc.Context)
	sel.SetWriter(p.cc.Output)
	sel.UpdateWinSize(p.cc.Input.GetWinSize())
	p.console.Run(sel)
	return
}

// Run run application
func (p *application) Run() error {
	p.console.Run(p.cc.Manager)
	debug.Println("exit application run")
	return nil
}

// EnterRawMode enter raw mode for read key press real time
func (p *application) EnterRawMode() error {
	return p.console.EnterRaw()
}

// ExitRawMode exit raw mode
func (p *application) ExitRawMode() error {
	return p.console.ExitRaw()
}

// Stdout return a wrap stdout writer. it can refersh view correct
func (p *application) Stdout() io.Writer {
	return terminal.NewWrapWriter(p.cc.Output, p.console)
}

// Stderr std err
func (p *application) Stderr() io.Writer {
	return terminal.NewWrapWriter(p.cc.Stderr, p.console)
}

// ClearScreen clears the screen.
func (p *application) ClearScreen() {
	p.cc.Output.EraseScreen()
	p.cc.Output.CursorGoTo(0, 0)
	debug.AssertNoError(p.cc.Output.Flush())
}

// SetTitle set window title
func (p *application) SetTitle(title string) {
	if len(title) < 1 {
		return
	}
	p.cc.Output.SetTitle(title)
	debug.AssertNoError(p.cc.Output.Flush())
}

// ClearTitle clear window title
func (p *application) ClearTitle() {
	p.cc.Output.ClearTitle()
	debug.AssertNoError(p.cc.Output.Flush())
}

// Print = fmt.Print
func (p *application) Print(v ...interface{}) {
	fmt.Fprint(p.Stdout(), v...)
}

// Printf = fmt.Printf
func (p *application) Printf(format string, v ...interface{}) {
	fmt.Fprintf(p.Stdout(), format, v...)
}

// Println = fmt.Println
func (p *application) Println(v ...interface{}) {
	fmt.Fprintln(p.Stdout(), v...)
}

// WPrint  print words
func (p *application) WPrint(words ...*Word) {
	for _, v := range words {
		p.cc.Output.SetColor(v.TextColor, v.BGColor, v.Bold)
		p.cc.Output.WriteStr(v.Text)
	}
	p.cc.Output.SetColor(output.DefaultColor, output.DefaultColor, false)
}

// WPrintln print words and newline
func (p *application) WPrintln(words ...*Word) {
	for _, v := range words {
		p.cc.Output.SetColor(v.TextColor, v.BGColor, v.Bold)
		p.cc.Output.WriteStr(v.Text)
	}
	p.cc.Output.SetColor(output.DefaultColor, output.DefaultColor, false)
	p.cc.Output.WriteRawStr("\n")
}

// Stop stop run
func (p *application) Stop() {
	p.console.Stop()
}

// SetPrompt update prompt string. prompt will auto add space suffix.
func (p *application) SetPrompt(prompt string) {
	if iface, ok := p.cc.Manager.(interface {
		SetPrompt(prompt string)
	}); ok {
		iface.SetPrompt(prompt)
	}
}

// SetPromptWords update prompt string. custom display.
func (p *application) SetPromptWords(words ...*Word) {
	if iface, ok := p.cc.Manager.(interface {
		SetPromptWords(words ...*Word)
	}); ok {
		iface.SetPromptWords(words...)
	}
}

// RemoveHistory remove from history
func (p *application) RemoveHistory(line string) {
	if iface, ok := p.cc.Manager.(interface {
		RemoveHistory(line string)
	}); ok {
		iface.RemoveHistory(line)
	}
}

// AddHistory add line to history
func (p *application) AddHistory(line string) {
	if iface, ok := p.cc.Manager.(interface {
		AddHistory(line string)
	}); ok {
		iface.AddHistory(line)
	}
}

// ResetHistoryFile reset history file
func (p *application) ResetHistoryFile(filename string) {
	if iface, ok := p.cc.Manager.(interface {
		ResetHistoryFile(fname string)
	}); ok {
		iface.ResetHistoryFile(filename)
	}
}

// GetPresetInputOptions get preset input options
func (p *application) GetPresetInputOptions() *InputOptions {
	return p.inputCC
}

// GetPresetSelectOptions get preset select options
func (p *application) GetPresetSelectOptions() *SelectOptions {
	return p.selectCC
}

// GetManager returns the BlocksManager used by this application
func (p *application) GetManager() BlocksManager {
	return p.cc.Manager
}
