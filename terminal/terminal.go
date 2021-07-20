package terminal

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/aggronmagi/promptx/input"
	"github.com/aggronmagi/promptx/internal/debug"
	"github.com/aggronmagi/promptx/output"
	"go.uber.org/atomic"
)

type App interface {
	Event(key input.Key, in []byte) (exit bool)
	Clear()
	Refresh()
	ExitSign(code int)
	UpdateWinSize(*input.WinSize)
	Setup(*input.WinSize)
	TearDown()
	IsInTask() bool
}

type TerminalApp struct {
	// terminal input
	in input.ConsoleParser
	// quit chan close when need exit
	quit chan struct{}
	// terminal size update chan
	sizeCh chan *input.WinSize
	// system sign ch
	signCh    chan int
	closeSign chan struct{}
	// read chan
	bufCh     chan []byte
	closeRead chan struct{}
	// rawmode flag
	rawmode atomic.Bool
	running atomic.Bool
	// running ptr
	appPtr  atomic.UnsafePointer
	m       sync.Mutex
	appList list.List
}

func NewTerminalApp(in input.ConsoleParser) *TerminalApp {
	return &TerminalApp{
		in:        in,
		sizeCh:    make(chan *input.WinSize),
		signCh:    make(chan int),
		closeSign: make(chan struct{}, 1),
		bufCh:     make(chan []byte),
		closeRead: make(chan struct{}, 1),
	}
}

func (t *TerminalApp) Run(app App) {
	debug.Println("running: ", fmt.Sprintf("%T", interface{}(app)))
	// start running
	if t.running.CAS(false, true) {
		t.quit = make(chan struct{})
		debug.Println("enter running", os.Getpid())
		t.EnterRaw()
		defer func() {
			debug.Println("exit running", os.Getpid())
			t.ExitRaw()
			if t.running.CAS(true, false) {
				close(t.quit)
			}
		}()
	}
	debug.Println("retry raw")
	// try recover raw mode
	t.EnterRaw()
	debug.Println("ready store app")
	t.appPtr.Store(unsafe.Pointer(&app))
	t.m.Lock()
	t.appList.PushBack(app)
	t.m.Unlock()
	defer func() {
		debug.Println("exit app", fmt.Sprintf("%T", interface{}(app)))
		t.m.Lock()
		defer t.m.Unlock()
		if t.appList.Len() < 1 {
			debug.Println("no more app left")
			return
		}
		t.appList.Remove(t.appList.Back())
		if t.appList.Len() == 0 {
			t.appPtr.Store(nil)
		} else {
			lastApp := t.appList.Back().Value.(App)
			debug.Println("revert app", fmt.Sprintf("%T", interface{}(lastApp)))
			t.appPtr.Store(unsafe.Pointer(&lastApp))
		}
	}()
	debug.Println("setup app")
	app.Setup(t.in.GetWinSize())
	defer app.TearDown()

	for {
		select {
		case w := <-t.sizeCh:
			app.UpdateWinSize(w)
		case code := <-t.signCh:
			app.ExitSign(code)
			if t.running.CAS(true, false) {
				close(t.quit)
			}
		case in := <-t.bufCh:
			key := input.GetKey(in)
			debug.Println("read from input", key, len(in))
			if app.Event(key, in) {
				debug.Println("recv exit app")
				return
			}
		case <-t.quit:
			return
		}
	}
}

func (t *TerminalApp) Stop() {
	if !t.running.CAS(true, false) {
		return
	}
	close(t.quit)
}

func (t *TerminalApp) ExitRaw() (err error) {
	// start read
	if !t.rawmode.CAS(true, false) {
		return
	}
	debug.Println("exit raw mode")
	t.closeSign <- struct{}{}
	t.closeRead <- struct{}{}
	err = t.in.TearDown()
	debug.AssertNoError(err)
	return
}

func (t *TerminalApp) EnterRaw() (err error) {
	// start read
	if !t.rawmode.CAS(false, true) {
		return
	}
	debug.Println("enter raw mode")
	// clean cache input
	for k := 0; k < len(t.bufCh); k++ {
		<-t.bufCh
	}
	err = t.in.Setup()
	debug.AssertNoError(err)
	go t.handleSignals(t.signCh, t.sizeCh, t.closeSign)
	go t.readBuffer(t.bufCh, t.closeRead)

	return
}

func (t *TerminalApp) readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	debug.Log("start reading buffer")
	for {
		select {
		case <-stopCh:
			debug.Log("stop reading buffer")
			return
		default:
			if b, err := t.in.Read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (t *TerminalApp) GetCurrentApp() App {
	ptr := t.appPtr.Load()
	if ptr != nil {
		return *((*App)(ptr))
	}
	return nil
}

type wrapWriter struct {
	t   *TerminalApp
	out output.ConsoleWriter
}

func NewWrapWriter(out output.ConsoleWriter, t *TerminalApp) io.Writer {
	return &wrapWriter{
		t:   t,
		out: out,
	}
}

func (w *wrapWriter) Write(b []byte) (n int, err error) {
	ptr := w.t.appPtr.Load()
	if ptr != nil {
		app := *((*App)(ptr))
		app.Clear()
		w.out.WriteRaw(b)
		err = w.out.Flush()
		n = len(b)
		if !app.IsInTask() {
			app.Refresh()
		}
		return
	}
	w.out.WriteRaw(b)
	err = w.out.Flush()
	n = len(b)

	return
}
