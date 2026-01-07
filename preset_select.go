package promptx

import (
	"strconv"
	"strings"
)

// SelectOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n SelectOption -f -o gen_options_select.go
func promptxSelectOptions() interface{} {
	return map[string]interface{}{
		"Options":   []*Suggest(nil),
		"Rows":      int(5),
		"OnFinish":  (func(sels []int))(nil),
		"Multi":     false,
		"Finish":    Key(Enter),
		"Cancel":    Key(ControlC),
		"Tip":       "",
		"TipColor":  Color(Yellow),
		"TipBG":     Color(DefaultColor),
		"ShowHelp":  false,
		"Help":      (SelHelpTextFunc)(defaultSelHelpText),
		"HelpColor": Color(DefaultColor),
		"HelpBG":    Color(DefaultColor),
		"Valid":     (func(sels []int) error)(nil),
		"ValidColor": Color(Red),
		"ValidBG":   Color(DefaultColor),
		"SuggestColor":    Color(White),
		"SuggestBG":       Color(Cyan),
		"SelSuggestColor": Color(Black),
		"SelSuggestBG":    Color(Turquoise),
		"DescColor":       Color(Black),
		"DescBG":          Color(Turquoise),
		"SelDescColor":    Color(White),
		"SelDescBG":       Color(Cyan),
		"BarColor":        Color(DarkGray),
		"BarBG":           Color(Cyan),
		"FinishText":      SelFinishTextFunc(defaultSelFinishText),
		"ShowItem":        true,
		"ResultColor":     Color(Blue),
		"ResultBG":        Color(DefaultColor),
		"Defaults":        []int(nil),
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
		Text:      cc.Tip + " ",
		TextColor: cc.TipColor,
		BGColor:   cc.TipBG,
	})
	opts := make([]string, 0, len(result)+1)
	if cc.ShowItem {
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
		TextColor: cc.ResultColor,
		BGColor:   cc.ResultBG,
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
	cc.Tip = deleteBreakLineCharacters(cc.Tip)
	m = &SelectBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		PreWords:          &BlocksWords{},
		Select:            &BlocksSelect{},
		Validate:          &BlocksNewLine{},
		cc:                cc,
	}
	m.Select.cc = cc
	m.Validate.TextColor = cc.ValidColor
	m.Validate.BGColor = cc.ValidBG

	m.PreWords.Words = append(m.PreWords.Words, &SelectWord)
	if len(cc.Tip) > 0 {
		m.PreWords.Words = append(m.PreWords.Words, &Word{
			Text:      cc.Tip,
			TextColor: cc.TipColor,
			BGColor:   cc.TipBG,
			Bold:      false,
		})
	}

	if cc.ShowHelp && cc.Help != nil {
		help := cc.Help(cc.Multi)
		if len(help) > 0 {
			m.PreWords.Words = append(m.PreWords.Words, &Word{
				Text:      help,
				TextColor: cc.HelpColor,
				BGColor:   cc.HelpBG,
				Bold:      false,
			})
		}
	}
	if len(m.PreWords.Words) == 1 {
		m.PreWords.Words = m.PreWords.Words[:0]
	} else {
		m.PreWords.Words = append(m.PreWords.Words, &NewLineWord)
	}

	m.SetCancelKey(cc.Cancel)
	m.SetFinishKey(cc.Finish)

	m.AddMirrorMode(m.PreWords)
	m.AddMirrorMode(m.Select)
	m.AddMirrorMode(m.Validate)

	m.SetCallBack(m.FinishCallBack)
	m.SetPreCheck(m.PreCheckCallBack)

	// plugin exit not exit
	m.SetCancelKeyAutoExit(false)
	// default select
	for len(m.cc.Defaults) > 0 {
		validSels := make([]int, 0, len(m.cc.Options))
		// 只传递了-1,那么就全选
		if len(m.cc.Defaults) == 1 && m.cc.Defaults[0] == -1 {
			for k := 0; k < len(m.cc.Options); k++ {
				validSels = append(validSels, k)
			}
		} else {
			// 传递多个参数. 检测有效性
			for _, k := range m.cc.Defaults {
				if k < 0 || k >= len(m.cc.Options) {
					continue
				}
				validSels = append(validSels, k)
			}
		}
		m.cc.Defaults = validSels
		if len(validSels) < 1 {
			break
		}
		if m.cc.Multi {
			m.Select.selects = m.cc.Defaults
		} else {
			m.Select.selected = m.cc.Defaults[0]
		}
		break
	}
	return
}

// FinishCallBack  call back
func (m *SelectBlockManager) FinishCallBack(status int, buf *Buffer) bool {
	if status == FinishStatus {
		// set not draw new line
		m.SetChangeStatus(1)

		if m.cc.OnFinish != nil {
			//m.ExecTask(func() {
			m.cc.OnFinish(m.Select.GetSelects())
			//})
		}
		return true
	}
	if status == CancelStatus {
		// set not draw new line
		m.SetChangeStatus(1)

		if m.cc.OnFinish != nil {
			//m.ExecTask(func() {
			m.cc.OnFinish(nil)
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
		if m.cc.Valid != nil {
			if err := m.cc.Valid(m.Select.GetSelects()); err != nil {
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

// TearDown to clear title and erasing.
func (m *SelectBlockManager) TearDown() {
	//m.BlocksBaseManager.Render(NormalStatus)
	m.BlocksBaseManager.TearDown()
	m.BlocksBaseManager.out.ShowCursor()
	m.out.Flush()

}
