package blocks

import (
	buffer "github.com/aggronmagi/promptx/v2/buffer"
	"github.com/aggronmagi/promptx/v2/input"
	"github.com/aggronmagi/promptx/v2/output"
)

////////////////////////////////////////////////////////////////////////////////

// PressContext export for logic loop
type PressContext interface {
	// GetBuffer get input buffer
	GetBuffer() *buffer.Buffer
	// console writer
	Writer() output.ConsoleWriter
	// GetKey get input keys
	GetKey() input.Key
	// GetInput get input bytes
	GetInput() []byte
}

var _ PressContext = &pressContext{}

type pressContext struct {
	buf   *buffer.Buffer
	out   output.ConsoleWriter
	key   input.Key
	input []byte
}

func (ctx *pressContext) GetBuffer() *buffer.Buffer {
	return ctx.buf
}

func (ctx *pressContext) Writer() output.ConsoleWriter {
	return ctx.out
}

// GetKey get input keys
func (ctx *pressContext) GetKey() input.Key {
	return ctx.key
}

// GetInput get input bytes
func (ctx *pressContext) GetInput() []byte {
	return ctx.input
}

////////////////////////////////////////////////////////////////////////////////

// PrintContext context
type PrintContext interface {
	// console writer
	Writer() output.ConsoleWriter

	// ToPos returns the relative position from the beginning of the string.
	ToPos(cursor int) (x, y int)
	// Columns return console window col count
	Columns() int
	// Rows return console window row count
	Rows() int

	// Backward moves cursor to Backward from a current cursor position
	// regardless there is a line break.
	Backward(from, n int) int
	// Move moves cursor to specified position from the beginning of input
	// even if there is a line break.
	Move(from, to int) int
	// PrepareArea prepare enough area to display info
	PrepareArea(lines int)
	// LineWrap line wrap
	LineWrap(cursor int)

	// InputCursor get input buffer cursor pos. if no cursor,return -1.
	InputCursor() int
	// GetBuffer get input buffer
	GetBuffer() *buffer.Buffer
	// SetBuffer set input buffer. if not set, cursor will be hided.
	SetBuffer(buf *buffer.Buffer)
	// Prepare pre calc line to prepare area
	Prepare() bool
	// Status return current status. 0: normal 1: finish 2:canel
	Status() int
}

var _ PrintContext = &consoleContext{}

type consoleContext struct {
	*BlocksBaseManager
	buf     *buffer.Buffer
	cursor  int
	prepare bool
	status  int
}

func (ctx *consoleContext) InputCursor() int {
	return ctx.cursor
}

// GetBuffer get input buffer
func (ctx *consoleContext) GetBuffer() *buffer.Buffer {
	return ctx.buf
}

// SetBuffer set input buffer. if not set, cursor will be hided.
func (ctx *consoleContext) SetBuffer(buf *buffer.Buffer) {
	ctx.buf = buf
}

// Prepare pre calc line to prepare area
func (ctx *consoleContext) Prepare() bool {
	return ctx.prepare
}

// Status return current status. 0: normal 1: finish 2:canel
func (ctx *consoleContext) Status() int {
	return ctx.status
}

const (
	// NormalStatus for normal state
	NormalStatus = iota
	// FinishStatus last input finish
	FinishStatus
	// CancelStatus cancel input
	CancelStatus
)
