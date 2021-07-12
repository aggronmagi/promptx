package promptx

import (
	"sync"
)

// InputOptionsOptionDeclareWithDefault promptx options
// generate by https://github.com/timestee/optiongen
//go:generate optionGen --option_with_struct_name=true --v=true
func SelectOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Options":                      []*Suggest(nil),
		"RowsLimit":                    int(5),
		"FinishFunc":                   (func(sels []int))(nil),
		"Multi":                        false,
		"FinishKey":                    Key(Enter),
		"CancelKey":                    Key(ControlC),
		"TipText":                      "",
		"TipTextColor":                 Color(Yellow),
		"TipBGColor":                   Color(DefaultColor),
		"ValidFunc":                    (func(sels []int) error)(nil),
		"ValidTextColor":               Color(Red),
		"ValidBGColor":                 Color(DefaultColor),
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
	}
}

// "ColumnsLimit":                 int(1),

// SelectBlockManager select mode
type SelectBlockManager struct {
	*BlocksBaseManager
	Tip      *BlocksSuffix
	Select   *BlocksSelect
	Validate *BlocksNewLine
	cc       *SelectOptions
	cond     *sync.Cond
	m        sync.Mutex
}

// NewSelectManager new input text
func NewSelectManager(cc *SelectOptions) (m *SelectBlockManager) {
	m = &SelectBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		Tip:               &BlocksSuffix{},
		Select:            &BlocksSelect{},
		Validate:          &BlocksNewLine{},
		cc:                cc,
		m:                 sync.Mutex{},
	}
	m.Select.cc = cc
	m.Tip.Text = cc.TipText
	m.Tip.TextColor = cc.TipTextColor
	m.Tip.BGColor = cc.TipBGColor
	m.Validate.TextColor = cc.ValidTextColor
	m.Validate.BGColor = cc.ValidBGColor

	m.SetCancelKey(cc.CancelKey)
	m.SetFinishKey(cc.FinishKey)

	m.AddMirrorMode(m.Tip)
	m.AddMirrorMode(m.Select)
	m.AddMirrorMode(m.Validate)

	m.SetCallBack(m.FinishCallBack)
	m.SetPreCheck(m.PreCheckCallBack)

	m.cond = sync.NewCond(&m.m)
	m.m.Lock()
	// plugin exit not exit
	m.SetCancelKeyAutoExit(false)
	return
}

// FinishCallBack  call back
func (m *SelectBlockManager) FinishCallBack(status int, buf *Buffer) {
	if status == FinishStatus {
		// set not draw new line
		m.SetChangeStatus(1)

		if m.cc.FinishFunc != nil {
			//m.ExecTask(func() {
			m.cc.FinishFunc(m.Select.GetSelects())
			//})
		}
		m.cond.Signal()
	}
	if status == CancelStatus {
		// set not draw new line
		m.SetChangeStatus(1)

		if m.cc.FinishFunc != nil {
			//m.ExecTask(func() {
			m.cc.FinishFunc(nil)
			//})
		}

		m.cond.Signal()
	}
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
	return
}
