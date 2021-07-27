package promptx

import (
	"fmt"
	"strings"

	buffer "github.com/aggronmagi/promptx/buffer"
	"github.com/aggronmagi/promptx/history"
	"github.com/aggronmagi/promptx/internal/debug"
)

// CommonOptionsOptionDeclareWithDefault promptx options
// generate by https://github.com/timestee/optiongen
//go:generate optionGen --option_with_struct_name=true --v=true
func CommonOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TipText":         "",
		"TipTextColor":    Color(Yellow),
		"TipBGColor":      Color(DefaultColor),
		"PrefixText":      ">>> ",
		"PrefixTextColor": Color(Green),
		"PrefixBGColor":   Color(DefaultColor),
		// check input valid
		"ValidFunc":      (func(status int, in *Document) error)(nil),
		"ValidTextColor": Color(Red),
		"ValidBGColor":   Color(DefaultColor),
		// exec input command
		"ExecFunc":   (func(ctx Context, command string))(nil),
		"FinishKey":  Key(Enter),
		"CancelKey":  Key(ControlC),
		"Completion": []CompleteOption(nil),
		// if command slice size > 0. it will ignore ExecFunc and ValidFunc options
		"Commands": []*Cmd(nil),
		// alway check input command
		"AlwaysCheckCommand": bool(false),
		// history file
		"History": string(""),
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
	root       *Cmd
	history    *history.History
	hf         string
}

// NewDefaultBlockManger default blocks manager.
func NewDefaultBlockManger(opts ...CommonOption) (m *CommonBlockManager) {
	cc := NewCommonOptions(opts...)
	cc.TipText = deleteBreakLineCharacters(cc.TipText)
	m = &CommonBlockManager{
		BlocksBaseManager: &BlocksBaseManager{},
		Tip:               &BlocksWords{},
		PreWords:          &BlocksWords{},
		Input:             &BlocksEmacsBuffer{},
		Validate:          &BlocksNewLine{},
		Completion:        &BlocksCompletion{},
		cc:                cc,
		history:           history.NewHistory(),
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
			m.Completion.Update(ctx.GetBuffer().Document())
		}
		return false
	}, ControlP, Up)
	m.Completion.BindKey(func(ctx PressContext) (exit bool) {
		buf := ctx.GetBuffer()
		if new, ok := m.history.Newer(buf.Text()); ok {
			buf.Reset()
			buf.InsertText(new, false, true)
			m.Completion.Update(ctx.GetBuffer().Document())
		}
		return false
	}, ControlN, Down)

	m.SetBeforeEvent(m.BeforeEvent)
	m.SetBehindEvent(m.BehindEvent)
	m.SetCancelKeyAutoExit(false)

	m.applyOptionModify()
	return
}

func (m *CommonBlockManager) applyOptionModify() {
	cc := m.cc

	if m.hf != cc.History {
		if len(m.hf) == 0 {
			debug.AssertNoError(m.history.Load(cc.History))
		} else {
			debug.AssertNoError(m.history.Save(m.hf))
			m.history.Reset()
		}
		m.hf = cc.History
	}

	if len(cc.TipText) > 0 {
		m.Tip.Words = append(m.Tip.Words, &Word{
			Text:      cc.TipText,
			TextColor: cc.TipTextColor,
			BGColor:   cc.TipBGColor,
			Bold:      false,
		})
		m.Tip.Words = append(m.Tip.Words, &NewLineWord)
	}

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
	// completion
	if m.Completion.Cfg == nil {
		m.Completion.Cfg = NewCompleteOptions(cc.Completion...)
		m.Completion.ApplyOptions()
		if m.Completion.Completions != nil {
			m.Completion.Update(buffer.NewDocument())
		}
	} else {
		m.Completion.ApplyOptions(cc.Completion...)
	}
	// command
	m.initCommand()
}

func (m *CommonBlockManager) initCommand() {
	cc := m.cc
	if len(cc.Commands) < 1 {
		// revert command
		m.root = nil
		return
	}
	m.root = &Cmd{}
	m.root.AddCmd(cc.Commands...)
	// replace completion
	m.Completion.ApplyOptions(
		WithCompleteOptionCompleter(m.completeCommand),
		WithCompleteOptionCompletionFillSpace(true),
	)
	// replace valid func
	m.cc.ValidFunc = m.validCommand

	// replace run action
	m.cc.ExecFunc = m.execCommand
}

