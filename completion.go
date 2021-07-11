package promptx

import (
	"fmt"

	"github.com/aggronmagi/promptx/internal/debug"
	runewidth "github.com/mattn/go-runewidth"
)

// Suggest is printed when completing.
type Suggest struct {
	Text        string
	Description string
}

// Completer should return the suggest item from Document.
type Completer func(in Document) []*Suggest

const (
	shortenSuffix = "..."
	leftPrefix    = " "
	leftSuffix    = " "
	rightPrefix   = " "
	rightSuffix   = " "
)

var (
	leftMargin       = runewidth.StringWidth(leftPrefix + leftSuffix)
	rightMargin      = runewidth.StringWidth(rightPrefix + rightSuffix)
	completionMargin = leftMargin + rightMargin
)

// CompletionManager manages which suggestion is now selected.
type CompletionManager struct {
	selected  int // -1 means nothing one is selected.
	tmp       []*Suggest
	max       uint16
	completer Completer

	verticalScroll int
	wordSeparator  string
	showAtStart    bool
}

// GetSelectedSuggestion returns the selected item.
func (c *CompletionManager) GetSelectedSuggestion() (s *Suggest) {
	if c.selected == -1 {
		if len(c.tmp) == 1 {
			return c.tmp[0]
		}
		return nil
	} else if c.selected < -1 {
		debug.Assert(false, "must not reach here")
		c.selected = -1
		return nil
	}

	if len(c.tmp) < 1 || c.selected > len(c.tmp) {
		return nil
	}

	return c.tmp[c.selected]
}

// GetSuggestions returns the list of suggestion.
func (c *CompletionManager) GetSuggestions() []*Suggest {
	return c.tmp
}

// Reset to select nothing.
func (c *CompletionManager) Reset() {
	c.selected = 0
	c.verticalScroll = 0
	// c.Update(*NewDocument())
	c.tmp = c.tmp[:0]
}

// Update to update the suggestions.
func (c *CompletionManager) Update(in Document) {
	c.tmp = c.completer(in)
	debug.Println("update compeltion", fmt.Sprintf("[%s]", in.Text), SugguestPrint(c.tmp))
	c.update()
}

// Previous to select the previous suggestion item.
func (c *CompletionManager) Previous() {
	if c.verticalScroll == c.selected && c.selected > 0 {
		c.verticalScroll--
	}
	c.selected--
	c.update()
}

// Next to select the next suggestion item.
func (c *CompletionManager) Next() {
	if c.verticalScroll+int(c.max)-1 == c.selected {
		c.verticalScroll++
	}
	c.selected++
	c.update()
}

func (c *CompletionManager) update() {
	max := int(c.max)
	if len(c.tmp) < max {
		max = len(c.tmp)
	}

	if c.selected >= len(c.tmp) {
		// c.Reset()
		c.selected = 0
		c.verticalScroll = 0
	} else if c.selected < 0 {
		c.selected = len(c.tmp) - 1
		c.verticalScroll = len(c.tmp) - max
	}
}

// NewCompletionManager returns initialized CompletionManager object.
func NewCompletionManager(completer Completer, max uint16) *CompletionManager {
	return &CompletionManager{
		selected:  -1,
		max:       max,
		completer: completer,

		verticalScroll: 0,
	}
}
