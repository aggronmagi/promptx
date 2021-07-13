package promptx

import (
	"fmt"
	"io"
	"sync"
	"syscall"

	"github.com/aggronmagi/promptx/internal/debug"
	"github.com/aggronmagi/promptx/internal/std"
	"go.uber.org/atomic"
	"golang.org/x/term"
)

// PromptOptionsOptionDeclareWithDefault promptx options
// generate by https://github.com/timestee/optiongen
//go:generate optionGen --option_with_struct_name=false --v=true
func PromptOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Stdin":  io.ReadCloser(std.Stdin),
		"Stdout": io.Writer(std.Stdout),
		// event chan size
		"ChanSize": int(256),
		// default global input options
		"InputOptions": []InputOption(nil),
		// default global select options
		"SelectOptions": []SelectOption(nil),
		// default common options. use to create default optins
		"CommonOpions": []CommonOption(nil),
		// default manager. if it is not nil, ignore CommonOpions.
		"BlocksManager": BlocksManager(nil),
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

	// terminal
	t *Terminal
	// current blocks manager
	mgr BlocksManager
	// next blocks manager
	next BlocksManager
	// backup blocks manager
	backup BlocksManager
	// has exchange status
	exchange atomic.Bool
	// mgr,next value protect
	rw sync.RWMutex
	// for exchange other mode, exec should run in another gorountine.
	execCh    chan func()
	m         sync.Mutex
	cond      *sync.Cond
	notRead   atomic.Bool
	syncCh    chan struct{}
	refreshCh chan struct{}

	// is start run
	start atomic.Bool
	// stop chan
	stop chan struct{}
	// wait finish
	wg sync.WaitGroup
}

// NewPromptx new prompt
func NewPromptx(opts ...PromptOption) *Promptx {
	cc := NewPromptOptions(opts...)
	p := new(Promptx)
	p.selectCC = NewSelectOptions(cc.SelectOptions...)
	p.inputCC = NewInputOptions(cc.InputOptions...)
	if cc.BlocksManager == nil {
		cc.BlocksManager = NewDefaultBlockManger(cc.CommonOpions...)
	}
	p.rw.Lock()
	defer p.rw.Unlock()
	p.cc = cc
	p.mgr = cc.BlocksManager
	p.cond = sync.NewCond(&p.m)
	p.m.Lock()
	p.t = NewTerminal(p.cc.Stdin, p.cc.Stdout, p.cc.ChanSize)
	p.mgr.SetWriter(NewConsoleWriter(p.t))
	p.syncCh = make(chan struct{}, 1)
	p.refreshCh = make(chan struct{}, 1)
	return p
}

// Start start run async
func (p *Promptx) Start() (err error) {
	// already running
	if !p.start.CAS(false, true) {
		return
	}
	p.wg.Add(1)
	go p.run()
	return
}

// Run run prompt
func (p *Promptx) Run() (err error) {
	err = p.Start()
	if err != nil {
		return
	}
	p.wg.Wait()
	return
}

func (p *Promptx) Stop() {
	// already stop
	if !p.start.CAS(true, false) {
		return
	}
	if p.stop == nil {
		return
	}
	p.stop <- struct{}{}
}

// run internal
func (p *Promptx) run() (err error) {

	defer func() {
		p.start.Store(false)
		p.wg.Done()
	}()
	if p.t == nil {
		p.t = NewTerminal(p.cc.Stdin, p.cc.Stdout, p.cc.ChanSize)
		p.mgr.SetWriter(NewConsoleWriter(p.t))
	}
	// update windows size
	w, h, err := term.GetSize(syscall.Stdout)
	if err != nil {
		return err
	}
	p.mgr.UpdateWinSize(w, h)
	// render pre
	p.mgr.Render(NormalStatus)
	p.mgr.SetExecFunc(p.Exec)
	p.mgr.SetExecContext(p)

	p.stop = make(chan struct{})
	p.t.Start()
	defer func() {
		close(p.stop)
		p.t.Close()
		p.stop = nil
		p.t = nil
	}()

	exitCh := make(chan int)
	winSize := make(chan *WinSize)
	p.execCh = make(chan func())
	go HandleSignals(exitCh, winSize, p.stop)

	go func() {
		for {
			select {
			case <-p.stop:
				return
			case f := <-p.execCh:
				go func() {
					f()
					if p.exchange.Load() {
						p.refreshCh <- struct{}{}
					}
					p.cond.Signal()
				}()
			}
		}
	}()

	// event chan
	for {
		select {
		case in, ok := <-p.t.InputChan():
			if !ok {
				return
			}
			key := GetKey(in)

			if p.getCurrent().Event(key, in) {
				if p.exchangeNext(true) {
					debug.Println("recv exit. but change next screen")
					break
				}
				debug.Println("recv exit and not change next screen")
				return
			}
			p.exchangeNext(true)
			p.exchange.Store(false)
		case <-p.refreshCh:
			p.getCurrent().Render(NormalStatus)
		case <-p.syncCh:
			p.exchangeNext(true)
			p.exchange.Store(false)
		case size := <-winSize:
			p.getCurrent().UpdateWinSize(size.Col, size.Row)
			p.getCurrent().Render(NormalStatus)
		case code := <-exitCh:
			p.getCurrent().Render(CancelStatus)
			fmt.Println("exit code", code)
			// os.Exit(code)
			return
		}
	}

	return
}

