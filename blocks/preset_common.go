package blocks

import (
	"fmt"
	"strings"

	buffer "github.com/aggronmagi/promptx/v2/buffer"
	"github.com/aggronmagi/promptx/v2/history"
	"github.com/aggronmagi/promptx/v2/input"
	"github.com/aggronmagi/promptx/v2/internal/debug"
	"github.com/aggronmagi/promptx/v2/output"
)

// CommonOptions promptx options
// generate by https://github.com/aggronmagi/gogen/
//
//go:generate gogen option -n CommonOption -f -o gen_options_common.go
func promptxCommonOptions() interface{} {
	return map[string]interface{}{
		"Tip":         "",
		"TipColor":    output.Color(output.Yellow),
		"TipBG":       output.Color(output.DefaultColor),
		"Prefix":      ">>> ",
		"PrefixColor": output.Color(output.Green),
		"PrefixBG":    output.Color(output.DefaultColor),
		// check input valid
		"Valid":      (func(status int, in *buffer.Document) error)(nil),
		"ValidColor": output.Color(output.Red),
		"ValidBG":    output.Color(output.DefaultColor),
		// exec input command
		"Exec":     (func(ctx Context, command string))(nil),
		"Finish":   input.Key(input.Enter),
		"Cancel":   input.Key(input.ControlC),
		"Complete": []CompleteOption(nil),
		// history file
		"History": string(""),
		// maximum history size
		"HistoryMaxSize": int(10000),
		// ignore consecutive duplicates
		"HistoryIgnoreDups": bool(true),
		// global deduplication
		"HistoryDedup": bool(false),
		// record timestamps
		"HistoryTimestamp": bool(false),
	}
}

// CommonBlockManager default block manager.
type CommonBlockManager struct {
	*BlocksBaseManager
	Tip        *BlocksWords
	PreWords   *BlocksWords
	Input      *BlocksEmacsBuffer
	Validate   *BlocksNewLine
	Completion *BlocksCompletion
	cc         *CommonOptions
	history    *history.History
	hf         string
}

// NewDefaultBlockManger default blocks manager.
func NewDefaultBlockManger(opts ...CommonOption) (m *CommonBlockManager) {
	cc := NewCommonOptions(opts...)
	cc.Tip = deleteBreakLineCharacters(cc.Tip)
	m = &CommonBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		Tip:               &BlocksWords{},
		PreWords:          &BlocksWords{},
		Input:             &BlocksEmacsBuffer{},
		Validate:          &BlocksNewLine{},
		Completion:        &BlocksCompletion{},
		cc:                cc,
		history: history.NewHistory(
			history.WithMaxSize(cc.HistoryMaxSize),
			history.WithIgnoreDups(cc.HistoryIgnoreDups),
			history.WithDeduplicate(cc.HistoryDedup),
			history.WithTimestamp(cc.HistoryTimestamp),
		),
	}

	m.AddMirrorMode(m.Tip)
	m.AddMirrorMode(m.PreWords)
	m.AddMirrorMode(m.Input)
	m.AddMirrorMode(m.Validate)
	m.AddMirrorMode(m.Completion)

	m.SetCallBack(m.FinishCallBack)
	m.SetPreCheck(m.PreCheckCallBack)

	m.Tip.SetIsDraw(func(status int) (draw bool) {
		return status == NormalStatus
	})
	m.Completion.BindKey(func(ctx PressContext) (exit bool) {
		buf := ctx.GetBuffer()
		if new, ok := m.history.Older(buf.Text()); ok {
			buf.Reset()
			buf.InsertText(new, false, true)
			//m.Completion.Update(ctx.GetBuffer().Document())
			m.Completion.resetCompletion(ctx)
		}
		return false
	}, input.ControlP, input.Up)
	m.Completion.BindKey(func(ctx PressContext) (exit bool) {
		buf := ctx.GetBuffer()
		if new, ok := m.history.Newer(buf.Text()); ok {
			buf.Reset()
			buf.InsertText(new, false, true)
			//m.Completion.Update(ctx.GetBuffer().Document())
			m.Completion.resetCompletion(ctx)
		}
		return false
	}, input.ControlN, input.Down)

	m.SetBeforeEvent(m.BeforeEvent)
	m.SetBehindEvent(m.BehindEvent)
	m.SetCancelKeyAutoExit(false)

	m.applyOptionModify()
	return
}

func (m *CommonBlockManager) applyOptionModify() {
	cc := m.cc

	m.history.ApplyOptions(
		history.WithMaxSize(cc.HistoryMaxSize),
		history.WithIgnoreDups(cc.HistoryIgnoreDups),
		history.WithDeduplicate(cc.HistoryDedup),
		history.WithTimestamp(cc.HistoryTimestamp),
	)

	if m.hf != cc.History {
		if len(m.hf) == 0 {
			debug.AssertNoError(m.history.Load(cc.History))
		} else {
			debug.AssertNoError(m.history.Save(m.hf))
			m.history.Reset()
		}
		m.hf = cc.History
	}

	if len(cc.Tip) > 0 {
		m.Tip.Words = append(m.Tip.Words, &Word{
			Text:      cc.Tip,
			TextColor: cc.TipColor,
			BGColor:   cc.TipBG,
			Bold:      false,
		})
		m.Tip.Words = append(m.Tip.Words, NewLineWord)
	}

	if len(m.PreWords.Words) == 0 {
		m.PreWords.Words = append(m.PreWords.Words, &Word{
			Text:      cc.Prefix,
			TextColor: cc.PrefixColor,
			BGColor:   cc.PrefixBG,
			Bold:      false,
		})
	}

	m.Validate.TextColor = cc.ValidColor
	m.Validate.BGColor = cc.ValidBG

	m.SetCancelKey(cc.Cancel)
	m.SetFinishKey(cc.Finish)
	// completion
	if m.Completion.Cfg == nil {
		m.Completion.Cfg = NewCompleteOptions(cc.Complete...)
		m.Completion.ApplyOptions()
		if m.Completion.Completions != nil {
			m.Completion.Update(buffer.NewDocument())
		}
	} else {
		m.Completion.ApplyOptions(cc.Complete...)
	}
}

