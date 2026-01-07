package promptx

import (
	"github.com/aggronmagi/promptx/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	// strMultiChoice       = " ❯"
	strMultiChoiceCur    = " >"
	strMultiChoiceNot    = "  "
	strMultiChoiceSpacer = " "
	strMultiChoiceOpen   = "⬡ "
	strMultiChoiceSelect = "⬢ "
)

// BlocksSelect render select
type BlocksSelect struct {
	EmptyBlocks
	cc             *SelectOptions
	selected       int
	verticalScroll int
	selects        []int

	SelectFunc func(sels []int)
}

func (c *BlocksSelect) InitBlocks() {
	c.SetActive(true)

	c.BindKey(c.Previous, Up, Left, ControlP, BackTab)
	c.BindKey(c.Next, Down, Right, ControlN, Tab)
	c.BindASCII(c.Next, 'j', 'J', 'l', 'L')
	c.BindASCII(c.Previous, 'k', 'K', 'h', 'H')
	if c.cc.Multi {
		// 空格选中
		c.BindASCII(c.Select, ' ')
	} else {
		c.BindKey(c.Select, c.cc.Finish)
	}

	// c.BindKey(c.Cancel, c.cc.CancelKey)
}

func (c *BlocksSelect) Select(ctx PressContext) (exit bool) {
	debug.Println("select", c.selected, c.selects, c.verticalScroll)
	find := false
	for k, v := range c.selects {
		if v == c.selected {
			find = true
			c.selects = append(c.selects[:k], c.selects[k+1:]...)
			break
		}
	}
	if find {
		return
	}
	c.selects = append(c.selects, c.selected)
	debug.Println("select", c.selected, c.selects, c.verticalScroll)
	return
}

func (c *BlocksSelect) GetSelects() []int {
	return c.selects
}

// Previous to select the previous suggestion item.
func (c *BlocksSelect) Previous(ctx PressContext) (exit bool) {
	if c.verticalScroll == c.selected && c.selected > 0 {
		c.verticalScroll--
	}
	c.selected--
	c.update()
	return
}

// Next to select the next suggestion item.
func (c *BlocksSelect) Next(ctx PressContext) (exit bool) {
	if c.verticalScroll+int(c.cc.Rows)-1 == c.selected {
		c.verticalScroll++
	}
	c.selected++
	c.update()
	return
}

func (c *BlocksSelect) update() {
	max := int(c.cc.Rows)
	if len(c.cc.Options) < max {
		max = len(c.cc.Options)
	}

	if c.selected >= len(c.cc.Options) {
		c.selected = 0
		c.verticalScroll = 0
	} else if c.selected < 0 {
		c.selected = len(c.cc.Options) - 1
		c.verticalScroll = len(c.cc.Options) - max
	}
	if c.verticalScroll+max > len(c.cc.Options) {
		c.verticalScroll = len(c.cc.Options) - max
	}
}

func (c *BlocksSelect) isAlreadySelect(id int) bool {
	for _, v := range c.selects {
		if v == id {
			return true
		}
	}
	return false
}

// Render render to console
func (c *BlocksSelect) Render(ctx PrintContext, preCursor int) int {
	if c.cc == nil || len(c.cc.Options) == 0 {
		return preCursor
	}
	col := ctx.Columns()
	// first change line
	firstMoveLeft := 0
	if x, _ := ctx.ToPos(preCursor); x != 0 {
		firstMoveLeft = x
		preCursor += col - x
	}
	debug.Println("select pre move:", firstMoveLeft)

	// no complete sugguestions
	suggestions := c.cc.Options
	if len(suggestions) == 0 {
		return preCursor
	}

	windowHeight := len(suggestions)
	if windowHeight > c.cc.Rows {
		windowHeight = c.cc.Rows
	}

	if ctx.Prepare() {
		return preCursor + ctx.Columns()*(windowHeight)
	}
	prefixLen := runewidth.StringWidth(strMultiChoiceCur)
	if c.cc.Multi {
		prefixLen += runewidth.StringWidth(strMultiChoiceOpen)
	}
	formatted, width := formatSuggestions(
		suggestions[c.verticalScroll:c.verticalScroll+windowHeight],
		ctx.Columns()-prefixLen-1, // -1 means a width of scrollbar
	)
	// +1 means a width of scrollbar.
	width++

	// ctx.PrepareArea(windowHeight)
	out := ctx.Writer()
	if firstMoveLeft > 0 {
		// out.CursorDown(1)
	} else {
		out.CursorUp(1)
	}
	out.CursorBackward(col)

	cursor := 0

	contentHeight := len(suggestions)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(c.verticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := c.selected - c.verticalScroll

	out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		out.CursorDown(1)
		ctx.Backward(cursor+width, width+prefixLen)

		if i == selected {
			out.SetColor(c.cc.SelSuggestColor, c.cc.SelSuggestBG, true)
			out.WriteStr(strMultiChoiceCur)
		} else {
			out.SetColor(c.cc.SuggestColor, c.cc.SuggestBG, false)
			out.WriteStr(strMultiChoiceNot)
		}
		// 多选
		if c.cc.Multi {
			if c.isAlreadySelect(i + c.verticalScroll) {
				out.WriteStr(strMultiChoiceSelect)
			} else {
				out.WriteStr(strMultiChoiceOpen)
			}
		}

		out.WriteStr(formatted[i].Text)

		if i == selected {
			out.SetColor(c.cc.SelDescColor, c.cc.SelDescBG, false)
		} else {
			out.SetColor(c.cc.DescColor, c.cc.DescBG, false)
		}
		out.WriteStr(formatted[i].Description)

		if isScrollThumb(i) {
			out.SetColor(DefaultColor, c.cc.BarColor, false)
		} else {
			out.SetColor(DefaultColor, c.cc.BarBG, false)
		}
		out.WriteStr(" ")
		out.SetColor(DefaultColor, DefaultColor, false)

		// ctx.LineWrap(cursor + width)

	}

	return preCursor + ctx.Columns()*(windowHeight-1) + width + prefixLen
}