func (m *CommonBlockManager) completeCommand(in Document) []*Suggest {
	return m.root.FindSuggest(&in)
}

func (m *CommonBlockManager) validCommand(status int, d *Document) error {
	// is check normal status
	if status == NormalStatus && !m.cc.AlwaysCheckCommand {
		return nil
	}
	if len(d.Text) == 0 {
		return nil
	}
	cmds, _, err := m.root.ParseInput(d.Text, false)
	if err != nil {
		return err
	}
	if len(cmds) < 1 {
		return fmt.Errorf("not found command[%s]", d.Text)
	}
	return err
}
func (m *CommonBlockManager) execCommand(oldCtx Context, command string) {
	if len(command) == 0 {
		return
	}
	ctx := &CommandContext{}
	ctx.Context = m.GetContext()
	ctx.Line = command
	ctx.Cmds, ctx.Args, _ = m.root.ParseInput(command, true)
	ctx.Root = m.root
	// debug.Println("find cmd size:", len(ctx.Cmds))
	// find last command which set exec func.
	find := false
	for i := len(ctx.Cmds) - 1; i >= 0; i-- {
		cmd := ctx.Cmds[i]
		// debug.Println("find ", cmd.Name)
		if cmd.Func != nil {
			// exec command func
			cmd.Func(ctx)
			find = true
			break
		}
	}
	if !find {
		oldCtx.Println("command set deal functions.", command)
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
			TextColor: m.cc.PrefixTextColor,
			BGColor:   m.cc.PrefixBGColor,
			Bold:      false,
		},
	}
	m.PreWords.test = nil

	return
}

// SetPromptWords update prompt string. custom display.
func (m *CommonBlockManager) SetPromptWords(words ...*Word) {
	if len(words) < 1 {
		return
	}
	m.PreWords.Words = words
	debug.Println("update prompts words", words)
	m.PreWords.test = func() {
		debug.Println("get prompts words", words)
	}
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

func (m *CommonBlockManager) BeforeEvent(ctx PressContext, key Key, in []byte) (exit bool) {
	// first deal input char event
	if key == NotDefined && ctx.GetBuffer() != nil {
		ctx.GetBuffer().InsertText(string(in), false, true)
	}

	return
}

func (m *CommonBlockManager) BehindEvent(ctx PressContext, key Key, in []byte) (exit bool) {

	if ctx.GetBuffer() != nil {
		if m.Input.IsBind(key) || key == NotDefined{
			m.history.Rebuild(ctx.GetBuffer().Text(), false)
		}
		if key == m.cc.CancelKey ||
			key == m.cc.FinishKey {
			m.history.Rebuild("", true)
		}
	}
	return
}

// FinishCallBack  call back
func (m *CommonBlockManager) FinishCallBack(status int, buf *Buffer) bool {
	if status == FinishStatus {
		if m.cc.ExecFunc != nil && buf != nil && buf.Text() != "" {
			text := buf.Document().Text
			ctx := m.GetContext()
			m.history.Add(buf.Text())
			m.cc.ExecFunc(ctx, text)
		}
	}
	return false
}

// PreCheckCallBack change status pre check
func (m *CommonBlockManager) PreCheckCallBack(status int, buf *Buffer) (success bool) {
	success = true
	if buf != nil && m.Completion.Active() && m.Completion.Completions != nil {
		if status == FinishStatus {
			// interrupt enter key press
			if m.Completion.EnterSelect(buf) {
				success = false
			}
		}
		if status == CancelStatus {
			m.Completion.Update(buffer.NewDocument())
		}
	}

	// check input
	if m.cc.ValidFunc != nil && buf != nil {
		switch status {
		case CancelStatus:
			m.Validate.Text = ""
		case FinishStatus, NormalStatus:
			if err := m.cc.ValidFunc(status, buf.Document()); err != nil {
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
}