// SetPrompt set prefix text
func (m *CommonBlockManager) SetPrompt(text string) {
	if !strings.HasSuffix(text, " ") {
		text += " "
	}
	if len(m.PreWords.Words) == 1 {
		m.PreWords.Words[0].Text = text
		return
	}
	m.PreWords.Words = []*Word{
		{
			Text:      text,
			TextColor: m.cc.PrefixColor,
			BGColor:   m.cc.PrefixBG,
			Bold:      false,
		},
	}
	m.PreWords.test = nil
}

// SetPromptWords update prompt string. custom display.
func (m *CommonBlockManager) SetPromptWords(words ...*Word) {
	if len(words) < 1 {
		return
	}
	// 自动追加空格
	last := words[len(words)-1]
	if !strings.HasSuffix(last.Text, " ") {
		last.Text += " "
	}
	m.PreWords.Words = words
	debug.Println("update prompts words", words)
	m.PreWords.test = func() {
		debug.Println("get prompts words", words)
	}
}

// RemoveHistory remove from history
func (m *CommonBlockManager) RemoveHistory(line string) {
	m.history.Remove(line)
}

// AddHistory add line to history
func (m *CommonBlockManager) AddHistory(line string) {
	m.history.Add(line)
}

func (m *CommonBlockManager) ResetHistoryFile(filename string) {
	cc := m.cc
	if filename == cc.History {
		m.history.Reset()
		return
	}
	if len(m.hf) > 0 {
		debug.AssertNoError(m.history.Save(m.hf))
	}

	m.history.Reset()

	cc.History = filename
	if len(filename) < 1 {
		return
	}
	debug.AssertNoError(m.history.Load(cc.History))
	m.hf = cc.History
}

func (m *CommonBlockManager) SetOption(opt CommonOption) {
	_ = opt(m.cc)
	m.applyOptionModify()
}

func (m *CommonBlockManager) ApplyOption(opts ...CommonOption) {
	for _, opt := range opts {
		_ = opt(m.cc)
	}
	m.applyOptionModify()
}

func (m *CommonBlockManager) BeforeEvent(ctx PressContext, key input.Key, in []byte) (exit bool) {
	// first deal input char event
	if key == input.NotDefined && ctx.GetBuffer() != nil {
		ctx.GetBuffer().InsertText(string(in), false, true)
	}

	return
}

func (m *CommonBlockManager) BehindEvent(ctx PressContext, key input.Key, in []byte) (exit bool) {

	if ctx.GetBuffer() != nil {
		if m.Input.IsBind(key) || key == input.NotDefined {
			m.history.Rebuild(ctx.GetBuffer().Text(), false)
		}
		if key == m.cc.Cancel ||
			key == m.cc.Finish {
			m.history.Rebuild("", true)
		}
		// when exit,reset completion.
		if key == input.ControlD && len(ctx.GetBuffer().Text()) == 0 && m.Completion != nil && m.Completion.Completions != nil {
			m.Completion.Completions.Reset()
		}
	}
	return
}

// FinishCallBack  call back
func (m *CommonBlockManager) FinishCallBack(status int, buf *buffer.Buffer) bool {
	if status == FinishStatus {
		if m.cc.Exec != nil && buf != nil && buf.Text() != "" {
			text := buf.Document().Text
			ctx := m.GetContext()
			m.history.Add(buf.Text())
			m.cc.Exec(ctx, text)
		}
	}
	return false
}

// PreCheckCallBack change status pre check
func (m *CommonBlockManager) PreCheckCallBack(status int, buf *buffer.Buffer) (success bool) {
	success = true
	if buf != nil && m.Completion.Active() && m.Completion.Completions != nil {
		if status == FinishStatus {
			// interrupt enter key press
			if m.Completion.EnterSelect(buf) {
				success = false
			}
		}
		if status == CancelStatus {
			// m.Completion.Update(buffer.NewDocument())
		}
	}

	// check input
	if m.cc.Valid != nil && buf != nil {
		switch status {
		case CancelStatus:
			m.Validate.Text = ""
		case FinishStatus, NormalStatus:
			if err := m.cc.Valid(status, buf.Document()); err != nil {
				m.Validate.Text = err.Error()
				success = false
			} else {
				m.Validate.Text = ""
			}
		}
		// if valid failed. close completion.
		if m.Completion.Active() && !success {
			m.Completion.Completions.Reset()
		}
	}
	if !success {
		return
	}

	return
}

// TearDown to clear title and erasing.
func (m *CommonBlockManager) TearDown() {
	m.BlocksBaseManager.TearDown()
	if len(m.hf) > 0 {
		debug.AssertNoError(m.history.Save(m.hf))
	}
	// Fix linux new line
	fmt.Println()
}
