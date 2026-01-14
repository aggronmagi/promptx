package blocks

import (
	"strings"

	"github.com/aggronmagi/promptx/v2/output"
	runewidth "github.com/mattn/go-runewidth"
)

// Word terminal displayed text
type Word struct {
	// content
	Text string
	// colors
	TextColor output.Color
	BGColor   output.Color
	// bold font
	Bold bool
}

// Width calc display pos
func (w *Word) Render(ctx PrintContext, preCursor int) (nextCursor int) {
	if len(w.Text) == 0 {
		return preCursor
	}
	nextCursor = preCursor
	list := strings.Split(w.Text, "\n")
	for k, v := range list {
		if k > 0 {
			nextCursor += ctx.Columns() - nextCursor%ctx.Columns()
		}
		nextCursor += runewidth.StringWidth(v)
	}
	if ctx.Prepare() {
		return
	}
	out := ctx.Writer()
	out.SetColor(w.TextColor, w.BGColor, w.Bold)
	out.WriteStr(w.Text)
	out.SetColor(output.DefaultColor, output.DefaultColor, false)
	// last one is '\n',backword text
	if len(list[len(list)-1]) == 0 {
		out.CursorBackward(ctx.Columns())
	}

	return
}

// WordDefault color text
func WordDefault(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.DefaultColor,
	}
}

// WordBlue color text
func WordBlue(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Blue,
	}
}

// WordBrown color text
func WordBrown(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Brown,
	}
}

// WordCyan color text
func WordCyan(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Cyan,
	}
}

// WordGreen color text
func WordGreen(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Green,
	}
}

// WordPurple color text
func WordPurple(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Purple,
	}
}

// WordRed color text
func WordRed(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Red,
	}
}

// WordTurquoise color text
func WordTurquoise(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Turquoise,
	}
}

// WordWhite color text
func WordWhite(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.White,
	}
}

// WordYellow color text
func WordYellow(str string) *Word {
	return &Word{
		Text:      str,
		TextColor: output.Yellow,
	}
}

// preset word for display select,input prefix word.
var (
	// SuccessWord success word
	SuccessWord = &Word{
		Text:      "✔ ",
		TextColor: output.Green,
		BGColor:   output.DefaultColor,
		Bold:      false,
	}
	// FailureWord failure word
	FailureWord = &Word{
		Text:      "✗ ",
		TextColor: output.Red,
		BGColor:   output.DefaultColor,
		Bold:      false,
	}
	AskWord = &Word{
		Text:      "? ",
		TextColor: output.Blue,
		BGColor:   output.DefaultColor,
		Bold:      false,
	}
	SelectWord = &Word{
		Text:      "▸ ",
		TextColor: output.DefaultColor,
		BGColor:   output.DefaultColor,
		Bold:      false,
	}
	NewLineWord = &Word{
		Text:      "\n",
		TextColor: output.DefaultColor,
		BGColor:   0,
		Bold:      false,
	}
)

// BlocksWords words display
type BlocksWords struct {
	EmptyBlocks
	Words []*Word
}

// Render render to console
func (c *BlocksWords) Render(ctx PrintContext, preCursor int) (nextCursor int) {
	if len(c.Words) == 0 {
		return preCursor
	}
	nextCursor = preCursor
	for _, v := range c.Words {
		nextCursor = v.Render(ctx, nextCursor)
	}
	return
}

// BlocksPrefix render line prefix
type BlocksPrefix struct {
	EmptyBlocks
	// context
	Text string
	// colors
	TextColor output.Color
	BGColor   output.Color
	Words     []*Word
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
	out.SetColor(output.DefaultColor, output.DefaultColor, false)
	return runewidth.StringWidth(c.Text) + preCursor
}

// BlocksSuffix render one line
type BlocksSuffix struct {
	EmptyBlocks
	// context
	Text string
	// colors
	TextColor output.Color
	BGColor   output.Color
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
	out.SetColor(output.DefaultColor, output.DefaultColor, false)
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
	TextColor output.Color
	BGColor   output.Color
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
	out.SetColor(output.DefaultColor, output.DefaultColor, false)
	// // backward cursor
	// out.CursorBackward(col)

	return newCursor
}
