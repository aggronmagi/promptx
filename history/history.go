package history

import (
	"bytes"
	"os"
	"strings"

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
}

// Add to add text in history.
func (h *History) Add(input string) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return
	}
	if len(h.histories) < 1 || h.histories[len(h.histories)-1] != input {
		h.histories = append(h.histories, input)
	}
	h.buf = ""
	h.Rebuild("", true)
}

func (h *History) Remove(input string) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return
	}
	for k := range h.histories {
		if h.histories[k] == input {
			h.histories = append(h.histories[:k], h.histories[k+1:]...)
			break
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
		copy(h.tmp, h.histories)

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
		if strings.HasPrefix(v, buf) {
			h.tmp = append(h.tmp, v)
		}
	}
	// // if not match any one histories,put all histories to tmp.
	// if len(h.tmp) < 1 {
	// 	h.tmp = append(h.tmp, h.histories...)
	// }
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

// Load read persistence data from file
func (h *History) Load(file string) (err error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}
	h.histories = strings.Split(string(data), "\n")
	h.Rebuild("", true)
	return
}

func (h *History) Reset() {
	h.histories = make([]string, 0, 128)
	h.Rebuild("", true)
}

// NewHistory returns new history object.
func NewHistory() *History {
	return &History{
		histories: make([]string, 0, 128),
		tmp:       []string{""},
		selected:  0,
	}
}