func (p *Promptx) Exec(f func()) {
	p.notRead.Store(true)
	p.execCh <- f
	p.cond.Wait()
	p.notRead.Store(false)
}

func (p *Promptx) getCurrent() BlocksManager {
	p.rw.RLock()
	defer p.rw.RUnlock()
	return p.mgr
}

func (p *Promptx) exchangeNext(render bool) (change bool) {
	p.rw.RLock()
	if p.next == nil {
		p.rw.RUnlock()
		return false
	}
	p.rw.RUnlock()
	p.rw.Lock()
	defer p.rw.Unlock()
	debug.Println(fmt.Sprintf("exchange %T => %T %t", p.mgr, p.next, render))
	// update and fix next data
	p.next.SetExecFunc(p.Exec)
	p.next.SetExecContext(p)
	p.next.SetWriter(p.mgr.Writer())
	p.next.UpdateWinSize(p.mgr.Columns(), p.mgr.Rows())
	// notify change and revert internal status
	p.next.ChangeStatus()
	// render
	if render {
		p.next.Render(NormalStatus)
	}

	// change next
	p.backup = p.mgr
	p.mgr = p.next
	p.next = nil
	p.exchange.Store(true)
	return true
}

// ChangeMode change current block manager.
func (p *Promptx) ChangeMode(next BlocksManager) {
	p.rw.RLock()
	if next == p.mgr {
		p.rw.RUnlock()
		return
	}
	p.rw.RUnlock()
	p.rw.Lock()
	defer p.rw.Unlock()
	p.mgr.SetChangeStatus(1)
	p.next = next
	p.backup = p.mgr
	p.cond.Signal()
	p.syncCh <- struct{}{}
}

// RevertMode revert last mode to next
func (p *Promptx) RevertMode() {
	if p.backup == nil || p.backup == p.next {
		p.backup = nil
		return
	}
	backup := p.backup
	p.backup = nil
	p.ChangeMode(backup)
}

// ResetDefaultMode reset default mode
func (p *Promptx) ResetDefaultMode() {
	p.ChangeMode(p.cc.BlocksManager)
}

// EnterRawMode enter raw mode for read key press real time
func (p *Promptx) EnterRawMode() (err error) {
	return p.t.EnterRawMode()
}

// ExitRawMode exit raw mode
//
//BUG(Terminal) use `exec` package to run interactive command will cause
// exception. `Terminal` will catch your `stdin` input even if you call the function.
// NOTE: use `reset` command to recover your terminal.
func (p *Promptx) ExitRawMode() (err error) {
	return p.t.ExitRawMode()
}

// Stdout return a wrap stdout writer. it can refersh view correct
func (p *Promptx) Stdout() io.Writer {
	return &wrapWriter{
		p:      p,
		target: p.t,
	}
}

// Stderr std err
func (p *Promptx) Stderr() io.Writer {
	return &wrapWriter{
		p:      p,
		target: std.Stderr,
	}
}

// ClearScreen clears the screen.
func (p *Promptx) ClearScreen() {
	out := p.mgr.Writer()
	out.EraseScreen()
	out.CursorGoTo(0, 0)
	debug.AssertNoError(out.Flush())
}

// SetTitle set title
func (p *Promptx) SetTitle(title string) {
	if len(title) < 1 {
		return
	}
	out := p.mgr.Writer()
	out.SetTitle(title)
	debug.AssertNoError(out.Flush())
}

// ClearTitle clear title
func (p *Promptx) ClearTitle() {
	out := p.mgr.Writer()
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
	p.ChangeMode(input)
	input.cond.Wait()
	p.ResetDefaultMode()
	//p.RevertMode()
	// Set exhange state to refresh when command exec finish
	p.exchange.Store(true)
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
	p.ChangeMode(sel)
	sel.cond.Wait()
	p.ResetDefaultMode()
	// p.RevertMode()
	// Set exhange state to refresh when command exec finish
	p.exchange.Store(true)
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
	p.ChangeMode(sel)
	sel.cond.Wait()
	p.ResetDefaultMode()
	// p.RevertMode()
	// Set exhange state to refresh when command exec finish
	p.exchange.Store(true)
	return
}

type wrapWriter struct {
	p      *Promptx
	target io.Writer
}

func (w *wrapWriter) Write(b []byte) (n int, err error) {
	// if w.p.notRead.Load() {
	// 	return w.target.Write(b)
	// }
	// 	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&Current)), unsafe.Pointer(Updating))

	// (BlocksManager)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&Current))))

	w.p.mgr.Clear()
	n, err = w.target.Write(b)
	if w.p.exchange.Load() {
		debug.Println("after call exchange next.")
		w.p.exchangeNext(false)
	}

	w.p.mgr.Render(NormalStatus)
	return n, err
}
