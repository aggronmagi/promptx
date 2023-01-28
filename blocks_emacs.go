package promptx

import (
	buffer "github.com/aggronmagi/promptx/buffer"
	"github.com/aggronmagi/promptx/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

// BlocksEmacsBuffer emacs operate like input ctx.GetBuffer()fer
//
// Moving the cursor
// -----------------
// | ok  | key       | description                                                  |
// |-----+-----------+--------------------------------------------------------------|
// | [x] | Ctrl + a  | Go to the beginning of the line (Home)                       |
// | [x] | Ctrl + e  | Go to the End of the line (End)                              |
// | [x] | Ctrl + p  | Previous command (Up arrow)                                  |
// | [x] | Ctrl + n  | Next command (Down arrow)                                    |
// | [x] | Ctrl + f  | Forward one character                                        |
// | [x] | Ctrl + b  | Backward one character                                       |
// | [x] | Meta + B  |                                                              |
// | [x] | Meta + F  |                                                              |

// Editing
// -------
// | ok  | key      | description                                             |
// |-----+----------+---------------------------------------------------------|
// | [x] | Ctrl + L | Clear the Screen, similar to the clear command          |
// | [x] | Ctrl + d | Delete character under the cursor                       |
// | [x] | Ctrl + h | Delete character before the cursor (Backspace)          |
// | [x] | Ctrl + w | Cut the Word before the cursor to the clipboard.        |
// | [x] | Ctrl + k | Cut the Line after the cursor to the clipboard.         |
// | [x] | Ctrl + u | Cut/delete the Line before the cursor to the clipboard. |
// | [ ] | Ctrl + t | Swap the last two characters before the cursor (typo).  |
// | [ ] | Esc  + t | Swap the last two words before the cursor.              |
// | [ ] | ctrl + y | Paste the last thing to be cut (yank)                   |
// | [ ] | ctrl + _ | Undo                                                    |
type BlocksEmacsBuffer struct {
	EmptyBlocks
	buf *Buffer
	// colors
	TextColor Color
	BGColor   Color
	// // Select
	// SelectTextColor Color
	// SelectBGColor   Color
	init bool
}

func (c *BlocksEmacsBuffer) ResetBuffer() {
	c.buf = buffer.NewBuffer()
	return
}

func (c *BlocksEmacsBuffer) GetBuffer() *Buffer {
	return c.buf
}

func (c *BlocksEmacsBuffer) InitBlocks() {
	if c.init {
		return
	}
	c.SetActive(true)
	c.ResetBuffer()
	for _, v := range emacsKeyBindings {
		c.BindKey(v.Fn, v.Key)
	}
	for _, v := range commonKeyBindings {
		c.BindKey(v.Fn, v.Key)
	}
	c.init = true
}

// Render render to console
func (c *BlocksEmacsBuffer) Render(ctx PrintContext, preCursor int) int {
	if c.buf == nil {
		return preCursor
	}
	if ctx.Prepare() {
		ctx.SetBuffer(c.buf)
		return runewidth.StringWidth(c.buf.Text()) + preCursor
	}
	// SetCtx.GetBuffer()fer to notify show cursor
	ctx.SetBuffer(c.buf)
	text := c.buf.Text()
	out := ctx.Writer()
	out.SetColor(c.TextColor, c.BGColor, false)
	out.WriteStr(text)
	out.SetColor(DefaultColor, DefaultColor, false)
	return runewidth.StringWidth(text) + preCursor
}

var emacsKeyBindings = []KeyBind{
	// Go to the End of the line
	{
		Key: ControlE,
		Fn: func(ctx PressContext) bool {
			x := []rune(ctx.GetBuffer().Document().TextAfterCursor())
			ctx.GetBuffer().CursorRight(len(x))
			return false
		},
	},
	// Go to the beginning of the line
	{
		Key: ControlA,
		Fn: func(ctx PressContext) bool {
			x := []rune(ctx.GetBuffer().Document().TextBeforeCursor())
			ctx.GetBuffer().CursorLeft(len(x))
			return false
		},
	},
	// Cut the Line after the cursor
	{
		Key: ControlK,
		Fn: func(ctx PressContext) bool {
			x := []rune(ctx.GetBuffer().Document().TextAfterCursor())
			ctx.GetBuffer().Delete(len(x))
			return false
		},
	},
	// Cut/delete the Line before the cursor
	{
		Key: ControlU,
		Fn: func(ctx PressContext) bool {
			x := []rune(ctx.GetBuffer().Document().TextBeforeCursor())
			ctx.GetBuffer().DeleteBeforeCursor(len(x))
			return false
		},
	},
	// Delete character under the cursor
	{
		Key: ControlD,
		Fn: func(ctx PressContext) bool {
			if ctx.GetBuffer().Text() != "" {
				ctx.GetBuffer().Delete(1)
			}else {
				// use control-d exit 
				return true 
			}
			return false
		},
	},
	// Backspace
	{
		Key: ControlH,
		Fn: func(ctx PressContext) bool {
			ctx.GetBuffer().DeleteBeforeCursor(1)
			return false
		},
	},
	// Left word arrow:
	{
		Key: MetaB,
		Fn: func(ctx PressContext) bool {
			buf := ctx.GetBuffer()
			doc := buf.Document()
			buf.CursorLeft(len([]rune(doc.GetWordBeforeCursorWithSpace())))
			return false
		},
	},
	// right word arrow:
	{
		Key: MetaF,
		Fn: func(ctx PressContext) bool {
			buf := ctx.GetBuffer()
			doc := buf.Document()
			buf.CursorRight(len([]rune(doc.GetWordAfterCursorWithSpace())))
			return false
		},
	},
	// Right allow: Forward one character
	{
		Key: ControlF,
		Fn: func(ctx PressContext) bool {
			ctx.GetBuffer().CursorRight(1)
			return false
		},
	},
	// Left allow: Backward one character
	{
		Key: ControlB,
		Fn: func(ctx PressContext) bool {
			ctx.GetBuffer().CursorLeft(1)
			return false
		},
	},
	// Cut the Word before the cursor.
	{
		Key: ControlW,
		Fn: func(ctx PressContext) bool {
			ctx.GetBuffer().DeleteBeforeCursor(len([]rune(
				ctx.GetBuffer().Document().GetWordBeforeCursorWithSpace(),
			)))
			return false
		},
	},
	// Clear the Screen, similar to the clear command
	{
		Key: ControlL,
		Fn: func(ctx PressContext) bool {
			out := ctx.Writer()
			out.EraseScreen()
			out.CursorGoTo(0, 0)
			debug.AssertNoError(out.Flush())
			return false
		},
	},
}

var commonKeyBindings = []KeyBind{
	// Go to the End of the line
	{
		Key: End,
		Fn: func(ctx PressContext) (exit bool) {
			buf := ctx.GetBuffer()
			x := []rune(buf.Document().TextAfterCursor())
			buf.CursorRight(len(x))
			return
		},
	},
	// Go to the beginning of the line
	{
		Key: Home,
		Fn: func(ctx PressContext) (exit bool) {
			buf := ctx.GetBuffer()
			x := []rune(buf.Document().TextBeforeCursor())
			buf.CursorLeft(len(x))
			return
		},
	},
	// Delete character under the cursor
	{
		Key: Delete,
		Fn: func(ctx PressContext) (exit bool) {
			buf := ctx.GetBuffer()
			buf.Delete(1)
			return
		},
	},
	// Backspace
	{
		Key: Backspace,
		Fn: func(ctx PressContext) (exit bool) {
			buf := ctx.GetBuffer()
			buf.DeleteBeforeCursor(1)
			return
		},
	},
	// Right allow: Forward one character
	{
		Key: Right,
		Fn: func(ctx PressContext) (exit bool) {
			buf := ctx.GetBuffer()
			buf.CursorRight(1)
			return
		},
	},
	// Left allow: Backward one character
	{
		Key: Left,
		Fn: func(ctx PressContext) (exit bool) {
			buf := ctx.GetBuffer()
			buf.CursorLeft(1)
			return
		},
	},
}
