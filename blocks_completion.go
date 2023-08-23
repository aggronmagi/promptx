package promptx

import (
	completion "github.com/aggronmagi/promptx/completion"
	"github.com/aggronmagi/promptx/internal/debug"
	"github.com/aggronmagi/promptx/stack"
	runewidth "github.com/mattn/go-runewidth"
)

// Completer should return the suggest item from Document.
type Completer func(in Document) []*Suggest

// CompleteOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n CompleteOption -f -o gen_options_complete.go
func promptxCompleteOptions() interface{} {
	return map[string]interface{}{
		"SuggestionTextColor":          Color(White),
		"SuggestionBGColor":            Color(Cyan),
		"SelectedSuggestionTextColor":  Color(Black),
		"SelectedSuggestionBGColor":    Color(Turquoise),
		"DescriptionTextColor":         Color(Black),
		"DescriptionBGColor":           Color(Turquoise),
		"SelectedDescriptionTextColor": Color(White),
		"SelectedDescriptionBGColor":   Color(Cyan),
		"ScrollbarThumbColor":          Color(DarkGray),
		"ScrollbarBGColor":             Color(Cyan),
		"Completer":                    Completer(nil),
		"CompleteMax":                  int(5),
		"CompletionFillSpace":          false,
	}
}

// BlocksCompletion completion blocks
type BlocksCompletion struct {
	EmptyBlocks
	// colors
	previewSuggestionTextColor Color
	previewSuggestionBGColor   Color

	Cfg *CompleteOptions

	Completions *CompletionManager

	init bool
}

func (c *BlocksCompletion) ApplyOptions(opt ...CompleteOption) {
	debug.Println("last:", c.Cfg.Completer != nil)
	c.Cfg.ApplyOption(opt...)
	debug.Println("new:", c.Cfg.Completer != nil)
	if c.Cfg.Completer != nil {
		c.Completions = completion.NewCompletionManager(c.Cfg.CompleteMax)
		c.SetActive(true)
		debug.Println("enable compeltion", stack.TakeStacktrace(1))
	} else {
		c.SetActive(false)
		debug.Println("disable compeltion", stack.TakeStacktrace(1))
	}
}

func (c *BlocksCompletion) InitBlocks() {
	if c.init {
		return
	}
	if c.Cfg == nil {
		c.Cfg = NewCompleteOptions()
	}
	c.ApplyOptions()

	c.BindKey(c.resetCompletion, Escape)
	c.BindKey(func(ctx PressContext) (exit bool) {
		comp := c.Completions
		if len(comp.GetSuggestions()) > 0 {
			comp.Previous()
		}
		return
	}, BackTab)
	c.BindKey(c.tabCompletion, Tab)
	c.BindKey(c.refreshCompletion, NotDefined)
	for _, v := range emacsKeyBindings {
		c.BindKey(c.refreshCompletion, v.Key)
	}
	for _, v := range commonKeyBindings {
		c.BindKey(c.refreshCompletion, v.Key)
	}
	c.BindKey(func(ctx PressContext) (exit bool) {
		if ctx.GetBuffer().Text() == "" {
			c.Completions.Reset()
		}
		return
	}, Backspace)
	c.init = true
}

