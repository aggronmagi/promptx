package promptx

import (
	"runtime"

	"github.com/aggronmagi/promptx/internal/debug"
	"go.uber.org/atomic"
)

type BlocksManager interface {
	// Setup to initialize console output.
	Setup(size *WinSize)
	// TearDown to clear title and erasing.
	TearDown()

	Refresh()
	ExitSign(code int)

	// SetWriter set console writer
	SetWriter(out ConsoleWriter)
	// Writer get console writer interface
	Writer() ConsoleWriter
	// SetChangeStatus  change mode status
	//
	// 0: not change(default value)
	// 1: ready to leave
	SetChangeStatus(st int)
	// change status notify
	ChangeStatus()

	// UpdateWinSize called when window size is changed.
	UpdateWinSize(size *WinSize)
	// Columns return console window col count
	Columns() int
	// Rows return console window row count
	Rows() int

	// Event deal console key press
	Event(key Key, in []byte) (exit bool)
	// Render renders to the console.
	Render(status int)
	// clear blocks last render
	Clear()

	// use for context
	SetExecContext(ctx Context)
	GetContext() Context

	// IsInTask is running event
	IsInTask() bool
}

var _ BlocksManager = &BlocksBaseManager{}

type EventCall func(ctx PressContext, key Key, in []byte) (exit bool)

// BlocksBaseManager manage some blocks to merge self view
type BlocksBaseManager struct {
	// console writer
	out ConsoleWriter

	//
	eventBefore EventCall
	eventBehind EventCall
	// action config
	cancelKey Key
	finishKey Key
	preCheck  func(status int, doc *Buffer) bool
	callback  func(status int, doc *Buffer) bool
	// cancel key touch auto exit
	cancelNotExit bool

	execFunc func(exec func())

	// windows size
	row int
	col int
	// previous cursor size
	previousCursor int
	// change mode status
	//
	// 0: not change(default value)
	// 1: ready to leave
	changeStatus int

	// children blocks
	children []ConsoleBlocks
	// major is input buffer
	major ConsoleBlocks

	inTask atomic.Bool

	// run context
	ctx Context
}

func (m *BlocksBaseManager) SetBeforeEvent(f EventCall) {
	m.eventBefore = f
}

func (m *BlocksBaseManager) SetBehindEvent(f EventCall) {
	m.eventBehind = f
}

func (m *BlocksBaseManager) SetExecContext(ctx Context) {
	m.ctx = ctx
}

func (m *BlocksBaseManager) GetContext() Context {
	return m.ctx
}

// SetChangeStatus modify change status
func (m *BlocksBaseManager) SetChangeStatus(st int) {
	m.changeStatus = st
}

// SetCancelKey set cancel key
func (m *BlocksBaseManager) SetCancelKey(k Key) {
	m.cancelKey = k
}

// SetFinishKey set finish key press
func (m *BlocksBaseManager) SetFinishKey(k Key) {
	m.finishKey = k
}

// SetPreCheck pre check function
// its arg status is CancelStatus or FinishStatus
// return true will call callback function later,
// otherwise render normal
func (m *BlocksBaseManager) SetPreCheck(f func(status int, doc *Buffer) bool) {
	m.preCheck = f
}

// SetCallBack set call back notify
func (m *BlocksBaseManager) SetCallBack(f func(status int, doc *Buffer) bool) {
	m.callback = f
}

// Setup to initialize console output.
func (m *BlocksBaseManager) Setup(size *WinSize) {
	// if title != "" {
	// 	m.out.SetTitle(title)
	// 	debug.AssertNoError(m.out.Flush())
	// }
	m.UpdateWinSize(size)
	m.Render(NormalStatus)
}

// SetCancelKeyAutoExit set exit if press cancel key and input buffer is ""
func (m *BlocksBaseManager) SetCancelKeyAutoExit(exit bool) {
	m.cancelNotExit = !exit
}

// TearDown to clear title and erasing.
func (m *BlocksBaseManager) TearDown() {
	// m.out.ClearTitle()
	// m.out.EraseDown()
	// debug.AssertNoError(m.out.Flush())
}

func (m *BlocksBaseManager) ExitSign(code int) {
	m.Render(CancelStatus)
}

func (m *BlocksBaseManager) SetWriter(out ConsoleWriter) {
	m.out = out
}

