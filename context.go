package promptx

import "io"

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
	// ChangeMode
	ChangeMode(next BlocksManager)
	// RevertMode revert last mode to next
	RevertMode()
	// ResetDefaultMode reset default mode
	ResetDefaultMode()
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
	// Print = fmt.Print
	Print(v ...interface{})
	// Printf = fmt.Printf
	Printf(fmt string, v ...interface{})
	// Println = fmt.Println
	Println(v ...interface{})
}

////////////////////////////////////////////////////////////////////////////////

// PressContext export for logic loop
type PressContext interface {
	// GetBuffer get input buffer
	GetBuffer() *Buffer
	// console writer
	Writer() ConsoleWriter
	// GetKey get input keys
	GetKey() Key
	// GetInput get input bytes
	GetInput() []byte
}

type pressContext struct {
	buf   *Buffer
	out   ConsoleWriter
	key   Key
	input []byte
}

func (ctx *pressContext) GetBuffer() *Buffer {
	return ctx.buf
}

func (ctx *pressContext) Writer() ConsoleWriter {
	return ctx.out
}

// GetKey get input keys
func (ctx *pressContext) GetKey() Key {
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
	Writer() ConsoleWriter

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
	GetBuffer() *Buffer
	// SetBuffer set input buffer. if not set, cursor will be hided.
	SetBuffer(buf *Buffer)
	// Prepare pre calc line to prepare area
	Prepare() bool
	// Status return current status. 0: normal 1: finish 2:canel
	Status() int
}

type consoleContext struct {
	*BlocksBaseManager
	buf     *Buffer
	cursor  int
	prepare bool
	status  int
}

func (ctx *consoleContext) InputCursor() int {
	return ctx.cursor
}

// GetBuffer get input buffer
func (ctx *consoleContext) GetBuffer() *Buffer {
	return ctx.buf
}

// SetBuffer set input buffer. if not set, cursor will be hided.
func (ctx *consoleContext) SetBuffer(buf *Buffer) {
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
