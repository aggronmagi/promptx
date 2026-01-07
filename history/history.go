package history

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aggronmagi/promptx/internal/debug"
)

// History stores the texts that are entered.
type History struct {
	// all history
	histories []string
	// history suggestion
	tmp      []string
	selected int
	buf      string

	// config options
	maxSize    int
	ignoreDups bool
	dedup      bool
	timestamp  bool
}

// HistoryOption configures the history.
type HistoryOption func(*History)

// WithMaxSize sets the maximum number of history records.
func WithMaxSize(size int) HistoryOption {
	return func(h *History) {
		h.maxSize = size
	}
}

// WithIgnoreDups ignores consecutive duplicate commands.
func WithIgnoreDups(ignore bool) HistoryOption {
	return func(h *History) {
		h.ignoreDups = ignore
	}
}

// WithDeduplicate removes all previous occurrences of the same command.
func WithDeduplicate(dedup bool) HistoryOption {
	return func(h *History) {
		h.dedup = dedup
	}
}

// WithTimestamp enables recording timestamps for each command.
func WithTimestamp(enable bool) HistoryOption {
	return func(h *History) {
		h.timestamp = enable
	}
}

// ApplyOptions configures the history with the given options.
func (h *History) ApplyOptions(opts ...HistoryOption) {
	for _, opt := range opts {
		opt(h)
	}
}

// Add to add text in history.
func (h *History) Add(input string) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return
	}

	// Check if we should ignore this input
	if h.ignoreDups && len(h.histories) > 0 {
		last := h.histories[len(h.histories)-1]
		if h.timestamp {
			last = h.extractCommand(last)
		}
		if last == input {
			return
		}
	}

	// Global deduplication
	if h.dedup {
		h.Remove(input)
	} else if len(h.histories) > 0 {
		// Default behavior: ignore consecutive duplicates
		last := h.histories[len(h.histories)-1]
		if h.timestamp {
			last = h.extractCommand(last)
		}
		if last == input {
			return
		}
	}

	item := input
	if h.timestamp {
		item = fmt.Sprintf(": %d:0;%s", time.Now().Unix(), input)
	}

	h.histories = append(h.histories, item)

	// Limit size
	if h.maxSize > 0 && len(h.histories) > h.maxSize {
		h.histories = h.histories[len(h.histories)-h.maxSize:]
	}

	h.buf = ""
	h.Rebuild("", true)
}

func (h *History) extractCommand(item string) string {
	if !strings.HasPrefix(item, ": ") {
		return item
	}
	idx := strings.Index(item, ";")
	if idx < 0 {
		return item
	}
	return item[idx+1:]
}

func (h *History) Remove(input string) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return
	}
	for i := 0; i < len(h.histories); i++ {
		cmd := h.histories[i]
		if h.timestamp {
			cmd = h.extractCommand(cmd)
		}
		if cmd == input {
			h.histories = append(h.histories[:i], h.histories[i+1:]...)
			i--
		}
	}
	h.buf = ""
	h.Rebuild("", true)
}

// rebulid tmp with buf prefix.
func (h *History) Rebuild(buf string, force bool) {
	buf = strings.TrimSpace(buf)
	debug.Println("rebuild-buf", buf)
	// add all history
	if force || (len(buf) == 0 && len(h.tmp) != len(h.histories)+1) {
		h.tmp = make([]string, len(h.histories)+1)
		for i, v := range h.histories {
			if h.timestamp {
				h.tmp[i] = h.extractCommand(v)
			} else {
				h.tmp[i] = v
			}
		}

		h.selected = len(h.tmp) - 1
		h.buf = buf
		return
	}
	if h.buf == buf {
		return
	}

	if cap(h.tmp) < len(h.histories)+1 {
		h.tmp = make([]string, 0, len(h.histories)+1)
	} else {
		h.tmp = h.tmp[:0]
	}

	for _, v := range h.histories {
		cmd := v
		if h.timestamp {
			cmd = h.extractCommand(v)
		}
		if strings.HasPrefix(cmd, buf) {
			h.tmp = append(h.tmp, cmd)
		}
	}
	h.tmp = append(h.tmp, "")
	h.selected = len(h.tmp) - 1
	h.buf = buf
	debug.Println("rebuild-update", h.buf, h.tmp)
}