// Writer get console writer interface
func (m *BlocksBaseManager) Writer() ConsoleWriter {
	return m.out
}

// AddMirrorMode add mirror mode
func (m *BlocksBaseManager) AddMirrorMode(mode ...ConsoleBlocks) {
	for _, v := range mode {
		v.InitBlocks()
		if v.GetBuffer() == nil {
			m.children = append(m.children, v)
			continue
		}
		// ignore invalid blocks
		if m.major != nil {
			continue
		}
		// store major blocks

		m.major = v
		m.children = append(m.children, v)
	}
}

func (m *BlocksBaseManager) Refresh() {
	m.Render(NormalStatus)
}

func (m *BlocksBaseManager) IsInTask() bool {
	return m.inTask.Load()
}

// Event deal console key press
func (m *BlocksBaseManager) Event(key Key, in []byte) (exit bool) {
	m.inTask.Store(true)
	defer m.inTask.Store(false)
	var ctx pressContext
	ctx.key = key
	ctx.input = in
	if m.major != nil {
		ctx.buf = m.major.GetBuffer()
	}
	// debug.Println("block mgr. event: buf:", ctx.buf != nil, "major:", m.major != nil)
	ctx.out = m.out
	if m.eventBefore != nil && m.eventBefore(&ctx, key, in) {
		exit = true
	}
	for _, v := range m.children {
		if !v.Active() {
			continue
		}
		if v.OnEvent(&ctx, key, in) {
			exit = true
		}
	}
	if m.eventBehind != nil && m.eventBehind(&ctx, key, in) {
		exit = true
	}

	// control key press
	if key == m.cancelKey {
		// cancel key press
		if m.preCheck != nil && !m.preCheck(CancelStatus, ctx.buf) {
			m.Render(NormalStatus)
			return
		}

		m.Render(CancelStatus)

		if m.callback != nil && m.callback(CancelStatus, ctx.buf) {
			exit = true
			return
		}
		if !m.cancelNotExit {
			if m.major != nil && m.major.GetBuffer() != nil &&
				len(m.major.GetBuffer().Text()) == 0 {
				exit = true
				return
			}
		}
		if m.major != nil {
			m.major.ResetBuffer()
		}
	} else if key == m.finishKey {
		// finish key press

		if m.preCheck != nil && !m.preCheck(FinishStatus, ctx.buf) {
			m.Render(NormalStatus)
			return
		}

		m.Render(FinishStatus)
		if m.callback != nil && m.callback(FinishStatus, ctx.buf) {
			exit = true
			return
		}
		if m.major != nil {
			m.major.ResetBuffer()
		}
	} else {
		// normal check
		if m.preCheck != nil && !m.preCheck(NormalStatus, ctx.buf) {
			m.Render(NormalStatus)
			return
		}
		if m.callback != nil && m.callback(NormalStatus, ctx.buf) {
			exit = true
			return
		}
	}

	// // will change next mode
	// if m.changeStatus == 1 {
	// 	return
	// }

	// normal draw
	m.Render(NormalStatus)

	return
}

// OnEventBefore deal console key press
func (m *BlocksBaseManager) OnEventBefore(ctx PressContext, key Key, in []byte) (exit bool) {
	return
}

// OnEventBehind deal console key press
func (m *BlocksBaseManager) OnEventBehind(ctx PressContext, key Key, in []byte) (exit bool) {
	return
}

// UpdateWinSize called when window size is changed.
func (m *BlocksBaseManager) UpdateWinSize(size *WinSize) {
	m.row = size.Row
	m.col = size.Col
}

// ChangeStatus change status notify. revert internal status
func (m *BlocksBaseManager) ChangeStatus() {
	m.previousCursor = 0
	m.changeStatus = 0
}

