package blocks

import (
	"fmt"
	"io"

	buffer "github.com/aggronmagi/promptx/v2/buffer"
	"github.com/aggronmagi/promptx/v2/input"
	"github.com/aggronmagi/promptx/v2/internal/debug"
	"github.com/aggronmagi/promptx/v2/output"
)

// InputOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n InputOption -f -o gen_options_input.go
func promptxInputOptions() interface{} {
	return map[string]interface{}{
		"Tip":         "",
		"TipColor":    output.Color(output.Yellow),
		"TipBG":       output.Color(output.DefaultColor),
		"Prefix":      ">> ",
		"PrefixColor": output.Color(output.Green),
		"PrefixBG":    output.Color(output.DefaultColor),
		"Valid":       (func(*buffer.Document) error)(nil),
		"ValidColor":  output.Color(output.Red),
		"ValidBG":     output.Color(output.DefaultColor),
		"OnFinish":    (func(input string, eof error))(nil),
		"Finish":      input.Key(input.Enter),
		"Cancel":      input.Key(input.ControlC),
		// result display
		"ResultText":   InputFinishTextFunc(defaultInputFinishText),
		"ResultColor":  output.Color(output.Blue),
		"ResultBG":     output.Color(output.DefaultColor),
		"Default":      "",
		"DefaultColor": output.Color(output.Brown),
		"DefaultBG":    output.Color(output.DefaultColor),
	}
}

// InputFinishTextFunc modify finish text display
type InputFinishTextFunc func(cc *InputOptions, status int, doc *buffer.Document, defaultText string) (words []*Word)

func defaultInputFinishText(cc *InputOptions, status int, doc *buffer.Document, defaultText string) (words []*Word) {

	if status == FinishStatus {
		words = append(words, SuccessWord)
	} else {
		words = append(words, FailureWord)
	}
	words = append(words, &Word{
		Text:      cc.Prefix,
		TextColor: cc.PrefixColor,
		BGColor:   cc.PrefixBG,
		Bold:      false,
	})

	if doc.Text != "" {
		defaultText = doc.Text
	}

	words = append(words, &Word{
		Text:      defaultText,
		TextColor: cc.ResultColor,
		BGColor:   cc.ResultBG,
		Bold:      false,
	})

	return
}

type InputBlockManager struct {
	*BlocksBaseManager
	PreWords   *BlocksWords
	Input      *BlocksEmacsBuffer
	Validate   *BlocksNewLine
	cc         *InputOptions
	useDefault bool
}

// NewInputManager new input text
func NewInputManager(cc *InputOptions) (m *InputBlockManager) {
	cc.Tip = deleteBreakLineCharacters(cc.Tip)
	m = &InputBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		PreWords:          &BlocksWords{},
		Input:             &BlocksEmacsBuffer{},
		Validate:          &BlocksNewLine{},
		cc:                cc,
	}
	if len(cc.Tip) > 0 {
		m.PreWords.Words = append(m.PreWords.Words, &Word{
			Text:      cc.Tip,
			TextColor: cc.TipColor,
			BGColor:   cc.TipBG,
			Bold:      false,
		})
		m.PreWords.Words = append(m.PreWords.Words, NewLineWord)
	}
	m.PreWords.Words = append(m.PreWords.Words, AskWord)
	m.PreWords.Words = append(m.PreWords.Words, &Word{
		Text:      cc.Prefix,
		TextColor: cc.PrefixColor,
		BGColor:   cc.PrefixBG,
		Bold:      false,
	})
	// 检测默认值是否有效
	if m.cc.Default != "" {
		doc := &buffer.Document{}
		doc.Text = m.cc.Default
		if err := m.cc.Valid(doc); err != nil {
			m.cc.Default = ""
		}
	}
	// 开启默认值显示
	if m.cc.Default != "" {
		m.PreWords.Words = append(m.PreWords.Words, &Word{
			Text:      fmt.Sprintf("[%s]", m.cc.Default),
			TextColor: m.cc.DefaultColor,
			BGColor:   m.cc.DefaultBG,
			Bold:      true,
		})
	}

	m.Validate.TextColor = cc.ValidColor
	m.Validate.BGColor = cc.ValidBG

	m.SetCancelKey(cc.Cancel)
	m.SetFinishKey(cc.Finish)

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
func (m *InputBlockManager) FinishCallBack(status int, buf *buffer.Buffer) bool {

	if status == FinishStatus {
		// set not draw new line
		m.SetChangeStatus(1)
		if m.cc.OnFinish != nil && buf != nil {
			text := buf.Document().Text
			if len(text) == 0 && m.cc.Default != "" {
				text = m.cc.Default
			}
			//m.ExecTask(func() {
			m.cc.OnFinish(text, nil)
			//})
		}
		debug.Println("recv input finish")
		return true
	}
	if status == CancelStatus {
		if m.cc.OnFinish != nil {
			// m.ExecTask(func() {
			m.cc.OnFinish("", io.EOF)
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
func (m *InputBlockManager) PreCheckCallBack(status int, buf *buffer.Buffer) (success bool) {
	success = true
	switch status {
	case CancelStatus:
		m.Validate.Text = ""
	case FinishStatus, NormalStatus:
		// 未输入参数,但是设置了默认值,不检测输入
		if buf.Document().Text == "" && m.cc.Default != "" {
			m.Validate.Text = ""
			break
		}
		// 检查输入
		if m.cc.Valid != nil && buf != nil {
			if err := m.cc.Valid(buf.Document()); err != nil {
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

		m.PreWords.Words = m.cc.ResultText(m.cc, status, buf.Document(), m.cc.Default)
	}

	return
}

func (m *InputBlockManager) BeforeEvent(ctx PressContext, key input.Key, in []byte) (exit bool) {
	// first deal input char event
	if key == input.NotDefined && ctx.GetBuffer() != nil {
		m.useDefault = false
		ctx.GetBuffer().InsertText(string(in), false, true)
	}
	return
}
