package completion

import (
	"github.com/aggronmagi/promptx/v2/internal/debug"
)

// Suggest is printed when completing.
type Suggest struct {
	Text        string
	Description string
}

// CompletionManager manages which suggestion is now selected.
type CompletionManager struct {
	Selected int // -1 means nothing one is selected.
	Max      int
	tmp      []*Suggest

	VerticalScroll int
	WordSeparator  string
}

// GetSelectedSuggestion returns the selected item.
func (c *CompletionManager) GetSelectedSuggestion() (s *Suggest) {
	if c.Selected == -1 {
		if len(c.tmp) == 1 {
			return c.tmp[0]
		}
		return nil
	} else if c.Selected < -1 {
		debug.Assert(false, "must not reach here")
		c.Selected = -1
		return nil
	}

	if len(c.tmp) < 1 || c.Selected > len(c.tmp) {
		return nil
	}

	return c.tmp[c.Selected]
}

// GetSuggestions returns the list of suggestion.
func (c *CompletionManager) GetSuggestions() []*Suggest {
	return c.tmp
}

// Reset to select nothing.
func (c *CompletionManager) Reset() {
	c.Selected = 0
	c.VerticalScroll = 0
	// c.Update(*NewDocument())
	c.tmp = c.tmp[:0]
}

// Update to update the suggestions.
func (c *CompletionManager) Update(in []*Suggest) {
	c.tmp = in
	c.update()
}

// Previous to select the previous suggestion item.
func (c *CompletionManager) Previous() {
	if c.VerticalScroll == c.Selected && c.Selected > 0 {
		c.VerticalScroll--
	}
	c.Selected--
	c.update()
}

// Next to select the next suggestion item.
func (c *CompletionManager) Next() {
	if c.VerticalScroll+c.Max-1 == c.Selected {
		c.VerticalScroll++
	}
	c.Selected++
	c.update()
}

func (c *CompletionManager) update() {
	max := c.Max
	if len(c.tmp) < max {
		max = len(c.tmp)
	}

	if c.Selected >= len(c.tmp) {
		// c.Reset()
		c.Selected = 0
		c.VerticalScroll = 0
	} else if c.Selected < 0 {
		c.Selected = len(c.tmp) - 1
		c.VerticalScroll = len(c.tmp) - max
	}
}

// NewCompletionManager returns initialized CompletionManager object.
func NewCompletionManager(max int) *CompletionManager {
	return &CompletionManager{
		Selected:       -1,
		Max:            max,
		VerticalScroll: 0,
	}
}
