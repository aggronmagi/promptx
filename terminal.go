package promptx

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"syscall"

	"github.com/aggronmagi/promptx/internal/std"
	"go.uber.org/atomic"
	"golang.org/x/term"
)

// // KeyPressListener tty event listener
// type KeyPressListener interface {
// 	Event(key Key, buf []byte) (exit bool)
// }

// // KeyPressListenFunc event func
// type KeyPressListenFunc func(key Key, buf []byte) (exit bool)

// // Event implement KeyPressListener interface
// func (v KeyPressListenFunc) Event(key Key, buf []byte) (exit bool) {
// 	return v(key, buf)
// }

// var _ KeyPressListener = KeyPressListenFunc(nil)

type RawMode struct {
	state *term.State
	fd    int
}

func NewRawMode(fd int) (r *RawMode) {
	return &RawMode{
		fd: fd,
	}
}

func (r *RawMode) Enter() (err error) {
	r.state, err = term.MakeRaw(r.fd)
	return err
}

func (r *RawMode) Exit() error {
	if r.state == nil {
		return nil
	}
	err := term.Restore(r.fd, r.state)
	r.state = nil
	return err
}

// Terminal interaction with console
type Terminal struct {
	m      sync.Mutex
	raw    *RawMode
	stdin  io.ReadCloser
	stdout io.Writer
	wg     sync.WaitGroup
	start  atomic.Bool
	ch     chan []byte
}

// NewTerminal new terminal
func NewTerminal(stdin io.ReadCloser, stdout io.Writer, size int) *Terminal {
	t := &Terminal{
		raw:    NewRawMode(syscall.Stdin),
		stdin:  std.NewCancelableStdin(stdin),
		stdout: stdout,
		ch:     make(chan []byte, size),
	}

	return t
}

// ioloop internal io loop
func (t *Terminal) ioloop() {
	defer func() {
		t.raw.Exit()
		t.wg.Done()
	}()

	buf := make([]byte, 4096)
	for {
		n, err := t.stdin.Read(buf)
		if err != nil {
			// stdin close
			if err == io.EOF {
				t.ch <- nil
				break
			}
			if strings.Contains(err.Error(), "interrupted system call") {
				continue
			}
			t.stdout.Write([]byte(
				"terminal read touch error. " + err.Error() + "\n"),
			)
			break
		}

		ev := make([]byte, n)
		copy(ev, buf[:n])
		t.ch <- ev
		// // notify event
		// if t.listener.Event(key, buf[:n]) {
		// 	break
		// }
	}
}

////////////////////////////////////////////////////////////////////////////////
// Terminal Export Func

// InputChan get console input chan
func (t *Terminal) InputChan() (ch <-chan []byte) {
	return t.ch
}

// Run Terminal
func (t *Terminal) Run() {
	if !t.start.CAS(false, true) {
		return
	}
	t.wg.Add(1)
	t.raw.Enter()
	t.ioloop()
}

// Start start terminal
func (t *Terminal) Start() {
	if !t.start.CAS(false, true) {
		return
	}
	t.wg.Add(1)
	t.raw.Enter()
	go t.ioloop()
}

// Close close terminal. io.Closer interface implement
func (t *Terminal) Close() error {
	if !t.start.CAS(true, false) {
		return nil
	}
	if t.stdin != nil {
		t.stdin.Close()
	}
	t.wg.Wait()
	return nil
}

// Write io.Writer interface implement
func (t *Terminal) Write(b []byte) (int, error) {
	return t.stdout.Write(b)
}

// Print print string
func (t *Terminal) Print(s string) {
	fmt.Fprintf(t.stdout, "%s", s)
}

// PrintRune print rune char
func (t *Terminal) PrintRune(r rune) {
	fmt.Fprintf(t.stdout, "%c", r)
}

// Bell bell sound
func (t *Terminal) Bell() {
	// ASCII Control Character 7 = BELL
	t.PrintRune(7)
	// fmt.Fprintf(t, "%c", CharBell)
}

// EnterRawMode enter raw mode for read key press real time
func (t *Terminal) EnterRawMode() (err error) {
	// clean last key press info
	for i := 0; i < len(t.ch); i++ {
		<-t.ch
	}
	return t.raw.Enter()
}

// ExitRawMode exit raw mode
func (t *Terminal) ExitRawMode() (err error) {
	return t.raw.Exit()
}

// // SleepToResume will sleep myself, and return only if I'm resumed.
// func (t *Terminal) SleepToResume() {
// 	if !atomic.CompareAndSwapInt32(&t.sleeping, 0, 1) {
// 		return
// 	}
// 	defer atomic.StoreInt32(&t.sleeping, 0)
// 	t.raw.Exit()
// 	ch := term.WaitForResume()
// 	term.SuspendMe()
// 	<-ch
// 	t.raw.Exit()
// }