// Older saves a buffer of current line and get a buffer of previous line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Older(buf string) (new string, changed bool) {
	if len(h.tmp) == 1 || h.selected == 0 {
		return buf, false
	}
	h.tmp[h.selected] = buf

	h.selected--
	new = h.tmp[h.selected]
	return new, true
}

// Newer saves a buffer of current line and get a buffer of next line by up-arrow.
// The changes of line buffers are stored until new history is created.
func (h *History) Newer(buf string) (new string, changed bool) {
	if h.selected >= len(h.tmp)-1 {
		return buf, false
	}
	h.tmp[h.selected] = buf

	h.selected++
	new = h.tmp[h.selected]
	return new, true
}

// Save save data persistence to file
func (h *History) Save(file string) (err error) {
	buf := &bytes.Buffer{}
	for _, v := range h.histories {
		buf.WriteString(v)
		buf.WriteByte('\n')
	}
	err = os.WriteFile(file, buf.Bytes(), 0644)
	return
}

// AppendToFile appends the last history item to the file.
func (h *History) AppendToFile(file string) error {
	if len(h.histories) == 0 {
		return nil
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	last := h.histories[len(h.histories)-1]
	_, err = f.WriteString(last + "\n")
	return err
}

// Load read persistence data from file
func (h *History) Load(file string) (err error) {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}

	lines := strings.Split(string(data), "\n")
	h.histories = h.histories[:0]

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Auto detect timestamp mode if not explicitly set
		if !h.timestamp && strings.HasPrefix(line, ": ") {
			h.timestamp = true
		}
		h.histories = append(h.histories, line)
	}

	// Deduplicate if enabled
	if h.dedup {
		h.Deduplicate()
	}

	// Trim if needed
	h.Trim()

	h.Rebuild("", true)
	return
}

func (h *History) Reset() {
	h.histories = make([]string, 0, 128)
	h.Rebuild("", true)
}

// Trim restricts the history size to maxSize.
func (h *History) Trim() {
	if h.maxSize > 0 && len(h.histories) > h.maxSize {
		h.histories = h.histories[len(h.histories)-h.maxSize:]
	}
}

// Deduplicate removes all but the latest occurrence of each command.
func (h *History) Deduplicate() {
	seen := make(map[string]bool)
	newHistories := make([]string, 0, len(h.histories))
	for i := len(h.histories) - 1; i >= 0; i-- {
		cmd := h.histories[i]
		orig := cmd
		if h.timestamp {
			cmd = h.extractCommand(cmd)
		}
		if !seen[cmd] {
			seen[cmd] = true
			newHistories = append(newHistories, orig)
		}
	}
	// Reverse to maintain order
	for i, j := 0, len(newHistories)-1; i < j; i, j = i+1, j-1 {
		newHistories[i], newHistories[j] = newHistories[j], newHistories[i]
	}
	h.histories = newHistories
}

// GetWithTimestamp returns the history with timestamps if enabled.
func (h *History) GetWithTimestamp() []string {
	return h.histories
}

// GetCommands returns only the commands without timestamps.
func (h *History) GetCommands() []string {
	cmds := make([]string, len(h.histories))
	for i, v := range h.histories {
		if h.timestamp {
			cmds[i] = h.extractCommand(v)
		} else {
			cmds[i] = v
		}
	}
	return cmds
}

// NewHistory returns new history object.
func NewHistory(opts ...HistoryOption) *History {
	h := &History{
		histories: make([]string, 0, 128),
		tmp:       []string{""},
		selected:  0,
		maxSize:   10000, // Default oh-my-zsh like limit
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
