package promptx

import (
	"io"

	"github.com/aggronmagi/promptx/internal/debug"
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
		// cresult display
		"ResultText":      InputFinishTextFunc(defaultInputFinishText),
		"ResultTextColor": Color(Blue),
		"ResultBGColor":   Color(DefaultColor),
	}
}

// InputFinishTextFunc modify finish text display
type InputFinishTextFunc func(cc *InputOptions, status int, doc *Document) (words []*Word)

func defaultInputFinishText(cc *InputOptions, status int, doc *Document) (words []*Word) {

	if status == FinishStatus {
		words = append(words, &SuccessWord)
	} else {
		words = append(words, &FailureWord)
	}
	words = append(words, &Word{
		Text:      cc.PrefixText,
		TextColor: cc.PrefixTextColor,
		BGColor:   cc.PrefixBGColor,
		Bold:      false,
	})

	words = append(words, &Word{
		Text:      doc.Text,
		TextColor: cc.ResultTextColor,
		BGColor:   cc.ResultBGColor,
		Bold:      false,
	})

	return
}

type InputBlockManager struct {
	*BlocksBaseManager
	PreWords *BlocksWords
	Input    *BlocksEmacsBuffer
	Validate *BlocksNewLine
	cc       *InputOptions
}

// NewInputManager new input text
func NewInputManager(cc *InputOptions) (m *InputBlockManager) {
	cc.TipText = deleteBreakLineCharacters(cc.TipText)
	m = &InputBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		PreWords:          &BlocksWords{},
		Input:             &BlocksEmacsBuffer{},
		Validate:          &BlocksNewLine{},
		cc:                cc,
	}
	if len(cc.TipText) > 0 {
		m.PreWords.Words = append(m.PreWords.Words, &Word{
			Text:      cc.TipText,
			TextColor: cc.TipTextColor,
			BGColor:   cc.TipBGColor,
			Bold:      false,
		})
		m.PreWords.Words = append(m.PreWords.Words, &NewLineWord)
	}
	m.PreWords.Words = append(m.PreWords.Words, &AskWord)
	m.PreWords.Words = append(m.PreWords.Words, &Word{
		Text:      cc.PrefixText,
		TextColor: cc.PrefixTextColor,
		BGColor:   cc.PrefixBGColor,
		Bold:      false,
	})

	m.Validate.TextColor = cc.ValidTextColor
	m.Validate.BGColor = cc.ValidBGColor

	m.SetCancelKey(cc.CancelKey)
	m.SetFinishKey(cc.FinishKey)

	m.AddMirrorMode(m.PreWords)
	m.AddMirrorMode(m.Input)
	m.AddMirrorMode(m.Validate)

	m.SetBeforeEvent(m.BeforeEvent)

	m.SetCallBack(m.FinishCallBack)
	m.SetPreCheck(m.PreCheckCallBack)

	// plugin exit not exit
	m.SetCancelKeyAutoExit(false)
	return
}

// FinishCallBack  call back
func (m *InputBlockManager) FinishCallBack(status int, buf *Buffer) bool {

	if status == FinishStatus {
		// set not draw new line
		m.SetChangeStatus(1)
		if m.cc.FinishFunc != nil && buf != nil {
			text := buf.Document().Text
			//m.ExecTask(func() {
			m.cc.FinishFunc(text, nil)
			//})
		}
		debug.Println("recv input finish")
		return true
	}
	if status == CancelStatus {
		if m.cc.FinishFunc != nil {
			// m.ExecTask(func() {
			m.cc.FinishFunc("", io.EOF)
			//})
		}
		debug.Println("recv input cancel")
		// set not draw new line
		m.SetChangeStatus(1)
		return true
	}
	return false
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
	if !success {
		return
	}
	if status == FinishStatus || status == CancelStatus {
		// hide blocks
		m.Input.SetActive(false)
		m.Validate.SetActive(false)

		m.PreWords.Words = m.cc.ResultText(m.cc, status, buf.Document())
	}

	return
}

func (m *InputBlockManager) BeforeEvent(ctx PressContext, key Key, in []byte) (exit bool) {
	// first deal input char event
	if key == NotDefined && ctx.GetBuffer() != nil {
		ctx.GetBuffer().InsertText(string(in), false, true)
	}
	return
}