// Render render to console
func (c *BlocksCompletion) Render(ctx PrintContext, preCursor int) int {
	// not input buffer. ignore
	if ctx.GetBuffer() == nil {
		return preCursor
	}
	// not completion
	if c.Completions == nil {
		return preCursor
	}

	// cancel and finsih status not print completion
	switch ctx.Status() {
	case CancelStatus,
		FinishStatus:
		return preCursor
	}

	completions := c.Completions

	// no complete sugguestions
	suggestions := completions.GetSuggestions()
	if len(suggestions) == 0 {
		return preCursor
	}

	windowHeight := len(suggestions)
	if windowHeight > int(completions.Max) {
		windowHeight = int(completions.Max)
	}

	if ctx.Prepare() {
		return preCursor + ctx.Columns()*(windowHeight+1) - preCursor%ctx.Columns()
	}

	prefix := ""
	formatted, width := formatSuggestions(
		suggestions[completions.VerticalScroll:completions.VerticalScroll+windowHeight],
		ctx.Columns()-runewidth.StringWidth(prefix)-1, // -1 means a width of scrollbar
	)

	if len(formatted) == 0 {
		debug.Println(SugguestPrint(suggestions), SugguestPrint(formatted))
		return preCursor
	}
	// +1 means a width of scrollbar.
	width++

	// ctx.PrepareArea(windowHeight)

	// real cursor postion
	cursor := ctx.InputCursor()
	// move to real cursor pos
	ctx.Move(preCursor, cursor)
	x, _ := ctx.ToPos(cursor)
	if x+width >= ctx.Columns() {
		cursor = ctx.Backward(cursor, x+width-ctx.Columns())
	}
	// pre cursor pos fix
	if cursor+windowHeight*ctx.Columns() > preCursor {
		preCursor = cursor + windowHeight*ctx.Columns()
	}

	contentHeight := len(suggestions)

	fractionVisible := float64(windowHeight) / float64(contentHeight)
	fractionAbove := float64(completions.VerticalScroll) / float64(contentHeight)

	scrollbarHeight := int(clamp(float64(windowHeight), 1, float64(windowHeight)*fractionVisible))
	scrollbarTop := int(float64(windowHeight) * fractionAbove)

	isScrollThumb := func(row int) bool {
		return scrollbarTop <= row && row <= scrollbarTop+scrollbarHeight
	}

	selected := completions.Selected - completions.VerticalScroll
	out := ctx.Writer()
	out.SetColor(White, Cyan, false)
	for i := 0; i < windowHeight; i++ {
		out.CursorDown(1)
		if i == selected {
			out.SetColor(c.Cfg.SelectedSuggestionTextColor, c.Cfg.SelectedSuggestionBGColor, true)
		} else {
			out.SetColor(c.Cfg.SuggestionTextColor, c.Cfg.SuggestionBGColor, false)
		}
		out.WriteStr(formatted[i].Text)

		if i == selected {
			out.SetColor(c.Cfg.SelectedDescriptionTextColor, c.Cfg.SelectedDescriptionBGColor, false)
		} else {
			out.SetColor(c.Cfg.DescriptionTextColor, c.Cfg.DescriptionBGColor, false)
		}
		out.WriteStr(formatted[i].Description)

		if isScrollThumb(i) {
			out.SetColor(DefaultColor, c.Cfg.ScrollbarThumbColor, false)
		} else {
			out.SetColor(DefaultColor, c.Cfg.ScrollbarBGColor, false)
		}
		out.WriteStr(" ")
		out.SetColor(DefaultColor, DefaultColor, false)

		ctx.LineWrap(cursor + width)
		ctx.Backward(cursor+width, width)
	}

	return preCursor
}

func (c *BlocksCompletion) resetCompletion(ctx PressContext) (exit bool) {
	if !c.Active() || c.Completions == nil {
		return
	}
	c.Completions.Reset()
	return
}

func (c *BlocksCompletion) refreshCompletion(ctx PressContext) (exit bool) {
	if !c.Active() || c.Completions == nil {
		return
	}
	c.Update(ctx.GetBuffer().Document())
	return
}

func (c *BlocksCompletion) tabCompletion(ctx PressContext) (exit bool) {
	if !c.Active() || c.Completions == nil {
		return
	}
	// check if need try update
	if len(c.Completions.GetSuggestions()) == 0 {
		c.Update(ctx.GetBuffer().Document())
	}
	suggestions := c.Completions.GetSuggestions()
	switch len(suggestions) {
	case 0:
		return
	case 1:
		buf := ctx.GetBuffer()
		w := buf.Document().GetWordBeforeCursorUntilSeparator(c.Completions.WordSeparator)
		if w != "" {
			buf.DeleteBeforeCursor(len([]rune(w)))
		}
		if c.Cfg.CompletionFillSpace {
			buf.InsertText(suggestions[0].Text+" ", false, true)
		} else {
			buf.InsertText(suggestions[0].Text, false, true)
		}

		c.Update(ctx.GetBuffer().Document())
	default:
		c.Completions.Next()
	}
	return
}

func (c *BlocksCompletion) EnterSelect(buf *Buffer) (ok bool) {
	if !c.Active() || c.Completions == nil || buf == nil {
		return
	}
	s := c.Completions.GetSelectedSuggestion()
	if s == nil {
		return
	}
	w := buf.Document().GetWordBeforeCursorUntilSeparator(c.Completions.WordSeparator)
	if w != "" {
		buf.DeleteBeforeCursor(len([]rune(w)))
	}
	if c.Cfg.CompletionFillSpace {
		buf.InsertText(s.Text+" ", false, true)
	} else {
		buf.InsertText(s.Text, false, true)
	}

	// c.Update(buf.Document())
	return true
}

func (c *BlocksCompletion) Update(doc *Document) {
	c.Completions.Update(c.Cfg.Completer(*doc))
}

func clamp(high, low, x float64) float64 {
	switch {
	case high < x:
		return high
	case x < low:
		return low
	default:
		return x
	}
}
