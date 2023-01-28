package promptx

import (
	"strconv"
	"strings"
)

// InputOptionsOptionDeclareWithDefault promptx options
// generate by https://github.com/aggronmagi/gogen/
//go:generate gogen option -n SelectOption -f -o gen_options_select.go
func SelectOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Options":    []*Suggest(nil),
		"RowsLimit":  int(5),
		"FinishFunc": (func(sels []int))(nil),
		"Multi":      false,
		"FinishKey":  Key(Enter),
		"CancelKey":  Key(ControlC),
		// select tip info
		"TipText":      "",
		"TipTextColor": Color(Yellow),
		"TipBGColor":   Color(DefaultColor),
		// help text info
		"ShowHelpText":  false,
		"HelpText":      (SelHelpTextFunc)(defaultSelHelpText),
		"HelpTextColor": Color(DefaultColor),
		"HelpBGColor":   Color(DefaultColor),
		// valid info
		"ValidFunc":      (func(sels []int) error)(nil),
		"ValidTextColor": Color(Red),
		"ValidBGColor":   Color(DefaultColor),
		// select options info
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
		// finish text show
		"FinishText": SelFinishTextFunc(defaultSelFinishText),
		// contrl selct result display select context
		"ResultShowItem":  true,
		"ResultTextColor": Color(Blue),
		"ResultBGColor":   Color(DefaultColor),
	}
}

// SelHelpTextFunc modify help text func
type SelHelpTextFunc func(mul bool) (help string)

func defaultSelHelpText(mul bool) (help string) {
	if mul {
		help = "use [j/k C-n/C-p] move. space select. enter finish"
	} else {
		help = "use [j/k C-n/C-p] move. press enter select."
	}
	return
}

// SelFinishTextFunc modify finish text display
type SelFinishTextFunc func(cc *SelectOptions, result []int) (words []*Word)

func defaultSelFinishText(cc *SelectOptions, result []int) (words []*Word) {
	if len(result) > 0 {
		words = append(words, &SuccessWord)
	} else {
		words = append(words, &FailureWord)
	}
	words = append(words, &Word{
		Text:      cc.TipText + " ",
		TextColor: cc.TipTextColor,
		BGColor:   cc.TipBGColor,
	})
	opts := make([]string, 0, len(result)+1)
	if cc.ResultShowItem {
		for _, v := range result {
			opts = append(opts, cc.Options[v].Text)
		}
	} else {
		for _, v := range result {
			opts = append(opts, strconv.Itoa(v))
		}
	}
	if len(opts) < 1 {
		opts = append(opts, "not select")
	}

	words = append(words, &Word{
		Text:      strings.Join(opts, ","),
		TextColor: cc.ResultTextColor,
		BGColor:   cc.ResultBGColor,
	})
	return
}

// SelectBlockManager select mode
type SelectBlockManager struct {
	*BlocksBaseManager
	PreWords *BlocksWords
	Select   *BlocksSelect
	Validate *BlocksNewLine
	cc       *SelectOptions
}

// NewSelectManager new input text
func NewSelectManager(cc *SelectOptions) (m *SelectBlockManager) {
	cc.TipText = deleteBreakLineCharacters(cc.TipText)
	m = &SelectBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		PreWords:          &BlocksWords{},
		Select:            &BlocksSelect{},
		Validate:          &BlocksNewLine{},
		cc:                cc,
	}
	m.Select.cc = cc
	m.Validate.TextColor = cc.ValidTextColor
	m.Validate.BGColor = cc.ValidBGColor

	m.PreWords.Words = append(m.PreWords.Words, &SelectWord)
	if len(cc.TipText) > 0 {
		m.PreWords.Words = append(m.PreWords.Words, &Word{
			Text:      cc.TipText,
			TextColor: cc.TipTextColor,
			BGColor:   cc.TipBGColor,
			Bold:      false,
		})
	}

	if cc.ShowHelpText && cc.HelpText != nil {
		help := cc.HelpText(cc.Multi)
		if len(help) > 0 {
			m.PreWords.Words = append(m.PreWords.Words, &Word{
				Text:      help,
				TextColor: cc.HelpTextColor,
				BGColor:   cc.HelpBGColor,
				Bold:      false,
			})
		}
	}
	if len(m.PreWords.Words) == 1 {
		m.PreWords.Words = m.PreWords.Words[:0]
	} else {
		m.PreWords.Words = append(m.PreWords.Words, &NewLineWord)
	}

	m.SetCancelKey(cc.CancelKey)
	m.SetFinishKey(cc.FinishKey)

	m.AddMirrorMode(m.PreWords)
	m.AddMirrorMode(m.Select)
	m.AddMirrorMode(m.Validate)

	m.SetCallBack(m.FinishCallBack)
	m.SetPreCheck(m.PreCheckCallBack)

	// plugin exit not exit
	m.SetCancelKeyAutoExit(false)
	return
}

// FinishCallBack  call back
func (m *SelectBlockManager) FinishCallBack(status int, buf *Buffer) bool {
	if status == FinishStatus {
		// set not draw new line
		m.SetChangeStatus(1)

		if m.cc.FinishFunc != nil {
			//m.ExecTask(func() {
			m.cc.FinishFunc(m.Select.GetSelects())
			//})
		}
		return true
	}
	if status == CancelStatus {
		// set not draw new line
		m.SetChangeStatus(1)

		if m.cc.FinishFunc != nil {
			//m.ExecTask(func() {
			m.cc.FinishFunc(nil)
			//})
		}

		return true
	}
	return false
}

// PreCheckCallBack change status pre check
func (m *SelectBlockManager) PreCheckCallBack(status int, buf *Buffer) (success bool) {
	success = true
	switch status {
	case CancelStatus:
		m.Validate.Text = ""
	case FinishStatus, NormalStatus:
		// 检查输入
		if m.cc.ValidFunc != nil {
			if err := m.cc.ValidFunc(m.Select.GetSelects()); err != nil {
				m.Validate.Text = err.Error()
				success = false
			} else {
				m.Validate.Text = ""
			}
		}
	}
	if !success {
		return
	}
	if status == FinishStatus || status == CancelStatus {
		// hide blocks
		m.Select.SetActive(false)
		m.Validate.SetActive(false)
		// modify finish display text
		m.PreWords.Words = m.cc.FinishText(m.cc, m.Select.GetSelects())
	}
	return
}
