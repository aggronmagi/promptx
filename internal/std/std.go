package std

import (
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/aggronmagi/promptx/internal/debug"
)

var (
	Stdin  io.ReadCloser  = os.Stdin
	Stdout io.WriteCloser = os.Stdout
	Stderr io.WriteCloser = os.Stderr
)

type CancelableStdin struct {
	r      io.Reader
	mutex  sync.Mutex
	stop   chan struct{}
	closed int32
	notify chan struct{}
	data   []byte
	read   int
	err    error
}

func NewCancelableStdin(r io.Reader) *CancelableStdin {
	c := &CancelableStdin{
		r:      r,
		notify: make(chan struct{}),
		stop:   make(chan struct{}),
	}
	go c.ioloop()
	return c
}

func (c *CancelableStdin) ioloop() {
loop:
	for {
		select {
		case <-c.notify:
			c.read, c.err = c.r.Read(c.data)
			select {
			case c.notify <- struct{}{}:
			case <-c.stop:
				debug.Println("recv stop  stdin sign")
				break loop
			}
		case <-c.stop:
			debug.Println("recv stop  stdin sign")
			break loop
		}
	}
}

func (c *CancelableStdin) Read(b []byte) (n int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if atomic.LoadInt32(&c.closed) == 1 {
		return 0, io.EOF
	}

	c.data = b
	select {
	case c.notify <- struct{}{}:
	case <-c.stop:
		debug.Println("recv stop  stdin sign")
		return 0, io.EOF
	}
	select {
	case <-c.notify:
		return c.read, c.err
	case <-c.stop:
		debug.Println("recv stop  stdin sign")
		return 0, io.EOF
	}
}

func (c *CancelableStdin) Close() error {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		close(c.stop)
	}
	debug.Println("close stdin", atomic.LoadInt32(&c.closed))
	return nil
}

// FillableStdin is a stdin reader which can prepend some data before
// reading into the real stdin
type FillableStdin struct {
	sync.Mutex
	stdin       io.Reader
	stdinBuffer io.ReadCloser
	buf         []byte
	bufErr      error
}

// NewFillableStdin gives you FillableStdin
func NewFillableStdin(stdin io.Reader) (io.ReadCloser, io.Writer) {
	r, w := io.Pipe()
	s := &FillableStdin{
		stdinBuffer: r,
		stdin:       stdin,
	}
	s.ioloop()
	return s, w
}

func (s *FillableStdin) ioloop() {
	go func() {
		for {
			bufR := make([]byte, 100)
			var n int
			n, s.bufErr = s.stdinBuffer.Read(bufR)
			if s.bufErr != nil {
				if s.bufErr == io.ErrClosedPipe {
					break
				}
			}
			s.Lock()
			s.buf = append(s.buf, bufR[:n]...)
			s.Unlock()
		}
	}()
}

// Read will read from the local buffer and if no data, read from stdin
func (s *FillableStdin) Read(p []byte) (n int, err error) {
	s.Lock()
	i := len(s.buf)
	if len(p) < i {
		i = len(p)
	}
	if i > 0 {
		n := copy(p, s.buf)
		s.buf = s.buf[:0]
		cerr := s.bufErr
		s.bufErr = nil
		s.Unlock()
		return n, cerr
	}
	s.Unlock()
	n, err = s.stdin.Read(p)
	return n, err
}

func (s *FillableStdin) Close() error {
	s.stdinBuffer.Close()
	if close, ok := s.stdin.(io.Closer); ok {
		close.Close()
	}
	return nil
}
