package promptx

import (
	"io"
	"sync"
)

// InputOptionsOptionDeclareWithDefault promptx options
// generate by https://github.com/timestee/optiongen
//go:generate optionGen --option_with_struct_name=true --v=true
func InputOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TipText":         "",
		"TipTextColor":    Color(Yellow),
		"TipBGColor":      Color(DefaultColor),
		"PrefixText":      ">> ",
		"PrefixTextColor": Color(Green),
		"PrefixBGColor":   Color(DefaultColor),
		"ValidFunc":       (func(*Document) error)(nil),
		"ValidTextColor":  Color(Red),
		"ValidBGColor":    Color(DefaultColor),
		"FinishFunc":      (func(input string, eof error))(nil),
		"FinishKey":       Key(Enter),
		"CancelKey":       Key(ControlC),
	}
}

type InputBlockManager struct {
	*BlocksBaseManager
	Tip      *BlocksSuffix
	Prefix   *BlocksPrefix
	Input    *BlocksEmacsBuffer
	Validate *BlocksNewLine
	cc       *InputOptions
	cond     *sync.Cond
	m        sync.Mutex
}

// NewInputManager new input text
func NewInputManager(cc *InputOptions) (m *InputBlockManager) {
	m = &InputBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		Tip:               &BlocksSuffix{},
		Prefix:            &BlocksPrefix{},
		Input:             &BlocksEmacsBuffer{},
		Validate:          &BlocksNewLine{},
		cc:                cc,
	}
	m.Tip.Text = cc.TipText
	m.Tip.TextColor = cc.TipTextColor
	m.Tip.BGColor = cc.TipBGColor
	m.Prefix.Text = cc.PrefixText
	m.Prefix.TextColor = cc.PrefixTextColor
	m.Prefix.BGColor = cc.PrefixBGColor
	m.Validate.TextColor = cc.ValidTextColor
	m.Validate.BGColor = cc.ValidBGColor

	m.SetCancelKey(cc.CancelKey)
	m.SetFinishKey(cc.FinishKey)

	m.AddMirrorMode(m.Tip)
	m.AddMirrorMode(m.Prefix)
	m.AddMirrorMode(m.Input)
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
func (m *InputBlockManager) FinishCallBack(status int, buf *Buffer) {
	if status == FinishStatus {
		// set not draw new line
		m.SetChangeStatus(1)
		if m.cc.FinishFunc != nil && buf != nil {
			text := buf.Document().Text
			//m.ExecTask(func() {
			m.cc.FinishFunc(text, nil)
			//})
		}
		m.cond.Signal()
	}
	if status == CancelStatus {
		if m.cc.FinishFunc != nil {
			// m.ExecTask(func() {
			m.cc.FinishFunc("", io.EOF)
			//})
		}
		// set not draw new line
		m.SetChangeStatus(1)
		m.cond.Signal()
	}
}

// PreCheckCallBack change status pre check
func (m *InputBlockManager) PreCheckCallBack(status int, buf *Buffer) (success bool) {
	success = true
	switch status {
	case CancelStatus:
		m.Validate.Text = ""
	case FinishStatus, NormalStatus:
		// 检查输入
		if m.cc.ValidFunc != nil && buf != nil {
			if err := m.cc.ValidFunc(buf.Document()); err != nil {
				m.Validate.Text = err.Error()
				success = false
			} else {
				m.Validate.Text = ""
			}
		}
	}
	return
}
