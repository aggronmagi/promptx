package promptx

import (
	"fmt"
	"strings"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"
)

func deleteBreakLineCharacters(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func formatTexts(o []string, max int, prefix, suffix string) (new []string, width int) {
	l := len(o)
	n := make([]string, l)

	lenPrefix := runewidth.StringWidth(prefix)
	lenSuffix := runewidth.StringWidth(suffix)
	lenShorten := runewidth.StringWidth(shortenSuffix)
	min := lenPrefix + lenSuffix + lenShorten
	for i := 0; i < l; i++ {
		n[i] = deleteBreakLineCharacters(o[i])

		w := runewidth.StringWidth(n[i])
		if width < w {
			width = w
		}
	}

	if width == 0 {
		return make([]string, l), 0
	}
	if min >= max {
		return make([]string, l), 0
	}
	if lenPrefix+width+lenSuffix > max {
		width = max - lenPrefix - lenSuffix
	}

	for i := 0; i < l; i++ {
		x := runewidth.StringWidth(n[i])
		if x <= width {
			spaces := strings.Repeat(" ", width-x)
			n[i] = prefix + n[i] + spaces + suffix
		} else if x > width {
			x := runewidth.Truncate(n[i], width, shortenSuffix)
			// When calling runewidth.Truncate("您好xxx您好xxx", 11, "...") returns "您好xxx..."
			// But the length of this result is 10. So we need fill right using runewidth.FillRight.
			n[i] = prefix + runewidth.FillRight(x, width) + suffix
		}
	}
	return n, lenPrefix + width + lenSuffix
}

func formatSuggestions(suggests []*Suggest, max int) (new []*Suggest, width int) {
	num := len(suggests)
	new = make([]*Suggest, num)

	left := make([]string, num)
	for i := 0; i < num; i++ {
		left[i] = suggests[i].Text
	}
	right := make([]string, num)
	for i := 0; i < num; i++ {
		right[i] = suggests[i].Description
	}

	left, leftWidth := formatTexts(left, max, leftPrefix, leftSuffix)
	if leftWidth == 0 {
		return nil, 0
	}
	right, rightWidth := formatTexts(right, max-leftWidth, rightPrefix, rightSuffix)

	for i := 0; i < num; i++ {
		new[i] = &Suggest{Text: left[i], Description: right[i]}
	}
	return new, leftWidth + rightWidth
}

func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func HasPrefix(r, prefix []rune) bool {
	if len(r) < len(prefix) {
		return false
	}
	return Equal(r[:len(prefix)], prefix)
}

func TrimSpaceLeft(in []rune) []rune {
	firstIndex := len(in)
	for i, r := range in {
		if unicode.IsSpace(r) == false {
			firstIndex = i
			break
		}
	}
	return in[firstIndex:]
}

func TrimFirstSpace(in []rune) []rune {
	firstIndex := len(in)
	for i, r := range in {
		if unicode.IsSpace(r) == true {
			firstIndex = i
			break
		}
	}
	return in[:firstIndex]
}

type SugguestPrint []*Suggest

func (s SugguestPrint) String() string {
	str := fmt.Sprintf("[%d](", len(s))
	for k, v := range s {
		if k > 0 {
			str += ","
		}
		str += v.Text
	}
	str += ")"
	return str
}