// BlocksManager renders to the console.
func (m *BlocksBaseManager) Render(status int) {
	// In situations where a pseudo tty is allocated (e.g. within a docker container),
	// window size via TIOCGWINSZ is not immediately available and will result in 0,0 dimensions.
	if m.col == 0 {
		return
	}
	defer func() { debug.AssertNoError(m.out.Flush()) }()
	ctx := &consoleContext{
		BlocksBaseManager: m,
		buf:               nil,
		cursor:            -1,
		status:            status,
		prepare:           true,
	}
	// prepare draw
	prepare := func() (line int) {
		newCursor, lastCursor := 0, 0
		for _, item := range m.children {
			// ignore inactive
			if !item.Active() {
				continue
			}
			if !item.IsDraw(status) {
				continue
			}
			// render windows
			newCursor = item.Render(ctx, lastCursor)

			lastCursor = newCursor
		}
		w, line := m.ToPos(newCursor)
		// has current cursor line. so calc result need -1
		if w == 0 && line > 0 {
			line--
		}
		debug.Assert(line >= 0, "prepare line")
		return
	}()

	ctx.buf = nil
	ctx.cursor = -1
	ctx.BlocksBaseManager = m
	ctx.prepare = false

	// clean
	m.Move(m.previousCursor, 0)
	m.out.CursorBackward(m.col)
	m.out.EraseDown()
	m.PrepareArea(prepare)

	// cursor
	if status == NormalStatus {
		m.out.HideCursor()
	}
	// Rendering
	newCursor, lastCursor := 0, 0
	for _, item := range m.children {
		// ignore inactive
		if !item.Active() {
			continue
		}
		if !item.IsDraw(status) {
			continue
		}
		// render windows
		newCursor = item.Render(ctx, lastCursor)
		if ctx.cursor == -1 && ctx.buf != nil {
			// calc real cursor pos
			ctx.cursor = lastCursor + ctx.buf.Document().DisplayCursorPosition()
		}
		//		debug.Log(fmt.Sprintf("last:%d new:%d buf:%t cursor:%d", lastCursor, newCursor, ctx.buf != nil, ctx.cursor))

		lastCursor = newCursor
	}

	if status == NormalStatus {
		if ctx.cursor != -1 {
			// recover cursor pos
			m.Move(newCursor, ctx.cursor)
			// save pos
			m.previousCursor = ctx.cursor
			m.out.ShowCursor()
		} else {
			m.previousCursor = newCursor
		}

	} else {
		m.out.WriteStr("\n")

		col, _ := m.ToPos(newCursor)
		if col > 0 {
			m.out.CursorBackward(col)
		}

		m.previousCursor = 0
	}
}

// // BreakLine to break line.
// func (r *BlocksManager) BreakLine() {
// 	r.Render(CancelStatus)
// }

// Clear erases the screen from a beginning of input
// even if there is line break which means input length exceeds a window's width.
func (m *BlocksBaseManager) Clear() {
	m.Move(m.previousCursor, 0)
	m.out.EraseDown()
	m.out.Flush()
	m.previousCursor = 0
}

// Backward moves cursor to Backward from a current cursor position
// regardless there is a line break.
func (m *BlocksBaseManager) Backward(from, n int) int {
	return m.Move(from, from-n)
}

// Move moves cursor to specified position from the beginning of input
// even if there is a line break.
func (m *BlocksBaseManager) Move(from, to int) int {
	fromX, fromY := m.ToPos(from)
	toX, toY := m.ToPos(to)

	m.out.CursorUp(fromY - toY)
	m.out.CursorBackward(fromX - toX)
	return to
}

// ToPos returns the relative position from the beginning of the string.
func (m *BlocksBaseManager) ToPos(cursor int) (x, y int) {
	col := int(m.col)
	return cursor % col, cursor / col
}

// Columns return console window col count
func (m *BlocksBaseManager) Columns() int {
	return m.col
}

// Rows return console window row count
func (m *BlocksBaseManager) Rows() int {
	return m.row
}

// LineWrap line wrap
func (m *BlocksBaseManager) LineWrap(cursor int) {
	if runtime.GOOS != "windows" && cursor > 0 && cursor%int(m.col) == 0 {
		m.out.WriteRaw([]byte{'\n'})
	}
}

func (m *BlocksBaseManager) renderWindowTooSmall() {
	m.out.CursorGoTo(0, 0)
	m.out.EraseScreen()
	m.out.SetColor(DarkRed, White, false)
	m.out.WriteStr("Your console window is too small...")
}

// PrepareArea prepare enough area to display info
func (m *BlocksBaseManager) PrepareArea(lines int) {
	for i := 0; i < lines; i++ {
		m.out.ScrollDown()
	}
	for i := 0; i < lines; i++ {
		m.out.ScrollUp()
	}
}

