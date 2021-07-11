package promptx

import (
	runewidth "github.com/mattn/go-runewidth"
)

// BlocksPrefix render line prefix
type BlocksPrefix struct {
	EmptyBlocks
	// context
	Text string
	// colors
	TextColor Color
	BGColor   Color
}

// Render render to console
func (c *BlocksPrefix) Render(ctx PrintContext, preCursor int) int {
	if len(c.Text) == 0 {
		return preCursor
	}
	if ctx.Prepare() {
		return runewidth.StringWidth(c.Text) + preCursor
	}
	out := ctx.Writer()
	out.SetColor(c.TextColor, c.BGColor, false)
	out.WriteStr(c.Text)
	out.SetColor(DefaultColor, DefaultColor, false)
	return runewidth.StringWidth(c.Text) + preCursor
}

// BlocksSuffix render one line
type BlocksSuffix struct {
	EmptyBlocks
	// context
	Text string
	// colors
	TextColor Color
	BGColor   Color
}

// Render render to console
func (c *BlocksSuffix) Render(ctx PrintContext, preCursor int) int {
	if len(c.Text) == 0 {
		return preCursor
	}

	preCursor = runewidth.StringWidth(c.Text) + preCursor
	col := ctx.Columns()
	newCursor := preCursor + int(col) - preCursor%int(col)

	if ctx.Prepare() {
		return newCursor
	}
	out := ctx.Writer()
	out.SetColor(c.TextColor, c.BGColor, false)
	out.WriteStr(c.Text + "\n")
	out.SetColor(DefaultColor, DefaultColor, false)
	// backward cursor
	out.CursorBackward(col)
	return newCursor
}

// BlocksNewLine render one new line
type BlocksNewLine struct {
	// BlocksSuffix
	EmptyBlocks
	// context
	Text string
	// colors
	TextColor Color
	BGColor   Color
}

// Render render to console
func (c *BlocksNewLine) Render(ctx PrintContext, preCursor int) int {
	if len(c.Text) == 0 {
		return preCursor
	}
	col := ctx.Columns()
	// first change line
	firstMoveLeft := 0
	if x, _ := ctx.ToPos(preCursor); x != 0 {
		firstMoveLeft = x
		preCursor += col - x
	}
	// add new context line
	newCursor := runewidth.StringWidth(c.Text) + preCursor
	// newCursor := preCursor + col - preCursor%int(col)

	if ctx.Prepare() {
		return newCursor
	}
	out := ctx.Writer()
	_ = firstMoveLeft
	//if firstMoveLeft > 0 {
	out.CursorDown(1)
	out.CursorBackward(col)
	// }
	out.SetColor(c.TextColor, c.BGColor, false)
	out.WriteStr(c.Text)
	out.SetColor(DefaultColor, DefaultColor, false)
	// // backward cursor
	// out.CursorBackward(col)

	return newCursor
}

// BlocksStatusPrefix render line prefix with status
type BlocksStatusPrefix struct {
	EmptyBlocks
	// context
	Status   bool
	FailText string
	SucText  string
	// colors
	FailTextColor Color
	FailBGColor   Color
	SucTextColor  Color
	SucBGColor    Color
}

// Render render to console
func (c *BlocksStatusPrefix) Render(ctx PrintContext, preCursor int) int {
	text := c.FailText
	textColor := c.FailTextColor
	bgColor := c.FailBGColor
	if c.Status {
		text = c.SucText
		textColor = c.SucTextColor
		bgColor = c.SucBGColor
	}
	if len(text) == 0 {
		return preCursor
	}
	if ctx.Prepare() {
		return runewidth.StringWidth(text) + preCursor
	}
	out := ctx.Writer()
	out.SetColor(textColor, bgColor, false)
	out.WriteStr(text)
	out.SetColor(DefaultColor, DefaultColor, false)
	return runewidth.StringWidth(text) + preCursor
}
