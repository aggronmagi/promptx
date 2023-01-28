package promptx

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

var _ PressContext = &pressContext{}

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

var _ PrintContext = &consoleContext{}

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
